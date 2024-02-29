package sql

import (
	"errors"
	"genchat/message/modules"
	"genchat/message/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InsertSql(user *modules.Userbasic, log *logrus.Entry) error {
	result := utils.DB.Create(user)
	if result.Error != nil {
		log.Error("Failed to insert user into database:", result.Error)
		return result.Error
	}
	return nil
}

func IfadminName(c *gin.Context, user *modules.Userbasic, log *logrus.Entry) bool {
	var findusr modules.Userbasic
	result := utils.DB.Where("name = ?", user.Name).First(&findusr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败,不存在该用户"})
		} else {
			log.Errorf("查询数据库失败: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		}
		return false
	}
	flag := utils.ValidPassword(user.Password, findusr.Salt, findusr.Password)
	if !flag {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return false
	}
	return true
}

func IfadminPhone(c *gin.Context, user *modules.Userbasic, log *logrus.Entry) bool {
	var findusr modules.Userbasic
	result := utils.DB.Where("phone = ?", user.Phone).First(&findusr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败,不存在该用户"})
		} else {
			log.Errorf("查询数据库失败: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务器错误"})
		}
		return false
	}
	flag := utils.ValidPassword(user.Password, findusr.Salt, findusr.Password)
	if !flag {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return false
	}
	return true
}
