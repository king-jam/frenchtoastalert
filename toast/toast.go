package toast

import (
	"encoding/json"
	"fmt"
	"github.com/king-jam/ft-alert-bot/models"
	"github.com/king-jam/ft-alert-bot/store"
	"log"
	"net/http"
)

func ToastApi() {

	http.HandleFunc("/toast", ToastHandler)
	http.ListenAndServe(":7777", nil)
	log.Println("Serving toast")
}

func ToastHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		log.Println("Handing POST")
		decoder := json.NewDecoder(req.Body)
		var areaData models.Area
		decoder.Decode(&areaData)
		toastLevel := CheckToast(&areaData)
		jsonResp, err := json.Marshal(toastLevel)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Fprintf(w, string(jsonResp))

	} else {
		fmt.Fprintf(w, "404?")
	}
}

func CheckToast(area *models.Area) *models.Location {
	// Setup DB connection
	db, err := store.NewDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()
	ss := store.New(db)

	location := &models.Location{
		Area: &models.Area{
			City:   area.City,
			State:  area.State,
			County: area.County,
		},
	}

	var locationForecast *models.Location
	fmt.Println("Santy1")
	locationForecast, err = ss.Repo.LatestForecast(location)
	fmt.Println("Santy2")
	fmt.Println(locationForecast)
	if err != nil {
		log.Fatalln(err)
	}

	toastLevel := SetLevel(locationForecast)
	toastAlert := makeFriendlyToastAlert(locationForecast, toastLevel)

	fmt.Println(toastAlert)
	return locationForecast
}

func makeFriendlyToastAlert(toastLocation *models.Location, toastLevel *models.Toast) models.ToastAlert {
	toastForecast := &models.ToastSnowForecast{
		TimeStamp : 				toastLocation.SnowForecasts[0].TimeStamp,
		LowEndSnowfall : 			toastLocation.SnowForecasts[0].LowEndSnowfall,
		ExpectedSnowfall : 			toastLocation.SnowForecasts[0].ExpectedSnowfall,
		HighEndSnowfall : 			toastLocation.SnowForecasts[0].HighEndSnowfall,
		ChanceMoreThanZero : 		toastLocation.SnowForecasts[0].ChanceMoreThanZero,
		ChanceMoreThanOne : 		toastLocation.SnowForecasts[0].ChanceMoreThanOne,
		ChanceMoreThanTwo : 		toastLocation.SnowForecasts[0].ChanceMoreThanTwo,
		ChanceMoreThanFour:         toastLocation.SnowForecasts[0].ChanceMoreThanFour,  
		ChanceMoreThanSix :     	toastLocation.SnowForecasts[0].ChanceMoreThanSix,      
		ChanceMoreThanEight:    	toastLocation.SnowForecasts[0].ChanceMoreThanEight,    
		ChanceMoreThanTwelve :  	toastLocation.SnowForecasts[0].ChanceMoreThanTwelve,   
		ChanceMoreThanEighteen :	toastLocation.SnowForecasts[0].ChanceMoreThanEighteen, 
	}

	toastAlert := models.ToastAlert{
		Area: *toastLocation.Area,
		ToastSnowForecast: *toastForecast,
		Toast: *toastLevel,
	}

	for forecast := range toastLocation.SnowForecasts{
		fmt.Println("--------------------", forecast)
	}
	return toastAlert
}

//SetLevel gets the alert level based on this highly, highly complex algorithm
func SetLevel(locationForecast *models.Location) *models.Toast {
	for _, forecast := range locationForecast.SnowForecasts {
		snowLevel := forecast.ExpectedSnowfall
		switch {
		case snowLevel < 0.1: // none
			return models.LevelZero
		case snowLevel >= 0.1 && snowLevel < 2.0: // low
			return models.LevelOne
		case snowLevel >= 2.0 && snowLevel < 6.0: // guarded
			return models.LevelTwo
		case snowLevel >= 6.0 && snowLevel < 12.0: // elevated
			return models.LevelThree
		case snowLevel >= 12.0 && snowLevel < 18.0: // high
			return models.LevelFour
		case snowLevel >= 18.0: // severe
			return models.LevelFive
		default: 
			//FIXME: unknown therefore internet is likely down due to ice age
			fmt.Println("Ice Age")
			return models.LevelFive
		}
	}
	return nil
}
