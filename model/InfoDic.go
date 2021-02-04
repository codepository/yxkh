package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// InfoDic 字典
type InfoDic struct {
	Model
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Type2       string `json:"type2"`
	Description string `json:"description"`
}

// FindSingleDict FindSingleDict
func FindSingleDict(query interface{}) (*InfoDic, error) {
	info := []*InfoDic{}
	err := db.Where(query).Find(&info).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询info_dic:%s", err.Error())
	}
	if len(info) > 0 {
		return info[0], nil
	}
	return nil, nil
}

// FindAllInfoDic FindAllInfoDic
func FindAllInfoDic(query interface{}, values ...interface{}) ([]*InfoDic, error) {
	info := []*InfoDic{}
	err := db.Where(query, values...).Find(&info).Error
	if err == gorm.ErrRecordNotFound {
		return make([]*InfoDic, 0), nil
	}
	return info, err
}

// FindAllInfoDicPaged 分页查询
func FindAllInfoDicPaged(query interface{}, limit, offset int) ([]*InfoDic, error) {
	info := []*InfoDic{}
	if limit == 0 {
		limit = 20
	}
	err := db.Where(query).Limit(limit).Offset(offset).Find(&info).Error
	if err == gorm.ErrRecordNotFound {
		return make([]*InfoDic, 0), nil
	}
	return info, err
}

// Updates 更新非空字段
func (d *InfoDic) Updates() error {
	return db.Model(d).Where("ID=?", d.ID).Updates(d).Error
}

// FirstOrCreate 存在就更新，不存在就创建
func (d *InfoDic) FirstOrCreate() error {
	d.CreateTime = time.Now()
	// value := d.Value
	return db.Where(InfoDic{
		Name:  d.Name,
		Type:  d.Type,
		Type2: d.Type2,
		Value: d.Value,
	}).Attrs(d).FirstOrCreate(d).Error
}

// DeleteDicsIDs 根据id批量删除
func DeleteDicsIDs(ids interface{}) error {
	return db.Where(ids).Delete(InfoDic{}).Error
}

// FindMarksBase 查询加减分基础分
func FindMarksBase() (string, error) {
	datas, err := FindAllInfoDic(map[string]interface{}{"name": "基础分", "type": "基本定格"})
	if err != nil {
		return "", err
	}
	if len(datas) > 1 {
		return "", fmt.Errorf("info_dic表中存在[%d]个,name='基础分'的结果,请删至一个", len(datas))
	}
	return datas[0].Value, nil

}
