package conmgr

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// ImportPublicAssess 导入群众评议
func ImportPublicAssess(c *model.Container) error {
	return GetPublicAssessFromXlsx(c.File)
}

// ImportMarks 导入加减分
func ImportMarks(c *model.Container) error {
	var checked = "0"
	if len(c.PostParams.Get("checked")) > 0 {
		checked = c.PostParams.Get("checked")
	}
	return GetMarksFromXlsx(c.File, checked)
}

// ExportAllEvalutionRank 导出半年、全年考核
func ExportAllEvalutionRank(c *model.Container) error {
	header := c.Body.Data[0].([]interface{})
	fields := c.Body.Data[1].([]interface{})
	c.Body.Data = c.Body.Data[:0]
	err := FindAllEvalutionRank(c)
	if err != nil {
		return err
	}
	if len(c.Body.Data) == 0 || c.Body.Data[0] == nil {
		return nil
	}
	ranks := c.Body.Data[0].([]*model.ResEvaluation)
	// 重新计算分数
	// 先查询基础分
	baseStr, err := model.FindMarksBase()
	if err != nil {
		return err
	}
	basemarks, err := strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return err
	}
	// 优秀、合格 得分
	assessmentMarks, err := model.FindAllInfoDic(map[string]interface{}{"type": "基本定格对应得分"})
	if err != nil {
		return fmt.Errorf("基本定格对应得分时报错:%s", err.Error())
	}
	// 部门属性表
	attMap := make(map[string]int)
	for _, e := range ranks {
		// 考核量化分重新计算
		total, err := model.SumMarks(e.StartDate, e.EndDate, map[string]interface{}{"userId": e.UID, "checked": 1})
		if err != nil {
			return fmt.Errorf("重新计算考核量化分:%s", err.Error())
		}
		total2, err := strconv.ParseFloat(total, 32)
		if err != nil {
			return fmt.Errorf("考核量化分转字符:%s", err.Error())
		}
		mark := fmt.Sprintf("%.2f", total2+basemarks)
		// 判断加减分是否是最新的
		if mark != e.Marks {
			// 更新加减分和总分
			e.Marks = mark
			// 查询用户所在部门的经营属性
			var attriute = attMap[e.Department]
			if attriute == 0 {
				attriute, err = FindDepartAttribute(map[string]interface{}{"name": e.Department})
				if err != nil {
					return fmt.Errorf("用户[%s]查询用户部门属性时失败:%s,请联系管理员", e.Username, err.Error())
				}
			}
			err = e.GenerateTotal(attriute, assessmentMarks)
			if err != nil {
				return fmt.Errorf("用户[%s]计算总分失败:%s,请稍后再试", e.Username, err.Error())
			}
			err = e.Updates()
			if err != nil {
				return fmt.Errorf("用户[%s]更新数据失败:%s,请稍后再试", e.Username, err.Error())
			}
		}
	}
	datas, err := transform2Csv(header, fields, ranks)
	c.Body.Data = c.Body.Data[:0]
	c.Body.Data = append(c.Body.Data, datas...)
	return nil
}
func transform2Csv(header []interface{}, fields []interface{}, datas interface{}) ([]interface{}, error) {
	if datas == nil {
		return []interface{}{}, nil
	}
	if len(header) == 0 || len(fields) == 0 {
		return nil, errors.New("数据转换成csv格式时,行标和字段名不能为空")
	}
	// 获取结果集
	s := reflect.ValueOf(datas)
	var result []interface{}
	// 遍历结果
	for i := 0; i < s.Len(); i++ {
		var row []string
		for _, f := range fields {

			item := s.Index(i)
			str, _ := util.ToJSONStr(item.Interface())
			data, err := util.Str2Map(str)
			if err != nil {
				return []interface{}{}, nil
			}
			value := data[f.(string)]
			if value != nil {
				row = append(row, fmt.Sprintf("%v", value))
			} else {
				row = append(row, fmt.Sprintf("%s", ""))

			}

		}
		result = append(result, row)
	}
	return result, nil
}
