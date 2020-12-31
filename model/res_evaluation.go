package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mumushuiding/util"
)

// ResEvaluationTableName ResEvaluationTableName
var ResEvaluationTableName = "res_evaluation"

// ResEvaluation 申请表内容
type ResEvaluation struct {
	EID int `gorm:"primary_key;column:eId" json:"eId,omitempty"`
	// UID 用户id
	UID        int    `gorm:"column:uid" json:"uid"`
	Username   string `json:"username"`
	Department string `json:"department"`
	// Position 用户填表时的职务
	Position           string `json:"position"`
	StartDate          string `gorm:"column:startDate" json:"startDate"`
	EndDate            string `gorm:"column:endDate" json:"endDate"`
	ProcessInstanceID  string `gorm:"column:processInstanceId" json:"processInstanceId,omitempty"`
	SelfEvaluation     string `gorm:"column:selfEvaluation" json:"selfEvaluation,omitempty"`
	Attendance         string `json:"attendance,omitempty"`
	OverseerEvaluation string `gorm:"column:overseerEvaluation" json:"overseerEvaluation,omitempty"`
	// 半年或全年考核 领导点评
	LeadershipEvaluation string `gorm:"column:leadershipEvaluation" json:"leadershipEvaluation,omitempty"`
	// 半年或全年考核 群众评议
	PublicEvaluation string `gorm:"column:publicEvaluation" json:"publicEvaluation,omitempty"`
	// 半年或全年考核 组织考核
	OrganizationEvaluation string `gorm:"column:organizationEvaluation" json:"organizationEvaluation,omitempty"`
	// 半年或全年考核 TotalMark 总分或考评总分
	TotalMark string `gorm:"column:totalMark" json:"totalMark,omitempty"`
	// 半年或全年考核 Marks 考效总分或考核量化分
	Marks             string    `gorm:"column:marks" json:"marks,omitempty"`
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

// EvaluationProcess EvaluationProcess
type EvaluationProcess struct {
	ResEvaluation
	ProcessBean Process `json:"processBean,omitempty"`
	Identity    string  `json:"identity,omiteymty"`
}

// GenerateTotal 重新计算并生成总分
// attribute 为部门属性：1-采编经营类,2-后勤行政类
// scoreshare 与info_dic表中type2=量化计分占比 相对应
// assessmentMarks 与info_dic表中type=基本定格对应得分 相对应
func (e *ResEvaluation) GenerateTotal(attribute int, assessmentMarks []*InfoDic) error {
	var leader, public, overseer, assess float64
	// 查询评分占比
	name := "采编经营类"
	if attribute == 2 {
		name = "行政后勤类"
	}
	// 从info_dic表中查询相应的评分占比数据
	scoreshare, err := FindAllInfoDic(map[string]interface{}{"name": name, "type2": "量化计分占比"})
	if err != nil {
		return fmt.Errorf("查询评分占比时报错:%s", err.Error())
	}

	// 领导点评
	if len(e.LeadershipEvaluation) == 0 {
		leader = 0.0
	} else {
		leader, err = getScore(e.LeadershipEvaluation, "领导", scoreshare, assessmentMarks)
		if err != nil {
			return err
		}

	}
	// 群众评议
	if len(e.PublicEvaluation) == 0 {
		public = 0.0
	} else {
		public, err = strconv.ParseFloat(e.PublicEvaluation, 64)
		if err != nil {
			return err
		}
		pshare, err := getScore("", "群众", scoreshare, assessmentMarks)
		if err != nil {
			return err
		}
		public = public * pshare
	}
	// 组织点评
	if len(e.OverseerEvaluation) == 0 {
		overseer = 0.0
	} else {
		overseer, err = getScore(e.OverseerEvaluation, "组织", scoreshare, assessmentMarks)
		if err != nil {
			return err
		}
	}
	// 考核量化分
	assess, err = strconv.ParseFloat(e.Marks, 64)
	if err != nil {
		return fmt.Errorf("考核量化分marks转化成数字报错:%s", err.Error())
	}
	ashare, err := getScore("", "考核量化", scoreshare, assessmentMarks)
	if err != nil {
		return err
	}
	e.TotalMark = fmt.Sprintf("%.2f", leader+public+overseer+assess*ashare)
	return nil
}

// getScore 获取分数
// assess 值有：优秀、合格、基本合格、不合格,如果值为空，则只返回分数占比
// role 值有：领导、组织、群众、考核量化
// scoreshare 为info_dic表中 type2=量化计分占比 所对应的结果
// assessment 为info_dic表中 type=基本定格对应得分 所对应的结果
func getScore(assess, role string, scoreshare []*InfoDic, assessment []*InfoDic) (float64, error) {
	// 转化对应的得分
	var mark, share float64
	var err error
	// 获取得分占比
	for _, s := range scoreshare {
		if strings.Contains(s.Type, role) {
			share, err = strconv.ParseFloat(s.Value, 64)
			if err != nil {
				return mark, fmt.Errorf("[%s]点评分数占比查询失败:%s，请通知管理员", role, err.Error())
			}
			break
		}
	}
	if share == 0.0 {
		return mark, fmt.Errorf("[%s]点评分数占比查询失败，请通知管理员,检查info_dic表，type2=量化计分占比", role)
	}
	if len(assess) == 0 {
		return share, nil
	}
	for _, a := range assessment {
		if assess == a.Name {
			mark, err = strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return mark, fmt.Errorf("基本定格[%s]转化成分数报错:%s", assess, err.Error())
			}
			break
		}
	}
	if mark == 0.0 {
		return mark, fmt.Errorf("基本定格[%s]未转化成分数，请通知管理员，检查info_dic表，type=基本定格对应得分", assess)
	}
	return mark * share, nil
}

// FromMap FromMap
func (e *ResEvaluation) FromMap(data map[string]interface{}) error {
	if data["startDate"] == nil || data["endDate"] == nil {
		return fmt.Errorf("startDate 和 endDate 不能为空")
	}
	s, _ := util.ToJSONStr(data)
	err := util.Str2Struct(s, e)
	if err != nil {
		return err
	}
	e.StartDate = data["startDate"].(string)
	e.EndDate = data["endDate"].(string)
	e.CreateTime = time.Now()
	return err
}

// FirstOrCreate 存在就更新，不存在就插入
func (e *ResEvaluation) FirstOrCreate() error {
	if len(e.ProcessInstanceID) == 0 {
		return fmt.Errorf("processInstanceId不能为空")
	}
	return db.Where(ResEvaluation{ProcessInstanceID: e.ProcessInstanceID}).Assign(e).FirstOrCreate(&ResEvaluation{}).Error
}

// FindSingleEvaluation2 查询流程，返回一条
func FindSingleEvaluation2(query interface{}, value ...interface{}) (*ResEvaluation, error) {
	var datas []*ResEvaluation
	err := db.Table(ResEvaluationTableName).Where(query, value...).Find(&datas).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if len(datas) > 0 {
		return datas[0], err
	}
	return nil, err
}

// FindSingleEvaluation 根据流程ID查询
func FindSingleEvaluation(processInstanceID interface{}) (*ResEvaluation, error) {
	var ep ResEvaluation
	err := db.Table(ResEvaluationTableName).Where("processInstanceId=?", processInstanceID).Find(&ep).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ep, nil
}

// DelEvaluationByProcessInstanceID 根据流程ID来删除
func DelEvaluationByProcessInstanceID(processInstanceID interface{}) error {
	return db.Where("processInstanceId=?", processInstanceID).Delete(&ResEvaluation{}).Error
}

// Updates 只会更新更改的和非空白字段
func (e *ResEvaluation) Updates() error {
	return db.Model(e).Where("eId=?", e.EID).Updates(e).Error
}

// FindAllEvaluation FindAllEvaluation
func FindAllEvaluation(fields, order string, limit, offset int, query interface{}) ([]*ResEvaluation, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	if len(order) == 0 {
		order = "createTime desc"
	}
	if limit == 0 {
		limit = 20
	}
	var datas []*ResEvaluation
	err := db.Select(fields).Where(query).Order(order).Limit(limit).Offset(offset).Find(&datas).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return make([]*ResEvaluation, 0), nil
	}
	return datas, err
}
