package player

type Player interface {
	Play(args *PlayArgs)
	Stop(callback chan<- struct{})
	Volume(args *VolumeArgs)
	Pause()
	Resume()
}

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

type Proxy interface {
	GetChannel() chan<- *Action
}

type Action struct {
	Act  Actions
	Args interface{}
}

type Actions int

const (
	_ Actions = iota
	Play
	Stop
	Volume
	Pause
	Resume
)
