package model

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrDataStoreFailed  = errors.New("数据保存失败")
	ErrDataQueryFailed  = errors.New("数据查询失败")
	ErrDataDeleteFailed = errors.New("数据删除失败")
	ErrDataNotFound     = errors.New("数据不存在")
)

type IMigrate interface {
	Migrate(db *gorm.DB) error
}

func Migrate(db *gorm.DB) error {
	models := []IMigrate{
		&User{},
	}

	for _, model := range models {
		if err := model.Migrate(db); err != nil {
			return err
		}
	}

	return nil
}
