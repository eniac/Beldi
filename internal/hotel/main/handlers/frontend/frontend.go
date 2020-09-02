package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/internal/hotel/main/frontend"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func Handler(env *beldilib.Env) interface{} {
	req := env.Input.(map[string]interface{})
	return frontend.SendRequest(env, req["userId"].(string), req["flightId"].(string), req["hotelId"].(string))
}

func main() {
	lambda.Start(beldilib.Wrapper(Handler))
}
