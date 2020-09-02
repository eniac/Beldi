package beldilib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	lambdaSdk "github.com/aws/aws-sdk-go/service/lambda"
	"time"
)

func RestartAll(lambdaId string) {
	intentTable := fmt.Sprintf("%s-collector", lambdaId)
	filter := expression.Name("ASYNC").Equal(expression.Value(true)).
		And(expression.Name("DONE").Equal(expression.Value(false)).
			And(expression.AttributeExists(expression.Name("ST")).
				And(expression.Name("ST").LessThan(expression.Value(time.Now().Unix() - T)))))
	items := LibScan(intentTable, []string{"InstanceId", "INPUT"}, filter)
	for _, item := range items {
		instanceId := item["InstanceId"].(string)
		input := item["INPUT"]
		iw := InputWrapper{
			InstanceId: instanceId,
			Async:      true,
			CallerName: "",
			Input:      input,
		}
		payload := iw.Serialize()
		_, err := LambdaClient.Invoke(&lambdaSdk.InvokeInput{
			FunctionName:   aws.String(fmt.Sprintf("beldi-dev-%s", lambdaId)),
			Payload:        payload,
			InvocationType: aws.String("Event"),
		})
		CHECK(err)
	}
}
