package net

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Router interface
type Router interface {
	Run(port string)
}

type router struct {
	router *gin.Engine
}

// NewRouter は、Routerを作成して返す。
func NewRouter() Router {
	r := &router{
		router: gin.Default(),
	}
	SetHandler(r.router)
	return r
}

func (r *router) Run(port string) {
	r.Run(fmt.Sprintf(":%s", port))
}
