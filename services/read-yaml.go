package services

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Node struct {
	Name   string `yaml:"-"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Path   string `yaml:"path"`
	Scheme string `yaml:"scheme"`
}

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Node []Node `yaml:"node"`
}

func ReadYaml() *Config {
	var filename string
	flag.StringVar(&filename, "f", "", "file name")
	flag.Parse()

	println(filename)
	// Check if the -f flag is set
	if filename == "" {
		log.Fatalf("error: -f flag is not set")
	} else {
		fmt.Println("-f flag value:", filename)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config := Config{}
	println("data", string(data))
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var temp = []Node{}

	for _, node := range config.Node {
		if node.Path == "*" {
			temp = append(temp, Node{
				Name:   node.Name,
				Host:   node.Host,
				Port:   node.Port,
				Path:   "",
				Scheme: node.Scheme,
			})
			break
		}
		temp = append(temp, node)

	}

	config.Node = temp
	return &config
}
