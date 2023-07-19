package handler

import (
	"github.com/abdtyx/JarvisGo/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func InitHandler() (*Handler, error) {
	var h Handler
	var err error

	h.svc, err = service.InitService()
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *Handler) Handle(c *gin.Context) {
	var msg service.Message
	err := c.BindJSON(&msg)
	if err != nil {
		h.svc.Log.Println("Handler: ", err)
	}

	h.MsgHandler(msg)

	// This is a placeholder to make gin happy:)
	c.JSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})

	return
}

func (h *Handler) MsgHandler(msg service.Message) {
	// handle blacklist

	// handle msg
	switch {
	case msg.RawMsg == "Jarvis":
		h.svc.Jarvis(msg)
	case msg.RawMsg == ".help":
		h.svc.Jhelp(msg)
	}
	return
}
