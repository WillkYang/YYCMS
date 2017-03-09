package controllers

import (
	"github.com/agelinazf/egb"
	"strings"
	m "YYCMS/models"
)

type NewsController struct {
	LoginController
}

//获取某个栏目的新闻列表
func (c *NewsController) List() {
	keyword := strings.TrimSpace(c.Str("keyword"))
	c.Msg["keyword"] = keyword
	startTime := c.Str("starttime")
	c.Msg["starttime"] = startTime
	endTime := c.Str("endtime")
	c.Msg["endtime"] = endTime
	cateId := c.MustInt("cateId")

	page := c.Int("page")
	pagesize := c.Int("pagesize")

	FormatPage(&page)
	FormatPageSize(&pagesize)

	count := m.GetNewssNum(cateId, keyword, startTime, endTime)
	data := m.GetNewss(cateId, keyword, startTime, endTime, pagesize, (page - 1) * pagesize)
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//查看新闻详情
func (c *NewsController) Profile() {
	id := c.MustInt("id")
	data := m.GetOneNewsById(id)
	c.Msg["data"] = data
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//添加新闻
func (c *NewsController) Add() {
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	source := c.Str("source")
	author := c.Str("author")
	thumb := c.Str("thumb")
	description := c.Str("description")
	content := c.MustStr("content")
	isPublish := c.Bool("isPublish")
	publishTime := c.Str("publishTime")

	if err := m.CreateOneNews(cateId, title, source, thumb, author, description, content, publishTime, isPublish); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//发布新闻
func (c *NewsController) Publish() {
	newsId := c.MustInt("id")
	if err := m.PublishOneNews(newsId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "发布成功")
}

//更新新闻
func (c *NewsController) Update() {
	id := c.MustInt("id")
	title := c.Str("title")
	source := c.Str("source")
	author := c.Str("author")
	thumb := c.Str("thumb")
	description := c.Str("description")
	content := c.MustStr("content")
	publishTime := c.MustStr("publishTime")
	if err := m.UpdateNews(id, title, source, thumb, author, description, content, publishTime); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//删除新闻
func (c *NewsController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneNews(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//对新闻进行排序
func (c *NewsController) Sort() {
	postdata := c.Ctx.Request.PostForm
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateNewsSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}