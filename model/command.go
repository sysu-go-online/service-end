package model

import (
	"fmt"

	"github.com/go-redis/redis"
)

// IsUUIDExists judge if the uuid exists in redis
func IsUUIDExists(name string, client *redis.Client) (bool, error) {
	_, err := client.Get(name).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
