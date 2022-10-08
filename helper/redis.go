package helper

import (
	"context"
	"github.com/go-redis/redis/v9"
	"kamari.com/config"
	"time"
)

type RedisHelper struct {
	Client *redis.Client
}

func NewRedisClient() RedisHelper {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisConfig.Host,
		Password: config.Config.RedisConfig.Password,
		DB:       0,
	})
	return RedisHelper{Client: client}
}

func (redis RedisHelper) SetData(key string, data string) error {
	ttl := time.Duration(5) * time.Minute
	op1 := redis.Client.Set(context.Background(), key, data, ttl)
	return op1.Err()
}

func (redis RedisHelper) GetData(key string) (string, error) {
	op2 := redis.Client.Get(context.Background(), key)
	if err := op2.Err(); err != nil {
		return "", err
	}
	res, err := op2.Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (redis RedisHelper) Delete(key string) error {
	op := redis.Client.Del(context.Background(), key)
	return op.Err()
}
