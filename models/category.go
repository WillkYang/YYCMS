package models

import (
	cnf "YYCMS/conf"
	"YYCMS/utils/YYLog"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/agelinazf/egb"
	"github.com/astaxie/beego/orm"
)

//栏目表
type Category struct {
	//主键
	Id int `orm:"column(Id);pk;auto"`
	//名字
	Title string `orm:"column(Title)"`
	//排序值
	Sort int `orm:"column(Sort);default(0)"`
	//模型Id
	ModelId int `orm:"column(ModelId)"`
	//父栏目Id
	ParentId int `orm:"column(ParentId)"`
}

func (t *Category) TableName() string {
	return "category"
}

//GetOneCategoryById 通过id获取一个栏目信息
//@params	id
//@return	*Category,error
func GetOneCategoryById(id int) (*Category, error) {
	category := new(Category)
	category.Id = id
	if err := ormer().Read(category, "Id"); err != nil {
		YYLog.Error("GetOneCategoryById : " + err.Error())
		return new(Category), fmt.Errorf(ErrInfo[InfoNotExistError])
	}
	YYLog.Debug(category)
	return category, nil
}

func GetOneCategoryRawById(id int) (orm.Params, error) {
	var results []orm.Params
	if _, err := ormer().Raw("select category.*, group_concat(Action) as Actions, group_concat(privilege.Id) as ActionIds from category, privilege where category.Id = privilege.CateId and category.Id = ? GROUP BY category.Id", id).Values(&results); err != nil {
		YYLog.Error("GetOneCategoryById : " + err.Error())
		return nil, fmt.Errorf(ErrInfo[InfoNotExistError])
	}
	if len(results) <= 0 {
		return nil, fmt.Errorf(ErrInfo[InfoNotExistError])
	}
	return results[0], nil
}

//func GetOneCategoryByIdAndRole(cateId, roleId int) (orm.Params, error) {
//	var results []orm.Params
//	if _, err := ormer().Raw(`select category.*, privilege2.Actions from
//(select category.*, group_concat(Action) as Actions, group_concat(privilege.Id) as ActionIds from category, privilege where category.Id = privilege.CateId AND category.Id = ? GROUP BY category.Id) category
//LEFT JOIN
//(SELECT CateId, group_concat(Action) as Actions FROM privilege, role_privilege WHERE privilege.Id = role_privilege.PrivilegeId AND role_privilege.RoleId = ? GROUP BY CateId) privilege2 ON category.Id = privilege2.CateId;`, cateId, roleId).Values(&results); err != nil {
//		YYLog.Error("GetOneCategoryById : " + err.Error())
//		return nil, fmt.Errorf(ErrInfo[InfoNotExistError])
//	}
//
//	if len(results) <= 0 {
//		return nil,fmt.Errorf(ErrInfo[InfoNotExistError])
//	}
//	return results[0], nil
//}

//CategoryHasChildById 判断某个栏目是否含有子栏目
//@params	id
//@return	bool
func CategoryHasChildById(id int) bool {
	num, err := ormer().QueryTable(&Category{}).Filter("ParentId", id).Count()
	if err != nil || num == 0 {
		return false
	}
	return true
}

//GetSubCategorysById 递归的获取一个栏目id的子栏目
//@params	parentid isrecursion(是否递归获取)
//@return	[]interface{}
func GetSubCategorysById(parentid, roleId int, isrecursion bool) []orm.Params {
	categorys := make([]orm.Params, 0)

	//超级管理员特殊处理
	if roleId == cnf.SuperAdminRoleId {
		ormer().Raw(`select category.*, privilege2.Actions as RoleActions from
  		(select category.*, group_concat(Action ORDER BY Action ASC) as Actions, group_concat(privilege.Id ORDER BY privilege.Id ASC) as ActionIds FROM
    	(SELECT * FROM category WHERE category.ParentId = ?) category LEFT JOIN privilege ON category.Id = privilege.CateId GROUP BY category.Id) category
  		LEFT JOIN
  		(SELECT CateId, group_concat(Action ORDER BY Action ASC) as Actions FROM privilege, role_privilege WHERE privilege.Id = role_privilege.PrivilegeId GROUP BY CateId) privilege2 ON category.Id = privilege2.CateId ORDER BY Sort DESC, category.Id ASC;`, parentid).Values(&categorys)
	} else {
		ormer().Raw(`select category.*, privilege2.Actions as RoleActions from
  		(select category.*, group_concat(Action ORDER BY Action ASC) as Actions, group_concat(privilege.Id ORDER BY privilege.Id ASC) as ActionIds FROM
    	(SELECT * FROM category WHERE category.ParentId = ?) category LEFT JOIN privilege ON category.Id = privilege.CateId GROUP BY category.Id) category
  		LEFT JOIN
  		(SELECT CateId, group_concat(Action ORDER BY Action ASC) as Actions FROM privilege, role_privilege WHERE privilege.Id = role_privilege.PrivilegeId AND role_privilege.RoleId = ? GROUP BY CateId) privilege2 ON category.Id = privilege2.CateId ORDER BY Sort DESC, category.Id ASC;`, parentid, roleId).Values(&categorys)
	}
	if parentid == 0 {
		YYLog.Warning(categorys)
	}
	for _, category := range categorys {
		parentId := egb.StringToInt(category["ParentId"].(string))
		cate, _ := GetOneCategoryRawById(parentId)
		if parentId != 0 {
			category["ParentCateName"] = cate["Title"]
		}
		if roleId == cnf.SuperAdminRoleId {
			category["RoleActions"] = category["Actions"]
		}
	}
	data := make([]orm.Params, 0)
	for _, v := range categorys {
		//temp := make(orm.Params, 0)
		temp := v
		if isrecursion {
			cateId := egb.StringToInt(v["Id"].(string))
			if CategoryHasChildById(cateId) {
				temp["children"] = GetSubCategorysByIdWithCache(cateId, roleId, isrecursion)
			} else {
				//temp["children"] = make([]interface{},0)
			}
		}
		data = append(data, temp)
	}
	return data
}

//GetSubCategorysByIdWithCache 递归的获取一个栏目id的子栏目
//@params	parentid isrecursion(是否递归获取)
//@return	[]interface{}
func GetSubCategorysByIdWithCache(parentid, roleId int, isrecursion bool) []orm.Params {
	var results []orm.Params
	oldMask := syscall.Umask(0)
	if err := os.MkdirAll(cnf.CategoryCachePath, os.ModePerm); err != nil {
		syscall.Umask(oldMask)
		YYLog.Error("GetSubCategorysByIdWithCache : " + err.Error())
		return nil
	}
	syscall.Umask(oldMask)

	cachePath := cnf.CategoryCachePath + fmt.Sprintf("GetSubCategorysByIdWithCache:%d.%d.%t.cache", parentid, roleId, isrecursion)
	result, err := egb.FileGetBytes(cachePath)
	if err != nil {
		//缓存文件为空
		results := GetSubCategorysById(parentid, roleId, isrecursion)
		if len(results) == 0 {
			return nil
		}
		str, err := json.Marshal(results)
		if err != nil {
			YYLog.Error("GetSubCategorysByIdWithCache : " + err.Error())
			return nil
		}
		if err := egb.FileSetBytes(cachePath, str); err != nil {
			YYLog.Error("GetSubCategorysByIdWithCache : " + err.Error())
			return nil
		}
		return results
	}

	//缓存文件不为空
	if err := json.Unmarshal(result, &results); err != nil {
		YYLog.Error("GetSubCategorysByIdWithCache : " + err.Error())
		return nil
	}
	return results
}

//CreateOneCategory 创建一个栏目
//@params	role modelid catname pid data
//@return	error
func CreateOneCategory(modelid int, catename string, parentId int) (int, error) {
	YYLog.Debug(parentId)
	if _, err := GetOneCategoryById(parentId); err != nil {
		YYLog.Error("CreateOneCategory GetParentCate error : " + err.Error())
		return 0, fmt.Errorf(ErrInfo[SystemError])
	} else {
		//todo 父栏目必须是目录
		//！=目录
		//pcategory.ModelId = +0
		//return fmt.Errorf(ErrInfo[SystemError])
	}
	YYLog.Debug("hehe")
	category := &Category{
		ModelId:  modelid,
		Title:    catename,
		ParentId: parentId,
	}

	if index, err := ormer().Insert(category); err != nil {
		YYLog.Error("CreateOneCategory : " + err.Error())
		return 0, fmt.Errorf(ErrInfo[SystemError])
	} else {
		return int(index), nil
	}
	//os.RemoveAll(CategoryCachePath)
}

//DeleteOneCategory 删除一个栏目
//@params	id
//@return	error
func DeleteOneCategory(id int) error {
	category, err := GetOneCategoryById(id)
	if err != nil {
		YYLog.Error("DeleteOneCategory : " + err.Error())
		return err
	}
	//category.Status = -1

	if _, err := ormer().Delete(category); err != nil {
		YYLog.Error("DeleteOneCategory : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	//os.RemoveAll(CategoryCachePath)
	return nil
}

//UpdateCategory 更新栏目信息
//@params	id name
//@return	error
func UpdateCategory(id, modelId, parentId int, title string) error {
	category, _ := GetOneCategoryById(id)
	category.Title = title
	category.ModelId = modelId
	category.ParentId = parentId

	if _, err := ormer().Update(category); err != nil {
		YYLog.Error("UpdateCategory : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateCategorySort 更新栏目排序
//@params	id sort
//@return	error
func UpdateCategorySort(id, sort int) error {
	category, _ := GetOneCategoryById(id)
	category.Sort = sort

	if _, err := ormer().Update(category); err != nil {
		YYLog.Error("UpdateCategorySort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//MustCreateOneCategory 创建一个栏目
//@params	id role modelid catname pid data
//@return	error
func MustCreateOneCategory(id, modelid int, catename, description string, parentId int) error {
	category := &Category{
		Id:       id,
		ModelId:  modelid,
		Title:    catename,
		ParentId: parentId,
	}

	if _, err := ormer().Insert(category); err != nil {
		YYLog.Error("CreateOneCategory : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	if model, err := GetOneModelByIdWithCache(modelid); err != nil {
		YYLog.Error(err)
	} else {
		YYLog.Debug(model.Actions)
		CreatePrivileges(id, model.Actions)
	}
	UpdateCache()
	return nil
}
