package config

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/sgostarter/libconfig"
	"github.com/sgostarter/libeasygo/stg/redisex"
)

type TeleConfig struct {
	Token       string `yaml:"Token"`
	APIEndPoint string `yaml:"APIEndPoint"`
}

type Config struct {
	TeleConfig *TeleConfig `yaml:"TeleConfig"`

	Listen   string `yaml:"Listen"`
	RedisDNS string `yaml:"RedisDNS"`

	//
	//
	//

	RedisCli *redis.Client `yaml:"-"`
}

var (
	_config Config
	_once   sync.Once
)

func GetConfig() *Config {
	_once.Do(func() {
		_, err := libconfig.Load("config.yaml", &_config)
		if err != nil {
			panic(err)
		}

		_config.RedisCli, err = redisex.InitRedis(_config.RedisDNS)
		if err != nil {
			panic(err)
		}
	})

	return &_config
}
