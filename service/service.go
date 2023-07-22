package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abdtyx/JarvisGo/config"
	"github.com/abdtyx/JarvisGo/message"
)

type Service struct {
	Cfg            *config.Config
	Log            *log.Logger
	userBlacklist  []string
	groupBlacklist []string
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

	// Load config
	svc.Cfg, err = config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	svc.Log = log.Default()
	f, err := os.OpenFile(svc.Cfg.WorkingDirectory+"jlogs/"+time.Now().String()+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	svc.Log.SetOutput(f)

	// Initialize blacklist, skip if error
	svc.readBlacklist()

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
		if svc.Cfg.EnableGroup {
			err = message.GroupMsg(msg.UserID, msg.GroupID, groupResp)
		} else {
			err = errors.New(strings.ToTitle("Service: Group message not enabled (should not reach here)"))
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

func (svc *Service) CheckBlacklist(msg Message) (userFlag, groupFlag bool) {
	// check member one by one
	userFlag = false
	groupFlag = false

	for _, v := range svc.userBlacklist {
		if v == strconv.FormatInt(msg.UserID, 10) {
			userFlag = true
			break
		}
	}

	if msg.MsgType == "group" && svc.Cfg.EnableGroup {
		for _, v := range svc.groupBlacklist {
			if v == strconv.FormatInt(msg.GroupID, 10) {
				groupFlag = true
				break
			}
		}
	}
	return userFlag, groupFlag
}

/**
* TODO:
* prior msg not sent to group
 */
func (svc *Service) checkMaster(msg Message) bool {
	for _, v := range svc.Cfg.Masters {
		if v == msg.UserID {
			return true
		}
	}
	return false
}

func (svc *Service) readBlacklist() {
	// open blacklist file, then read from blacklist
	userBlacklistByte, err := ioutil.ReadFile(svc.Cfg.WorkingDirectory + "jdata/UserBlacklist.txt")
	if err != nil {
		svc.userBlacklist = nil
		svc.Log.Println("CheckBlacklist: ", err)
		return
	}
	blacklist := hex.EncodeToString(userBlacklistByte)
	svc.userBlacklist = strings.Split(blacklist, "\n")

	groupBlacklistByte, err := ioutil.ReadFile(svc.Cfg.WorkingDirectory + "jdata/GroupBlacklist.txt")
	if err != nil {
		svc.groupBlacklist = nil
		svc.Log.Println("CheckBlacklist: ", err)
		return
	}
	blacklist = hex.EncodeToString(groupBlacklistByte)
	svc.groupBlacklist = strings.Split(blacklist, "\n")
}

/*
* TODO:
* Blacklist write back on terminating with signal handler
 */
