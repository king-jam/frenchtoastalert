package store

import (
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
	//pg, err := gorm.Open("postgres", "host=localhost port=54320 user=snow dbname=snow password=snow123 sslmode=disable")
	pg, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
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
			logrus.Infoln("Inserting new forecast data")
			if err := db.Repo.Insert(snowForecast); err != nil {
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

func (s *Store) LatestForecast(query *models.Location) (*models.Location, error) {
	location := new(models.Location)
	toastAlert := new(models.ToastAlert)
	snowForecasts := make([]models.SnowForecast, 0)
	if result := s.DB.Last(location, query).Related(&snowForecasts).Select(location, toastAlert); result.Error != nil {
		//if result := s.DB.Last(location, query).Related(&snowForecasts); result.Error != nil {
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
