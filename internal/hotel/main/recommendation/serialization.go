package recommendation

type Request struct {
	Require string
	Lat     float64
	Lon     float64
}

type Result struct {
	HotelIds []string
}
