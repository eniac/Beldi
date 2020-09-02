package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"os"
	"time"
)

func main() {
	if len(os.Args) >= 2 {
		option := os.Args[1]
		if option == "clean" {
			beldilib.DeleteLambdaTables("gctest")
			return
		}
		if option == "txn" {
			beldilib.DeleteLambdaTables("gctest")
			time.Sleep(60 * time.Second)
			beldilib.CreateTxnTables("gctest")
			time.Sleep(60 * time.Second)
			beldilib.LibWrite("gctest", aws.JSONValue{"K": "K"},
				map[expression.NameBuilder]expression.OperandBuilder{
					expression.Name("V"): expression.Value(1),
				})
			return
		}
	}
	beldilib.DeleteLambdaTables("gctest")
	time.Sleep(60 * time.Second)
	beldilib.CreateLambdaTables("gctest")
	time.Sleep(60 * time.Second)
	beldilib.Populate("gctest", "K", 1, false)
}
