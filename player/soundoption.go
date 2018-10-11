package player

type PlayArgs struct {
	Src  string `json:"src"`
	Loop bool   `json:"loop"`
	Stop bool   `json:"stop"`
}

type VolumeArgs struct {
	val float64
}

const (
	dir = "assets/audio/"
)

type optPlay func(*PlayArgs)

func NewPlayArgs(opts ...optPlay) *PlayArgs {
	u := &PlayArgs{Src: "", Loop: false, Stop: false}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func Wave(name string) optPlay {
	return func(a *PlayArgs) {
		a.Src = dir + name
	}
}

func Loop(enable bool) optPlay {
	return func(a *PlayArgs) {
		a.Loop = enable
	}
}

func Stop(enable bool) optPlay {
	return func(a *PlayArgs) {
		a.Stop = enable
	}
}
