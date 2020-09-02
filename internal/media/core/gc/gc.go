package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func Handler(input aws.JSONValue) {
	if _, ok := input["account"]; !ok {
		service := input["service"].(string)
		static := input["static"].(bool)
		//fmt.Printf("Start GC: %s\n", service)
		if static {
			beldilib.StaticGC(service)
		} else {
			beldilib.GC(service)
		}
		return
	}
	services := []string{"ComposeReview", "UserReview", "MovieReview", "ReviewStorage"}
	statics := []string{"Frontend", "MovieId", "UniqueId", "Plot", "MovieInfo", "User", "Rating", "Text"}
	for _, service := range services {
		args := aws.JSONValue{
			"service": service,
			"static":  false,
		}
		stream, err := json.Marshal(args)
		beldilib.CHECK(err)
		_, err = beldilib.LambdaClient.Invoke(&lambdaSdk.InvokeInput{
			FunctionName:   aws.String("beldi-dev-mediagc"),
			Payload:        stream,
			InvocationType: aws.String("Event"),
		})
		beldilib.CHECK(err)
	}
	for _, service := range statics {
		args := aws.JSONValue{
			"service": service,
			"static":  true,
		}
		stream, err := json.Marshal(args)
		beldilib.CHECK(err)
		_, err = beldilib.LambdaClient.Invoke(&lambdaSdk.InvokeInput{
			FunctionName:   aws.String("beldi-dev-mediagc"),
			Payload:        stream,
			InvocationType: aws.String("Event"),
		})
		beldilib.CHECK(err)
	}
}

func main() {
	lambda.Start(Handler)
}
