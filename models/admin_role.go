package models

import (
	"YYCMS/conf"
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/astaxie/beego/orm"
)

//角色表
type AdminRole struct {
	//主键
	Id int `orm:"column(Id);pk;auto"`
	//栏目Id
	CateId int `orm:"column(CateId)"`
	//名字
	Title string `orm:"column(Title)"`
	//排序值
	Sort int `orm:"column(Sort);default(0)"`
}

func (t *AdminRole) TableName() string {
	return "admin_role"
}

//GetOneAdminRoleById 获取一个角色
//@params	Id
//@return	*AdminRole
func GetOneAdminRoleById(Id int) (*AdminRole, error) {
	adminrole := new(AdminRole)
	adminrole.Id = Id

	if err := ormer().Read(adminrole, "Id"); err != nil {
		YYLog.Error("GetOneAdminRoleById : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return adminrole, nil
}

//GetAdminRolesNum 获取角色的数量
//@params	cateId keyword(搜索title的关键词)
//@return	int
func GetAdminRolesNum(catId int, keyword string) int {
	var data []orm.Params
	sql := "SELECT Id FROM admin_role WHERE CateId = ? AND Title LIKE ? "
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword).Values(&data)
	return len(data)
}

//GetAdminRoles 获取角色
//@params	catId keyword(搜索title的关键词) pagesize offset
//@return	[]orm.Params
func GetAdminRoles(catId int, keyword string, pagesize, offset int) []orm.Params {

	var data []orm.Params
	sql := `SELECT * FROM admin_role WHERE CateId = ? AND Title LIKE ? ORDER BY admin_role.Sort DESC LIMIT ?,?`
	keyword = "%" + keyword + "%"
	ormer().Raw(sql, catId, keyword, offset, pagesize).Values(&data)
	return data
}

//ReadOneRole 读取某个角色
//@params	id
//@retrun	*Role error
func ReadOneRole(id int) (orm.Params, error) {

	var results []orm.Params
	if _, err := ormer().Raw(`SELECT * FROM admin_role WHERE Id = ?`, id).Values(&results); err != nil {
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return results[0], nil
}

func ReadOneRoleCates(id int) map[string]string {
	var roleCates []orm.Params
	if id == conf.SuperAdminRoleId {
		ormer().Raw(`SELECT CateId, group_concat(privilege.PrivilegeId) as RoleActionIds, group_concat(privilege.Action) as RoleActions from (SELECT privilege.*, role_privilege.PrivilegeId, role_privilege.RoleId FROM privilege, role_privilege WHERE privilege.Id = role_privilege.PrivilegeId GROUP BY privilege.Id) privilege GROUP BY CateId;`).Values(&roleCates)
	} else {
		ormer().Raw(`SELECT CateId, group_concat(role_privilege.PrivilegeId) as RoleActionIds, group_concat(privilege.Action) as RoleActions FROM privilege, role_privilege WHERE privilege.Id = role_privilege.PrivilegeId AND role_privilege.RoleId = ? GROUP BY CateId;`, id).Values(&roleCates)
	}
	roleCatesMap := make(map[string]string, 0)
	for _, roleCate := range roleCates {
		roleCatesMap[roleCate["CateId"].(string)] = roleCate["RoleActionIds"].(string)
	}
	return roleCatesMap
}

func ReadOneRoleCatesWithActions(id int) map[string]string {
	var roleCates []orm.Params
	if id == conf.SuperAdminRoleId {
		ormer().Raw(`SELECT CateId, group_concat(privilege.Action) as RoleActions FROM privilege GROUP BY CateId;`).Values(&roleCates)
	} else {
		ormer().Raw(`SELECT CateId, group_concat(role_privilege.PrivilegeId) as RoleActionIds, group_concat(privilege.Action) as RoleActions FROM privilege, role_privilege WHERE privilege.Id = role_privilege.PrivilegeId AND role_privilege.RoleId = ? GROUP BY CateId;`, id).Values(&roleCates)
	}
	roleCatesMap := make(map[string]string, 0)
	for _, roleCate := range roleCates {
		roleCatesMap[roleCate["CateId"].(string)] = roleCate["RoleActions"].(string)
	}
	return roleCatesMap
}

//CreateOneAdminRole 新建一篇角色
//@params	catId title thumb source description content
//@return	error
func CreateOneAdminRole(cateId int, title string) (int, error) {
	adminrole := new(AdminRole)
	adminrole.CateId = cateId
	adminrole.Title = title

	if index, err := ormer().Insert(adminrole); err != nil {
		YYLog.Error("CreateOneAdminRole : " + err.Error())
		return 0, fmt.Errorf(ErrInfo[SystemError])
	} else {
		return int(index), nil
	}

}

//UpdateAdminRole 更新角色
//@params	id title description
//@return	error
func UpdateAdminRole(id int, title string) error {
	adminrole, err := GetOneAdminRoleById(id)
	if err != nil {
		return err
	}
	adminrole.Title = title

	if _, err := ormer().Update(adminrole); err != nil {
		YYLog.Error("UpdateAdminRole : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		return nil
	}
}

//UpdateAdminRoleSort 更新角色排序
//@params	id sort
//@return	error
func UpdateAdminRoleSort(id, sort int) error {
	adminrole, err := GetOneAdminRoleById(id)
	if err != nil {
		return fmt.Errorf(ErrInfo[SystemError])
	}
	adminrole.Sort = sort
	if _, err := ormer().Update(adminrole, "Sort"); err != nil {
		YYLog.Error("UpdateAdminRoleSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneAdminRole 删除一个角色
//@params	id
//@return	error
func DeleteOneAdminRole(id int) error {
	adminrole, err := GetOneAdminRoleById(id)
	if err != nil {
		return err
	}
	if _, err := ormer().Delete(adminrole); err != nil {
		YYLog.Error("DeleteOneAdminRole : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//SearchAdminRolesNum 搜索角色数目
//@params	keyword
//@return	int
func SearchAdminRolesNum(keyword string) int {
	var data []orm.Params
	ormer().Raw(`SELECT a.Id FROM admin_role AS a WHERE a.Title LIKE ? `, "%"+keyword+"%").Values(&data)
	return len(data)
}

//SearchAdminRoles 搜索角色
//@params	keyword pagesize offset
//@return	[]orm.Params
func SearchAdminRoles(cateId int, keyword string, pagesize, offset int) []orm.Params {
	var data []orm.Params
	sql := `select Id, CateId, ModelId, Title, Description from admin_role where CateId = ? and Title like concat('%',?,'%') limit ?, ?`
	ormer().Raw(sql, cateId, keyword, offset, pagesize).Values(&data)
	return data
}

//MustCreateOneAdminRole 新建角色
//@params	catId title thumb source description content
//@return	error
func MustCreateOneAdminRole(id, cateId int, title string) error {
	adminrole := new(AdminRole)
	adminrole.Id = id
	adminrole.CateId = cateId
	adminrole.Title = title
	ormer().Insert(adminrole)
	return nil
}
