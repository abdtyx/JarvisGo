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

func (svc *Service) Jhelp(msg Message) {
	resp := "I'm Jarvis, Mr.A1pha's personal assistant. What can I do for you?%0A---menu---%0A.api --See functions available now%0A.pic --setu ps: a little delay and may fail to transmit%0A.suggest --Make your suggestions, including but not limited to the functions you want Jarvis to implement.%0AExample: .suggest Jarvis, I want you to ..."
	sig := "Jhelp"

	svc.SendAndLogMsg(msg, resp, resp, sig)
}

func (svc *Service) Api(msg Message) {
	resp := "Sir, you have chosen to see functions available now. I list them out for you."
	requireTable := "%0A%0AGet your timetable(only for xjtu): %0AFormat is as follows: %0A.require table%0AYour username%0AYour password%0AExample: %0A.require table%0A123456%0A123456"
	weather := "%0A%0AGet region weather: %0AFormat is as follows: %0A.weather region (superior administrative division of the region)%0AExamples: %0A.weather 西安%0A.weather 海淀 北京"
	endl := "%0A"
	_log := endl + endl + "Make a log: " + endl + "Format is as follows: " + endl + ".log Your words here" + endl + "Example: " + endl + ".log I have a record to record."
	clog := endl + endl + "Check your log: " + endl + "Format is as follows: " + endl + ".clog (year-month-day)" + endl + ".clog [six characters in front of a specific log]" + endl + "Examples: " + endl + ".clog" + endl + ".clog 2022-1-5" + endl + ".clog 99a5b3"
	dlog := endl + endl + "Delete a log: " + endl + "Format is as follows: " + endl + ".dlog [six characters in front of a specific log]" + endl + "Example: " + endl + ".dlog 99a5b3"
	resp = resp + requireTable + weather + _log + clog + dlog
	sig := "Api"

	svc.SendAndLogMsg(msg, resp, resp, sig)
}
