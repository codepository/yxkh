package conmgr

import (
	"fmt"
	"log"
	"time"

	"github.com/codepository/yxkh/model"
)

// RefreshHomeData 刷新首页数据
func RefreshHomeData(c *model.Container) error {
	// 返回字段设置
	c.Body.Fields = []string{"上月月度考核排名前10", "月度考核提交情况", "任务审批数排行", "月度考核未交清单", "加减分排行", "年度考评排行", "半年考核排行", "一线考核嘉奖通报", "一线考核相关文件", "考核组"}
	// 查询考核组
	var assessGroup []string
	datas, err := FindAllTags(map[string]interface{}{"field": "tagName", "type": "考核组"})
	if err != nil {
		return fmt.Errorf("查询标签:%s", err.Error())
	}
	if len(datas) == 0 {
		return fmt.Errorf("查询标签:不存在type为[考核组]的标签")
	}
	for _, g := range datas {
		x := g.(map[string]interface{})
		assessGroup = append(assessGroup, x["tagName"].(string))
	}
	// 上月月度考核排名前10
	usernames, err := FindLastMonthMarkRankTop10()
	if err != nil {
		return fmt.Errorf("上月月度考核排名前10:%s", err.Error())
	}
	if len(usernames) == 0 {
		c.Body.Data = append(c.Body.Data, []interface{}{})
	} else {
		c.Body.Data = append(c.Body.Data, usernames)
	}
	now := time.Now()
	year := now.Year()
	month := now.Month() - 1
	if month == 0 {
		month = 12
		year--
	}
	titleLike := fmt.Sprintf("%d年%d月份-月度考核", year, month)
	println("查询:", titleLike)
	// 上个月一线考核提交情况
	completedDescribe, err := FindTaskCompletedDescribe(titleLike)
	if err != nil {
		return fmt.Errorf("提交率:%s", err.Error())
	}
	if completedDescribe == nil {
		c.Body.Data = append(c.Body.Data, []interface{}{})

	} else {
		c.Body.Data = append(c.Body.Data, completedDescribe)
	}
	// 上月月度考核未审批任务排行和已审批任务排行
	taskRank, err := FindTaskRank(titleLike)
	if err != nil {
		return fmt.Errorf("月度考核审批排行:%s", err.Error())
	}
	if taskRank == nil {
		c.Body.Data = append(c.Body.Data, []interface{}{})
	} else {
		c.Body.Data = append(c.Body.Data, taskRank)
	}
	// 未提交一线考核的员工
	users, err := FindPersonApplyYxkh("", titleLike, 0, 5)
	if err != nil {
		return fmt.Errorf("月度考核未交清单:%s", err.Error())
	}
	if users == nil {
		c.Body.Data = append(c.Body.Data, []interface{}{})
	} else {
		c.Body.Data = append(c.Body.Data, users)
	}
	// 加减分排行
	userMarks, err := FindMarkRankCurrentYearByGroup(assessGroup, [][]int{{1, 2}, {0}}, []string{"项目舞台"})
	if err != nil {
		return fmt.Errorf("加减分排行:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, userMarks)

	month = now.Month()
	year = now.Year()
	sparation1 := ""
	sparation2 := ""
	if month > 6 {
		sparation1 = fmt.Sprintf("%d年-半年考核", year)
	} else {
		sparation1 = fmt.Sprintf("%d年-半年考核", year-1)
	}
	sparation2 = fmt.Sprintf("%d年-年度考核", year-1)
	// 半年考核排行
	halfyear, err := model.FindAllEvaluation("department,sparation,result,publicEvaluation,leadershipEvaluation,eId,uid,username,totalMark,marks,overseerEvaluation", "res_evaluation.totalMark+0 desc,res_evaluation.marks+0 desc", 20, 0, map[string]interface{}{"sparation": sparation1})
	if err != nil {
		return fmt.Errorf("查询半年考核:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, halfyear)
	// 年度考评排行
	fullyear, err := model.FindAllEvaluation("department,sparation,result,publicEvaluation,leadershipEvaluation,eId,uid,username,totalMark,marks,overseerEvaluation", "res_evaluation.totalMark+0 desc,res_evaluation.marks+0 desc", 20, 0, map[string]interface{}{"sparation": sparation2})
	if err != nil {
		return fmt.Errorf("查询年度考核:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, fullyear)
	// 一线考核嘉奖通报
	uploadfiles, err := FindAllUploadfiles("remarks", "id,filename,uid,username")
	if err != nil {
		return fmt.Errorf("查询嘉奖通报:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, uploadfiles)
	// 一线考核相关文件
	public, err := FindAllUploadfiles("public", "id,filename,uid,username")
	if err != nil {
		return fmt.Errorf("查询公共文件:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, public)
	// 一线考核组
	c.Body.Data = append(c.Body.Data, assessGroup)
	Conmgr.cacheMap[HomeDataCache] = c.Body.Data
	return nil

}

// GetHomeData 获取所有首页数据
func GetHomeData(c *model.Container) error {
	if Conmgr.cacheMap[HomeDataCache] == nil {
		err := RefreshHomeData(c)
		if err != nil {
			return err
		}
		Conmgr.cacheMap[HomeDataCache] = c.Body.Data
	} else {
		log.Println("get homedata from cache")
		c.Body.Data = Conmgr.cacheMap[HomeDataCache].([]interface{})
	}
	return nil
}
