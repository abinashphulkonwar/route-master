package services_test

import (
	"fmt"
	"testing"

	"github.com/abinashphulkonwar/route-master/services"
)

func TestReadYaml(t *testing.T) {
	config := services.ReadYaml()
	fmt.Println(config.Server.Host, config.Server.Port)
}
