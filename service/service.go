package service

import (
	"errors"
	"fmt"
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

func (svc *Service) SendAndLogMsg(msg Message, privateResp, groupResp, signature string) {
	var err error
	isGroup := true
	signature += ": "

	// Send message
	switch msg.MsgType {
	case "private":
		err = message.PrivateMsg(msg.UserID, privateResp)
		isGroup = false
	case "group":
		if svc.cfg.EnableGroup {
			err = message.GroupMsg(msg.UserID, msg.GroupID, groupResp)
		} else {
			err = errors.New(strings.ToTitle("Service: Group message not enabled"))
		}
	default:
		err = errors.New(strings.ToTitle("Service: Invalid message type: " + msg.MsgType))
	}

	// Log message
	if err != nil {
		svc.Log.Println(signature, err)
	} else {
		if isGroup {
			svc.Log.Println(signature, fmt.Sprintf("Group(%v) message sent to %v: ", msg.GroupID, msg.UserID), groupResp)
		} else {
			svc.Log.Println(signature, fmt.Sprintf("Private message sent to %v: ", msg.UserID), privateResp)
		}
	}

}

func (svc *Service) Jarvis(msg Message) {
	privateResp := "I'm here, sir."
	groupResp := "?"
	sig := "Jarvis"

	svc.SendAndLogMsg(msg, privateResp, groupResp, sig)
}
