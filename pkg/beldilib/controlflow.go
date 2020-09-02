package beldilib

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/lithammer/shortuuid"
	"github.com/mitchellh/mapstructure"
	"strings"
	"time"
)

type InputWrapper struct {
	CallerName  string      `mapstructure:"CallerName"`
	CallerId    string      `mapstructure:"CallerId"`
	CallerStep  int32       `mapstructure:"CallerStep"`
	InstanceId  string      `mapstructure:"InstanceId"`
	Input       interface{} `mapstructure:"Input"`
	TxnId       string      `mapstructure:"TxnId"`
	Instruction string      `mapstructure:"Instruction"`
	Async       bool        `mapstructure:"Async"`
}

func (iw *InputWrapper) Serialize() []byte {
	stream, err := json.Marshal(*iw)
	CHECK(err)
	return stream
}

func (iw *InputWrapper) Deserialize(stream []byte) {
	err := json.Unmarshal(stream, iw)
	CHECK(err)
}

type StackTraceCall struct {
	Label string `json:"label"`
	Line  int    `json:"line"`
	Path  string `json:"path"`
}

func (ie *InvokeError) Deserialize(stream []byte) {
	err := json.Unmarshal(stream, ie)
	CHECK(err)
	if ie.ErrorMessage == "" {
		panic(errors.New("never happen"))
	}
}

type InvokeError struct {
	ErrorMessage string           `json:"errorMessage"`
	ErrorType    string           `json:"errorType"`
	StackTrace   []StackTraceCall `json:"stackTrace"`
}

type OutputWrapper struct {
	Status string
	Output interface{}
}

func (ow *OutputWrapper) Serialize() []byte {
	stream, err := json.Marshal(*ow)
	CHECK(err)
	return stream
}

func (ow *OutputWrapper) Deserialize(stream []byte) {
	err := json.Unmarshal(stream, ow)
	CHECK(err)
	if ow.Status != "Success" && ow.Status != "Failure" {
		ie := InvokeError{}
		ie.Deserialize(stream)
		panic(ie)
	}
}

func ParseInput(raw interface{}) *InputWrapper {
	var iw InputWrapper
	if body, ok := raw.(map[string]interface{})["body"]; ok {
		CHECK(json.Unmarshal([]byte(body.(string)), &iw))
	} else {
		CHECK(mapstructure.Decode(raw, &iw))
	}
	return &iw
}

func PrepareEnv(iw *InputWrapper) *Env {
	s := strings.Split(lambdacontext.FunctionName, "-")
	lambdaId := s[len(s)-1]
	if iw.InstanceId == "" {
		iw.InstanceId = shortuuid.New()
	}
	return &Env{
		LambdaId:    lambdaId,
		InstanceId:  iw.InstanceId,
		LogTable:    fmt.Sprintf("%s-log", lambdaId),
		IntentTable: fmt.Sprintf("%s-collector", lambdaId),
		LocalTable:  fmt.Sprintf("%s-local", lambdaId),
		StepNumber:  0,
		Input:       iw.Input,
		TxnId:       iw.TxnId,
		Instruction: iw.Instruction,
	}
}

func SyncInvoke(env *Env, callee string, input interface{}) (interface{}, string) {
	if TYPE == "BASELINE" {
		iw := InputWrapper{
			InstanceId:  "",
			Input:       input,
			CallerName:  "",
			Async:       false,
			TxnId:       env.TxnId,
			Instruction: env.Instruction,
		}
		if iw.Instruction == "EXECUTE" {
			LibWrite(env.LocalTable, aws.JSONValue{"K": env.TxnId}, map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("CALLEES"): expression.Name("CALLEES").ListAppend(expression.Value([]string{callee})),
			})
		}
		payload := iw.Serialize()
		res, err := LambdaClient.Invoke(&lambdaSdk.InvokeInput{
			FunctionName: aws.String(fmt.Sprintf("beldi-dev-%s", callee)),
			Payload:      payload,
		})
		CHECK(err)
		ow := OutputWrapper{}
		ow.Deserialize(res.Payload)
		switch ow.Status {
		case "Success":
			return ow.Output, iw.InstanceId
		default:
			panic("never happens")
		}
	}
	iw := InputWrapper{
		CallerName:  env.LambdaId,
		CallerId:    env.InstanceId,
		CallerStep:  env.StepNumber,
		Async:       false,
		InstanceId:  shortuuid.New(),
		Input:       input,
		TxnId:       env.TxnId,
		Instruction: env.Instruction,
	}
	pk := aws.JSONValue{"InstanceId": env.InstanceId, "StepNumber": env.StepNumber}
	ok := LibPut(env.LogTable, pk, aws.JSONValue{"Callee": iw.InstanceId})
	if !ok {
		item := LibRead(env.LogTable, pk, []string{"Callee", "RET"})
		if val, exist := item["Callee"].(string); exist {
			iw.InstanceId = val
		} else {
			panic("error")
		}
		if val, exist := item["RET"]; exist {
			env.StepNumber += 1
			return val, iw.InstanceId
		}
	}
	env.StepNumber += 1
	if iw.Instruction == "EXECUTE" {
		EOSWrite(env, env.LocalTable, env.TxnId, map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("CALLEES"): expression.Name("CALLEES").ListAppend(expression.Value([]string{callee})),
		})
	}
	payload := iw.Serialize()
	res, err := LambdaClient.Invoke(&lambdaSdk.InvokeInput{
		FunctionName: aws.String(fmt.Sprintf("beldi-dev-%s", callee)),
		Payload:      payload,
	})
	CHECK(err)
	ow := OutputWrapper{}
	ow.Deserialize(res.Payload)
	switch ow.Status {
	case "Success":
		return ow.Output, iw.InstanceId
	default:
		panic("never happens")
	}
}

func AssignedSyncInvoke(env *Env, callee string, input interface{}, stepNumber int32) (interface{}, string) {
	if TYPE == "BASELINE" {
		return SyncInvoke(env, callee, input)
	}
	iw := InputWrapper{
		CallerName:  env.LambdaId,
		CallerId:    env.InstanceId,
		CallerStep:  stepNumber,
		Async:       false,
		InstanceId:  shortuuid.New(),
		Input:       input,
		TxnId:       env.TxnId,
		Instruction: env.Instruction,
	}
	pk := aws.JSONValue{"InstanceId": env.InstanceId, "StepNumber": stepNumber}
	ok := LibPut(env.LogTable, pk, aws.JSONValue{"Callee": iw.InstanceId})
	if !ok {
		item := LibRead(env.LogTable, pk, []string{"Callee", "RET"})
		if val, exist := item["Callee"].(string); exist {
			iw.InstanceId = val
		} else {
			panic("error")
		}
		if val, exist := item["RET"]; exist {
			return val, iw.InstanceId
		}
	}
	if iw.Instruction == "EXECUTE" {
		EOSWrite(env, env.LocalTable, env.TxnId, map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("CALLEES"): expression.Name("CALLEES").ListAppend(expression.Value([]string{callee})),
		})
	}
	payload := iw.Serialize()
	res, err := LambdaClient.Invoke(&lambdaSdk.InvokeInput{
		FunctionName: aws.String(fmt.Sprintf("beldi-dev-%s", callee)),
		Payload:      payload,
	})
	CHECK(err)
	ow := OutputWrapper{}
	ow.Deserialize(res.Payload)
	switch ow.Status {
	case "Success":
		return ow.Output, iw.InstanceId
	default:
		panic("never happens")
	}
}

func AsyncInvoke(env *Env, callee string, input interface{}) string {
	if TYPE == "BASELINE" {
		iw := InputWrapper{
			InstanceId: "",
			Async:      true,
			CallerName: "",
			Input:      input,
		}
		payload := iw.Serialize()
		_, err := LambdaClient.Invoke(&lambdaSdk.InvokeInput{
			FunctionName:   aws.String(fmt.Sprintf("beldi-dev-%s", callee)),
			Payload:        payload,
			InvocationType: aws.String("Event"),
		})
		CHECK(err)
		return ""
	}

	iw := InputWrapper{
		CallerName: env.LambdaId,
		CallerId:   env.InstanceId,
		CallerStep: env.StepNumber,
		Async:      true,
		InstanceId: shortuuid.New(),
		Input:      input,
	}

	pk := aws.JSONValue{"InstanceId": env.InstanceId, "StepNumber": env.StepNumber}
	ok := LibPut(env.LogTable, pk, aws.JSONValue{"Callee": iw.InstanceId})
	if !ok {
		item := LibRead(env.LogTable, pk, []string{"Callee", "RET"})
		if val, exist := item["Callee"].(string); exist {
			iw.InstanceId = val
		} else {
			panic("error")
		}
		if _, exist := item["RET"]; exist {
			env.StepNumber += 1
			return iw.InstanceId
		}
	}

	ok = LibPut(fmt.Sprintf("%s-collector", callee), aws.JSONValue{"InstanceId": iw.InstanceId},
		aws.JSONValue{"DONE": false, "ASYNC": true, "INPUT": iw.Input, "ST": time.Now().Unix()})

	if !ok {
		env.StepNumber += 1
		return iw.InstanceId
	}

	LibWrite(env.LogTable, pk, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("RET"): expression.Value(1),
	})

	payload := iw.Serialize()
	_, err := LambdaClient.Invoke(&lambdaSdk.InvokeInput{
		FunctionName:   aws.String(fmt.Sprintf("beldi-dev-%s", callee)),
		Payload:        payload,
		InvocationType: aws.String("Event"),
	})
	CHECK(err)
	env.StepNumber += 1
	return iw.InstanceId
}

func TPLCommit(env *Env) {
	item := EOSRead(env, env.LocalTable, env.TxnId, []string{})
	var callees []string
	for k, v := range item {
		if k == "CALLEES" {
			CHECK(mapstructure.Decode(v, &callees))
			continue
		}
		ks := strings.Split(k, "-")
		if len(ks) != 2 {
			continue
		}
		tablename, key := ks[0], ks[1]
		update := map[expression.NameBuilder]expression.OperandBuilder{}
		for kk, vv := range v.(map[string]interface{}) {
			update[expression.Name(kk)] = expression.Value(vv)
		}
		update[expression.Name("HOLDER")] = expression.Value(AVAILABLE)
		EOSWrite(env, tablename, key, update)
	}
	LibDelete(env.LocalTable, aws.JSONValue{"K": env.TxnId, "ROWHASH": "HEAD"})
	for _, callee := range callees {
		if callee == " " {
			continue
		}
		SyncInvoke(env, callee, aws.JSONValue{})
	}
}

func TPLAbort(env *Env) {
	item := EOSRead(env, env.LocalTable, env.TxnId, []string{"CALLEES"})
	var callees []string
	for k, v := range item {
		if k == "CALLEES" {
			CHECK(mapstructure.Decode(v, &callees))
			continue
		}
		ks := strings.Split(k, "-")
		if len(ks) != 2 {
			continue
		}
		tablename, key := ks[0], ks[1]
		update := map[expression.NameBuilder]expression.OperandBuilder{}
		update[expression.Name("HOLDER")] = expression.Value(AVAILABLE)
		EOSWrite(env, tablename, key, update)
	}
	LibDelete(env.LocalTable, aws.JSONValue{"K": env.TxnId, "ROWHASH": "HEAD"})
	for _, callee := range callees {
		if callee == " " {
			continue
		}
		SyncInvoke(env, callee, aws.JSONValue{})
	}
}
func Wrapper(f func(env *Env) interface{}) func(iw interface{}) (OutputWrapper, error) {
	return func(raw interface{}) (OutputWrapper, error) {
		iw := ParseInput(raw)
		env := PrepareEnv(iw)
		if TYPE != "BASELINE" {
			if iw.Async == false || iw.CallerName == "" {
				LibPut(env.IntentTable, aws.JSONValue{"InstanceId": env.InstanceId},
					aws.JSONValue{"DONE": false, "ASYNC": iw.Async, "INPUT": iw.Input, "ST": time.Now().Unix()})
			} else {
				LibWrite(env.IntentTable, aws.JSONValue{"InstanceId": env.InstanceId},
					map[expression.NameBuilder]expression.OperandBuilder{
						expression.Name("ST"): expression.Value(time.Now().Unix()),
					})
			}
			//ok := LibPut(env.IntentTable, aws.JSONValue{"InstanceId": env.InstanceId},
			//	aws.JSONValue{"DONE": false, "ASYNC": iw.Async})
			//if !ok {
			//	res := LibRead(env.IntentTable, aws.JSONValue{"InstanceId": env.InstanceId}, []string{"RET"})
			//	output, exist := res["RET"]
			//	if exist {
			//		return OutputWrapper{
			//			Status: "Success",
			//			Output: output,
			//		}, nil
			//	}
			//}
		}

		var output interface{}
		if env.Instruction == "COMMIT" {
			TPLCommit(env)
			output = 0
		} else if env.Instruction == "ABORT" {
			TPLAbort(env)
			output = 0
		} else if env.Instruction == "EXECUTE" {
			EOSWrite(env, env.LocalTable, env.TxnId, map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("CALLEES"): expression.Value([]string{" "}),
			})
			output = f(env)
		} else {
			output = f(env)
		}

		if TYPE != "BASELINE" {
			if iw.CallerName != "" {
				LibWrite(fmt.Sprintf("%s-log", iw.CallerName),
					aws.JSONValue{"InstanceId": iw.CallerId, "StepNumber": iw.CallerStep},
					map[expression.NameBuilder]expression.OperandBuilder{
						expression.Name("RET"): expression.Value(output),
					})
			}
			LibWrite(env.IntentTable, aws.JSONValue{"InstanceId": env.InstanceId},
				map[expression.NameBuilder]expression.OperandBuilder{
					expression.Name("DONE"): expression.Value(true),
					expression.Name("TS"):   expression.Value(time.Now().Unix()),
					//expression.Name("RET"):  expression.Value(output),
				})
		}
		return OutputWrapper{
			Status: "Success",
			Output: output,
		}, nil
	}
}
