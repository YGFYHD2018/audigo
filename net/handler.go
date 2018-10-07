package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/code560/audigo/player"
	"github.com/gin-gonic/gin"
)

var (
	once sync.Once
	h    = newHandler()
)

const (
	INIT_PLAYER_COUNT = 10
)

type handler struct {
	players map[string]*player.Proxy
}

// SetHandler は、ginにapiハンドラーを設定します。
func SetHandler(r *gin.Engine) {
	setV1(r)
}

func setV1(r *gin.Engine) {
	v1 := r.Group("audigo/v1")
	{
		v1.GET("/players", h.ids)
		v1.POST("/init:content_id", h.create)
		v1.POST("/play:content_id", h.play)
		v1.POST("/stop:content_id", h.stop)
		v1.POST("/volume:content_id", h.volume)
		v1.POST("/pause:content_id", h.pause)
		v1.POST("/resume:content_id", h.resume)
	}
}

func newHandler() *handler {
	var inst *handler
	once.Do(func() {
		inst = &handler{
			players: make(map[string]*player.Proxy, INIT_PLAYER_COUNT), // TODO
		}
	})
	return inst
}

func (h *handler) create(c *gin.Context) {
	code := http.StatusNoContent
	id := c.Param("content_id")
	_, ok := h.players[id]
	if !ok {
		h.players[id] = player.NewProxy()
		code = http.StatusCreated
	}
	c.JSON(code, nil)
}

func (h *handler) ids(c *gin.Context) {
	code := http.StatusBadRequest
	keys := make([]string, len(h.players))
	for k := range h.players {
		keys = append(keys, k)
	}
	res, err := json.Marshal(keys)
	if err != nil {
		code = http.StatusInternalServerError
		c.JSON(code, err)
		return
	}
	c.JSON(code, res)
}

func (h *handler) play(c *gin.Context) {
	code := http.StatusAccepted
	p, err := h.getPlayer(c)
	if err != nil {
		return
	}
	args := player.PlayArgs{}
	if err := c.ShouldBindJSON(args); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	p.Play <- args
	c.JSON(code, nil)
}

func (h *handler) stop(c *gin.Context) {
	code := http.StatusAccepted
	p, err := h.getPlayer(c)
	if err != nil {
		return
	}
	p.Stop <- struct{}{}
	c.JSON(code, nil)
}

func (h *handler) volume(c *gin.Context) {
	code := http.StatusAccepted
	p, err := h.getPlayer(c)
	if err != nil {
		return
	}
	args := player.VolumeArgs{}
	if err := c.ShouldBindJSON(args); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	p.Volume <- args
	c.JSON(code, nil)
}

func (h *handler) pause(c *gin.Context) {
	code := http.StatusAccepted
	p, err := h.getPlayer(c)
	if err != nil {
		return
	}
	p.Pause <- struct{}{}
	c.JSON(code, nil)
}

func (h *handler) resume(c *gin.Context) {
	code := http.StatusAccepted
	p, err := h.getPlayer(c)
	if err != nil {
		return
	}
	p.Resume <- struct{}{}
	c.JSON(code, nil)
}

func (h *handler) getPlayer(c *gin.Context) (*player.Proxy, error) {
	id := c.Param("content_id")
	p, ok := h.players[id]
	if !ok {
		c.JSON(http.StatusNotFound, id)
		return nil, fmt.Errorf("not found id player: %s", id)
	}
	return p, nil
}
