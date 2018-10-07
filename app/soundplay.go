package app

import (
	"flag"
	"os"
	"runtime/trace"
	"sync"
	"time"

	"github.com/code560/audigo/player"
	"github.com/code560/audigo/util"
)

var log = util.GetLogger()

func SoundPlay() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return
	}

	// trace
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	// trace end

	w := sync.WaitGroup{}
	_ = playFiles(args, &w)
	w.Wait()
}

func playFiles(files []string, w *sync.WaitGroup) []*player.Player {
	plist := make([]*player.Player, len(files))
	for i, arg := range files {
		p := player.NewPlayer()
		plist[i] = p

		w.Add(1)
		go func(p *player.Player, name string) {
			p.Play(player.NewPlayArgs(
				player.Wave(name),
				player.Loop(false)))
			w.Done()
		}(p, arg)
	}
	return plist
}

func stopPlayers(plist []*player.Player) {
	time.Sleep(time.Second * 2)
	for _, p := range plist {
		p.Stop()
	}
}
