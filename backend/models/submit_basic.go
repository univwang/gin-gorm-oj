package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model
	Identity        string        `gorm:"column:identity;type:varchar(100);" json:"identity"`
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(100);" json:"problem_identity"`
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(100);" json:"user_identity"`
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"`
	Path            string        `gorm:"column:path;type:varchar(100);" json:"path"`
	Status          int           `gorm:"column:status;type:tinyint(1);" json:"status"` //-1待判断1正确2错误3运行超时4内存溢出
}

func (s *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity string, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content")
	}).Preload("UserBasic")
	if problemIdentity != "" {
		tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx.Where("status = ?", status)
	}
	return tx
}
