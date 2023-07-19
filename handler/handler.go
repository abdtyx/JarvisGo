package handler

import (
	"fmt"

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
	// handle group msg not enabled
	if msg.MsgType == "group" && h.svc.Cfg.EnableGroup == false {
		h.svc.Log.Println("Handler: Group message not enabled")
	}

	// handle blacklist
	if h.svc.CheckBlacklist(msg) {
		groupTag := ""
		if msg.MsgType == "group" {
			groupTag += fmt.Sprintf("From group %v:", msg.GroupID)
		} else {
			groupTag += "Not from group: "
		}
		h.svc.Log.Println(groupTag + fmt.Sprintf("Sir, a prohibited user %v tried to access my service", msg.UserID))
	}

	// handle msg
	switch {
	case msg.RawMsg == "Jarvis":
		h.svc.Jarvis(msg)
	case msg.RawMsg == ".help":
		h.svc.Jhelp(msg)
	case msg.RawMsg == ".api":
		h.svc.Api(msg)
	}

	return
}
