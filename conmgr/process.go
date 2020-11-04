package conmgr

import (
	"errors"
	"fmt"
	"time"

	"github.com/codepository/yxkh/model"
	"github.com/mumushuiding/util"
)

// YXKHYdkh 月度考核
var YXKHYdkh = "月度考核"

// YXKHNdkh 年度考核
var YXKHNdkh = "年度考核"

// YXKHBnkh 半年考核
var YXKHBnkh = "半年考核"

// Zrqd 责任清单
var Zrqd = "责任清单"

// StartProcess 启动流程
func StartProcess(c *model.Container) error {
	errstr := `参数缺一不可，参数格式:{"header":{"token":""},"body":{"params":{"perform":1,"title":"张三-6月一线考核","templateId":"002d5df2a737dd36a2e78314b07d0bb1_1591669930","businessType":"月度考核","data":{}}}} perform:1表示保存先不提交,默认值0表示直接提交`
	// 参数检查
	// 判断token
	if len(c.Header.Token) == 0 {
		return errors.New(errstr)
	}
	// 判断任务类型
	var businessType string
	var data map[string]interface{}
	if c.Body.Params["businessType"] == nil || len(c.Body.Params["businessType"].(string)) == 0 {
		return errors.New(errstr)
	}
	if c.Body.Params["templateId"] == nil || len(c.Body.Params["templateId"].(string)) == 0 {
		return errors.New(errstr)
	}
	if c.Body.Params["title"] == nil || len(c.Body.Params["title"].(string)) == 0 {
		return errors.New(errstr)
	}
	if c.Body.Params["data"] == nil {
		return errors.New(errstr)
	}
	businessType = c.Body.Params["businessType"].(string)
	// templateID = c.Body.Params["templateId"].(string)
	// title = c.Body.Params["title"].(string)
	data = c.Body.Params["data"].(map[string]interface{})
	if len(data) == 0 {
		return errors.New("data 不能为空")
	}
	// 检测远程用户端口是否可用
	ok := IsUserAPIAlive()
	if !ok {
		return errors.New("远程用户API接口不可用")
	}
	// 启动流程：调用远程接口，返回流程ID
	flow, err := StartFlowByToken(c.Header.Token, c.Body.Params)
	if err != nil {
		return fmt.Errorf("启动流程失败:%s", err.Error())
	}
	f := flow.(map[string]interface{})
	flowdata := f["data"].(map[string]interface{})
	if flowdata["ThirdNo"] == nil && len(flowdata["ThirdNo"].(string)) == 0 {
		return fmt.Errorf("流程返回数据ThirdNo为空")
	}
	data["processInstanceId"] = flowdata["ThirdNo"].(string)
	// log.Println("流程ID", flowdata["ThirdNo"].(string))
	// 存储考核数据,若失败则删除远程流程
	err = FirstOrCreateFlowData(businessType, data)
	if err != nil {
		// 删除流程数据
		DeleteFlowByID(flowdata["ThirdNo"])
		return fmt.Errorf("保存数据失败:%s", err.Error())
	}
	// 是否自动提交
	commit := true
	if c.Body.Params["perform"] != nil {
		perform, _ := util.Interface2Int(c.Body.Params["perform"])
		if perform == 1 {
			commit = false
		}
	}
	// log.Println("perform:", c.Body.Params["perform"])
	// log.Println("commit:", commit)
	c.Body.Data = append(c.Body.Data, f["data"])
	if commit {

		c.Body.Params["thirdNo"] = flowdata["ThirdNo"].(string)
		c.Body.Params["perform"] = 2
		result, err := CompleteProcessByToken(c.Header.Token, c.Body.Params)

		if err != nil {
			return fmt.Errorf("流程提交失败:%s", err.Error())
		}
		c.Body.Data = append(c.Body.Data, result...)
	}

	return nil
}

// CompleteProcess 审批流程
func CompleteProcess(c *model.Container) error {
	errstr := `参数格式:{"header":{"token":""},"body":{"params":{"businessType":"月度考核","data":{"processInstanceId":"B4MttX5xHnfKdf"},"perform":2,"speech":""}}}`
	if c.Body.Params["businessType"] == nil || len(c.Body.Params["businessType"].(string)) == 0 {
		return fmt.Errorf("businessType不能为空,%s", errstr)
	}
	// 检测远程用户端口是否可用
	ok := IsUserAPIAlive()
	if !ok {
		return errors.New("远程用户API接口不可用")
	}

	// 更新数据
	var data map[string]interface{}
	var thirdNo string
	if c.Body.Params["data"] != nil {
		data = c.Body.Params["data"].(map[string]interface{})
		thirdNo = data["processInstanceId"].(string)
	}
	if len(thirdNo) == 0 {
		return fmt.Errorf("processInstanceId 不能为空,%s", errstr)
	}
	if len(data) > 0 {
		err := FirstOrCreateFlowData(c.Body.Params["businessType"].(string), data)
		if err != nil {
			return fmt.Errorf("更新流程数据失败:%s", err.Error())
		}
	}
	c.Body.Params["thirdNo"] = thirdNo
	result, err := CompleteProcessByToken(c.Header.Token, c.Body.Params)
	if err != nil {
		return fmt.Errorf("流程提交失败:%s", err.Error())
	}
	// 一线考核流程结束之后让分数生效，以及根据评价进行扣分
	process := result[0].(map[string]interface{})
	completed, err := util.Interface2Int(process["completed"])
	if err != nil {
		return err
	}
	// 这里涉及分布式事务,需要用到日志进行辅助
	if completed == 1 {
		data, err := ExecWhenProcessCompleted(process["processInstanceId"].(string))
		if err != nil {
			errlog := &model.ErrLog{}
			errlog.CreateTime = time.Now()
			errlog.BusinessType = process["businessType"].(string)
			errlog.Data = data.ToString()
			errlog.Err = err.Error()
			errlog.Key = process["processInstanceId"].(string)
			err1 := errlog.Create()
			if err1 != nil {
				return fmt.Errorf("保存错误日志失败，请通知管理员，原因:%s", err1.Error())
			}
			return fmt.Errorf("执行流程结束程序失败，请通知管理员，原因:%s", err.Error())
		}
	}
	c.Body.Data = append(c.Body.Data, result...)
	return nil
}

// ExecWhenProcessCompleted 当流程结束时
func ExecWhenProcessCompleted(processInstanceID string) (*ReExecData, error) {
	// 查询流程
	println("执行ExecWhenProcessCompleted函数")
	red := &ReExecData{Key: processInstanceID, FuncName: "ExecWhenProcessCompleted"}
	ps, err := FindAllProcess(map[string]interface{}{"processInstanceId": processInstanceID})
	if err != nil {
		return red, fmt.Errorf("查询流程报错:%s", err.Error())
	}
	if len(ps) == 0 {
		return red, fmt.Errorf("流程%s不存在", processInstanceID)
	}
	switch ps[0].BusinessType {
	// 月度考核
	case YXKHYdkh:
		// 根据组织考核结果进行评分,data作用是保存需要存储的数据和类型
		data, err := RemarkEvaluationByProcessInstanceID(processInstanceID)
		if err != nil {
			return data, fmt.Errorf("月度考核评分报错:%s", err.Error())
		}
		// 项目加分设置为已经确认
		data, err = CheckProjectByProcessInstanceID(processInstanceID)
		if err != nil {
			return data, fmt.Errorf("项目确认失败:%s", err.Error())
		}
		break
	default:
	}
	return nil, nil
}

// FirstOrCreateFlowData 保存流程数据
func FirstOrCreateFlowData(businessType string, data map[string]interface{}) error {
	switch businessType {
	// 月度考核
	case YXKHYdkh:
		e := model.ResEvaluation{}
		err := e.FromMap(data)
		if err != nil {
			return err
		}
		err = e.FirstOrCreate()
		if err != nil {
			return err
		}
		break
	// 年度考核
	case YXKHNdkh:
		e := model.ResEvaluation{}
		err := e.FromMap(data)
		if err != nil {
			return err
		}
		err = e.FirstOrCreate()
		if err != nil {
			return err
		}
		break
	// 半年考核
	case YXKHBnkh:
		e := model.ResEvaluation{}
		err := e.FromMap(data)
		if err != nil {
			return err
		}
		err = e.FirstOrCreate()
		if err != nil {
			return err
		}
		break
	// 责任清单
	case Zrqd:
		break
	default:
		return fmt.Errorf("流程类型[%s]不存在", businessType)
	}
	return nil

}

// DelFlow 删除流程
func DelFlow(c *model.Container) error {
	errStr := fmt.Sprintf(`参数格式{"body":{"params":{"processInstanceId":"a","businessType":"月度考核"}}}`)
	if c.Body.Params == nil || len(c.Body.Params) == 0 {
		return errors.New(errStr)
	}
	if c.Body.Params["processInstanceId"] == nil {
		return fmt.Errorf("processInstanceId 不能为空")
	}
	if c.Body.Params["businessType"] == nil {
		return fmt.Errorf("businessType 不能为空")
	}
	businessType := c.Body.Params["businessType"].(string)
	// 删除流程
	var err error
	err = DeleteFlowByID(c.Body.Params["processInstanceId"])
	if err != nil {
		return err
	}
	// 删除流程数据
	switch businessType {
	// 月度考核
	case YXKHYdkh:
		err = model.DelEvaluationByProcessInstanceID(c.Body.Params["processInstanceId"])
		break
	// 年度考核
	case YXKHNdkh:
		err = model.DelEvaluationByProcessInstanceID(c.Body.Params["processInstanceId"])
		break
	// 半年考核
	case YXKHBnkh:
		err = model.DelEvaluationByProcessInstanceID(c.Body.Params["processInstanceId"])
		break
	// 责任清单
	case Zrqd:
		break
	default:
		return fmt.Errorf("流程类型[%s]不存在", businessType)
	}
	return err
}

// FindFlowDatas 查询流程数据
func FindFlowDatas(c *model.Container) error {
	errStr := fmt.Sprintf(`参数格式{"body":{"params":{"processInstanceId":"a","businessType":"月度考核"}}}`)
	if c.Body.Params == nil || len(c.Body.Params) == 0 {
		return errors.New(errStr)
	}
	if c.Body.Params["processInstanceId"] == nil {
		return fmt.Errorf("processInstanceId 不能为空")
	}
	if c.Body.Params["businessType"] == nil {
		return fmt.Errorf("businessType 不能为空")
	}
	businessType := c.Body.Params["businessType"].(string)
	var data interface{}
	var err error
	switch businessType {
	// 月度考核
	case YXKHYdkh:
		data, err = model.FindSingleEvaluation(c.Body.Params["processInstanceId"])
		break
	// 年度考核
	case YXKHNdkh:
		data, err = model.FindSingleEvaluation(c.Body.Params["processInstanceId"])
		break
	// 半年考核
	case YXKHBnkh:
		data, err = model.FindSingleEvaluation(c.Body.Params["processInstanceId"])
		break
	// 责任清单
	case Zrqd:
		break
	default:
		return fmt.Errorf("流程类型[%s]不存在", businessType)
	}
	if err != nil {
		return err
	}
	c.Body.Data = append(c.Body.Data, data)
	return nil

}
