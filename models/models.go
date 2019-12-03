package models

import "encoding/xml"

// Forecast is all things from allen
type Forecast struct {
	XMLName   xml.Name `xml:"forecast"`
	Text      string   `xml:"text"`
	TimeStamp string   `xml:"timestamp"`
}

// Snow is amazing stuff
type Snow struct {
	Place                  string
	State                  string
	County                 string
	LowEndSnowfall         float64
	ExpectedSnowfall       float64
	HighEndSnowfall        float64
	ChanceMoreThanZero     float64
	ChanceMoreThanOne      float64
	ChanceMoreThanTwo      float64
	ChanceMoreThanFour     float64
	ChanceMoreThanSix      float64
	ChanceMoreThanEight    float64
	ChanceMoreThanTwelve   float64
	ChanceMoreThanEighteen float64
}

// SnowPlaces is a list of places that get snow
type SnowPlaces []Snow

// // SnowCity is a easy lookup for snow data
// type SnowCity map[string]Snow

// // SnowPlaces is a easy lookup for snow data
// type SnowPlaces map[string]SnowCity