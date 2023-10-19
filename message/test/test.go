package test

import (
	"fmt"
	"genchat/message/modules"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("xuty:password@tcp(127.0.0.1:3306)/genchat?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic("can not open mysql" + err.Error())
	}
	db.AutoMigrate((&modules.Userbasic{}))
	user := &modules.Userbasic{}
	user.Name = "nihao"
	db.Create(user)
	fmt.Println(db.First(user, 1))

}
