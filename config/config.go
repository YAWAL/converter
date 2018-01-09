package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type Index struct {
	Page      int
	Name      int
	RoleStart int
	RoleEnd   int
}

type IndexString struct {
	Page      string `json:"page"`
	Name      string `json:"name"`
	RoleStart string `json:"role_start"`
	RoleEnd   string `json:"role_end"`
}

type Config struct {
	Indices map[string]Index
}

var config *Config
var once sync.Once

func Get() Config {
	once.Do(loadConfig)
	return *config
}

var loadConfig = func() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Cannot load config file: %s", err)
	}
	temp := struct {
		Indices map[string]struct {
			Page      string `json:"page"`
			Name      string `json:"name"`
			RoleStart string `json:"role_start"`
			RoleEnd   string `json:"role_end"`
		} `json:"indices"`
	}{}
	if err = json.Unmarshal(data, &temp); err != nil {
		log.Fatalf("Config file corrupted: %s", err)
	}
	config = &Config{Indices: make(map[string]Index)}
	for k, v := range temp.Indices {
		(*config).Indices[k] = Index{
			Page:      stringIndexToInt(v.Page),
			Name:      stringIndexToInt(v.Name),
			RoleStart: stringIndexToInt(v.RoleStart),
			RoleEnd:   stringIndexToInt(v.RoleEnd),
		}
	}
}

func stringIndexToInt(index string) int {
	var res int
	for i, r := range index {
		if i == 0 {
			res += int(r - 'A')
		} else {
			res += int(r - 'A' + 1)
		}
	}
	return res
}
