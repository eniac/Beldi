package profile

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func GetProfiles(env *beldilib.Env, req Request) Result {
	var hotels []data.Hotel
	for _, i := range req.HotelIds {
		hotel := data.Hotel{}
		res := beldilib.Read(env, data.Tprofile(), i)
		err := mapstructure.Decode(res, &hotel)
		beldilib.CHECK(err)
		hotels = append(hotels, hotel)
	}
	return Result{Hotels: hotels}
}
