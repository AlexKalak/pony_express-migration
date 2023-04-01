package standardDB

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var database *sql.DB

func Init() *sql.DB {
	godotenv.Load()
	dbUser := os.Getenv("DB_USR")
	dbPass := os.Getenv("DB_PASS")
	dbServ := os.Getenv("DB_SERV")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbFullHost := dbServ + ":" + dbPort

	// dsn := fmt.Sprintf("root:rootroot@tcp(localhost:3306)/pony_express_dev")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbFullHost, dbName)
	var err error
	database, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	fmt.Println(dsn)
	return database
}

func GetDB() *sql.DB {
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
