package main

import (
	"fmt"
	"log"

	"github.com/king-jam/ft-alert-bot/scraper"
)

func main() {
	fmt.Println("scraping")
	err := scraper.Scraper()
	if err != nil {
		log.Fatalf(err.Error())
	}

}
