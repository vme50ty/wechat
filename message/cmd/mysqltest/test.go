package main

import (
	"genchat/message/modules"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log := logrus.WithFields(logrus.Fields{
		"func": "connect DB",
	})
	db, err := gorm.Open(mysql.Open("xuty:password@tcp(127.0.0.1:3306)/genchat?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("can not open mysql" + err.Error())
	}
	// if err := db.AutoMigrate(&modules.Userbasic{}); err != nil {
	// 	log.Error("Failed to auto migrate:", err)
	// }
	if err := db.AutoMigrate(&modules.GroupBasic{}); err != nil {
		log.Error("Failed to auto migrate:", err)
	}
	if err := db.AutoMigrate(&modules.Contact{}); err != nil {
		log.Error("Failed to auto migrate:", err)
	}
	// user := modules.NewUser("nihao", "123456", "15736992173", "233323123")
	// db.Create(user)
	// fmt.Println(db.First(user, 1))

}
