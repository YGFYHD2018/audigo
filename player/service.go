package player

import (
	"sync"
)

type mpp map[*Proxy][]*Player

var (
	playerPool sync.Pool
	proxyPool  sync.Pool
	xProxyMap  mpp
)

func NewProxy() *Proxy {
	proxy := proxyPool.Get().(*Proxy)
	player := playerPool.Get().(*Player)
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

func mapPlayer(proxy *Proxy, players ...*Player) {
	list := make([]*Player, len(players))
	for i, p := range players {
		list[i] = p
	}

	already, ok := xProxyMap[proxy]
	if ok {
		list = append(already, list...)
	}
	xProxyMap[proxy] = list
}
