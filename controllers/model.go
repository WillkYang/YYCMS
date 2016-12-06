package controllers

import (
	m "YYCMS/models"
	"github.com/agelinazf/egb"
)

type ModelController struct {
	LoginController
}

//获取全部模型信息
func (c *ModelController) List() {
	models := m.GetAllModelWithCache()
	c.Msg["keyword"] = ""
	c.Msg["count"] = len(models)
	c.Msg["lists"] = models
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//Profile 查看模型详情
//@params	id
//@return 	model
func (c *ModelController) Profile() {
	id := c.MustInt("id")
	if data, err := m.GetOneModelByIdWithCache(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}
}

//Add 添加模型动作
//@params   title tableName description
//@return   success/error
func (c *ModelController) Add() {
	cateId := c.MustInt("cateId")
	modeName := c.MustStr("title")
	description := c.Str("description")
	tableName := c.MustStr("tableName")

	if err := m.CreateOneModel(cateId, modeName, description, tableName); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "添加模型成功")
}

//Update 编辑模型动作
//@params	id name
//@return	success/error
func (c *ModelController) Update() {
	id := c.MustInt("id")
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	description := c.Str("description")
	tableName := c.MustStr("tableName")

	if err := m.UpdateModel(id, cateId, title, description, tableName); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "编辑模型成功")
}

//Delete 删除模型
//@params	id
//@return	success/error
func (c *ModelController) Delete() {
	id := c.MustInt("id")

	if err := m.DeleteModel(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//Sort 模型排序
//params	int : int 模型Id : 新sort值
//@return	success/error
func (c *ModelController) Sort() {
	postdata := c.Ctx.Request.PostForm
	for k, v := range postdata {
		modelId := egb.StringToInt(k)

		_, err := m.GetOneModelByIdWithCache(modelId)
		if err != nil {
			c.AjaxMsg(nil, m.InfoNotExistError, m.ErrInfo[m.InfoNotExistError], "")
			return
		}

		if err := m.UpdateModelSort(modelId, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}

