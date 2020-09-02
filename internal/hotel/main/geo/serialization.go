package geo

type Request struct {
	Lat float64
	Lon float64
}

type Result struct {
	HotelIds []string
}
