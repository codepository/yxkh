package model

import (
	"fmt"
	"strings"
	"time"
)

// ResEvaluation 申请表内容
type ResEvaluation struct {
	EID                int       `gorm:"primary_key column:eId" json:"eId,omitempty"`
	StartDate          time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate            time.Time `gorm:"column:endDate" json:"endDate"`
	ProcessInstanceID  int       `gorm:"column:processInstanceId" json:"processInstanceId,omitempty"`
	SelfEvaluation     string    `gorm:"column:selfEvaluation" json:"selfEvaluation,omitempty"`
	Attendance         string    `json:"attendance,omitempty"`
	OverseerEvaluation string    `gorm:"column:overseerEvaluation" json:"overseerEvaluation,omitempty"`
	// 半年或全年考核 领导点评
	LeadershipEvaluation string `gorm:"column:leadershipEvaluation" json:"leadershipEvaluation,omitempty"`
	// 半年或全年考核 群众评议
	PublicEvaluation string `gorm:"column:publicEvaluation" json:"publicEvaluation,omitempty"`
	// 半年或全年考核 组件考核
	OrganizationEvaluation string `gorm:"column:organizationEvaluation" json:"organizationEvaluation,omitempty"`
	// 半年或全年考核 TotalMark 总分
	TotalMark string `gorm:"column:totalMark" json:"totalMark,omitempty"`
	// 半年或全年考核 Marks 考效总分
	Marks             string    `json:"marks,omitempty"`
	EvaluationType    string    `gorm:"column:evaluationType" json:"evaluationType,omitempty"`
	CreateTime        time.Time `gorm:"column:createTime" json:"createTime"`
	ShortComesAndPlan string    `gorm:"column:shortComesAndPlan" json:"shortComesAndPlan,omitempty"`
	Sparation         string    `json:"sparation,omitempty"`
	LeadershipRemark  string    `gorm:"column:leadershipRemark" json:"leadershipRemark,omitempty"`
	Committed         string    `json:"committed,omitempty"`
	OverseerRemark    string    `gorm:"column:overseerRemark" json:"overseerRemark,omitempty"`
	PublicRemark      string    `gorm:"column:publicRemark" json:"publicRemark,omitempty"`
	// 半年或全年考核 考核结果
	Result string `json:"result,omitempty"`
}

// EvaluationUser EvaluationUser
type EvaluationUser struct {
	ResEvaluation
	UserID   int    `gorm:"column:userId" json:"userId,omitempty"`
	Username string `json:"username,omitempty"`
	DeptName string `gorm:"column:deptName" json:"deptName,omitempty"`
}

// FindAllEvaluationPaged 查询所有的申请表
// select res_evaluation.processInstanceId,process.userId,process.username,process.deptName,res_evaluation.marks from res_evaluation join process on process.businessType like '一线干部-半年考核' and process.title like '%2020年上半年-半年考核%' and res_evaluation.processInstanceId=process.processInstanceId  order by res_evaluation.totalMark+0 desc,res_evaluation.marks+0 desc limit 10;
func FindAllEvaluationPaged(fields string, limit, offset int, order string, joins string, query interface{}, values ...interface{}) ([]*EvaluationUser, int, error) {
	var total int
	var datas []*EvaluationUser
	if len(fields) == 0 {
		fields = "*"
	}
	if len(order) == 0 {
		order = "res_evaluation.totalMark+0 desc,res_evaluation.marks+0 desc"
	}
	if limit == 0 {
		limit = 10
	}
	if len(joins) == 0 {
		joins = "join process on res_evaluation.processInstanceId=process.processInstanceId"
	}
	err := db.Table("res_evaluation").Select(fields).Joins(joins).Where(query, values...).
		Count(&total).Order(order).Limit(limit).Offset(offset).
		Find(&datas).Error
	if err != nil {
		return nil, 0, nil
	}
	return datas, total, nil
}

// FindAllEvaluationPagedByType 查询所有的申请表
// typename 与 数据库 Sparation 字段对应
// select res_evaluation.processInstanceId,process.userId,process.username,process.deptName,res_evaluation.marks from res_evaluation join process on process.businessType like '一线干部-半年考核' and process.title like '%2020年上半年-半年考核%' and res_evaluation.processInstanceId=process.processInstanceId  order by res_evaluation.totalMark+0 desc,res_evaluation.marks+0 desc limit 10;
func FindAllEvaluationPagedByType(fields string, limit, offset int, order string, typename string) ([]*EvaluationUser, int, error) {
	return FindAllEvaluationPaged(fields, limit, offset, order, "", "sparation=?", typename)
}

// FindAllEvaluationPagedByTypeAndUserids 查询所有的申请表
func FindAllEvaluationPagedByTypeAndUserids(fields string, limit, offset int, order string, typename string, userids []interface{}) ([]*EvaluationUser, int, error) {
	var idbuffer strings.Builder
	for _, id := range userids {
		idbuffer.WriteString(fmt.Sprintf(",%v", id))
	}
	joins := fmt.Sprintf("join process on res_evaluation.processInstanceId=process.processInstanceId and userId in (%s)", idbuffer.String()[1:])
	return FindAllEvaluationPaged(fields, limit, offset, order, joins, "sparation=?", typename)
}

// FindAllEvaluationPagedByTypeAndUsername 查询所有的申请表
// typename 与 数据库 Sparation 字段对应
func FindAllEvaluationPagedByTypeAndUsername(fields string, limit, offset int, order string, typename, username string) ([]*EvaluationUser, int, error) {
	return FindAllEvaluationPaged(fields, limit, offset, order, "", "sparation=? and username=?", typename, username)
}
