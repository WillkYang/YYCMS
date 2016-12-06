package controllers

import (
	m "YYCMS/models"
	"github.com/agelinazf/egb"
	cnf "YYCMS/conf"
)

type AdminRoleController struct {
	LoginController
}

//func (c *AdminRoleController) Prepare() {
//	beego.Debug("找到控制器")
//	c.LoginController.Prepare()
//
//	//TODO - 添加用户权限验证
//	beego.Debug("完成父控制器初始化")
//}

//List 获取角色列表
//@params   keyword(搜索title的关键词)
//@return   []AdminRole
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
	data := m.GetAdminRoles(cateId, keyword, pagesize, (page - 1) * pagesize)

	c.Msg["keyword"] = keyword
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}


//Profile 查看角色详情
//@params	id
//@return	AdminRole
func (c *AdminRoleController) Profile () {
	id := c.MustInt("id")
	if data,err := m.GetOneAdminRoleById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}

}

//AddAdminRole 添加角色
//@params	cateId modelId title description permission
//@return	success/error
func (c *AdminRoleController) Add() {
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	description := c.Str("description")
	permission := c.MustStr("permission")

	if err := m.CreateOneAdminRole(cateId, title, description, permission); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//UpdateAdminRole 更新角色
//@params	id cateId modelId title description permission
//@return	success/error
func (c *AdminRoleController) Update() {
	id := c.MustInt("id")
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	description := c.Str("description")
	permission := c.MustStr("permission")

	if err := m.UpdateAdminRole(id, cateId, title, description, permission); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//DelAdminRole 删除角色
//@params	id
//@return	error
func (c *AdminRoleController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneAdminRole(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
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
