package player

// Proxy は、sound player proxyです。
type Proxy struct {
	Play   chan PlayArgs
	Stop   chan struct{}
	Volume chan VolumeArgs
	Pause  chan struct{}
	Resume chan struct{}

	closing chan struct{}
}

// NewProxy は、Playerを生成して返します。
func NewProxy() *Proxy {
	return &Proxy{
		Play:   make(chan PlayArgs),
		Stop:   make(chan struct{}),
		Volume: make(chan VolumeArgs),
		Pause:  make(chan struct{}),
		Resume: make(chan struct{}),

		closing: make(chan struct{}),
	}
}
