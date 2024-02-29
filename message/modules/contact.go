package modules

import "gorm.io/gorm"

// 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint64 //谁的关系
	TargetId uint64 //对应的谁
	Type     int    //对应关系类型 0，1，2
	Desc     string //描述信息
}

func (table *Contact) TableName() string {
	return "Contact"
}
