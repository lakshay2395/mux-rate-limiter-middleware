package middleware

import (
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/lakshay2395/rate-limiting-algorithms/limiters/leakybucket"
)

func LeakyBucket(client *redis.Client, limit int, expiry time.Duration, identifierCallback func() string) func(http.Handler) http.Handler {
	limiter := leakybucket.NewLeakyBucket(client, limit, expiry, identifierCallback)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok, _ := limiter.CanPass()
			if ok {
				next.ServeHTTP(w, r)
			} else {
				TooManyRequests(w)
			}
		})
	}
}
