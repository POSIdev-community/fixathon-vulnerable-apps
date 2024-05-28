package models

import (
	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/sirupsen/logrus"
)

var db *Database

type Database struct {
	*sql.DB
	closeDb bool
}

func initDbConnection() {
	if db != nil {
		if err := db.Ping(); err != nil {
			SetDbConnection(initCosmicDb(), true)
		}
		return
	}
	SetDbConnection(initCosmicDb(), true)
}

func SetDbConnection(dbConnection *sql.DB, close bool) {
	db = &Database{DB: dbConnection, closeDb: close}
}

func RemoveDbConnection() {
	if db != nil {
		db = nil
	}
}

func closeDbConnection() {
	if db != nil && db.closeDb {
		db.Close()
	}
}

func initCosmicDb() *sql.DB {
	dbConnection, err := sql.Open("sqlite3",
		"../cosmic_db.sqlite")
	dbConnection.SetMaxOpenConns(1)
	if err != nil {
		logrus.Fatal(err)
	}

	return dbConnection
}
