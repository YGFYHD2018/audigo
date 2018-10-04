package main

import (
	"flag"
	"math"
	"os"
	"runtime/trace"
	"time"

	"github.com/code560/audigo"
	"github.com/code560/audigo/util"
)

var l = util.GetLogger()

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return
	}

	// trace
	f, err := os.Create("trace.out")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	// trace end

	plist := make([]*audigo.Player, len(args))
	for i, arg := range args {
		p := audigo.NewPlayer()
		// l.Debug("make player: %p", p)
		// plist = append(plist, p)
		plist[i] = p

		go func(p *audigo.Player, name string) {
			p.Play(name, math.MaxInt64)
		}(p, arg)
	}

	// l.Debug("started wav")
	time.Sleep(time.Second * 2)
	for _, p := range plist {
		// l.Debug("call stop")
		// l.Debug("player: %p", p)
		p.Stop()
	}

	// l.Debug("done routines")
}
