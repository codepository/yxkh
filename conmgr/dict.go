package conmgr

import (
	"errors"
	"strings"

	"github.com/codepository/yxkh/model"
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
