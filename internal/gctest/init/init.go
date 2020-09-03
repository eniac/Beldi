package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
	"os"
)

func main() {
	if len(os.Args) >= 2 {
		option := os.Args[1]
		if option == "clean" {
			beldilib.DeleteLambdaTables("gctest")
			beldilib.WaitUntilAllDeleted([]string{"gctest", "gctest-log", "gctest-collector"})
			return
		}
		if option == "txn" {
			beldilib.DeleteLambdaTables("gctest")
			beldilib.WaitUntilAllDeleted([]string{"gctest", "gctest-log", "gctest-collector"})
			for ; ; {
				beldilib.CreateTxnTables("gctest")
				if beldilib.WaitUntilAllActive([]string{"gctest", "gctest-log", "gctest-collector"}) {
					break
				}
			}
			beldilib.LibWrite("gctest", aws.JSONValue{"K": "K"},
				map[expression.NameBuilder]expression.OperandBuilder{
					expression.Name("V"): expression.Value(1),
				})
			return
		}
	}
	beldilib.DeleteLambdaTables("gctest")
	beldilib.WaitUntilAllDeleted([]string{"gctest", "gctest-log", "gctest-collector"})
	for ; ; {
		beldilib.CreateLambdaTables("gctest")
		if beldilib.WaitUntilAllActive([]string{"gctest", "gctest-log", "gctest-collector"}) {
			break
		}
	}
	beldilib.Populate("gctest", "K", 1, false)
}
