package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/codepository/yxkh/service"

	"github.com/mumushuiding/util"
)

// SaveMark SaveMark
func SaveMark(w http.ResponseWriter, r *http.Request) {
	var receiver = service.MarkReceiver{}
	if err := util.GetBody(&receiver, w, r); err != nil {
		util.ResponseErr(w, err)
		return
	}
	if receiver.ProjectID == 0 {
		util.ResponseErr(w, "字段 projectId 不能为空")
		return
	}
	if len(receiver.UserID) == 0 {
		util.ResponseErr(w, "字段 userId 不能为空")
		return
	}
	if len(receiver.MarkNumber) == 0 {
		util.ResponseErr(w, "字段 markNumber 不能为空")
		return
	}
	if len(receiver.StartDate) == 0 {
		util.ResponseErr(w, "字段 startDate 不能为空")
		return
	}
	if len(receiver.EndDate) == 0 {
		util.ResponseErr(w, "字段 endDate 不能为空")
		return
	}
	id, err := receiver.Save()
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.Response(w, fmt.Sprintf("%d", id), true)
}

// FindAllMark FindAllMark
func FindAllMark(w http.ResponseWriter, r *http.Request) {
	var receiver = service.MarkReceiver{}
	if err := util.GetBody(&receiver, w, r); err != nil {
		util.ResponseErr(w, err)
		return
	}
	data, err := receiver.FindAll()
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	fmt.Fprintf(w, data)
}

// DelMark DelMark
func DelMark(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ResponseNo(w, "只支持GET方法")
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
	if err = service.DelMarkByID(id); err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseOk(w)

}

// UpdateMark UpdateMark
func UpdateMark(w http.ResponseWriter, r *http.Request) {
	var receiver = service.MarkReceiver{}
	if err := util.GetBody(&receiver, w, r); err != nil {
		util.ResponseErr(w, err)
		return
	}
	if receiver.MarkID == 0 {
		util.ResponseErr(w, "字段 markId 不能为空")
		return
	}
	err := receiver.Update()
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseOk(w)
}

// FindMarkRankForHome 首页加减分累计排行数据
func FindMarkRankForHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		util.ResponseNo(w, "只支持GET方法")
		return
	}
	r.ParseForm()
	redis := true
	if len(r.Form["redis"]) > 0 {
		flag, err := strconv.ParseBool(r.Form["redis"][0])
		if err != nil {
			util.ResponseErr(w, err)
			return
		}
		redis = flag
	}
	nowDate := util.FormatDate(time.Now(), util.YYYY_MM_DD)
	firstOfYear := util.FirstDayOfCurrentYearAsString()
	result, err := service.FindUserMarkRank(firstOfYear, nowDate, redis)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	fmt.Fprint(w, result)
}
