package controllers

import (
	m "YYCMS/models"
	cnf "YYCMS/conf"
	"github.com/agelinazf/egb"
	"github.com/astaxie/beego"
)

type CategoryController struct {
	LoginController
}

//
//func (c *CategoryController) Prepare() {
//	beego.Debug("找到控制器")
//	c.LoginController.Prepare()
//
//	//TODO - 添加用户权限验证
//	beego.Debug("完成父控制器初始化")
//}

//获取全部栏目信息
func (c *CategoryController) List() {
	c.MustGet()
	topCateId := cnf.TopCateId
	beego.Debug(c.Ctx.Input.Query("isRecursion"))
	isrecursion := c.Int("isRecursion")
	if c.Int("topCateId") > 0 {
		topCateId = c.Int("topCateId")
	}
	cates := m.GetSubCategorysById(topCateId, isrecursion > 0)
	beego.Debug(cates)
	c.Msg["keyword"] = ""
	c.Msg["count"] = len(cates)
	c.Msg["lists"] = cates
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

////获取后台顶部栏目信息
//func (c *CategoryController) TopCates() {
//	cates := m.GetSubCategorysById(cnf.TopCateId, false)
//	c.AjaxMsg(cates, m.NoError, "", "")
//}

////获取全部栏目类型
//func (c *CategoryController) Models() {
//	models := m.GetAllModel()
//	c.Msg["models"] = models
//	c.AjaxMsg(c.Msg, m.NoError, "", "")
//}


//Profile 查看栏目详情
//@params	id
//@return	AdminRole
func (c *CategoryController) Profile () {
	c.MustGet()
	id := c.MustInt("id")
	if data,err := m.GetOneCategoryById(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
	} else {
		c.AjaxMsg(data, m.NoError, "", "")
	}
}

//Add 添加栏目动作
//@params   pid name catetype modelid
//@return   success/error
func (c *CategoryController) Add() {
	c.MustPost()

	parentId := c.MustInt("parentId")
	cateName := c.MustStr("title")
	modelId := c.MustInt("modelId")
	description := c.Str("description")

	if err := m.CreateOneCategory(modelId,cateName,description,parentId); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "添加成功")
		return
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "添加栏目成功")
}

//Update 编辑栏目动作
//@params	id name
//@return	success/error
func (c *CategoryController) Update() {
	c.MustPost()
	id := c.MustInt("id")
	title := c.MustStr("title")
	modelId := c.MustInt("modelId")
	parentId := c.MustInt("parentId")
	description := c.Str("description")

	if err := m.UpdateCategory(id,modelId,parentId,title,description); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(),"")
		return
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "编辑栏目成功")
}

//Delete 删除栏目
//@params	id
//@return	success/error
func (c *CategoryController) Delete() {
	c.MustPost()
	id := c.MustInt("id")
	if err := m.DeleteOneCategory(id); err != nil {
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}
	m.UpdateCache()
	c.AjaxMsg(nil, m.NoError, "", "删除成功")
}

//Sort 栏目排序
//params	int : int 栏目Id : 新sort值
//@return	success/error
func (c *CategoryController) Sort() {
	c.MustPost()
	postdata := c.Ctx.Request.PostForm
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

