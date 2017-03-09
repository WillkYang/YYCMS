package models

import (
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

//单页面表
type SinglePage struct {
	//主键
	Id int `orm:"column(Id);pk;auto"`
	//栏目Id
	CateId int `orm:"column(CateId)"`
	//名字
	Title string `orm:"column(Title)"`
	//描述
	Description string `orm:"column(Description);null"`
	//创建时间
	CreateTime string `orm:"column(CreateTime)"`
	//更新时间
	UpdateTime string `orm:"column(UpdateTime)"`
	//内容
	Content string `orm:"column(Content);type(text)"`
}

func (t *SinglePage) TableName() string {
	return "singlepage"
}

//ReadOrCreateOneSinglePage 读取一条单页面数据信息(不存在则创建)
//@params	cateId
//@return	*SinglePage,error
func ReadOrCreateOneSinglePage(cateId int) (*SinglePage, error) {
	singlePage := new(SinglePage)
	singlePage.CateId = cateId
	if err := orm.NewOrm().Read(singlePage, "CateId"); err != nil {
		if err == orm.ErrNoRows {
			newPage := new(SinglePage)
			newPage.CateId = cateId
			newPage.CreateTime = egb.TimeNowUnix()
			newPage.UpdateTime = egb.TimeNowUnix()
			if _, err := orm.NewOrm().Insert(newPage); err != nil {
				YYLog.Error("ReadOrCreateOneSinglePage : " + err.Error())
				return nil, fmt.Errorf(ErrInfo[DataBaseUpdateError])
			}
			return ReadOrCreateOneSinglePage(newPage.CateId)
		}
		YYLog.Error("ReadOrCreateOneSinglePage : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return singlePage, nil
}

//UpdateOneSinglePage 更新单页面数据
//@params	cateId sourceId title content
//@return	error
func UpdateOneSinglePage(cateId int, title, content string) error {
	singlePage, err := ReadOrCreateOneSinglePage(cateId)
	if err != nil {
		return err
	}
	singlePage.Title = title
	singlePage.Content = content
	singlePage.UpdateTime = egb.TimeNowUnix()
	if _, err := orm.NewOrm().Update(singlePage); err != nil {
		YYLog.Error("UpdateOneSinglePage : " + err.Error())
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}
