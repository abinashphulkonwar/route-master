package services

import (
	"net/http"
	"sync"
)

type Health struct {
	config *Config
	hMap   sync.Map
}

func NewHealth(config *Config) *Health {

	health := Health{
		config: config,
		hMap:   sync.Map{},
	}
	go health.check()
	return &health
}

func (h *Health) check() error {
	for _, node := range h.config.Node {
		for _, url := range node.Target {
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				println(err.Error())
				continue
			}
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				println(err.Error())
				continue
			}
			if res.StatusCode != 200 {
				println("status code is not 200")
				continue
			}

		}
	}
	return nil
}
