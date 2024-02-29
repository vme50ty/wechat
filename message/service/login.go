package service

import (
	"genchat/message/modules"
	"genchat/message/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Login
// @Tags 首页
// @Success 200 {string} json{"code,"login"}
// @Router /user/Login [get]
func Login(c *gin.Context) {
	log := logrus.WithFields(logrus.Fields{
		"func": "SignIn",
	})
	log.Debugf("Received a request to login.")
	var user modules.Userbasic
	if err := c.Bind(&user); err != nil {
		log.Errorf("Failed to parse request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "In request body"})
		return
	}
	if user.Name == "" && user.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or phone must be provided"})
		return
	}

	if user.Name != "" {
		if judge := sql.IfadminName(c, &user, log); judge {
			c.JSON(200, gin.H{"message": "Login successful"})
			return
		}
	}

	if user.Phone != "" {
		if judge := sql.IfadminPhone(c, &user, log); judge {
			c.JSON(200, gin.H{"message": "Login successful"})
			return
		}
	}

}
