package player

import "testing"

func TestOption(t *testing.T) {
	msg := "何かおかしい newPlayArgs()の"
	a := newPlayArgs(Wave("bgm_wave.wav"), Loop(true), Stop(true))
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

	a = newPlayArgs(Wave("foo"))
	if a.wav != dir+"foo" {
		t.Error(msg + "Wave only")
	}
	a = newPlayArgs(Loop(true))
	if !a.loop {
		t.Error(msg + "Loop only")
	}
	a = newPlayArgs(Stop(true))
	if !a.stop {
		t.Error(msg + "Stop only")
	}
}
