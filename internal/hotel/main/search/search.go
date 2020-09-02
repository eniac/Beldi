package search

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/internal/hotel/main/geo"
	"github.com/eniac/Beldi/internal/hotel/main/rate"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/mitchellh/mapstructure"
)

func Nearby(env *beldilib.Env, req Request) Result {
	res, _ := beldilib.SyncInvoke(env, data.Tgeo(), geo.Request{Lat: req.Lat, Lon: req.Lon})
	var geoRes geo.Result
	beldilib.CHECK(mapstructure.Decode(res, &geoRes))
	res, _ = beldilib.SyncInvoke(env, data.Trate(), rate.Request{
		HotelIds: geoRes.HotelIds,
		Indate:   req.InDate,
		Outdate:  req.OutDate,
	})
	var rateRes rate.Result
	beldilib.CHECK(mapstructure.Decode(res, &rateRes))
	var hts []string
	for _, r := range rateRes.RatePlans {
		hts = append(hts, r.HotelId)
	}
	return Result{HotelIds: hts}
}
