package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mumushuiding/util"
)

// ResProject 项目
type ResProject struct {
	ProjectID int `gorm:"primary_key;column:projectId" json:"projectId"`
	// 一个中文长度为：3，英文为： 1
	ProjectContent string     `gorm:"size:1000;column:projectContent" json:"projectContent"`
	UserID         int        `gorm:"column:userId" json:"userId"`
	StartDate      time.Time  `gorm:"column:startDate" json:"startDate"`
	EndDate        time.Time  `gorm:"column:endDate" json:"endDate"`
	Progress       string     `gorm:"size:1000" json:"progress"`
	Createtime     time.Time  `gorm:"column:createTime" json:"createTime"`
	Marks          []*ResMark `gorm:"FOREIGNKEY:res_mark_projectId_res_project_projectId_foreign;ASSOCIATION_FOREIGNKEY:projectId" json:"marks"`
}

// FromJSON FromJSON
func (p *ResProject) FromJSON(json map[string]interface{}) error {
	if json["projectId"] != nil {
		id, err := util.Interface2Int(json["projectId"])
		if err != nil {
			return err
		}
		p.ProjectID = id
	}
	if json["progress"] != nil {
		p.Progress = json["progress"].(string)
	}
	if json["endDate"] != nil {
		end, err := util.ParseDate3(json["endDate"].(string))
		if err != nil {
			return err
		}
		p.EndDate = end
	}
	if json["startDate"] != nil {
		start, err := util.ParseDate3(json["startDate"].(string))
		if err != nil {
			return err
		}
		p.StartDate = start
	}
	if json["projectContent"] != nil {
		p.ProjectContent = json["projectContent"].(string)
	}
	if json["userId"] != nil {
		uid, err := util.Interface2Int(json["userId"])
		if err != nil {
			return err
		}
		p.UserID = uid
	}
	return nil
}

// FirstOrCreate 不存在就创建
func (p *ResProject) FirstOrCreate() error {
	// 检查参数
	if len(p.ProjectContent) == 0 {
		return errors.New("projectContent 项目内容不能为空")
	}
	if p.UserID == 0 {
		return errors.New("userId 用户id不能为空")
	}
	if p.StartDate.IsZero() {
		return errors.New("startDate 不能为空")
	}
	if p.EndDate.IsZero() {
		return errors.New("endDate 不能为空")
	}
	p.Createtime = time.Now()
	return db.Where(ResProject{
		UserID: p.UserID, StartDate: p.StartDate, EndDate: p.EndDate,
		ProjectContent: p.ProjectContent,
	}).Assign(p).FirstOrCreate(p).Error
}

// UpdatesProject 只更新更改的字段
func UpdatesProject(id int, params map[string]interface{}) error {
	if params["startDate"] != nil {
		start, err := util.ParseDate3(params["startDate"].(string))
		if err != nil {
			return err
		}
		params["startDate"] = start
	}
	if params["endDate"] != nil {
		end, err := util.ParseDate3(params["endDate"].(string))
		if err != nil {
			return err
		}
		params["endDate"] = end
	}
	return db.Model(&ResProject{ProjectID: id}).Updates(params).Error
}

// FindProjectWithMarks 查询项目和分数
func FindProjectWithMarks(query interface{}, values ...interface{}) ([]*ResProject, error) {
	var projects []*ResProject
	// var marks []*ResMark
	err := db.Where(query, values...).Find(&projects).Error
	if err != nil {
		return nil, err
	}
	if len(projects) > 0 {
		for _, p := range projects {
			marks, err := FindAllMark("projectId=?", p.ProjectID)
			if err != nil {
				return nil, err
			}
			p.Marks = marks
		}
	}
	return projects, nil
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

// DelProjectByIDs 根据id删除项目
func DelProjectByIDs(id []int) error {
	return db.Where(id).Delete(&ResProject{}).Error
}

// Update Update
func (p *ResProject) Update() error {
	return db.Model(&ResProject{}).Updates(p).Error
}

// AddProjectWithMark 添加项目并评分
// startDate、endDate、userID 必不可少
func AddProjectWithMark(startDate, endDate time.Time, projectContent string, userID int, checked, markNumber, markReason string) error {
	// 先添加项目，如果已经存在就不添加
	p := &ResProject{
		ProjectContent: projectContent,
		StartDate:      startDate,
		EndDate:        endDate,
		UserID:         userID,
	}
	err := p.FirstOrCreate()
	if err != nil {
		return fmt.Errorf("添加项目失败:%s", err.Error())
	}
	// 为项目添加评分，如果已经存在就不添加
	mark := &ResMark{
		ProjectID:  p.ProjectID,
		UserID:     userID,
		Checked:    checked,
		MarkReason: markReason,
		MarkNumber: markNumber,
		StartDate:  startDate,
		EndDate:    endDate,
	}
	err = mark.FirstOrCreate()
	if err != nil {
		return fmt.Errorf("添加分数失败:%s", err.Error())
	}
	return nil
}
