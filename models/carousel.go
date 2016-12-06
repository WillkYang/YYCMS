package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/agelinazf/egb"
)

//轮换图表
type Carousel struct {
	//主键
	Id          int `orm:"column(Id);pk"`
	//栏目Id
	CateId      int    `orm:"column(CateId)"`
	//名字
	Title       string `orm:"column(Title)"`
	//描述
	Description string `orm:"column(Description);null"`
	//创建时间
	CreateTime  string `orm:"column(CreateTime)"`
	//更新时间
	UpdateTime  string `orm:"column(UpdateTime)"`
	//排序值
	Sort        int `orm:"column(Sort);default(0)"`

	//链接地址
	Url         string `orm:"column(Url)"`
	//图片
	Path 		string 	`orm:"column(Path)"`
}

//GetOneCarouselById 获取一个轮换图
//@params	Id
//@return	*Carousel
func GetOneCarouselById(Id int) (*Carousel, error) {
	carousel := new(Carousel)
	carousel.Id = Id

	if err := ormer().Read(carousel,"Id"); err != nil {
		beego.Error("GetOneCarouselById : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return carousel,nil
}

//GetCarouselsNum 获取轮换图的数量
//@params	cateId keyword(搜索title的关键词)
//@return	int
func GetCarouselsNum(catId int, keyword string) int {
	var data []orm.Params
	//todo count
	sql := "SELECT Id FROM carousel WHERE CateId = ? AND Title LIKE ? "
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword).Values(&data)
	return len(data)
}

//GetCarousels 获取轮换图
//@params	catId keyword(搜索title的关键词) pagesize offset
//@return	[]orm.Params
func GetCarousels(catId int, keyword string, pagesize, offset int) []orm.Params {

	var data []orm.Params
	sql := `SELECT * FROM carousel WHERE CateId = ? AND Title LIKE ? ORDER BY carousel.Sort DESC,UpdateTime DESC LIMIT ?,?`
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword, offset, pagesize).Values(&data)
	return data
}

//CreateOneCarousel 新建一篇轮换图
//@params	catId title thumb source description content
//@return	error
func CreateOneCarousel(cateId int, title, url, path, description string) error {
	carousel := new(Carousel)
	carousel.CateId = cateId
	carousel.Title = title
	carousel.Path = path
	carousel.Description = description
	carousel.Url = url
	carousel.CreateTime = egb.TimeNowUnix()
	carousel.UpdateTime = egb.TimeNowUnix()

	if _,err := ormer().Insert(carousel); err != nil {
		beego.Error("CreateOneCarousel : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateCarousel 更新轮换图
//@params	id title description
//@return	error
func UpdateCarousel(id, cateId int, title, url, path, description string) error {
	carousel,err := GetOneCarouselById(id)
	if err != nil {
		return err
	}

	carousel.CateId = cateId
	carousel.Title = title
	carousel.Path = path
	carousel.Url = url
	carousel.Description = description
	carousel.CreateTime = egb.TimeNowUnix()
	carousel.UpdateTime = egb.TimeNowUnix()

	if _,err := ormer().Update(carousel); err != nil {
		beego.Error("UpdateCarousel : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		return nil
	}
}

//UpdateCarouselSort 更新轮换图排序
//@params	id sort
//@return	error
func UpdateCarouselSort(id, sort int) error {
	carousel,err := GetOneCarouselById(id)
	if err != nil {
		return fmt.Errorf(ErrInfo[SystemError])
	}
	carousel.Sort = sort
	beego.Debug(carousel)
	if _,err := ormer().Update(carousel,"Sort"); err != nil {
		beego.Error("UpdateCarouselSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneCarousel 删除一个轮换图
//@params	id
//@return	error
func DeleteOneCarousel(id int) error {
	carousel,err := GetOneCarouselById(id)
	if err != nil {
		return err
	}
	if _,err := ormer().Delete(carousel); err != nil {
		beego.Error("DeleteOneCarousel : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}
