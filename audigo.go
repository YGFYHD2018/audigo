package main

import (
	"os"
	"path"

	"github.com/code560/audigo/app"
)

var execDir = path.Dir(os.Args[0])

func main() {
	os.Chdir(execDir)

	// app.Serve()
	app.SoundPlay()
}
