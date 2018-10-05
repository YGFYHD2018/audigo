package main

import (
	"flag"

	"github.com/code560/audigo/net"
)

func main() {
	flag.Parse()
	port := flag.Arg(0)

	r := net.NewRouter(port)
	r.Run()
}
