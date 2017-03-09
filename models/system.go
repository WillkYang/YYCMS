package models

import (
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/astaxie/beego/orm"
)

//系统信息表
type System struct {
	//主键
	Id int `orm:"column(Id);pk;auto"`
	//访问量
	HitTime int64 `orm:"column(HitTime)"`
	//名字
	Title string `orm:"column(Title)"`
	//描述
	Description string `orm:"column(Description);null"`
}

func (t *System) TableName() string {
	return "system"
}

func AutoAddHitTime() {
	ormer().QueryTable(&System{}).Update(orm.Params{
		"HitTime": orm.ColValue(orm.ColAdd, 1),
	})
}

//ReadOrCreateOneSystem 读取一条单页面数据信息(不存在则创建)
//@params	cateId
//@return	*System,error
func ReadOrCreateOneSystem() (*System, error) {
	system := new(System)
	if err := ormer().QueryTable(&System{}).One(system); err != nil {
		if err == orm.ErrNoRows {
			system := new(System)
			if _, err := orm.NewOrm().Insert(system); err != nil {
				YYLog.Error("ReadOrCreateOneSystem : " + err.Error())
				return nil, fmt.Errorf(ErrInfo[DataBaseUpdateError])
			}
			return ReadOrCreateOneSystem()
		}
		YYLog.Error("ReadOrCreateOneSystem : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return system, nil
}
