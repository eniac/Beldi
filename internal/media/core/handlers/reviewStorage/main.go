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
	case "StoreReview":
		var review core.Review
		beldilib.CHECK(mapstructure.Decode(rpcInput.Input, &review))
		core.StoreReview(env, review)
		return 0
	case "ReadReviews":
		var reviewIds []string
		beldilib.CHECK(mapstructure.Decode(rpcInput.Input, &reviewIds))
		return core.ReadReviews(env, reviewIds)
	}
	panic("no such function")
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
