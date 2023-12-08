package handler

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abinashphulkonwar/route-master/logger"
	"github.com/abinashphulkonwar/route-master/services"
)

func Director(config *services.Config, Logger *logger.Logger) func(request *http.Request) {
	return func(request *http.Request) {
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

				request.URL.Host = currentNode
				// create logger.log

				Logger.Log(&logger.Log{
					Method:  request.Method,
					Path:    request.URL.Path,
					Address: request.RemoteAddr,
					Scheme:  node.Scheme,
					Name:    node.Name,
					Time:    time.Now(),
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

}
