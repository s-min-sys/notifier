package config

import (
	"sync"

	"github.com/GizmoVault/gotools/configx"
)

type Config struct {
	Listens string `yaml:"Listens"`

	Senders map[string]string `yaml:"Senders"`
}

var (
	_config Config
	_once   sync.Once
)

func GetConfig() *Config {
	_once.Do(func() {
		_, err := configx.Load("config.yaml", &_config)
		if err != nil {
			panic(err)
		}
	})

	return &_config
}
