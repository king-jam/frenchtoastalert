package models

import (
	"encoding/xml"

	"github.com/jinzhu/gorm"
)

// Forecast is all things from allen
type Forecast struct {
	XMLName   xml.Name `xml:"forecast"`
	Text      string   `xml:"text"`
	TimeStamp string   `xml:"timestamp"`
}

// SnowForecast is a individual snow forcast for a given place, state, county at a time
type SnowForecast struct {
	gorm.Model
	ID                     string
	TimeStamp              string
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

// SnowPlace is a list of places that get snow
type SnowPlace struct {
	gorm.Model
	ID           string
	Place        string
	State        string
	County       string
	SnowForecast *SnowForecast
}

// SnowPlaces is a list of snowplace objects which all have a forcast
type SnowPlaces []*SnowPlace

// Toast has all the ingredients
type Toast struct {
	Slices int
	Status string
	Alert  string
}
