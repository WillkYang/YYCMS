package controllers

import (
	m "YYCMS/models"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strings"
)

type BaseController struct {
	beego.Controller
}

//Display 渲染模版
func (c *BaseController) Display(tpl string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl + ".html"
	}
	c.TplName = tplname
}

//获取用户IP地址
func (c *BaseController) GetClientIp() string {
	s := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

//ShowMsg 显示错误信息
func (c *BaseController) ShowErrorPage(args ...string) {
	c.Data["message"] = args[0]
	redirect := c.Ctx.Request.Referer()
	if len(args) > 1 {
		redirect = args[1]
	}
	c.Data["redirect"] = redirect
	c.Data["pageTitle"] = "系统提示"
	c.Display("error/error")
}

//AjaxMsg 返回数据
func (c *BaseController) AjaxMsg(msg interface{}, errCode int, errMsg, sucMsg string) {
	out := make(map[string]interface{})
	out["errCode"] = errCode
	out["errMsg"] = errMsg
	out["sucMsg"] = sucMsg
	out["data"] = msg
	c.jsonResult(out)
}


func (c *BaseController) jsonResult(out interface{}) {
	c.Data["json"] = out
	c.ServeJSON()
}

func (c *BaseController) ajaxMsg(msg interface{}, errCode int, errMsg string, sucMsg string) {
	out := make(map[string]interface{})
	out["errCode"] = errCode
	out["errMsg"] = errMsg
	out["sucMsg"] = sucMsg
	out["data"] = msg
	c.Data["json"] = out
	c.ServeJSON()
	c.StopRun()
}

//MustInt 获取int类型参数 ，参数不能为空
func (c *BaseController) MustInt(param string) int {
	if strv := c.Ctx.Input.Query(param); strv != "" {
		i, err := strconv.Atoi(strv)
		if err != nil {
			paramerror := "参数错误 : " + err.Error()
			c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
		}
		return i
	}
	paramerror := "参数 " + param + " 不能为空"
	c.ajaxMsg(nil, m.ParamsMissError, paramerror, "")
	return 0
}

//Int 获取int类型参数,参数可以为空
func (c *BaseController) Int(param string) int {
	if strv := c.Ctx.Input.Query(param); strv != "" {
		i, err := strconv.Atoi(strv)
		if err != nil {
			paramerror := "参数错误 : " + err.Error()
			c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
		}
		return i
	}
	return 0
}

//MustFloat 获取float64类型参数.,参数不能为空
func (c *BaseController) MustFloat(param string) float64 {
	if strv := c.Ctx.Input.Query(param); strv != "" {
		f, err := strconv.ParseFloat(strv, 64)
		if err != nil {
			paramerror := "参数错误 : " + err.Error()
			c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
		}
		return f
	}
	paramerror := "参数 " + param + " 不能为空"
	c.ajaxMsg(nil, m.ParamsMissError, paramerror, "")
	return 0.0
}

//Float 获取float64类型参数.,参数可以为空
func (c *BaseController) Float(param string) float64 {
	if strv := c.Ctx.Input.Query(param); strv != "" {
		f, err := strconv.ParseFloat(strv, 64)
		if err != nil {
			paramerror := "参数错误 : " + err.Error()
			c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
		}
		return f
	}
	return 0.0
}

//MustBool 获取bool类型参数,参数不可以为空
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
func (c *BaseController) MustBool(param string) bool {
	if strv := c.Ctx.Input.Query(param); strv != "" {
		b, err := strconv.ParseBool(strv)
		if err != nil {
			paramerror := "参数错误 : " + err.Error()
			c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
		}
		return b
	}
	paramerror := "参数 " + param + " 不能为空"
	c.ajaxMsg(nil, m.ParamsMissError, paramerror, "")
	return true
}

//Bool 获取bool类型参数,参数可以为空    如果为空则为true
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
func (c *BaseController) Bool(param string) bool {
	if strv := c.Ctx.Input.Query(param); strv != "" {
		b, err := strconv.ParseBool(strv)
		if err != nil {
			paramerror := "参数错误 : " + err.Error()
			c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
		}
		return b
	}
	return true
}

//MustStr 获取string类型参数 ，参数不能为空
func (c *BaseController) MustStr(param string) string {
	str := c.Ctx.Input.Query(param)
	if str == "" {
		paramerror := "参数 " + param + " 不能为空"
		c.ajaxMsg(nil, m.ParamsTypeError, paramerror, "")
	}
	return str
}

//Str 获取string类型参数,参数可以为空
func (c *BaseController) Str(param string) string {
	str := c.Ctx.Input.Query(param)
	return str
}

func (c *BaseController) isGET() bool {
	if c.Ctx.Request.Method == "GET" {
		return true
	}
	return false
}

func (c *BaseController) isPOST() bool {
	if c.Ctx.Request.Method == "POST" {
		return true
	}
	return false
}

func (c *BaseController) MustPost() {
	if c.Ctx.Request.Method != "POST" {
		paramerror := "必须为post请求"
		c.ajaxMsg(nil, m.SystemError, paramerror, "")
	}
}

func (c *BaseController) MustGet() {
	if c.Ctx.Request.Method == "POST" {
		paramerror := "必须为get请求"
		c.ajaxMsg(nil, m.SystemError, paramerror, "")
	}
}

func (c *BaseController) finish (result interface{}) {
	beego.Debug("处理返回")
	if err,ok := result.(error); ok {
		//err := error(result.(error))
		errcode,_ := strconv.Atoi(err.Error())
		beego.Debug("进入err处理方法")
		c.AjaxMsg(nil, errcode, m.ErrInfo[errcode], "")
	} else {
		beego.Debug("进入成功回调")
		c.AjaxMsg(result, m.NoError, "", "")
	}
}

func (c *BaseController) handleValidParam(obj interface{}) {
	valid := validation.Validation{}
	b,err := valid.Valid(obj)
	if err != nil {
		beego.Error(err)
	}
	if !b {
		for _,err := range valid.Errors {
			beego.Error(err)
			c.AjaxMsg(nil,m.ParamsTypeError,err.Message,"")
		}
	}
}