package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/rate"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func Handler(env *beldilib.Env) interface{} {
	req := rate.Request{}
	err := mapstructure.Decode(env.Input, &req)
	beldilib.CHECK(err)
	return rate.GetRates(env, req)
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
