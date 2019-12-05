package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/king-jam/ft-alert-bot/models"
)

type DB struct {
	DB *gorm.DB
}

func NewDB() (*DB, error) {
	pg, err := gorm.Open("postgres", "host=localhost port=54320 user=snow dbname=snow password=snow123 sslmode=disable")
	if err != nil {
		return nil, err
	}
	db := &DB{DB: pg}
	return db, nil
}

func (db *DB) StorePlaces(snowPlace models.SnowPlace) error {
	db.DB.Create(&snowPlace)
	return nil
}
