package controller

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

var (
	// 文件 key
	uploadFileKey = "upload-key"
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
	if !par.Body.Paged {
		util.ResponseData(w, par.ToString())
		return
	}
	result, err := par.ToPageJSON()
	if err != nil {
		util.ResponseErr(w, err)
		return
	}
	fmt.Fprintf(w, result)

}

// Import 导入
func Import(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		util.ResponseErr(w, "只支持POST方法")
		return
	}
	//解析 form 中的file 上传名字
	file, filehead, err := req.FormFile("filename")
	if err != nil {
		fmt.Fprintf(w, "file upload fail:%s", err)
	}
	if len(filehead.Filename) == 0 {
		fmt.Fprintf(w, "filename不能为空")
		return
	}
	ss := strings.Split(filehead.Filename, ".")
	if ss[len(ss)-1] != "xlsx" {
		fmt.Fprintf(w, "只支持xlsx")
		file.Close()
		return
	}
	filesave := fmt.Sprintf("%s%d", filehead.Filename, time.Now().Nanosecond())
	//打开 已只读,文件不存在创建 方式打开  要存放的路径资源
	f, err := os.OpenFile(filesave, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		f.Close()
		file.Close()
		os.Remove(filesave)
		fmt.Fprintf(w, "file open fail:%s", err)
		return
	}
	//文件 copy
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Fprintf(w, "file copy fail:%s", err)
		f.Close()
		file.Close()
		os.Remove(filesave)
		return
	}
	//  文件导入之后，执行操作
	par := model.Container{}
	values := req.URL.Query()
	method := values.Get("method")
	token := values.Get("token")
	if len(token) == 0 {
		token = req.PostFormValue("token")
		if len(token) == 0 {
			fmt.Fprintf(w, "需要添加token,可以以get或者post方式添加")
			return
		}
	}
	if len(method) == 0 {
		method = req.PostFormValue("method")
		if len(method) == 0 {
			fmt.Fprintf(w, "需要添加method参数,可以以get或者post方式添加")
			return
		}
	}
	par.Body.Method = method
	par.Header.Token = token
	par.File = f
	funcs, err := GetRoute(par.Body.Method, par.Header.Token)
	if err != nil {
		util.ResponseErr(w, err)
		f.Close()
		file.Close()
		os.Remove(filesave)
		return
	}
	err = funcs(&par)

	// err = conmgr.GetPublicAssessFromXlsx(f)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		f.Close()
		file.Close()
		os.Remove(filesave)
		return
	}
	//关闭对应打开的文件
	f.Close()
	file.Close()
	os.Remove(filesave)
	fmt.Fprintf(w, "成功")

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
	w.Header().Set("Content-Type", "text/csv;charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	w.Write(b.Bytes())
}
