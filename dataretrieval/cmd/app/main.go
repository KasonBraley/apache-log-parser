package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Define command-line flags
	serverHost := flag.String("serverHost", "", "HTTP server host name")
	serverPort := flag.Int("serverPort", 4001, "HTTP server network port")
	flag.Parse()

	serverAddr := fmt.Sprintf("%s:%d", *serverHost, *serverPort)

	if err := run(serverAddr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(serverAddr string) error {
	// setup database
	// TODO: Move this to a function
	dsn := "host=database user=kason password=pass dbname=apache_logs port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	srv := newServer(db)

	return http.ListenAndServe(serverAddr, srv)
}
