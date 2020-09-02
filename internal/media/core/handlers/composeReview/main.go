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
	case "UploadReq":
		core.UploadReq(env, req["reqId"].(string))
		return 0
	case "UploadUniqueId":
		core.UploadUniqueId(env, req["reqId"].(string), req["reviewId"].(string))
		return 0
	case "UploadText":
		core.UploadText(env, req["reqId"].(string), req["text"].(string))
		return 0
	case "UploadRating":
		core.UploadRating(env, req["reqId"].(string), int32(req["rating"].(float64)))
		return 0
	case "UploadUserId":
		core.UploadUserId(env, req["reqId"].(string), req["userId"].(string))
		return 0
	case "UploadMovieId":
		core.UploadMovieId(env, req["reqId"].(string), req["movieId"].(string))
		return 0
	}
	panic("no such function")
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
