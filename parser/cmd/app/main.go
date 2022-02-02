package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define command-line flags
	serverHost := flag.String("serverHost", "", "HTTP server host name")
	serverPort := flag.Int("serverPort", 4000, "HTTP server network port")
	flag.Parse()

	app := new(logLine)

	// Initialize a new http.Server struct.
	serverAddr := fmt.Sprintf("%s:%d", *serverHost, *serverPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: app.routes(),
	}

	log.Fatal(srv.ListenAndServe())
}
