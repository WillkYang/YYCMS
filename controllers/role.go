package controllers

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"YYCMS/utils"
	"math"

	"github.com/agelinazf/egb"
)

type AdminRoleController struct {
	LoginController
}

//角色列表
func (c *AdminRoleController) List() {
	page := c.Int("page")
	keyword := c.Str("keyword")
	cateId := c.MustInt("cateId")
	pagesize := c.Int("pagesize")

	if page == 0 {
		page = 1
	}

	if pagesize <= 0 {
		pagesize = cnf.DefaultPageSize
	}

	count := m.GetAdminRolesNum(cateId, keyword)
	datas := m.GetAdminRoles(cateId, keyword, pagesize, (page-1)*pagesize)
	for _, data := range datas {
		roleId := egb.StringToInt(data["Id"].(string))
		//cates := m.GetSubCategorysByIdWithCache(cnf.TopCateId, roleId, true)
		data["RoleCates"] = m.ReadOneRoleCates(roleId)
	}
	c.Msg["keyword"] = keyword
	c.Msg["count"] = count
	c.Msg["lists"] = datas
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

func (c *AdminRoleController) All() {
	cateId := c.MustInt("cateId")
	data := m.GetAdminRoles(cateId, "", math.MaxInt64, 0)
	c.Msg["lists"] = data
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//角色详情
func (c *AdminRoleController) Profile() {
	id := c.MustInt("id")
	if data, err := m.ReadOneRole(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}
}

//添加角色
func (c *AdminRoleController) Add() {
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	privIds := c.MustStr("privIds")

	//1.添加角色
	roleId, err := m.CreateOneAdminRole(cateId, title)
	if err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	//删除旧权限
	if err := m.DeleteAllRolePrivileges(roleId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	//添加新的权限
	for _, privId := range utils.StringToIntArray(privIds) {
		if err := m.CreateOneRolePrivilege(roleId, privId); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	m.UpdateCache()
	m.CreateOneLog(c.User.Account, "添加角色")
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//UpdateAdminRole 更新角色
//@params	id cateId modelId title description permission
//@return	success/error
func (c *AdminRoleController) Update() {
	roleId := c.MustInt("id")
	title := c.MustStr("title")
	privIds := c.Str("privIds")

	//更新角色
	if err := m.UpdateAdminRole(roleId, title); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	if roleId != cnf.SuperAdminRoleId {
		//删除旧权限
		if err := m.DeleteAllRolePrivileges(roleId); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}

		//添加新的权限
		for _, privId := range utils.StringToIntArray(privIds) {
			if err := m.CreateOneRolePrivilege(roleId, privId); err != nil {
				c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
				return
			}
		}
	}
	m.UpdateCache()
	m.CreateOneLog(c.User.Account, "编辑角色")
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//DelAdminRole 删除角色
//@params	id
//@return	error
func (c *AdminRoleController) Delete() {
	roleId := c.MustInt("id")

	//超级管理员不允许删除
	if roleId == -1 {
		c.AjaxMsg(nil, m.OperateError, "操作失败，超级管理员不可删除", "")
		return
	}

	//判断是否存在该角色的管理员，若有则拒绝删除
	if m.GetAdminUserNumByRole(roleId) > 0 {
		c.AjaxMsg(nil, m.OperateError, "操作失败，请先删除此角色下的管理员", "")
		return
	}

	//删除角色
	if err := m.DeleteOneAdminRole(roleId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	//删除权限
	if err := m.DeleteAllRolePrivileges(roleId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.UpdateCache()
	m.CreateOneLog(c.User.Account, "删除角色")
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//SortAdminRole 对角色进行排序
//@params	id sort
//@return	success/error
func (c *AdminRoleController) Sort() {
	postdata := c.Ctx.Request.PostForm
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateAdminRoleSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}
