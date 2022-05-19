package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"column:name;type:varchar(64);not null;default:'';comment:名字"`
	Age   uint8  `gorm:"column:name;type:tinyint(3);not null;default:0,comment:年龄"`
	Phone string `gorm:"column:phone;type:varchar(11);not null;default:'';comment:手机号码"`
}

func (u User) TableName() string {
	return "user"
}

func (u User) Migrate(db *gorm.DB) error {
	return nil
}
