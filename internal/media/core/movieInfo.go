package core

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func WriteMovieInfo(env *beldilib.Env, info MovieInfo) {
	beldilib.Write(env, TMovieInfo(), info.MovieId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(info),
	})
}

func ReadMovieInfo(env *beldilib.Env, movieId string) MovieInfo {
	var movieInfo MovieInfo
	item := beldilib.Read(env, TMovieId(), movieId)
	beldilib.CHECK(mapstructure.Decode(item, &movieInfo))
	return movieInfo
}

func UpdateRating(env *beldilib.Env, movieId string, sumUncommittedRating int32, numUncommittedRating int32) {
	var movieInfo MovieInfo
	item := beldilib.Read(env, TMovieId(), movieId)
	beldilib.CHECK(mapstructure.Decode(item, &movieInfo))
	movieInfo.AvgRating = (movieInfo.AvgRating*float64(movieInfo.NumRating) + float64(sumUncommittedRating)) / float64(movieInfo.NumRating+numUncommittedRating)
	movieInfo.NumRating += numUncommittedRating
	beldilib.Write(env, TMovieId(), movieId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(movieInfo),
	})
}
