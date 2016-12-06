package models

import (
	"github.com/astaxie/beego/orm"
	"os"
	cnf "YYCMS/conf"
)

// 注 status为99表示系统保留，非删除却不可见


//模型表
type Model struct {
	//主键
	Id          int `orm:"column(Id);pk"`
	//名字
	Title       string `orm:"column(Title)"`
	//描述
	Description string `orm:"column(Description);null"`
	//排序值
	Sort        int `orm:"column(Sort);default(0)"`
	//表名
	DBTableName string `orm:"column(DBTableName)"`
}






//管理员表
type AdminUser struct {
	//主键
	Id              int `orm:"column(Id);pk"`
	//栏目Id
	CateId          int    `orm:"column(CateId)"`
	//排序值
	Sort            int `orm:"column(Sort);default(0)"`

	//账号
	Account         string `orm:"column(Account)"`
	//密码
	Password        string `orm:"column(Password)"`
	//最后登陆IP
	LatestLoginIP   string `orm:"column(LatestLoginIP)"`
	//Email
	Email           string `orm:"column(Email)"`
	//真实姓名
	RealName        string    `orm:"column(RealName)"`
	//最后登陆时间
	LastLoginTime   string `orm:"column(LastTime)"`
	//最新登陆时间
	LatestLoginTime string `orm:"column(LatestTime)"`
	//用户角色
	Role            int        `orm:"column(Role)"`
}

//角色表
type AdminRole struct {
	//主键
	Id                  int `orm:"column(Id);pk"`
	//栏目Id
	CateId              int    `orm:"column(CateId)"`
	//名字
	Title               string    `orm:"column(Title)"`
	//描述
	Description         string `orm:"column(Description); null"`
	//排序值
	Sort                int `orm:"column(Sort);default(0)"`
	//操作权限
	OperationPermission string    `orm:"column(OperationPermission)"`
}

////栏目权限表 - 暂不使用
//type CategoryPrivilege struct {
//	//栏目Id
//	CateId      int    `orm:"column(CateId)"`
//	//角色Id
//	RoleId      int    `orm:"column(RoleId)"`
//	//方法
//	Action		string 	`orm:"column(Action)"`
//}

//栏目表
type Category struct {
	//主键
	Id          int `orm:"column(Id);pk"`
	//名字
	Title       string `orm:"column(Title)"`
	//描述
	Description string `orm:"column(Description);null"`
	//排序值
	Sort        int `orm:"column(Sort);default(0)"`
	//模型Id
	ModelId     int `orm:"column(ModelId)"`
	//状态 - 当值为-1时表示数据已删除
	Status      int    `orm:"column(Status);default(0)"`

	//父栏目Id
	ParentId    int    `orm:"column(ParentId)"`
}

//单页面表
type Single struct {
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
	//状态 - 当值为-1时表示数据已删除
	Status      int    `orm:"column(Status);default(0)"`

	//内容
	Content     string    `orm:"column(Content);type(BLOB)"`
}

//文章表
type Article struct {
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
	Sort        int `orm:"column(Sort)"`
	//状态 - 当值为-1时表示数据已删除
	Status      int    `orm:"column(Status);default(0)"`

	//内容
	Content     string    `orm:"column(Content);type(BLOB)"`
}

//视频表
type Video struct {
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
	//状态 - 当值为-1时表示数据已删除
	Status      int    `orm:"column(Status);default(0)"`

	//内容
	Content     string    `orm:"column(Content);type(BLOB)"`
}

//系统信息表
type System struct {
	//主键
	Id          int `orm:"column(Id);pk"`

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
	//状态 - 当值为-1时表示数据已删除
	Status      int    `orm:"column(Status);default(0)"`
}



type Url struct {
	//主键
	Id          int `orm:"column(Id);pk"`
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
	//状态 - 当值为-1时表示数据已删除
	Status      int    `orm:"column(Status);default(0)"`

	//链接地址
	Url         string `orm:"column(Url)"`
}

func ormer() orm.Ormer {
	return orm.NewOrm()
}

//UpdateCache 更新缓存
//@params	nil
//@return	nil
func UpdateCache() {
	os.RemoveAll(cnf.ModelCachePath)
	os.RemoveAll(cnf.CategoryCachePath)
}