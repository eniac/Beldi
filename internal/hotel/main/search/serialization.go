package search

type Request struct {
	Lat     float64
	Lon     float64
	InDate  string
	OutDate string
}

type Result struct {
	HotelIds []string
}
