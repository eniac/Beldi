package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/media/core"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func Handler(env *beldilib.Env) interface{} {
	var rpcInput core.RPCInput
	beldilib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "UploadText2":
		core.UploadText2(env, req["reqId"].(string), req["text"].(string))
		return 0
	}
	panic("no such function")
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}