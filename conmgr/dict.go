package conmgr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// FindallDict 字典查询
func FindallDict(c *model.Container) error {
	errstr := `{"body":"params":{"name":"评分依据","value":"","type":"","type2"}} 不能全为空`
	if len(c.Body.Params) == 0 {
		return errors.New(errstr)
	}
	var buf strings.Builder
	if c.Body.Params["name"] != nil {
		buf.WriteString(" and name='" + c.Body.Params["name"].(string) + "'")
	}
	if c.Body.Params["type"] != nil {
		buf.WriteString(" and type='" + c.Body.Params["type"].(string) + "'")
	}
	if c.Body.Params["type2"] != nil {
		buf.WriteString(" and type2='" + c.Body.Params["type2"].(string) + "'")
	}
	if c.Body.Params["value"] != nil {
		buf.WriteString(" and value like '%" + c.Body.Params["value"].(string) + "%'")
	}
	if buf.Len() == 0 {
		return errors.New(errstr)
	}
	data, err := model.FindAllInfoDic(buf.String()[4:])
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, data)
	return nil
}

// ScoreShare 查询半年、年度考核评分占比
func ScoreShare(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"departmentId":34}}部门ID`
	// 参数判断，必须有部门ID
	if len(c.Body.Params) == 0 || c.Body.Params["departmentId"] == nil {
		return fmt.Errorf(errstr)
	}
	did, err := util.Interface2Int(c.Body.Params["departmentId"])
	if err != nil {
		return err
	}
	// 根据部门id，查询部门属性，是采编经营类，还是行政后勤类
	attribute, err := FindDepartAttributeByID(did)
	if err != nil {
		return err
	}
	// println("部门属性：", attribute)
	name := "采编经营类"
	if attribute == 2 {
		name = "行政后勤类"
	}
	// 从info_dic表中查询相应的评分占比数据
	dics, err := model.FindAllInfoDic(map[string]interface{}{"name": name, "type2": "量化计分占比"})
	if err != nil {
		return fmt.Errorf("查询字典失败:%s", err.Error())
	}
	var vals []string
	var names []string
	for _, d := range dics {
		vals = append(vals, d.Value)
		names = append(names, d.Type)
	}
	c.Body.Data = append(c.Body.Data, vals)
	c.Body.Data = append(c.Body.Data, names)
	return nil
}
