package controllers

import (
	cnf "YYCMS/conf"
	"YYCMS/helper"
	m "YYCMS/models"
	"YYCMS/utils"
	"YYCMS/utils/YYLog"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/agelinazf/egb"
)

type LoginController struct {
	BaseController
	User m.AdminUser
}

//权限访问
var OpenPerm = map[string]bool{
	"LoginController.*":        true,
	"SystemController.Profile": true,
	"ModelController.List":     true,
}

//登陆访问
var OpenAuth = map[string]bool{
	"LoginController.*": true,
}

func (c *LoginController) Prepare() {
	c.BaseController.Prepare()
	c.EnableRender = false

	//身份认证
	c.checkIdentity()

	//方法检查
	c.CheckRequestMothod()
}

//检查权限
func (c *LoginController) checkIdentity() {
	//是否需要登陆
	if c.isOpen(OpenAuth) {
		return
	}

	//登陆检查
	c.Auth()

	//是否为开放权限
	if c.isOpen(OpenPerm) {
		return
	}

	//权限检查
	c.checkPriv()

	//requestURI := strings.ToLower(c.Ctx.Request.RequestURI)
	////身份认证
	//if !(strings.Contains(requestURI, "/login/admin") || strings.Contains(requestURI, "/cms") || strings.Contains(requestURI, "/login/bgimage")) {
	//	YYLog.Info("登录检查")
	//	c.Auth()
	//	YYLog.Info(fmt.Sprintf("用户%s访问了%s", c.User.Account, requestURI))
	//	YYLog.Info("权限检查")
	//	if !(strings.Contains(requestURI, "/login/admin") || strings.Contains(requestURI, "/cms") || strings.Contains(requestURI, "/login/bgimage") || strings.Contains(requestURI, "/login/exit") || strings.Contains(requestURI, "/model/list") || strings.Contains(requestURI, "/profile/10")) {
	//		c.checkPriv()
	//	}
	//}
}

//判断是否开放
func (c *BaseController) isOpen(openDict map[string]bool) bool {
	if c.isOpenOfController(openDict) || c.isOpenPermOfAction(openDict) {
		return true
	}
	return false
}

func (c *BaseController) isOpenOfController(openPerm map[string]bool) bool {
	YYLog.Warning(c.controllerName, c.actionName)
	if openPerm[fmt.Sprint(c.controllerName, ".*")] {
		return true
	}
	return false
}

func (c *BaseController) isOpenPermOfAction(openPerm map[string]bool) bool {
	YYLog.Warning(c.controllerName, c.actionName)
	if openPerm[fmt.Sprint(c.controllerName, ".", c.actionName)] {
		return true
	}
	return false
}

func (c *LoginController) CheckRequestMothod() {
	requestURI := strings.ToLower(c.Ctx.Request.RequestURI)
	//方法检查
	if strings.Contains(requestURI, "/list/") || strings.Contains(requestURI, "/profile/") {
		c.MustGet()
		return
	}
	if strings.Contains(requestURI, "/add/") || strings.Contains(requestURI, "/update/") || strings.Contains(requestURI, "/delete/") || strings.Contains(requestURI, "/sort/") || strings.Contains(requestURI, "/login/admin") {
		c.MustPost()
		return
	}
}

//Login 执行登录
//@params   username password key word
//@return   error
func (c *LoginController) Admin() {

	account := c.MustStr("username")
	password := c.MustStr("password")
	//codekey := c.MustStr("key")
	//codeword := c.MustStr("word")

	////验证验证码
	//if isSuc, result := VerifyCaptcha(codekey, codeword); !isSuc {
	//	YYLog.Debug(isSuc,result)
	//	c.finish(fmt.Errorf(strconv.Itoa(m.CaptchaError)))
	//	return
	//}
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	user, err := m.AdminUserLogin(account, password, ip)
	if err != nil {
		YYLog.Error(err)
		c.AjaxMsg(nil, m.ErrCode[err.Error()], err.Error(), "")
		return
	}

	c.User = *user
	c.UpdateCookie()
	msg := make(map[string]interface{})
	userInfo := c.User
	userInfo.Password = ""
	msg["userInfo"] = userInfo
	msg["systemInfo"], _ = m.ReadOrCreateOneSystem()
	msg["roleInfo"] = m.ReadOneRoleCatesWithActions(c.User.Role)
	m.CreateOneLog(userInfo.Account, "登录成功")
	c.AjaxMsg(msg, m.NoError, "", "登录成功")
}

func (c *LoginController) Auth() {
	arr := strings.Split(c.Ctx.GetCookie(cnf.AdminAuthCookieKey), "|")
	if len(arr) == 2 {
		idstr, authScrect := arr[0], arr[1]
		userId, _ := strconv.Atoi(idstr)
		if userId > 0 {
			user, err := m.GetOneAdminUserById(userId)

			if err != nil {
				c.ClearCookie()
				c.AjaxMsg(nil, m.CookieError, m.ErrInfo[m.CookieError], "")
				return
			}

			if authScrect == egb.StringMD5Hex(c.GetClientIp()+"|"+user.Password) {
				c.User = *user
				c.Data["User"] = c.User
				c.UpdateCookie()
				//c.Ctx.SetCookie(cnf.AdminAuthCookieKey, strconv.Itoa(user.Id)+"|"+authScrect, 60 * 20)
			} else {
				YYLog.Info("校验不通过")
				c.AjaxMsg(nil, m.CookieError, m.ErrInfo[m.CookieError], "")
				c.StopRun()
			}
		} else {
			YYLog.Info("未获取到用户")
			c.AjaxMsg(nil, m.CookieError, m.ErrInfo[m.CookieError], "")
			c.StopRun()
		}
	} else {
		YYLog.Info("未获取到校验信息")
		c.AjaxMsg(nil, m.CookieError, m.ErrInfo[m.CookieError], "")
		c.StopRun()
	}
}

func (c *LoginController) checkPriv() {
	if c.User.Role == cnf.SuperAdminRoleId {
		return
	}

	cateId := c.MustInt("cateId")
	//权限检查
	result, _ := m.GetOneCategoryRawById(cateId)
	YYLog.Warning(result)
	if result["Actions"] != nil {
		for _, action := range utils.StringToStringArray(result["Actions"].(string)) {
			if strings.Contains(c.Ctx.Request.RequestURI, action) {
				if !checkPrivilege(c.User.Role, cateId, action) {
					c.AjaxMsg(nil, m.NoPermissonError, m.ErrInfo[m.NoPermissonError], "")
					c.StopRun()
				}
			}
		}
	}
}

func GetAllActions() []string {
	dict := helper.SystemAllAction()
	var allKeys []string
	for k, _ := range dict {
		allKeys = append(allKeys, k)
	}
	return allKeys
	//return []string{
	//	"add",
	//	//"list",
	//	"update",
	//	"delete",
	//	//"sort",
	//	//"profile",
	//	"review",
	//	"mark",
	//}

}

func checkPrivilege(roleId, cateId int, action string) bool {
	privilege := m.ReadOnePrivilege(cateId, m.Action(action))
	if privilege.Id > 0 {
		if m.RolePrivilegeIsExist(roleId, privilege.Id) {
			return true
		}
	}
	return false
}

func (c *LoginController) UpdateCookie() {
	authScrect := egb.StringMD5Hex(c.GetClientIp() + "|" + c.User.Password)
	c.Ctx.SetCookie(cnf.AdminAuthCookieKey, strconv.Itoa(c.User.Id)+"|"+authScrect, 60*60*2)
}

func (c *LoginController) ClearCookie() {
	c.Ctx.SetCookie(cnf.AdminAuthCookieKey, "", 60*60*2)
}

func (c *LoginController) Page() {
	c.EnableRender = true
	c.TplName = "admin/index.html"
}

//BgImage bing背景图
func (c *LoginController) BgImage() {
	res, err := egb.HttpGetBytes(&http.Client{}, "http://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1", nil)
	if err != nil {
		YYLog.Error("BgImage :" + err.Error())
		c.AjaxMsg(nil, m.SystemError, "请求bing接口失败", "")
		return
	}
	c.Msg["data"] = string(res)
	c.AjaxMsg(c.Msg, m.NoError, "", "")
}

func (c *LoginController) Exit() {
	c.Auth()
	m.CreateOneLog(c.User.Account, "退出登录")
	c.ClearCookie()
	c.AjaxMsg(nil, m.NoError, "", "退出登录成功")
}
