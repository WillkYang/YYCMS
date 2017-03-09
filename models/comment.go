package models

import (
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//CommentType 留言类型
type CommentType int8

const (
	//Question 留言的问
	Question CommentType = 0
	//Answer 留言的回复
	Answer CommentType = 1
)

//Comment 评论表模型
type Comment struct {
	//评论标示
	Id int `orm:"column(Id);pk;auto"`
	//所属项Id
	RelativeId int `orm:"column(RelativeId)"`
	//用户 - 不一定是纯数字
	UserId string `orm:"column(UserId)"`

	//用户名
	UserName string `orm:"column(UserName)"`
	//头像
	Avatar string `orm:"column(Avatar)"`
	//类型  问Question 或者  答Answer
	Type CommentType `orm:"column(Type);"`
	//内容
	Content string `orm:"column(Content);type(text)"`
	//时间
	CreateTime string `orm:"column(CreateTime);size(20)"`
	//发布时间
	PublishTime string `orm:"column(PublishTime);size(20)"`
	//如果是回复，则此处指的是问题的id
	ParentId int `orm:"column(ParentId)"`
	//IsPublish
	IsPublish bool `orm:"column(IsPublish)"`
}

//TableName 表名
func (model *Comment) TableName() string {
	return "comment"
}

//Query 查找语句
func (model *Comment) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(model)
}

//Insert 插入
func (model *Comment) Insert() error {
	if _, err := orm.NewOrm().Insert(model); err != nil {
		return err
	}
	return nil
}

//Read 读取
func (model *Comment) Read(id ...string) error {
	if err := orm.NewOrm().Read(model, id...); err != nil {
		return err
	}
	return nil
}

//Update 更新
func (model *Comment) Update(field ...string) error {
	if _, err := orm.NewOrm().Update(model, field...); err != nil {
		return err
	}
	return nil
}

//Delete 删除
func (model *Comment) Delete() error {
	if _, err := orm.NewOrm().Delete(model); err != nil {
		return err
	}
	return nil
}

//GetOneCommentById 通过id获取评论
//@params	id
//@return	*Comment
func GetOneCommentById(id int) *Comment {
	comment := new(Comment)
	comment.Id = id
	if err := comment.Read("Id"); err != nil {
		return new(Comment)
	}
	return comment
}

//AddOneComment 添加一个评论
//@params	relativeId Type content pId
//@return	error
func AddOneComment(relativeId int, cType CommentType, userName, avatar, userId, content string, pId int) error {
	comment := new(Comment)
	comment.RelativeId = relativeId
	comment.Type = cType
	comment.UserName = userName
	comment.UserId = userId
	comment.Avatar = avatar
	if cType == Answer {
		if pId == 0 {
			return fmt.Errorf(ErrInfo[SystemError])
		}
		comment.ParentId = pId
	}
	comment.Content = content
	comment.CreateTime = egb.TimeNowUnix()
	comment.PublishTime = comment.CreateTime
	if err := comment.Insert(); err != nil {
		beego.Error("AddOneComment : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateOneComment 编辑一个评论
//@params	relativeId Type content pId
//@return	error
func UpdateOneComment(id int, content string) error {
	comment := GetOneCommentById(id)
	comment.Content = content

	if _, err := ormer().Update(comment); err != nil {
		YYLog.Error("UpdateComment : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		return nil
	}
}

//PublishOneComment 发布评论
func PublishOneComment(id int) error {
	comment := GetOneCommentById(id)
	if comment.Id == 0 {
		return fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	comment.PublishTime = egb.TimeNowUnix()
	comment.IsPublish = true
	if _, err := orm.NewOrm().Update(comment, "IsPublish", "PublishTime"); err != nil {
		beego.Error("PublishOneComment :" + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneComment 删除评论
//@params	id
//@return	error
func DeleteOneComment(id int) error {
	comment := GetOneCommentById(id)
	if comment.Type == Answer {
		if err := comment.Delete(); err != nil {
			beego.Error("DeleteOneComment : " + err.Error())
			return fmt.Errorf(ErrInfo[SystemError])
		}
	} else {
		o := orm.NewOrm()
		o.Begin()
		//先删除问
		if _, err := o.QueryTable(&Comment{}).Filter("Id", id).Delete(); err != nil {
			beego.Error("DeleteOneComment : " + err.Error())
			o.Rollback()
			return fmt.Errorf(ErrInfo[SystemError])
		}
		//再删除这个问题下面的全部回答
		comments := GetAnswersByQuestionId(id)
		for _, v := range comments {
			if _, err := o.QueryTable(&Comment{}).Filter("Id", v.Id).Delete(); err != nil {
				beego.Error("DeleteOneComment : " + err.Error())
				o.Rollback()
				return fmt.Errorf(ErrInfo[SystemError])
			}
		}
		o.Commit()
		return nil
	}
	return nil
}

//GetCommentNum 获取评论问数量
//@params	relativeId
//@return	int
func GetCommentNum() int {
	num, _ := (&Comment{}).Query().Filter("Type", Question).Count()
	return int(num)
}

//GetAnswersByQuestionId 获取问题的回复
//@params	id
//@return	[]Comment
func GetAnswersByQuestionId(id int) []Comment {
	comments := make([]Comment, 0)
	(&Comment{}).Query().Filter("ParentId", id).Filter("Type", Answer).OrderBy("-CreateTime").All(&comments)
	return comments
}

//GetComments 获取评论
//@params	relativeId pagesize offset
//@return	[]orm.Params
func GetComments(keyword string, startTime, endTime string, pagesize, offset int) []orm.Params {
	comments := make([]orm.Params, 0)

	keyword = "%" + keyword + "%"
	if startTime == "" {
		startTime = "0"
	}
	if endTime == "" {
		endTime = "9999999999"
	}

	ormer().Raw(`SELECT comment.*, blog.Title FROM comment,blog WHERE Type = 0 AND RelativeId = blog.Id AND comment.Content LIKE ? AND comment.PublishTime > ? AND comment.PublishTime <= ? ORDER BY IsPublish ASC, CreateTime DESC LIMIT ?,?`, keyword, startTime, endTime, offset, pagesize).Values(&comments)
	return comments
}

