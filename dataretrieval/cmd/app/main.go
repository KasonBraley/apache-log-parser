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
	var (
		serverHost = flag.String("serverHost", "", "HTTP server host name")
		serverPort = flag.Int("serverPort", 4001, "HTTP server network port")

		postgresHost     = flag.String("postgresHost", "localhost", "PostgreSQL host name")
		postgresPort     = flag.Int("postgresPort", 5432, "PostgreSQL port")
		postgresUser     = flag.String("postgresUser", "kason", "PostgreSQL user name")
		postgresPassword = flag.String("postgresPassword", "pass", "PostgreSQL user password")
		postgresDatabase = flag.String("postgresDatabase", "apache_logs", "PostgreSQL host name")
	)
	flag.Parse()

	serverAddr := fmt.Sprintf("%s:%d", *serverHost, *serverPort)

	if err := run(serverAddr,
		*postgresHost,
		*postgresPort,
		*postgresUser,
		*postgresPassword,
		*postgresDatabase); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(serverAddr, dbHost string, dbPort int, dbUser, dbPassword, dbDatabase string) error {
	// setup database
	// TODO: Move this to a function
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbDatabase, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	srv := newServer(db)

	return http.ListenAndServe(serverAddr, srv)
}
