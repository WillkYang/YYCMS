package controllers


import (
	"github.com/agelinazf/egb"
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"github.com/astaxie/beego"
)


type GalleryController struct {
	LoginController
}


//List 获取相册列表
//@params   keyword(搜索title的关键词)
//@return   []Gallery
func (c *GalleryController) List() {
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

	count := m.GetGallerysNum(cateId, keyword)
	data := m.GetGallerys(cateId, keyword, pagesize, (page - 1) * pagesize)

	c.Msg["keyword"] = keyword
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}


//Profile 查看相册详情
//@params	id
//@return	Gallery
func (c *GalleryController) Profile () {
	id := c.MustInt("id")
	if data,err := m.GetOneGalleryById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}

}

//AddGallery 添加相册
//@params	cateId modelId title description
//@return	success/error
func (c *GalleryController) Add() {
	cateId := c.MustInt("cateId")
	modelId := c.Int("modelId")
	if cate,err := m.GetOneCategoryById(cateId); err != nil {
		c.AjaxMsg(nil,m.DataBaseGetError,m.ErrInfo[m.DataBaseGetError],"")
	} else {
		modelId = cate.ModelId
	}
	title := c.MustStr("title")
	description := c.Str("description")

	if err := m.CreateOneGallery(cateId, modelId, title, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//UpdateGallery 更新相册
//@params	id cateId modelId title description
//@return	success/error
func (c *GalleryController) Update() {
	id := c.MustInt("id")
	cateId := c.MustInt("cateId")
	modelId := c.Int("modelId")
	if cate,err := m.GetOneCategoryById(cateId); err != nil {
		c.AjaxMsg(nil,m.DataBaseGetError,m.ErrInfo[m.DataBaseGetError],"")
	} else {
		modelId = cate.ModelId
	}
	title := c.MustStr("title")
	description := c.Str("description")

	if err := m.UpdateGallery(id, cateId, modelId, title, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//DelGallery 删除相册
//@params	id
//@return	error
func (c *GalleryController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneGallery(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//SortGallery 对相册进行排序
//@params	id sort
//@return	success/error
func (c *GalleryController) Sort() {
	postdata := c.Ctx.Request.PostForm
	beego.Debug(postdata)
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateGallerySort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}

