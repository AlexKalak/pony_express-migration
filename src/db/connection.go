package db

import (
	"fmt"
	"log"
	"time"

	"github.com/alexkalak/migration/src/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB

func Init() *gorm.DB {
	godotenv.Load()
	// dbUser := os.Getenv("DB_USR")
	// dbPass := os.Getenv("DB_PASS")
	// dbServ := os.Getenv("DB_SERV")
	// dbPort := os.Getenv("DB_PORT")
	// dbName := os.Getenv("DB_NAME")
	// dbFullHost := dbServ + ":" + dbPort

	dsn := fmt.Sprintf("root:rootroot@tcp(localhost:3306)/pony_express?charset=utf8mb4&parseTime=True&loc=Local")
	// dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbFullHost, dbName)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//asd
	db.AutoMigrate(
		&models.ProductType{},
		&models.CountryCode{},
		&models.Country{},
		&models.City{},
		&models.DeliveryType{},
		&models.Region{},
		&models.PackageType{},
		&models.SenderCityWithOffice{},
		&models.SenderCity{},
		&models.SenderRegion{},
	)
	return db
}

func GetDB() *gorm.DB {
	if database == nil {
		database = Init()
		for database == nil {
			var sleep = time.Second
			sleep = sleep * 2
			fmt.Printf("database id unvaliable. Wait for %s seconds \n", sleep.String())

			time.Sleep(sleep)
		}
	}
	return database
}
