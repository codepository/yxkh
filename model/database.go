package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/codepository/yxkh/config"
	"github.com/jinzhu/gorm"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

// 配置
var conf = *config.Config

// Model 基类
type Model struct {
	ID         int       `gorm:"primary_key" json:"ID,omitempty"`
	CreateTime time.Time `gorm:"column:createTime" json:"createTime"`
}

// StartDB 启动数据库
func StartDB() {
	setup()
}

// StopDB 关闭数据库
func StopDB() {
	CloseDB()
	log.Println("关闭数据库")
}

// Setup 初始化一个db连接
func setup() {
	var err error
	log.Println("数据库初始化！！")
	db, err = gorm.Open(conf.DbType, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.DbUser, conf.DbPassword, conf.DbHost, conf.DbPort, conf.DbName))
	if err != nil {
		log.Fatalf("数据库连接失败 err: %v", err)
	}
	// 启用Logger，显示详细日志
	mode, _ := strconv.ParseBool(conf.DbLogMode)

	db.LogMode(mode)
	db.SingularTable(true) //全局设置表名不可以为复数形式
	// db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	idle, err := strconv.Atoi(conf.DbMaxIdleConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(idle)
	open, err := strconv.Atoi(conf.DbMaxOpenConns)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(open)

	db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(&ResProject{}).AutoMigrate(&ResMark{}).AutoMigrate(&Label{}).AutoMigrate(&UserLabel{}).
		AutoMigrate(&InfoDic{}).AutoMigrate(&Uploadfile{})
	db.Model(&ResProject{}).ModifyColumn("startDate", "date").ModifyColumn("endDate", "date").AddIndex("projectidindex", "projectId")
	db.Model(&ResMark{}).ModifyColumn("startDate", "date").ModifyColumn("endDate", "date").AddForeignKey("projectId", "res_project(projectId)", "CASCADE", "CASCADE")
	db.Model(&UserLabel{}).AddIndex("labelidindex", "label_id").AddForeignKey("label_id", "label(id)", "CASCADE", "CASCADE")
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}

// GetDB getdb
func GetDB() *gorm.DB {
	return db
}

// GetTx GetTx
func GetTx() *gorm.DB {
	return db.Begin()
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName("CreateTime"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		// if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
		// 	if modifyTimeField.IsBlank {
		// 		modifyTimeField.Set(nowTime)
		// 	}
		// }
	}
}
