package db

import (
	"Ecommerce/pkg/config"
	"Ecommerce/pkg/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func ConnectDB(cfg config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	d, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db = d
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	err := db.AutoMigrate(&models.User{}, &models.Address{}, &models.ProductsToOrder{}, &models.Order{}, &models.Product{})
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
