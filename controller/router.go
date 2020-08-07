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
		// 查询半年全年考核
		{route: "visit/yxkh/findAllEvalution", handler: conmgr.FindAllEvalution},
		// 半年和年度考核排行
		{route: "visit/yxkh/findAllEvalutionRank", handler: conmgr.FindAllEvalutionRank},
		// 导出半年和年度考核排行
		{route: "export/yxkh/findAllEvalutionRank", handler: conmgr.ExportAllEvalutionRank},
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
