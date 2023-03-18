package learnReflect

import (
	"fmt"
	"reflect"
)

func structReflect(entity any) (map[string]any, error) {
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	fmt.Println(typ.Name(), "的type为", typ)
	fmt.Println(typ.Name(), "的value为", val)

	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldValue := val.Field(i)
		if fieldTyp.IsExported() {
			res[fieldTyp.Name] = fieldValue.Interface()
			fmt.Println(fieldTyp.Name, "的value为", fieldValue.Interface())
		} else {
			res[fieldTyp.Name] = reflect.Zero(fieldValue.Type()).Interface()
		}
	}
	return res, nil
}
