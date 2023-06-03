package service

import (
	"errors"
	"strings"

	"github.com/abdtyx/JarvisGo/message"
)

type Message struct {
	MsgType string `json:"message_type"`
	UserID  int64  `json:"user_id"`
	GroupID int64  `json:"group_id"`
	RawMsg  string `json:"raw_message"`
}

func Jarvis(msg Message) error {
	if msg.MsgType == "private" {
		return message.PrivateMsg(msg.UserID, "I'm here, sir.")
	} else if msg.MsgType == "group" {
		return message.GroupMsg(msg.UserID, msg.GroupID, "?")
	} else {
		return errors.New(strings.ToTitle("Jarvis: Invalid message type"))
	}
}
