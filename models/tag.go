package models

import (
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

//标签表
type Tag struct {
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
	//排序值
	Sort int `orm:"column(Sort);default(0)"`
	//点击量
	Hit int64 `orm:"column(Hit);default(0)"`
	//下载量
	PicsDownload int64 `orm:"column(PicsDownload);default(0)"`
}

func (t *Tag) TableName() string {
	return "tag"
}

//GetOneTagById 获取一个标签
//@params	Id
//@return	*Tag
func GetOneTagById(Id int) (*Tag, error) {
	tag := new(Tag)
	tag.Id = Id

	if err := ormer().Read(tag, "Id"); err != nil {
		YYLog.Error("GetOneTagById : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return tag, nil
}

//GetTagsNum 获取标签的数量
//@params	cateId keyword(搜索title的关键词)
//@return	int
func GetTagsNum(catId int, keyword string) int {
	var data []orm.Params
	sql := "SELECT Id FROM tag WHERE CateId = ? AND Title LIKE ? "
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword).Values(&data)
	return len(data)
}

//GetTags 获取标签
//@params	catId keyword(搜索title的关键词) pagesize offset
//@return	[]orm.Params
func GetTags(catId int, keyword string, pagesize, offset int) []orm.Params {

	var data []orm.Params
	sql := `SELECT * FROM tag WHERE CateId = ? AND Title LIKE ? ORDER BY tag.Sort DESC LIMIT ?,?`
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword, offset, pagesize).Values(&data)
	return data
}

//GetOneTagByTitle 通过账号获取一个管理员用户
//@params   title
//@return   *Tag
func GetOneTagByTitle(cateId int, title string) (*Tag, error) {
	tag := new(Tag)
	tag.CateId = cateId
	tag.Title = title
	if err := ormer().Read(tag, "CateId", "Title"); err != nil {
		return nil, err
	}
	return tag, nil
}

//CreateOneTag 新建一篇标签
//@params	catId title thumb source description content
//@return	error
func CreateOneTag(cateId int, title, description string) (int, error) {
	tag := new(Tag)
	tag.CateId = cateId
	tag.Title = title
	tag.Description = description
	tag.CreateTime = egb.TimeNowUnix()

	if index, err := ormer().Insert(tag); err != nil {
		YYLog.Error("CreateOneTag : " + err.Error())
		return 0, fmt.Errorf(ErrInfo[SystemError])
	} else {
		return int(index), nil
	}
}

//UpdateTag 更新标签
//@params	id title description
//@return	error
func UpdateTag(id, cateId int, title, description string) error {
	tag, err := GetOneTagById(id)
	if err != nil {
		return err
	}

	tag.CateId = cateId
	tag.Title = title
	tag.Description = description

	if _, err := ormer().Update(tag); err != nil {
		YYLog.Error("UpdateTag : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		return nil
	}
}

//UpdateTagSort 更新标签排序
//@params	id sort
//@return	error
func UpdateTagSort(id, sort int) error {
	tag, err := GetOneTagById(id)
	if err != nil {
		return fmt.Errorf(ErrInfo[SystemError])
	}
	tag.Sort = sort
	if _, err := ormer().Update(tag, "Sort"); err != nil {
		YYLog.Error("UpdateTagSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneTag 删除一个标签
//@params	id
//@return	error
func DeleteOneTag(id int) error {
	tag, err := GetOneTagById(id)
	if err != nil {
		return err
	}
	if _, err := ormer().Delete(tag); err != nil {
		YYLog.Error("DeleteOneTag : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}
