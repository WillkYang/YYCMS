package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"encoding/json"
	cnf "YYCMS/conf"
	"YYCMS/utils"
	"github.com/agelinazf/egb"
)

func init() {
	initDataBase()
	initSystemLog()
	//initSystemAdmin()
	//initSystemModel()
	initViewFunc()
}

func initDataBase() {

	beego.Debug("initDB")

	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	orm.RegisterModel(
		new(Model),
		new(Picture),
		new(Tag),
		new(Gallery),
		new(AdminUser),
		new(Category),
		new(Single),
		new(Article),
		new(Video),
		new(System),
		new(Carousel),
		new(AdminRole),
	)
	beego.Debug("initDB OK")
}

func initSystemLog() {
	beego.Debug("initLog")
	runmode := beego.AppConfig.String("runmode")
	logpath := beego.AppConfig.String("logpath")
	os.Mkdir(logpath, os.ModePerm)
	conf := make(map[string]interface{}, 0)
	conf["filename"] = logpath + "log.log"
	conf["maxlines"] = 1000000
	conf["maxsize"] = 1 << 28
	conf["daily"] = true
	conf["maxdays"] = 30
	conf["rotate"] = true
	confjson, _ := json.Marshal(conf)
	beego.SetLogger("file", string(confjson))
	beego.SetLogFuncCall(true)
	if runmode == "pro" {
		beego.BeeLogger.DelLogger("console")
	}
}

func InitSystemAdmin() {
	beego.Debug("initSystemAdmin")
	if err := CreateOneAdminUser("yeeyun_root","PTpGWn4E",0); err != nil {
		return
	}
	if err := CreateOneAdminUser("yeeyun","PTpGWn4E",0); err != nil {
		return
	}
	beego.Debug("initSystemAdmin OK")
}

func InitSystemModel() {
	beego.Debug("initSystemModel")
	if result,err := utils.ReadFileToArray(cnf.ModelConfigFilePath); err != nil {
		beego.Error(err)
	} else {
		for _,v := range result {
			v1 := v.(map[string]interface{})
			if bool((v1["isNeed"]).(bool)) {
				if err := CreateOneModel(int(v1["id"].(float64)), v1["name"].(string), v1["name"].(string), v1["tablename"].(string)); err != nil {
					return
				}
			}
		}
	}
	beego.Debug("initSystemModel OK")
}


func initViewFunc() {
	beego.AddFuncMap("StringToInt",egb.StringToInt)
	beego.AddFuncMap("StringToInt64",egb.StringToInt64)
	beego.AddFuncMap("StringFromInt",egb.StringFromInt)
	beego.AddFuncMap("StringFromInt64",egb.StringFromInt64)
	beego.AddFuncMap("TimeFromUnix",egb.TimeFromUnix)
}

