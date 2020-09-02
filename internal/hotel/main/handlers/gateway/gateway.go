package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func Handler(env *beldilib.Env) interface{} {
	var rpcInput data.RPCInput
	beldilib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	//req := rpcInput.Input.(map[string]interface{})
	switch rpcInput.Function {
	case "search":
		res, _ := beldilib.SyncInvoke(env, data.Tsearch(), rpcInput.Input)
		return res
	case "recommend":
		res, _ := beldilib.SyncInvoke(env, data.Trecommendation(), rpcInput.Input)
		return res
	case "user":
		res, _ := beldilib.SyncInvoke(env, data.Tuser(), rpcInput.Input)
		return res
	case "reserve":
		res, _ := beldilib.SyncInvoke(env, data.Tfrontend(), rpcInput.Input)
		return res
	}
	return 0
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
