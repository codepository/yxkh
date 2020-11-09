package model

import "github.com/jinzhu/gorm"

// InfoDic 字典
type InfoDic struct {
	Model
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
	Type2 string `json:"type2"`
}

// FindAllInfoDic FindAllInfoDic
func FindAllInfoDic(query interface{}, values ...interface{}) ([]*InfoDic, error) {
	info := []*InfoDic{}
	err := db.Where(query, values...).Find(&info).Error
	if err == gorm.ErrRecordNotFound {
		return make([]*InfoDic, 0), nil
	}
	return info, err
}
