package controllers

import (
	cnf "YYCMS/conf"
	m "YYCMS/models"
	"YYCMS/utils/YYLog"
)

type LogController struct {
	LoginController
}

func (c *LogController) List() {
	page := c.Int("page")
	pagesize := c.Int("pagesize")

	if page == 0 {
		page = 1
	}

	if pagesize <= 0 {
		pagesize = cnf.DefaultPageSize
	}

	count := m.GetAllLogsNums()

	c.Msg["count"] = count
	c.Msg["page"] = page

	if logs, err := m.ReadAllLogs(pagesize, (page-1)*pagesize); err != nil {
		YYLog.Error(err)
		c.AjaxMsg(nil, m.DataBaseGetError, m.ErrInfo[m.DataBaseGetError], "")
	} else {
		c.Msg["lists"] = logs
		c.AjaxMsg(c.Msg, m.NoError, "", "")
	}
}
