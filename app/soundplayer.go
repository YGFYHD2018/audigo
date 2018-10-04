package main

import (
	"flag"
	"math"
	"os"
	"sync"
	"time"

	"github.com/code560/audigo"
	"github.com/code560/audigo/util"
)

var l = util.GetLogger()

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		return
	}

	wg := sync.WaitGroup{}

	plist := make([]audigo.Player, len(os.Args))
	for _, arg := range flag.Args() {
		wg.Add(1)
		p := audigo.NewPlayer()
		plist = append(plist, *p)
		go func(p *audigo.Player, name string) {
			defer wg.Done()
			p.Play(name, math.MaxInt64)
		}(p, arg)
	}

	l.Debug("started wav")
	time.Sleep(time.Second * 2)
	for _, p := range plist[0:] {
		l.Debug("call stop")
		p.Stop()
	}

	wg.Wait()
	l.Debug("done routines")
}
