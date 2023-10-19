package service

import (
	"genchat/message/modules"

	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tags 首页
// @Success 200 {string} json{"code,"message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*modules.Userbasic, 10)
	data = modules.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}
