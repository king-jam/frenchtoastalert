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
	//SnowForecastID uint
	SnowPlaceID uint
	SnowPlace   *SnowPlace `gorm:"foreignkey:SnowPlaceID"`
	//ToastID     uint64
	//Toast       *Toast `gorm:"foreignkey:ToastID"` //;association_foreignkey:slices"` //;foreignkey:snowPlaceId"`//`gorm:"foreignkey:Slices"`

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
	Place         string
	State         string
	County        string
	SnowForecasts []SnowForecast //`gorm:"foreignkey:SnowForecastID"` //`gorm:"many2many:snowplace_snowforecast;association_foreignkey:ToastID;foreignkey:snowPlaceID"` //`gorm:"foreignkey:SnowForecastID"`
	// 	//SnowForecasts []*SnowForecast //`gorm:"many2many:snowplace_snowforecast"` //;association_foreignkey:snowForecastId;foreignkey:snowPlaceId"` //;foreignkey:SnowPlaceID
}

// SnowPlaces is a list of snowplace objects which all have a forcast
type SnowPlaces []*SnowPlace

// SnowForecasts is a list of snowplace objects which all have a forcast
type SnowForecasts []*SnowForecast

// Toast has all the ingredients
type Toast struct {
	gorm.Model
	Slices uint `gorm:"primary_key:true"`
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
	Last(query *SnowPlace) (*SnowPlace, error)
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
