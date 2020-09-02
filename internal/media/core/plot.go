package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func WritePlot(env *beldilib.Env, plotId string, plot string) {
	beldilib.Write(env, TPlot(), plotId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(aws.JSONValue{"plotId": plotId, "plot": plot}),
	})
}

func ReadPlot(env *beldilib.Env, plotId string) string {
	item := beldilib.Read(env, TPlot(), plotId)
	return item.(map[string]interface{})["plot"].(string)
}
