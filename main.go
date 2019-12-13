package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/scraper"
	"github.com/king-jam/ft-alert-bot/store"
)

// SCRAPEINTERVAL will be set with flags
const SCRAPEINTERVAL time.Duration = 3

func main() {
	fmt.Printf("scraping every %d\n", SCRAPEINTERVAL)

	dataChan := make(chan models.SnowPlaces, 1)
	// somehow error handle this
	go func() { scraper.ScrapeAndParse(SCRAPEINTERVAL*time.Second, dataChan) }()

	// make a place to subscribe to alerts from
	place := "Waterville"
	state := "ME"
	county := "Kennebec"

	snowPlace := &models.SnowPlace{
		Place:  place,
		State:  state,
		County: county,
	}

	// post snow place into database
	s, err := store.NewDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer s.DB.Close()

	t := &models.Toast{}
	s.DB.FirstOrCreate(t, models.LevelTwo)
	sp := &models.SnowPlace{}
	s.DB.FirstOrCreate(sp, snowPlace)

	f := &models.SnowForecast{
		Toast:     t,
		TimeStamp: "VALARIE",
		SnowPlace: sp,
	}
	snowForecasts := make([]models.SnowForecast, 0)
	// snowForecasts = append(snowForecasts, f)
	//f.SnowPlace = snowPlace
	// snowPlace.SnowForecasts = snowForecasts
	s.DB.Create(f)
	data := &models.SnowForecast{}
	// toast := &models.Toast{}
	splace := &models.SnowPlace{}
	data.SnowPlace = splace
	// data.Toast = toast

	// snowPlaceBAD := &models.SnowPlace{
	// 	Place:  "nope",
	// 	State:  "MA",
	// 	County: "NOT",
	// }
	//s.DB.First(data, &models.SnowForecast{TimeStamp: "VALARIE"}).Related(toast).Related(splace)
	// s.DB.Set("gorm:auto_preload", true)
	data.SnowPlace = snowPlace
	// snowPlace = snowPlaceBAD
	// I WANT THE LAST SNOWFORECAST PER SNOW PLACE

	// get all forecasts for a give snow place and then order them and pick last one
	s.DB.Find(&splace)

	s.DB.Preload("Toast").Preload("SnowPlace", "state = ?", "MA").Find(&snowForecasts)

	// s.DB.Preload("Toast").Preload("SnowPlace").Where("TimeStamp").Find(&snowForecasts)

	// s.DB.Preload("SnowForecasts").First(splace)

	// s.DB.First(splace, snowPlace).Related(&snowForecasts).Related(toast)
	// splace.SnowForecasts = snowForecasts
	// fmt.Printf("Forecast %+v\n\n\n", data.ID)

	// // s.DB.Model(data).Related(toast)
	// fmt.Printf("Toast %+v\n\n\n", toast.Slices)
	//data.Toast = toast

	// // call alert logic here
	// ts := toast.New(s)
	// ts.SetLevel(dataChan)

	// data, err := ts.Repo.Last(snowPlace)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	//data := models.SnowPlace{}

	// s.DB.Last(data, snowPlace)
	//s.DB.Model(&data).Related(&snowForecasts, "SnowForecasts")
	spew.Dump(snowForecasts)

}
