package conmgr

import (
	"encoding/json"
	"errors"

	"github.com/codepository/yxkh/model"
)

// GetDataFromUserAPI GetDataFromUserAPI
func GetDataFromUserAPI(method string, params map[string]interface{}) (interface{}, error) {
	var par model.Container
	par.Body.Method = method
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
	par.Body.Data = append(par.Body.Data, params)
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

// FindUsersUncompleteTask 未完成任务的用户
func FindUsersUncompleteTask(taskname, start string) (interface{}, error) {
	method := "visit/task/uncomplete"
	params := map[string]interface{}{"task": taskname, "start": start}
	result, err := GetDataFromUserAPI(method, params)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
}

// FindTaskCompleteRates 查询任务完成率
func FindTaskCompleteRates(taskname, start string) (interface{}, error) {
	method := "visit/task/completeRate"
	params := map[string]interface{}{"task": taskname, "start": start, "max_results": 20}
	return GetDataFromUserAPI(method, params)
}

// FindUseridsByTagsAndLevel 根据标签和职级查询用户id
// tagMethod 有两值,and 表示同时包含多个指定标签，or 表示包含任一标签
func FindUseridsByTagsAndLevel(tags []string, tagMethods []string, level []int) ([]interface{}, error) {
	method := "visit/user/getUseridsByTagAndLevel"
	params := map[string]interface{}{"tags": tags, "methods": tagMethods, "levels": level}
	result, err := GetDataFromUserAPI(method, params)
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
	result, err := GetDataFromUserAPI(method, params)
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
