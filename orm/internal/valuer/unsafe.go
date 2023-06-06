package valuer

import (
	"GoStudy/orm/internal/errs"
	"GoStudy/orm/model"
	"database/sql"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model
	//T的指针
	val any
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(model *model.Model, val any) Value {
	return &unsafeValue{
		model: model,
		val:   val,
	}
}
func (u *unsafeValue) SetColumns(row *sql.Rows) error {
	//拿到列名后肯定要借助model元数据
	cols, err := row.Columns()
	if err != nil {
		return err
	}
	vals := make([]any, 0, len(cols))
	address := reflect.ValueOf(u.val).UnsafePointer()
	for _, c := range cols {
		fd, ok := u.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
		val := reflect.NewAt(fd.Typ, fdAddress)
		vals = append(vals, val.Interface())
	}
	row.Scan(vals)
	return nil
}
