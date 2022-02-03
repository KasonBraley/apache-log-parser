package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (s *server) handleDataGet() http.HandlerFunc {
	// advantage for this pattern is that it allows you to set things up for the handler if needed
	// thing := prepareThing()
	// setup dependencies
	// lazy initialization with sync.Once

	// type request struct{}
	// TODO: Should this type instead be a slice of these properties?
	type response struct {
		Method      string `json:"method"`
		Status      int    `json:"status"`
		HTTPVersion int    `json:"httpversion"`
	}

	// Lazy setup. Speeds up slow starts. Only if needed
	// var init sync.Once

	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		// init.Do(func(){})

		if r.URL.Path != "/retrieve" {
			s.respond(w, r, nil, http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			s.respond(w, r, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		response := []response{}

		// Get all records. Keeping only the columns specified by the response type.
		// TODO: Does this need an error check?
		s.db.Model(&logLine{}).Find(&response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		s.respond(w, r, response, http.StatusOK)
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
