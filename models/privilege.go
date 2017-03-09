package models

import (
	"YYCMS/utils"
	"YYCMS/utils/YYLog"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Action string

const (
	Add     Action = "add"
	Delete  Action = "delete"
	Profile Action = "profile"
	Sort    Action = "sort"
	Update  Action = "update"
	Publish Action = "publish"
)

//获取全部的可执行操作
func getAllAction() []Action {
	operations := make([]Action, 0)
	operations = append(operations, Profile, Add, Update, Delete, Publish, Sort)
	return operations
}

////栏目权限表
type Privilege struct {
	Id int `orm:"column(Id)"`
	//栏目Id
	CateId int `orm:"column(CateId)"`
	//方法
	Action Action `orm:"column(Action)"`
}

func (t *Privilege) TableName() string {
	return "privilege"
}

//批量创建栏目权限
func CreatePrivileges(cateId int, privileges string) error {
	for _, action := range utils.StringToStringArray(privileges) {
		if CategoryOperationIsExist(cateId, Action(action)) {
			continue
		}
		privilege := new(Privilege)
		privilege.CateId = cateId
		privilege.Action = Action(action)
		if _, err := orm.NewOrm().Insert(privilege); err != nil {
			YYLog.Error("CreateOnePrivilege : " + err.Error())
			return fmt.Errorf(ErrInfo[DataBaseUpdateError])
		}
	}
	return nil
}

//批量删除栏目权限   privileges: add,delete,update...
func DeletePrivileges(cateId int, privileges string) error {
	for _, action := range utils.StringToStringArray(privileges) {
		if _, err := ormer().QueryTable(&Privilege{}).Filter("CateId", cateId).Filter("Action", Action(action)).Delete(); err != nil {
			return fmt.Errorf(ErrInfo[DataBaseUpdateError])
		}
	}
	return nil
}

//批量删除栏目权限
func DeleteAllPrivileges(cateId int) error {
	if _, err := ormer().QueryTable(&Privilege{}).Filter("CateId", cateId).Delete(); err != nil {
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//CreateOnePrivilege	创建一个新的权限
//@params	cateId action
//@return	error
func CreateOnePrivilege(cateId int, action Action) error {
	if CategoryOperationIsExist(cateId, action) {
		return nil
	}
	privilege := new(Privilege)
	privilege.CateId = cateId
	privilege.Action = action
	if _, err := orm.NewOrm().Insert(privilege); err != nil {
		YYLog.Error("CreateOnePrivilege : " + err.Error())
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//CreateCategoryOperation 为某个栏目添加全部操作权限
//@params	cateId
//@return	error
func CreateCategoryOperation(cateId int) error {
	operations := getAllAction()
	for _, v := range operations {
		if err := CreateOnePrivilege(cateId, v); err != nil {
			return err
		}
	}
	return nil
}

//DeleteCateAllAction 删除栏目的全部权限
//@params	cateId
//@return	error
func DeleteCateAllAction(cateId int) error {
	_, err := orm.NewOrm().QueryTable(&Privilege{}).Filter("CateId", cateId).Delete()
	if err != nil {
		YYLog.Error("DeleteCategoryOperation : " + err.Error())
		return fmt.Errorf(ErrInfo[DataBaseUpdateError])
	}
	return nil
}

//CategoryOperationIsExist
//@params	cateId action
//@return	bool
func CategoryOperationIsExist(cateId int, action Action) bool {
	return orm.NewOrm().QueryTable(&Privilege{}).Filter("CateId", cateId).Filter("Action", action).Exist()
}

//ReadCategoryPrivileges 获取栏目对应的全部权限id
//@params	cateId
//@return	[]int(read,create,update,delete,publish)
func ReadCategoryPrivileges(cateId int) []Action {
	var actionArr []Action
	var privilege Privilege
	for _, action := range getAllAction() {
		if err := orm.NewOrm().QueryTable(&Privilege{}).Filter("CateId", cateId).Filter("Action", action).One(&privilege); err == nil {
			actionArr = append(actionArr, action)
		}
	}
	return actionArr
}

//ReadPrivilege 通过cateId和operation获取权限信息
//@params	cateId operation
//@return	*Privilege
func ReadOnePrivilege(cateId int, action Action) *Privilege {
	privilege := new(Privilege)
	privilege.CateId = cateId
	privilege.Action = action
	if err := orm.NewOrm().Read(privilege, "CateId", "Action"); err != nil {
		YYLog.Error("ReadPrivilege : " + err.Error())
		return new(Privilege)
	}
	return privilege
}
