package services

import (
	"Asgard/constants"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/logger"
)

func GetCache(key string) string {
	redis := cache.Get(constants.CACHE_NAME)
	data, err := redis.Get(key).Result()
	if err != nil && err.Error() != "redis: nil" {
		logger.Error("GetCache Error:", err, key)
	}
	return data
}

func SetCache(key, data string) {
	redis := cache.Get(constants.CACHE_NAME)
	err := redis.SetNX(key, data, constants.CACHE_TTL).Err()
	if err != nil {
		logger.Error("SetCache Error:", err, key)
	}
}

func DelCache(key string) {
	redis := cache.Get(constants.CACHE_NAME)
	err := redis.Del(key).Err()
	if err != nil {
		logger.Error("DelCache Error:", err, key)
	}
}
