package controllers

import (
	m "YYCMS/models"
	"strings"

	"github.com/agelinazf/egb"
)

type BlogController struct {
	LoginController
}

//获取某个栏目的博客列表
func (c *BlogController) List() {
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

	count := m.GetBlogsNum(cateId, keyword, startTime, endTime)
	data := m.GetBlogs(cateId, keyword, startTime, endTime, pagesize, (page-1)*pagesize)
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//查看博客详情
func (c *BlogController) Profile() {
	id := c.MustInt("id")
	data := m.GetOneBlogById(id)
	c.Msg["data"] = data
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//添加博客
func (c *BlogController) Add() {
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	class := c.Str("class")
	author := c.Str("author")
	thumb := c.Str("thumb")
	bigImage := c.Str("bigImage")
	publishTime := c.Str("publishTime")
	if len(publishTime) == 0 {
		publishTime = egb.TimeNowUnix()
	}
	//hit := c.Int("hit")
	description := c.Str("description")
	content := c.MustStr("content")
	if err := m.CreateOneBlog(cateId, title, class, thumb, bigImage, author, description, content, publishTime, false); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//发布博客
func (c *BlogController) Publish() {
	newsId := c.MustInt("id")
	if err := m.PublishOneBlog(newsId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "发布成功")
}

//更新博客
func (c *BlogController) Update() {
	id := c.MustInt("id")
	title := c.Str("title")
	class := c.Str("class")
	author := c.Str("author")
	thumb := c.Str("thumb")
	bigImage := c.Str("bigImage")
	publishTime := c.Str("publishTime")
	//hit := c.Int("hit")
	description := c.Str("description")
	content := c.MustStr("content")
	if err := m.UpdateBlog(id,  title, class, thumb, bigImage, author, description, content, publishTime); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//删除博客
func (c *BlogController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneBlog(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//对博客进行排序
func (c *BlogController) Sort() {
	postdata := c.Ctx.Request.PostForm
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateBlogSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}
