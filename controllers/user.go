package controllers

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"fmt"

	"github.com/agelinazf/egb"
)

type AdminUserController struct {
	LoginController
}

//List 获取用户列表
//@params	roleId pagesize page
//@return	[]adminuser
func (c *AdminUserController) List() {
	pagesize := c.Int("pagesize")
	page := c.Int("page")
	roleId := c.Int("roleId")

	if page == 0 {
		page = 1
	}

	if pagesize <= 0 {
		pagesize = cnf.DefaultPageSize
	}

	count := m.GetAdminUserNumByRole(roleId)
	data := m.GetAdminUserByRole(roleId, pagesize, (page-1)*pagesize)
	c.Msg["count"] = count
	c.Msg["page"] = page
	c.Msg["lists"] = data
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//Profile 获取单个用户信息
//@params	userId
//@return	adminuser
func (c *AdminUserController) Profile() {
	id := c.MustInt("id")
	if user, err := m.GetOneAdminUserById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		user.Password = ""
		c.Msg["user"] = user
		c.AjaxMsg(c.Msg, m.NoError, "", "")
	}
}

//Add 添加用户页面 / 添加用户动作
//params	account  pwd roleId
//return	success/error
func (c *AdminUserController) Add() {
	account := c.MustStr("account")
	//pwd := c.MustStr("pwd")
	pwd := egb.StringSHA1Hex(cnf.DefaultAdminPassword)

	roleId := c.MustInt("roleId")
	if err := m.CreateOneAdminUser(account, pwd, roleId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.CreateOneLog(c.User.Account, "添加用户："+account)
	c.AjaxMsg(nil, m.NoError, "", fmt.Sprintf("添加成功,密码为：%s,请登录后立即修改密码", cnf.DefaultAdminPassword))
}

//Update  更新用户信息操作
//params	userId name code pwd newpwd(if pwd != "" ,update pwd) roleId releaseRoleId
//return	success/error
func (c *AdminUserController) Update() {
	id := c.MustInt("id")
	name := c.MustStr("account")
	roleId := c.MustInt("roleId")
	oldpwd := c.Str("pwd")
	newpwd := c.Str("newpwd")
	if err := m.UpdateAdminUserInfo(id, roleId, name, oldpwd, newpwd); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.CreateOneLog(c.User.Account, "编辑用户："+name)
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//Del 删除管理员操作
//@params	userId
//@return	success/error
func (c *AdminUserController) Delete() {
	id := c.MustInt("id")
	if id == c.User.Id {
		c.AjaxMsg(nil, m.SystemError, "无法删除自身用户", "")
	}
	if err := m.DeleteOneAdminUser(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.CreateOneLog(c.User.Account, "删除用户")
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//ResetPwd 重置用户密码
//@params	userId
//@return	success/error
func (c *AdminUserController) ResetPwd() {
	userId := c.MustInt("id")
	if err := m.ResetAdminUserPwd(userId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.CreateOneLog(c.User.Account, "重置用户密码："+egb.StringFromInt(userId))
	c.AjaxMsg(nil, m.NoError, "", "重置密码成功,默认密码为:"+cnf.DefaultAdminPassword)
}

//UpdatePwd 修改密码
//@params	password
//@return	success/error
func (c *AdminUserController) UpdatePwd() {
	oldPwd := c.MustStr("oldPwd")
	newPwd := c.MustStr("newPwd")

	if oldPwd == newPwd {
		c.AjaxMsg(nil, m.PwdRepeatError, m.ErrInfo[m.PwdRepeatError], "")
		return
	}

	if !egb.RegexpIsStrongPassword(newPwd) {
		c.AjaxMsg(nil, m.PwdWeakError, m.ErrInfo[m.PwdWeakError], "")
		return
	}

	if user, err := m.UpdateAdminUserPwd(c.User.Id, oldPwd, newPwd); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	} else {
		c.User = *user
		c.UpdateCookie()
	}
	c.AjaxMsg(nil, m.NoError, "", "修改密码成功")
}

//SortAdminRole 对角色进行排序
//@params	id sort
//@return	success/error
func (c *AdminUserController) Sort() {
	postdata := c.Ctx.Request.PostForm
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateAdminUserSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}
