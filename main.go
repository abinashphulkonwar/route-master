package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/abinashphulkonwar/route-master/logger"
	"github.com/abinashphulkonwar/route-master/services"
)

func main() {

	config := services.ReadYaml()
	Logger := logger.NewLogger()
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
				// create logger.log

				Logger.Log(&logger.Log{
					Method:  request.Method,
					Path:    request.URL.Path,
					Address: request.RemoteAddr,
					Scheme:  node.Scheme,
					Name:    node.Name,
				})
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
		Transport: &customTransport{
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

type customTransport struct {
	Transport http.RoundTripper
}

func (c *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		resp = &http.Response{
			Status:        "Internal Server Error",
			StatusCode:    500,
			Proto:         req.Proto,
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          io.NopCloser(bytes.NewBufferString("hii")),
			ContentLength: int64(len("hii")),
			Request:       req,
			Header:        make(http.Header, 0),
		}

		return resp, nil

	}

	return resp, err
}
