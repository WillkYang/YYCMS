package controllers

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"YYCMS/utils"
	"YYCMS/utils/YYLog"
	"strings"

	"github.com/agelinazf/egb"
	"Mogu/conf"
	"YYCMS/helper"
)

type CategoryController struct {
	LoginController
}

//添加栏目
func (c *CategoryController) Add() {
	parentId := c.MustInt("parentId")
	cateName := c.MustStr("title")
	modelId := c.MustInt("modelId")
	privileges := c.MustStr("privileges")

	privileges = sortPrivileges(privileges)

	if cateId, err := m.CreateOneCategory(modelId, cateName, parentId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	} else {
		if err := m.CreatePrivileges(cateId, privileges); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		} else {
			m.UpdateCache()
			m.CreateOneLog(c.User.Account, "添加栏目："+cateName)
			c.AjaxMsg(nil, m.NoError, "", "添加栏目成功")
		}
	}
}

//编辑栏目
func (c *CategoryController) Update() {
	id := c.MustInt("id")
	title := c.MustStr("title")
	privileges := c.MustStr("privileges")
	privileges = sortPrivileges(privileges)
	YYLog.Warning(privileges)
	oldCategory, err := m.GetOneCategoryById(id)
	if err != nil {
		YYLog.Error(err)
		c.AjaxMsg(nil, m.DataBaseGetError, m.ErrInfo[m.DataBaseGetError], "")
		return
	}

	//更新栏目信息
	if err := m.UpdateCategory(id, oldCategory.ModelId, oldCategory.ParentId, title); err != nil {
		m.UpdateCache()
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	//删除旧的权限
	if err := m.DeleteAllPrivileges(id); err != nil {
		m.UpdateCache()
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	if oldCategory.ModelId != conf.DirModelId {
		//添加新的权限
		if err := m.CreatePrivileges(id, privileges); err != nil {
			m.UpdateCache()
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}

	m.UpdateCache()
	m.CreateOneLog(c.User.Account, "编辑栏目")
	c.AjaxMsg(nil, m.NoError, "", "编辑栏目成功")
}

//栏目列表
func (c *CategoryController) List() {
	topCateId := cnf.TopCateId
	YYLog.Debug(c.Ctx.Input.Query("isRecursion"))
	isrecursion := c.Int("isRecursion")
	if c.Int("topCateId") > 0 {
		topCateId = c.Int("topCateId")
	}
	cates := m.GetSubCategorysByIdWithCache(topCateId, c.User.Role, isrecursion <= 0)
	c.Msg["keyword"] = ""
	c.Msg["count"] = len(cates)
	c.Msg["lists"] = cates

	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//栏目详情
func (c *CategoryController) Profile() {
	id := c.MustInt("id")
	if data, err := m.GetOneCategoryRawById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}
}

//删除栏目
func (c *CategoryController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneCategory(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		m.UpdateCache()
		return
	}
	if err := m.DeleteAllPrivileges(id); err != nil {
		m.UpdateCache()
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.UpdateCache()
	m.CreateOneLog(c.User.Account, "删除栏目")
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//栏目排序
func (c *CategoryController) Sort() {
	postdata := c.Ctx.Request.PostForm

	YYLog.Warning(postdata)
	for k, v := range postdata {
		catId := egb.StringToInt(k)
		_, err := m.GetOneCategoryById(catId)
		if err != nil {
			c.AjaxMsg(nil, m.InfoNotExistError, m.ErrInfo[m.InfoNotExistError], "")
			return
		}
		if err := m.UpdateCategorySort(catId, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}

//全部权限
func (c *CategoryController) AllActions() {
	actions := helper.SystemAllAction()
	c.AjaxMsg(actions, m.NoError, m.ErrInfo[m.NoError], "")
}

//整理权限
func sortPrivileges(privileges string) string {
	actions := GetAllActions()
	var privs []string
	for _, action := range actions {
		if strings.Contains(privileges, action) {
			privs = append(privs, action)
		}
	}
	return utils.StringArrayToString(privs)
}
