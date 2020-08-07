package controller

import (
	"net/http"

	"github.com/codepository/yxkh/service"
	"github.com/mumushuiding/util"
)

// FindAllInfoDic FindAllInfoDic
func FindAllInfoDic(w http.ResponseWriter, r *http.Request) {
	m, err := util.Body2MapWithDecode(r)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	fields := make(map[string]interface{})
	i := 0
	cols := []string{"type", "type2", "name"}
	for _, field := range cols {
		if len(m.Get(field)) > 0 {
			i++
			fields[field] = m.Get(field)
		}
	}
	if i == 0 {
		util.ResponseErr(w, "类型和名称不能全为空")
		return
	}
	data, err := service.FindAllInfoDicByMap(fields)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseData(w, data)
}
