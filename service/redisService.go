package service

import (
	"encoding/json"
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
	//fmt.Println("鉴权token", token)
	result, err := RedisClient.Exists(token).Result()
	if err != nil {
		panic(err)
	}
	if result == 0 {
		return false
	}
	return true
}

type resultStruct struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func GetNameByToken(token string) string {
	result, err := RedisClient.Get(token).Result()
	if err != nil {
		if err == redis.Nil {
			// Handle the case where the token does not exist in the cache
			fmt.Println("Token not found in cache:", token)
			return ""
		}
		panic(err)
	}

	var res resultStruct
	if err := json.Unmarshal([]byte(result), &res); err != nil {
		panic(err)
	}
	return res.Name
}

func GetIdByToken(token string) int64 {
	result, err := RedisClient.Get(token).Result()
	if err != nil {
		if err == redis.Nil {
			// Handle the case where the token does not exist in the cache
			fmt.Println("Token not found in cache:", token)
			return 0
		}
		panic(err)
	}

	var res resultStruct
	if err := json.Unmarshal([]byte(result), &res); err != nil {
		panic(err)
	}
	return res.ID
}
