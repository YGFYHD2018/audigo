package util

import (
	"fmt"
	"sync"
)

func recoveGo(f func() error) error {
	e := make(chan error, 1)
	go func() {
		if err := recove(f); err != nil {
			e <- err
		}
		e <- nil
	}()
	return <-e
}

func recove(f func() error) error {
	log := GetLogger()
	defer func() error {
		if r := recover(); r != nil {
			log.Error("recorvered from ", r)
			return fmt.Errorf("error %s", r)
		}
		return nil
	}()
	return f()
}

var muClosing sync.Mutex

type Closing struct {
	Done chan bool
}

func NewClosing() *Closing {
	muClosing.Lock()
	defer muClosing.Unlock()
	c := &Closing{
		Done: make(chan bool),
	}
	return c
}

func (c *Closing) Reset() {
	c.Close()
	c.Done = make(chan bool)
}

func (c *Closing) Close() {
	muClosing.Lock()
	defer muClosing.Unlock()

	if !IsDone(c.Done) {
		close(c.Done)
	}
}

func (c *Closing) GetDone() <-chan bool {
	return c.Done
}

func IsDone(c chan bool) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}
