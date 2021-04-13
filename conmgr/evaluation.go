package conmgr

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
	"github.com/tealeg/xlsx"
)

// FindAllEvalutionRank 半年和年度考核排行
func FindAllEvalutionRank(c *model.Container) error {
	errStr := fmt.Sprintf(`参数格式:{"body":{"params":{"fields":"","tags":"第一考核组成员,第二考核组成员","departments":"A部门,B部门","username":"张三","limit":20,"offset":0,"sparation":"2020年-半年考核"}}} tags 为用户标签,username按用户名模糊查询,fields为查询字段,sparation为模糊查询`)
	var err error
	if len(c.Body.Params) == 0 {
		return fmt.Errorf(errStr)
	}
	order := "res_evaluation.totalMark+0 desc,res_evaluation.marks+0 desc"
	limit := 20
	if c.Body.Params["limit"] != nil {
		limit, err = util.Interface2Int(c.Body.Params["limit"])
		if err != nil {
			return fmt.Errorf("limit 转数字:%s", err.Error())
		}
	}
	offset := 0
	if c.Body.Params["offset"] != nil {
		offset, err = util.Interface2Int(c.Body.Params["offset"])
		if err != nil {
			return fmt.Errorf("offset 转数字:%s", err.Error())
		}
	}
	fields := ""
	if c.Body.Params["fields"] != nil {
		fields = c.Body.Params["fields"].(string)
	}
	if c.Body.Params["order"] != nil {
		order = c.Body.Params["order"].(string)
	}
	var usersqlbuff strings.Builder
	var sqlbuff strings.Builder
	if c.Body.Params["tags"] != nil {
		tstr, ok := c.Body.Params["tags"].(string)
		if !ok {
			return fmt.Errorf("tags必须为字符串如'第一考核组成员,第二考核组成员,管理员'")
		}
		tarr := strings.Split(tstr, ",")
		tres := ""
		for _, t := range tarr {
			tres += fmt.Sprintf(",'%s'", t)
		}
		usersqlbuff.WriteString(fmt.Sprintf("and id in (select uId from weixin_oauser_taguser where tagId in (select id from weixin_oauser_tag where tagName in (%s))) ", tres[1:]))
	}
	if c.Body.Params["departments"] != nil {
		dstr, ok := c.Body.Params["departments"].(string)
		if !ok {
			return fmt.Errorf("departments 必须为字符串如'A部门,B部门'")
		}
		darr := strings.Split(dstr, ",")
		dres := ""
		for _, d := range darr {
			dres += fmt.Sprintf(",'%s'", d)
		}
		if len(dres) > 0 {
			usersqlbuff.WriteString(fmt.Sprintf("and departmentname in (%s) ", dres[1:]))
		}
	}
	if c.Body.Params["username"] != nil {
		name, ok := c.Body.Params["username"].(string)
		if !ok {
			return fmt.Errorf("username 必须为字符串")
		}
		usersqlbuff.WriteString(fmt.Sprintf("and name like '%s'", "%"+name+"%"))
	}
	// 查询用户id
	if usersqlbuff.Len() != 0 {
		users, err := FindAllUsers("id", map[string]interface{}{"where": usersqlbuff.String()[4:]})
		if err != nil {
			return fmt.Errorf("从fznewsuser微服务查询用户:%s", err.Error())
		}
		ids := ""
		for _, u := range users {
			umap := u.(map[string]interface{})
			id, err := util.Interface2Int(umap["id"])
			if err != nil {
				return err
			}
			ids += fmt.Sprintf(",%d", id)
		}
		if len(ids) > 0 {
			sqlbuff.WriteString(fmt.Sprintf("and uid in (%s) ", ids[1:]))
		}
	}
	where := ""
	if c.Body.Params["sparation"] != nil {
		sqlbuff.WriteString(fmt.Sprintf("and sparation like '%s' ", "%"+c.Body.Params["sparation"].(string)+"%"))
	}
	if sqlbuff.Len() != 0 {
		where = sqlbuff.String()[4:]
	}

	// 根据用户id查询考核表
	datas, err := model.FindAllEvaluation(fields, order, limit, offset, where)
	if err != nil {
		return fmt.Errorf("查询考核表:%s", err.Error())
	}
	c.Body.Data = append(c.Body.Data, datas)
	return nil
}

// RemarkEvaluationByProcessInstanceID 根据流程检查月度考核评价，并根据评价添加加减分
func RemarkEvaluationByProcessInstanceID(processInstanceID string) (*ReExecData, error) {
	red := &ReExecData{Key: processInstanceID, FuncName: "RemarkEvaluationByProcessInstanceID"}
	// 查询对应的月度考核
	e, err := model.FindSingleEvaluation(processInstanceID)
	users, err := FindAllUsers("name", map[string]interface{}{"id": e.UID})
	if err != nil {
		return red, fmt.Errorf("查询用户名:%s", err.Error())
	}
	user := users[0].(map[string]interface{})
	if err != nil {
		return red, fmt.Errorf("查询月度考核失败:%s", err.Error())
	}
	// 查询 organizationEvaluation 字段所对应的加减分
	dicName := ""
	switch e.OverseerEvaluation {
	case "优秀":
		dicName = "月考评优"
		break
	case "基本合格":
		dicName = "月考基本合格"
		break
	case "不合格":
		dicName = "月考不合格"
		break
	default:
		return red, fmt.Errorf("月度考核评价【%s】不存在，请务必联系管理员", e.OverseerEvaluation)
	}
	dic, err := model.FindAllInfoDic(map[string]interface{}{"type": "月考自动加减分", "name": dicName})
	if err != nil || len(dic) == 0 {
		return red, fmt.Errorf("查询字典【%s】失败:%s", dicName, err.Error())
	}
	// 添加加减分
	_, err = model.AddProjectWithMark(e.StartDate, e.EndDate, dicName, e.UID, "1", dic[0].Value, dicName, user["name"].(string))
	if err != nil {
		return red, fmt.Errorf("月度考核自动加分失败:%s", err.Error())
	}
	return nil, nil
}

// 确定是否在规定时间内导入文件
func checkImportPublicAssessDate(processName string) error {
	// 判断格式是否正确
	var err error
	if len(processName) > 20 {
		return fmt.Errorf("流程类型格式必须是: [2018年-半年考核]或[2018年-年度考核]")
	}
	yes, err := regexp.MatchString(`[0-9]{4}((年-半年考核)|(年-年度考核))`, processName)
	if !yes {
		return fmt.Errorf("流程类型格式必须是: [2018年-半年考核]或[2018年-年度考核]")
	}
	year := strings.TrimSpace(processName)[0:4]
	var start, end time.Time
	typename := strings.Split(processName, "-")[1]
	now := time.Now().Unix()
	switch typename {
	case YXKHBnkh:
		start, err = util.ParseDate3(year + "-07-01")
		if err != nil {
			return err
		}
		end, _ = util.ParseDate3(year + "-12-31")
		// 判断导入时间是否是在 7月1日和12月31号之间
		if start.Unix() <= now && now <= end.Unix() {
			return nil
		}
		return fmt.Errorf("[%s]群众评议只能在[%s]和[%s]之间导入", processName, year+"-07-01", year+"-12-31")
	case YXKHNdkh:
		y, err := strconv.Atoi(year)
		if err != nil {
			return err
		}
		y++
		start, err = util.ParseDate3(fmt.Sprintf("%d-01-01", y))
		if err != nil {
			return err
		}
		end, _ := util.ParseDate3(fmt.Sprintf("%d-06-30", y))
		if start.Unix() <= now && now <= end.Unix() {
			return nil
		}
		return fmt.Errorf("[%s]群众评议只能在[%s]和[%s]之间导入", processName, fmt.Sprintf("%d-01-01", y), fmt.Sprintf("%d-06-30", y))
	default:
		return fmt.Errorf("流程类型格式:【2099年-半年考核】或【2011年-年度考核】")
	}
}

// GetMarksFromXlsx 从xlsx文件读取加减分
func GetMarksFromXlsx(file *os.File, checked string) error {
	if file == nil {
		return fmt.Errorf("file值为nil")
	}
	var err error
	haserr := false
	var buff strings.Builder
	buff.WriteString("可重复导入,请保留原始数据\n")
	// 用户id
	idmap := make(map[string]int)
	// 用户名
	usernameMap := make(map[string]string)
	// 读取数据
	xlFile, err := xlsx.OpenFile(file.Name())
	if err != nil {
		return err
	}
	datas, err := xlFile.ToSlice()
	if err != nil {
		return fmt.Errorf("文件数据转换成数组失败:%s", err.Error())
	}
	// 一共6列
	if len(datas[0][0]) < 6 {
		return fmt.Errorf("一共为6列:1.电话或姓名 2.开始日期 3.结束日期 4.项目内容 5.加分原因 6.加分分数")
	}
	// 参数检查:1.电话或姓名有效性,2.日期有效性,3.加分分数检查
	for i, r := range datas[0] {
		if i == 0 {
			continue
		}
		// 检查日期
		yes, err := util.IsDate3(r[1])
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第2列:%s\n", i+1, err.Error()))
		}
		if err == nil && !yes {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第2列日期格式必须是:yyyy-mm-dd\n", i+1))
		}
		yes, err = util.IsDate3(r[2])
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第2列:%s\n", i+1, err.Error()))
		}
		if err == nil && !yes {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第2列日期格式必须是:yyyy-mm-dd\n", i+1))
		}
		start, _ := util.ParseDate3(r[1])
		end, _ := util.ParseDate3(r[2])
		if start.Unix() > end.Unix() {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第3列日期必须大于第2列日期\n", i+1))
		}
		// 检查第6列分数
		_, err = strconv.ParseFloat(r[5], 64)
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第6列必须为数字:%s\n", i+1, err.Error()))
		}
		// 第1列必须是电话号码或者用户名，当存在重名时报错
		id, name, err := GetUseridAndNameByMobileOrName(r[0])
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第1列用户有误:%s\n", i+1, err.Error()))
		}
		idmap[r[0]] = id
		usernameMap[r[0]] = name
	}
	if haserr {
		return fmt.Errorf(buff.String())
	}
	// 导入数据
	for i, r := range datas[0] {
		if i == 0 {
			continue
		}
		_, err = model.AddProjectWithMark(r[1], r[2], r[3], idmap[r[0]], checked, r[5], r[4], usernameMap[r[0]])
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行导入数据有误:%s\n", i+1, err.Error()))
		}
	}
	if haserr {
		return fmt.Errorf(buff.String())
	}
	return nil

}

// GetPublicAssessFromXlsx 从xlsx文件读取群众评议
func GetPublicAssessFromXlsx(file *os.File) error {
	if file == nil {
		return fmt.Errorf("file值为nil")
	}
	var err error
	haserr := false
	// 先查询基础分
	baseStr, err := model.FindMarksBase()
	if err != nil {
		return err
	}
	basemarks, err := strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return err
	}
	// 优秀、合格 得分
	assessmentMarks, err := model.FindAllInfoDic(map[string]interface{}{"type": "基本定格对应得分"})
	if err != nil {
		return fmt.Errorf("基本定格对应得分时报错:%s", err.Error())
	}
	xlFile, err := xlsx.OpenFile(file.Name())
	if err != nil {
		return err
	}
	datas, err := xlFile.ToSlice()
	if err != nil {
		return fmt.Errorf("文件数据转换成数组失败:%s", err.Error())
	}
	// 验证有效性
	var buff strings.Builder
	buff.WriteString("可重复导入,请保留原始数据\n")
	buff.WriteString("第4列为部门,为社直下属部门的填'社直',其它可不填\n")
	// 第14列必须是最终得分
	if datas[0][0][13] != "最终得分" {
		return fmt.Errorf("第14列必须是【最终得分】")
	}
	// 用户id
	idmap := make(map[string]int)
	departmentidMap := make(map[string]int)
	for i, r := range datas[0] {
		if i == 0 {
			continue
		}
		// 限制导入:只能导入最近一次的群众评议
		err = checkImportPublicAssessDate(r[1])
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第2列:%s\n", i+1, err.Error()))
		}
		// 第1列必须是电话号码或者用户名，当存在重名时报错

		id, did, err := GetUseridAndDeptByMobileOrName(r[0])
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第1列用户有误:%s\n", i+1, err.Error()))
		}
		idmap[r[0]] = id
		departmentidMap[r[0]] = did
		// 第2列必须是流程类型
		yes, err := regexp.MatchString(`[0-9]{4}((年-半年考核)|(年-年度考核))`, r[1])
		if !yes {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第2列数据格式有误必须是：【2099年-半年考核】或【2011年-年度考核】\n", i+1))
		}
		// 第14列必须是数字
		_, err = strconv.ParseFloat(r[13], 32)
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行第14列最终得分有误：%s\n", i+1, err.Error()))
		}
	}
	if haserr {
		return fmt.Errorf(buff.String())
	}
	for i, r := range datas[0] {
		if i == 0 {
			continue
		}
		// 查询流程
		e, err := model.FindSingleEvaluation2(map[string]interface{}{"sparation": r[1], "uid": idmap[r[0]]})
		if err != nil {
			buff.WriteString(fmt.Sprintf("第%d行查询用户[%s]流程数据报错：%s\n", i+1, r[0], err.Error()))
			continue
		}
		if e == nil {
			buff.WriteString(fmt.Sprintf("第%d行用户【%s】还没填写【%s】,待其提交后务必重新导入此表\n", i+1, r[0], r[1]))
			continue
		}
		// 群众评价赋值
		if strings.Contains(r[3], "社直") {
			e.PublicRemark = fmt.Sprintf("%s，分别召集社直员工、社直服务对象开展社直工作人员群众测评。应参加人数(%s)人，实参加人数(%s)人。（1）社直评议结果：发出推荐票(%s)张，收回(%s)张，有效票(%s)张。该同志得票为：优秀(%s)票，合格(%s)票，基本合格(%s)票，不合格(%s)票。（2）服务对象评议结果：发出推荐票(%s)张，收回(%s)张，有效票(%s)张。该同志得票为：优秀(%s)票，合格(%s)票，基本合格(%s)票，不合格(%s)票。",
				r[2], r[4], r[5], r[6], r[7], r[8], r[9], r[10], r[11], r[12], r[14], r[15], r[16], r[17], r[18], r[19], r[20])
		} else {
			e.PublicRemark = fmt.Sprintf("%s，，报社一线考核群众评议会，对%s一般工作人员进行群众评议。应参加人数(%s)人，实参加人数(%s)人。（1）社直评议结果：发出推荐票(%s)张，收回(%s)张，有效票(%s)张。该同志得票为：优秀(%s)票，合格(%s)票，基本合格(%s)票，不合格(%s)票。",
				r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9], r[10], r[11], r[12])
		}
		// 群众评分
		e.PublicEvaluation = r[13]
		// 考核量化分重新计算
		total, err := model.SumMarks(e.StartDate, e.EndDate, map[string]interface{}{"userId": idmap[r[0]], "checked": 1})
		if err != nil {
			buff.WriteString(fmt.Sprintf("第%d行用户[%s]计算考核量化分失败:%s,请稍后再试\n", i+1, r[0], err.Error()))
		}
		total2, _ := strconv.ParseFloat(total, 32)
		e.Marks = fmt.Sprintf("%.2f", total2+basemarks)
		// 总分重新计算
		// 查询用户所在部门的经营属性
		attriute, err := FindDepartAttribute(map[string]interface{}{"id": departmentidMap[r[0]]})
		if err != nil {
			buff.WriteString(fmt.Sprintf("第%d行用户[%s]查询用户部门属性时失败:%s,请联系管理员\n", i+1, r[0], err.Error()))
		}
		err = e.GenerateTotal(attriute, assessmentMarks)
		if err != nil {
			buff.WriteString(fmt.Sprintf("第%d行用户[%s]计算总分失败:%s,请稍后再试\n", i+1, r[0], err.Error()))
		}
		err = e.Updates()
		if err != nil {
			haserr = true
			buff.WriteString(fmt.Sprintf("第%d行用户[%s]更新数据失败:%s,请稍后再试\n", i+1, r[0], err.Error()))
		}

	}

	if haserr {
		return fmt.Errorf(buff.String())
	}
	return nil
}

// SaveEvaluation 存储考核数据
func SaveEvaluation(c *model.Container) error {
	// 获取参数
	if len(c.Body.Params) == 0 {
		return fmt.Errorf(`参数格式:{"body":"params":{}} params为需要保存的数据`)
	}
	// 转换格式
	e := model.ResEvaluation{}
	err := e.FromMap(c.Body.Params)
	if err != nil {
		return err
	}
	// 保存数据
	err = e.FirstOrCreate()
	if err != nil {
		return err
	}
	return nil
}
