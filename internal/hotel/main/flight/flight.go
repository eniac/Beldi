package flight

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

type Flight struct {
	FlightId  string
	Cap       int32
	Customers []string
}

func BaseReserveFlight(env *beldilib.Env, flightId string, userId string) bool {
	item := beldilib.Read(env, data.Tflight(), flightId)
	var flight Flight
	beldilib.CHECK(mapstructure.Decode(item, &flight))
	if flight.Cap == 0 {
		return false
	}
	beldilib.Write(env, data.Tflight(), flightId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.cap"): expression.Value(flight.Cap),
	})
	return true
}

func ReserveFlight(env *beldilib.Env, flightId string, userId string) bool {
	ok, item := beldilib.TPLRead(env, data.Tflight(), flightId, []string{"V"})
	if !ok {
		return false
	}
	var flight Flight
	beldilib.CHECK(mapstructure.Decode(item["V"], &flight))
	if flight.Cap == 0 {
		return false
	}
	ok = beldilib.TPLWrite(env, data.Tflight(), flightId,
		aws.JSONValue{"V.Cap": flight.Cap})
	return ok
}

func AddFlight(env *beldilib.Env, flightId string, cap int32) {
	beldilib.Write(env, data.Tflight(), flightId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(Flight{
			FlightId:  flightId,
			Cap:       cap,
			Customers: []string{},
		}),
	})
}
