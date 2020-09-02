package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func Handler() {
	beldilib.GC("gctest")
}

func main() {
	lambda.Start(Handler)
}
