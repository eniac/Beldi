package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
	"sync"
)

func UploadReq(env *beldilib.Env, reqId string) {
	beldilib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(aws.JSONValue{"reqId": reqId, "counter": 0}),
	})
}

func UploadUniqueId(env *beldilib.Env, reqId string, reviewId string) {
	beldilib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.reviewId"): expression.Value(reviewId),
		expression.Name("V.counter"):  expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadText(env *beldilib.Env, reqId string, text string) {
	beldilib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.text"):    expression.Value(text),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadRating(env *beldilib.Env, reqId string, rating int32) {
	beldilib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.rating"):  expression.Value(rating),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadUserId(env *beldilib.Env, reqId string, userId string) {
	beldilib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.userId"):  expression.Value(userId),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func UploadMovieId(env *beldilib.Env, reqId string, movieId string) {
	beldilib.Write(env, TComposeReview(), reqId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.movieId"): expression.Value(movieId),
		expression.Name("V.counter"): expression.Name("V.counter").Plus(expression.Value(1)),
	})
	TryComposeAndUpload(env, reqId)
}

func Cleanup(reqId string) {
	// Debugging
	if beldilib.TYPE == "BASELINE" {
		beldilib.LibDelete(TComposeReview(), aws.JSONValue{"K": reqId})
		return
	}
	beldilib.LibDelete(TComposeReview(), aws.JSONValue{
		"K":       reqId,
		"ROWHASH": "HEAD",
	})
	//cond := expression.Key("K").Equal(expression.Value(reqId))
	//expr, err := expression.NewBuilder().
	//	WithProjection(beldilib.BuildProjection([]string{"K", "ROWHASH"})).
	//	WithKeyCondition(cond).Build()
	//beldilib.CHECK(err)
	//res, err := beldilib.DBClient.Query(&dynamodb.QueryInput{
	//	TableName:                 aws.String(TComposeReview()),
	//	KeyConditionExpression:    expr.KeyCondition(),
	//	ProjectionExpression:      expr.Projection(),
	//	ExpressionAttributeNames:  expr.Names(),
	//	ExpressionAttributeValues: expr.Values(),
	//	ConsistentRead:            aws.Bool(true),
	//})
	//beldilib.CHECK(err)
	//var items []aws.JSONValue
	//err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &items)
	//beldilib.CHECK(err)
	//for _, item := range items {
	//	beldilib.LibDelete(TComposeReview(), item)
	//}
}

func TryComposeAndUpload(env *beldilib.Env, reqId string) {
	item := beldilib.Read(env, TComposeReview(), reqId)
	if item == nil {
		return
	}
	res := item.(map[string]interface{})
	if counter, ok := res["counter"].(float64); ok {
		if int32(counter) == 5 {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				Cleanup(reqId)
			}()
			var review Review
			beldilib.CHECK(mapstructure.Decode(res, &review))
			//beldilib.AsyncInvoke(env, TReviewStorage(), RPCInput{
			//	Function: "StoreReview",
			//	Input:    review,
			//})
			//beldilib.AsyncInvoke(env, TUserReview(), RPCInput{
			//	Function: "UploadUserReview",
			//	Input: aws.JSONValue{
			//		"userId":    review.UserId,
			//		"reviewId":  review.ReviewId,
			//		"timestamp": review.Timestamp,
			//	},
			//})
			//beldilib.AsyncInvoke(env, TMovieReview(), RPCInput{
			//	Function: "UploadMovieReview",
			//	Input: aws.JSONValue{
			//		"movieId":   review.MovieId,
			//		"reviewId":  review.ReviewId,
			//		"timestamp": review.Timestamp,
			//	},
			//})
			wg.Wait()
		}
	} else {
		panic("counter not found")
	}
}
