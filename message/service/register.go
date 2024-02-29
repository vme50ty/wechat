package service

import (
	"fmt"
	"genchat/message/modules"
	"genchat/message/sql"
	"genchat/message/utils"
	"math/rand"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// Register
// @Tags 首页
// @Success 200 {string} json{"code,"Register"}
// @Router /user/Login [post]

func Register(c *gin.Context) {
	log := logrus.WithFields(logrus.Fields{
		"func": "Register",
	})
	log.Debugf("Received a request to Registration.")
	// 获取解析后的表单数据
	var user modules.Userbasic
	if err := c.BindJSON(&user); err != nil {
		log.Errorf("Failed to parse request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Debugf("%v", user)
	salt := fmt.Sprintf("%06d", rand.Int31())
	if user.Name == "" || user.Password == "" || user.Phone == "" || user.Email == "" {
		c.JSON(401, gin.H{"error": "请填入用户名，密码和手机号和邮箱"})
		return
	}
	_, err := govalidator.ValidateStruct(user)
	user.Password = utils.MakePassword(user.Password, salt)
	user.Salt = salt
	if err != nil {
		c.JSON(402, gin.H{"message": fmt.Sprintf("请输入有效的手机号或邮箱: %v", err)})
		return
	}
	if err := sql.InsertSql(&user, log); err != nil {
		whicherror(err, c)
		return
	} else {
		c.JSON(200, gin.H{"message": "Registration successful"})
	}
}

// 用于找出注册中哪一步出了问题，例如用户名重复，密码重复，将错误返回给用户
func whicherror(err error, c *gin.Context) {
	gormErr, _ := err.(*mysql.MySQLError)
	if gormErr.Number == 1062 {
		if strings.Contains(gormErr.Message, "user_basic.phone") {
			c.JSON(404, gin.H{"error": "Duplicate phone number, please try again"})
			return
		} else if strings.Contains(gormErr.Message, "user_basic.name") {
			c.JSON(403, gin.H{"error": "Duplicate username, please try again"})
			return
		} else if strings.Contains(gormErr.Message, "user_basic.email") {
			c.JSON(404, gin.H{"error": "Duplicate email, please try again"})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
}
