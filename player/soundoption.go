package player

type PlayArgs struct {
	Src  string `json:"src"`
	Loop bool   `json:"loop"`
	Stop bool   `json:"stop"`
}

type VolumeArgs struct {
	Vol float64 `json:"vol"`
}

const (
	dir = "asset/audio/"
)

type OptPlay func(*PlayArgs)

func NewPlayArgs(opts ...OptPlay) *PlayArgs {
	u := &PlayArgs{Src: "", Loop: false, Stop: false}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func Src(name string) OptPlay {
	return func(a *PlayArgs) {
		a.Src = dir + name
	}
}

func Loop(enable bool) OptPlay {
	return func(a *PlayArgs) {
		a.Loop = enable
	}
}

func Stop(enable bool) OptPlay {
	return func(a *PlayArgs) {
		a.Stop = enable
	}
}
