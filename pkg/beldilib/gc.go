package beldilib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"strings"
	"sync"
	"time"
)

var DEBUG = false

func ClearRow(tablename string, key string, prevRow string, currentRow string, ts int64) {
	currentPk := aws.JSONValue{"K": key, "ROWHASH": currentRow}
	prevPk := aws.JSONValue{"K": key, "ROWHASH": prevRow}
	res := LibRead(tablename, currentPk, []string{"NEXTROW", "GCSIZE"})
	if nextRow, exists := res["NEXTROW"].(string); exists {
		if int(res["GCSIZE"].(float64)) == GLOGSIZE() {
			LibWrite(tablename, prevPk, map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("NEXTROW"): expression.Value(nextRow),
			})
			LibWrite(tablename, currentPk, map[expression.NameBuilder]expression.OperandBuilder{
				expression.Name("TS"): expression.Value(ts),
			})
			ClearRow(tablename, key, prevRow, nextRow, ts)
		} else {
			ClearRow(tablename, key, currentRow, nextRow, ts)
		}
	} else {
		// never remove the last row
		return
	}
}

func ScanIntent(lambdaId string) []aws.JSONValue {
	filter := expression.Name("DONE").Equal(expression.Value(true)).
		And(expression.AttributeExists(expression.Name("TS")).
			And(expression.Name("TS").LessThan(expression.Value(time.Now().Unix() - T))))
	intentTable := fmt.Sprintf("%s-collector", lambdaId)
	return LibScan(intentTable, []string{"InstanceId"}, filter)
}

func QueryHeads(lambdaId string) []aws.JSONValue {
	filter := expression.Name("ROWHASH").Equal(expression.Value("HEAD"))
	return LibScan(lambdaId, []string{"K", "NEXTROW"}, filter)
}

func ScanDangle(lambdaId string) []aws.JSONValue {
	filter := expression.AttributeExists(expression.Name("TS")).
		And(expression.Name("TS").LessThan(expression.Value(time.Now().Unix() - T)))
	return LibScan(lambdaId, []string{"K", "ROWHASH"}, filter)
}

func ClearReadLog(wg *sync.WaitGroup, lambdaId string, instanceId string) {
	projection := []string{"InstanceId", "StepNumber"}
	cond := expression.Key("InstanceId").Equal(expression.Value(instanceId))
	logTable := fmt.Sprintf("%s-log", lambdaId)
	items := LibQuery(logTable, cond, projection)
	for _, item := range items {
		wg.Add(1)
		go func(item_ aws.JSONValue) {
			defer wg.Done()
			LibDelete(logTable, item_)
		}(item)
	}
}

func ClearRowDAAL(row aws.JSONValue, idx map[string]bool, lambdaId string) {
	_, Key := GeneratePK(row["K"].(string), row["ROWHASH"].(string))
	update := expression.UpdateBuilder{}
	count := 0
	logs, ok := row["LOGS"].(map[string]interface{})
	if !ok {
		return
	}
	for cid, _ := range logs {
		cidPath := fmt.Sprintf("LOGS.%s", cid)
		for instanceId, _ := range idx {
			if strings.HasPrefix(cid, instanceId) {
				update = update.Remove(expression.Name(cidPath))
				count += 1
				if count >= 250 {
					update = update.
						Set(expression.Name("GCSIZE"), expression.Name("GCSIZE").Plus(expression.Value(count)))
					expr, err := expression.NewBuilder().WithUpdate(update).
						WithCondition(expression.AttributeExists(expression.Name(cidPath))).Build()
					CHECK(err)
					_, err = DBClient.UpdateItem(&dynamodb.UpdateItemInput{
						TableName:                 aws.String(lambdaId),
						Key:                       Key,
						ConditionExpression:       expr.Condition(),
						ExpressionAttributeNames:  expr.Names(),
						ExpressionAttributeValues: expr.Values(),
						UpdateExpression:          expr.Update(),
					})
					if err != nil {
						AssertConditionFailure(err)
					}
					update = expression.UpdateBuilder{}
					count = 0
				}
				break
			}
		}
	}
	if count > 0 {
		update = update.
			Set(expression.Name("GCSIZE"), expression.Name("GCSIZE").Plus(expression.Value(count)))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		CHECK(err)
		_, err = DBClient.UpdateItem(&dynamodb.UpdateItemInput{
			TableName:                 aws.String(lambdaId),
			Key:                       Key,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		})
		CHECK(err)
	}
}

func ClearDAAL(wg *sync.WaitGroup, lambdaId string, idx map[string]bool) {
	filter := expression.Value(true).Equal(expression.Value(true))
	rows := LibScan(lambdaId, []string{"K", "ROWHASH", "LOGS", "GCSIZE"}, filter)
	for _, row := range rows {
		wg.Add(1)
		go func(row_ aws.JSONValue) {
			defer wg.Done()
			ClearRowDAAL(row_, idx, lambdaId)
		}(row)
	}
}

func ClearIntent(wg *sync.WaitGroup, lambdaId string, ids []string) {
	intentTable := fmt.Sprintf("%s-collector", lambdaId)
	for _, instanceId := range ids {
		wg.Add(1)
		go func(id_ string) {
			defer wg.Done()
			LibDelete(intentTable, aws.JSONValue{"InstanceId": id_})
		}(instanceId)
	}
}

func MarkDangling(wg *sync.WaitGroup, lambdaId string) {
	heads := QueryHeads(lambdaId)
	ts := time.Now().Unix()
	for _, head := range heads {
		wg.Add(1)
		go func(head_ aws.JSONValue) {
			defer wg.Done()
			if nextRow, exists := head_["NEXTROW"].(string); exists {
				ClearRow(lambdaId, head_["K"].(string), "HEAD", nextRow, ts)
			}
		}(head)
	}
}

func ClearDangling(wg *sync.WaitGroup, lambdaId string) {
	items := ScanDangle(lambdaId)
	for _, item := range items {
		wg.Add(1)
		go func(item_ aws.JSONValue) {
			defer wg.Done()
			LibDelete(lambdaId, item_)
		}(item)
	}
}

func GC(lambdaId string) {
	var wg sync.WaitGroup

	start := time.Now()

	items := ScanIntent(lambdaId)

	var ids []string
	idx := make(map[string]bool)
	for _, item := range items {
		var tmp = item["InstanceId"].(string)
		ids = append(ids, tmp)
		idx[tmp] = true
	}

	for _, instanceId := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ClearReadLog(&wg, lambdaId, id)
		}(instanceId)
	}

	ClearDAAL(&wg, lambdaId, idx)

	wg.Wait()
	if DEBUG {
		fmt.Printf("1: %s\n", time.Since(start))
	}

	start = time.Now()
	ClearIntent(&wg, lambdaId, ids)
	wg.Wait()
	if DEBUG {
		fmt.Printf("2: %s\n", time.Since(start))
	}

	start = time.Now()

	MarkDangling(&wg, lambdaId)
	wg.Wait()
	if DEBUG {
		fmt.Printf("3: %s\n", time.Since(start))
	}

	start = time.Now()

	ClearDangling(&wg, lambdaId)
	wg.Wait()
	if DEBUG {
		fmt.Printf("4: %s\n", time.Since(start))
	}
}

func StaticGC(lambdaId string) {
	var wg sync.WaitGroup
	items := ScanIntent(lambdaId)

	var ids []string
	for _, item := range items {
		var tmp = item["InstanceId"].(string)
		ids = append(ids, tmp)
	}

	for _, instanceId := range ids {
		//fmt.Print("Recyling: ")
		//fmt.Println(instanceId)
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ClearReadLog(&wg, lambdaId, id)
		}(instanceId)
	}
	wg.Wait()

	ClearIntent(&wg, lambdaId, ids)
	wg.Wait()
}
