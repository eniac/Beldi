package hotel

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

type Hotel struct {
	HotelId   string
	Cap       int32
	Customers []string
}

func BaseReserveHotel(env *beldilib.Env, hotelId string, userId string) bool {
	item := beldilib.Read(env, data.Thotel(), hotelId)
	var hotel Hotel
	beldilib.CHECK(mapstructure.Decode(item, &hotel))
	if hotel.Cap == 0 {
		return false
	}
	beldilib.Write(env, data.Thotel(), hotelId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V.Cap"): expression.Value(hotel.Cap),
	})
	return true
}

func ReserveHotel(env *beldilib.Env, hotelId string, userId string) bool {
	ok, item := beldilib.TPLRead(env, data.Thotel(), hotelId, []string{"V"})
	if !ok {
		return false
	}
	var hotel Hotel
	beldilib.CHECK(mapstructure.Decode(item["V"], &hotel))
	if hotel.Cap == 0 {
		return false
	}
	ok = beldilib.TPLWrite(env, data.Thotel(), hotelId,
		aws.JSONValue{"V.Cap": hotel.Cap})
	return ok
}

func AddHotel(env *beldilib.Env, hotelId string, cap int32) {
	beldilib.Write(env, data.Thotel(), hotelId, map[expression.NameBuilder]expression.OperandBuilder{
		expression.Name("V"): expression.Value(Hotel{
			HotelId:   hotelId,
			Cap:       cap,
			Customers: []string{},
		}),
	})
}
