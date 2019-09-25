package model

import (
	"github.com/jinzhu/gorm"
)

// Label 标签
type Label struct {
	Model
	Type string `json:"type"`
	Name string `json:"name"`
}

// UserLabel 用户标签,如部门、考核组、职级
type UserLabel struct {
	Model
	UserID   int    `json:"userID"`
	UserName string `json:"userName"`
	LabelID  int    `json:"labelID"`
}

// Save 保存
func (l *Label) Save() error {
	return db.Create(l).Error
}

// Update 更新
func (l *Label) Update() error {
	return db.Model(&Label{}).Updates(l).Error
}

// DelLabelByID 根据ID删除标签
func DelLabelByID(ID int) error {
	return db.Where("ID=?", ID).Delete(&Label{}).Error
}

// FindLabelByType 根据label类型查询
func FindLabelByType(types string) ([]*Label, error) {
	var labels []*Label
	err := db.Where("type=?", types).Find(&labels).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return labels, nil
}

// --------------- UserLabel ------------------

// Save 保存
func (u *UserLabel) Save() error {
	return db.Create(u).Error
}

// Update 更新
func (u *UserLabel) Update() error {
	return db.Model(&UserLabel{}).Updates(u).Error
}

// DelUserLabelByID 根据ID删除用户标签
func DelUserLabelByID(ID int) error {
	return db.Where("ID=?", ID).Delete(&UserLabel{}).Error
}

// FindUserLabelsByUserID 查询用户所有标签
func FindUserLabelsByUserID(userID int) ([]*UserLabel, error) {
	var datas []*UserLabel
	err := db.Where("userID=?", userID).Find(&datas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return datas, nil
}
