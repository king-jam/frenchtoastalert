package store

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/king-jam/ft-alert-bot/models"
	"github.com/sirupsen/logrus"
)

// DbService wraps the store interface funcs
type DbService struct {
	Repo models.Repository
}

// New returns an initialized DbService for making toast
func New(repo models.Repository) *DbService {
	return &DbService{Repo: repo}
}

type Store struct {
	DB *gorm.DB
}

func NewDB() (*Store, error) {
	var err error
	var pg *gorm.DB

	if os.Getenv("ENV") == "TEST" {
		pg, err = gorm.Open("postgres", "host=localhost port=54320 user=snow dbname=snow password=snow123 sslmode=disable")
	} else {
		pg, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}
	if err != nil {
		log.Fatal(err)
	}
	db := &Store{DB: pg}
	//db.DB.DropTableIfExists(&models.Location{})
	//db.DB.DropTableIfExists(&models.SnowForecasts{})
	// db.DB.DropTableIfExists(&models.SnowForecast{})
	// db.DB.DropTableIfExists(&models.Toast{})
	// db.DB.DropTableIfExists("snowplace_snowforecast")
	// db.DB.DropTableIfExists("snowforecast_toast")

	// Make SnowPlace
	if !db.DB.HasTable(&models.Location{}) {
		db.DB.CreateTable(&models.Location{}).AddUniqueIndex("index_location", "city", "state", "county")
	}

	db.DB.AutoMigrate(&models.Location{})
	// Make SnowForecast
	if !db.DB.HasTable(&models.SnowForecast{}) {
		db.DB.CreateTable(&models.SnowForecast{})
	}
	db.DB.AutoMigrate(&models.SnowForecast{})
	// Make Toast
	if !db.DB.HasTable(&models.Toast{}) {
		db.DB.CreateTable(&models.Toast{}).AddUniqueIndex("index_slices", "slices")
	}
	db.DB.AutoMigrate(&models.Toast{})

	return db, nil
}

func ListenAndStore(dataChan chan models.SnowForecasts) error {
	log := logrus.New()
	log.Infoln("Starting listen and store")
	// post snow place into database
	db, err := NewDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()
	ss := New(db)

	for {
		payload := <-dataChan
		log.Infoln("Took payload off channel")
		if payload != nil {
			go func() {
				err := ss.saveForecast(payload)
				if err != nil {
					log.Fatalln(err)
				}
			}()
		}
	}
}

func (db *DbService) saveForecast(payload models.SnowForecasts) error {
	logrus.Infoln("Started save forecast, ranging over payload")

	//Iterate through the payload
	for _, snowForecast := range payload {
		lastLocation, err := db.Repo.Last(snowForecast.Location)
		if err != nil && err != models.ErrRecordNotFound {
			return err
		}
		if lastLocation == nil || len(lastLocation.SnowForecasts) == 0 {
			logrus.Infoln("Initializing forecast table data")
			if err := db.Repo.Insert(snowForecast); err != nil {
				return err
			}
		} else if lastLocation.SnowForecasts[0].TimeStamp != snowForecast.TimeStamp {
			snowForecast.LocationID = lastLocation.SnowForecasts[0].LocationID
			if err := db.Repo.UpdateRow(snowForecast); err != nil {
				return err
			}
		}
	}
	return nil
}

func initToast(db *gorm.DB) error {
	if result := db.Create(models.LevelZero); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	if result := db.Create(models.LevelOne); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	if result := db.Create(models.LevelTwo); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	if result := db.Create(models.LevelThree); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	if result := db.Create(models.LevelFour); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	if result := db.Create(models.LevelFive); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (s *Store) Insert(snowForecast *models.SnowForecast) error {

	sp := &models.Location{Area: &models.Area{}}
	s.DB.FirstOrCreate(sp, snowForecast.Location)
	snowForecast.LocationID = sp.ID
	snowForecast.Location = sp

	if result := s.DB.Create(snowForecast); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (s *Store) UpdateRow(snowForecast *models.SnowForecast) error {
	updateTemplate := `UPDATE snow_forecasts SET 
	"time_stamp" = '%s',
	"low_end_snowfall" = %.1f, 
	"expected_snowfall" = %.1f,
	"high_end_snowfall" = %.1f, 
	"chance_more_than_zero" = %.1f,
	"chance_more_than_one" = %.1f,
	"chance_more_than_two" = %.1f,
	"chance_more_than_four" = %.1f,
	"chance_more_than_six" = %.1f,
	"chance_more_than_eight" = %.1f,
	"chance_more_than_twelve" = %.1f,
	"chance_more_than_eighteen" = %.1f
	WHERE location_id=%d`

	//FIXME: Any better way to do this?
	updateQuery := fmt.Sprintf(updateTemplate, snowForecast.TimeStamp, snowForecast.LowEndSnowfall, snowForecast.ExpectedSnowfall, snowForecast.HighEndSnowfall,
		snowForecast.ChanceMoreThanZero, snowForecast.ChanceMoreThanOne, snowForecast.ChanceMoreThanTwo, snowForecast.ChanceMoreThanFour, snowForecast.ChanceMoreThanSix,
		snowForecast.ChanceMoreThanEight, snowForecast.ChanceMoreThanTwelve, snowForecast.ChanceMoreThanEighteen, snowForecast.LocationID)
	if result := s.DB.Debug().Exec(updateQuery); result.Error != nil {

		//s.DB.Debug().Select(sf, snowForecast).Where("snowForecast.location_id =="+string(snowForecast.ID)).Update("time_stamp", snowForecast.TimeStamp); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (s *Store) LatestForecast(query *models.Location) (*models.Location, error) {
	location := new(models.Location)
	toastAlert := new(models.ToastAlert)
	snowForecasts := make([]models.SnowForecast, 0)
	if result := s.DB.Last(location, query).Related(&snowForecasts).Select(location, toastAlert); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, models.ErrRecordNotFound
		}
		return nil, models.ErrDatabaseGeneral(result.Error.Error())
	}
	location.SnowForecasts = snowForecasts
	return location, nil
}

// Last gets the last entry into the db table of snowPlaces
func (s *Store) Last(query *models.Location) (*models.Location, error) {
	location := new(models.Location)
	snowForecasts := make([]models.SnowForecast, 0)
	if result := s.DB.Last(location, query).Related(&snowForecasts); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, models.ErrRecordNotFound
		}
		return nil, models.ErrDatabaseGeneral(result.Error.Error())
	}
	location.SnowForecasts = snowForecasts
	return location, nil
}
