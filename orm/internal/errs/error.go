package errs

import (
	"errors"
	"fmt"
)

//一般不带参数直接声明

var (
	ErrPointerOnly = errors.New("orm:只支持一级指针作为输入")
	ErrNoRows      = errors.New("orm:没有数据")
	ErrMultiCols   = errors.New("orm:列过多")
)

//带参数的可声明为函数

func NewErrUnsupportedExpression(expr any) error {
	return errors.New(fmt.Sprintf("orm:不支持的表达式 %v", expr))
}

func NewErrUnknownField(expr any) error {
	return errors.New(fmt.Sprintf("orm:未知字段 %v", expr))
}
func NewErrInvalidTagContent(pair any) error {
	return errors.New(fmt.Sprintf("orm:无效tag %v", pair))
}
func NewErrUnknownColumn(expr any) error {
	return errors.New(fmt.Sprintf("orm:未知列 %v", expr))
}
