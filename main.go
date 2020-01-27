package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/scraper"
	"github.com/king-jam/ft-alert-bot/store"
	"github.com/king-jam/ft-alert-bot/toast"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// SCRAPEINTERVAL will be set with flags
const SCRAPEINTERVAL time.Duration = 1

func main() {
	log := logrus.New()
	log.Infoln("Starting main")

	dataChan := make(chan models.SnowForecasts, 1)
	go func() { scraper.ScrapeAndParse(SCRAPEINTERVAL*time.Second, dataChan) }()
	go func() { store.ListenAndStore(dataChan) }()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("$PORT must be set")
	}

	router := httprouter.New()
	fs := http.FileSystem(http.Dir("data/source"))
	router.ServeFiles("/*filepath", fs)
	router.POST("/toast", toast.ToastHandler)
	log.Fatal(http.ListenAndServe(":"+port, router))
	log.Infof("SERVING")
}
