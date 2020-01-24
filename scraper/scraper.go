package scraper

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/king-jam/ft-alert-bot/models"
	//"github.com/sirupsen/logrus"
)

var timeStamp string

func doEvery(d time.Duration, dataChan chan models.SnowForecasts, f func() (*models.Forecast, error)) error {
	for range time.Tick(d) {
		forecast, err := f()
		if err != nil {
			return err
		}
		if forecast.TimeStamp != timeStamp {
			Locations, err := Parser(forecast)
			if err != nil {
				return err
			}

			timeStamp = forecast.TimeStamp
			dataChan <- Locations
		}
	}
	return nil
}

// ScrapeAndParse goes and gets the data every set duration
func ScrapeAndParse(d time.Duration, dataChan chan models.SnowForecasts) error {
	return doEvery(d, dataChan, Scraper)
}

// Scraper scrapes
func Scraper() (*models.Forecast, error) {
	//resp, err := http.Get("https://www.weather.gov/source/box/winter/snow_prob.xml")
	forecastxml := fmt.Sprintf("http://localhost" + ":" + os.Getenv("PORT") + "/snow_prob.xml")
	//logrus.Infoln(forecastxml)
	resp, err := http.Get(forecastxml)

	if err != nil {
		// handle error
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var forecast models.Forecast
	xml.Unmarshal(body, &forecast)
	return &forecast, nil
}

// Parser parses
func Parser(forecast *models.Forecast) (models.SnowForecasts, error) {
	//snowPlaces := make(models.SnowPlaces, 0)
	snowForecasts := make([]*models.SnowForecast, 0)

	line := strings.Split(forecast.Text, "\n")
	for _, v := range line {

		// input validation for bad strings from xml parse
		if v == " " || v == "" {
			continue
		}
		lineItems := strings.Split(v, ",")

		Location := &models.Location{
			Area: &models.Area{
				City:   lineItems[0],
				State:  lineItems[1],
				County: lineItems[2],
			},
		}

		lowEndSnowfall, err := strconv.ParseFloat(lineItems[4], 64)
		if err != nil {
			return nil, err
		}
		expectedSnowfall, err := strconv.ParseFloat(lineItems[5], 64)
		if err != nil {
			return nil, err
		}
		highEndSnowfall, err := strconv.ParseFloat(lineItems[6], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanZero, err := strconv.ParseFloat(lineItems[7], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanOne, err := strconv.ParseFloat(lineItems[8], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanTwo, err := strconv.ParseFloat(lineItems[9], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanFour, err := strconv.ParseFloat(lineItems[10], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanSix, err := strconv.ParseFloat(lineItems[11], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanEight, err := strconv.ParseFloat(lineItems[12], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanTwelve, err := strconv.ParseFloat(lineItems[13], 64)
		if err != nil {
			return nil, err
		}
		chanceMoreThanEighteen, err := strconv.ParseFloat(lineItems[14], 64)
		if err != nil {
			return nil, err
		}
		snowForecast := &models.SnowForecast{
			Location:               Location,
			TimeStamp:              forecast.TimeStamp,
			LowEndSnowfall:         lowEndSnowfall,
			ExpectedSnowfall:       expectedSnowfall,
			HighEndSnowfall:        highEndSnowfall,
			ChanceMoreThanZero:     chanceMoreThanZero,
			ChanceMoreThanOne:      chanceMoreThanOne,
			ChanceMoreThanTwo:      chanceMoreThanTwo,
			ChanceMoreThanFour:     chanceMoreThanFour,
			ChanceMoreThanSix:      chanceMoreThanSix,
			ChanceMoreThanEight:    chanceMoreThanEight,
			ChanceMoreThanTwelve:   chanceMoreThanTwelve,
			ChanceMoreThanEighteen: chanceMoreThanEighteen,
		}
		snowForecasts = append(snowForecasts, snowForecast)
	}
	return snowForecasts, nil
}

// ScraperService wraps the store interface funcs
type ScraperService struct {
	Repo models.Repository
}

// New returns an initialized ScraperService for making toast
func New(repo models.Repository) *ScraperService {
	return &ScraperService{Repo: repo}
}

func (ss *ScraperService) Store(dataChan chan models.SnowForecasts) error {
	for snowForecasts := range dataChan {
		//Loop:
		for _, snowForecast := range snowForecasts {
			if err := ss.Repo.Insert(snowForecast); err != nil {
				return err
			}
		}
	}
	return nil
}
