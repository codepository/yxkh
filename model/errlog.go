package model

import "github.com/jinzhu/gorm"

// ErrLog 错误日志
type ErrLog struct {
	Model
	// 定位关键词
	Key string
	// 错误类型
	BusinessType string `gorm:"column:businessType" json:"businessType,omitempty"`
	// 包含数据,用于重新执行
	Data string `gorm:"size:1000;" json:"data"`
	// 具体错误
	Err string `gorm:"size:500" json:"err"`
}

// Create 创建
func (e *ErrLog) Create() error {
	if len(e.Data) > 999 {
		e.Data = e.Data[:999]
	}
	if len(e.Err) > 499 {
		e.Err = e.Err[:499]
	}
	return db.Create(e).Error
}

// FindErrLog 查询
func FindErrLog(query interface{}, value ...interface{}) ([]*ErrLog, error) {
	var datas []*ErrLog
	err := db.Where(query, value).Find(&datas).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return make([]*ErrLog, 0), nil
	}
	return datas, err
}

// DelByTx 删除日志
func (e *ErrLog) DelByTx(tx *gorm.DB) error {
	return tx.Delete(e).Error
}

// DelByID 删除日志
func (e *ErrLog) DelByID() error {
	return db.Delete(e).Error
}
