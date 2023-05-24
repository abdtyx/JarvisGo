package msg

import (
	"github.com/gin-gonic/gin"
)

type Msg struct {
	c *gin.Context
}

func (m *Msg) PrivateMsg(userID string, message string) error {
	m.c.JSON(200, gin.H{
		"message": "helloworld",
	})
	return nil
}

func (m *Msg) GroupMsg(groupID string, message string) error {
	return nil
}
