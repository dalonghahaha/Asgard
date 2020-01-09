package models

import (
	"reflect"

	"github.com/dalonghahaha/avenger/components/db"
)

func All(list interface{}) error {
	return db.Get(DB_NAME).Find(list).Error
}

func Search(list interface{}, where map[string]interface{}) error {
	return db.Get(DB_NAME).Where(where).Find(list).Error
}

func Count(object interface{}, where map[string]interface{}, count *int) error {
	return db.Get(DB_NAME).Model(object).Where(where).Count(count).Error
}

func AllList(object interface{}, where map[string]interface{}, list interface{}, count *int) error {
	err := db.Get(DB_NAME).Where(where).Find(list).Error
	if err != nil {
		return err
	}
	return db.Get(DB_NAME).Model(object).Where(where).Count(count).Error
}

func PageList(object interface{}, where map[string]interface{}, page int, pageSize int, list interface{}, count *int) error {
	err := db.Get(DB_NAME).Where(where).Limit(pageSize).Offset((page - 1) * pageSize).Find(list).Error
	if err != nil {
		return err
	}
	return db.Get(DB_NAME).Model(object).Where(where).Count(count).Error
}

func Get(id int64, object interface{}) error {
	return db.Get(DB_NAME).Where("id = ? ", id).First(object).Error
}

func Find(where map[string]interface{}, object interface{}) (err error) {
	return db.Get(DB_NAME).Where(where).First(object).Error
}

func Create(object interface{}) (err error) {
	return db.Get(DB_NAME).Create(object).Error
}

func Update(object interface{}) (err error) {
	return db.Get(DB_NAME).Save(object).Error
}

func Delete(object interface{}) (err error) {
	return db.Get(DB_NAME).Delete(object).Error
}

func ModelToMap(inStructPtr interface{}) map[string]interface{} {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}
	info := map[string]interface{}{}
	for i := 0; i < rType.NumField(); i++ {
		key := rType.Field(i).Name
		value := rVal.Field(i).Interface()
		info[key] = value
	}
	return info
}
