package main

import (
	"genchat/message/router"
	"genchat/message/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	log := logrus.WithFields(logrus.Fields{
		"func": "oppenhttp",
	})
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	var err error
	go func() {
		err = r.Run(":8081")
	}()
	if err != nil {
		log.Error("启动页面失败")
	} else {
		log.Debugf("启动页面成功")
	}
	select {}
}
