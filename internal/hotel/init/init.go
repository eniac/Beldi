package main

import (
	"fmt"
	"github.com/eniac/Beldi/internal/hotel/main/data"
	"github.com/eniac/Beldi/internal/hotel/main/flight"
	"github.com/eniac/Beldi/internal/hotel/main/hotel"
	"github.com/eniac/Beldi/pkg/beldilib"
	"os"
	"strconv"
)

var services = []string{"user", "search", "flight", "frontend", "geo", "order",
	"hotel", "profile", "rate", "recommendation", "gateway"}

func tables(baseline bool) {
	if baseline {
		for _, service := range services {
			tablename := fmt.Sprintf("b%s", service)
			for ; ; {
				beldilib.CreateBaselineTable(tablename)
				if beldilib.WaitUntilActive(tablename) {
					break
				}
			}
		}
	} else {
		for _, service := range services {
			for ; ; {
				beldilib.CreateLambdaTables(service)
				if beldilib.WaitUntilAllActive([]string{service, fmt.Sprintf("%s-collector", service), fmt.Sprintf("%s-log", service)}) {
					break
				}
			}
		}
		ss := []string{"flight", "frontend", "order", "hotel"}
		for _, service := range ss {
			for ; ; {
				beldilib.CreateMainTable(fmt.Sprintf("%s-local", service))
				if beldilib.WaitUntilActive(fmt.Sprintf("%s-local", service)) {
					break
				}
			}
		}
	}
}

func delete_tables(baseline bool) {
	if baseline {
		for _, service := range services {
			beldilib.DeleteTable(fmt.Sprintf("b%s", service))
			beldilib.WaitUntilDeleted(fmt.Sprintf("b%s", service))
		}
	} else {
		for _, service := range services {
			beldilib.DeleteLambdaTables(service)
			beldilib.WaitUntilAllDeleted([]string{service, fmt.Sprintf("%s-collector", service), fmt.Sprintf("%s-log", service)})
		}
		ss := []string{"flight", "frontend", "order", "hotel"}
		for _, service := range ss {
			beldilib.DeleteTable(fmt.Sprintf("%s-local", service))
			beldilib.WaitUntilDeleted(fmt.Sprintf("%s-local", service))
		}
	}
}

func geo(baseline bool) {
	beldilib.Populate("geo", "1", data.Point{Pid: "1", Plat: 37.7867, Plon: 0}, baseline)
	beldilib.Populate("geo", "2", data.Point{Pid: "2", Plat: 37.7854, Plon: -122.4005}, baseline)
	beldilib.Populate("geo", "3", data.Point{Pid: "3", Plat: 37.7867, Plon: -122.4071}, baseline)
	beldilib.Populate("geo", "4", data.Point{Pid: "4", Plat: 37.7936, Plon: -122.3930}, baseline)
	beldilib.Populate("geo", "5", data.Point{Pid: "5", Plat: 37.7831, Plon: -122.4181}, baseline)
	beldilib.Populate("geo", "6", data.Point{Pid: "6", Plat: 37.7863, Plon: -122.4015}, baseline)
	for i := 7; i <= 80; i++ {
		hotelId := strconv.Itoa(i)
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		beldilib.Populate("geo", hotelId, data.Point{Pid: hotelId, Plat: lat, Plon: lon}, baseline)
	}
}

func profile(baseline bool) {
	beldilib.Populate("profile", "1", data.Hotel{
		Id:          "1",
		Name:        "Clift Hotel",
		PhoneNumber: "(415) 775-4700",
		Description: "A 6-minute walk from Union Square and 4 minutes from a Muni Metro station, this luxury hotel designed by Philippe Starck features an artsy furniture collection in the lobby, including work by Salvador Dali.",
		Address: data.Address{
			StreetNumber: "495",
			StreetName:   "Geary St",
			City:         "San Francisco",
			State:        "CA",
			Country:      "United States",
			PostalCode:   "94102",
			Lat:          37.7867,
			Lon:          -122.4112,
		},
	}, baseline)
	beldilib.Populate("profile", "2", data.Hotel{
		Id:          "2",
		Name:        "W San Francisco",
		PhoneNumber: "(415) 777-5300",
		Description: "Less than a block from the Yerba Buena Center for the Arts, this trendy hotel is a 12-minute walk from Union Square.",
		Address: data.Address{
			StreetNumber: "181",
			StreetName:   "3rt St",
			City:         "San Francisco",
			State:        "CA",
			Country:      "United States",
			PostalCode:   "94103",
			Lat:          37.7854,
			Lon:          -122.4005,
		},
	}, baseline)
	beldilib.Populate("profile", "3", data.Hotel{
		Id:          "3",
		Name:        "Hotel Zetta",
		PhoneNumber: "(415) 543-8555",
		Description: "A 3-minute walk from the Powell Street cable-car turnaround and BART rail station, this hip hotel 9 minutes from Union Square combines high-tech lodging with artsy touches.",
		Address: data.Address{
			StreetNumber: "55",
			StreetName:   "5th St",
			City:         "San Francisco",
			State:        "CA",
			Country:      "United States",
			PostalCode:   "94103",
			Lat:          37.7834,
			Lon:          -122.4071,
		},
	}, baseline)
	beldilib.Populate("profile", "4", data.Hotel{
		Id:          "4",
		Name:        "Hotel Vitale",
		PhoneNumber: "(415) 278-3700",
		Description: "This waterfront hotel with Bay Bridge views is 3 blocks from the Financial District and a 4-minute walk from the Ferry Building.",
		Address: data.Address{
			StreetNumber: "8",
			StreetName:   "Mission St",
			City:         "San Francisco",
			State:        "CA",
			Country:      "United States",
			PostalCode:   "94105",
			Lat:          37.7936,
			Lon:          -122.3930,
		},
	}, baseline)
	beldilib.Populate("profile", "5", data.Hotel{
		Id:          "5",
		Name:        "Phoenix Hotel",
		PhoneNumber: "(415) 776-1380",
		Description: "Located in the Tenderloin neighborhood, a 10-minute walk from a BART rail station, this retro motor lodge has hosted many rock musicians and other celebrities since the 1950s. It’s a 4-minute walk from the historic Great American Music Hall nightclub.",
		Address: data.Address{
			StreetNumber: "601",
			StreetName:   "Eddy St",
			City:         "San Francisco",
			State:        "CA",
			Country:      "United States",
			PostalCode:   "94109",
			Lat:          37.7831,
			Lon:          -122.4181,
		},
	}, baseline)
	beldilib.Populate("profile", "6", data.Hotel{
		Id:          "6",
		Name:        "St. Regis San Francisco",
		PhoneNumber: "(415) 284-4000",
		Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
		Address: data.Address{
			StreetNumber: "125",
			StreetName:   "3rd St",
			City:         "San Francisco",
			State:        "CA",
			Country:      "United States",
			PostalCode:   "94109",
			Lat:          37.7863,
			Lon:          -122.4015,
		},
	}, baseline)
	for i := 7; i < 80; i++ {
		hotelId := strconv.Itoa(i)
		phoneNum := "(415) 284-40" + hotelId
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		beldilib.Populate("profile", hotelId, data.Hotel{
			Id:          hotelId,
			Name:        "St. Regis San Francisco",
			PhoneNumber: phoneNum,
			Description: "St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			Address: data.Address{
				StreetNumber: "125",
				StreetName:   "3rd St",
				City:         "San Francisco",
				State:        "CA",
				Country:      "United States",
				PostalCode:   "94109",
				Lat:          lat,
				Lon:          lon,
			},
		}, baseline)
	}
}

func rate(baseline bool) {
	beldilib.Populate("rate", "1", data.RatePlan{
		HotelId: "1",
		Code:    "RACK",
		Indate:  "2015-04-09",
		Outdate: "2015-04-10",
		RoomType: data.RoomType{
			BookableRate:       190.00,
			Code:               "KNG",
			RoomDescription:    "King sized bed",
			TotalRate:          109.00,
			TotalRateInclusive: 123.17,
		},
	}, baseline)
	beldilib.Populate("rate", "2", data.RatePlan{
		HotelId: "2",
		Code:    "RACK",
		Indate:  "2015-04-09",
		Outdate: "2015-04-10",
		RoomType: data.RoomType{
			BookableRate:       139.00,
			Code:               "QN",
			RoomDescription:    "Queen sized bed",
			TotalRate:          139.00,
			TotalRateInclusive: 153.09,
		},
	}, baseline)
	beldilib.Populate("rate", "3", data.RatePlan{
		HotelId: "3",
		Code:    "RACK",
		Indate:  "2015-04-09",
		Outdate: "2015-04-10",
		RoomType: data.RoomType{
			BookableRate:       109.00,
			Code:               "KNG",
			RoomDescription:    "King sized bed",
			TotalRate:          109.00,
			TotalRateInclusive: 123.17,
		},
	}, baseline)
	for i := 4; i < 80; i++ {
		if i%3 == 0 {
			hotelId := strconv.Itoa(i)
			endDate := "2015-04-"
			rate := 109.00
			rateInc := 123.17
			if i%2 == 0 {
				endDate = endDate + "17"
			} else {
				endDate = endDate + "24"
			}
			if i%5 == 1 {
				rate = 120.00
				rateInc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rateInc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rateInc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rateInc = 258.00
			}
			beldilib.Populate("rate", hotelId, data.RatePlan{
				HotelId: hotelId,
				Code:    "RACK",
				Indate:  "2015-04-09",
				Outdate: endDate,
				RoomType: data.RoomType{
					BookableRate:       rate,
					Code:               "KNG",
					RoomDescription:    "King sized bed",
					TotalRate:          rate,
					TotalRateInclusive: rateInc,
				},
			}, baseline)
		}
	}
}

func recommendation(baseline bool) {
	beldilib.Populate("recommendation", "1", data.Recommend{
		HId:    "1",
		HLat:   37.7867,
		HLon:   -122.4112,
		HRate:  109.00,
		HPrice: 150.00,
	}, baseline)
	beldilib.Populate("recommendation", "2", data.Recommend{
		HId:    "2",
		HLat:   37.7854,
		HLon:   -122.4005,
		HRate:  139.00,
		HPrice: 120.00,
	}, baseline)
	beldilib.Populate("recommendation", "3", data.Recommend{
		HId:    "3",
		HLat:   37.7834,
		HLon:   -122.4071,
		HRate:  109.00,
		HPrice: 190.00,
	}, baseline)
	beldilib.Populate("recommendation", "4", data.Recommend{
		HId:    "4",
		HLat:   37.7936,
		HLon:   -122.3930,
		HRate:  129.00,
		HPrice: 160.00,
	}, baseline)
	beldilib.Populate("recommendation", "5", data.Recommend{
		HId:    "5",
		HLat:   37.7831,
		HLon:   -122.4181,
		HRate:  119.00,
		HPrice: 140.00,
	}, baseline)
	beldilib.Populate("recommendation", "6", data.Recommend{
		HId:    "6",
		HLat:   37.7863,
		HLon:   -122.4015,
		HRate:  149.00,
		HPrice: 200.00,
	}, baseline)
	for i := 7; i < 80; i++ {
		hotelId := strconv.Itoa(i)
		lat := 37.7835 + float64(i)/500.0*3
		lon := -122.41 + float64(i)/500.0*4
		rate := 135.00
		rateInc := 179.00
		if i%3 == 0 {
			if i%5 == 0 {
				rate = 109.00
				rateInc = 123.17
			} else if i%5 == 1 {
				rate = 120.00
				rateInc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rateInc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rateInc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rateInc = 258.00
			}
		}
		beldilib.Populate("recommendation", hotelId, data.Recommend{
			HId:    hotelId,
			HLat:   lat,
			HLon:   lon,
			HRate:  rate,
			HPrice: rateInc,
		}, baseline)
	}
}

func user(baseline bool) {
	for i := 0; i <= 500; i++ {
		suffix := strconv.Itoa(i)
		username := "Cornell_" + suffix
		password := ""
		for j := 0; j < 10; j++ {
			password += suffix
		}
		beldilib.Populate("user", username, data.User{
			Username: username,
			Password: password,
		}, baseline)
	}
}

func add_hotels(baseline bool) {
	for i := 0; i < 100; i++ {
		hotelId := strconv.Itoa(i)
		beldilib.Populate("hotel", hotelId, hotel.Hotel{
			HotelId:   hotelId,
			Cap:       10,
			Customers: []string{},
		}, baseline)
	}
}

func add_flights(baseline bool) {
	for i := 0; i < 100; i++ {
		flightId := strconv.Itoa(i)
		beldilib.Populate("flight", flightId, flight.Flight{
			FlightId:  flightId,
			Cap:       10,
			Customers: []string{},
		}, baseline)
	}
}

func populate(baseline bool) {
	geo(baseline)
	profile(baseline)
	rate(baseline)
	recommendation(baseline)
	user(baseline)
	add_hotels(baseline)
	add_flights(baseline)
}

func main() {
	option := os.Args[1]
	baseline := os.Args[2] == "baseline"
	if option == "create" {
		tables(baseline)
	} else if option == "populate" {
		populate(baseline)
	} else if option == "clean" {
		delete_tables(baseline)
	}
}
