package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/abinashphulkonwar/route-master/services"
)

func main() {

	config := services.ReadYaml()
	director := func(request *http.Request) {
		isFound := false
		for _, node := range config.Node {
			if strings.Contains(request.URL.Path, node.Path) {

				request.URL.Scheme = node.Scheme

				currentConfig := node.Config

				if currentConfig.Index < (currentConfig.Length - 1) {
					currentConfig.Index++
				} else {
					currentConfig.Index = 0
				}

				currentNode := node.Target[currentConfig.Index]
				println("currentConfig.Index", "🚀", currentConfig.Index, node.Target[currentConfig.Index], currentNode, " 🚀")

				request.URL.Host = currentNode
				log.Printf("%s %s %s %s %s %s", request.RemoteAddr, request.Method, request.URL.Path, node.Name, node.Scheme, currentNode)
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
		Rewrite: func(r *httputil.ProxyRequest) {

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
