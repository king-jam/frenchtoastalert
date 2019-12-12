package main

import (
	"fmt"
	"log"
	"time"

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

	f := models.SnowForecast{
		Toast:     models.LevelTwo,
		TimeStamp: "VAL",
	}
	snowForecasts := make([]models.SnowForecast, 0)
	snowForecasts = append(snowForecasts, f)
	//f.SnowPlace = snowPlace
	snowPlace.SnowForecasts = snowForecasts
	s.DB.Create(snowPlace)

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

	//fmt.Printf("%+v", data)
}
