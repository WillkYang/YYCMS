package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"fmt"
	"os"
	"github.com/agelinazf/egb"
)

//News 新闻表模型
type News struct {
	//新闻Id
	Id int `orm:"column(Id);pk;auto"`
	//新闻所属的栏目Id
	CateId int `orm:"column(CateId)"`
	//标题
	Title string `orm:"column(Title)"`
	//缩略图
	Thumb string `orm:"column(Thumb)"`
	//来源
	Source string  `orm:"column(Source)"`
	//作者
	Author string `orm:"column(Author)"`
	//摘要
	Description string `orm:"column(Description);type(text)"`
	//内容
	Content string `orm:"column(Content);type(text)"`
	//是否发布
	IsPublish bool `orm:"column(IsPublish)"`
	//创建时间
	CreateTime string `orm:"column(CreateTime);size(20)"`
	//更新时间
	UpdateTime string `orm:"column(UpdateTime);size(20)"`
	//发布时间
	PublishTime string `orm:"column(PublishTime);size(20)"`
	//排序
	Sort int `orm:"column(Sort);default(0)"`
	//点击量
	Hit int64 `orm:"column(Hit);default(0)"`
}

//TableName 表名
func (model *News) TableName() string {
	return "news"
}

//GetOneNewsById 获取一篇新闻
func GetOneNewsById(id int) *News {
	news := new(News)
	news.Id = id
	if err := orm.NewOrm().Read(news, "Id"); err != nil {
		beego.Error("GetOneNewsById : " + err.Error())
		return new(News)
	}
	return news
}

//GetNewssNum 获取新闻的数量
func GetNewssNum(cateId int, keyword string, startTime, endTime string) int {
	var data []orm.Params
	sql := "SELECT Id FROM news WHERE CateId = ? AND Title LIKE ? AND UpdateTime > ? AND UpdateTime < ?"
	keyword = "%" + keyword + "%"
	if startTime == "" {
		startTime = "0"
	}
	if endTime == "" {
		endTime = "9999999999"
	}

	orm.NewOrm().Raw(sql, cateId, keyword, startTime, endTime).Values(&data)
	return len(data)
}

//GetNewss 获取新闻
func GetNewss(cateId int, keyword string, startTime, endTime string, pagesize, offset int) []orm.Params {
	var data []orm.Params
	var sql string
	sql = `SELECT Id,Title,Thumb,Author,Description,Source,UpdateTime,PublishTime,news.Sort FROM news WHERE CateId = ? AND Title LIKE ? AND PublishTime > ? AND PublishTime < ? ORDER BY news.Sort DESC,UpdateTime DESC LIMIT ?,?`
	keyword = "%" + keyword + "%"
	if startTime == "" {
		startTime = "0"
	}
	if endTime == "" {
		endTime = "9999999999"
	}
	orm.NewOrm().Raw(sql, cateId, keyword, startTime, endTime, offset, pagesize).Values(&data)
	return data
}

//CreateOneNews 新建一篇新闻
func CreateOneNews(cateId int, title, source, thumb, author, description, content, publishTime string, isPublish bool) error {
	news := new(News)
	news.CateId = cateId
	news.Title = title
	news.Thumb = thumb
	news.Source = source
	news.Author = author
	news.Description = description
	news.Content = content
	news.CreateTime = egb.TimeNowUnix()
	news.UpdateTime = egb.TimeNowUnix()
	news.PublishTime = publishTime
	news.IsPublish = isPublish
	if _, err := orm.NewOrm().Insert(news); err != nil {
		beego.Error("CreateOneNews : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//PublishOneNews 发布新闻
func PublishOneNews(id int) error {
	news := GetOneNewsById(id)
	if news.Id == 0 {
		return fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	news.IsPublish = true
	news.UpdateTime = egb.TimeNowUnix()
	if _, err := orm.NewOrm().Update(news, "IsPublish", "UpdateTime"); err != nil {
		beego.Error("PublishOneNews :" + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateNews 更新新闻
func UpdateNews(id int, title, source, thumb, author, description, content, publishTime string) error {
	news := GetOneNewsById(id)
	news.Title = title
	news.Source = source
	if thumb != "" {
		news.Thumb = thumb
	}
	if author != "" {
		news.Author = author
	}
	if description != "" {
		news.Description = description
	}
	news.Content = content
	news.UpdateTime = egb.TimeNowUnix()
	news.PublishTime = publishTime
	if _, err := orm.NewOrm().Update(news); err != nil {
		beego.Error("UpdateNews : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateNewsSort 更新新闻排序
func UpdateNewsSort(id, order int) error {
	news := GetOneNewsById(id)
	news.Sort = order
	if _, err := orm.NewOrm().Update(news, "Sort"); err != nil {
		beego.Error("UpdateNewsSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneNews 删除一篇新闻
func DeleteOneNews(id int) error {
	news := GetOneNewsById(id)
	path := news.Thumb
	if _, err := orm.NewOrm().Delete(news); err != nil {
		beego.Error("DeleteOneNews : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	os.Remove("." + path)
	return nil
}

//SearchNewssNum 搜索新闻数目
func SearchNewssNum(keyword string) int {
	var data []orm.Params
	orm.NewOrm().Raw(`SELECT a.Id FROM news AS a WHERE a.Title LIKE ? `, "%"+keyword+"%").Values(&data)
	return len(data)
}

//SearchNewss 搜索新闻
func SearchNewss(keyword string, pagesize, offset int) []orm.Params {
	var data []orm.Params
	orm.NewOrm().Raw(`SELECT a.Id,a.CateId,a.Title,a.Thumb,a.Source,a.UpdateTime,a.Description,c.ParentId
    FROM news AS a ,category AS c
    WHERE a.CateId = c.Id AND a.Title LIKE ? ORDER BY a.UpdateTime DESC LIMIT ?,?`, "%"+keyword+"%", offset, pagesize).Values(&data)
	return data
}
