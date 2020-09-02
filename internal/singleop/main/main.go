package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
	"time"
)

var TXN = "DISABLE"

func Handler(env *beldilib.Env) interface{} {
	if TXN == "ENABLE" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.SyncInvoke(env, "tnop", "")
		fmt.Printf("DURATION Call %s\n", time.Since(start))

		start = time.Now()
		beldilib.TWrite(env, "tsingleop", "K", a)
		fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.TCondWrite(env, "tsingleop", "K", a, true)
		fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.TCondWrite(env, "tsingleop", "K", a, false)
		fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.TRead(env, "tsingleop", "K")
		fmt.Printf("DURATION Read %s\n", time.Since(start))
		return 0
	}
	if beldilib.TYPE == "BELDI" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.Write(env, "singleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(a),
		})
		fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "singleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(1),
		}, expression.Name("V").Equal(expression.Value(a)))
		fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "singleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(a),
		}, expression.Name("V").Equal(expression.Value(2)))
		fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.Read(env, "singleop", "K")
		fmt.Printf("DURATION Read %s\n", time.Since(start))

		start = time.Now()
		beldilib.SyncInvoke(env, "nop", "")
		fmt.Printf("DURATION Call %s\n", time.Since(start))
		return 0
	}
	if beldilib.TYPE == "BASELINE" {
		a := shortuuid.New()
		start := time.Now()
		beldilib.Write(env, "bsingleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(a),
		})
		fmt.Printf("DURATION DWrite %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "bsingleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(1),
		}, expression.Name("V").Equal(expression.Value(a)))
		fmt.Printf("DURATION CWriteT %s\n", time.Since(start))

		start = time.Now()
		beldilib.CondWrite(env, "bsingleop", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V2"): expression.Value(a),
		}, expression.Name("V").Equal(expression.Value(2)))
		fmt.Printf("DURATION CWriteF %s\n", time.Since(start))

		start = time.Now()
		beldilib.Read(env, "bsingleop", "K")
		fmt.Printf("DURATION Read %s\n", time.Since(start))

		start = time.Now()
		beldilib.SyncInvoke(env, "bnop", "")
		fmt.Printf("DURATION Call %s\n", time.Since(start))
		return 0
	}
	return 1
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
