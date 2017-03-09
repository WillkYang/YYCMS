package controllers

import (
	m "YYCMS/models"
	"YYCMS/utils/YYLog"
)

type SystemController struct {
	LoginController
}

func (c *SystemController) Profile() {
	if system, err := m.ReadOrCreateOneSystem(); err != nil {
		YYLog.Error(err)
		c.AjaxMsg(nil, m.DataBaseGetError, m.ErrInfo[m.DataBaseGetError], "")
	} else {
		msg := make(map[string]interface{})
		userInfo := c.User
		userInfo.Password = ""
		msg["userInfo"] = userInfo
		msg["systemInfo"] = system
		msg["roleInfo"] = m.ReadOneRoleCatesWithActions(c.User.Role)
		c.AjaxMsg(msg, m.NoError, "", "")
	}
}
