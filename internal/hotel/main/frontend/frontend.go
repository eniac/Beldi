package frontend

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
)

func SendRequest(env *beldilib.Env, userId string, flightId string, hotelId string) string {
	if beldilib.TYPE == "BASELINE" {
		input := map[string]string{
			"hotelId": hotelId,
			"userId":  userId,
		}
		beldilib.SyncInvoke(env, data.Thotel(), data.RPCInput{
			Function: "BaseReserveHotel",
			Input:    input,
		})
		input = map[string]string{
			"flightId": flightId,
			"userId":   userId,
		}
		beldilib.SyncInvoke(env, data.Tflight(), data.RPCInput{
			Function: "BaseReserveFlight",
			Input:    input,
		})
		input = map[string]string{
			"flightId": flightId,
			"hotelId":  hotelId,
			"userId":   userId,
		}
		beldilib.AsyncInvoke(env, data.Torder(), data.RPCInput{
			Function: "PlaceOrder",
			Input:    input,
		})
		return ""
	}
	beldilib.BeginTxn(env)
	input := map[string]string{
		"hotelId": hotelId,
		"userId":  userId,
	}
	res, _ := beldilib.SyncInvoke(env, data.Thotel(), data.RPCInput{
		Function: "ReserveHotel",
		Input:    input,
	})
	if !res.(bool) {
		beldilib.AbortTxn(env)
		return "Place Order Fails"
	}
	input = map[string]string{
		"flightId": flightId,
		"userId":   userId,
	}
	res, _ = beldilib.SyncInvoke(env, data.Tflight(), data.RPCInput{
		Function: "ReserveFlight",
		Input:    input,
	})
	if !res.(bool) {
		beldilib.AbortTxn(env)
		return "Place Order Fails"
	}
	input = map[string]string{
		"flightId": flightId,
		"hotelId":  hotelId,
		"userId":   userId,
	}
	beldilib.CommitTxn(env)
	beldilib.AsyncInvoke(env, data.Torder(), data.RPCInput{
		Function: "PlaceOrder",
		Input:    input,
	})
	return "Place Order Success"
}
