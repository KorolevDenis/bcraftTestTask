package DBController

import (
	"bcraftTestTask/internal/logging"
	"bcraftTestTask/internal/properties"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBController struct {
	db *sql.DB
}

func NewDBController() *DBController {
	db, err := setDBConnection()
	if err != nil {
		logger := logging.GetLogger()
		logger.Fatal(err)
	}

	return &DBController{
		db: db,
	}
}

var DbCon *DBController

func setDBConnection() (*sql.DB, error) {
	config, err := properties.GetConfig()

	logger := logging.GetLogger()
	if err != nil {
		logger.Fatal(err)
	}

	username := config.DBSettings.DBUsername
	password := config.DBSettings.DBPassword
	dbName := config.DBSettings.DBName
	dbHost := config.DBSettings.DBHost
	dbPort := config.DBSettings.DBPort

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		dbHost, username, dbName, password, dbPort) //Создать строку подключения

	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		logger.Fatal(err)
	}
	return db, nil
}

func (db_ *DBController) GetDB() *sql.DB {
	return db_.db
}
