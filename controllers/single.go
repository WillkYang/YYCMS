package controllers

import (
	m "YYCMS/models"
)

//SinglePageController 单页面控制器
type SinglePageController struct {
	LoginController
}

//Profile 单个数据
//@params	cateId
//@return	data
func (c *SinglePageController) Profile() {
	cateId := c.MustInt("cateId")
	singlePage, _ := m.ReadOrCreateOneSinglePage(cateId)
	c.Msg["data"] = singlePage
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//Update 编辑
//@params	cateId title content
//@return	error
func (c *SinglePageController) Update() {
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	content := c.MustStr("content")
	if err := m.UpdateOneSinglePage(cateId, title, content); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}
