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
			log.Errorf("recorvered from %s", r)
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
	// log := GetLogger()
	muClosing.Lock()
	defer muClosing.Unlock()
	c := &Closing{
		Done: make(chan bool, 1),
	}
	// log.Debug("make closing: %p", c)
	// log.Debug("make chan(Done): %p", c.Done)
	return c
}

func (c *Closing) Close() {
	// log := GetLogger()
	muClosing.Lock()
	defer muClosing.Unlock()

	// log.Debug("closing: %p", c)
	// log.Debug("chan(Done): %p", c.Done)
	if !isDone(c.Done) {
		close(c.Done)
	}

	// recove(func() error {
	// 	log.Debug("make chan(Done): %p", c.Done)
	// 	if !isDone(c.Done) {
	// 		close(c.Done)
	// 	}
	// 	return nil
	// }) // TODO: 握りつぶし、そのうち対応
}

func (c *Closing) GetDone() <-chan bool {
	return c.Done
}

func isDone(c chan bool) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}
