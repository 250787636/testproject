package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type config struct {
	System struct {
		Version string `ini:"version"`
		Port    string `ini:"port"`
	} `ini:"system"`
	Mysql struct {
		Host     string `ini:"server"`
		Port     string `ini:"port"`
		UserName string `ini:"username"`
		Password string `ini:"password"`
		Database string `ini:"database"`
	} `ini:"mysql"`
	CMIC4AService struct {
		AppCode     string `ini:"appcode"`
		Tenant     string `ini:"tenant"`
	} `ini:"CMIC4AService"`
}

func LoadConfig() config {
	var conf config
	configPath := "config/config.ini"
	_, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	err = ini.MapTo(&conf, configPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
	return conf
}
