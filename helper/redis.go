package helper

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"kamari.com/config"
	"time"
)

type RedisHelper struct {
	Client *redis.Client
}

/*
instanceNumber 0 : OTP
instanceNumber 1 : step up token
instanceNumber 2 : refresh token
instansceNumber 3: add bank account
instansceNumber 4: Top Up Wallet
instanceNumber 5: something
instanceNumber 6: Pay SLP
*/
func NewRedisClient(instanceNumber int) RedisHelper {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisConfig.Host,
		Password: config.Config.RedisConfig.Password,
		DB:       instanceNumber,
	})
	return RedisHelper{Client: client}
}

func (redis RedisHelper) SetData(key string, data string) error {
	ttl := time.Duration(5) * time.Minute
	op1 := redis.Client.Set(context.Background(), key, data, ttl)
	return op1.Err()
}

func (redis RedisHelper) SetData7Days(key string, data string) error {
	ttl := time.Duration(24) * time.Hour * 7
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

func (redis RedisHelper) GenerateKey(objectType string, uniqueKey string, field string) string {
	return fmt.Sprintf("%s:%s:%s", objectType, uniqueKey, field)
}

func (redis RedisHelper) Set(key string, data interface{}) error {
	ttl := time.Duration(5) * time.Minute
	op3 := redis.Client.Set(context.Background(), key, data, ttl)
	return op3.Err()
}
