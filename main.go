package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/abinashphulkonwar/route-master/services"
)

// api/v1/coupon

func main() {
	services.StartNode()
	config := services.ReadYaml()
	director := func(request *http.Request) {
		isFound := false
		for _, node := range config.Node {

			if strings.Contains(request.URL.Path, node.Path) {
				log.Printf("%s %s %s %s %s %s %s", request.RemoteAddr, request.Method, request.URL.Path, node.Name, node.Scheme, node.Host, node.Port)

				request.URL.Scheme = node.Scheme
				request.URL.Host = node.Host + ":" + node.Port
				isFound = true
				break
			}
		}
		if !isFound {
			log.Panicln(isFound, " 🚀")
			log.Printf("%s %s %s", request.RemoteAddr, request.Method, request.URL.Path)

			request.Write(bytes.NewBufferString("not found"))
		}

	}

	rp := &httputil.ReverseProxy{
		Director: director,
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
