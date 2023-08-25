package handler

import (
	"fmt"
	"regexp"

	"github.com/abdtyx/JarvisGo/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

var (
	regBlacklist = regexp.MustCompile(`^\.blacklist`)
	regPicRSA    = regexp.MustCompile(`^\.pic rsa`)
	regPic       = regexp.MustCompile(`^\.pic`)
	regTimetable = regexp.MustCompile(`^\.timetable`)
	regWeather   = regexp.MustCompile(`^\.weather`)
	regSuggest   = regexp.MustCompile(`^\.suggest`)
	regLog       = regexp.MustCompile(`^\.log`)
	regClog      = regexp.MustCompile(`^\.clog`)
	regDlog      = regexp.MustCompile(`^\.dlog`)
	regHocation  = regexp.MustCompile(`^\.Hocation`)
	regUsd       = regexp.MustCompile(`^\.usd`)
)

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
		return
	}

	if msg.MsgType == "" {
		// Heartbeat event
		service.TimedMsgHandler()
	} else {
		// Message event
		h.MsgHandler(msg)
	}

	// This is a placeholder to make gin happy:)
	c.JSON(200, gin.H{
		"html": "<b>Jarvis! Go!</b>",
	})

	return
}

func (h *Handler) Shutdown() {
	h.svc.Shutdown()
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
	switch {
	case msg.RawMsg == "Jarvis":
		h.svc.Jarvis(msg)
	case msg.RawMsg == ".help":
		h.svc.Jhelp(msg)
	case msg.RawMsg == ".api":
		h.svc.Api(msg)
	case msg.RawMsg == ".Jeminder":
		h.svc.Jeminder(msg)
	case regBlacklist.MatchString(msg.RawMsg):
		h.svc.Blacklist(msg)
	case regPicRSA.MatchString(msg.RawMsg):
		h.svc.PicRsa(msg)
	case regPic.MatchString(msg.RawMsg):
		h.svc.Pic(msg)
	case regTimetable.MatchString(msg.RawMsg):
		h.svc.Timetable(msg)
	case regWeather.MatchString(msg.RawMsg):
		h.svc.Weather(msg)
	case regSuggest.MatchString(msg.RawMsg):
		h.svc.Suggest(msg)
	case regLog.MatchString(msg.RawMsg):
		h.svc.Jlog(msg)
	case regClog.MatchString(msg.RawMsg):
		h.svc.Clog(msg)
	case regDlog.MatchString(msg.RawMsg):
		h.svc.Dlog(msg)
	case regHocation.MatchString(msg.RawMsg):
		h.svc.Hocation(msg)
	case regUsd.MatchString(msg.RawMsg):
		h.svc.Usd(msg)
	}

	return
}
