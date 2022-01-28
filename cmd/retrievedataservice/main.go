package main

import (
	"apache-log-parser/logger"
	"apache-log-parser/registry"
	"apache-log-parser/retrieve"
	"apache-log-parser/service"
	"context"
	"fmt"
	"log"
)

func main() {
	// bring in from config file / env
	host, port := "localhost", "4003"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.Registration
	r.ServiceName = registry.RetrieveDataService
	r.ServiceURL = serviceAddress
	r.RequiredServices = []registry.ServiceName{registry.LogService}
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), host, port, r, retrieve.RegisterHandlers)
	if err != nil {
		log.Fatal(err)
	}

	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %v\n", logProvider)
		logger.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down data retrieval service")
}
