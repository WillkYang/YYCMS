package test

import (
	"YYCMS/models"
	_ "YYCMS/routers"
	"YYCMS/utils/YYLog"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestMain is a sample to run an endpoint test
func TestMain(t *testing.T) {
	privilege := &models.RolePrivilege{
		RoleId:      1,
		PrivilegeId: 1,
	}
	YYLog.Debug(privilege)
}
