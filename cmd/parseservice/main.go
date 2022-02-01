package main

import (
	"apache-log-parser/logger"
	"apache-log-parser/registry"
	"apache-log-parser/service"
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Define command-line flags
	serverHost := flag.String("serverHost", "localhost", "HTTP server host name")
	serverPort := flag.Int("serverPort", 4002, "HTTP server network port")
	flag.Parse()

	var r registry.Registration
	r.ServiceName = registry.ParseService
	r.ServiceURL = fmt.Sprintf("http://%v:%v", *serverHost, *serverPort)
	r.RequiredServices = []registry.ServiceName{registry.LogService}
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), *serverHost, *serverPort, r, registerHandlers)
	if err != nil {
		log.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %v\n", logProvider)
		logger.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down parse service")
}
