package model

import (
	"github.com/jinzhu/gorm"
)

// UserMark 用户和分数
type UserMark struct {
	UserID   int    `json:"userID"`
	UserName string `json:"userName"`
	Mark     int    `json:"mark"`
}

// FindUserMarkPagedWithSQL 根据SQL语句查询用户总分
func FindUserMarkPagedWithSQL(sql string, paged bool) ([]*UserMark, int, error) {
	var datas []*UserMark
	var count int
	err := db.Raw(sql).Scan(&datas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	if paged {
		err = db.Raw(sql).Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return datas, count, nil
}
