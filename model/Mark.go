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
	MarkID     int    `gorm:"primary_key;column:markId" json:"markId,omitempty"`
	ProjectID  int    `gorm:"column:projectId" json:"projectId,omitempty"`
	MarkNumber string `gorm:"size:32;column:markNumber" json:"markNumber,omitempty"`
	// MarkReason 加分原因
	MarkReason string `gorm:"size:1000;column:markReason" json:"markReason,omitempty"`
	// Accordingly 加分依据的规则
	Accordingly string    `gorm:"size:2000" json:"accordingly,omitempty"`
	StartDate   string    `gorm:"column:startDate" json:"startDate,omitempty"`
	EndDate     string    `gorm:"column:endDate" json:"endDate,omitempty"`
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
	_, err = util.ParseDate3(json["startDate"].(string))
	if err != nil {
		return err
	}
	m.StartDate = json["startDate"].(string)
	_, err = util.ParseDate3(json["endDate"].(string))
	if err != nil {
		return err
	}
	m.EndDate = json["endDate"].(string)
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
	} else {
		m.Checked = "0"
	}
	return nil
}

// FirstOrCreate 不存在就创建
func (m *ResMark) FirstOrCreate() error {
	m.CreateTime = time.Now()
	return db.Where(ResMark{ProjectID: m.ProjectID, MarkReason: m.MarkReason, Accordingly: m.Accordingly}).Attrs(m).FirstOrCreate(m).Error
}

// DelMarkByIDs DelMarkByIDs
func DelMarkByIDs(ids []int) error {
	return db.Where(ids).Delete(&ResMark{}).Error
}

// UpdatesMark 只更新更改的字段
func UpdatesMark(id int, params interface{}) error {

	return db.Model(&ResMark{MarkID: id}).Updates(params).Error
}

// FindSingleMark 查询分数
func FindSingleMark(query interface{}) (*ResMark, error) {
	var datas []*ResMark
	err := db.Where(query).Find(&datas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询mark:%s", err.Error())
	}
	if len(datas) > 0 {
		datas[0].StartDate = datas[0].StartDate[0:10]
		datas[0].EndDate = datas[0].EndDate[0:10]
		return datas[0], nil
	}
	return nil, nil
}

// FindAllMark FindAllMark
func FindAllMark(fields string, query interface{}, values ...interface{}) ([]*ResMark, error) {

	var datas []*ResMark
	if len(fields) == 0 {
		fields = "*"
	}
	err := db.Select(fields).Where(query, values...).Order("createTime desc").Find(&datas).Error
	// err := db.Raw(rawSQL, values...).Scan(&datas).Error
	return datas, err
}

// FindAllMarkBySQL FindAllMarkBySQL
func FindAllMarkBySQL(sql string) ([]*ResMark, error) {

	var datas []*ResMark

	err := db.Raw(sql).Scan(&datas).Error
	return datas, err
}

// FindAllMarkPaged FindAllMarkPaged
func FindAllMarkPaged(fields string, limit, offset int, order string, query interface{}) ([]*ResMark, error) {
	var datas []*ResMark
	if limit == 0 {
		limit = 20
	}
	if len(fields) == 0 {
		fields = "*"
	}
	if len(order) == 0 {
		order = "createTime desc"
	}
	err := db.Select(fields).Where(query).Limit(limit).Offset(offset).Order(order).Find(&datas).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return make([]*ResMark, 0), err
	}
	return datas, nil
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

// delMark 删除加减分
func delMark(query interface{}, args ...interface{}) error {
	if query == nil {
		return fmt.Errorf("query 不能为空")
	}
	return db.Where(query, args...).Delete(&ResMark{}).Error
}

// // DelMarkDelaySubmit 删除延迟提交的
// // 日期格式: 2020-01-12
// func DelMarkDelaySubmit(userID int,startDate string,endDate string) error{
// 	// 参数检查
// 	s,err:=util.ParseDate3(startDate)
// 	if err!=nil{
// 		return fmt.Errorf("startDate格式必须为'2020-01-12':%s",err.Error())
// 	}
// 	e,err:=util.ParseDate3(endDate)
// 	if err!=nil{
// 		return fmt.Errorf("endDate格式必须为'2020-01-12':%s",err.Error())
// 	}
// 	if s.Unix()-e.Unix()>0{
// 		return fmt.Errorf("startDate 不能大于 endDate")
// 	}
// 	if userID==0{
// 		return fmt.Errorf("userid 不能不为空")
// 	}
// 	// sql语句构造
// 	accordingly:="月度考核延迟提交产生扣分"
// 	delMark("startDate=? and endDate=? and userId=? and accordingly=?",startDate,endDate,userID,accordingly)

// }
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

// UpdateMarks 更新值
func UpdateMarks(query interface{}, values interface{}) error {
	return db.Model(&ResMark{}).Where(query).Updates(values).Error
}

// SumMarks 加减分合计
func SumMarks(startDate, endDate string, query interface{}, values ...interface{}) (string, error) {
	var mark ResMark

	err := db.Select("sum(markNumber) as markNumber").Where(query, values...).Where("startDate>=? and endDate<=?", startDate, endDate).
		Find(&mark).Error
	return mark.MarkNumber, err
}

// FindMarksRankPaged 分页查询加减分排行
func FindMarksRankPaged(limit, offset int, query interface{}) ([]*ResMark, error) {
	var datas []*ResMark
	err := db.Select("userId,username,ifnull(round(sum(markNumber),2),0) as markNumber").Where(query).
		Group("userId,username").Order("markNumber desc").
		Limit(limit).Offset(offset).Find(&datas).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return make([]*ResMark, 0), nil
	}
	return datas, err
}
