package beldilib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func CreateMainTable(lambdaId string) {
	_, _ = DBClient.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("K"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("ROWHASH"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("TS"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("K"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("ROWHASH"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			&dynamodb.GlobalSecondaryIndex{
				IndexName: aws.String("rowidx"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("ROWHASH"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					NonKeyAttributes: []*string{aws.String("K"), aws.String("NEXTROW")},
					ProjectionType:   aws.String("INCLUDE"),
				},
			},
			&dynamodb.GlobalSecondaryIndex{
				IndexName: aws.String("tsidx"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("TS"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					NonKeyAttributes: []*string{aws.String("K"), aws.String("NEXTROW")},
					ProjectionType:   aws.String("INCLUDE"),
				},
			},
		},
		TableName: aws.String(lambdaId),
	})
}

func CreateLogTable(lambdaId string) {
	_, _ = DBClient.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("InstanceId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("StepNumber"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("InstanceId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("StepNumber"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName: aws.String(fmt.Sprintf("%s-log", lambdaId)),
	})
}

func CreateCollectorTable(lambdaId string) {
	_, _ = DBClient.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("InstanceId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("TS"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("InstanceId"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName: aws.String(fmt.Sprintf("%s-collector", lambdaId)),
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			&dynamodb.GlobalSecondaryIndex{
				IndexName: aws.String("tsidx"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("TS"),
						KeyType:       aws.String("HASH"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
			},
		},
	})
}

func CreateBaselineTable(lambdaId string) {
	_, _ = DBClient.CreateTable(&dynamodb.CreateTableInput{
		BillingMode: aws.String("PAY_PER_REQUEST"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("K"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("K"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName: aws.String(lambdaId),
	})
}

func CreateLambdaTables(lambdaId string) {
	CreateMainTable(lambdaId)
	CreateLogTable(lambdaId)
	CreateCollectorTable(lambdaId)
}

func CreateTxnTables(lambdaId string) {
	CreateBaselineTable(lambdaId)
	CreateLogTable(lambdaId)
	CreateCollectorTable(lambdaId)
}

func DeleteTable(tablename string) {
	_, _ = DBClient.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(tablename)})
}

func DeleteLambdaTables(lambdaId string) {
	DeleteTable(lambdaId)
	DeleteTable(fmt.Sprintf("%s-log", lambdaId))
	DeleteTable(fmt.Sprintf("%s-collector", lambdaId))
}

func WriteHead(tablename string, key string) {
	LibWrite(tablename, aws.JSONValue{"K": key, "ROWHASH": "HEAD"},
		map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("GCSIZE"):  expression.Value(0),
			expression.Name("LOGSIZE"): expression.Value(GLOGSIZE()),
			expression.Name("NEXTROW"): expression.Value("ROW2"),
			expression.Name("LOGS"):    expression.Value(aws.JSONValue{"ignore": nil}),
		})
}

func WriteTail(tablename string, key string, row string) {
	LibWrite(tablename, aws.JSONValue{"K": key, "ROWHASH": row},
		map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("GCSIZE"):  expression.Value(0),
			expression.Name("LOGSIZE"): expression.Value(0),
			expression.Name("LOGS"):    expression.Value(aws.JSONValue{"ignore": nil}),
		})
}

func WriteNRows(tablename string, key string, n int) {
	WriteHead(tablename, key)
	i := 2
	for ; i < n; i++ {
		LibWrite(tablename, aws.JSONValue{"K": key, "ROWHASH": fmt.Sprintf("ROW%d", i)},
			map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("GCSIZE"):  expression.Value(0),
				expression.Name("LOGSIZE"): expression.Value(GLOGSIZE()),
				expression.Name("NEXTROW"): expression.Value(fmt.Sprintf("ROW%d", i+1)),
				expression.Name("LOGS"):    expression.Value(aws.JSONValue{"ignore": nil}),
			})
	}
	WriteTail(tablename, key, fmt.Sprintf("ROW%d", i))
}

func Populate(tablename string, key string, value interface{}, baseline bool) {
	if baseline {
		btable := fmt.Sprintf("b%s", tablename)
		LibWrite(btable, aws.JSONValue{"K": key}, map[expression.NameBuilder]expression.OperandBuilder{
			expression.Name("V"): expression.Value(value),
		})
	} else {
		LibWrite(tablename, aws.JSONValue{"K": key, "ROWHASH": "HEAD"},
			map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("GCSIZE"):  expression.Value(0),
				expression.Name("LOGSIZE"): expression.Value(0),
				expression.Name("LOGS"):    expression.Value(aws.JSONValue{"ignore": nil}),
				expression.Name("V"):       expression.Value(value),
			})
	}
}
