package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// bring in from config file / env
	serverHost := flag.String("serverHost", "", "HTTP server host name")
	serverPort := flag.Int("serverPort", 4001, "HTTP server network port")

	app := new(logLine)

	// Initialize a new http.Server struct.
	serverAddr := fmt.Sprintf("%s:%d", *serverHost, *serverPort)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: app.routes(),
	}

	log.Fatal(srv.ListenAndServe())
}
