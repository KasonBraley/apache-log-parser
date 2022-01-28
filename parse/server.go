package parse

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
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

	http.Handle("/upload", handler)
}

func (l LogLine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, fileHeader, err := r.FormFile("file")

	// The file cannot be received.
	if err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := readLog(fileHeader)
	if err != nil {
		http.Error(w, "Unable to open and read the log", http.StatusBadRequest)
		return
	}

	logLines, err := parseLog(content)
	if err != nil {
		http.Error(w, "Unable to parse values from the log", http.StatusBadRequest)
		return
	}

	db := connectDB()
	storeData(db, logLines)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(201)
	fmt.Fprintf(w, "File uploaded and parsed")
}

func parseLog(lines []string) ([]LogLine, error) {
	var items []LogLine

	for _, line := range lines {

		fields := strings.Fields(line)
		remoteHost := fields[0]
		method := strings.Trim(fields[5], "\"")
		route := fields[6]

		status, err := strconv.Atoi(fields[8])
		if err != nil {
			log.Printf("Unable to convert status code to int %s", err)
			return []LogLine{}, err
		}

		layout := "02/Jan/2006:15:04:05 -0700"
		value := strings.Trim(fields[3], "[") + " " + strings.Trim(fields[4], "]")
		datetime, err := time.Parse(layout, value)
		if err != nil {
			log.Printf("unable to parse date %s", err)
			return []LogLine{}, err
		}

		httpVersion, err := strconv.Atoi(strings.Trim(fields[7], "HTTP/.0\""))
		if err != nil {
			log.Printf("Unable to convert http Version string to int %s", err)
			return []LogLine{}, err
		}

		lineData := LogLine{
			RemoteHost:  remoteHost,
			DateTime:    datetime,
			Method:      method,
			Route:       route,
			Status:      status,
			HTTPVersion: httpVersion,
		}

		items = append(items, lineData)
	}

	return items, nil
}

func readLog(file *multipart.FileHeader) ([]string, error) {
	f, err := file.Open()
	if err != nil {
		log.Printf("unable to open uploaded file %s", err)
		return []string{}, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func connectDB() *gorm.DB {
	dsn := "host=localhost user=kason password=pass dbname=apache_logs port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func storeData(db *gorm.DB, logLines []LogLine) {

	// Migrate the schema
	db.AutoMigrate(&LogLine{})

	for _, line := range logLines {
		db.Create(&line)
	}
}
