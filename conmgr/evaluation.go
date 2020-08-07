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
	fields := "process.userId,process.username,process.deptName,res_evaluation.marks,publicEvaluation,leadershipEvaluation,overseerEvaluation,totalMark,result"
	c.Body.Order = "totalMark+0 desc,marks+0 desc"

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
