package controllers

import (
	"strings"
	"github.com/astaxie/beego"
	"fmt"
	"strconv"
	m "YYCMS/models"
	"github.com/agelinazf/egb"
	cnf "YYCMS/conf"
)

type LoginController struct {
	BaseController

	User m.AdminUser

	Msg map[string]interface{}
}

func (c *LoginController) Prepare() {
	c.EnableRender = false

	c.Msg = make(map[string]interface{})

	if strings.Contains(c.Ctx.Input.URL(),"list") || strings.Contains(c.Ctx.Input.URL(),"profile") {
		c.MustGet()
	}

	if strings.Contains(c.Ctx.Input.URL(),"add") || strings.Contains(c.Ctx.Input.URL(),"update")  || strings.Contains(c.Ctx.Input.URL(),"delete") || strings.Contains(c.Ctx.Input.URL(),"sort") {
		c.MustPost()
	}

}


//登陆页面
func (c *LoginController) Get() {

}

//Login 执行登录
//@params   username password key word
//@return   error
func (c *LoginController) Post() {

	account := c.MustStr("username")
	password := c.MustStr("password")
	codekey := c.MustStr("key")
	codeword := c.MustStr("word")

	//验证验证码
	if isSuc, result := VerifyCaptcha(codekey, codeword); !isSuc {
		beego.Debug(isSuc,result)
		c.finish(fmt.Errorf(strconv.Itoa(m.CaptchaError)))
		return
	}
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	user,err := m.AdminUserLogin(account, password, ip);
	if err != nil {
		c.finish(err)
		return
	}

	authkey := egb.StringMD5Hex(ip + "|" + user.Password)
	c.Ctx.SetCookie(cnf.AdminAuthCookieKey, egb.StringFromInt(user.Id) + "|" + authkey, 60 * 60 * 2)
	c.User = *user

	msg := make(map[string]interface{})
	msg["key"] = cnf.AdminAuthCookieKey
	msg["value"] = egb.StringFromInt(user.Id) + "|" + authkey
	c.finish(msg)
}


func (c *LoginController) List() {
	c.ajaxMsg("hehe",0,"","")
}


