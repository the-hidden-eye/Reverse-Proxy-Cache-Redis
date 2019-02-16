package main

import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"strconv"
	"time"
	"net/http"
)

type Cache struct {
}

var redisClient *redis.Client
var redisClient2 *redis.Client
var redisClient3 *redis.Client

// CreateCache - Create Redis Connection
func CreateCache() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
		Password: "",
		DB:       0,
	})
	redisClient2 = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
		Password: "",
		DB:       1,
	})
	redisClient3 = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
		Password: "",
		DB:       2,
	})
}

func (c *Cache) has(key string) bool {
	total, err := redisClient.Exists(key).Result()
	if err != nil {
		panic(err)
	}
	if total > 0 {
		return true
	}
	return false
}

func (c *Cache) get(key string) ([]byte, map[string]string, error) {
	value, err := redisClient.Get(key).Result()
	headers, _ := redisClient3.HGetAll(key).Result()
	return []byte(value), headers, err
}

func (c *Cache) put(key string, header http.Header, content []byte) error {
	err := redisClient.Set(key, content, 0).Err()
	for k, vs := range header {
		for _, v := range vs {
			redisClient3.HSet(key, k, v)
		}
	}
	if err != nil {
		return err
	}
	setExpiration(key)
	return nil
}

func setExpiration(key string) {
	exp, _ := strconv.Atoi(getEnv("CACHE_EXP", "60"))
	err := redisClient2.Set(key, key, time.Duration(exp)*time.Second).Err()
	if err != nil {
		Error.Printf("Occured an error on set expiration cache: %s", key)
	}
	return
}

func (c *Cache) checkUpdate(key string) bool {
	total, err := redisClient2.Exists(key).Result()
	if err != nil {
		Error.Printf("Occured an erro on check valid expiration cache: %s", key)
		return false
	}
	if total > 0 {
		return false
	}
	Info.Printf("Cache %s is expired, recreating", key)
	go recreate(key)
	return true
}

func recreate(key string) {
	Info.Printf("Updating cache from %s", key)
	setExpiration(key)
	headers, _ := redisClient3.HGetAll(key).Result()
	response, err := gatewayRequestUpdate(key, headers)
	body, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		Error.Printf("Occured error on update cache: %s - Error: %v", key, err)
		return
	}
	err = cache.put(key, response.Header, body)
	if err != nil {
		Error.Printf("Occured error on update cache: %s - Error: %v", key, err)
		return
	}
}
