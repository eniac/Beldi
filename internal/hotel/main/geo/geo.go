package geo

import (
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/pkg/beldilib"
	"github.com/hailocab/go-geoindex"
	"github.com/mitchellh/mapstructure"
)

func newGeoIndex(env *beldilib.Env) *geoindex.ClusteringIndex {
	var ps []data.Point
	res := beldilib.Scan(env, data.Tgeo())
	err := mapstructure.Decode(res, &ps)
	beldilib.CHECK(err)
	index := geoindex.NewClusteringIndex()
	for _, e := range ps {
		index.Add(e)
	}
	return index
}

func getNearbyPoints(env *beldilib.Env, lat float64, lon float64) []geoindex.Point {
	center := &geoindex.GeoPoint{
		Pid:  "",
		Plat: lat,
		Plon: lon,
	}
	index := newGeoIndex(env)
	res := index.KNearest(
		center,
		5,
		geoindex.Km(10), func(p geoindex.Point) bool {
			return true
		},
	)
	return res
}

func Nearby(env *beldilib.Env, req Request) Result {
	var (
		points = getNearbyPoints(env, req.Lat, req.Lon)
	)
	res := Result{HotelIds: []string{}}
	for _, p := range points {
		res.HotelIds = append(res.HotelIds, p.Id())
	}
	return res
}
