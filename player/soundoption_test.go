package player

import "testing"

func TestOption(t *testing.T) {
	msg := "何かおかしい newPlayArgs()の"
	a := NewPlayArgs(Wave("bgm_wave.wav"), Loop(true), Stop(true))
	if a.wav != dir+"bgm_wave.wav" {
		t.Error(msg + "Wave")
		t.Errorf("set: %s, get: %s", "bgm_wave.wav", a.wav)
	}
	if !a.loop {
		t.Error(msg + "Loop")
	}
	if !a.stop {
		t.Error(msg + "Stop")
	}

	a = NewPlayArgs(Wave("foo"))
	if a.wav != dir+"foo" {
		t.Error(msg + "Wave only")
	}
	a = NewPlayArgs(Loop(true))
	if !a.loop {
		t.Error(msg + "Loop only")
	}
	a = NewPlayArgs(Stop(true))
	if !a.stop {
		t.Error(msg + "Stop only")
	}
}
