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
	switch rpcInput.Function {
	case "WriteCastInfo":
		var info core.CastInfo
		beldilib.CHECK(mapstructure.Decode(rpcInput.Input, &info))
		core.WriteCastInfo(env, info)
		return 0
	case "ReadCastInfo":
		var castInfos []string
		beldilib.CHECK(mapstructure.Decode(rpcInput.Input, &castInfos))
		return core.ReadCastInfo(env, castInfos)
	}
	panic("no such function")
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
