package config

import (
	"github.com/ftheodoro/exchange-dollar/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./database/db.sqlite"), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	db.AutoMigrate(model.ExchangeRate{})
	return db, nil
}
