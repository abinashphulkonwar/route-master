package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	file   *os.File
}

func NewHealth(config *Config) *Health {
	file_ref, err := os.OpenFile("health.log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	health := Health{
		config: config,
		hMap:   sync.Map{},
		file:   file_ref,
	}
	go health.checkHealth()
	return &health
}

func (h *Health) getKey(node string, url string) string {
	return node + ":" + url
}

func (h *Health) check(node string, url string, healthCheckPath string) {
	key := h.getKey(node, url)
	request, err := http.NewRequest("GET", url+healthCheckPath, nil)
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
	defer h.file.Close()
	for {
		for _, node := range h.config.Node {
			for _, url := range node.Target {
				h.check(node.Name, node.Scheme+"://"+url, node.Health)
				h.checkStatus(node.Name, node.Scheme+"://"+url, node.Health)
			}
		}
		h.updateTheService()
		time.Sleep(25 * time.Second)
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

	if nodeHealth.Error > 0 && ((nodeHealth.Error / total_check * 100) > 20) {
		println("error: request fails", (total_check / nodeHealth.Error * 100))
		message := fmt.Sprintf("error: request fails  %d%; "+key+"\n", (total_check / nodeHealth.Error * 100))
		_, err := h.file.WriteString(message)
		if err != nil {
			println(err)
		}

	}

	if nodeHealth.Success > 0 && ((nodeHealth.Success / total_check * 100) < 70) {
		println("error: request success only ", (total_check / nodeHealth.Success * 100))
		message := fmt.Sprintf("error: request success only %d%; "+key+"\n", (total_check / nodeHealth.Success * 100))
		_, err := h.file.WriteString(message)
		if err != nil {
			println(err)
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

func (h *Health) updateTheService() {
	if h.config.Internal.Target == "" || h.config.Internal.Scheme == "" {
		return
	}
	url := h.config.Internal.Scheme + "://" + h.config.Internal.Target
	var nodes []*NodeHealth
	h.hMap.Range(func(key, value interface{}) bool {
		node := value.(*NodeHealth)
		nodes = append(nodes, node)
		return true
	})
	body, err := json.Marshal(nodes)
	if err != nil {
		println(err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		println(err)
		return

	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println(err)
		return
	}

	if res.StatusCode != 200 {
		println("error: request fails to updates health status at: ", url)
		h.file.WriteString("error: request fails to updates health status at: " + url + "/n")
		return
	}
	println("success: request updates health status at: ", url, req.Response.StatusCode)

}
