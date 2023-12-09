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

func (h *Health) getKey(node string, url string) string {
	return node + ":" + url
}

func (h *Health) check(node string, url string, healthCheckPath string) {

	request, err := http.NewRequest("GET", url+healthCheckPath, nil)
	key := h.getKey(node, url)
	if err != nil {
		h.hSet(key, false)
		println(err.Error())
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		h.hSet(key, false)
		println(err.Error())
		return
	}
	if res.StatusCode == 200 {
		h.hSet(key, true)
	} else {
		h.hSet(key, false)
		println("status code is not 200")

	}

}

func (h *Health) checkHealth() {
	for {
		for _, node := range h.config.Node {
			for _, url := range node.Target {
				h.check(node.Name, node.Scheme+url, node.Health)
				h.checkStatus(node.Name, node.Scheme+url, node.Health)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (h *Health) checkStatus(node string, url string, healthCheckPath string) {
	key := h.getKey(node, url)
	val, isFound := h.hMap.Load(key)
	if !isFound {
		return
	}
	nodeHealth := val.(*NodeHealth)

	if h.isExpired(nodeHealth.updatedAt) {
		nodeHealth.Success = 0
		nodeHealth.Error = 0
		nodeHealth.updatedAt = time.Now()
		h.hSetNode(key, nodeHealth)
		return
	}

	total_check := nodeHealth.Success + nodeHealth.Error

	if (nodeHealth.Error / total_check * 100) > 20 {
		println("error: request fails", (total_check / nodeHealth.Error * 100))
	}

	if (nodeHealth.Success / total_check * 100) < 70 {
		println("error: request success only ", (total_check / nodeHealth.Success * 100))
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
		if h.isExpired(node.updatedAt) {
			node.updatedAt = time.Now()
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
func (h *Health) hSetNode(key string, node *NodeHealth) {
	h.hMap.Store(key, node)

}

func (h *Health) isExpired(arg time.Time) bool {
	current_time := time.Now()

	expired_time := arg.Add(5 * time.Minute)

	return current_time.Unix() > expired_time.Unix()
}
