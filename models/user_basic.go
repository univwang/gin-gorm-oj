package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(100);" json:"identity"`
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password string `gorm:"column:password;type:varchar(100);" json:"password"`
	Phone    string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Mail     string `gorm:"column:mail;type:varchar(100)" json:"mail"`
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
