package main

import (
	"os"
	"path/filepath"

	"github.com/code560/audigo/app"
	"github.com/code560/audigo/util"
	"github.com/urfave/cli"
)

var log = util.GetLogger()

func main() {
	// reset cd
	execDir, _ := os.Executable()
	log.Debugf("executable: %s\n", execDir)
	execDir = filepath.Dir(execDir)
	os.Chdir(execDir)
	log.Debugf("change current directory: %s\n", execDir)

	cl := cli.NewApp()
	cl.Name = "audigo"
	cl.Usage = "Audio service by LED CUBU"
	cl.Version = "1.0.0"
	cl.Flags = []cli.Flag{
		cli.BoolTFlag{
			Name:  "server, s",
			Usage: "Instant server mode.",
		},
		cli.BoolFlag{
			Name:  "client, c",
			Usage: "Instant client mode.",
		},
		cli.BoolFlag{
			Name:  "repl, r",
			Usage: "Instant REPL mode.",
		},
		cli.StringFlag{
			Name:  "cd",
			Usage: "change current directory by repl",
			Value: "",
		},
	}
	cl.Action = func(ctx *cli.Context) error {
		if ctx.Bool("c") {
			app.Client()
		} else if ctx.Bool("r") {
			dir := ctx.String("cd")
			app.Repl(dir)
		} else if ctx.Bool("s") {
			app.Serve(ctx.Args().Get(0))
		}
		return nil
	}

	cl.Run(os.Args)
}
