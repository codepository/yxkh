package model

// Uploadfile Uploadfile
type Uploadfile struct {
	ID int `gorm:"primary_key" json:"ID,omitempty"`
	Filetype string `json:"filetype,omitempty"`
	Filename string `json:"filename,omitempty"`
	Filepath string `json:"filepath,omitempty"`
}

// FindAllUploadFiles 查询上传文件
func FindAllUploadFiles(query interface{}, limit, offset int) ([]*Uploadfile, error) {
	var datas []*Uploadfile
	if limit == 0 {
		limit = 50
	}
	err := db.Where(query).Order("createdDate desc").Limit(limit).Offset(offset).Find(&datas).Error
	return datas, err
}
