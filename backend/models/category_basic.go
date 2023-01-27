package models

import "gorm.io/gorm"

type CategoryBasic struct {
	gorm.Model
	Identity string `gorm:"column:identity;type:varchar(100);" json:"identity"` //分类标识
	Name     string `gorm:"column:name;type:varchar(100);" json:"name"`         //分类名称
	ParentId int    `gorm:"column:parent_id;type:int(11);" json:"parent_id"`    //父亲ID
}

func (c *CategoryBasic) TableName() string {
	return "category_basic"
}
