package leakybucket

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/lakshay2395/rate-limiting-algorithms/limiters"
)

type limiter struct {
	client             *redis.Client
	expiry             time.Duration
	limit              int
	identifierCallback func() string
}

func NewLeakyBucket(client *redis.Client, limit int, expiry time.Duration, identifierCallback func() string) limiters.Limiter {
	return limiter{
		client:             client,
		expiry:             expiry,
		limit:              limit,
		identifierCallback: identifierCallback,
	}
}

func (l limiter) CanPass() (bool, error) {
	identifier := l.identifierCallback()
	currentMinute := getCurrentMinute()
	value, err := l.client.Get(identifier + currentMinute).Result()
	count := 0
	if err == redis.Nil {
		count = 0
	} else {
		count, _ = strconv.Atoi(value)
	}
	if count < l.limit {
		pipe := l.client.Pipeline()
		pipe.Incr(identifier + currentMinute)
		pipe.Expire(identifier+currentMinute, l.expiry)
		pipe.Exec()
		return true, nil
	}
	return false, nil
}

func getCurrentMinute() string {
	return strconv.FormatInt(time.Now().Unix()/60, 10)
}
