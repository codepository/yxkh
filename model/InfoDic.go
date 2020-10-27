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
func FindAllInfoDic(query interface{}, values ...interface{}) ([]*InfoDic, error) {
	info := []*InfoDic{}
	err := db.Where(query, values...).Find(&info).Error
	return info, err
}
