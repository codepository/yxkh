package conmgr

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// AutoDeductWithMonthProcess 月度考核自动减分程序
func AutoDeductWithMonthProcess() {

	errlog := &model.ErrLog{
		BusinessType: "月度考核自动减分程序",
	}
	now := time.Now()
	errlog.CreateTime = now
	// 查询设置,判断是否启用该程序
	dic, err := model.FindSingleDict("name='月度考核延期扣分是否启用' and type='月度考核'")
	if err != nil {
		errlog.Err = err.Error()
		errlog.Create()
		return
	}
	if dic.Value != "1" {
		print("无需扣分，若是需要开启将info_dic表中字段name为[月度考核自动减分程序]的value值设置为1")
		return
	}
	// 判断是否处于需减分的时间段
	start, err := FindStartDayOfAutoDedcutOfMonthProcess()
	if err != nil {
		errlog.Err = err.Error()
		errlog.Create()
		return
	}
	end, err := FindEndDayOfAutoDedcutOfMonthProcess()
	if err != nil {
		errlog.Err = err.Error()
		errlog.Create()
		return
	}
	day := now.Day()
	if day < start || day > end {
		return
	}
	year := now.Year()
	month := now.Month()
	if month == 1 {
		year--
		month = 12
	}
	// 查询需要减分的用户
	users, err := FindPersonApplyYxkh("id,name", fmt.Sprintf("%d年%d月份-月度考核", year, month), 0, 10000)
	if err != nil {
		errlog.Err = fmt.Sprintf("查询未提交月度考核的用户:%s", err.Error())
		errlog.Create()
		return
	}
	if users == nil {
		return
	}
	startDate, endDate := util.GetLastMonthStartAndEnd()
	startStr := util.FormatDate3(startDate)
	endStr := util.FormatDate3(endDate)
	datas := users.([]map[string]interface{})
	// 生成减分的项
	for _, u := range datas {
		err := model.AddProjectWithMark(startStr, endStr, "系统导入", u["id"].(int), "1", "0.5", "月度考核延迟提交产生扣分", u["name"].(string))
		if err != nil {
			errlog.Err = fmt.Sprintf("添加减分:%s", err.Error())
			errlog.Create()
			return
		}
	}
}

// FindLastMonthMarkRankTop10 查询上月一线考核得分排名前10的
func FindLastMonthMarkRankTop10() ([]string, error) {
	// 获取上月日期起止
	start, _ := util.GetLastMonthStartAndEnd()
	enddate := start.Add(-24 * time.Hour)
	startdate := time.Date(enddate.Year(), enddate.Month(), 1, 0, 0, 0, 0, enddate.Location())
	// 查询参数
	sql := fmt.Sprintf("select userId,username from res_mark where startDate >='%s' and endDate<='%s' and checked=1 group by userid,username  order by ifnull(round(sum(markNumber),2),0) desc limit 10", util.FormatDate3(startdate), util.FormatDate3(enddate))
	marks, err := model.FindAllMarkBySQL(sql)
	var result []string
	if err != nil {
		return result, nil
	}
	for _, v := range marks {
		result = append(result, v.Username)
	}
	return result, nil

}

// FindMarkRankCurrentYearByGroup 分组查询用户今年的考核分数
// groups表示考核组名称数组，level表示职级数组，tags表示用户标签数组
// 查询：group+level,group+tag
func FindMarkRankCurrentYearByGroup(groups []string, level [][]int, tags []string) ([]interface{}, error) {
	var datas1 []interface{}
	var datas2 []interface{}
	var fields []string
	now := time.Now()
	start := util.FirstDayOfCurrentYearAsString()
	end := util.FormatDate3(time.Now())
	// 如何现在月份小于3月，那么显示去年的加减分排行，因为今年的还没统计
	if now.Month() < 3 {
		start = fmt.Sprintf("%d-01-01", now.Year()-1)
		end = fmt.Sprintf("%d-12-31", now.Year()-1)
	}
	for _, post := range level {
		// 组合 groups 和 level
		for _, group := range groups {
			// print("考核组:", fmt.Sprintf("%s-%v", group, post))
			fields = append(fields, fmt.Sprintf("%s-%v", group, post))
			// 远程查询用户id，要排除已经退休的
			ids, err := FindUseridsByTagsAndLevel([]string{group}, []string{"and", "or"}, post)
			if err != nil {
				return nil, fmt.Errorf("查询用户id:%s", err.Error())
			}
			// 根据ids查询结果
			marks, err := FindCurrentYearRankByUserids(ids, start, end)
			if err != nil {
				return nil, err
			}
			datas2 = append(datas2, marks)
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
			marks, err := FindCurrentYearRankByUserids(ids, start, end)
			if err != nil {
				return nil, err
			}
			datas2 = append(datas2, marks)
		}

	}
	datas1 = append(datas1, datas2)
	datas1 = append(datas1, fields)
	datas1 = append(datas1, groups)
	datas1 = append(datas1, []string{start, end})
	return datas1, nil
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

// FindMarksRank 查询加减分排行
func FindMarksRank(c *model.Container) error {
	errstr := `参数格式:{"body":{"params":{"limit":50,"offset":0,"startDate":"2020-01-01","endDate":"2020-01-02","username":"用户姓名","group":"第一考核组成员,第二考核组成员","tags":"项目舞台,系统管理员","level":"0,1,2"}}} group 表示考核组 tags 表示用户标签，level表示用户职级：0-普通员工，1-中层副职,2-中层正职`
	var err error
	if len(c.Body.Params) == 0 {
		return fmt.Errorf(errstr)
	}
	// 获取参数
	var sqlbuff strings.Builder
	var userbuff strings.Builder
	// 获取考核组数组
	if c.Body.Params["group"] != nil {
		groupstr := c.Body.Params["group"].(string)
		grouparr := strings.Split(groupstr, ",")
		gsql := ""
		for _, g := range grouparr {
			gsql += fmt.Sprintf(",'%s'", g)
		}
		if len(gsql) > 0 {
			userbuff.WriteString(fmt.Sprintf("and id in (select uId from weixin_oauser_taguser where tagId in (select id from weixin_oauser_tag where tagName in (%s) ) ) ", gsql[1:]))
		}
	}
	// 获取用户标签数组
	if c.Body.Params["tags"] != nil {
		tagstr := c.Body.Params["tags"].(string)
		tagsarr := strings.Split(tagstr, ",")
		gsql := ""
		for _, g := range tagsarr {
			gsql += fmt.Sprintf(",'%s'", g)
		}
		if len(gsql) > 0 {
			userbuff.WriteString(fmt.Sprintf("and id in (select uId from weixin_oauser_taguser where tagId in (select id from weixin_oauser_tag where tagName in (%s) ) ) ", gsql[1:]))
		}
	}
	// 获取用户职级数组
	if c.Body.Params["level"] != nil {
		lstr := c.Body.Params["level"].(string)
		larr := strings.Split(lstr, ",")
		lsql := ""
		for _, l := range larr {
			lsql += fmt.Sprintf(",%s", l)
		}
		if len(lsql) > 0 {
			userbuff.WriteString(fmt.Sprintf("and level in (%s) ", lsql[1:]))
		}
	}
	// 获取用户姓名
	if c.Body.Params["username"] != nil {
		username := c.Body.Params["username"].(string)
		if len(username) > 0 {
			userbuff.WriteString(fmt.Sprintf("and name='%s' ", username))
		}
	}
	// 考核组 and 用户标签 and 用户职级 and 用户姓名 查询用户id
	if userbuff.Len() == 0 {
		return fmt.Errorf("参数 group,tags,level,username 不能同时为空")
	}
	users, err := FindAllUsers("id", map[string]interface{}{"where": userbuff.String()[4:]})
	if err != nil {
		return fmt.Errorf("查询用户:%s", err.Error())
	}
	var idsql strings.Builder
	for _, u := range users {
		umap := u.(map[string]interface{})
		idsql.WriteString(fmt.Sprintf(",%v", umap["id"]))
	}
	if idsql.Len() == 0 {
		return fmt.Errorf("未查询到用户")
	}
	sqlbuff.WriteString(fmt.Sprintf("and userId in (%s) ", idsql.String()[1:]))
	// 查询加减分排行
	if c.Body.Params["startDate"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and startDate>='%s' ", c.Body.Params["startDate"].(string)))
	}
	if c.Body.Params["endDate"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and endDate<='%s' ", c.Body.Params["endDate"].(string)))
	}
	sqlbuff.WriteString("and checked=1 ")
	limit := 50
	if c.Body.Params["limit"] != nil {
		limit, err = util.Interface2Int(c.Body.Params["limit"])
		if err != nil {
			return fmt.Errorf("参数 limit:%s", err.Error())
		}
	}
	offset := 0
	if c.Body.Params["offset"] != nil {
		offset, err = util.Interface2Int(c.Body.Params["offset"])
		if err != nil {
			return fmt.Errorf("参数 offset:%s", err.Error())
		}
	}
	datas, err := model.FindMarksRankPaged(limit, offset, sqlbuff.String()[4:])
	if err != nil {
		return fmt.Errorf("查询加减分排行:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, datas)
	return nil

}

// FindCurrentYearRankByUserids FindCurrentYearRankByUserids
func FindCurrentYearRankByUserids(userids []interface{}, start, end string) ([]*model.ResMark, error) {
	if len(userids) == 0 {
		return nil, nil
	}
	// 查询参数
	var idbuffer strings.Builder
	for _, id := range userids {
		idbuffer.WriteString(fmt.Sprintf(",%v", id))
	}
	sql := fmt.Sprintf("select userId,username, ifnull(round(sum(markNumber),2),0) as markNumber from res_mark where startDate >='%s' and endDate<='%s' and checked=1 and userId in ("+idbuffer.String()[1:]+") group by userid,username  order by markNumber desc limit 10", start, end)
	return model.FindAllMarkBySQL(sql)
}

// FindAllMarks 根据条件查询所有加减分
func FindAllMarks(c *model.Container) error {
	errstr := `参数格式:{"body":{"params":{"offset":0,"limit":20,"fields":"markId,userId,username,markReason,markNumber,accordingly","markId":1,"projectId":3,"markReason":"加分原因","accordingly":"加分依据的规则","startDate":"2020-04-03","endDate":"2020-05-06","userId":19,"username":"用户名","checked":"1"}}} ，checked为字符串：1-已生效、0-未生效`
	var err error
	if len(c.Body.Params) == 0 {
		return fmt.Errorf(errstr)
	}
	var sqlbuff strings.Builder
	if c.Body.Params["checked"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and checked='%s' ", c.Body.Params["checked"].(string)))
	}
	if c.Body.Params["username"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and username='%s' ", c.Body.Params["username"].(string)))
	}
	if c.Body.Params["userId"] != nil {
		userid, err := util.Interface2Int(c.Body.Params["userId"])
		if err != nil {
			return fmt.Errorf("参数userId:%s", err.Error())
		}
		sqlbuff.WriteString(fmt.Sprintf("and userId=%d ", userid))
	}
	if c.Body.Params["endDate"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and endDate<='%s' ", c.Body.Params["endDate"].(string)))
	}
	if c.Body.Params["startDate"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and startDate>='%s' ", c.Body.Params["startDate"].(string)))
	}
	if c.Body.Params["accordingly"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and accordingly like '%s' ", "%"+c.Body.Params["accordingly"].(string)+"%"))
	}
	if c.Body.Params["markReason"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and markReason like '%s' ", "%"+c.Body.Params["markReason"].(string)+"%"))
	}
	if c.Body.Params["projectId"] != nil {
		projectID, err := util.Interface2Int(c.Body.Params["projectId"])
		if err != nil {
			return fmt.Errorf("参数projectId:%s", err.Error())
		}
		sqlbuff.WriteString(fmt.Sprintf("and projectId=%d ", projectID))
	}
	if c.Body.Params["markId"] != nil {
		markID, err := util.Interface2Int(c.Body.Params["markId"])
		if err != nil {
			return fmt.Errorf("参数markId:%s", err.Error())
		}
		sqlbuff.WriteString(fmt.Sprintf("and markId=%d ", markID))
	}
	if sqlbuff.Len() == 0 {
		return fmt.Errorf(errstr)
	}
	fields := ""
	if c.Body.Params["fields"] != nil {
		fields = c.Body.Params["fields"].(string)
	}
	limit := 20
	if c.Body.Params["limit"] != nil {
		limit, err = util.Interface2Int(c.Body.Params["limit"])
		if err != nil {
			return fmt.Errorf("参数 limit:%s", err.Error())
		}
	}
	offset := 0
	if c.Body.Params["offset"] != nil {
		offset, err = util.Interface2Int(c.Body.Params["offset"])
		if err != nil {
			return fmt.Errorf("参数 offset:%s", err.Error())
		}
	}
	datas, err := model.FindAllMarkPaged(fields, limit, offset, "", sqlbuff.String()[4:])
	if err != nil {
		return fmt.Errorf("查询加减分:%s", err)
	}
	c.Body.Data = append(c.Body.Data, datas)
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
