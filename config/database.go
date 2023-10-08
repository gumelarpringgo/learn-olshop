package config

import (
	"fmt"
	"learn/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseStruct struct {
	DB *gorm.DB
}

func ConnectDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Database Connected")

	db.AutoMigrate(
		model.User{},
		model.Address{},
		model.Product{},
		model.ProductImage{},
	)
	return db
}
