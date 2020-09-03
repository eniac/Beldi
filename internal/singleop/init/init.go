package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"os"
	"time"
)

func ClearAll() {
	beldilib.DeleteLambdaTables("singleop")
	beldilib.DeleteLambdaTables("nop")
	beldilib.DeleteTable("bsingleop")
	beldilib.DeleteTable("bnop")
	beldilib.DeleteLambdaTables("tsingleop")
	beldilib.DeleteLambdaTables("tnop")
}

func main() {
	if len(os.Args) >= 2 {
		option := os.Args[1]
		if option == "clean" {
			ClearAll()
			return
		}
	}
	ClearAll()
	beldilib.WaitUntilAllDeleted([]string{
		"singleop", "singleop-log", "singleop-collector",
		"nop", "nop-log", "nop-collector",
		"tsingleop", "tsingleop-log", "tsingleop-collector",
		"tnop", "tnop-log", "tnop-collector",
		"bsingleop", "bnop",
	})
	beldilib.CreateLambdaTables("singleop")
	beldilib.CreateLambdaTables("nop")

	beldilib.CreateBaselineTable("bsingleop")
	beldilib.CreateBaselineTable("bnop")

	beldilib.CreateTxnTables("tsingleop")
	beldilib.CreateTxnTables("tnop")

	time.Sleep(60 * time.Second)
	beldilib.WriteNRows("singleop", "K", 20)

	beldilib.LibWrite("bsingleop", aws.JSONValue{"K": "K"}, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(1),
	})

	beldilib.LibWrite("tsingleop", aws.JSONValue{"K": "K"}, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(1),
	})
}
