package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
	"sync"
)

func ReadPage(env *beldilib.Env, movieId string) Page {
	var movieInfo MovieInfo
	var reviews []Review
	var castInfos []CastInfo
	var plot string
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		res, _ := beldilib.AssignedSyncInvoke(env, TMovieInfo(), RPCInput{
			Function: "ReadMovieInfo",
			Input:    aws.JSONValue{"movieId": movieId},
		}, env.StepNumber)
		beldilib.CHECK(mapstructure.Decode(res, &movieInfo))
		var ids []string
		for _, cast := range movieInfo.Casts {
			ids = append(ids, cast.CastInfoId)
		}
		go func() {
			defer wg.Done()
			res, _ := beldilib.AssignedSyncInvoke(env, TCastInfo(), RPCInput{
				Function: "ReadCastInfo",
				Input:    ids,
			}, env.StepNumber+1)
			beldilib.CHECK(mapstructure.Decode(res, &castInfos))
		}()
		go func() {
			defer wg.Done()
			res, _ := beldilib.AssignedSyncInvoke(env, TPlot(), RPCInput{
				Function: "ReadPlot",
				Input:    aws.JSONValue{"plotId": movieInfo.PlotId},
			}, env.StepNumber+2)
			beldilib.CHECK(mapstructure.Decode(res, &plot))
		}()
	}()
	go func() {
		defer wg.Done()
		res, _ := beldilib.AssignedSyncInvoke(env, TMovieReview(), RPCInput{
			Function: "ReadMovieReviews",
			Input:    aws.JSONValue{"movieId": movieId},
		}, env.StepNumber+3)
		beldilib.CHECK(mapstructure.Decode(res, &reviews))
	}()
	wg.Wait()
	env.StepNumber += 4
	return Page{CastInfos: castInfos, Reviews: reviews, MovieInfo: movieInfo, Plot: plot}
}
