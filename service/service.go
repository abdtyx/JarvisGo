package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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
	svc.db, err = gorm.Open(mysql.Open(svc.Cfg.DSN), &gorm.Config{})
	if err != nil {
		panic("Service: Failed to open database connection, error: " + err.Error())
	}
	svc.db = svc.db.Debug()

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
	resp := "I'm Jarvis, Mr.A1pha's personal assistant. What can I do for you?\n---menu---\n.api --See functions available now\n.pic --setu ps: a little delay and may fail to transmit\n.suggest --Make your suggestions, including but not limited to the functions you want Jarvis to implement.\nExample: .suggest Jarvis, I want you to ..."
	sig := "Jhelp"

	svc.SendAndLogMsg(msg, resp, resp, sig)
}

func (svc *Service) Api(msg Message) {
	resp := "Sir, you have chosen to see functions available now. I list them out for you."
	timetable := "\n\nGet your timetable(only for xjtu): \nFormat is as follows: \n.timetable\nYour username\nYour password\nExample: \n.require table\n123456\n123456"
	weather := "\n\nGet region weather: \nFormat is as follows: \n.weather region (superior administrative division of the region)\nExamples: \n.weather 西安\n.weather 海淀 北京"
	endl := "\n"
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
			// modify MsgType to send errors to masters only, preventing secrets from being revealed
			msg.MsgType = "private"
			svc.SendAndLogMsg(msg, err.Error(), "", sig)
			return
		}
		svc.SendAndLogMsg(msg, "", "Jeminder on.", sig)
		return
	}

	if err := svc.db.Delete(&data).Error; err != nil {
		// Found error when delete
		// modify MsgType to send errors to masters only, preventing secrets from being revealed
		msg.MsgType = "private"
		svc.SendAndLogMsg(msg, err.Error(), "", sig)
		return
	}
	svc.SendAndLogMsg(msg, "", "Jeminder off.", sig)
	return
}

func (svc *Service) Blacklist(msg Message) {
	// TODO:
	sig := "blacklist"
	if !svc.checkMaster(msg) {
		resp := errdefs.ErrPermissionDenied{
			BaseErr: errdefs.BaseErr{
				Sig: sig, User: msg.UserID, Group: msg.GroupID,
			},
		}
		svc.SendAndLogMsg(msg, resp.String(), resp.String(), sig)
		return
	}

	//			  0				1			   2								 3	   4
	// example: [".blacklist", "user, group", "add, remove, inspect, list(ls)", "id", "comment"]
	args := parseArgs(msg.RawMsg)
	if args[0] != ".blacklist" {
		// Recognized command isn't ".blacklist"
		return
	}

	if !((len(args) == 5 && args[2] == "add") || (len(args) == 4 && (args[2] == "remove" || args[2] == "inspect")) || (len(args) == 3 && (args[2] == "list" || args[2] == "ls"))) {
		// Must have at least 4 args, but no more than 5 args
		goto wrongArgs
	}

	if args[1] != "user" && args[1] != "group" {
		goto wrongArgs
	}

	// The most special edge cases
	if args[1] == "user" && args[2] == "add" {
		// A master should not be put onto the blacklist
		userid, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			goto wrongArgs
		}
		if svc.checkMaster(Message{UserID: userid}) {
			resp := errdefs.ErrPermissionDenied{
				BaseErr: errdefs.BaseErr{
					Sig: sig, User: msg.UserID, Group: msg.GroupID,
				},
			}
			svc.SendAndLogMsg(msg, resp.String(), resp.String(), sig)
			return
		}
	}

	switch args[2] {
	case "add":
		id, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			goto wrongArgs
		}
		var blacklistEntry model.Blacklist
		err = svc.db.Where("id = ? AND type = ?", id, args[1]).First(&blacklistEntry).Error

		if err == gorm.ErrRecordNotFound {
			// Record not found
			// Add a new record
			blacklistEntry.Id = uint(id)
			blacklistEntry.Type = args[1]
			blacklistEntry.Comment = args[4]
			err = svc.db.Create(&blacklistEntry).Error
			if err != nil {
				// modify MsgType to send errors to masters only, preventing secrets from being revealed
				msg.MsgType = "private"
				svc.SendAndLogMsg(msg, err.Error(), "", sig)
				return
			}
			resp := fmt.Sprintf("Sir, the %v (%v) is added to the blacklist.", args[1], args[3])
			svc.SendAndLogMsg(msg, resp, resp, sig)
			return
		} else if err != nil {
			// modify MsgType to send errors to masters only, preventing secrets from being revealed
			msg.MsgType = "private"
			svc.SendAndLogMsg(msg, err.Error(), "", sig)
			return
		}
		// Record already exists
		resp := fmt.Sprintf("Sir, the %v (%v) is already in the blacklist.", args[1], args[3])
		svc.SendAndLogMsg(msg, resp, resp, sig)
		return

	case "remove":
		id, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			goto wrongArgs
		}
		err = svc.db.Where("id = ? AND type = ?", id, args[1]).Delete(&model.Blacklist{}).Error
		if err == gorm.ErrRecordNotFound {
			// Record not found
			resp := fmt.Sprintf("Sir, the %v (%v) is not in the blacklist.", args[1], args[3])
			svc.SendAndLogMsg(msg, resp, resp, sig)
			return
		} else if err != nil {
			// modify MsgType to send errors to masters only, preventing secrets from being revealed
			msg.MsgType = "private"
			svc.SendAndLogMsg(msg, err.Error(), "", sig)
			return
		}

		// Successfully deleted
		resp := fmt.Sprintf("Sir, the %v (%v) is removed from blacklist.", args[1], args[3])
		svc.SendAndLogMsg(msg, resp, resp, sig)
		return

	case "inspect":
		// the result of inspect should only be sent through private messages
		msg.MsgType = "private"
		id, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			goto wrongArgs
		}
		var blacklistEntry model.Blacklist
		err = svc.db.Where("id = ? AND type = ?", id, args[1]).First(&blacklistEntry).Error
		if err == gorm.ErrRecordNotFound {
			// Record not found
			resp := fmt.Sprintf("Sir, the %v (%v) is not in the blacklist.", args[1], args[3])
			svc.SendAndLogMsg(msg, resp, resp, sig)
			return
		} else if err != nil {
			svc.SendAndLogMsg(msg, err.Error(), "", sig)
			return
		}

		// Successfully found
		resp := "Sir, I found this record for you: \n"
		resp += blacklistEntry.String()
		svc.SendAndLogMsg(msg, resp, resp, sig)
		return

	case "list", "ls":
		// the result of inspect should only be sent through private messages
		msg.MsgType = "private"
		var blacklistEntries []model.Blacklist
		err := svc.db.Where("type = ?", args[1]).Find(&blacklistEntries).Error
		if err == gorm.ErrRecordNotFound {
			// Record not found
			resp := fmt.Sprintf("I'm sorry sir, I found no record under category (%v).", args[1])
			svc.SendAndLogMsg(msg, resp, resp, sig)
			return
		} else if err != nil {
			svc.SendAndLogMsg(msg, err.Error(), "", sig)
			return
		}

		// Successfully found
		resp := "Sir, I found these records for you: \n"
		for _, v := range blacklistEntries {
			resp += v.String()
		}
		svc.SendAndLogMsg(msg, resp, resp, sig)

	default:
		// modify MsgType to send errors to masters only, preventing secrets from being revealed
		msg.MsgType = "private"
		svc.SendAndLogMsg(msg, "**WARNING**: args[2] is tampered.", "", sig)
		return
	}

wrongArgs:
	resp := errdefs.ErrWrongArgs{
		BaseErr: errdefs.BaseErr{
			Sig: sig, User: msg.UserID, Group: msg.GroupID,
		},
	}
	svc.SendAndLogMsg(msg, resp.String(), resp.String(), sig)
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
	if err := svc.db.Where("id = ? AND type = ?", msg.UserID, "user").First(&data).Error; err == nil {
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
 * DONE:
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

/**
 * args parser
 */
func parseArgs(cmd string) []string {
	var args []string
	cmds := strings.Split(cmd, " ")
	for _, v := range cmds {
		if len(v) > 0 {
			args = append(args, v)
		}
	}
	return args
}
