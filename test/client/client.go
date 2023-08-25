package client

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type User struct {
	UserId  int64
	GroupId int64
}

type PrivateMsgEvent struct {
	Time     int64  `json:"time"`
	SelfId   int64  `json:"self_id"`
	PostType string `json:"post_type"`
	RawMsg   string `json:"raw_message"`
	MsgType  string `json:"message_type"`
	UserID   int64  `json:"user_id"`
}

type GroupMsgEvent struct {
	Time     int64  `json:"time"`
	SelfId   int64  `json:"self_id"`
	PostType string `json:"post_type"`
	RawMsg   string `json:"raw_message"`
	MsgType  string `json:"message_type"`
	UserID   int64  `json:"user_id"`
	GroupID  int64  `json:"group_id"`
}

func (user *User) PrivateMsg(msg string) {
	_ = PrivateMsgEvent{
		Time:     time.Now().Unix(),
		SelfId:   123456,
		PostType: "message",
		RawMsg:   msg,
		MsgType:  "private",
		UserID:   user.UserId,
	}
	http.PostForm("http://127.0.0.1:5701/", url.Values{
		"time":         {strconv.FormatInt(time.Now().Unix(), 10)},
		"self_id":      {"123456"},
		"post_type":    {"message"},
		"raw_message":  {msg},
		"message_type": {"private"},
		"user_id":      {strconv.FormatInt(user.UserId, 10)},
	})
}

func (user *User) GroupMsg(msg string) {
	_ = GroupMsgEvent{
		Time:     time.Now().Unix(),
		SelfId:   123456,
		PostType: "message",
		RawMsg:   msg,
		MsgType:  "group",
		UserID:   user.UserId,
		GroupID:  user.GroupId,
	}
	http.PostForm("http://127.0.0.1:5701/", url.Values{
		"time":         {strconv.FormatInt(time.Now().Unix(), 10)},
		"self_id":      {"123456"},
		"post_type":    {"message"},
		"raw_message":  {msg},
		"message_type": {"group"},
		"user_id":      {strconv.FormatInt(user.UserId, 10)},
		"group_id":     {strconv.FormatInt(user.GroupId, 10)},
	})
}

func (user *User) Heartbeat() {
	http.PostForm("http://127.0.0.1:5701/", url.Values{
		"time":            {strconv.FormatInt(time.Now().Unix(), 10)},
		"self_id":         {"123456"},
		"post_type":       {"meta_event"},
		"meta_event_type": {"heartbeat"},
	})
}
