package orm

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

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}
