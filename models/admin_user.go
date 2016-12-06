package models

import (
	"fmt"
	"github.com/agelinazf/egb"
	"github.com/astaxie/beego"
)

//todo 返回字段进行筛选

//GetOneAdminUserById 通过id获取一个管理员用户
//@params   id
//@return   *AdminUser
func GetOneAdminUserById(id int) (*AdminUser,error) {
	admin := new(AdminUser)
	admin.Id = id

	if err := ormer().Read(admin,"Id"); err != nil {
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return admin,nil
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
//@params   account password
//@return   error
func AdminUserLogin(account, password, ip string) (*AdminUser, error) {
	user := GetOneAdminUserByAccount(account)
	if user.Id == 0 {
		return nil, fmt.Errorf(ErrInfo[LoginError])
	}
	beego.Warn(egb.StringSHA256Hex(password))
	if user.Password != egb.StringSHA256Hex(password) {
		return nil, fmt.Errorf(ErrInfo[LoginError])
	}
	user.LastLoginTime = user.LatestLoginTime
	user.LatestLoginIP = ip
	user.LatestLoginTime = egb.TimeNowUnix()

	if _,err := ormer().Update(user,"LastLoginTime","LatestLoginTime","LatestLoginIP"); err != nil {
		beego.Error("AdminUserLogin : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[SystemError])
	}
	return user, nil
}

//CreateOneAdminUser 新建一个管理员用户
//@params	account pwd role permission
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

	if _,err := ormer().Insert(admin); err != nil {
		beego.Error("CreateOneAdminUser : " + err.Error())
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
		num, _ = ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Filter("Type", role).Count()
	}
	return num
}

//GetAdminuserByRole 通过类型获取用户
//@params	role
//@return	[]Adminuser
func GetAdminUserByRole(roleId, pagesize, offset int) []AdminUser {
	var users []AdminUser
	if roleId == 0 {
		ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Limit(pagesize, offset).All(&users)
	} else {
		ormer().QueryTable(&AdminUser{}).Exclude("Account", "yeeyun_root").Limit(pagesize, offset).Filter("RoleId", roleId).All(&users)
	}
	for i := 0; i < len(users); i++ {
		users[i].Password = ""
	}
	return users
}

//UpdateAdminUserInfo 更新用户信息
//@params	userid name oldpwd newpwd
//@return	error
func UpdateAdminUserInfo(userid, role int, name, oldpwd, newpwd string) error {
	idUser,err := GetOneAdminUserById(userid)
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

	if _,err := ormer().Update(idUser); err != nil {
		beego.Error("UpdateAdminUserInfo : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneAdminUser 删除管理员
//@params	id
//@return	error
func DeleteOneAdminUser(id int) error {
	user,err := GetOneAdminUserById(id)
	if err != nil {
		return err
	}

	if _,err := ormer().Delete(user); err != nil {
		beego.Error("DeleteOneAdminUser : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//ResetAdminUserPwd 重置管理员密码
//@params	id
//@return	error
func ResetAdminUserPwd(id int) error {
	user,err := GetOneAdminUserById(id)
	if err != nil {
		return err
	}
	user.Password = egb.StringSHA256Hex("zaq1xsw2")

	if _,err := ormer().Update(user,"Password"); err != nil {
		beego.Error("ResetAdminUserPwd : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateAdminUserSort 更新用户排序
//@params	id sort
//@return	error
func UpdateAdminUserSort(id, sort int) error {

	adminuser,err := GetOneAdminUserById(id)
	if err != nil {
		return err
	}
	adminuser.Sort = sort

	if _,err := ormer().Update(adminuser,"Sort"); err != nil {
		beego.Error("UpdateAdminUserSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}
