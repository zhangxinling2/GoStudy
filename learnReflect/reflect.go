package learnReflect

import (
	"errors"
	"fmt"
	"reflect"
)

func SetColumn(entity any, field string, newVal any) error {
	if entity == nil {
		return errors.New("不支持 nil")
	}
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return errors.New("非法类型")
	}
	vals := reflect.ValueOf(entity)
	vals = vals.Elem()
	val := vals.FieldByName(field)
	if !val.CanSet() {
		return errors.New(fmt.Sprintf("%s不能被设置", field))
	}
	val.Set(reflect.ValueOf(newVal))
	return nil
}
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

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	typ := reflect.TypeOf(entity)                                  //得到类型信息
	if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct { //判断是否为结构体或指针
		return nil, errors.New("非法类型")
	}
	numFunc := typ.NumMethod() //得到方法数量
	result := make(map[string]FuncInfo, numFunc)
	for i := 0; i < numFunc; i++ {

		m := typ.Method(i)                                       //typ.Method(i)得到Method
		num := m.Type.NumIn()                                    //.Type得到方法信息 .NumIn()得到输入数量
		fn := m.Func                                             //.Func是方法的Value
		input := make([]reflect.Type, 0, num)                    //input是输入参数的类型
		inputValue := make([]reflect.Value, 0, num)              //inputValue是输入参数的值
		inputValue = append(inputValue, reflect.ValueOf(entity)) //输入的第一个永远是结构体本身，就如同java的this
		for j := 0; j < num; j++ {

			fnInType := fn.Type().In(j) //In返回的是第j个参数的类型
			input = append(input, fnInType)
			if j > 0 {
				inputValue = append(inputValue, reflect.Zero(fnInType)) //输入都用0值即可，用来测试
			}
		}
		outNum := m.Type.NumOut()
		output := make([]reflect.Type, 0, outNum)
		for j := 0; j < outNum; j++ {
			output = append(output, fn.Type().Out(j))
		}

		resValues := fn.Call(inputValue)
		results := make([]any, 0, len(resValues))
		for _, v := range resValues {
			results = append(results, v.Interface())
		}
		funcInfo := FuncInfo{
			Name:   m.Name,
			Input:  input,
			Output: output,
			Result: results,
		}
		result[m.Name] = funcInfo
	}
	return result, nil
}

type FuncInfo struct {
	Name   string
	Input  []reflect.Type
	Output []reflect.Type
	Result []any
}

func IterateSlice(entity any) error {
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	switch typ.Kind() {
	case reflect.Slice:
		res := make([]any, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			res = append(res, val.Index(i).Interface())
		}
		return nil
	case reflect.Map:
		key := make([]any, 0, val.Len())
		value := make([]any, 0, val.Len())
		for _, k := range val.MapKeys() {
			key = append(key, k.Interface())
			value = append(value, k.MapIndex(k).Interface())
		}
		return nil
	}
	return nil
}
