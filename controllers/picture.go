package controllers

import (
	"github.com/agelinazf/egb"
	cnf "YYCMS/conf"
	m "YYCMS/models"
)

type PictureController struct {
	LoginController
}

//List 获取图片列表
//@params   keyword(搜索title的关键词)
//@return   []Picture
func (c *PictureController) List() {

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

	count := m.GetPicturesNum(cateId, keyword)
	data := m.GetPictures(cateId, keyword, pagesize, (page - 1) * pagesize)

	c.Msg["keyword"] = keyword
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}


//Profile 查看图片详情
//@params	id
//@return	Picture
func (c *PictureController) Profile () {
	id := c.MustInt("id")
	if data,err := m.GetOnePictureById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}

}

//AddPicture 添加图片
//@params	cateId modelId title description permission
//@return	success/error
func (c *PictureController) Add() {
	//cateId := c.MustInt("cateId")
	//title := c.MustStr("title")
	//description := c.Str("description")
	//
	//if err := m.CreateOnePicture(cateId, title, description); err != nil {
	//	c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	//	return
	//}
	//c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//UpdatePicture 更新图片
//@params	id cateId modelId title description permission
//@return	success/error
func (c *PictureController) Update() {
	id := c.MustInt("id")
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	description := c.Str("description")

	if err := m.UpdatePicture(id, cateId, title, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//DelPicture 删除图片
//@params	id
//@return	error
func (c *PictureController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOnePicture(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//SortPicture 对图片进行排序
//@params	id sort
//@return	success/error
func (c *PictureController) Sort() {
	postdata := c.Ctx.Request.PostForm
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdatePictureSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}

