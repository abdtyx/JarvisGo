package message

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func Msg(userID int64, groupID int64, message string, enableGroup bool) error {
	if enableGroup {
		return GroupMsg(userID, groupID, message)
	}

	if userID != 0 {
		return PrivateMsg(userID, message)
	}

	return errors.New(strings.ToTitle("Msg: Invalid message type"))
}

func PrivateMsg(userID int64, message string) error {
	message = url.QueryEscape(message)
	_, err := http.Get(fmt.Sprintf("http://127.0.0.1:5700/send_private_msg?user_id=%v&message=%v&auto_escape=false", userID, message))
	return err
}

func GroupMsg(userID, groupID int64, message string) error {
	message = url.QueryEscape(message)
	_, err := http.Get(fmt.Sprintf("http://127.0.0.1:5700/send_group_msg?group_id=%v&message=[CQ:at,qq=%v]%%0A%v&auto_escape=false", groupID, userID, message))
	return err
}
