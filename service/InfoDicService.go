package service

import (
	"github.com/codepository/yxkh/model"
	"github.com/jinzhu/gorm"
	"github.com/mumushuiding/util"
)

// InfoDicReceiver InfoDicReceiver
type InfoDicReceiver struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Type2 string `json:"type2"`
}

// FindAllInfoDicByMap FindAllInfoDicByMap
func FindAllInfoDicByMap(fields map[string]interface{}) (string, error) {
	result, err := model.FindAllInfoDic(fields)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	return util.ToJSONStr(result)
}

// FindAllInfoDic FindAllInfoDic
func (i InfoDicReceiver) FindAllInfoDic() (string, error) {
	fields := map[string]interface{}{}
	if len(i.Type) > 0 {
		fields["type"] = i.Type
	}
	if len(i.Type2) > 0 {
		fields["type2"] = i.Type2
	}
	if len(i.Name) > 0 {
		fields["name"] = i.Name
	}
	result, err := model.FindAllInfoDic(fields)
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	return util.ToJSONStr(result)
}
