package modules

import (
	"genchat/message/utils"
	"time"

	"gorm.io/gorm"
)

type Userbasic struct {
	gorm.Model
	Identity      string
	Name          string
	Password      string
	Phone         string
	Email         string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	Heartbeattime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *Userbasic) TableName() string {
	return "user_basic"
}
func GetUserList() []*Userbasic {
	data := make([]*Userbasic, 10)
	utils.DB.Find(&data)
	return data
}
