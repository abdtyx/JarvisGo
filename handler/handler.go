package handler

import (
	"log"

	"github.com/abdtyx/JarvisGo/service"
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	var msg service.Message
	l := log.Default()
	err := c.BindJSON(&msg)
	if err != nil {
		l.Println(err)
	}

	MsgHandler(msg)

	return
}

func MsgHandler(msg service.Message) {
	l := log.Default()
	switch {
	case msg.RawMsg == "Jarvis":
		err := service.Jarvis(msg)
		if err != nil {
			l.Println(err)
		}
	}
	return
}
