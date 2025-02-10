package db

import (
	"APP4/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=123456 dbname=app4 port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db

	err = DB.AutoMigrate(
		&models.User{},
		&models.ConnectedAccounts{},
		// &models.TwitterAccessToken{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
