package main

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	HTTPVersion float64
}

func (l logLine) routes() *http.ServeMux {
	// Register handler functions.
	r := http.NewServeMux()
	r.HandleFunc("/upload", l.upload)

	return r
}

func (l logLine) upload(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/upload" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, fileHeader, err := r.FormFile("file")

	// The file cannot be received.
	if err != nil {
		http.Error(w, "File upload issue. Please upload a file that's less than 1MB in size", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer file.Close()

	content, err := readLog(fileHeader)
	if err != nil {
		http.Error(w, "Unable to open and read the log", http.StatusBadRequest)
		log.Println(err)
		return
	}

	logLines, err := parseLog(content)
	if err != nil {
		http.Error(w, "Unable to parse values from the log", http.StatusBadRequest)
		log.Println(err)
		return
	}

	db := connectDB()
	storeData(db, logLines)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded and parsed")
}

func parseLog(lines []string) ([]logLine, error) {
	var items []logLine

	for _, line := range lines {

		fields := strings.Fields(line)
		remoteHost := fields[0]
		method := strings.Trim(fields[5], "\"")
		route := fields[6]

		status, err := strconv.Atoi(fields[8])
		if err != nil {
			log.Printf("Unable to convert status code to int %s", err)
			return []logLine{}, err
		}

		layout := "02/Jan/2006:15:04:05 -0700"
		value := strings.Trim(fields[3], "[") + " " + strings.Trim(fields[4], "]")
		datetime, err := time.Parse(layout, value)
		if err != nil {
			log.Printf("unable to parse date %s", err)
			return []logLine{}, err
		}

		re := regexp.MustCompile(`([A-Z])\w+\/\d\.\d`)
		found := re.Find([]byte(fields[7]))
		trimmed := strings.Trim(string(found), "HTTP/")
		httpVersion, err := strconv.ParseFloat(trimmed, 64)

		if err != nil {
			log.Printf("Unable to convert http Version string to int %s", err)
			return []logLine{}, err
		}

		lineData := logLine{
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
	dsn := "host=database user=kason password=pass dbname=apache_logs port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func storeData(db *gorm.DB, logLines []logLine) {

	// Migrate the schema
	db.AutoMigrate(&logLine{})

	for _, line := range logLines {
		db.Create(&line)
	}
}
