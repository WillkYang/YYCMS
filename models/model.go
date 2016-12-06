package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"syscall"
	"os"
	"encoding/json"
	"github.com/agelinazf/egb"
	cnf "YYCMS/conf"
)

func (t *Model) TableName() string {
	return "model"
}

//CreateOneModel 创建一个新的model
//@params	cateId name dbname
//@return	error
func CreateOneModel(id int, title, description, dbname string) error {
	model := new(Model)
	model.Id = id
	model.Title = title
	model.Description = description
	model.DBTableName = dbname
	if _, err := ormer().Insert(model); err != nil {
		beego.Error("CreateOneModel : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//GetOneModelById 通过id获取一个模型信息
//@params	id
//@return	*Model,error
func GetOneModelById(id int) (*Model, error) {
	model := new(Model)
	model.Id = id

	if err := ormer().Read(model,"Id"); err != nil {
		beego.Error("GetOneModelById : " + err.Error())
		return new(Model), fmt.Errorf(ErrInfo[InfoNotExistError])
	}
	return model, nil
}

//GetOneModelByIdWithCache 通过id获取一个模型信息(从缓存取)
//@params	id
//@return	*Model,error
func GetOneModelByIdWithCache(id int) (*Model, error) {
	model := new(Model)
	oldMask := syscall.Umask(0)
	if err := os.MkdirAll(cnf.ModelCachePath, os.ModePerm); err != nil {
		syscall.Umask(oldMask)
		beego.Error("GetOneModelByIdWithCache : " + err.Error())
		return model, fmt.Errorf(ErrInfo[SystemError])
	}
	syscall.Umask(oldMask)
	cachePath := cnf.ModelCachePath + egb.StringFromInt(id) + ".cache"
	result, err := egb.FileGetBytes(cachePath)
	if err != nil {
		//无缓存文件，先生成缓存文件
		//首先从数据库取model
		model, err := GetOneModelById(id)
		if err != nil {
			beego.Error("GetOneModelByIdWithCache : " + err.Error())
			return model, fmt.Errorf(ErrInfo[SystemError])
		}
		//JSON序列化model
		str, err := json.Marshal(model)
		if err != nil {
			beego.Error("GetOneModelByIdWithCache : " + err.Error())
			return model, fmt.Errorf(ErrInfo[SystemError])
		}
		//存入缓存文件
		if err := egb.FileSetBytes(cachePath, str); err != nil {
			beego.Error("GetOneModelByIdWithCache : " + err.Error())
			return model, fmt.Errorf(ErrInfo[SystemError])
		}
		return model, nil
	}
	//存在缓存文件
	if err := json.Unmarshal(result, model); err != nil {
		beego.Error("GetOneModelByIdWithCache : " + err.Error())
		return new(Model), fmt.Errorf(ErrInfo[JSONUnMarsalError])
	}
	return model, nil
}

//GetAllModel 获取全部的模型信息
//@params	nil
//@return	[]Model
func GetAllModel() []Model {
	var models []Model
	ormer().QueryTable(&Model{}).All(&models)
	return models
}

//UpdateModel 更新模型信息
//@params	id cateId title description tableName
//@return	error
func UpdateModel(id, cateId int, title, description, tableName string) error {
	model := &Model{
		Id:id,
		Title:title,
		Description:description,
		DBTableName:tableName,
	}

	if _,err := ormer().Update(model); err != nil {
		beego.Error("UpdateModel : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//GetAllModelWithCache 通过缓存获取全部的模型信息
//@params	nil
//@return	[]Model
func GetAllModelWithCache() []Model {
	var models []Model
	oldMask := syscall.Umask(0)
	if err := os.MkdirAll(cnf.ModelCachePath, os.ModePerm); err != nil {
		syscall.Umask(oldMask)
		beego.Error("GetAllModelWithCache : " + err.Error())
		return nil
	}
	syscall.Umask(oldMask)
	cachePath := cnf.ModelCachePath + "all.cache"
	result, err := egb.FileGetBytes(cachePath)
	if err != nil {
		//缓存文件为空
		models := GetAllModel()
		if len(models) == 0 {
			return nil
		}
		str, err := json.Marshal(models)
		if err != nil {
			beego.Error("GetAllModelWithCache : " + err.Error())
			return nil
		}
		if err := egb.FileSetBytes(cachePath, str); err != nil {
			beego.Error("GetAllModelWithCache : " + err.Error())
			return nil
		}
		return models
	}
	//缓存文件不为空
	if err := json.Unmarshal(result, &models); err != nil {
		beego.Error("GetAllModelWithCache : " + err.Error())
		return nil
	}
	return models
}

//DeleteModel 删除模型信息
//@params	id
//@return	error
func DeleteModel(id int) error {
	model, _ := GetOneModelByIdWithCache(id)

	if _,err := ormer().Delete(model); err != nil {
		beego.Error("UpdateModel : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}

//UpdateModelSort 更新栏目排序
//@params	id sort
//@return	error
func UpdateModelSort(id, sort int) error {

	model, _ := GetOneModelByIdWithCache(id)
	model.Sort = sort

	if _,err := ormer().Update(model); err != nil {
		beego.Error("UpdateModelSort : " + err.Error())
		return fmt.Errorf(ErrInfo[SystemError])
	}
	return nil
}