package helper

import (
	"YYCMS/utils/YYLog"
	"YYCMS/utils"
	cnf "YYCMS/conf"
)

func SystemAllAction() map[string]interface{} {
	result, err := utils.ReadFileToMap(cnf.AllActionsPath)
	if err != nil {
		YYLog.Error(err)
	}
	return result
}

