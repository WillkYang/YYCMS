package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func (t *Category) TableName() string {
	return "category"
}

//GetOneCategoryById 通过id获取一个栏目信息
//@params	id
//@return	*Category,error
func GetOneCategoryById(id int) (*Category, error) {
	category := new(Category)
	category.Id = id

	if err := ormer().Read(category,"Id"); err != nil {
		beego.Error("GetOneCategoryById : " + err.Error())
		return new(Category), fmt.Errorf(ErrInfo[InfoNotExistError])
	}
	beego.Debug(category)
	return category, nil
}

//CategoryHasChildById 判断某个栏目是否含有子栏目
//@params	id
//@return	bool
func CategoryHasChildById(id int) bool {
	num, err := ormer().QueryTable(&Category{}).Exclude("Status",-1).Filter("ParentId", id).Count()
	if err != nil || num == 0 {
		return false
	}
	return true
}

//GetSubCategorysById 递归的获取一个栏目id的子栏目
//@params	parentid isrecursion(是否递归获取)
//@return	[]interface{}
func GetSubCategorysById(parentid int, isrecursion bool) []interface{} {
	categorys := make([]orm.Params, 0)
	ormer().QueryTable(&Category{}).Filter("ParentId", parentid).OrderBy("-Sort").Values(&categorys)
	for i := 0; i < len(categorys); i++ {
		if categorys[i]["ParentId"].(int64) != 0 {
			cate, _ := GetOneCategoryById(int(categorys[i]["ParentId"].(int64)))
			categorys[i]["ParentCateName"] = (*cate).Title
		}
	}
	data := make([]interface{}, 0)
	for _, v := range categorys {
		temp := make(orm.Params, 0)
		temp = v
		if isrecursion {
			if CategoryHasChildById(int(v["Id"].(int64))) {
				temp["children"] = GetSubCategorysById(int(v["Id"].(int64)), isrecursion)
			} else {
				//temp["children"] = make([]interface{},0)
			}
		}
		data = append(data, temp)
	}
	return data
}

//CreateOneCategory 创建一个栏目
//@params	role modelid catname pid data
//@return	error
func CreateOneCategory(modelid int, catename, description string, parentId int) error {
	beego.Debug(parentId)
	if _, err := GetOneCategoryById(parentId); err != nil {
		beego.Error("CreateOneCategory GetParentCate error : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	} else {
		//todo 父栏目必须是目录
		//！=目录
		//pcategory.ModelId = +0
		//return fmt.Errorf(ErrInfo[SystemError])
	}
	beego.Debug("hehe")
	category := &Category{
		ModelId:modelid,
		Title:catename,
		Description:description,
		ParentId:parentId,
	}

	if _, err := ormer().Insert(category); err != nil {
		beego.Error("CreateOneCategory : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	//os.RemoveAll(CategoryCachePath)
	return nil
}

//DeleteOneCategory 删除一个栏目
//@params	id
//@return	error
func DeleteOneCategory(id int) error {
	category, err := GetOneCategoryById(id)
	if err != nil {
		beego.Error("DeleteOneCategory : " + err.Error())
		return err
	}
	category.Status = -1

	if _,err := ormer().Update(category,"Status"); err != nil {
		beego.Error("DeleteOneCategory : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	//os.RemoveAll(CategoryCachePath)
	return nil
}

//UpdateCategory 更新栏目信息
//@params	id name
//@return	error
func UpdateCategory(id, modelId, parentId int, title, description string) error {
	category, _ := GetOneCategoryById(id)
	category.Title = title
	category.ModelId = modelId
	category.ParentId = parentId
	category.Description = description

	if _,err := ormer().Update(category); err != nil {
		beego.Error("UpdateCategory : " + err.Error())
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

	if _,err := ormer().Update(category); err != nil {
		beego.Error("UpdateCategorySort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}


//MustCreateOneCategory 创建一个栏目
//@params	id role modelid catname pid data
//@return	error
func MustCreateOneCategory(id, modelid int, catename, description string, parentId int) error {
	category := &Category{
		Id:id,
		ModelId:modelid,
		Title:catename,
		Description:description,
		ParentId:parentId,
	}

	if _, err := ormer().Insert(category); err != nil {
		beego.Error("CreateOneCategory : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	//os.RemoveAll(CategoryCachePath)
	return nil
}