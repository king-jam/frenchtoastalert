package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/scraper"
	"github.com/king-jam/ft-alert-bot/store"
	"github.com/king-jam/ft-alert-bot/toast"
)

// SCRAPEINTERVAL will be set with flags
const SCRAPEINTERVAL time.Duration = 3

func main() {

	dataChan := make(chan models.SnowForecasts, 1)
	go func() { scraper.ScrapeAndParse(SCRAPEINTERVAL*time.Second, dataChan) }()
	fmt.Printf("scraping every %d\n", SCRAPEINTERVAL)
	go func() { store.ListenAndStore(dataChan) }()
	go func() { toast.ToastApi() }()
	go func ()  { serving()	} ()
	//Don't exit main
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	// fmt.Printf("", snowPlace)

	// t := &models.Toast{}
	// s.DB.FirstOrCreate(t, models.LevelTwo)

	//  ss := scraper.New(s)
	//  err = ss.Store(dataChan)

	// // if err != nil {
	// // 	log.Fatalln(err)
	// // }

	// data, err := s.Last(snowPlace)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// spew.Dump(data)

	// sp := &models.SnowPlace{}
	// s.DB.FirstOrCreate(sp, testSnowPlace)

	// f := &models.SnowForecast{
	// 	Toast:     t,
	// 	TimeStamp: "VALARIE",
	// 	SnowPlace: sp,
	// }
	// s.DB.Create(f)

	// OUT OBJECTS
	//snowForecasts := make([]models.SnowForecast, 0)
	// snowForecast := models.SnowForecast{}

	// snowPlaces := make([]models.SnowPlace, 0)
	//snowPlace := models.SnowPlace{}

	// toasts := make([]models.Toast, 0)
	// toast := models.Toast{}

	// QUERIES

	// GIVEN a snow place, find the forecasts

	// s.DB.First(&snowPlace).Find(&snowForecasts)
	// FIRST -> SELECT * FROM "snow_places"  WHERE "snow_places"."deleted_at" IS NULL ORDER BY "snow_places"."id" ASC LIMIT 1
	// FIND -> SELECT * FROM "snow_forecasts"  WHERE "snow_forecasts"."deleted_at" IS NULL ORDER BY "snow_forecasts"."id" ASC

	// s.DB.First(&snowPlace).Related(&snowForecasts)
	// FIRST -> SELECT * FROM "snow_places"  WHERE "snow_places"."deleted_at" IS NULL ORDER BY "snow_places"."id" ASC LIMIT 1
	// RELATED -> SELECT * FROM "snow_forecasts"  WHERE "snow_forecasts"."deleted_at" IS NULL AND (("snow_place_id" = $1)) ORDER BY "snow_forecasts"."id" ASC parameters: $1 = '1'

	// snowPlace.SnowForecasts = snowForecasts
	//s.DB.First(&snowPlace).Related(&snowForecasts)
	//snowPlace.SnowForecasts = snowForecasts

	// snowForecasts = append(snowForecasts, f)
	//f.SnowPlace = snowPlace
	// snowPlace.SnowForecasts = snowForecasts
	// data := &models.SnowForecast{}
	// toast := &models.Toast{}
	// splace := &models.SnowPlace{}
	// data.SnowPlace = splace
	// data.Toast = toast

	// snowPlaceBAD := &models.SnowPlace{
	// 	Place:  "nope",
	// 	State:  "MA",
	// 	County: "NOT",
	// }
	//s.DB.First(data, &models.SnowForecast{TimeStamp: "VALARIE"}).Related(toast).Related(splace)
	//s.DB.Set("gorm:auto_preload", true)
	// data.SnowPlace = snowPlace
	// snowPlace = snowPlaceBAD
	// I WANT THE LAST SNOWFORECAST PER SNOW PLACE

	// get all forecasts for a give snow place and then order them and pick last one
	//s.DB.Set("gorm:auto_preload", true).Find(&snowForecasts)

	// s.DB.Preload("Toast").Preload("SnowPlace").Find(&snowForecasts)

	// s.DB.First(&splace).Related(&snowForecasts, "SnowForecasts")
	// splace.SnowForecasts = snowForecasts

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
	//spew.Dump(snowPlace)

}

func serving() {
	fs := http.FileServer(http.Dir("data/source"))
	http.Handle("/", fs)

	log.Println("Serving weather data...")
	http.ListenAndServe(":7000", nil)
}
