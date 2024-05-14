package lib

type Awaiter struct {
	ready chan struct{}
}

func NewAwaiter() *Awaiter {
	return &Awaiter{ready: make(chan struct{})}
}

func (a *Awaiter) Await() <-chan struct{} {
	return a.ready
}

func (a *Awaiter) Done() {
	close(a.ready)
}
