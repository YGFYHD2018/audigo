package app

import (
	"github.com/code560/audigo/net"
)

func Serve(port string) {
	r := net.NewRouter()
	r.Run(port)
}
