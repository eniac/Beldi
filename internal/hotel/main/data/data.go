package data

import "github.com/eniac/Beldi/pkg/beldilib"

type Address struct {
	StreetNumber string
	StreetName   string
	City         string
	State        string
	Country      string
	PostalCode   string
	Lat          float64
	Lon          float64
}

type Hotel struct {
	Id          string
	Name        string
	PhoneNumber string
	Description string
	Address     Address
}

type RoomType struct {
	BookableRate       float64
	TotalRate          float64
	TotalRateInclusive float64
	Code               string
	RoomDescription    string
}

type RatePlan struct {
	HotelId  string
	Code     string
	Indate   string
	Outdate  string
	RoomType RoomType
}

type Recommend struct {
	HId    string
	HLat   float64
	HLon   float64
	HRate  float64
	HPrice float64
}

type Reservation struct {
	HotelId      string
	CustomerName string
	InDate       string
	OutDate      string
	Number       int
}

type Number struct {
	HotelId string
	Num     int32
}

type User struct {
	Username string
	Password string
}

type Point struct {
	Pid  string  `mapstructure:"hotelId" json:"hotelId"`
	Plat float64 `mapstructure:"lat" json:"lat"`
	Plon float64 `mapstructure:"lon" json:"lon"`
}

func (p Point) Lat() float64 { return p.Plat }
func (p Point) Lon() float64 { return p.Plon }
func (p Point) Id() string   { return p.Pid }

type RPCInput struct {
	Function string
	Input    interface{}
}

func Tgeo() string {
	if beldilib.TYPE == "BASELINE" {
		return "bgeo"
	} else {
		return "geo"
	}
}

func Tflight() string {
	if beldilib.TYPE == "BASELINE" {
		return "bflight"
	} else {
		return "flight"
	}
}

func Tfrontend() string {
	if beldilib.TYPE == "BASELINE" {
		return "bfrontend"
	} else {
		return "frontend"
	}
}

func Tgateway() string {
	if beldilib.TYPE == "BASELINE" {
		return "bgateway"
	} else {
		return "gateway"
	}
}

func Thotel() string {
	if beldilib.TYPE == "BASELINE" {
		return "bhotel"
	} else {
		return "hotel"
	}
}

func Torder() string {
	if beldilib.TYPE == "BASELINE" {
		return "border"
	} else {
		return "order"
	}
}

func Tprofile() string {
	if beldilib.TYPE == "BASELINE" {
		return "bprofile"
	} else {
		return "profile"
	}
}

func Trate() string {
	if beldilib.TYPE == "BASELINE" {
		return "brate"
	} else {
		return "rate"
	}
}

func Trecommendation() string {
	if beldilib.TYPE == "BASELINE" {
		return "brecommendation"
	} else {
		return "recommendation"
	}
}

func Tsearch() string {
	if beldilib.TYPE == "BASELINE" {
		return "bsearch"
	} else {
		return "search"
	}
}

func Tuser() string {
	if beldilib.TYPE == "BASELINE" {
		return "buser"
	} else {
		return "user"
	}
}
