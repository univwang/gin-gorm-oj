package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(100);" json:"identity"`
	Name      string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password  string `gorm:"column:password;type:varchar(100);" json:"password"`
	Phone     string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Mail      string `gorm:"column:mail;type:varchar(100)" json:"mail"`
	PassNum   int64  `gorm:"column:pass_num;type:int(11)" json:"pass_num"`
	SubmitNum int64  `gorm:"column:submit_num;type:int(11)" json:"submit_num"`
	IsAdmin   int    `gorm:"column:is_admin;type:tinyint(1)" json:"is_admin"` //是否是管理员1是0否
	Photo     string `gorm:"column:photo;type:varchar(200)" json:"photo"`
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
