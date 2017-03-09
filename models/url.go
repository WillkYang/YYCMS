package models

type Url struct {
	//主键
	Id int `orm:"column(Id);pk;auto"`
	//名字
	Title string `orm:"column(Title)"`
	//描述
	Description string `orm:"column(Description);null"`
	//创建时间
	CreateTime string `orm:"column(CreateTime)"`
	//更新时间
	UpdateTime string `orm:"column(UpdateTime)"`
	//排序值
	Sort int `orm:"column(Sort);default(0)"`
	//状态 - 当值为-1时表示数据已删除
	Status int `orm:"column(Status);default(0)"`

	//链接地址
	Url string `orm:"column(Url)"`
}
