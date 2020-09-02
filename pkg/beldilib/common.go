package beldilib

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/lambda"
	"strconv"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var LambdaClient = lambda.New(sess)

//var url = "http://133.130.115.39:8000"
//var DBClient = dynamodb.New(sess, &aws.Config{Endpoint: aws.String(url),
//	Region:                        aws.String("us-east-1"),
//	CredentialsChainVerboseErrors: aws.Bool(true)})

var DBClient = dynamodb.New(sess)

var DLOGSIZE = "1000"

func GLOGSIZE() int {
	r, _ := strconv.Atoi(DLOGSIZE)
	return r
}

var T = int64(60)

var TYPE = "BELDI"

func CHECK(err error) {
	if err != nil {
		panic(err)
	}
}
