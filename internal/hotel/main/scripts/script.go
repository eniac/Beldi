package main

//
//import (
//	"fmt"
//	"txn/utility"
//)
//
//var services = []string{"user", "search", "flight", "frontend", "geo", "order",
//	"hotel", "profile", "rate", "recommendation", "gateway"}
//
//func CreateAll() {
//	for _, service := range services {
//		utility.CreateLambdaTables(service)
//	}
//}
//
//func DeleteAll() {
//	for _, service := range services {
//		utility.DeleteLambdaTables(service)
//	}
//}
//
//func CreateBaseline() {
//	for _, service := range services {
//		utility.CreateBaselineTable(service)
//	}
//}
//
//func DeleteBaseline() {
//	for _, service := range services {
//		utility.DeleteTable(service)
//	}
//}
//
//func CreateLocal() {
//	ss := []string{"flight", "frontend", "order", "hotel"}
//	for _, service := range ss {
//		utility.CreateMainTable(fmt.Sprintf("%s-local", service))
//	}
//}
//
//func DeleteLocal() {
//	ss := []string{"flight", "frontend", "order", "hotel"}
//	for _, service := range ss {
//		utility.DeleteTable(fmt.Sprintf("%s-local", service))
//	}
//}
//
//func main() {
//	DeleteAll()
//	//DeleteLocal()
//	//CreateAll()
//	//CreateLocal()
//	//CreateBaseline()
//	//DeleteBaseline()
//}
