package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func main() {
	beldilib.CreateLambdaTables("rwtest")
	env := &beldilib.Env{
		LambdaId:    "rwtest",
		InstanceId:  "t1",
		LogTable:    "rwtest-log",
		IntentTable: "rwtest-collector",
		LocalTable:  "rwtest-local",
		StepNumber:  0,
	}
	for i := 0; i < 20; i++ {
		beldilib.Write(env, "rwtest", "K", map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(i),
		})
	}
	v := beldilib.Read(env, "rwtest", "K")
	fmt.Println(v)
	beldilib.CondWrite(env, "rwtest", "K", map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(100),
	}, expression.Name("V").Equal(expression.Value(19)))
	v = beldilib.Read(env, "rwtest", "K")
	fmt.Println(v)
	beldilib.GC("rwtest")
}
