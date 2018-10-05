package player

type playArgs struct {
	wav  string
	loop bool
	stop bool
}

type volumeArgs struct {
	val float64
}

// Proxy は、sound player proxyです。
type Proxy struct {
	closing chan struct{}

	play   chan playArgs
	stop   chan struct{}
	volume chan volumeArgs
	pause  chan struct{}
	resume chan struct{}
}

// NewProxy は、Playerを生成して返します。
func NewProxy() *Proxy {
	return &Proxy{
		closing: make(chan struct{}),
		play:    make(chan playArgs),
		stop:    make(chan struct{}),
		volume:  make(chan volumeArgs),
		pause:   make(chan struct{}),
		resume:  make(chan struct{}),
	}
}
