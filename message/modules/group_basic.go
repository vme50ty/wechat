package modules

import "gorm.io/gorm"

// 群信息
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint64
	Type    int
	Icon    string //群头像
	Desc    string //描述
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
