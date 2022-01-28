package retrieve

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Represents some of the values of a line in a Common Apache log
type LogLine struct {
	gorm.Model
	RemoteHost  string
	DateTime    time.Time
	Method      string
	Route       string
	Status      int
	HTTPVersion int
}

func RegisterHandlers() {
	handler := new(LogLine)

	http.Handle("/retrieve", handler)
}

func (l LogLine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	logLines := retrieveAllRowsFromDB(db)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logLines)
}

func connectDB() *gorm.DB {
	dsn := "host=localhost user=kason password=pass dbname=apache_logs port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func retrieveAllRowsFromDB(db *gorm.DB) []LogLine {
	var logLines []LogLine

	// Get all records
	db.Find(&logLines)

	return logLines
}
