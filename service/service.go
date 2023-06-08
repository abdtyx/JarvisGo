package service

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/abdtyx/JarvisGo/config"
	"github.com/abdtyx/JarvisGo/message"
)

type Service struct {
	cfg *config.Config
	Log *log.Logger
}

type Message struct {
	MsgType string `json:"message_type"`
	UserID  int64  `json:"user_id"`
	GroupID int64  `json:"group_id"`
	RawMsg  string `json:"raw_message"`
}

func InitService() (*Service, error) {
	var svc Service
	var err error

	svc.cfg, err = config.LoadConfig()
	if err != nil {
		return nil, err
	}

	svc.Log = log.Default()
	f, err := os.OpenFile("./jlogs/"+time.Now().String()+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	svc.Log.SetOutput(f)

	return &svc, nil
}

func (svc *Service) Jarvis(msg Message) {
	var err error
	if msg.MsgType == "private" {
		err = message.PrivateMsg(msg.UserID, "I'm here, sir.")
	} else if svc.cfg.EnableGroup && msg.MsgType == "group" {
		err = message.GroupMsg(msg.UserID, msg.GroupID, "?")
	} else {
		err = errors.New(strings.ToTitle("Jarvis: Invalid message type"))
	}

	if err != nil {
		svc.Log.Println(err)
	} else {
		svc.Log.Println(err.Error())
	}
}
