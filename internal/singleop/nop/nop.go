package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func Handler(env *beldilib.Env) interface{} {
	return 0
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
