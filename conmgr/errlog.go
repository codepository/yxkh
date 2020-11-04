package conmgr

import (
	"fmt"
	"log"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// ReExecData 重新执行数据格式
type ReExecData struct {
	// Key 关键字
	Key string
	// FuncName 执行函数名
	FuncName string
}

// ReExecFunction 重新执行失败任务函数
type ReExecFunction func(string) (*ReExecData, error)

// ErrRouteHandler 失败路由对象
type ErrRouteHandler struct {
	route   string
	handler ReExecFunction
}

var reExecRouter []*ErrRouteHandler = []*ErrRouteHandler{
	// 月度考核自动评分
	{route: "RemarkEvaluationByProcessInstanceID", handler: RemarkEvaluationByProcessInstanceID},
	// 流程结束后触发事件
	{route: "ExecWhenProcessCompleted", handler: ExecWhenProcessCompleted},
	// 月度考核相关评分使生效
	{route: "CheckProjectByProcessInstanceID", handler: CheckProjectByProcessInstanceID},
}

// ToString ToString
func (r *ReExecData) ToString() string {
	str, _ := util.ToJSONStr(r)
	return str
}

// getRoute 获取路由
func getRoute(router string) (*ErrRouteHandler, error) {
	var f *ErrRouteHandler
	for _, r := range reExecRouter {
		if r.route == router {
			f = r
			break
		}
	}
	if f == nil {
		return nil, fmt.Errorf("reExecRouter method:" + router + ",不存在")
	}
	return f, nil
}

// ReExeProcessByErrLog 根据失败日志，重新执行事务
func ReExeProcessByErrLog() error {
	log.Println("根据失败日志，重新执行事务")
	logs, err := model.FindErrLog(map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
	}
	for _, log := range logs {
		red := &ReExecData{}
		err := util.Str2Struct(log.Data, red)
		if err != nil {
			return fmt.Errorf("字符串转ReExecData失败:%s", err.Error())
		}
		fun, err := getRoute(red.FuncName)
		if err != nil {
			return err
		}
		_, err = fun.handler(red.Key)
		if err != nil {
			return err
		}
		err = log.DelByID()
		if err != nil {
			return err
		}
	}
	return nil
}
