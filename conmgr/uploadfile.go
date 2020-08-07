package conmgr

import "github.com/codepository/yxkh/model"

// FindAllUploadfiles FindAllUploadfiles
func FindAllUploadfiles(query interface{}, limit, offset int) ([]*model.Uploadfile, error) {
	return model.FindAllUploadFiles(query, limit, offset)
}
