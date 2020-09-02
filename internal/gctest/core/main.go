package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
)

func Handler(env *beldilib.Env) interface{} {
	a := shortuuid.New()
	if beldilib.DLOGSIZE != "101" {
		beldilib.Write(env, "gctest", "K",
			map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("V"): expression.Value(a),
			})
	} else {
		beldilib.TWrite(env, "gctest", "K", a)
	}
	return 0
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
