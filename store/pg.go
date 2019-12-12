package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/king-jam/ft-alert-bot/models"
)

type Store struct {
	DB *gorm.DB
}

func NewDB() (*Store, error) {
	pg, err := gorm.Open("postgres", "host=localhost port=54320 user=snow dbname=snow password=snow123 sslmode=disable")
	if err != nil {
		return nil, err
	}
	db := &Store{DB: pg}
	// db.DB.DropTableIfExists(&models.SnowPlace{})
	// db.DB.DropTableIfExists(&models.SnowForecast{})
	// db.DB.DropTableIfExists(&models.Toast{})
	// db.DB.DropTableIfExists("snowplace_snowforecast")
	// db.DB.DropTableIfExists("snowforecast_toast")
	// Make SnowPlace
	if !db.DB.HasTable(&models.SnowPlace{}) {
		db.DB.CreateTable(&models.SnowPlace{}).AddUniqueIndex("index_place", "place", "state", "county")
	}
	db.DB.AutoMigrate(&models.SnowPlace{})
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

// Store yolo puts things into the db
func (s *Store) Insert(snowPlace *models.SnowPlace) error {
	if result := s.DB.Create(snowPlace); result.Error != nil {
		return models.ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

// Last gets the last entry into the db table of snowPlaces
func (s *Store) Last(query *models.SnowPlace) (*models.SnowPlace, error) {
	snowPlace := new(models.SnowPlace)
	if result := s.DB.Where("place = ? AND state = ? AND county = ?", query.Place, query.State, query.County).First(snowPlace); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, models.ErrRecordNotFound
		}
		return nil, models.ErrDatabaseGeneral(result.Error.Error())
	}
	return snowPlace, nil
}
