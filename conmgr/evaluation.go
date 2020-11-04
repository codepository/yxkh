package conmgr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codepository/yxkh/model"
)

// FindAllEvalution 查询所有申请表
func FindAllEvalution(c *model.Container) error {
	errStr := fmt.Sprintf(`参数格式:{"body":{"fields":["eId,marks","2020年上半年-半年考核"],"max_results":10,"start_index":0,"order":"marks desc"}}`)
	if c.Body.Fields == nil || len(c.Body.Fields) < 2 {
		return errors.New(errStr)
	}

	e, total, err := model.FindAllEvaluationPagedByType(c.Body.Fields[0], c.Body.MaxResults, c.Body.StartIndex, c.Body.Order, c.Body.Fields[1])
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, e)
	c.Body.Total = total
	return nil
}

// FindAllEvalutionRank 半年和年度考核排行
func FindAllEvalutionRank(c *model.Container) error {
	errStr := fmt.Sprintf(`参数格式:{"body":{"fields":["2020年上半年-半年考核"],"max_results":10,"start_index":0,"order":"marks desc","metrics":"第一考核组成员,第二考核组成员","username":"小明"}} metrics为用户标签,可以为空;username可以为空`)
	if c.Body.Fields == nil || len(c.Body.Fields) < 1 {
		return errors.New(errStr)
	}
	if len(c.Body.Fields[0]) == 0 {
		return errors.New(errStr)
	}
	fields := "process.userId,process.username,process.deptName,res_evaluation.marks,publicEvaluation,leadershipEvaluation,overseerEvaluation,totalMark,result,startDate,endDate"
	c.Body.Order = "totalMark+0 desc,marks+0 desc"
	c.Body.Data = c.Body.Data[:0]
	if len(c.Body.UserName) != 0 { // 用户名不为空
		e, total, err := model.FindAllEvaluationPagedByTypeAndUsername(fields, c.Body.MaxResults, c.Body.StartIndex, c.Body.Order, c.Body.Fields[0], c.Body.UserName)
		if err != nil {
			return err
		}
		c.Body.Data = append(c.Body.Data, e)
		c.Body.Total = total
	} else if len(c.Body.Metrics) != 0 { // 即用户标签不为空
		userids, err := FindUseridsByTags(strings.Split(c.Body.Metrics, ","), "")
		if len(userids) == 0 {
			c.Body.Data = append(c.Body.Data, []interface{}{})
			c.Body.Total = 0
			return nil
		}
		if err != nil {
			return fmt.Errorf("根据标签【%s】查询用户id时报错:%s", c.Body.Metrics, err.Error())
		}

		e, total, err := model.FindAllEvaluationPagedByTypeAndUserids(fields, c.Body.MaxResults, c.Body.StartIndex, c.Body.Order, c.Body.Fields[0], userids)
		if err != nil {
			return err
		}
		c.Body.Data = append(c.Body.Data, e)
		c.Body.Total = total

	} else {
		e, total, err := model.FindAllEvaluationPagedByType(fields, c.Body.MaxResults, c.Body.StartIndex, c.Body.Order, c.Body.Fields[0])
		if err != nil {
			return err
		}
		c.Body.Data = append(c.Body.Data, e)
		c.Body.Total = total
	}

	return nil
}

// RemarkEvaluationByProcessInstanceID 根据流程检查月度考核评价，并根据评价添加加减分
func RemarkEvaluationByProcessInstanceID(processInstanceID string) (*ReExecData, error) {
	red := &ReExecData{Key: processInstanceID, FuncName: "RemarkEvaluationByProcessInstanceID"}
	// 查询对应的月度考核
	e, err := model.FindSingleEvaluation(processInstanceID)
	if err != nil {
		return red, fmt.Errorf("查询月度考核失败:%s", err.Error())
	}
	// 查询 organizationEvaluation 字段所对应的加减分
	dicName := ""
	switch e.OrganizationEvaluation {
	case "优秀":
		dicName = "月考评优"
		break
	case "基本合格":
		dicName = "月考基本合格"
		break
	case "不合格":
		dicName = "月考不合格"
		break
	default:
		return red, fmt.Errorf("月度考核评价【%s】不存在，请务必联系管理员", e.OverseerEvaluation)
	}
	dic, err := model.FindAllInfoDic(map[string]interface{}{"type": "月考自动加减分", "name": dicName})
	if err != nil || len(dic) == 0 {
		return red, fmt.Errorf("查询字典【%s】失败:%s", dicName, err.Error())
	}
	// 添加加减分
	err = model.AddProjectWithMark(e.StartDate, e.EndDate, "系统导入", e.UID, "1", dic[0].Value, dicName)
	if err != nil {
		return red, fmt.Errorf("月度考核自动加分失败:%s", err.Error())
	}
	return nil, nil
}
