package model

import "reflect"

//Model 元数据用来构建sql和处理结果集
//元数据设计 一个模型：用来存储表的信息 一个列用来存储列的信息
//在Selector中引入元数据，最直接的需求就是校验字段正确与否
type Model struct {
	TableName string
	//字段名到字段的映射
	FieldMap map[string]*Field
	//列名到字段的映射
	ColumnMap map[string]*Field
}

//Field 保存字段信息
type Field struct {
	//go中的名字
	GoName string
	//列名
	ColName string
	//代表字段的类型
	Typ reflect.Type

	Offset uintptr
}

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}
