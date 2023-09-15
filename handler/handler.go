package handler

import (
	"fmt"

	"github.com/abdtyx/JarvisGo/service"
	"github.com/gin-gonic/gin"
)

type HandlerNode struct {
	next      map[byte]*HandlerNode
	handler   func(service.Message)
	fullEqual bool
}

func IsSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

/**
 * route should not contain any leading or trailing space. Otherwise the route won't work
 * fullEqual: if this flag is set, the handler requires the msg to be like "    .help      "
 */
func (node *HandlerNode) Add(route string, handler func(service.Message), fullEqual bool) {
	curNode := node
	for i := 0; i < len(route); i++ {
		if nxtNode, ok := curNode.next[route[i]]; ok {
			curNode = nxtNode
		} else {
			newNode := &HandlerNode{
				next:      make(map[byte]*HandlerNode),
				handler:   nil,
				fullEqual: false,
			}
			curNode = newNode
		}
	}
	curNode.handler = handler
	curNode.fullEqual = fullEqual
}

func (node *HandlerNode) Find(msg service.Message) {
	route := msg.RawMsg
	i := 0
	curNode := node
	// strip leading spaces
	for ; i < len(route); i++ {
		if !IsSpace(route[i]) {
			break
		}
	}

	// Find route
	for ; i < len(route); i++ {
		if IsSpace(route[i]) {
			break
		}
		if nxtNode, ok := curNode.next[route[i]]; ok {
			curNode = nxtNode
		} else {
			return
		}
	}
	if curNode.handler != nil {
		if !curNode.fullEqual {
			curNode.handler(msg)
		} else {
			// fullEqual required
			fullEqual := true
			for ; i < len(route); i++ {
				if !IsSpace(route[i]) {
					// msg not full equal
					fullEqual = false
				}
			}
			if fullEqual {
				curNode.handler(msg)
			}
		}
	}
}

type Handler struct {
	svc      *service.Service
	handlers *HandlerNode
}

func InitHandler() (*Handler, error) {
	var h Handler
	var err error

	h.svc, err = service.InitService()
	if err != nil {
		return nil, err
	}

	fmt.Println("Starting register handlers")

	// make trie
	h.handlers = &HandlerNode{
		next:      make(map[byte]*HandlerNode),
		handler:   nil,
		fullEqual: false,
	}
	// handler registry
	h.handlers.Add(`Jarvis`, h.svc.Jarvis, true)
	h.handlers.Add(`.help`, h.svc.Jhelp, true)
	h.handlers.Add(`.api`, h.svc.Api, true)
	h.handlers.Add(`.Jeminder`, h.svc.Jeminder, true)
	h.handlers.Add(`.blacklist`, h.svc.Blacklist, false)
	h.handlers.Add(`.picrsa`, h.svc.PicRsa, false)
	h.handlers.Add(`.pic`, h.svc.Pic, false)
	h.handlers.Add(`.timetable`, h.svc.Timetable, false)
	h.handlers.Add(`.weather`, h.svc.Weather, false)
	h.handlers.Add(`.suggest`, h.svc.Suggest, false)
	h.handlers.Add(`.log`, h.svc.Jlog, false)
	h.handlers.Add(`.clog`, h.svc.Clog, false)
	h.handlers.Add(`.dlog`, h.svc.Dlog, false)
	h.handlers.Add(`.Hocation`, h.svc.Hocation, false)
	h.handlers.Add(`.usd`, h.svc.Usd, false)

	fmt.Println("Handlers registered")

	// Timed message handler
	go h.svc.TimedMsgHandler()

	return &h, nil
}

func (h *Handler) Handle(c *gin.Context) {
	var msg service.Message
	err := c.BindJSON(&msg)
	if err != nil {
		h.svc.Log.Println("Handler: ", err)
		return
	}

	if msg.MsgType != "" {
		h.MsgHandler(msg)
	}

	// This is a placeholder to make gin happy:)
	c.JSON(200, gin.H{
		"html": "<b>Jarvis! Go!</b>",
	})

	return
}

func (h *Handler) Shutdown() error {
	return h.svc.Shutdown()
}

func (h *Handler) MsgHandler(msg service.Message) {
	// handle group msg not enabled
	if msg.MsgType == "group" && h.svc.Cfg.EnableGroup == false {
		h.svc.Log.Println("Handler: Group message not enabled")
		return
	}

	// handle blacklist
	if userFlag, groupFlag := h.svc.CheckBlacklist(msg); groupFlag {
		return
	} else if userFlag {
		groupTag := ""
		if msg.MsgType == "group" {
			groupTag += fmt.Sprintf("From group %v:", msg.GroupID)
		} else {
			groupTag += "Not from group: "
		}
		h.svc.Log.Println(groupTag + fmt.Sprintf("Sir, a prohibited user %v tried to access my service", msg.UserID))
		return
	}

	// handle msg
	h.handlers.Find(msg)

	return
}
