package conmgr

import (
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
