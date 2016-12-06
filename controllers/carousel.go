package controllers

import (
	"github.com/astaxie/beego"
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"github.com/agelinazf/egb"
)

type CarouselController struct {
	LoginController
}

//List 获取轮换图列表
//@params   keyword(搜索title的关键词)
//@return   []Carousel
func (c *CarouselController) List() {
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

	count := m.GetCarouselsNum(cateId, keyword)
	data := m.GetCarousels(cateId, keyword, pagesize, (page - 1) * pagesize)

	c.Msg["keyword"] = keyword
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}


//Profile 查看轮换图详情
//@params	id
//@return	Carousel
func (c *CarouselController) Profile () {
	id := c.MustInt("id")
	if data,err := m.GetOneCarouselById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}

}

//AddCarousel 添加轮换图
//@params	cateId modelId title description
//@return	success/error
func (c *CarouselController) Add() {
	cateId := c.MustInt("cateId")
	title := c.Str("title")
	description := c.Str("description")
	url :=  c.Str("url")
	path := c.MustStr("path")

	if err := m.CreateOneCarousel(cateId, title, url, path, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//UpdateCarousel 更新轮换图
//@params	id cateId modelId title description
//@return	success/error
func (c *CarouselController) Update() {
	id := c.MustInt("id")
	cateId := c.MustInt("cateId")
	title := c.Str("title")
	url :=  c.Str("url")
	description := c.Str("description")
	path := c.Str("path")

	if err := m.UpdateCarousel(id, cateId, title, url, path, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//DelCarousel 删除轮换图
//@params	id
//@return	error
func (c *CarouselController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneCarousel(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//SortCarousel 对轮换图进行排序
//@params	id sort
//@return	success/error
func (c *CarouselController) Sort() {
	postdata := c.Ctx.Request.PostForm
	beego.Debug(postdata)
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateCarouselSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}