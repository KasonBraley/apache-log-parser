package main

import (
	"apache-log-parser/logger"
	"apache-log-parser/registry"
	"apache-log-parser/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	logger.Run("./app.log")

	// bring in from config file / env
	host, port := "logger", "4001"
	serviceAddress := fmt.Sprintf("http://%v:%v", host, port)

	var r registry.Registration
	r.ServiceName = registry.LogService
	r.ServiceURL = serviceAddress
	// does not require any services
	r.RequiredServices = make([]registry.ServiceName, 0)
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), host, port, r, logger.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
