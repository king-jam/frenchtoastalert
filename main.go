package main

import (
	"net/http"
	"os"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/scraper"
	"github.com/king-jam/ft-alert-bot/store"
	"github.com/king-jam/ft-alert-bot/toast"
)

// SCRAPEINTERVAL will be set with flags
const SCRAPEINTERVAL time.Duration = 3

func main() {
	log := logrus.New()
	dataChan := make(chan models.SnowForecasts, 1)
	go func() { scraper.ScrapeAndParse(SCRAPEINTERVAL*time.Second, dataChan) }()
	go func() { store.ListenAndStore(dataChan) }()
	go func() { toast.ToastApi() }()
	//go func() { serving() }()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatalf("$PORT must be set")
	}
	log.Infof("SERVING")
	router := httprouter.New()
	router.GET("/", dummyForecast)
	router.POST("/toast", toast.ToastHandler)
	log.Fatal(http.ListenAndServe(":"+port, router))
	log.Infof("SERVING")
}

func dummyForecast(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fs := http.FileServer(http.Dir("data/source"))
	fs.ServeHTTP(w, req)
}