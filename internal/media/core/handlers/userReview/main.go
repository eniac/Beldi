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
	case "UploadUserReview":
		core.UploadUserReview(env, req["userId"].(string),
			req["reviewId"].(string), req["timestamp"].(string))
		return 0
	case "ReadUserReviews":
		return core.ReadUserReviews(env, req["userId"].(string))
	}
	panic("no such function")
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
