package message

import (
	"fmt"
	"net/http"
)

func PrivateMsg(userID int64, message string) error {
	http.Get(fmt.Sprintf("http://127.0.0.1:5700/send_private_msg?user_id=%v&message=%v&auto_escape=false", userID, message))
	return nil
}

func GroupMsg(userID, groupID int64, message string) error {
	http.Get(fmt.Sprintf("http://127.0.0.1:5700/send_group_msg?group_id=%v&message=[CQ:at,qq=%v]%%0A%v&auto_escape=false", groupID, userID, message))
	return nil
}
