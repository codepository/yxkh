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
	if len(c.Body.Params) == 0 {
		return fmt.Errorf(errstr)
	}
	// 根据部门id，查询部门属性，是采编经营类，还是行政后勤类
	attribute, err := FindDepartAttribute(c.Body.Params)
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
		return fmt.Errorf("查询量化计分占比失败:%s", err.Error())
	}
	var vals []string
	var names []string
	for _, d := range dics {
		vals = append(vals, d.Value)
		names = append(names, d.Type)
	}
	c.Body.Data = append(c.Body.Data, vals)
	c.Body.Data = append(c.Body.Data, names)
	// 优秀、合格 得分
	dics2, err := model.FindAllInfoDic(map[string]interface{}{"type": "基本定格对应得分"})
	if err != nil {
		return fmt.Errorf("查询基本定格对应得分失败:%s", err.Error())
	}
	var vals2 []string
	var names2 []string
	for _, d := range dics2 {
		vals2 = append(vals2, d.Value)
		names2 = append(names2, d.Name)
	}
	c.Body.Data = append(c.Body.Data, vals2)
	c.Body.Data = append(c.Body.Data, names2)
	return nil
}

// UpdateDict 更新日志
func UpdateDict(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"data":[{"id":1}]}},id不可或缺,注意字段中值为"",0和false的字段不更新`
	if len(c.Body.Params) == 0 || c.Body.Params["data"] == nil {
		return fmt.Errorf(errstr)
	}
	datas, ok := c.Body.Params["data"].([]interface{})
	if !ok {
		return fmt.Errorf("更新字典:data 必须为数组")
	}
	var deptMaps []map[string]interface{}
	for _, d := range datas {
		deptMap, ok := d.(map[string]interface{})
		if !ok {
			return fmt.Errorf(errstr)
		}
		deptMaps = append(deptMaps, deptMap)
	}
	var depts []*model.InfoDic
	for _, m := range deptMaps {
		var dept model.InfoDic
		str, _ := util.ToJSONStr(m)
		err := util.Str2Struct(str, &dept)
		if err != nil {
			return fmt.Errorf("更新字典:%s", err.Error())
		}
		depts = append(depts, &dept)

	}
	for _, dept := range depts {
		if dept.ID == 0 {
			return fmt.Errorf("更新字典:字典的id不能为空")
		}
		err := dept.Updates()
		if err != nil {
			return fmt.Errorf("更新字典:%s", err.Error())
		}
	}
	return nil
}

// AddDict 添加日志
func AddDict(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"data":[{"name":"优秀","value":"3","type":"","type2":""}]}},注意字段中值为"",0和false的字段不更新`
	if len(c.Body.Params) == 0 || c.Body.Params["data"] == nil {
		return fmt.Errorf(errstr)
	}
	datas, ok := c.Body.Params["data"].([]interface{})
	if !ok {
		return fmt.Errorf("添加字典:data 必须为数组")
	}
	ids, err := AddDictBatch(datas)
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, ids)
	return nil
}

// AddDictBatch 批量添加，返回添加成功的id
func AddDictBatch(datas []interface{}) ([]int, error) {
	errstr := `参数格式:{"body":"params":{"data":[{"name":"优秀","value":"3","type":"","type2":""}]}},注意字段中值为"",0和false的字段不更新`
	var ids []int
	var deptMaps []map[string]interface{}
	for _, d := range datas {
		deptMap, ok := d.(map[string]interface{})
		if !ok {
			return ids, fmt.Errorf(errstr)
		}
		deptMaps = append(deptMaps, deptMap)
	}
	var depts []*model.InfoDic
	for _, m := range deptMaps {
		var dept model.InfoDic
		str, _ := util.ToJSONStr(m)
		err := util.Str2Struct(str, &dept)
		if err != nil {
			return ids, fmt.Errorf("添加字典:%s", err.Error())
		}
		depts = append(depts, &dept)

	}
	for _, dept := range depts {
		err := dept.FirstOrCreate()
		if err != nil {
			return ids, fmt.Errorf("添加字典:%s", err.Error())
		}
		ids = append(ids, dept.ID)
	}
	return ids, nil
}

// DelDict 删除日志
func DelDict(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"data":[1,2,3]}} 数字为字典id`
	if len(c.Body.Params) == 0 || c.Body.Params["data"] == nil {
		return fmt.Errorf(errstr)
	}

	datas, ok := c.Body.Params["data"].([]interface{})
	if !ok {
		return fmt.Errorf("data必须为整数数组")
	}
	var ids []int
	for _, d := range datas {
		r, err := util.Interface2Int(d)
		if err != nil {
			return err
		}
		ids = append(ids, r)
	}
	return model.DeleteDicsIDs(ids)
}
