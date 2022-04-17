package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/nikitamirzani323/wl_apiagen/helpers"
)

var db *sql.DB
var err error

func Init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}
	var conString string = ""

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSchema := os.Getenv("DB_SCHEMA")

	switch dbDriver {
	case "cloudsql":
		var (
			instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		)

		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}

		conString = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPass, socketDir, instanceConnectionName, dbName)
		db, err = sql.Open("mysql", conString)

	case "cloudpostgres":
		var (
			instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		)

		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}
		// conString := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", dbHost, dbUser, dbPass, dbPort, dbName)

		conString = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", dbUser, dbPass, dbName, socketDir, instanceConnectionName)

		// dbPool is the pool of database connections.
		db, err = sql.Open("postgres", conString)

	case "postgres":
		log.Printf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s search_path=%s", dbHost, dbPort, dbUser, dbName, dbPass, dbSchema)
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s search_path=%s ", dbHost, dbPort, dbUser, dbName, dbPass, dbSchema)
		db, err = sql.Open(dbDriver, DBURL)
	default:
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		db, err = sql.Open(dbDriver, DBURL)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	err = db.Ping()
	helpers.ErrorCheck(err)
}
func CreateCon() *sql.DB {
	return db
}
