package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

func (t *AdminRole) TableName() string {
	return "admin_role"
}

//GetOneAdminRoleById 获取一个角色
//@params	Id
//@return	*AdminRole
func GetOneAdminRoleById(Id int) (*AdminRole, error) {
	adminrole := new(AdminRole)
	adminrole.Id = Id

	if err := ormer().Read(adminrole,"Id"); err != nil {
		beego.Error("GetOneAdminRoleById : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[DataBaseGetError])
	}
	return adminrole,nil
}

//GetAdminRolesNum 获取角色的数量
//@params	cateId keyword(搜索title的关键词)
//@return	int
func GetAdminRolesNum(catId int, keyword string) int {
	var data []orm.Params
	//todo count
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

//CreateOneAdminRole 新建一篇角色
//@params	catId title thumb source description content
//@return	error
func CreateOneAdminRole(cateId int, title, description, permission string) error {
	adminrole := new(AdminRole)
	adminrole.CateId = cateId
	adminrole.Title = title
	adminrole.Description = description
	adminrole.OperationPermission = permission

	if _,err := ormer().Insert(adminrole); err != nil {
		beego.Error("CreateOneAdminRole : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateAdminRole 更新角色
//@params	id title description
//@return	error
func UpdateAdminRole(id, cateId int, title, description, permission string) error {
	adminrole,err := GetOneAdminRoleById(id)
	if err != nil {
		return err
	}

	adminrole.Title = title
	adminrole.CateId = cateId
	adminrole.OperationPermission = permission

	if description != "" {
		adminrole.Description = description
	}



	if _,err := ormer().Update(adminrole); err != nil {
		beego.Error("UpdateAdminRole : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		return nil
	}
}

//UpdateAdminRoleSort 更新角色排序
//@params	id sort
//@return	error
func UpdateAdminRoleSort(id, sort int) error {
	adminrole,err := GetOneAdminRoleById(id)
	if err != nil {
		return fmt.Errorf(ErrInfo[SystemError])
	}
	adminrole.Sort = sort
	if _,err := ormer().Update(adminrole,"Sort"); err != nil {
		beego.Error("UpdateAdminRoleSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//DeleteOneAdminRole 删除一个角色
//@params	id
//@return	error
func DeleteOneAdminRole(id int) error {
	adminrole,err := GetOneAdminRoleById(id)
	if err != nil {
		return err
	}
	if _,err := ormer().Delete(adminrole); err != nil {
		beego.Error("DeleteOneAdminRole : " + err.Error())
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
func MustCreateOneAdminRole(id, cateId int, title, description, permission string) error {
	adminrole := new(AdminRole)
	adminrole.Id = id
	adminrole.CateId = cateId
	adminrole.Title = title
	adminrole.Description = description
	adminrole.OperationPermission = permission
	ormer().Insert(adminrole)
	return nil
}