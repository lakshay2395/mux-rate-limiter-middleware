package limiters

type Limiter interface {
	CanPass() (bool, error)
}
