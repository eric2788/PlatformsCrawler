package crawling

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/eric2788/PlatformsCrawler/file"
	"github.com/go-redis/redis/v8"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	cli      *redis.Client
	mu       = &sync.Mutex{}
	notFound = mapset.NewSet()
)

// InitRedis also for test
func InitRedis() {

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

func GetString(key string) (string, error) {
	if cli == nil {
		return "", fmt.Errorf("redis client does not initalized")
	}
	s, err := cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else {
		return s, err
	}
}

func SetString(key, value string) error {
	if cli == nil {
		return fmt.Errorf("redis client does not initalized")
	}
	return cli.Set(ctx, key, value, 0).Err()
}

func SetStringTemp(key, value string, duration time.Duration) error {
	if cli == nil {
		return fmt.Errorf("redis client does not initalized")
	}
	return cli.Set(ctx, key, value, duration).Err()
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

func ClearMap(key string, fields ...string) (int64, error) {
	if cli == nil {
		return -1, fmt.Errorf("redis client does not initalized")
	}

	if len(fields) > 0 {
		return cli.HDel(ctx, key, fields...).Result()
	} else {
		return cli.Del(ctx, key).Result()
	}
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

func LoadWithCache(key string, fetch func(toFetch []string) (map[string]string, error), isNotFound func(err error) bool, fields ...string) (map[string]string, error) {
	var cache map[string]string

	if notInCaches, err := GetMap(key, &cache, fields...); err != nil {
		return nil, err
	} else if len(notInCaches) == 0 { // all names are in the cache
		return cache, nil
	} else { // some names can't get

		toFetch := make([]string, 0)

		// 透過快取作過濾
		for _, name := range notInCaches {
			if notFound.Contains(name) {
				continue
			}
			toFetch = append(toFetch, name)
		}

		if len(toFetch) == 0 {
			return cache, nil
		}

		fetched, err := fetch(toFetch)

		if err != nil {

			if isNotFound(err) {
				logger.Warnf("查無這些用戶: %v", strings.Join(toFetch, ", "))
				for _, screen := range toFetch {
					notFound.Add(screen)
				}
			} else if fetched == nil {
				return nil, err
			} else {
				logger.Errorf("查找用戶 %s 時出現錯誤: %v", strings.Join(toFetch, ", "), err)
			}
		}

		if fetched != nil {

			if err := SaveMap(key, fetched); err != nil {
				logger.Warnf("Redis 儲存 用戶資訊 時 出現錯誤: %v", err)
				logger.Warnf("儲存的資訊: %v", fetched)
			}

			for screen, id := range fetched {
				cache[screen] = id
			}

			//把得出的結果加到 set 來檢查某些沒有加到結果的 screenNames
			nameSet := mapset.NewSet()

			for screen := range cache {
				nameSet.Add(screen)
			}

			// 然後加到 notFound 中作為例外
			for _, name := range fields {
				if !nameSet.Contains(name) {
					notFound.Add(name)
				}
			}

		}

		return cache, nil
	}
}
