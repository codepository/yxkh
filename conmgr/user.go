package conmgr

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// IsUserAPIAlive 远程用户接口是否可用
func IsUserAPIAlive() bool {
	url := model.APIClient.UserAPIURL + "/alive"
	result, err := model.APIClient.Get(url)
	if err != nil {
		return false
	}
	if result == "1" {
		return true
	}
	return false
}

// GetDataFromUserAPI2 GetDataFromUserAPI2
func GetDataFromUserAPI2(par model.Container) (interface{}, error) {
	if par.Body.Params["max_results"] != nil {
		switch par.Body.Params["max_results"].(type) {
		case float64:
			par.Body.MaxResults = int(par.Body.Params["max_results"].(float64))
			break
		case int:
			par.Body.MaxResults = par.Body.Params["max_results"].(int)
			break
		default:
			return nil, errors.New("max_results 必须为数字")
		}
		delete(par.Body.Params, "max_results")
	}
	if par.Body.Params["start_index"] != nil {
		switch par.Body.Params["start_index"].(type) {
		case float64:
			par.Body.StartIndex = int(par.Body.Params["start_index"].(float64))
			break
		case int:
			par.Body.StartIndex = par.Body.Params["start_index"].(int)
			break
		default:
			return nil, errors.New("start_index 必须为数字")
		}
		delete(par.Body.Params, "start_index")
	}
	url := model.APIClient.UserAPIURL + "/getData"
	result, err := model.APIClient.Post(url, par)
	// log.Println("result:", string(result))
	if err != nil {
		return nil, err
	}
	map1 := make(map[string]interface{})
	err = json.Unmarshal(result, &map1)
	if err != nil {
		return nil, err
	}
	if map1["status"].(float64) != 200 {
		return nil, errors.New(map1["message"].(string))
	}
	msg := map1["message"].(map[string]interface{})
	body := msg["body"].(map[string]interface{})
	// err = util.Str2Struct(map1["message"], &d)
	// if err != nil {
	// 	return nil, err
	// }
	return body["data"], nil
}

// GetDataFromUserAPI GetDataFromUserAPI
func GetDataFromUserAPI(token, method string, params map[string]interface{}) (interface{}, error) {
	var par model.Container
	par.Body.Method = method
	par.Header.Token = token
	if params["max_results"] != nil {
		switch params["max_results"].(type) {
		case float64:
			par.Body.MaxResults = int(params["max_results"].(float64))
			break
		case int:
			par.Body.MaxResults = params["max_results"].(int)
			break
		default:
			return nil, errors.New("max_results 必须为数字")
		}
		delete(params, "max_results")
	}
	if params["start_index"] != nil {
		switch params["start_index"].(type) {
		case float64:
			par.Body.StartIndex = int(params["start_index"].(float64))
			break
		case int:
			par.Body.StartIndex = params["start_index"].(int)
			break
		default:
			return nil, errors.New("start_index 必须为数字")
		}
		delete(params, "start_index")
	}
	par.Body.Params = params
	url := model.APIClient.UserAPIURL + "/getData"
	result, err := model.APIClient.Post(url, par)
	// log.Println("result:", string(result))
	if err != nil {
		return nil, err
	}
	map1 := make(map[string]interface{})
	err = json.Unmarshal(result, &map1)
	if err != nil {
		return nil, err
	}
	if map1["status"].(float64) != 200 {
		return nil, errors.New(map1["message"].(string))
	}
	msg := map1["message"].(map[string]interface{})
	body := msg["body"].(map[string]interface{})
	// err = util.Str2Struct(map1["message"], &d)
	// if err != nil {
	// 	return nil, err
	// }
	return body["data"], nil
}

// StartFlowByToken 启动流程
func StartFlowByToken(token string, params map[string]interface{}) ([]interface{}, error) {
	method := "exec/flow/startByToken"
	result, err := GetDataFromUserAPI(token, method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	// if datas[0] == nil {
	// 	return nil, nil
	// }

	// return datas[0].(interface{}), nil
	return datas, nil
}

// CompleteProcessByToken 审批流程
func CompleteProcessByToken(token string, params map[string]interface{}) ([]interface{}, error) {
	method := "exec/flow/completeFlowTask"
	result, err := GetDataFromUserAPI(token, method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	return datas, nil

}

// FindallProcess 查询所有流程
func FindallProcess(c *model.Container) error {
	data, err := FindAllProcess(c.Body.Params)
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, data)
	return nil
}

// FindAllProcess 查询所有流程
func FindAllProcess(params map[string]interface{}) ([]*model.Process, error) {
	method := "visit/flow/findall"
	data, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return nil, err
	}
	datas := data.([]interface{})
	if len(datas) == 0 {
		return make([]*model.Process, 0), nil
	}
	var result []*model.Process
	if datas[0] == nil {
		return make([]*model.Process, 0), nil
	}
	for _, d := range datas[0].([]interface{}) {
		// fmt.Println(reflect.TypeOf(d))
		process := &model.Process{}
		str, _ := util.ToJSONStr(d)
		err := util.Str2Struct(str, process)
		if err != nil {
			return nil, err
		}
		result = append(result, process)
	}
	return result, nil
}

// DeleteFlowByID 删除流程
func DeleteFlowByID(id interface{}) error {
	method := "exec/flow/delete"
	params := map[string]interface{}{"ThirdNo": id}
	_, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return err
	}
	return nil
}

// FindUsersUncompleteTask 未完成任务的用户
func FindUsersUncompleteTask(taskname, start string) (interface{}, error) {
	method := "visit/task/uncomplete"
	params := map[string]interface{}{"task": taskname, "start": start}
	result, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindTaskCompleteRates 查询任务完成率
func FindTaskCompleteRates(taskname, start string) (interface{}, error) {
	method := "visit/task/completeRate"
	params := map[string]interface{}{"task": taskname, "start": start, "max_results": 20}
	return GetDataFromUserAPI("", method, params)
}

// FindTaskCompletedDescribe FindTaskCompletedDescribe
// 任务人数，提交情况，审结情况
func FindTaskCompletedDescribe(titleLike string) (interface{}, error) {
	method := "visist/task/completedescribe"
	params := map[string]interface{}{"titleLike": titleLike}
	return GetDataFromUserAPI("", method, params)
}

// FindTaskRank 查询任务审批排行
func FindTaskRank(titleLike string) (interface{}, error) {
	method := "visit/task/userTaskRank"
	params := map[string]interface{}{"titleLike": titleLike}
	return GetDataFromUserAPI("", method, params)
}

// FindPersonApplyYxkh 查询已经提交或未提交流程的用户
func FindPersonApplyYxkh(titleLike string, apply int) (interface{}, error) {
	method := "visit/task/personApplyYxkh"
	params := map[string]interface{}{"titleLike": titleLike, "apply": apply, "limit": 5}
	return GetDataFromUserAPI("", method, params)
}

// FindUseridsByTagsAndLevel 根据标签和职级查询用户id
// tagMethod 有两值,and 表示同时包含多个指定标签，or 表示包含任一标签
func FindUseridsByTagsAndLevel(tags []string, tagMethods []string, level []int) ([]interface{}, error) {
	method := "visit/user/getUseridsByTagAndLevel"
	params := map[string]interface{}{"tags": tags, "methods": tagMethods, "levels": level}
	result, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	if datas[0] == nil {
		return nil, nil
	}
	return datas[0].([]interface{}), nil
}

// FindUseridsByTags 根据标签查询用户id
func FindUseridsByTags(tags []string, tagMethod string) ([]interface{}, error) {
	if len(tagMethod) == 0 {
		tagMethod = "or"
	}
	method := "visit/user/getUseridsByTagAndLevel"
	params := map[string]interface{}{"tags": tags, "methods": []string{tagMethod}}
	result, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	if datas[0] == nil {
		return nil, nil
	}
	return datas[0].([]interface{}), nil
}

// FindDepartments 查询部门
func FindDepartments(params map[string]interface{}) ([]interface{}, error) {
	method := "visit/department/find"
	result, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	if datas[0] == nil {
		return nil, nil
	}
	return datas[0].([]interface{}), nil
}

// FindDepartAttribute 查询部门属性，部门属性有:1-采编经营类，2-行政后勤类
func FindDepartAttribute(params map[string]interface{}) (int, error) {
	datas, err := FindDepartments(params)
	if err != nil {
		return 0, fmt.Errorf("查询部门属性报错:%s", err.Error())
	}
	dept, ok := datas[0].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("返回结果不是json格式")
	}
	att, err := util.Interface2Int(dept["attribute"])
	if err != nil {
		return 0, err
	}
	return att, nil
}

// FindAllUsers 查询所有用户
// 参数格式:{"body":{"metrics":"id,name","params":{"where":"id=1 and name='xx'","userid":"","name":"","departmentid":"","departmentname":"","mobile":"","email":"",}}}metrics为显示的字段,where不为空时忽略其它查询条件
func FindAllUsers(fields string, params map[string]interface{}) ([]interface{}, error) {
	var par model.Container
	par.Body.Method = "visit/user/findAll"
	par.Body.Params = params
	par.Body.Metrics = fields
	result, err := GetDataFromUserAPI2(par)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	if datas[0] == nil {
		return nil, nil
	}
	return datas[0].([]interface{}), nil
}

// GetUseridAndDeptByMobileOrName 根据电话或名字查询用户id和部门id
func GetUseridAndDeptByMobileOrName(mobileOrName string) (int, int, error) {
	fields := "id,departmentid"
	params := make(map[string]interface{})
	// 判断是否是电话
	if util.IsMobile(mobileOrName) {
		params["mobile"] = mobileOrName
	} else {
		params["name"] = mobileOrName
	}
	datas, err := FindAllUsers(fields, params)
	if err != nil {
		return 0, 0, err
	}
	if len(datas) == 0 {
		return 0, 0, fmt.Errorf("用户【%s】不存在", mobileOrName)
	}
	if len(datas) > 1 {
		return 0, 0, fmt.Errorf("用户【%s】存在重名,请用电话重试", mobileOrName)
	}
	user := datas[0].(map[string]interface{})
	id, err := util.Interface2Int(user["id"])
	if err != nil {
		return 0, 0, err
	}
	departmentid, err := util.Interface2Int(user["departmentid"])
	if err != nil {
		return 0, 0, err
	}
	return id, departmentid, nil
}

// GetUseridAndNameByMobileOrName 根据电话或名字查询用户
func GetUseridAndNameByMobileOrName(mobileOrName string) (int, string, error) {
	fields := "id"
	params := make(map[string]interface{})
	// 判断是否是电话
	if util.IsMobile(mobileOrName) {
		params["mobile"] = mobileOrName
	} else {
		params["name"] = mobileOrName
	}
	datas, err := FindAllUsers(fields, params)
	if err != nil {
		return 0, "", err
	}
	if len(datas) == 0 {
		return 0, "", fmt.Errorf("用户【%s】不存在", mobileOrName)
	}
	if len(datas) > 1 {
		return 0, "", fmt.Errorf("用户【%s】存在重名,请用电话重试", mobileOrName)
	}
	user := datas[0].(map[string]interface{})
	id, err := util.Interface2Int(user["id"])
	name := user["name"].(string)
	if err != nil {
		return 0, "", err
	}
	return id, name, nil
}

// FindAllUploadFiles 查询上传的文件
func FindAllUploadFiles(filetype string, fields string) (interface{}, error) {
	method := "visit/uploadfile/find"
	params := map[string]interface{}{"filetype": filetype, "fields": fields}
	return GetDataFromUserAPI("", method, params)
}

// FindAllTags 根据条件查询标签
func FindAllTags(params map[string]interface{}) ([]interface{}, error) {
	method := "visit/lable/find"
	result, err := GetDataFromUserAPI("", method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	datas := result.([]interface{})
	if datas[0] == nil {
		return nil, nil
	}
	return datas[0].([]interface{}), nil
}
