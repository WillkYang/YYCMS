package models

import (
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

//图片表
type Picture struct {
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
	//发布时间
	PublishTime string `orm:"column(PublishTime)"`
	//路径
	Path string `orm:"column(Path)"`
	//排序值
	Sort int `orm:"column(Sort);default(0)"`
	//状态 - 当值为-1时表示数据已删除
	Status int `orm:"column(Status);default(0)"`
	//点击量
	Hit int64 `orm:"column(Hit);default(0)"`
	//下载量
	Download int64 `orm:"column(Download);default(0)"`
}

func (t *Picture) TableName() string {
	return "picture"
}

//GetOnePictureById 获取一个图片
//@params	Id
//@return	*Picture
func GetOnePictureById(Id int) (*Picture, error) {
	picture := new(Picture)
	picture.Id = Id

	if err := ormer().Read(picture, "Id"); err != nil {
		YYLog.Error("GetOnePictureById : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return picture, nil
}

//GetPicturesNum 获取图片的数量
//@params	cateId keyword(搜索title的关键词)
//@return	int
func GetPicturesNum(catId int, keyword string) int {
	var data []orm.Params
	sql := "SELECT Id FROM picture WHERE CateId = ? AND Title LIKE ? "
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword).Values(&data)
	return len(data)
}

func GetPicturesNumByTagName(cateId int, keyTag string) {
	//sql := "select Id from picture, pic"
}

//GetPictures 获取图片
//@params	catId keyword(搜索title的关键词) pagesize offset
//@return	[]orm.Params
func GetPictures(catId int, keyword string, pagesize, offset int) []orm.Params {
	var data []orm.Params
	sql := `SELECT * FROM picture WHERE CateId = ? AND Title LIKE ? ORDER BY picture.Sort DESC,UpdateTime DESC LIMIT ?,?`
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword, offset, pagesize).Values(&data)
	return data
}

//CreateOnePicture 新建一篇图片
//@params	catId title thumb source description content
//@return	error
func CreateOnePicture(cateId, downloadCount int, title, publishTime, path, description string) (int, error) {
	picture := new(Picture)
	picture.CateId = cateId
	picture.Title = title
	picture.Path = path
	picture.Description = description
	picture.CreateTime = egb.TimeNowUnix()
	picture.UpdateTime = egb.TimeNowUnix()
	picture.PublishTime = publishTime
	picture.Download = int64(downloadCount)

	if index, err := ormer().Insert(picture); err != nil {
		YYLog.Error("CreateOnePicture : " + err.Error())
		return 0, fmt.Errorf(ErrInfo[SystemError])
	} else {
		YYLog.Debug(index)
		return int(index), nil
	}
}

//UpdatePicture 更新图片
//@params	id title description
//@return	error
func UpdatePicture(id, cateId int, title, description string) error {
	picture, err := GetOnePictureById(id)
	if err != nil {
		return err
	}

	picture.CateId = cateId
	picture.Title = title
	picture.Description = description
	picture.CreateTime = egb.TimeNowUnix()
	picture.UpdateTime = egb.TimeNowUnix()

	if _, err := ormer().Update(picture); err != nil {
		YYLog.Error("UpdatePicture : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		return nil
	}
}

//UpdatePictureSort 更新图片排序
//@params	id sort
//@return	error
func UpdatePictureSort(id, sort int) error {
	picture, err := GetOnePictureById(id)
	if err != nil {
		return fmt.Errorf(ErrInfo[SystemError])
	}
	picture.Sort = sort
	if _, err := ormer().Update(picture, "Sort"); err != nil {
		YYLog.Error("UpdatePictureSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOnePicture 删除一个图片
//@params	id
//@return	error
func DeleteOnePicture(id int) error {
	picture, err := GetOnePictureById(id)
	if err != nil {
		return err
	}
	if _, err := ormer().Delete(picture); err != nil {
		YYLog.Error("DeleteOnePicture : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}
