package conmgr

import (
	"errors"
	"fmt"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// FindAllProjectWithMarks 查询项目及分数
func FindAllProjectWithMarks(c *model.Container) error {
	errstr := `参数:{"body":{"params":{"userId":29,"startDate":"2020-01-07","endDate":"2020-01-31"}}} 三个参数不能为空`
	if len(c.Body.Params) < 3 {
		return fmt.Errorf(errstr)
	}
	if c.Body.Params["startDate"] == nil || c.Body.Params["endDate"] == nil || c.Body.Params["userId"] == nil {
		return fmt.Errorf(errstr)
	}
	endDate := fmt.Sprintf("%s 00:00:00", c.Body.Params["endDate"].(string))
	datas, err := model.FindProjectWithMarks("startDate>=? and endDate<=? and userId=?", c.Body.Params["startDate"], endDate, c.Body.Params["userId"])
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, datas)
	return nil
}

// UpdateProject 更新项目
func UpdateProject(c *model.Container) error {
	errstr := `参数格式：{"body":"params":{"data":{"projectId":2,"projectContent":"内容","progress":"进展"}}} `
	if len(c.Body.Params) == 0 || c.Body.Params["data"] == nil {
		return errors.New(errstr)
	}
	data, yes := c.Body.Params["data"].(map[string]interface{})

	if !yes {
		return errors.New(errstr)
	}
	if data["projectId"] == nil {
		return errors.New("projectId 不能为空")
	}
	id, err := util.Interface2Int(data["projectId"])
	if err != nil {
		return err
	}
	err = model.UpdatesProject(id, data)
	return err
}

// DelProject 删除项目
func DelProject(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"ids":[1,3]}} ids是要删除的评分`
	if len(c.Body.Params) == 0 || c.Body.Params["ids"] == nil {
		return errors.New(errstr)
	}
	var ids []int
	for _, id := range c.Body.Params["ids"].([]interface{}) {
		i, err := util.Interface2Int(id)
		if err != nil {
			return err
		}
		ids = append(ids, i)
	}
	err := model.DelProjectByIDs(ids)
	return err
}

// AddProject 添加项目
func AddProject(c *model.Container) error {
	errstr := `参数格式:{"body":{"params":{"data":[{"startDate":"2020-09-01","endDate":"2020-09-30","userId":12,"projectContent":"项目内容"}]}}}`
	if len(c.Body.Params) == 0 || c.Body.Params["data"] == nil {
		return fmt.Errorf(errstr)
	}
	datas, yes := c.Body.Params["data"].([]interface{})
	if !yes {
		return fmt.Errorf(errstr)
	}
	var ps []*model.ResProject
	for _, d := range datas {
		if d == nil {
			return fmt.Errorf(errstr)
		}
		json, yes := d.(map[string]interface{})
		if !yes {
			return fmt.Errorf(errstr)
		}
		p := &model.ResProject{}
		err := p.FromJSON(json)
		if err != nil {
			return err
		}
		ps = append(ps, p)
	}
	var ids []int
	for _, p := range ps {
		err := p.FirstOrCreate()
		if err != nil {
			return err
		}
		ids = append(ids, p.ProjectID)
	}
	c.Body.Data = append(c.Body.Data, ids)
	return nil

}

// CheckProjectByProcessInstanceID 使月度考核项目的评分生效
func CheckProjectByProcessInstanceID(processInstanceID string) (*ReExecData, error) {
	red := &ReExecData{Key: processInstanceID, FuncName: "CheckProjectByProcessInstanceID"}
	// 查询对应的月度考核
	e, err := model.FindSingleEvaluation(processInstanceID)
	if err != nil {
		return red, fmt.Errorf("查询月度考核失败:%s", err.Error())
	}
	query := map[string]interface{}{"startDate": util.FormatDate3(e.StartDate), "endDate": util.FormatDate3(e.EndDate), "userId": e.UID}
	values := map[string]interface{}{"checked": "1"}
	err = model.UpdateMarks(query, values)
	if err != nil {
		return red, fmt.Errorf("使加减分生效失败:%s", err.Error())
	}
	return nil, nil
}
