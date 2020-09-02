package profile

import "github.com/eniac/Beldi/internal/hotel/main/data"

type Request struct {
	HotelIds []string
	Locale   string
}

type Result struct {
	Hotels []data.Hotel
}
