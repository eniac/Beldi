package core

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/beldilib"
	"sync"
)

type ComposeInput struct {
	Username string
	Password string
	Title    string
	Rating   int
	Text     string
}

func Compose(env *beldilib.Env, input ComposeInput) {
	reqId := env.InstanceId
	res, _ := beldilib.SyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadReq",
		Input:    aws.JSONValue{"reqId": reqId},
	})
	if res.(float64) != 0 {
		fmt.Println(fmt.Sprintf("DEBUG: result is %s", res))
	}
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		beldilib.AssignedSyncInvoke(env, TUniqueId(), RPCInput{
			Function: "UploadUniqueId2",
			Input:    aws.JSONValue{"reqId": reqId},
		}, env.StepNumber)
	}()
	go func() {
		defer wg.Done()
		beldilib.AssignedSyncInvoke(env, TUser(), RPCInput{
			Function: "UploadUser",
			Input:    aws.JSONValue{"reqId": reqId, "username": input.Username},
		}, env.StepNumber+1)
	}()
	go func() {
		defer wg.Done()
		beldilib.AssignedSyncInvoke(env, TMovieId(), RPCInput{
			Function: "UploadMovie",
			Input:    aws.JSONValue{"reqId": reqId, "title": input.Title, "rating": input.Rating},
		}, env.StepNumber+2)
	}()
	go func() {
		defer wg.Done()
		beldilib.AssignedSyncInvoke(env, TText(), RPCInput{
			Function: "UploadText2",
			Input:    aws.JSONValue{"reqId": reqId, "text": input.Text},
		}, env.StepNumber+3)
	}()
	wg.Wait()
	env.StepNumber += 4
}
