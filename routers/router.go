package routers

import (
	"YYCMS/controllers"
	"YYCMS/utils/YYLog"

	"github.com/astaxie/beego"
)

func init() {
	YYLog.Info("initYYCMSRoutes")
	beego.AutoRouter(&controllers.AdminUserController{})
	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.AuthCodeController{})
	beego.AutoRouter(&controllers.AdminRoleController{})
	beego.AutoRouter(&controllers.CategoryController{})
	beego.AutoRouter(&controllers.ModelController{})
	beego.AutoRouter(&controllers.CarouselController{})
	beego.AutoRouter(&controllers.ImageController{})
	//ueditor
	beego.Router("/ueditor", &controllers.UeditorController{}, "*:ControllerUE")
	YYLog.Info("initYYCMSRoutes OK")
}
