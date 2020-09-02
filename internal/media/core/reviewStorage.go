package core

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func StoreReview(env *beldilib.Env, review Review) {
	beldilib.Write(env, TReviewStorage(), review.ReviewId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(review),
	})
}

func ReadReviews(env *beldilib.Env, ids []string) []Review {
	var reviews []Review
	for _, id := range ids {
		var review Review
		res := beldilib.Read(env, TReviewStorage(), id)
		beldilib.CHECK(mapstructure.Decode(res, &review))
		reviews = append(reviews, review)
	}
	return reviews
}
