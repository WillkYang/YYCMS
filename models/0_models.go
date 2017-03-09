package models

import (
	cnf "YYCMS/conf"
	"os"

	"github.com/astaxie/beego/orm"
)

// 注 status为99表示系统保留，非删除却不可见

func ormer() orm.Ormer {
	return orm.NewOrm()
}

//UpdateCache 更新缓存
//@params	nil
//@return	nil
func UpdateCache() {
	os.RemoveAll(cnf.ModelCachePath)
	os.RemoveAll(cnf.CategoryCachePath)
}
