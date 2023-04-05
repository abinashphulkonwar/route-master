package services

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Count struct {
	Count  int
	Index  int
	Length int
}

type Node struct {
	Name string `yaml:"-"`
	// Host   string `yaml:"host"`
	// Port   string `yaml:"port"`
	Target []string `yaml:"target"`
	Path   string   `yaml:"path"`
	Scheme string   `yaml:"scheme"`
	Config *Count
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
	//	Check if the -f flag is set
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

	for index, node := range config.Node {
		println("node", node.Target[0], node.Target[1], node.Path, node.Scheme)

		if len(node.Target) == 0 {
			log.Fatalf("target is empty")
		}

		if node.Scheme == "" {
			node.Scheme = "http"

		}

		if node.Path == "*" {
			node.Path = ""

		}

		node.Config = &Count{Length: len(node.Target), Count: 0, Index: 0}
		config.Node[index] = node

	}

	return &config
}
