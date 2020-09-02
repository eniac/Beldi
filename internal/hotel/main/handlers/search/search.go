package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/internal/hotel/main/search"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func Handler(env *beldilib.Env) interface{} {
	req := search.Request{}
	err := mapstructure.Decode(env.Input, &req)
	beldilib.CHECK(err)
	return aws.JSONValue{"search": search.Nearby(env, req)}
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
