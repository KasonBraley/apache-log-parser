package main

import (
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

	ctx, err := service.Start(context.Background(), host, port, r, retrieve.RegisterHandlers)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down data retrieval service")
}
