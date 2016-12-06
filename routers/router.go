package routers

import (
	"github.com/astaxie/beego"
	"YYCMS/controllers"
)


func init() {

	beego.Debug("initCMSRoutes")

	beego.AutoRouter(&controllers.AdminUserController{})
	beego.Router("/adminlogin", &controllers.LoginController{})
	beego.AutoRouter(&controllers.AuthCodeController{})
	beego.AutoRouter(&controllers.AdminRoleController{})
	beego.AutoRouter(&controllers.CategoryController{})
	beego.AutoRouter(&controllers.ModelController{})
	beego.AutoRouter(&controllers.CarouselController{})

	beego.AutoRouter(&controllers.ImageController{})
	beego.AutoRouter(&controllers.ImageController{})
	beego.Debug("initCMSRoutes OK")
}



