package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// ResProject 项目
type ResProject struct {
	ProjectID int `gorm:"primary_key;column:projectId" json:"projectId"`
	// 一个中文长度为：3，英文为： 1
	ProjectContent string    `gorm:"size:1000;column:projectContent" json:"projectContent"`
	UserID         int       `gorm:"column:userId" json:"userId"`
	StartDate      time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate        time.Time `gorm:"column:endDate" json:"endDate"`
	Progress       string    `gorm:"size:1000" json:"progress"`
	Createtime     time.Time `gorm:"column:createTime" json:"createTime"`
}

// FindAllProjectPaged 分页查询
func FindAllProjectPaged(pageIndex, pageSize int, sql string) ([]*ResProject, int, error) {
	var datas []*ResProject
	var count int
	if pageIndex == 0 {
		pageIndex = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	err := db.Where(sql).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&datas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	err = db.Model(&ResProject{}).Where(sql).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return datas, count, nil
}

// Save Save
func (p *ResProject) Save() error {
	if len(p.ProjectContent) > 1000 {
		return errors.New("projectcontent 长度不能超过1000")
	}
	if len(p.Progress) > 1000 {
		return errors.New("progress 长度不能超过1000")
	}
	p.Createtime = time.Now()
	return db.Create(p).Error
}

// DelProjectByID DelProjectByID
func DelProjectByID(id int) error {
	return db.Where("projectId=?", id).Delete(&ResProject{}).Error
}

// Update Update
func (p *ResProject) Update() error {
	return db.Model(&ResProject{}).Updates(p).Error
}
