package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type server struct {
	db     *gorm.DB
	router *http.ServeMux
}

func newServer(db *gorm.DB) *server {
	srv := &server{
		db:     db,
		router: http.NewServeMux(),
	}
	srv.routes()
	return srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

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

func (s *server) handleUploadPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/upload" {
			s.respond(w, r, nil, http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			s.respond(w, r, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			s.respond(w, r, "Invalid file uploaded", http.StatusBadRequest)
			log.Println(err)
			return
		}

		file, fileHeader, err := r.FormFile("file")

		// The file cannot be received.
		if err != nil {
			s.respond(w, r, "File upload issue. Please upload a file that's less than 1MB in size", http.StatusBadRequest)
			log.Println(err)
			return
		}
		defer file.Close()

		content, err := readLog(fileHeader)
		if err != nil {
			s.respond(w, r, "Unable to open and read the log", http.StatusBadRequest)
			log.Println(err)
			return
		}

		logLines, err := parseLog(content)
		if err != nil {
			s.respond(w, r, "Unable to parse values from the log", http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := s.storeData(logLines); err != nil {
			s.respond(w, r, "Internal server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		s.respond(w, r, "File uploaded and parsed", http.StatusCreated)
	}
}

// respond is a helper function for responding to http requests.
func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println("error in respond")
			// TODO: Handle err
		}
	}
}

func (s *server) storeData(logLines []logLine) error {

	// Migrate the schema
	if err := s.db.AutoMigrate(&logLine{}); err != nil {
		return err
	}

	for _, line := range logLines {
		s.db.Create(&line)
	}

	return nil
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
