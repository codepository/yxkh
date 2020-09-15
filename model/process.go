package model

// ProcessTabelName ProcessTabelName
var ProcessTabelName = "process"

// Process 流程
type Process struct {
	ProcessInstanceID     string `gorm:"primary_key;column:processInstanceId" json:"processInstanceId,omitempty"`
	UserID                int    `gorm:"column:userId" json:"userId,omitempty"`
	ExecutionID           string `gorm:"column:executionId" json:"executionId,omitempty"`
	RequestedDate         string `gorm:"column:requestedDate" json:"requestedDate,omitempty"`
	Title                 string `gorm:"column:title" json:"title,omitempty"`
	BusinessType          string `gorm:"column:businessType" json:"businessType,omitempty"`
	DeploymentID          string `gorm:"column:deploymentId" json:"deploymentId,omitempty"`
	BusinessKey           string `gorm:"column:businessKey" json:"businessKey,omitempty"`
	Completed             string `gorm:"column:completed" json:"completed,omitempty"`
	Committed             string `gorm:"column:committed" json:"committed,omitempty"`
	DeptName              string `gorm:"column:deptName" json:"deptName,omitempty"`
	CurrentCandidateGroup string `gorm:"column:currentCandidateGroup" json:"currentCandidateGroup,omitempty"`
	Candidate             string `json:"candidate,omitempty"`
	Username              string `gorm:"column:username" json:"username,omitempty"`
}

// FindAllProcessPaged FindAllProcessPaged
func FindAllProcessPaged(limit, offset int, query interface{}) ([]*Process, int, error) {
	var datas []*Process
	var total int
	err := db.Table(ProcessTabelName).Where(query).Count(&total).Limit(limit).Offset(offset).
		Order("requestedDate desc").Find(&datas).Error
	return datas, total, err
}
