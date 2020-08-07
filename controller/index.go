package controller

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// Index 首页
func Index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(writer, "Hello world!")
}

// GetToken 获取token
func GetToken(request *http.Request) (string, error) {
	token := request.Header.Get("Authorization")
	if len(token) == 0 {
		request.ParseForm()
		if len(request.Form["token"]) == 0 {
			return "", errors.New("header Authorization 没有保存 token, url参数也不存在 token， 访问失败 ！")
		}
		token = request.Form["token"][0]
	}
	return token, nil
}

// GetData 查询接口
func GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持POST方法")
		return
	}
	var par model.Container
	err := util.Body2Struct(r, &par)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	if len(par.Body.Method) == 0 {
		util.ResponseErr(w, "method不能为空")
		return
	}
	token := par.Header.Token
	if len(token) == 0 {
		token, _ = GetToken(r)
	}
	f, err := GetRoute(par.Body.Method, par.Header.Token)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	err = f(&par)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	util.ResponseData(w, par.ToString())
}

// Export 导出数据
func Export(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.ResponseErr(w, "只支持POST方法")
		return
	}
	var par model.Container
	err := util.Body2Struct(r, &par)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	errstr := `参数格式:{"body":{"method":"export/yxkh/findAllEvalutionRank","data":[["用户名","部门"],["username","deptName"]]}} data[0]导出文件行标、data[1]表示行标对应的字段，俩者一一对应，都不能为空`
	if len(par.Body.Method) == 0 {
		util.ResponseErr(w, errstr)
		return
	}

	if par.Body.Data == nil || len(par.Body.Data) != 2 {
		util.ResponseErr(w, errstr)
		return
	}
	token := par.Header.Token
	if len(token) == 0 {
		token, _ = GetToken(r)
	}
	categoryHeader := par.Body.Data[0].([]interface{})
	var header []string
	for _, h := range categoryHeader {
		header = append(header, h.(string))
	}
	// if !ok {
	// 	util.ResponseErr(w, errstr)
	// 	return
	// }
	f, err := GetRoute(par.Body.Method, par.Header.Token)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	err = f(&par)
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	datas := par.Body.Data

	// 导出
	fileName := "export.csv"
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	wr.Write(header)
	for i := 0; i < len(datas); i++ {
		wr.Write(datas[i].([]string))
	}
	wr.Flush()
	w.Header().Set("Content-Type", "application/vnd.ms-excel;charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	w.Write(b.Bytes())
}
