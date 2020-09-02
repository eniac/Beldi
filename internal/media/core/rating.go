package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func UploadRating2(env *beldilib.Env, reqId string, rating int32) {
	beldilib.AsyncInvoke(env, TComposeReview(), RPCInput{
		Function: "UploadRating",
		Input: aws.JSONValue{
			"reqId":  reqId,
			"rating": rating,
		},
	})
}
