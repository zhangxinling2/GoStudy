package orm

type Op string

var (
	EQ  Op = "="
	NOT Op = "NOT"
	LT  Op = "<"
	RT  Op = ">"
	AND Op = "AND"
)

//OpString 为了能够直接将Op作为字符串写入
func OpString(op Op) string {
	return string(op)
}

//Expression 可以把Where语句看做Expression Opt Expression   Expression可以是Predicate也可以是Column也可以是arg
//把where语句作为一个二叉树
//所以需要一个标记接口expression来把Predicate，Column，arg标记为expression
type Expression interface {
	expr()
}

//expr 注意不要用指针做接收器
func (p Predicate) expr() {}
func (c Column) expr()    {}

//Value 需要实现expr来标记arg所以arg需要改造成结构体
type Value struct {
	Arg any
}

func (v Value) expr() {}

//这样就全都标记为了Expression

//Predicate 设置来做WHERE的语句
//type Predicate struct {
//	Column Column
//	Opt    Op
//	args   any
//}

//Predicate 完成标记后就要改造Predicate
type Predicate struct {
	left  Expression
	Opt   Op
	right Expression
}

//And 用于Predicate之间的and 即where中的and
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		Opt:   AND,
		right: right,
	}
}

//Not 用来构造where Not
func Not(right Predicate) Predicate {
	return Predicate{
		Opt:   NOT,
		right: right,
	}
}

//Column 提出来单独做个结构体用来做链式调用
type Column struct {
	Name string
}

func C(name string) Column {
	return Column{
		Name: name,
	}
}

//Eq 用来构造And语句
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		Opt:  EQ,
		right: Value{
			Arg: arg,
		},
	}
}

// Lt 用来构造<
func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left: c,
		Opt:  LT,
		right: Value{
			Arg: arg,
		},
	}
}

// Rt 用来构造>
func (c Column) Rt(arg any) Predicate {
	return Predicate{
		left: c,
		Opt:  RT,
		right: Value{
			Arg: arg,
		},
	}
}
