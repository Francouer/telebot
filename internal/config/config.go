package config

import (
	"encoding/json"
	"io/ioutil"
	"telebot/telebot/CA/pkg/infrastructure/config"
)

type Config struct {
	Token string          `json:"token"`
	DB    config.DBconfig `json:"db"`
}

func New(path string) (*Config, error) {
	p, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)

	err = json.Unmarshal(p, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
