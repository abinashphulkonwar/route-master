package services

import (
	"net/http"
	"sync"
	"time"
)

type NodeHealth struct {
	Name      string
	Success   int
	Error     int
	updatedAt time.Time
}

type Health struct {
	config *Config
	hMap   sync.Map
}

func NewHealth(config *Config) *Health {
	health := Health{
		config: config,
		hMap:   sync.Map{},
	}
	go health.checkHealth()
	return &health
}

func (h *Health) check(node string, url string, healthCheckPath string) {

	request, err := http.NewRequest("GET", url+healthCheckPath, nil)
	if err != nil {
		h.hSet(node+":"+url, false)
		println(err.Error())
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		h.hSet(node+":"+url, false)
		println(err.Error())
		return
	}
	if res.StatusCode == 200 {
		h.hSet(node+":"+url, true)
	} else {
		h.hSet(node+":"+url, false)
		println("status code is not 200")

	}

}

func (h *Health) checkHealth() {
	for {
		for _, node := range h.config.Node {
			for _, url := range node.Target {
				h.check(node.Name, node.Scheme+url, node.Health)
			}
		}
	}
}

func (h *Health) hSet(key string, status bool) {

	val, isFound := h.hMap.Load(key)
	if isFound {
		node := val.(*NodeHealth)
		if status {
			node.Success++
		} else {
			node.Error++
		}
		return
	}

	node := NodeHealth{
		Name:      key,
		Success:   0,
		Error:     0,
		updatedAt: time.Now(),
	}
	if status {
		node.Success++
	} else {
		node.Error++
	}

	h.hMap.Store(key, &node)

}
