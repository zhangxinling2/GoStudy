package orm

import (
	"errors"
	"reflect"
	"strings"
	"unicode"
)

//元数据用来构建sql和处理结果集
//元数据设计 一个模型：用来存储表的信息 一个列用来存储列的信息
//在Selector中引入元数据，最直接的需求就是校验字段正确与否
type model struct {
	tableName string
	fields    map[string]*field
}

//field 保存字段信息
type field struct {
	colName string
}

// parseModel 解析模型
func parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	//限制输入
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errors.New("只支持一级指针作为输入")
	}
	typ = typ.Elem()
	//获取字段数量
	numField := typ.NumField()
	fields := make(map[string]*field, numField)
	//解析字段名作为列名
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		fields[fdType.Name] = &field{
			colName: TransferName(fdType.Name),
		}
	}
	return &model{
		tableName: TransferName(typ.Name()),
		fields:    fields,
	}, nil
}

func TransferName(name string) string {
	var s strings.Builder
	n := []rune(name)
	for i := 0; i < len(name); i++ {
		//判断是否是大写
		if unicode.IsUpper(n[i]) {
			//如果是开头的大写那么只转换成小写，如果不是则在前面加个_
			if i == 0 {
				s.WriteRune(unicode.ToLower(n[i]))
			} else {
				s.WriteByte('_')
				s.WriteRune(unicode.ToLower(n[i]))
			}
		} else {
			s.WriteRune(n[i])
		}
	}
	return s.String()
}
