package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"os"
)

type redisInfo struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

var redisConfig struct {
	Redis redisInfo `yaml:"redis"`
}

var Client *redis.Client

func InitRedis() {
	configFile, err := os.ReadFile("RedisConfig.yaml")
	if err != nil {
		panic(err)
	}

	//config := redisConfig{}

	err = yaml.Unmarshal(configFile, &redisConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Redis配置%+v", redisConfig)

	Client = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Redis.Addr,     // Redis服务器地址
		Password: redisConfig.Redis.Password, // Redis密码
		DB:       redisConfig.Redis.DB,       // Redis数据库索引
	})
}
