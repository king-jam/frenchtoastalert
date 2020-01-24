package models

import (
	"encoding/xml"
	"errors"
	"fmt"

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
	LocationID uint
	Location   *Location `gorm:"foreignkey:LocationID"`
	//snowForecast.LocationsnowForecast.LocationArea 					*Area `gorm:"foreignkey:LocationID"`
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
	//SnowForecastID uint
	//ToastID     uint64
	//Toast       *Toast `gorm:"foreignkey:ToastID"` //;association_foreignkey:slices"` //;foreignkey:LocationId"`//`gorm:"foreignkey:Slices"`
}

type Area struct {
	City   string
	State  string
	County string
}

// Location is a list of places that get snow
type Location struct {
	gorm.Model
	*Area
	SnowForecasts []SnowForecast //`gorm:"foreignkey:SnowForecastID"` //`gorm:"many2many:Location_snowforecast;association_foreignkey:ToastID;foreignkey:LocationID"` //`gorm:"foreignkey:SnowForecastID"`
	// 	//SnowForecasts []*SnowForecast //`gorm:"many2many:Location_snowforecast"` //;association_foreignkey:snowForecastId;foreignkey:LocationId"` //;foreignkey:LocationID
}

// Locations is a list of Location objects which all have a forcast
type Locations []*Location

// SnowForecasts is a list of Location objects which all have a forcast
type SnowForecasts []*SnowForecast

type ToastAlert struct {
	Area              Area
	ToastSnowForecast ToastSnowForecast
	Toast             Toast
}

type SlackAlert struct {
	City             string
	State            string
	ToastLevel       uint
	ExpectedSnowfall float64
	LowSnowfall      float64
	HighSnowfall     float64
}

type ToastSnowForecast struct {
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

// Toast has all the ingredients
type Toast struct {
	//gorm.Model
	Slices uint //`gorm:"primary_key:true"`
	Status string
	Alert  string
}

var LevelZero = &Toast{Slices: 0, Status: "None", Alert: ""}
var LevelOne = &Toast{Slices: 1, Status: "Low", Alert: ""}
var LevelTwo = &Toast{Slices: 2, Status: "Guarded", Alert: ""}
var LevelThree = &Toast{Slices: 3, Status: "Elevated", Alert: ""}
var LevelFour = &Toast{Slices: 4, Status: "High", Alert: ""}
var LevelFive = &Toast{Slices: 5, Status: "Severe", Alert: ""}

// Repository represent the gif usecases
type Repository interface {
	Insert(snowForecast *SnowForecast) error
	Last(query *Location) (*Location, error)
	LatestForecast(query *Location) (*Location, error)
}

var (
	// ErrRecordNotFound happens when we haven't found any matched data
	ErrRecordNotFound = errors.New("record not found")
)

// ErrDatabaseGeneral is a generic error wrapper for unexplained errors
type ErrDatabaseGeneral string

func (edg ErrDatabaseGeneral) Error() string {
	return fmt.Sprintf("General Database Error: %s", edg)
}
