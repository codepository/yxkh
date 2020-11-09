package conmgr

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// FindLastMonthMarkRankTop10 查询上月一线考核得分排名前10的
func FindLastMonthMarkRankTop10() ([]string, error) {
	// 获取上月日期起止
	start, _ := util.GetLastMonthStartAndEnd()
	enddate := start.Add(-24 * time.Hour)
	startdate := time.Date(enddate.Year(), enddate.Month(), 1, 0, 0, 0, 0, enddate.Location())
	// 查询参数
	sql := "select userId,username from res_mark where startDate >= ? and endDate<=? and checked=1 group by userid,username  order by ifnull(round(sum(markNumber),2),0) desc limit 10"
	marks, err := model.FindAllMark(sql, startdate, enddate)
	var result []string
	if err != nil {
		return result, nil
	}
	for _, v := range marks {
		result = append(result, v.Username)
	}
	return result, nil

}

// FindCurrentYearRankByUserids FindCurrentYearRankByUserids
func FindCurrentYearRankByUserids(userids []interface{}) ([]*model.ResMark, error) {
	if len(userids) == 0 {
		return nil, nil
	}
	// 获取今年开始日期
	start := util.FirstDayOfCurrentYearAsString()
	end := util.FormatDate3(time.Now())
	// 查询参数
	var idbuffer strings.Builder
	for _, id := range userids {
		idbuffer.WriteString(fmt.Sprintf(",%v", id))
	}
	sql := "select userId,username, ifnull(round(sum(markNumber),2),0) markNumber from res_mark where startDate >= ? and endDate<=? and checked=1 and userId in (" + idbuffer.String()[1:] + ") group by userid,username  order by markNumber desc limit 10"

	return model.FindAllMark(sql, start, end)
}

// FindMarkRankCurrentYearByGroup 分组查询用户今年的考核分数
// groups表示考核组名称数组，level表示职级数组，tags表示用户标签数组
// 查询：group+level,group+tag
func FindMarkRankCurrentYearByGroup(groups []string, level [][]int, tags []string) ([]interface{}, error) {
	var datas []interface{}
	var fields []string
	for _, post := range level {
		// 组合 groups 和 level
		for _, group := range groups {
			fields = append(fields, fmt.Sprintf("%s-%v", group, post))
			// 远程查询用户id，要排除已经退休的
			ids, err := FindUseridsByTagsAndLevel([]string{group}, []string{"and", "or"}, post)
			if err != nil {
				return nil, err
			}
			// 根据ids查询结果
			marks, err := FindCurrentYearRankByUserids(ids)
			if err != nil {
				return nil, err
			}
			datas = append(datas, marks)
		}
	}
	for _, tag := range tags {
		for _, group := range groups {
			// groups 和 tags 组合

			fields = append(fields, group+"-"+tag)
			ids, err := FindUseridsByTags([]string{group, tag}, "and")
			if err != nil {
				return nil, err
			}
			// log.Println("查询：", group+"-"+tag)
			// log.Println("结果：", ids)
			// 根据ids查询结果
			marks, err := FindCurrentYearRankByUserids(ids)
			if err != nil {
				return nil, err
			}
			datas = append(datas, marks)
		}

	}
	datas = append(datas, fields)
	datas = append(datas, groups)

	return datas, nil
}

// UpdateMark 修改评分
func UpdateMark(c *model.Container) error {
	errstr := `参数格式：{"body":"params":{"mark":{"markId":2,"markNumber":"3","markReason":"原因"}}} markNumber必须为字符串`
	if len(c.Body.Params) == 0 || c.Body.Params["mark"] == nil {
		return errors.New(errstr)
	}
	mark, yes := c.Body.Params["mark"].(map[string]interface{})

	if !yes {
		return errors.New(errstr)
	}
	if mark["markId"] == nil {
		return errors.New("markId 不能为空")
	}
	id, err := util.Interface2Int(mark["markId"])
	if err != nil {
		return err
	}
	err = model.UpdatesMark(id, c.Body.Params["mark"])
	return err
}

// DelMark 删除
func DelMark(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"ids":[1,3]}} ids是要删除的评分`
	if len(c.Body.Params) == 0 || c.Body.Params["ids"] == nil {
		return errors.New(errstr)
	}
	var ids []int
	for _, id := range c.Body.Params["ids"].([]interface{}) {
		i, err := util.Interface2Int(id)
		if err != nil {
			return err
		}
		ids = append(ids, i)
	}
	err := model.DelMarkByIDs(ids)
	if err != nil {
		return err
	}
	return nil

}

// AddMark 添加分数
func AddMark(c *model.Container) error {
	errstr := `参数格式:{"body":"data":[[{"userId":2,"username":"张三","projectId":1,"markNumber":"2.0","markReason":"加分原因","accordingly":"评分依据","startDate":"2020-07-01","endDate":"2020-07-31"}]]}`
	// 验证参数
	if len(c.Body.Data) == 0 || c.Body.Data[0] == nil {
		return fmt.Errorf(errstr)
	}
	var marks []*model.ResMark
	for _, d := range c.Body.Data[0].([]interface{}) {
		m := d.(map[string]interface{})
		mark := &model.ResMark{}
		err := mark.FromJSON(m)
		if err != nil {
			return err
		}
		marks = append(marks, mark)
	}
	// 存储
	var ids []int
	for _, m := range marks {
		err := m.FirstOrCreate()
		if err != nil {
			return err
		}
		ids = append(ids, m.MarkID)
	}
	c.Body.Data = c.Body.Data[:0]
	c.Body.Data = append(c.Body.Data, ids)
	return nil
}

// SumMarks 合计加减分
func SumMarks(c *model.Container) error {
	errstr := `参数格式:{"body":"params":{"startDate":"2020-01-02","endDate":"2020-01-31","userId":114}} 参数不能全为空`
	if len(c.Body.Params) == 0 {
		return fmt.Errorf(errstr)
	}
	if c.Body.Params["startDate"] == nil {
		return fmt.Errorf("startDate 不能为空，如：2020-01-31")
	}
	start, ok := c.Body.Params["startDate"].(string)
	if !ok {
		return fmt.Errorf("startDate 必须为字符串，如：2020-01-31")
	}
	delete(c.Body.Params, "startDate")
	c.Body.Params["checked"] = 1
	var end string
	if c.Body.Params["endDate"] != nil {
		end, ok = c.Body.Params["endDate"].(string)
		if !ok {
			return fmt.Errorf("endDate 必须为字符串,如：2020-01-31")
		}
		delete(c.Body.Params, "endDate")
	} else {
		end = util.FormatDate3(time.Now())
	}
	if len(end) == 10 {
		end = end + " 23:59:59"
	}
	total, err := model.SumMarks(start, end, c.Body.Params)
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, total)
	return nil

}
