package model

import "fmt"

// Flowdata 流程数据
type Flowdata struct {
	Model
	ProcessInstanceID string `gorm:"column:processInstanceId" json:"processInstanceId,omitempty"`
	Data              string `gorm:"size:1024" json:"data"`
}

// Updates 只会更新更改的和非空白字段
func (f *Flowdata) Updates() error {
	return db.Model(f).Where("processInstanceId=?", f.ProcessInstanceID).Updates(f).Error
}

// FirstOrCreate 存在就更新，不存在就插入
func (f *Flowdata) FirstOrCreate() error {
	if len(f.ProcessInstanceID) == 0 {
		return fmt.Errorf("processInstanceId不能为空")
	}
	err := db.Where(Flowdata{ProcessInstanceID: f.ProcessInstanceID}).Assign(f).FirstOrCreate(&Flowdata{}).Error
	if err != nil {
		return fmt.Errorf("保存流程数据:%s", err.Error())
	}
	return nil

}

// Create 创建
func (f *Flowdata) Create() error {
	return f.Create()
}

// FindFlowdataByProcessInstanceID 根据流程ID查询流程数据
func FindFlowdataByProcessInstanceID(processInstanceID interface{}) (*Flowdata, error) {
	var fd Flowdata
	err := db.Where("processInstanceId=?", processInstanceID).Find(&fd).Error
	return &fd, err
}

// DelFlowdataByProcessInstanceID DelFlowdataByProcessInstanceID
func DelFlowdataByProcessInstanceID(processInstanceID interface{}) error {
	return db.Where("processInstanceId=?", processInstanceID).Delete(&Flowdata{}).Error
}
