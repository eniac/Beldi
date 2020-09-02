package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
)

func UploadUniqueId2(env *beldilib.Env, reqId string) {
	reviewId := shortuuid.New()
	beldilib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadUniqueId",
		Input:    aws.JSONValue{"reqId": reqId, "reviewId": reviewId},
	})
}
