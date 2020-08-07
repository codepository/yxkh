package conmgr

import (
	"fmt"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// RefreshHomeData 刷新首页数据
func RefreshHomeData(c *model.Container) error {
	// 返回字段设置
	c.Body.Fields = []string{"上月月度考核排名前10", "提交率", "月度考核未交清单", "加减分排行", "年度考评排行", "半年考核排行", "一线考核嘉奖通报", "一线考核相关文件", "考核组"}
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
	// 提交率
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Local().Location())
	start := util.FormatDate3(startDate)
	completeRates, err := FindTaskCompleteRates("一线考核", start)
	if err != nil {
		return fmt.Errorf("提交率:%s", err.Error())
	}
	if completeRates == nil {
		c.Body.Data = append(c.Body.Data, []interface{}{})

	} else {
		c.Body.Data = append(c.Body.Data, completeRates)
	}

	// 月度考核未交清单,只显示30个
	users, err := FindUsersUncompleteTask("一线考核", start)
	if err != nil {
		return fmt.Errorf("月度考核未交清单:%s", err.Error())
	}
	if users == nil {
		c.Body.Data = append(c.Body.Data, []interface{}{})
	} else {
		c.Body.Data = append(c.Body.Data, users)
	}
	// 加减分排行
	assessGroup := []string{"第一考核组成员", "第二考核组成员", "第三考核组成员", "第四考核组成员"}
	userMarks, err := FindMarkRankCurrentYearByGroup(assessGroup, [][]int{{1, 2}, {0}}, []string{"项目舞台"})
	if err != nil {
		return fmt.Errorf("加减分排行:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, userMarks)
	// 半年考核排行
	c.Body.Data = append(c.Body.Data, []interface{}{})
	// 年度考评排行
	c.Body.Data = append(c.Body.Data, []interface{}{})
	// 一线考核嘉奖通报
	uploadfiles, err := FindAllUploadfiles(map[string]interface{}{"filetype": "remarks"}, 10, 0)
	if len(uploadfiles) == 0 {
		c.Body.Data = append(c.Body.Data, []interface{}{})
	} else {
		c.Body.Data = append(c.Body.Data, uploadfiles)
	}
	// 一线考核相关文件
	public, err := FindAllUploadfiles(map[string]interface{}{"filetype": "public"}, 10, 0)
	if len(uploadfiles) == 0 {
		c.Body.Data = append(c.Body.Data, []interface{}{})
	} else {
		c.Body.Data = append(c.Body.Data, public)
	}
	// 一线考核组
	c.Body.Data = append(c.Body.Data, assessGroup)
	return nil
}

// GetHomeData 获取所有首页数据
func GetHomeData(c *model.Container) error {
	if Conmgr.cacheMap[HomeDataCache] == nil {
		err := RefreshHomeData(c)
		if err != nil {
			return err
		}
	} else {
		c.Body.Data = Conmgr.cacheMap[HomeDataCache].([]interface{})
	}
	return nil
}
