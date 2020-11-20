package router

import (
	"fmt"
	"net/http"

	"github.com/codepository/yxkh/config"
	"github.com/codepository/yxkh/controller"
)

// Mux 路由
var Mux = http.NewServeMux()
var conf = *config.Config

func init() {
	setMux()
}
func intercept(h http.HandlerFunc) http.HandlerFunc {
	return crossOrigin(h)
}
func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", conf.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", conf.AccessControlAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AccessControlAllowHeaders)
		r.Header.Add("Accept-Charset", "utf-8")
		h(w, r)
	}
}
func setMux() {
	prefixV1 := "/api/v1/yxkh"
	Mux.HandleFunc("/api/v1/yxkh/index", controller.Index)
	Mux.HandleFunc("/api/v1/yxkh/getData", intercept(controller.GetData))
	// ------------------ 导出数据 ------------------------
	Mux.HandleFunc("/api/v1/yxkh/export", intercept(controller.Export))
	// 导入
	Mux.HandleFunc("/api/v1/yxkh/import", intercept(controller.Import))
	Mux.HandleFunc("/api/v1/yxkh/import/publicAssessment", intercept(controller.ImportPublicAssessment))
	// -------------------------- project ----------------
	Mux.HandleFunc("/api/v1/yxkh/project/save", intercept(controller.SaveProject))
	Mux.HandleFunc("/api/v1/yxkh/project/findall", intercept(controller.FindAllProject))
	Mux.HandleFunc("/api/v1/yxkh/project/del", intercept(controller.DelProjectByID))
	Mux.HandleFunc("/api/v1/yxkh/project/update", intercept(controller.UpdateProject))
	// -------------------------- mark ---------------------
	Mux.HandleFunc("/api/v1/yxkh/mark/save", intercept(controller.SaveMark))
	Mux.HandleFunc("/api/v1/yxkh/mark/findall", intercept(controller.FindAllMark))
	Mux.HandleFunc("/api/v1/yxkh/mark/del", intercept(controller.DelMark))
	Mux.HandleFunc("/api/v1/yxkh/mark/update", intercept(controller.UpdateMark))
	Mux.HandleFunc(fmt.Sprintf("%s/mark/findMarkRankForHome", prefixV1), intercept(controller.FindMarkRankForHome))

}
