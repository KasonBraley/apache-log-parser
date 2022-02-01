package main

import (
	"apache-log-parser/logger"
	"apache-log-parser/registry"
	"apache-log-parser/service"
	"context"
	"flag"
	"fmt"
	stlog "log"
)

func main() {
	logger.Run("./app.log")

	serverHost := flag.String("serverHost", "localhost", "HTTP server host name")
	serverPort := flag.Int("serverPort", 4001, "HTTP server network port")

	var r registry.Registration
	r.ServiceName = registry.LogService
	r.ServiceURL = fmt.Sprintf("http://%v:%v", *serverHost, *serverPort)
	// does not require any services
	r.RequiredServices = make([]registry.ServiceName, 0)
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), *serverHost, *serverPort, r, logger.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
