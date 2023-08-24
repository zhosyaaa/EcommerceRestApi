package db

import (
	"Ecommerce/internal/pkg/config"
	"Ecommerce/internal/pkg/db/models"
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
	sqldb, err := d.DB()
	if err != nil {
		log.Fatal(dbErr)
	}
	maxIdleConns := 10
	maxOpenConns := 100
	sqldb.SetMaxIdleConns(maxIdleConns)
	sqldb.SetMaxOpenConns(maxOpenConns)
	db = d
	errr := db.AutoMigrate(&models.User{}, &models.Address{}, &models.Order{}, &models.ProductsToOrder{}, &models.Product{})
	if errr != nil {
		return
	}
}

func GetDB() *gorm.DB {
	return db
}
