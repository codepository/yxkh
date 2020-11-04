package model

import "github.com/mumushuiding/util"

// ProcessTabelName ProcessTabelName
var ProcessTabelName = "process"

// Process 流程
type Process struct {
	ProcessInstanceID string `gorm:"primary_key;column:processInstanceId" json:"processInstanceId,omitempty"`
	// UID 对应用户ID
	UID int `gorm:"column:uid" json:"uid,omitempty"`
	// UserID 对应微信ID
	UserID        string `gorm:"column:userId" json:"userId,omitempty"`
	RequestedDate string `gorm:"column:requestedDate" json:"requestedDate,omitempty"`
	Title         string `gorm:"column:title" json:"title,omitempty"`
	BusinessType  string `gorm:"column:businessType" json:"businessType,omitempty"`
	Completed     int    `gorm:"column:completed" json:"completed"`
	DeptName      string `gorm:"column:deptName" json:"deptName,omitempty"`
	Candidate     string `json:"candidate,omitempty"`
	Username      string `gorm:"column:username" json:"username,omitempty"`
	DeploymentID  string `gorm:"column:deploymentId" json:"deploymentId,omitempty"`
	// step 当前执行步骤
	Step int `json:"step"`
}

// ProcessErrLogData 错误数据日志
type ProcessErrLogData struct {
	ProcessInstanceID string `gorm:"primary_key;column:processInstanceId" json:"processInstanceId,omitempty"`
	UID               int    `gorm:"column:uid" json:"uid,omitempty"`
	// UserID 对应微信ID
	UserID       string `gorm:"column:userId" json:"userId,omitempty"`
	BusinessType string `gorm:"column:businessType" json:"businessType,omitempty"`
}

// toString toString
func (p *ProcessErrLogData) toString() string {
	str, _ := util.ToJSONStr(p)
	return str
}
