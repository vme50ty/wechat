package service

import (
	"fmt"
	"genchat/message/modules"
	"genchat/message/utils"
	"math/rand"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// GetUserList
// @Tags 首页
// @Success 200 {string} json{"code,"message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := modules.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

func DeleteUser(c *gin.Context) {
	user := modules.Userbasic{}
	name := c.Query("Name")
	if name == "" {
		c.JSON(400, gin.H{"error": "用户名未提供"})
		return
	}
	user.Name = name
	modules.DeleteUser(user)
	c.JSON(200, gin.H{"message": "删除用户成功"})
}

func UpdateUser(c *gin.Context) {
	user := modules.Userbasic{}
	c.Bind(&user)
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(200, gin.H{"message": "修改参数不匹配"})
		return
	}
	if user.Password != "" {
		salt := fmt.Sprintf("%06d", rand.Int31())
		user.Password = utils.MakePassword(user.Password, salt)
		user.Salt = salt
	}
	modules.UpdateUser(user)
	c.JSON(200, gin.H{"message": "修改用户成功"})
}

func FindUserByName(c *gin.Context) {
	name := c.Query("Name")
	if name == "" {
		c.JSON(200, gin.H{
			"message": "请输入有效名字",
		})
	}
	data := modules.FindUserByName(name)
	c.JSON(200, gin.H{
		"message": data,
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendUserMsg(c *gin.Context) {
	log := logrus.WithFields(logrus.Fields{
		"func": "sendusermessage",
	})
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			log.Error(err)
		}
		log.Debug("websocket连接关闭")
	}(ws)
	MsgHandler(ws, c, log)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context, log *logrus.Entry) {
	msg, err := utils.Subscribe(c, utils.PublishKey, log)
	if err != nil {
		log.Error(err)
	}
	tm := time.Now().Format("2008-01-02 15:01:02")
	log.Debugf("[ws][%v]:%v", tm, msg)
	err = ws.WriteMessage(1, []byte(msg))
	if err != nil {
		log.Error(err)
	}
}

func SendMsg(c *gin.Context) {
	log := logrus.WithFields(logrus.Fields{
		"func": "sendmessage",
	})
	modules.Chat(c.Writer, c.Request, log)
}
