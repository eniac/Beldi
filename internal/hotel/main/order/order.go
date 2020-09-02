package order

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/lithammer/shortuuid"
)

type Order struct {
	OrderId  string
	FlightId string
	HotelId  string
	UserId   string
}

func PlaceOrder(env *beldilib.Env, userId string, flightId string, hotelId string) {
	orderId := shortuuid.New()
	beldilib.Write(env, data.Torder(), orderId,
		map[expression.NameBuilder]expression.OperandBuilder{expression.Name("V"): expression.Value(Order{
			OrderId: orderId, FlightId: flightId, HotelId: hotelId, UserId: userId,
		})})
}
