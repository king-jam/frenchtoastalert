package toast

import "fmt"

import "github.com/king-jam/ft-alert-bot/models"

// GetLevel gets the alert level based on this highly complex algorithm
func GetLevel(dataChan chan models.SnowPlaces, place models.SnowPlace) {
	for snowPlaces := range dataChan {
	Loop:
		for _, snowPlace := range snowPlaces {
			if snowPlace.ID == place.ID {
				switch snowLevel := snowPlace.SnowForecast.ExpectedSnowfall; {
				case snowLevel < 0.1: // low
					fmt.Println("low")
					break Loop
				case snowLevel >= 0.1 && snowLevel < 2.0: // guarded
					fmt.Println("guarded")
					break Loop
				case snowLevel >= 2.0 && snowLevel < 4.0: // elevated
					fmt.Println("elevated")
					break Loop
				case snowLevel >= 4.0 && snowLevel < 8.0: // high
					fmt.Println("high")
					break Loop
				case snowLevel >= 8.0: // severe
					fmt.Println("severe")
					break Loop
				default: // unknown therefore internet is likely down due to ice age
					fmt.Println("Ice Age")
					break Loop
				}
			}
		}
	}

}
