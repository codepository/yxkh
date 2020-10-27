package model

import (
	"errors"
	"fmt"
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
	UserID      int       `gorm:"column:userId;size:32" json:"userId,omitempty"`
	Username    string    `gorm:"column:username;size:32" json:"username,omitempty"`
	Checked     string    `gorm:"size:8;column:checked" json:"checked,omitempty"`
}

// FromJSON FromJSON
func (m *ResMark) FromJSON(json map[string]interface{}) error {

	if json["projectId"] == nil {
		return fmt.Errorf("projectId 不能为空")
	}
	if json["markNumber"] == nil {
		return errors.New("markNumber 分数不能为空")
	}
	if json["accordingly"] == nil {
		return errors.New("accordingly 加分依据不能为空")
	}
	if json["markReason"] == nil {
		return errors.New("markReason 加分原因不能为空")
	}
	if json["startDate"] == nil || json["endDate"] == nil {
		return errors.New("startDate 和 endDate 加分时间段不能为空")
	}
	projectID, err := util.Interface2Int(json["projectId"])
	if err != nil {
		return err
	}
	m.ProjectID = projectID
	m.MarkNumber = json["markNumber"].(string)
	m.Accordingly = json["accordingly"].(string)
	m.MarkReason = json["markReason"].(string)
	start, err := util.ParseDate3(json["startDate"].(string))
	if err != nil {
		return err
	}
	m.StartDate = start
	end, err := util.ParseDate3(json["endDate"].(string))
	if err != nil {
		return err
	}
	m.EndDate = end
	if json["userId"] != nil {
		uid, err := util.Interface2Int(json["userId"])
		if err != nil {
			return err
		}
		m.UserID = uid
	}
	if json["username"] != nil {
		m.Username = json["username"].(string)
	}
	if json["checked"] != nil {
		check, yes := json["checked"].(string)
		if !yes {
			return errors.New("checked 应为string")
		}
		m.Checked = check
	}
	return nil
}

// FirstOrCreate 不存在就创建
func (m *ResMark) FirstOrCreate() error {
	m.CreateTime = time.Now()
	return db.Where(ResMark{ProjectID: m.ProjectID, MarkReason: m.MarkReason, Accordingly: m.Accordingly}).Assign(m).FirstOrCreate(m).Error
}

// DelMarkByIDs DelMarkByIDs
func DelMarkByIDs(ids []int) error {
	return db.Where(ids).Delete(&ResMark{}).Error
}

// UpdatesMark 只更新更改的字段
func UpdatesMark(id int, params interface{}) error {

	return db.Model(&ResMark{MarkID: id}).Updates(params).Error
}

// FindAllMark FindAllMark
func FindAllMark(query interface{}, values ...interface{}) ([]*ResMark, error) {

	var datas []*ResMark
	err := db.Where(query, values...).Find(&datas).Error
	// err := db.Raw(rawSQL, values...).Scan(&datas).Error
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
