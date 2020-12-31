package conmgr

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// ImportPublicAssess 导入群众评议
func ImportPublicAssess(c *model.Container) error {
	return GetPublicAssessFromXlsx(c.File)
}

// ImportMarks 导入加减分
func ImportMarks(c *model.Container) error {

	return GetMarksFromXlsx(c.File)
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
	datas, err := transform2Csv(header, fields, ranks)
	str, _ := util.ToJSONStr(datas)
	println("csv:", str)
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
