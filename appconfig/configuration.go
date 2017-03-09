package appconfig

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ClientIDPath string `json:"clientIDPath"`
	PivotalToken string `json:"pivotalToken"`
	ProjectID    string `json:projectID"`
}

func NewConfig() *Config {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}

	config := &Config{}
	json.Unmarshal(data, config)
	return config
}
