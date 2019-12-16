package toast

// SetLevel gets the alert level based on this highly complex algorithm
// func (t *ToastService) SetLevel(dataChan chan models.SnowPlaces) {
// 	for snowPlaces := range dataChan {
// 	Loop:
// 		for _, snowPlace := range snowPlaces {
// 			for j, snowForecast := range snowPlace.SnowForecasts {
// 				switch snowLevel := snowForecast.ExpectedSnowfall; {
// 				case snowLevel < 0.1: // none
// 					fmt.Println("none")
// 					snowPlace.SnowForecasts[j].Toast = models.LevelZero
// 					t.Repo.Insert(snowPlace)
// 					break Loop
// 				case snowLevel >= 0.1 && snowLevel < 2.0: // low
// 					fmt.Println("low")
// 					snowPlace.SnowForecasts[j].Toast = models.LevelOne
// 					t.Repo.Insert(snowPlace)
// 					break Loop
// 				case snowLevel >= 2.0 && snowLevel < 6.0: // guarded
// 					fmt.Println("guarded")
// 					snowPlace.SnowForecasts[j].Toast = models.LevelTwo
// 					t.Repo.Insert(snowPlace)
// 					break Loop
// 				case snowLevel >= 6.0 && snowLevel < 12.0: // elevated
// 					fmt.Println("elevated")
// 					snowPlace.SnowForecasts[j].Toast = models.LevelThree
// 					t.Repo.Insert(snowPlace)
// 					fmt.Printf("%+v\n", snowPlace)
// 					break Loop
// 				case snowLevel >= 12.0 && snowLevel < 18.0: // high
// 					fmt.Println("high")
// 					snowPlace.SnowForecasts[j].Toast = models.LevelFour
// 					t.Repo.Insert(snowPlace)
// 					break Loop
// 				case snowLevel >= 18.0: // severe
// 					fmt.Println("severe")
// 					snowPlace.SnowForecasts[j].Toast = models.LevelFive
// 					t.Repo.Insert(snowPlace)
// 					break Loop
// 				default: // unknown therefore internet is likely down due to ice age
// 					fmt.Println("Ice Age")
// 					break Loop
// 				}
// 			}
// 		}
// 		return
// 	}

// }
