package test

import (
	"gin_oj/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestCreateTable(t *testing.T) {
	dsn := "root:123456@(39.101.1.158:3003)/gorm_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(models.TestCase{})
	db.AutoMigrate(models.CategoryBasic{})
	db.AutoMigrate(models.ProblemCategory{})
	db.AutoMigrate(models.ProblemBasic{})
	db.AutoMigrate(models.UserBasic{})
	db.AutoMigrate(models.SubmitBasic{})
}
