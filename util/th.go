package util

import (
	"fmt"
)

func recorveGo(f func() error) error {
	e := make(chan error, 1)
	go func() {
		if err := recorve(f); err != nil {
			e <- err
		}
		e <- nil
	}()
	return <-e
}

func recorve(f func() error) error {
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
