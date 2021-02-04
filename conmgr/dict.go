package conmgr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// FindallDict 字典查询
func FindallDict(c *model.Container) error {
	errstr := `{"body":"params":{"limit":20,"name":"评分依据","value":"","type":"","type2"}} 不能全为空`
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
	if c.Body.Params["limit"] == nil {
		data, err := model.FindAllInfoDic(buf.String()[4:])
		if err != nil {
			return err
		}
		c.Body.Data = append(c.Body.Data, data)
	} else {
		limit, err := util.Interface2Int(c.Body.Params["limit"])
		if err != nil {
			return err
		}
		offset := 0
		if c.Body.Params["offset"] != nil {
			offset, err = util.Interface2Int(c.Body.Params["offset"])
			if err != nil {
				return err
			}
		}
		data, err := model.FindAllInfoDicPaged(buf.String()[4:], limit, offset)
		if err != nil {
			return err
		}
		c.Body.Data = append(c.Body.Data, data)
	}
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

// FindStartDayOfAutoDedcutOfMonthProcess 月度考核开始扣分日期
func FindStartDayOfAutoDedcutOfMonthProcess() (int, error) {
	start, err := model.FindSingleDict("name='延迟提交扣分日期'")
	if err != nil {
		return 0, fmt.Errorf("查询info_dic表name为[延迟提交扣分日期]的值:%s", err.Error())
	}
	s, err := strconv.Atoi(start.Value)
	if err != nil {
		return 0, fmt.Errorf("info_dic表name为[延迟提交扣分日期]的值转数字:%s", err.Error())
	}
	return s, nil

}

// FindEndDayOfAutoDedcutOfMonthProcess 月度考核结束扣分日期
func FindEndDayOfAutoDedcutOfMonthProcess() (int, error) {
	start, err := model.FindSingleDict("name='限期未交扣分日期'")
	if err != nil {
		return 0, fmt.Errorf("查询info_dic表name为[延迟提交扣分日期]的值:%s", err.Error())
	}
	s, err := strconv.Atoi(start.Value)
	if err != nil {
		return 0, fmt.Errorf("info_dic表name为[限期未交扣分日期]的值转数字:%s", err.Error())
	}
	return s, nil

}

// ExportMarksPriciple 导出加减分规则
func ExportMarksPriciple(c *model.Container) error {
	header := c.Body.Data[0].([]interface{})
	fields := c.Body.Data[1].([]interface{})
	// 查询数据
	datas, err := model.FindAllInfoDic("name='评分依据'")
	if err != nil {
		return fmt.Errorf("导出加减分：%s", err.Error())
	}
	// 将数据转换成csv格式
	result, err := util.Transform2Csv(header, fields, datas)
	if err != nil {
		return fmt.Errorf("数据转换成csv:%s", err.Error())
	}
	c.Body.Data = c.Body.Data[:0]
	c.Body.Data = append(c.Body.Data, result...)
	return nil
}

// ImportMarksPriciple 导入加减分
func ImportMarksPriciple(c *model.Container) error {
	var buff strings.Builder
	var datas [][]string
	var err error
	ss := strings.Split(c.File.Name(), ".")
	switch ss[len(ss)-1] {
	case "xlsx":
		datas, err = GetDatasFromXlsx(c.File)
		break
	case "csv":
		datas, err = GetDatasFromCSV(c.File)
		break
	default:
		return fmt.Errorf("暂时只支持xlsx和csv后缀的文件")
	}
	if err != nil {
		return fmt.Errorf("解析xlsx或csv文件:%s", err.Error())
	}
	// 验证行头,第1列是ID,第2列是“修改前",第3列是"修改后"
	if len(datas) < 2 {
		return fmt.Errorf("确认导入数据是否为空")
	}
	for i, row := range datas[1:] {
		// 若第1列ID为空，判断是否已经存在，不存在就保存
		if len(row[0]) == 0 {
			dic := &model.InfoDic{
				Name:  "评分依据",
				Value: row[1],
			}
			err := dic.FirstOrCreate()
			if err != nil {
				buff.WriteString(fmt.Sprintf("第[%d]行:%s", i+2, err.Error()))
			}
		} else if len(row[2]) > 0 { // 获取修改后不为空的数据

			dic := &model.InfoDic{
				Value: row[2],
			}

			id, err := strconv.Atoi(row[0])
			if err != nil {
				buff.WriteString(fmt.Sprintf("第[%d]行:%s", i+2, err.Error()))
			}
			dic.ID = id
			err = dic.Updates()
			if err != nil {
				buff.WriteString(fmt.Sprintf("第[%d]行:%s", i+2, err.Error()))
			}

		}
	}
	if buff.Len() > 0 {
		return fmt.Errorf("以下导入失败:\n%s", buff.String())
	}
	return nil

}
