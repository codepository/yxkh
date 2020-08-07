package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/mumushuiding/util"
)

// ResMark ResMark
type ResMark struct {
	MarkID      int       `gorm:"primary_key;column:markId" json:"markId,omitempty"`
	ProjectID   int       `gorm:"column:projectId" json:"projectId,omitempty"`
	MarkNumber  string    `gorm:"size:32;column:markNumber" json:"markNumber,omitempty"`
	MarkReason  string    `gorm:"size:1000;column:markReason" json:"markReason,omitempty"`
	Accordingly string    `gorm:"size:2000" json:"accordingly,omitempty"`
	StartDate   time.Time `gorm:"column:startDate" json:"startDate,omitempty"`
	EndDate     time.Time `gorm:"column:endDate" json:"endDate,omitempty"`
	CreateTime  time.Time `gorm:"column:createTime" json:"createTime,omitempty"`
	UserID      string    `gorm:"column:userId;size:32" json:"userId,omitempty"`
	Username    string    `gorm:"column:username;size:32" json:"username,omitempty"`
	Checked     string    `gorm:"size:8;column:checked" json:"checked,omitempty"`
}

// FindAllMark FindAllMark
func FindAllMark(rawSQL string, values ...interface{}) ([]*ResMark, error) {

	var datas []*ResMark
	err := db.Raw(rawSQL, values...).Scan(&datas).Error
	return datas, err
}

// FindAllMarkPaged FindAllMarkPaged
func FindAllMarkPaged(pageIndex, pageSize int, sql string) ([]*ResMark, int, error) {
	var datas []*ResMark
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
	err = db.Model(&ResMark{}).Where(sql).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return datas, count, nil
}

// Save Save
func (m *ResMark) Save() error {
	if err := m.prepareData(); err != nil {
		return err
	}
	m.CreateTime = time.Now()
	return db.Create(m).Error
}

// SaveTx 事务
func (m *ResMark) SaveTx(tx *gorm.DB) error {
	if err := m.prepareData(); err != nil {
		return err
	}
	m.CreateTime = time.Now()
	return tx.Create(m).Error
}

// Update Update
func (m *ResMark) Update() error {
	if err := m.prepareData(); err != nil {
		return err
	}
	return db.Model(&ResMark{}).Updates(m).Error
}

// DelMarkByID DelMarkByID
func DelMarkByID(id int) error {
	return db.Where("markId=?", id).Delete(&ResMark{}).Error
}

func (m *ResMark) prepareData() error {
	if len(m.MarkNumber) == 0 || m.MarkNumber == "0" { // 评分不能为空或者0
		return errors.New("字段【markNumber】不能为0或者空")
	}
	yes, err := util.IsDoubleStr(m.MarkNumber)
	if err != nil {
		return err
	}
	if !yes {
		return errors.New("字段【markNumber】只能为浮点数")
	}
	if len(m.MarkReason) > 1000 {
		return errors.New("字段【markReason】长度不能超过1000")
	}
	if len(m.Accordingly) > 2000 {
		return errors.New("字段【accordingly】长度不能超过2000")
	}
	if len(m.Checked) == 0 {
		m.Checked = "0"
	}
	return nil
}
