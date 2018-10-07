package app

import (
	"flag"

	"github.com/code560/audigo/net"
)

func Serve() {
	flag.Parse()
	port := flag.Arg(0)

	r := net.NewRouter()
	r.Run(port)
}
