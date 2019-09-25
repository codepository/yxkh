package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mumushuiding/util"
	"github.com/codepository/yxkh/service"
)

// SaveProject SaveProject
func SaveProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持Post方法！！Only support Post ")
		return
	}
	var project = service.Project{}
	err := util.Body2Struct(r, &project)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	if len(project.EndDate) == 0 {
		util.ResponseErr(w, "字段 enddate 不能为空")
		return
	}
	if len(project.Progress) > 1000 {
		util.ResponseErr(w, "字段 progress 长度不能超过1000")
		return
	}
	if len(project.ProjectContent) > 1000 {
		util.ResponseErr(w, "字段 projectcontent 长度不能超过1000")
		return
	}
	if len(project.StartDate) == 0 {
		util.ResponseErr(w, "字段 startdate 不能为空")
		return
	}
	if project.UserID == 0 {
		util.ResponseErr(w, "字段 userid 不能为空")
		return
	}
	id, err := project.Save()
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.Response(w, fmt.Sprintf("%d", id), true)
}

// FindAllProject 分页查询
func FindAllProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持Post方法！！Only support Post ")
		return
	}
	var receiver = service.Project{}
	err := util.Body2Struct(r, &receiver)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	datas, err := receiver.FindAllPaged()
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	fmt.Fprintf(w, datas)
}

// DelProjectByID DelProjectByID
func DelProjectByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ResponseErr(w, "只支持GET方法")
		return
	}
	var id int
	var err error
	r.ParseForm()
	if len(r.Form["id"]) == 0 {
		util.ResponseErr(w, "字段 id 不能为空！")
		return
	}
	id, err = strconv.Atoi(r.Form["id"][0])
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	err = service.DelProjectByID(id)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseOk(w)
}

// UpdateProject 更新项目
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持POST方法！")
		return
	}
	var receiver = service.Project{}
	err := util.Body2Struct(r, &receiver)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	if receiver.ProjectID == 0 {
		util.ResponseErr(w, "字段 projectId 不能为空")
		return
	}
	if err = receiver.Update(); err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseOk(w)
}
