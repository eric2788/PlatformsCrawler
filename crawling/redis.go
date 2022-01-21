package crawling

import (
	"fmt"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/go-redis/redis/v8"
	"reflect"
	"sync"
)

var (
	cli *redis.Client
	mu  = &sync.Mutex{}
)

func initRedis() {

	rConfig := file.ApplicationYaml.Redis

	cli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rConfig.Host, rConfig.Port),
		Password: rConfig.Password,
		DB:       rConfig.Database,
	})
}

func DoRedis(do func(client *redis.Client) error) error {
	if cli == nil {
		return fmt.Errorf("redis client does not initalized")
	}
	mu.Lock()
	defer mu.Unlock()
	return do(cli)
}

func SaveMap(key string, dict interface{}) error {
	if reflect.TypeOf(dict).Kind() != reflect.Map {
		return fmt.Errorf("the dict value should be a map")
	}
	return DoRedis(func(cli *redis.Client) error {
		return cli.HSet(ctx, key, dict).Err()
	})
}

func GetAllMap(key string, dict *map[string]string) error {
	if cli == nil {
		return fmt.Errorf("redis client does not initalized")
	}
	d, err := cli.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}
	*dict = d
	return nil
}

func GetMap(key string, dict *map[string]string, fields ...string) ([]string, error) {

	var errorField []string

	if cli == nil {
		return errorField, fmt.Errorf("redis client does not initalized")
	}

	values := make(map[string]string)

	for _, field := range fields {
		value, err := cli.HGet(ctx, key, field).Result()
		if err != nil {
			if err != redis.Nil {
				logger.Warnf("Redis 獲取 %s 中的 %s 值時出現錯誤: %v", key, field, err)
				return nil, err
			}
			errorField = append(errorField, field)
		} else {
			values[field] = value
		}
	}

	*dict = values
	return errorField, nil
}
