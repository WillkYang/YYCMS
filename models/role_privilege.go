package models

import (
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

type RolePrivilege struct {
	//Id
	Id int `orm:"column(Id);pk;auto"`
	//角色Id
	RoleId int `orm:"column(RoleId)"`
	//权限Id
	PrivilegeId int `orm:"column(PrivilegeId)"`
}

func (m *RolePrivilege) TableName() string {
	return "role_privilege"
}

//CreateOneRolePrivilege 创建一个角色和权限的对应信息
//@params	roleId privilegeId
//@return	error
func CreateOneRolePrivilege(roleId, privilegeId int) error {
	if RolePrivilegeIsExist(roleId, privilegeId) {
		return nil
	}
	rolePrivilege := new(RolePrivilege)
	rolePrivilege.RoleId = roleId
	rolePrivilege.PrivilegeId = privilegeId
	if _, err := orm.NewOrm().Insert(rolePrivilege); err != nil {
		YYLog.Error("CreateOneRolePrivilege :" + err.Error())
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//DeleteOneRolePrivilege 删除一条角色和权限的中间表信息
//@params	roleId privilegeId
//@return	error
func DeleteOneRolePrivilege(roleId, privilegeId int) error {
	if !RolePrivilegeIsExist(roleId, privilegeId) {
		return nil
	}
	rolePrivilege := new(RolePrivilege)
	rolePrivilege.RoleId = roleId
	rolePrivilege.PrivilegeId = privilegeId
	if _, err := orm.NewOrm().Delete(rolePrivilege, "RoleId", "PrivilegeId"); err != nil {
		YYLog.Error("DeleteOneRolePrivilege : " + err.Error())
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//批量删除角色权限
func DeleteAllRolePrivileges(roleId int) error {
	if _, err := ormer().QueryTable(&RolePrivilege{}).Filter("RoleId", roleId).Delete(); err != nil {
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//CreateRolePrivileges 为一个角色添加权限
//@params	roleId privilegeId
//@return	error
func CreateRolePrivileges(roleId int, privilegeId string) error {
	if orm.NewOrm().QueryTable(RolePrivilege{}).Filter("RoleId", roleId).Exist() {
		return nil
	}
	privilegeIds, err := egb.StringToIntArray(privilegeId)
	if err != nil {
		YYLog.Error("CreateRolePrivileges : " + err.Error())
		return fmt.Errorf(ErrInfo[DataStructError])
	}
	for _, v := range privilegeIds {
		if err := CreateOneRolePrivilege(roleId, v); err != nil {
			return err
		}
	}
	return nil
}

//
////UpdateRolePrivileges 更新角色的权限
////@params	roleId privilegeId
////@return	error
//func UpdateRolePrivileges(roleId int, privilegeId string) error {
//	//以前的权限id数据
//	oldRolePrivileges := ReadRolePrivileges(roleId)
//	//现在的权限id数据
//	nowRolePrivileges, err := egb.StringToIntArray(privilegeId)
//	if err != nil {
//		YYLog.Error("UpdateRolePrivileges : " + err.Error())
//		return fmt.Errorf(ErrInfo[DataStructError])
//	}
//	oldHasButNowNot, nowHasButNowNot := FindIntArrayDifference(oldRolePrivileges, nowRolePrivileges)
//	//删除oldHasButNowNot
//	for _, v := range oldHasButNowNot {
//		if err := DeleteOneRolePrivilege(roleId, v); err != nil {
//			return err
//		}
//	}
//	//添加nowHasButNowNot
//	for _, v := range nowHasButNowNot {
//		if err := CreateOneRolePrivilege(roleId, v); err != nil {
//			return err
//		}
//	}
//	return nil
//}

//ReadRolePrivileges 获取某个角色的权限id数组
//@params	roleId
//@return	[]int
func ReadRolePrivileges(roleId int) []int {
	var mids []RolePrivilege
	orm.NewOrm().QueryTable(&RolePrivilege{}).Filter("RoleId", roleId).All(&mids)
	intArr := make([]int, 0)
	for _, v := range mids {
		intArr = append(intArr, v.PrivilegeId)
	}
	return intArr
}

//DeleteRolePrivileges 删除角色的权限中间表信息
//@params	roleId
//@return	error
func DeleteRolePrivileges(roleId int) error {
	_, err := orm.NewOrm().QueryTable(&RolePrivilege{}).Filter("RoleId", roleId).Delete()
	if err != nil {
		YYLog.Error("DeleteRolePrivileges : " + err.Error())
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//IfRolePrivilegeExist 角色权限是否存在
//@params	roleId privilegeId
//@return	bool
func RolePrivilegeIsExist(roleId, privilegeId int) bool {
	return orm.NewOrm().QueryTable(&RolePrivilege{}).Filter("RoleId", roleId).Filter("PrivilegeId", privilegeId).Exist()
}
