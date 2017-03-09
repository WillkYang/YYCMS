package models

import "YYCMS/utils/YYLog"

type Kvdb struct {
	//Id
	Id    int    `orm:"column(Id);pk;auto"`
	Key   string `orm:"column(Key);size(500)"`
	Value string `orm:"column(Value);size(500)"`
}

func (kvdv *Kvdb) TableName() string {
	return "kvdb"
}

func CreateOneKvdb(key, value string) {
	kv := &Kvdb{Key: key, Value: value}
	if _, _, err := ormer().ReadOrCreate(kv, "Key", "Value"); err != nil {
		YYLog.Error(err)
	}
}

func GetOneKvdbByKey(key string) (value string, result error) {
	kv := &Kvdb{Key: key}
	if result = ormer().Read(kv, "Key"); result != nil {
		value = key
	} else {
		value = kv.Value
	}
	return
}
