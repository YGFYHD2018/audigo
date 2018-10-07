package player

type PlayArgs struct {
	wav  string
	loop bool
	stop bool
}

type VolumeArgs struct {
	val float64
}

const (
	dir = "assets/audio/"
	ext = "wav"
)

type optPlay func(*PlayArgs)

func NewPlayArgs(opts ...optPlay) *PlayArgs {
	u := &PlayArgs{wav: "", loop: false, stop: false}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func Wave(name string) optPlay {
	return func(a *PlayArgs) {
		a.wav = dir + name
	}
}

func Loop(enable bool) optPlay {
	return func(a *PlayArgs) {
		a.loop = enable
	}
}

func Stop(enable bool) optPlay {
	return func(a *PlayArgs) {
		a.stop = enable
	}
}
