package core

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func WriteCastInfo(env *beldilib.Env, info CastInfo) {
	beldilib.Write(env, TCastInfo(), info.CastInfoId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(info),
	})
}

func ReadCastInfo(env *beldilib.Env, castIds []string) []CastInfo {
	var res []CastInfo
	for _, id := range castIds {
		var castInfo CastInfo
		item := beldilib.Read(env, TCastInfo(), id)
		beldilib.CHECK(mapstructure.Decode(item, &castInfo))
		res = append(res, castInfo)
	}
	return res
}
