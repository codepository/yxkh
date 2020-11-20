package controller

import (
	"errors"

	"github.com/codepository/yxkh/conmgr"
	"github.com/codepository/yxkh/model"
)

const (
	// SysManagerAuthority SysManagerAuthority
	SysManagerAuthority = "系统管理员"
	// AdvertiseAuthority 广告管理权限
	AdvertiseAuthority = "广告管理"
)

// RouteFunction 根据路径指向方法
type RouteFunction func(*model.Container) error

// RouterMap 路由
var RouterMap map[string]*RouteHandler

// RouteHandler 路由
type RouteHandler struct {
	handler RouteFunction
	route   string
	meta    *RouteMeta
}

// RouteMeta 路由参数
type RouteMeta struct {
	// 可以访问该路径的所有角色的id
	authority []string
}

var routers []*RouteHandler

// SetRouters 设置路由
func SetRouters() {
	routers = []*RouteHandler{
		{route: "visit/yxkh/homedata", handler: conmgr.GetHomeData},
		{route: "visit/yxkh/refreshhomedata", handler: conmgr.RefreshHomeData},
		// ***************** 工作考核 ***************************
		// 查询半年全年考核
		{route: "visit/yxkh/findAllEvalution", handler: conmgr.FindAllEvalution},
		// 半年和年度考核排行
		{route: "visit/yxkh/findAllEvalutionRank", handler: conmgr.FindAllEvalutionRank},
		// 导出半年和年度考核排行
		{route: "export/yxkh/findAllEvalutionRank", handler: conmgr.ExportAllEvalutionRank},
		// 查询半年、年度考核评分占比
		{route: "visit/yxkh/scoreShare", handler: conmgr.ScoreShare},
		// 导入群众评议
		{route: "exec/yxkh/importPublicAssess", handler: conmgr.ImportPublicAssess},
		// 导入加减分
		{route: "exec/yxkh/importMarks", handler: conmgr.ImportMarks},
		// ****************** 流程 *****************************
		// 启动流程
		{route: "exec/yxkh/startProcess", handler: conmgr.StartProcess},
		{route: "exec/yxkh/completeProcess", handler: conmgr.CompleteProcess},
		// 删除流程
		{route: "exec/yxkh/delFlow", handler: conmgr.DelFlow},
		// 查询流程数据
		{route: "visit/yxkh/findFlowDatas", handler: conmgr.FindFlowDatas},
		{route: "visit/yxkh/findAllProcess", handler: conmgr.FindallProcess},
		// ****************** 项目 *****************************
		// 项目和评分查询
		{route: "visit/yxkh/findAllProjectWithMarks", handler: conmgr.FindAllProjectWithMarks},
		// 添加项目
		{route: "exec/yxkh/addProject", handler: conmgr.AddProject},
		// 删除项目
		{route: "exec/yxkh/delProject", handler: conmgr.DelProject},
		// 修改项目
		{route: "exec/yxkh/updateProject", handler: conmgr.UpdateProject},
		// 添加评分
		{route: "exec/yxkh/addMark", handler: conmgr.AddMark},
		// 删除评分
		{route: "exec/yxkh/delMark", handler: conmgr.DelMark},
		// 修改评分
		{route: "exec/yxkh/updateMark", handler: conmgr.UpdateMark},
		// 合计加减分
		{route: "visit/yxkh/sumMarks", handler: conmgr.SumMarks},
		// ************************ 字典 *****************************
		// 字典查询
		{route: "visit/yxkh/findDict", handler: conmgr.FindallDict},
		{route: "exec/yxkh/updateDict", handler: conmgr.UpdateDict},
		{route: "exec/yxkh/addDict", handler: conmgr.AddDict},
		{route: "exec/yxkh/delDict", handler: conmgr.DelDict},
	}
}

// GetRoute 获取执行函数
func GetRoute(route, token string) (func(*model.Container) error, error) {
	var f *RouteHandler
	for _, r := range routers {
		if r.route == route {
			f = r
			break
		}
	}
	if f == nil {
		return nil, errors.New("method:" + route + ",不存在")
	}
	return f.handler, nil
}
