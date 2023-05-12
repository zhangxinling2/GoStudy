package learnUnsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	fields  map[string]fieldMeta
	address unsafe.Pointer
}
type fieldMeta struct {
	typ    reflect.Type
	Offset uintptr
}

//NewUnsafeAccessor entity是结构体指针
func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numField := typ.NumField()
	fields := make(map[string]fieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = fieldMeta{
			Offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		fields: fields,
		//不直接用UnsafeAddr，因为它对应的地址不是稳定的，Gc之后地址会变化
		//UnsafePointer会帮助维持指针
		address: val.UnsafePointer(),
	}
}
func (u *UnsafeAccessor) Field(field string) (any, error) {
	fd, ok := u.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}
	//这样不能加，需要对unsafePointer进行转化
	//fdAddress:=u.address+fd.Offset
	fdAddress := uintptr(u.address) + fd.Offset
	//如果知道类型那么就
	//用(*)(*int)(unsafe.Pointer(fdAddress))来读
	//不知道类型
	return reflect.NewAt(fd.typ, unsafe.Pointer(fdAddress)).Elem().Interface(), nil
}
