package config

import (
	"fmt"
	"testing"
)

func TestFile(receiver *testing.T)  {
	loadConfig := LoadConfig()
	fmt.Println(loadConfig)
}
