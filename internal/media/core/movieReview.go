package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func UploadMovieReview(env *beldilib.Env, movieId string, reviewId string, timestamp string) {
	reviewInfo := ReviewInfo{ReviewId: reviewId, Timestamp: timestamp}
	item := beldilib.Read(env, TMovieReview(), movieId)
	if item == nil {
		beldilib.Write(env, TMovieReview(), movieId, map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(aws.JSONValue{"reviews": []ReviewInfo{reviewInfo}}),
		})
	} else {
		beldilib.Write(env, TMovieReview(), movieId, map[expression.NameBuilder]expression.OperandBuilder{
			//expression.Name("V.reviews"): expression.Name("V.reviews").ListAppend(expression.Value([]ReviewInfo{reviewInfo})),
			expression.Name("V.reviews"): expression.Name("V.reviews"),
		})
	}
}

func ReadMovieReviews(env *beldilib.Env, movieId string) []Review {
	item := beldilib.Read(env, TMovieReview(), movieId)
	if item == nil {
		return []Review{}
	}
	var reviewInfos []ReviewInfo
	beldilib.CHECK(mapstructure.Decode(item.(map[string]interface{})["reviews"], &reviewInfos))
	var reviewIds []string
	for _, review := range reviewInfos {
		reviewIds = append(reviewIds, review.ReviewId)
	}
	var res []Review
	output, _ := beldilib.SyncInvoke(env, TReviewStorage(), RPCInput{
		Function: "ReadReviews",
		Input:    reviewIds,
	})
	beldilib.CHECK(mapstructure.Decode(output, &res))
	return res
}
