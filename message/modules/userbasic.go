package modules

import (
	"genchat/message/utils"

	"gorm.io/gorm"
)

type Userbasic struct {
	gorm.Model
	Identity   string
	Name       string `gorm:"unique"`
	Password   string
	Phone      string `gorm:"unique;size:13" valid:"matches(^1[3-9]{1}\\d)"`
	Email      string `valid:"email"`
	ClientIp   string
	ClientPort string
	Salt       string
	// LoginTime     time.Time
	// Heartbeattime time.Time
	// LoginOutTime  time.Time `gorm:"column:login_out_time"`
	IsLogout   bool
	DeviceInfo string
}

func (table *Userbasic) TableName() string {
	return "user_basic"
}
func GetUserList() []*Userbasic {
	data := make([]*Userbasic, 10)
	utils.DB.Find(&data)
	return data
}

func NewUser(userName string, password string, phone string, email string) *Userbasic {
	newUser := &Userbasic{
		Identity: "some_identity",
		Name:     userName,
		Password: password,
		Phone:    phone,
		Email:    email,
	}
	return newUser
}

func DeleteUser(user Userbasic) *gorm.DB {
	return utils.DB.Unscoped().Where("name = ?", user.Name).Delete(&user)
}

func UpdateUser(user Userbasic) *gorm.DB {
	// 假设 Userbasic 中有一个名为 Name 的字段表示唯一标识符
	return utils.DB.Model(&Userbasic{}).Where("Name = ?", user.Name).Updates(Userbasic{
		Phone:    user.Phone,
		Password: user.Password,
		Email:    user.Email,
		Salt:     user.Salt,
	})
}

func FindUserByName(name string) Userbasic {
	user := Userbasic{}
	utils.DB.Where("Name = ?", name).First(&user)
	return user
}
