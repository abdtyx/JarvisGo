package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/abdtyx/JarvisGo/config"
	"github.com/abdtyx/JarvisGo/errdefs"
	"github.com/abdtyx/JarvisGo/message"
	"github.com/abdtyx/JarvisGo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Service struct {
	Cfg *config.Config
	Log *log.Logger
	db  *gorm.DB
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

	// Print config for confirmation
	svc.Cfg.PrintConfig()

	// Initialize logger
	svc.Log = log.Default()
	f, err := os.OpenFile(svc.Cfg.WorkingDirectory+"jlogs/"+time.Now().String()+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	svc.Log.SetOutput(f)

	// Initialize database
	db, err := gorm.Open(mysql.Open(svc.Cfg.DSN), &gorm.Config{})
	if err != nil {
		panic("Service: Failed to open database connection, error: " + err.Error())
	}
	db = db.Debug()

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
	timetable := "%0A%0AGet your timetable(only for xjtu): %0AFormat is as follows: %0A.timetable%0AYour username%0AYour password%0AExample: %0A.require table%0A123456%0A123456"
	weather := "%0A%0AGet region weather: %0AFormat is as follows: %0A.weather region (superior administrative division of the region)%0AExamples: %0A.weather 西安%0A.weather 海淀 北京"
	endl := "%0A"
	_log := endl + endl + "Make a log: " + endl + "Format is as follows: " + endl + ".log Your words here" + endl + "Example: " + endl + ".log I have a record to record."
	clog := endl + endl + "Check your log: " + endl + "Format is as follows: " + endl + ".clog (year-month-day)" + endl + ".clog [six characters in front of a specific log]" + endl + "Examples: " + endl + ".clog" + endl + ".clog 2022-1-5" + endl + ".clog 99a5b3"
	dlog := endl + endl + "Delete a log: " + endl + "Format is as follows: " + endl + ".dlog [six characters in front of a specific log]" + endl + "Example: " + endl + ".dlog 99a5b3"
	resp = resp + timetable + weather + _log + clog + dlog
	sig := "Api"

	svc.SendAndLogMsg(msg, resp, resp, sig)
}

func (svc *Service) Jeminder(msg Message) {
	sig := "Jeminder"
	// Group msg only
	if msg.MsgType != "group" {
		svc.SendAndLogMsg(msg, "Sir, please use this command in a group.", "", sig)
		return
	}
	if !svc.checkMaster(msg) {
		resp := errdefs.ErrPermissionDenied{
			BaseErr: errdefs.BaseErr{
				Sig: sig, User: msg.UserID, Group: msg.GroupID,
			},
		}
		svc.SendAndLogMsg(msg, "", resp.String(), sig)
		return
	}

	var data model.Jeminder
	if err := svc.db.Where("id = ?", msg.GroupID).First(&data).Error; err != nil {
		// not found
		data.Id = uint(msg.GroupID)
		if err = svc.db.Create(&data).Error; err != nil {
			// Found error when create
			// modify MsgType to send errors to master only, preventing secrets from being revealed
			msg.MsgType = "private"
			svc.SendAndLogMsg(msg, "", err.Error(), sig)
			return
		}
		svc.SendAndLogMsg(msg, "", "Jeminder on.", sig)
		return
	}

	if err := svc.db.Delete(&data).Error; err != nil {
		// Found error when delete
		// modify MsgType to send errors to master only, preventing secrets from being revealed
		msg.MsgType = "private"
		svc.SendAndLogMsg(msg, "", err.Error(), sig)
		return
	}
	svc.SendAndLogMsg(msg, "", "Jeminder off.", sig)
	return
}

func (svc *Service) Blacklist(msg Message) {
	// TODO:
	// sig := "blacklist"

	return
}

func (svc *Service) Pic(msg Message) {
	// TODO:
	return
}

func (svc *Service) PicRsa(msg Message) {
	// TODO:
	return
}

func (svc *Service) Timetable(msg Message) {
	// TODO:
	return
}

func (svc *Service) Weather(msg Message) {
	// TODO:
	return
}

func (svc *Service) Suggest(msg Message) {
	// TODO:
	return
}

func (svc *Service) Jlog(msg Message) {
	// TODO:
	return
}

func (svc *Service) Clog(msg Message) {
	// TODO:
	return
}

func (svc *Service) Dlog(msg Message) {
	// TODO:
	return
}

func (svc *Service) Hocation(msg Message) {
	// TODO:
	return
}

func (svc *Service) Usd(msg Message) {
	// TODO:
	return
}

func (svc *Service) CheckBlacklist(msg Message) (userFlag, groupFlag bool) {
	userFlag = false
	groupFlag = false
	var data model.Blacklist
	if err := svc.db.Where("id = ? AND type = ?", msg.UserID, msg.MsgType).First(&data).Error; err == nil {
		userFlag = true
	}

	if svc.Cfg.EnableGroup {
		if err := svc.db.Where("id = ? AND type = ?", msg.GroupID, msg.MsgType).First(&data).Error; err == nil {
			groupFlag = true
		}
	}

	return userFlag, groupFlag
}

/**
 * TODO:
 * prior msg not sent to group
 * i.e., permission handler
 */
func (svc *Service) checkMaster(msg Message) bool {
	for _, v := range svc.Cfg.Masters {
		if v == msg.UserID {
			return true
		}
	}
	return false
}

/*
 * DONE:
 * Blacklist write back on terminating with signal handler
 */
func (svc *Service) Shutdown() error {
	svc.Log.Println("Received interrupt or terminate signal.")
	// Close database
	if svc.db != nil {
		db, err := svc.db.DB()
		if err != nil {
			return err
		}
		return db.Close()
	}
	return nil
}
