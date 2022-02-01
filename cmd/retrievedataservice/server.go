package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Represents some of the values of a line in a Common Apache log
type logLine struct {
	gorm.Model
	RemoteHost  string
	DateTime    time.Time
	Method      string
	Route       string
	Status      int
	HTTPVersion int
}

func registerHandlers() {
	handler := new(logLine)

	http.Handle("/retrieve", handler)
}

func (l logLine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	logLines := retrieveAllRowsFromDB(db)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logLines)
}

func connectDB() *gorm.DB {
	dsn := "host=database user=kason password=pass dbname=apache_logs port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func retrieveAllRowsFromDB(db *gorm.DB) []logLine {
	var logLines []logLine

	// Get all records
	db.Find(&logLines)

	return logLines
}