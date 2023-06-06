package valuer

import (
	"GoStudy/orm/internal/errs"
	"GoStudy/orm/model"
	"database/sql"
	"reflect"
)

type reflectValue struct {
	model *model.Model
	//T的指针
	val any
}

var _ Creator = NewReflectValue

func NewReflectValue(model *model.Model, val any) Value {
	return &reflectValue{
		model: model,
		val:   val,
	}
}
func (r *reflectValue) SetColumns(row *sql.Rows) error {
	//得到列，就可以在之后得到列名
	cols, err := row.Columns()
	if err != nil {
		return err
	}
	vals := make([]any, 0, len(cols))
	//row.Scan(vals)
	valElem := make([]reflect.Value, 0, len(cols))
	for _, c := range cols {
		//没有了selector，那么元数据从哪来？
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			//说明根本没有这个列，查错了
			return errs.NewErrUnknownColumn(c)
		}
		//反射创建了新的实例
		//这里创建的时原本类型的指针 例如fd.typ=int那么val就是int的指针
		val := reflect.New(fd.Typ)
		vals = append(vals, val.Interface())
		valElem = append(valElem, val.Elem())
	}
	//判断是否列过多
	if len(cols) > len(r.model.FieldMap) {
		return errs.ErrMultiCols
	}
	//把值传入vals后再放入t
	err = row.Scan(vals...)
	if err != nil {
		return err
	}
	tpValue := reflect.ValueOf(r.val)
	for i, c := range cols {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			//说明根本没有这个列，查错了
			return errs.NewErrUnknownColumn(c)
		}
		tpValue.Elem().FieldByName(fd.GoName).Set(valElem[i])
	}
	return nil
}
