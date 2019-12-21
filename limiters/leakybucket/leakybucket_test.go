package leakybucket

import (
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
)

func TestLimiter(t *testing.T) {
	ch := make(chan int)
	for i := 0; i < 3; i++ {
		go (func(i int) {
			ticker := time.NewTicker(100 * time.Millisecond)
			for ; true; <-ticker.C {
				redisClient, _ := getRedisClient()
				limiter := limiter{client: redisClient, expiry: 10 * time.Second, limit: 10, identifierCallback: func() string {
					return "U2"
				}}
				result, err := limiter.CanPass()
				if err != nil {
					log.Fatal("Error Occurred : ", err)
				}
				redisClient.Close()
				log.Println((i + 1), " : user = ", limiter.identifierCallback(), " can pass = ", result)
			}
		})(i)
	}
	<-ch
}

func getRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
