package learnReflect

import (
	"errors"
	"fmt"
	"reflect"
)

func structReflect(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, errors.New("不支持0值")
	}
	fmt.Println(typ.Name(), "的type为", typ)
	fmt.Println(typ.Name(), "的value为", val)
	if typ.Kind() != reflect.Ptr {
		return nil, errors.New("只支持结构体的指针类型")
	}
	typ = typ.Elem()
	val = val.Elem()
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("只支持结构体的指针类型")
	}
	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldValue := val.Field(i)
		if fieldTyp.IsExported() {
			res[fieldTyp.Name] = fieldValue.Interface()
			fmt.Println("字段"+fieldTyp.Name, "的value为", fieldValue.Interface())
		} else {
			res[fieldTyp.Name] = reflect.Zero(fieldValue.Type()).Interface()
		}
	}
	return res, nil
}
