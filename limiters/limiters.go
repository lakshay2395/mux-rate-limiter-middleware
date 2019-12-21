package limiters

import "net/http"

type Limiter interface {
	CanPass(r *http.Request) (bool, error)
}
