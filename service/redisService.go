package service

import (
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

type redisInfo struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

var redisConfig struct {
	Redis redisInfo `yaml:"redis"`
}

var RedisClient *redis.Client

func InitRedis() {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	//config := redisConfig{}

	err = yaml.Unmarshal(configFile, &redisConfig)
	if err != nil {
		panic(err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Redis.Addr,     // Redis服务器地址
		Password: redisConfig.Redis.Password, // Redis密码
		DB:       redisConfig.Redis.DB,       // Redis数据库索引
	})
}

func IsTokenExist(token string) bool {
	result, err := RedisClient.Exists(token).Result()
	if err != nil {
		panic(err)
	}
	if result == 0 {
		return false
	}
	return true
}

func GetIdByToken(token string) int64 {
	result, err := RedisClient.Get(token).Result()
	if err != nil {
		panic(err)
	}
	result64, _ := strconv.ParseInt(result, 10, 64)
	return result64
}
