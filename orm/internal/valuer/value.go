package valuer

import (
	"GoStudy/orm/model"
	"database/sql"
)

//Value 不在函数里面传entity，而是在创建Valuer时传入
//也可以使用在函数里传入entity的设计
type Value interface {
	SetColumns(row *sql.Rows) error
}
type Creator func(model *model.Model, entity any) Value
