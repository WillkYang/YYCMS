package YYCMS

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	_ "YYCMS/routers"
	"YYCMS/utils"
	"YYCMS/utils/YYLog"
	"encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"fmt"
	"os/exec"
)

func InitCMS() {
	RegisterModel()
	InitSystemLog()
	InitViewFunc()
	//InitDir()
}

func getAllModels() map[string]interface{} {
	ms := make(map[string](interface{}))
	ms["kvdb"] = new(m.Kvdb)
	ms["log"] = new(m.Log)
	ms["model"] = new(m.Model)
	ms["picture"] = new(m.Picture)
	ms["tag"] = new(m.Tag)
	ms["gallery"] = new(m.Gallery)
	ms["admin_user"] = new(m.AdminUser)
	ms["category"] = new(m.Category)
	ms["single"] = new(m.SinglePage)
	ms["article"] = new(m.Article)
	ms["video"] = new(m.Video)
	ms["system"] = new(m.System)
	ms["carousel"] = new(m.Carousel)
	ms["admin_role"] = new(m.AdminRole)
	ms["privilege"] = new(m.Privilege)
	ms["role_privilege"] = new(m.RolePrivilege)
	ms["news"] = new(m.News)
	ms["blog"] = new(m.Blog)
	ms["comment"] = new(m.Comment)
	return ms
}

func RegisterModel() {
	YYLog.Info("YYCMS:RegisterModel")
	//defer YYLog.Info("YYCMS:RegisterModel OK")
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	orm.RegisterModel(
		new(m.Model),
		new(m.Kvdb),
		new(m.Privilege),
		new(m.RolePrivilege),
	)
	//可选注册
	if result, err := utils.ReadFileToArray(cnf.ModelConfigFilePath); err != nil {
		YYLog.Error(err)
	} else {
		ms := getAllModels()
		for _, v := range result {
			v1 := v.(map[string]interface{})
			if bool((v1["isNeed"]).(bool)) {
				if len(v1["tablename"].(string)) > 0 {
					orm.RegisterModel(ms[v1["tablename"].(string)])
				}
			}
		}
	}
	YYLog.Info("YYCMS:RegisterModel OK")
}

func InitSystemLog() {
	YYLog.Info("YYCMS:InitSystemLog")
	defer YYLog.Info("YYCMS:InitSystemLog OK")
	runmode := beego.AppConfig.String("runmode")
	logpath := beego.AppConfig.String("logpath")
	os.Mkdir(logpath, os.ModePerm)
	conf := make(map[string]interface{}, 0)
	conf["filename"] = logpath + "log.log"
	conf["maxlines"] = 10000000
	conf["maxsize"] = 1 << 28
	conf["daily"] = true
	conf["maxdays"] = 30
	conf["rotate"] = true
	conf["separate"] = []string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}
	confjson, _ := json.Marshal(conf)
	beego.SetLogger("multifile", string(confjson))
	//beego.SetLevel(beego.LevelInformational)
	//beego.SetLevel(beego.LevelWarning)
	beego.SetLogFuncCall(true)
	if runmode == "pro" {
		beego.BeeLogger.DelLogger("console")
	}
}

func InitSystemAdmin() {
	YYLog.Info("YYCMS:InitSystemAdmin")
	//defer YYLog.Info("YYCMS:InitSystemAdmin OK")
	if err := m.CreateOneAdminUser("yeeyun_root", egb.StringSHA1Hex(cnf.DefaultAdminPassword), cnf.SuperAdminRoleId); err != nil {
		YYLog.Error(err)
		return
	}
	if err := m.CreateOneAdminUser("yeeyun", egb.StringSHA1Hex(cnf.DefaultAdminPassword), cnf.SuperAdminRoleId); err != nil {
		YYLog.Error(err)
		return
	}
	YYLog.Info("YYCMS:InitSystemAdmin OK")
}

func InitSystemModel() {
	YYLog.Info("YYCMS:initSystemModel")
	//defer YYLog.Info("YYCMS:initSystemModel OK")
	if result, err := utils.ReadFileToArray(cnf.ModelConfigFilePath); err != nil {
		YYLog.Error(err)
	} else {
		for _, v := range result {
			v1 := v.(map[string]interface{})
			if bool((v1["isNeed"]).(bool)) {
				if err := m.CreateOneModel(int(v1["id"].(float64)), v1["name"].(string), v1["name"].(string), v1["tablename"].(string), v1["actions"].(string), bool((v1["isShow"]).(bool))); err != nil {
					YYLog.Error(err)
					//return
				}
			}
		}
	}
	YYLog.Info("YYCMS:initSystemModel OK")

}

func InitViewFunc() {
	YYLog.Info("YYCMS:InitViewFunc")
	//defer YYLog.Info("YYCMS:InitViewFunc OK")
	beego.AddFuncMap("StringToInt", egb.StringToInt)
	beego.AddFuncMap("StringToInt64", egb.StringToInt64)
	beego.AddFuncMap("StringFromInt", egb.StringFromInt)
	beego.AddFuncMap("StringFromInt64", egb.StringFromInt64)
	beego.AddFuncMap("TimeFromUnix", egb.TimeFromUnix)

	YYLog.Info("YYCMS:InitViewFunc OK")

}

func InitDir() {
	YYLog.Info("YYCMS:InitDir")

	os.MkdirAll("upload/public/", os.ModePerm)
	os.MkdirAll("upload/private/", os.ModePerm)
	os.MkdirAll(cnf.TransferThumbImagePath, os.ModePerm)
	os.MkdirAll(cnf.TransferHtmlImagePath, os.ModePerm)
	os.MkdirAll("logs", os.ModePerm)
	os.MkdirAll(cnf.CategoryCachePath, os.ModePerm)

	YYLog.Info("YYCMS:InitDir OK")
}

func CreateDataBase() {
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	dbname := beego.AppConfig.String("dbname")
	dbaddr := beego.AppConfig.String("dbaddr")
	dbport := beego.AppConfig.String("dbport")
	sql := fmt.Sprintf("create database IF NOT EXISTS %s default charset utf8mb4 collate utf8mb4_general_ci", dbname)
	create_db_shell := fmt.Sprintf("mysql -h%s -P%s -u%s -p%s -e \"%s\"", dbaddr, dbport, username, password, sql)

	cmd := exec.Command("/bin/bash", "-c", create_db_shell)

	_, err := cmd.StdoutPipe()
	if err != nil {
		YYLog.Error("Error:can not obtain the stdout pipe for command No.0: %s", err)
		return
	}
	err = cmd.Start()

	if err != nil {
		YYLog.Error(err)
	}

	err = cmd.Wait()
	if err != nil {
		YYLog.Error("Command finished with error: %v", err)
	}
}
