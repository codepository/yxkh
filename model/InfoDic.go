package model

// InfoDic 字典
type InfoDic struct {
	Model
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
	Type2 string `json:"type2"`
}

// FindAllInfoDic FindAllInfoDic
func FindAllInfoDic(fields map[string]interface{}) ([]*InfoDic, error) {
	info := []*InfoDic{}
	err := db.Where(fields).Find(&info).Error
	return info, err
}
