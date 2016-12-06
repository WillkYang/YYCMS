package controllers
import (
	"github.com/agelinazf/egb"
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"github.com/astaxie/beego"
)

type TagController struct {
	LoginController
}

//List 获取标签列表
//@params   keyword(搜索title的关键词)
//@return   []Tag
func (c *TagController) List() {
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

	count := m.GetTagsNum(cateId, keyword)
	beego.Info(count)
	data := m.GetTags(cateId, keyword, pagesize, (page - 1) * pagesize)
	beego.Info(data)

	c.Msg["keyword"] = keyword
	c.Msg["count"] = count
	c.Msg["lists"] = data
	c.Msg["page"] = page
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}


//Profile 查看标签详情
//@params	id
//@return	Tag
func (c *TagController) Profile () {
	id := c.MustInt("id")
	if data,err := m.GetOneTagById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}

}

//AddTag 添加标签
//@params	cateId modelId title description
//@return	success/error
func (c *TagController) Add() {
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	description := c.Str("description")

	//判读是否已有相同标签-若有则直接返回
	if _,err := m.GetOneTagByTitle(cateId,title); err == nil {
		c.AjaxMsg(nil, m.NoError, "", "添加成功")
		return
	}

	if _,err := m.CreateOneTag(cateId, title, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "添加成功")
}

//UpdateTag 更新标签
//@params	id cateId modelId title description
//@return	success/error
func (c *TagController) Update() {
	id := c.MustInt("id")
	cateId := c.MustInt("cateId")
	title := c.MustStr("title")
	description := c.Str("description")

	//判读是否已有相同标签-若有则直接返回
	if _,err := m.GetOneTagByTitle(cateId,title); err == nil {
		c.AjaxMsg(nil, m.SystemError, "存在相同的标签", "")
		return
	}

	if err := m.UpdateTag(id, cateId, title, description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "编辑成功")
}

//DelTag 删除标签
//@params	id
//@return	error
func (c *TagController) Delete() {
	id := c.MustInt("id")
	if err := m.DeleteOneTag(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//SortTag 对标签进行排序
//@params	id sort
//@return	success/error
func (c *TagController) Sort() {
	postdata := c.Ctx.Request.PostForm
	beego.Debug(postdata)
	for k, v := range postdata {
		id := egb.StringToInt(k)
		if err := m.UpdateTagSort(id, egb.StringToInt(v[0])); err != nil {
			c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
			return
		}
	}
	c.AjaxMsg(nil, m.NoError, "", "排序成功")
}