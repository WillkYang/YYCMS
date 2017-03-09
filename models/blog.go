package models

import (
	"fmt"
	"os"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//Blog 博客表模型
type Blog struct {
	//博客Id
	Id int `orm:"column(Id);pk;auto"`
	//博客栏目Id
	CateId int `orm:"column(CateId)"`
	//类别
	Class string `orm:"column(Class)"`
	//标题
	Title string `orm:"column(Title)"`
	//缩略图
	Thumb string `orm:"column(Thumb)"`
	//顶部大图
	BigImage string `orm:"column(BigImage)"`
	//来源
	Source string `orm:"column(Source)"`
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
func (model *Blog) TableName() string {
	return "blog"
}

//GetOneBlogById 获取一篇博客
func GetOneBlogById(id int) *Blog {
	blog := new(Blog)
	blog.Id = id
	if err := orm.NewOrm().Read(blog, "Id"); err != nil {
		beego.Error("GetOneBlogById : " + err.Error())
		return new(Blog)
	}
	return blog
}

//GetBlogsNum 获取博客的数量
func GetBlogsNum(cateId int, keyword string, startTime, endTime string) int {
	var data []orm.Params
	sql := "SELECT Id FROM blog WHERE CateId = ? AND Title LIKE ? AND UpdateTime > ? AND UpdateTime < ?"
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

//GetBlogs 获取博客
func GetBlogs(cateId int, keyword string, startTime, endTime string, pagesize, offset int) []orm.Params {
	var data []orm.Params
	var sql string
	sql = `SELECT Id,Title,Thumb,Author,Description,Source,UpdateTime,PublishTime,blog.Sort FROM blog WHERE CateId = ? AND Title LIKE ? AND PublishTime > ? AND PublishTime < ? ORDER BY blog.Sort DESC,UpdateTime DESC LIMIT ?,?`
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

//CreateOneBlog 新建一篇博客
func CreateOneBlog(cateId int, title, class, thumb, bigImage, author, description, content, publishTime string, isPublish bool) error {
	blog := new(Blog)
	blog.CateId = cateId
	blog.Title = title
	blog.Thumb = thumb
	blog.BigImage = bigImage
	blog.Class = class
	blog.Author = author
	blog.Description = description
	blog.Content = content
	blog.CreateTime = egb.TimeNowUnix()
	blog.UpdateTime = egb.TimeNowUnix()
	blog.PublishTime = publishTime
	blog.IsPublish = isPublish
	if _, err := orm.NewOrm().Insert(blog); err != nil {
		beego.Error("CreateOneBlog : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//PublishOneBlog 发布博客
func PublishOneBlog(id int) error {
	blog := GetOneBlogById(id)
	if blog.Id == 0 {
		return fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	blog.IsPublish = true
	blog.UpdateTime = egb.TimeNowUnix()
	if _, err := orm.NewOrm().Update(blog, "IsPublish", "UpdateTime"); err != nil {
		beego.Error("PublishOneBlog :" + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateBlog 更新博客
func UpdateBlog(id int, title, class, thumb, bigImage, author, description, content, publicTime string) error {
	blog := GetOneBlogById(id)
	blog.Title = title
	blog.Class = class
	//blog.Hit = int64(hit)
	blog.Author = author
	blog.Description = description
	blog.Content = content
	blog.UpdateTime = egb.TimeNowUnix()
	blog.PublishTime = publicTime

	if thumb != "" {
		blog.Thumb = thumb
	}

	if bigImage != "" {
		blog.BigImage = bigImage
	}

	if _, err := orm.NewOrm().Update(blog); err != nil {
		beego.Error("UpdateBlog : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateBlogSort 更新博客排序
func UpdateBlogSort(id, order int) error {
	blog := GetOneBlogById(id)
	blog.Sort = order
	if _, err := orm.NewOrm().Update(blog, "Sort"); err != nil {
		beego.Error("UpdateBlogSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneBlog 删除一篇博客
func DeleteOneBlog(id int) error {
	blog := GetOneBlogById(id)
	path := blog.Thumb
	if _, err := orm.NewOrm().Delete(blog); err != nil {
		beego.Error("DeleteOneBlog : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	os.Remove("." + path)
	return nil
}

//SearchBlogsNum 搜索博客数目
func SearchBlogsNum(keyword string) int {
	var data []orm.Params
	orm.NewOrm().Raw(`SELECT a.Id FROM blog AS a WHERE a.Title LIKE ? `, "%"+keyword+"%").Values(&data)
	return len(data)
}

//SearchBlogs 搜索博客
func SearchBlogs(keyword string, pagesize, offset int) []orm.Params {
	var data []orm.Params
	orm.NewOrm().Raw(`SELECT a.Id,a.CateId,a.Title,a.Thumb,a.Source,a.UpdateTime,a.Description,c.ParentId
    FROM blog AS a ,category AS c
    WHERE a.CateId = c.Id AND a.Title LIKE ? ORDER BY a.UpdateTime DESC LIMIT ?,?`, "%"+keyword+"%", offset, pagesize).Values(&data)
	return data
}

func AddBlogHitTime(cateId,id int) {
	ormer().QueryTable(&Blog{}).Filter("Id", id).Filter("CateId",cateId).Update(orm.Params{
		"Hit": orm.ColValue(orm.ColAdd, 1),
	})
}
