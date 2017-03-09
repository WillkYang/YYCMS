package controllers

import (
	m "YYCMS/models"
	"YYCMS/utils/YYLog"
	"strings"
)

type CommentController struct {
	LoginController
}

//获取某个栏目的评论列表
func (c *CommentController) List() {
	keyword := strings.TrimSpace(c.Str("keyword"))
	c.Msg["keyword"] = keyword

	page := c.Int("page")
	pagesize := c.Int("pagesize")
	startTime := c.Str("starttime")
	c.Msg["starttime"] = startTime
	endTime := c.Str("endtime")
	c.Msg["endtime"] = endTime

	FormatPage(&page)
	FormatPageSize(&pagesize)
	count := m.GetCommentNum()
	data := m.GetComments(keyword, startTime, endTime, pagesize, (page-1)*pagesize)
	YYLog.Warning(data)
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//查看评论详情
func (c *CommentController) Profile() {
	id := c.MustInt("id")
	data := m.GetOneCommentById(id)
	c.Msg["data"] = data
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

//添加评论
func (c *CommentController) Add() {
	relId := c.MustInt("relId")
	userName := c.Str("username")
	userId := c.Str("userId")
	content := c.MustStr("content")
	avatar := c.Str("avatar")
	if err := m.AddOneComment(relId, m.Question, userName, avatar, userId, content, 0); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//发布评论
func (c *CommentController) Publish() {
	newsId := c.MustInt("id")
	if err := m.PublishOneComment(newsId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "发布成功")
}

//更新评论
func (c *CommentController) Update() {
	id := c.MustInt("id")
	content := c.MustStr("content")
	if err := m.UpdateOneComment(id, content); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//删除评论
func (c *CommentController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneComment(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//排序评论
func (c *CommentController) Sort() {
	c.AjaxMsg(nil, m.SystemError, "", "暂不提供")
}
