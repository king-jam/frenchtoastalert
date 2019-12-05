package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/scraper"
	"github.com/king-jam/ft-alert-bot/store"
	"github.com/king-jam/ft-alert-bot/toast"
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

	hash := sha1.New()
	hash.Write([]byte(place + state + county))
	snowPlace := models.SnowPlace{
		ID:           fmt.Sprintf("%x", hash.Sum(nil)),
		Place:        place,
		State:        state,
		County:       county,
		SnowForecast: nil,
	}

	// call alert logic here
	toast.GetLevel(dataChan, snowPlace)

	// post snow place into database
	db, err := store.NewDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()
	err = db.StorePlaces(snowPlace)
	if err != nil {
		log.Fatalln(err)
	}
	getPlace := &models.SnowPlace{
		ID: fmt.Sprintf("%x", hash.Sum(nil)),
	}
	data := &models.SnowPlace{}
	db.DB.First(data, getPlace)
	fmt.Printf("%+v", data)
}
