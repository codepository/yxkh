package conmgr

import (
	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

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
	ranks := c.Body.Data[0].([]*model.EvaluationUser)
	datas, err := util.Transform2Csv(header, fields, ranks)
	c.Body.Data = datas
	return nil
}
