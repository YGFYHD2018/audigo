package player

import "testing"

func TestOption(t *testing.T) {
	msg := "何かおかしい newPlayArgs()の"
	a := NewPlayArgs(Src("bgm_wave.wav"), Loop(true), Stop(true))
	if a.Src != dir+"bgm_wave.wav" {
		t.Error(msg + "Wave")
		t.Errorf("set: %s, get: %s", "bgm_wave.wav", a.Src)
	}
	if !a.Loop {
		t.Error(msg + "Loop")
	}
	if !a.Stop {
		t.Error(msg + "Stop")
	}

	a = NewPlayArgs(Src("foo"))
	if a.Src != dir+"foo" {
		t.Error(msg + "Wave only")
	}
	a = NewPlayArgs(Loop(true))
	if !a.Loop {
		t.Error(msg + "Loop only")
	}
	a = NewPlayArgs(Stop(true))
	if !a.Stop {
		t.Error(msg + "Stop only")
	}
}
