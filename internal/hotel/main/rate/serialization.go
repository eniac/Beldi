package rate

import "github.com/eniac/Beldi/internal/hotel/main/data"

type RatePlans []data.RatePlan

func (r RatePlans) Len() int {
	return len(r)
}

func (r RatePlans) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
	return r[i].RoomType.TotalRate > r[j].RoomType.TotalRate
}

type Request struct {
	HotelIds []string
	Indate   string
	Outdate  string
}

type Result struct {
	RatePlans []data.RatePlan
}
