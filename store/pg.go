package store

import (
	"log"
	//"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/king-jam/ft-alert-bot/models"
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
	 pg, err := gorm.Open("postgres", "host=localhost port=54320 user=snow dbname=snow password=snow123 sslmode=disable")
	// if err != nil {
	// 	return nil, err
	// }
	//pg, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
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

	// init toast table
	// err = initToast(db.DB)
	// if err != nil {
	// 	return nil, err
	// }

	return db, nil
}

func ListenAndStore(dataChan chan models.SnowForecasts) error {

	// post snow place into database
	db, err := NewDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.DB.Close()
	ss := New(db)

	for {
		payload := <-dataChan
		if payload != nil {
			go func() {
				err := ss.saveForecast(payload)
				if err != nil {
					log.Fatalln(err)
				}
			}()
		}
	}
	return nil
}

func (db *DbService) saveForecast(payload models.SnowForecasts) error {

	//Iterate through the payload
	for _, snowForecast := range payload {
		lastLocation, err := db.Repo.Last(snowForecast.Location)
		if err != nil && err != models.ErrRecordNotFound {
			return err
		}
		if lastLocation == nil || len(lastLocation.SnowForecasts) == 0 {
			if err := db.Repo.Insert(snowForecast); err != nil {
				return err
			}
		} else if lastLocation.SnowForecasts[0].TimeStamp != snowForecast.TimeStamp {
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

// func (s *Store) InsertLocation(location *models.Location) error {

// 	sp := &models.Location{}
// 	s.DB.FirstOrCreate(sp, Location)

// 	if result := s.DB.Create(snowForecast); result.Error != nil {
// 		return models.ErrDatabaseGeneral(result.Error.Error())
// 	}
// 	return nil
// }

func (s *Store) Insert(snowForecast *models.SnowForecast) error {

	sp := &models.Location{Area: &models.Area{}}
	s.DB.FirstOrCreate(sp, snowForecast.Location)
	// s.DB.Last(sp, snowPlace)
	// sp.SnowForecasts = snowPlace.SnowForecasts
	// s.DB.FirstOrCreate(sp, testSnowPlace)

	// snowPlace.SnowForecasts[0].SnowPlaceID = sp.ID
	// &snowPlace.SnowForecasts[0]
	snowForecast.LocationID = sp.ID
	snowForecast.Location = sp
	//snowForecast.Location.Area = sp.Area

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
	//area := new(models.Area)
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
