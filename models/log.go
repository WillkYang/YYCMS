package models

import (
	"YYCMS/utils/YYLog"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

type Log struct {
	//Id
	Id          int64  `orm:"column(Id);pk;auto"`
	CreateTime  string `orm:"column(CreateTime)"`
	Title       string `orm:"column(Title);size(500)"`
	Description string `orm:"column(Description);size(5000)"`
}

func (log *Log) TableName() string {
	return "log"
}

func CreateOneLog(title, description string) {
	log := &Log{Title: title, Description: description, CreateTime: egb.TimeNowUnix()}
	if _, _, err := ormer().ReadOrCreate(log, "Title", "Description", "CreateTime"); err != nil {
		YYLog.Error(err)
	}
}

func ReadAllLogs(pagesize, offset int) ([]orm.Params, error) {
	var result []orm.Params
	if _, err := ormer().QueryTable(&Log{}).Exclude("Title", "yeeyun_root").OrderBy("-CreateTime").Limit(pagesize, offset).Values(&result); err != nil {
		YYLog.Error(err)
		return nil, err
	} else {
		return result, nil
	}
}

func GetAllLogsNums() int64 {
	if count, err := ormer().QueryTable(&Log{}).Exclude("Title", "yeeyun_root").Count(); err != nil {
		YYLog.Error(err)
		return 0
	} else {
		return count
	}
}
