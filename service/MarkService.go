package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// MarkReceiver MarkReceiver
type MarkReceiver struct {
	MarkID      int    `json:"markId"`
	ProjectID   int    `json:"projectId"`
	MarkNumber  string `json:"markNumber"`
	MarkReason  string `json:"markReason"`
	Accordingly string `json:"accordingly"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	CreateTime  string `json:"createTime"`
	UserID      string `json:"userId"`
	Checked     string `json:"checked"`
	PageIndex   int    `json:"pageIndex"`
	PageSize    int    `json:"pageSize"`
}

// RankConfig 排名设置
type RankConfig struct {
	Post []string
}

// Save Save
func (r *MarkReceiver) Save() (int, error) {
	entity, err := r.prepareData()
	if err != nil {
		return 0, err
	}
	err = entity.Save()
	if err != nil {
		return 0, err
	}
	return entity.MarkID, nil
}
func (r *MarkReceiver) prepareData() (*model.ResMark, error) {
	var start, end time.Time
	var err error
	if len(r.StartDate) > 0 {
		start, err = util.ParseDate(r.StartDate, util.YYYY_MM_DD)
		if err != nil {
			return nil, err
		}
	}
	if len(r.EndDate) > 0 {
		end, err = util.ParseDate(r.EndDate, util.YYYY_MM_DD)
		if err != nil {
			return nil, err
		}
	}
	var result = model.ResMark{
		MarkID:      r.MarkID,
		ProjectID:   r.ProjectID,
		MarkNumber:  r.MarkNumber,
		MarkReason:  r.MarkReason,
		Accordingly: r.Accordingly,
		StartDate:   start,
		EndDate:     end,
		UserID:      r.UserID,
		Checked:     r.Checked,
	}
	return &result, nil
}

// FindAll FindAll
func (r *MarkReceiver) FindAll() (string, error) {
	datas, count, err := model.FindAllMarkPaged(r.PageIndex, r.PageSize, r.getSQL())
	if err != nil {
		return "", err
	}
	return util.ToPageJSON(datas, count, r.PageIndex, r.PageSize)
}
func (r *MarkReceiver) getSQL() string {
	var sql string
	if len(r.EndDate) > 0 {
		// maps["endDate"] = p.EndDate
		sql += " and endDate <='" + r.EndDate + "'"
	}
	if len(r.StartDate) > 0 {
		// maps["startDate"] = p.StartDate
		sql += " and startDate>='" + r.StartDate + "'"
	}
	if r.ProjectID != 0 {
		// maps["projectId"] = p.ProjectId
		sql += " and projectId=" + fmt.Sprintf("%d", r.ProjectID)
	}
	if len(r.UserID) > 0 {
		sql += " and userId=" + r.UserID
	}
	if len(sql) > 0 && sql[0:4] == " and" {
		return sql[5:]
	}
	return sql
}

// DelMarkByID DelMarkByID
func DelMarkByID(id int) error {

	return model.DelMarkByID(id)
}

// Update Update
func (r *MarkReceiver) Update() error {
	entity, err := r.prepareData()
	if err != nil {
		return err
	}
	return entity.Update()
}

type userLabelWithGroup struct {
	userLabelType  string
	userLabelName  string
	groupLabelID   int
	groupLabelName string
}

// FindUserMarkRank 加减分排行
func FindUserMarkRank(startDate string, endDate string) (string, error) {
	// ---------- 用户类型标签获取 ----------
	var userLables []*model.Label
	// 从rank.json读取用户职级
	postLabels, err := getPostLabelFromRankFile()
	if err != nil {
		return "", nil
	}
	// 从label表查询加减分排行员工分类
	rankLabels, err := getUserLabelNeedRank()
	if err != nil {
		return "", nil
	}
	userLables = append(userLables, postLabels...)
	userLables = append(userLables, rankLabels...)
	if len(userLables) == 0 {
		return "", errors.New("用户分类标签为空！")
	}
	// ----------- 考核组标签获取 -----------
	groupLabels, err := getGroupLables()
	if err != nil {
		return "", nil
	}
	if len(groupLabels) == 0 {
		return "", errors.New("考核组为空！")
	}
	// ----------- 用户类型标签和考核组标签组合 --------------
	total := len(userLables) * len(groupLabels)
	inChannel := make(chan *userLabelWithGroup, total)
	for _, u := range userLables {
		for _, g := range groupLabels {
			ug := &userLabelWithGroup{
				userLabelType:  u.Type,
				userLabelName:  u.Name + "%",
				groupLabelID:   g.ID,
				groupLabelName: g.Name,
			}
			inChannel <- ug
		}
	}
	for x := range inChannel {
		// fmt.Printf("type:%s  name:%s  groupID: %d groupName:%s\n", x.userLabelType, x.userLabelName, x.groupLabelID, x.groupLabelName)
		userSQL := getUserSQL(x)
		userRankSQL := fmt.Sprintf("select r.userId,u.user_name,ifnull(round(sum(r.markNumber),2),0) as mark from res_mark r join (%s) u on u.user_id=r.userId and r.startDate>='%s' and r.endDate<='%s' group by r.userId order by mark desc", userSQL, startDate, endDate)
		fmt.Println(userRankSQL)

	}
	// ----------- 查询 ---------------------

	// ---------- 结果合并生成json数组------------
	return "", nil
}

// getUserSQL getUserSQL
func getUserSQL(ug *userLabelWithGroup) string {
	return fmt.Sprintf("select ul.user_id, ul.user_name from user_label ul where ul.label_id in (select id from label where label_id=%d or (type='%s' and name like '%s')) group by ul.user_id having count(ul.user_id)=2", ug.groupLabelID, ug.userLabelType, ug.userLabelName)
}

// getPostLabelFromRankFile 从rank.json读取用户职级
func getPostLabelFromRankFile() ([]*model.Label, error) {
	file, err := os.Open("rank.json")
	if err != nil {
		panic(err)
	}
	var rc RankConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&rc)
	if err != nil {
		return nil, err
	}
	var result []*model.Label
	for _, labelName := range rc.Post {
		label := &model.Label{
			Name: labelName,
			Type: "职级",
		}
		result = append(result, label)
	}
	return result, nil
}

// getUserLabelNeedRank 从label表查询 type='加减分排行员工分类'
func getUserLabelNeedRank() ([]*model.Label, error) {
	result, err := model.FindLabelByType("加减分排行员工分类")
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getGroupLables 从label表获取type="考核组"的标签
func getGroupLables() ([]*model.Label, error) {
	result, err := model.FindLabelByType("考核组")
	if err != nil {
		return nil, err
	}
	return result, nil
}
