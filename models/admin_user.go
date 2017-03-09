package models

import (
	"YYCMS/conf"
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

//管理员表
type AdminUser struct {
	//主键
	Id int `orm:"column(Id);pk;auto"`
	//栏目Id
	//CateId          int    `orm:"column(CateId)"`
	//排序值
	Sort int `orm:"column(Sort);default(0)"`
	//账号
	Account string `orm:"column(Account)"`
	//密码
	Password string `orm:"column(Password)"`
	//最后登录IP
	LatestLoginIP string `orm:"column(LatestLoginIP)"`
	//Email
	Email string `orm:"column(Email)"`
	//真实姓名
	RealName string `orm:"column(RealName)"`
	//最后登录时间
	LastLoginTime string `orm:"column(LastLoginTime)"`
	//最新登录时间
	LatestLoginTime string `orm:"column(LatestLoginTime)"`
	//用户角色
	Role int `orm:"column(Role)"`
}

//todo 返回字段进行筛选

//GetOneAdminUserById 通过id获取一个管理员用户
//@params   id
//@return   *AdminUser
func GetOneAdminUserById(id int) (*AdminUser, error) {
	admin := new(AdminUser)
	admin.Id = id

	if err := ormer().Read(admin, "Id"); err != nil {
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return admin, nil
}

//GetOneAdminUserByAccount 通过账号获取一个管理员用户
//@params   account
//@return   *AdminUser
func GetOneAdminUserByAccount(account string) *AdminUser {
	admin := new(AdminUser)
	admin.Account = account
	if err := ormer().Read(admin, "Account"); err != nil {
		return new(AdminUser)
	}
	return admin
}

//AdminUserLogin 管理员登录
//@params   account password ip
//@return   error
func AdminUserLogin(account, password, ip string) (*AdminUser, error) {
	user := GetOneAdminUserByAccount(account)
	if user.Id == 0 {
		return nil, fmt.Errorf(ErrInfo[LoginError])
	}
	if user.Password != egb.StringSHA256Hex(password) {
		return nil, fmt.Errorf(ErrInfo[LoginError])
	}
	user.LastLoginTime = user.LatestLoginTime
	user.LatestLoginIP = ip
	user.LatestLoginTime = egb.TimeNowUnix()

	if _, err := ormer().Update(user, "LastLoginTime", "LatestLoginTime", "LatestLoginIP"); err != nil {
		YYLog.Error("AdminUserLogin : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[SystemError])
	}
	//user.Password = ""
	return user, nil
}

//CreateOneAdminUser 新建一个管理员用户
//@params	account pwd role
//@return	error
func CreateOneAdminUser(account, pwd string, role int) error {
	user := GetOneAdminUserByAccount(account)
	if user.Id != 0 {
		return fmt.Errorf(ErrInfo[AccountExistError])
	}
	admin := new(AdminUser)
	admin.Account = account
	admin.Password = egb.StringSHA256Hex(pwd)
	admin.Role = role

	if _, err := ormer().Insert(admin); err != nil {
		YYLog.Error("CreateOneAdminUser : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//GetAdminUserNumByRole 通过类型获取用户的数量
//@params	role
//@return	int64
func GetAdminUserNumByRole(role int) int64 {
	var num int64
	if role == 0 {
		num, _ = ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Count()
	} else {
		num, _ = ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Filter("Role", role).Count()
	}
	return num
}

//GetAdminuserByRole 通过类型获取用户
//@params	role pagesize offset
//@return	[]Adminuser
func GetAdminUserByRole(roleId, pagesize, offset int) []AdminUser {
	var users []AdminUser
	if roleId == 0 {
		ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Limit(pagesize, offset).All(&users,"Id", "Sort", "Account", "LatestLoginIP", "Role", "LastLoginTime")
	} else {
		ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Limit(pagesize, offset).Filter("Role", roleId).All(&users, "Id", "Sort", "Account", "LatestLoginIP", "Role", "LastLoginTime")
	}
	for i := 0; i < len(users); i++ {
		users[i].Password = ""
	}
	return users
}

//GetAdminUserRawByRole 通过类型获取用户
func GetAdminUserRawByRole(roleId, pagesize, offset int) []orm.Params {
	//var users []AdminUser
	var result []orm.Params
	if roleId == 0 {
		//ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Limit(pagesize, offset).All(&users,"Id", "Sort", "Account", "LatestLoginIP", "Role")
		ormer().Raw("select * FROM admin_user WHERE Account != \"yeeyun_root\" LIMIT ?, ?", offset, pagesize).Values(&result, "Id", "Sort", "Account", "LatestLoginIP", "Role")
	} else {
		//ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Limit(pagesize, offset).Filter("Role", roleId).All(&users, "Id", "Sort", "Account", "LatestLoginIP", "Role")
		ormer().Raw("select * FROM admin_user WHERE Role = ? AND Account != \"yeeyun_root\" LIMIT ?, ?", roleId, offset, pagesize).Values(&result, "Id", "Sort", "Account", "LatestLoginIP", "Role")
	}
	return result
}

//UpdateAdminUserInfo 更新用户信息
//@params	userid name oldpwd newpwd
//@return	error
func UpdateAdminUserInfo(userid, role int, name, oldpwd, newpwd string) error {
	idUser, err := GetOneAdminUserById(userid)
	if err != nil {
		return err
	}
	nameUser := GetOneAdminUserByAccount(name)
	if nameUser.Id != 0 && nameUser.Id != idUser.Id {
		return fmt.Errorf(ErrInfo[AccountExistError])
	}
	idUser.Account = name
	idUser.Role = role
	if oldpwd != "" {
		if oldpwd == newpwd {
			return fmt.Errorf(ErrInfo[PwdRepeatError])
		}
		//需要修改密码
		if !egb.RegexpIsStrongPassword(newpwd) {
			return fmt.Errorf(ErrInfo[PwdWeakError])
		}
		if idUser.Password != egb.StringSHA256Hex(oldpwd) {
			return fmt.Errorf(ErrInfo[OldPwdError])
		}
		idUser.Password = egb.StringSHA256Hex(newpwd)
	}

	if _, err := ormer().Update(idUser); err != nil {
		YYLog.Error("UpdateAdminUserInfo : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneAdminUser 删除管理员
//@params	id
//@return	error
func DeleteOneAdminUser(id int) error {
	user, err := GetOneAdminUserById(id)
	if err != nil {
		return err
	}

	if _, err := ormer().Delete(user); err != nil {
		YYLog.Error("DeleteOneAdminUser : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//ResetAdminUserPwd 重置管理员密码
//@params	id
//@return	error
func ResetAdminUserPwd(id int) error {
	user, err := GetOneAdminUserById(id)
	if err != nil {
		return err
	}
	user.Password = egb.StringSHA256Hex(egb.StringSHA1Hex(conf.DefaultAdminPassword))
	if _, err := ormer().Update(user, "Password"); err != nil {
		YYLog.Error("ResetAdminUserPwd : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateAdminUserPwd 更新管理员密码
//@params	id
//@return	error
func UpdateAdminUserPwd(id int, oldPwd, newPwd string) (*AdminUser, error) {
	user, err := GetOneAdminUserById(id)
	if err != nil {
		return nil, err
	}
	if egb.StringSHA256Hex(oldPwd) == user.Password {
		user.Password = egb.StringSHA256Hex(newPwd)
		if _, err := ormer().Update(user, "Password"); err != nil {
			YYLog.Error("UpdateAdminUserPwd : " + err.Error())
			return user, fmt.Errorf(ErrInfo[SystemError])
		}
		return user, nil
	} else {
		return user, fmt.Errorf(ErrInfo[LoginPwdError])
	}
}

//UpdateAdminUserSort 更新用户排序
//@params	id sort
//@return	error
func UpdateAdminUserSort(id, sort int) error {

	adminuser, err := GetOneAdminUserById(id)
	if err != nil {
		return err
	}
	adminuser.Sort = sort

	if _, err := ormer().Update(adminuser, "Sort"); err != nil {
		YYLog.Error("UpdateAdminUserSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//roleIds:1,2,3
//根据Ids字符串获取管理员用户
func GetAdminUserByRoleIds(roleIds string) (result []orm.Params) {
	if len(roleIds) == 0 {
		return
	}
	ormer().Raw(fmt.Sprintf("SELECT Id,Account,Role FROM admin_user WHERE Role IN (%s)",roleIds)).Values(&result)
	YYLog.Warning(result)
	return
}