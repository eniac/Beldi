package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eniac/Beldi/pkg/beldilib"
	"sync"
)

func Handler() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		beldilib.RestartAll("Frontend")
	}()
	go func() {
		defer wg.Done()
		beldilib.RestartAll("ComposeReview")
	}()
	wg.Wait()
}

func main() {
	lambda.Start(Handler)
}
