package player

import (
	"sync"
)

type mpp map[*Proxy][]*player

var (
	playerPool sync.Pool
	proxyPool  sync.Pool
	xProxyMap  mpp
)

func NewProxy() *Proxy {
	proxy := proxyPool.Get().(*Proxy)
	player := playerPool.Get().(*player)
	mapPlayer(proxy, player)
	return proxy
}

func init() {
	initMap()
	initPool()
}

func initMap() {
	xProxyMap = make(mpp, 50)
}

func initPool() {
	playerPool = sync.Pool{
		New: func() interface{} {
			return newPlayerImpl() // todo trans internal
		},
	}

	proxyPool = sync.Pool{
		New: func() interface{} {
			return newProxyImpl() // todo trans internal
		},
	}
}

func mapPlayer(proxy *Proxy, players ...*player) {
	list := make([]*player, len(players))
	for i, p := range players {
		list[i] = p
	}

	already, ok := xProxyMap[proxy]
	if ok {
		list = append(already, list...)
	}
	xProxyMap[proxy] = list
}
