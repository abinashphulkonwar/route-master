package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/abinashphulkonwar/route-master/handler"
	"github.com/abinashphulkonwar/route-master/logger"
	"github.com/abinashphulkonwar/route-master/services"
)

func main() {

	config := services.ReadYaml()
	Logger := logger.NewLogger()
	services.NewHealth(config)

	rp := &httputil.ReverseProxy{
		Director: handler.Director(config, Logger),
		Transport: &handler.CustomTransport{
			Transport: http.DefaultTransport,
		},
	}

	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	if config.Server.Port == "" {
		config.Server.Port = "3002"
	}

	host := config.Server.Host + ":" + config.Server.Port
	s := http.Server{
		Addr:    host,
		Handler: rp,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
