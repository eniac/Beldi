package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func UploadText2(env *beldilib.Env, reqId string, text string) {
	beldilib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadText",
		Input:    aws.JSONValue{"reqId": reqId, "text": text},
	})
}
