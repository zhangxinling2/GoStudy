package orm

import (
	"GoStudy/internal/errs"
	"context"
	"database/sql"
	"strings"
)

//Selector TableName是为了表名
//sb 是为了拼接字符串
//为了使用where 定义一个Predicate切片
//为了接收参数设置一个any数组
type Selector[T any] struct {
	TableName string
	//为了在多个函数中拼接字符串将string加入struct中
	sb    *strings.Builder
	where []Predicate
	args  []any

	model *model //有了元数据后，selector中就可以加入元数据,Build中的数据就可以使用元数据中的数据

	db *DB //设计出DB后，在selector中加入DB
}

//NewSelector 新建selector实例，自定义db
func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: &strings.Builder{},
		db: db,
	}
}
func (s *Selector[T]) Build() (*Query, error) {
	//有了元数据后就可以改造selector
	var t T //用来解析
	var err error
	//定义元数据注册中心后selector使用它的get方法即可
	//s.model,err=parseModel(&t)
	s.model, err = s.db.r.Get(&t)
	if err != nil {
		return nil, err
	}

	s.sb = &strings.Builder{}
	sb := s.sb
	sb.WriteString("SELECT * FROM ")
	//把表名加进去
	if s.TableName == "" {
		//通过反射获取T的名称，需先定义一个T
		//var t T 获取元数据后可注释掉
		sb.WriteByte('`')
		//利用反射获得表名
		//sb.WriteString(TransferName(reflect.TypeOf(t).Name()))
		//使用元数据解析的表名
		sb.WriteString(s.model.tableName)
		sb.WriteByte('`')
	} else {
		sb.WriteString(s.TableName)
	}
	//添加完表名之后，继续拼接where条件
	//判断是否有条件
	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		//先取第一个用于和后面的where组合
		pw := s.where[0]
		for i := 1; i < len(s.where); i++ {
			//合并predicate
			pw = pw.And(s.where[i])
		}
		//where有三种情况需要处理Predicate，Column和Value
		//构建where，直接使用不能断言，需要一个函数以expression形式接受之后断言
		//switch typ := pw.(type) {
		//}
		if err = s.buildExpression(pw); err != nil {
			return nil, err
		}
	}
	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil
}

//From 添加表名
func (s *Selector[T]) From(tableName string) *Selector[T] {
	s.TableName = tableName
	return s
}

//Where 添加Predicate
func (s *Selector[T]) Where(predicates ...Predicate) *Selector[T] {
	s.where = predicates
	return s
}
func (s *Selector[T]) buildExpression(expr Expression) error {
	switch e := expr.(type) {
	case nil:
		return nil
	//处理expression为列的情况
	case Column:
		//有了元数据后就可以校验列存不存在
		fd, ok := s.model.fields[e.Name]
		if !ok {
			return errs.NewErrUnknownField(e.Name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
	case Value:
		//需要先初始化一下arg切片
		if s.args == nil {
			s.args = make([]any, 0, 8)
		}
		s.sb.WriteByte('?')
		s.args = append(s.args, e.Arg)
	//最后来处理Predicate情况
	case Predicate:
		//由于Predicate是二叉树形态，所以可以用递归来构建
		//构建左边Predicate
		//断言left是否是Predicate,如果是Predicate则证明仍然是一个式子需要用括号扩起来
		_, ok := e.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(e.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(OpString(e.Opt))
		s.sb.WriteByte(' ')
		//构建右侧Predicate
		//断言right是否是Predicate，如果是Predicate则证明仍然是一个式子需要用括号扩起来
		_, ok = e.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(e.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var db sql.DB
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	row := db.QueryRowContext(ctx, q.SQL, q.Args...)
	var t *T
	err = row.Scan(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	var db sql.DB
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

	}
	return nil, nil
}
