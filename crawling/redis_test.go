package crawling

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"testing"
	"time"
)

var cache = make(map[string]string)

func initRedisTest() {
	cli = redis.NewClient(&redis.Options{
		Addr:     "192.168.0.127:6379",
		DB:       1,
		Password: "",
	})
}

func MapIO(t *testing.T) {

	initRedisTest()

	key := "test:test_map"

	rand.Seed(time.Now().UnixMicro())

	values := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O"}

	k, v := fmt.Sprintf("%v", rand.Int()), values[rand.Intn(len(values))]

	fmt.Printf("assigning %v to %v from cache \n", k, v)

	cache[k] = v

	fmt.Println(cache)

	if err := SaveMap(key, cache); err != nil {
		t.Fatal(err)
	}

	var take map[string]string

	if err := GetAllMap(key, &take); err != nil && err != redis.Nil {
		t.Fatal(err)
	}

	fmt.Println(take)
}

func MapMultiKey(t *testing.T) {
	initRedisTest()

	key := "test:test_map"

	insert := map[string]string{
		"A": "123",
		"B": "456",
		"C": "789",
	}

	if err := SaveMap(key, insert); err != nil {
		t.Fatal(err)
	}

	var take map[string]string

	if er, err := GetMap(key, &take, "A", "B", "C", "D"); err != nil {
		t.Fatal(err)
	} else if len(er) > 0 {
		fmt.Println("can't get: ", er)
	}

	fmt.Println(take)

}
