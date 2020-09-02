package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/media/core"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func Handler(env *beldilib.Env) interface{} {
	var rpcInput core.RPCInput
	beldilib.CHECK(mapstructure.Decode(env.Input, &rpcInput))
	switch rpcInput.Function {
	case "Compose":
		var input core.ComposeInput
		beldilib.CHECK(mapstructure.Decode(rpcInput.Input, &input))
		core.Compose(env, input)
		return 0
	}
	fmt.Println("ERROR: no such function")
	panic(rpcInput)
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
