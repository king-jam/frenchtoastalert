package main

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/scraper"
	"github.com/king-jam/ft-alert-bot/toast"
)

// SCRAPEINTERVAL will be set with flags
const SCRAPEINTERVAL time.Duration = 10

func main() {
	fmt.Printf("scraping every %d", SCRAPEINTERVAL)

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

	toast.GetLevel(dataChan, snowPlace)

	// call alert logic here

}
