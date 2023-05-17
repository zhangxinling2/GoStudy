### ORM

#### reflect

##### reflect设置值

一般使用指针，先使用

```go
vals:=reflect.ValueOf(entity)//得到值信息
vals=vals.Elem()//传递指针时不论是type还是value都要使用elem得到指针指向的东西
val:=vals.FieldByName(field)//根据字段名得到字段值的信息
if !val.CanSet(){			//CanSet()判断能否赋值
		errors.New(fmt.Sprintf("%s不能被设置",field))
	}
val.Set(reflect.ValueOf(newVal))//赋值需要Value类型，所以取传入值的valueOf
```

goland中依赖爆红：

[goland 解决 cannot resolve directory 'xxxx'问题_Lucky小黄人的博客-CSDN博客](https://blog.csdn.net/qq_41767116/article/details/126863153)

****

##### reflect输出方法

方法的接收器有结构体和指针，定义在结构体上的方法使用指针也可以访问。

```go
func IterateFunc(entity any)(map[string]FuncInfo,error){
   typ:=reflect.TypeOf(entity)//得到类型信息
   if typ.Kind()!=reflect.Ptr&&typ.Kind()!=reflect.Struct{//判断是否为结构体或指针
      return nil, errors.New("非法类型")
   }
   numFunc := typ.NumMethod()//得到方法数量
   result:=make(map[string]FuncInfo,numFunc)
   for i := 0; i < numFunc; i++ {
      m:=typ.Method(i)//typ.Method(i)得到Method
      num:=m.Type.NumIn()//.Type得到方法信息 .NumIn()得到输入数量
      fn:=m.Func//.Func是方法的Value
      input:=make([]reflect.Type,0,num)//input是输入参数的类型
      inputValue:=make([]reflect.Value,0,num)//inputValue是输入参数的值
      inputValue=append(inputValue,reflect.ValueOf(entity))//输入的第一个永远是结构体本身，就如同java的this
      for j := 0; j < num; j++ {
         fnInType:=fn.Type().In(j)//In返回的是第j个参数的类型，与m.Type.In()等价
         input= append(input, fnInType)
         if j>0{
            inputValue=append(inputValue,reflect.Zero(fnInType))//输入都用0值即可，用来测试
         }
      }
      outNum:=m.Type.NumOut()
      output:=make([]reflect.Type,0,outNum)
      for j := 0; j < outNum; j++ {
         output= append(output, fn.Type().Out(j))
      }
      resValues:=fn.Call(inputValue)//执行方法，返回的是Value切片
      results:=make([]any,0,len(resValues))
      for _,v:=range resValues{
         results=append(results,v.Interface())
      }
      funcInfo:=FuncInfo{
         Name:m.Name,
         Input: input,
         Output: output,
         Result: results,
      }
      result[m.Name]=funcInfo
   }
   return result,nil
}
type FuncInfo struct {
   Name string
   Input []reflect.Type
   Output []reflect.Type
   Result []any
}
```

#### SELECT起步

##### 核心接口定义

###### 设计一

大而全的核心接口,把所有的需求都放进Orm接口中。

###### 设计二

大一统的Query接口，增删改查都放一起，使用Builder模式。

###### 设计三

直接定义Selector接口，需要构造复杂查询就向里面加方法。使用单一职责的Builder模式。

###### 设计四

只定义Builder模式的终结方法。依旧是Builder模式。

使用泛型做约束：例如SELECT和INSERT语句。

QueryBuilder：作为构建SQL这一个单独步骤的顶级抽象。

```go
// 此接口用于查
type Querier[T any] interface {
   //终结方法
   Get(ctx context.Context) (*T, error)
   GetMulti(ctx context.Context) ([]*T, error)
}
// 此接口用于增删改
type Executor interface {
   Exec(ctx context.Context) (sql.Result, error)
}
// 代表语句
type Query struct {
	SQL  string
	Args []any
}
// Builder接口用来Build语句
type QueryBuilder interface {
	Build() (*Query, error)
}
```

##### SELECTOR定义

先定义出整体再一点点丰富，此结构体实现QueryBuilder接口。

```go
type Selector[T any] struct {
}
```

##### FROM语句

根据mysql规范，先准备构造from。

1.由于Selector本身有泛型参数，可以用泛型的类型名作表名。

2.添加一个From方法，用From方法传入的参数作表名。

使用From的话那么就在Selector中维持一个TableName作表名。

```go
//From 添加表名
func (s *Selector[T]) From(tableName string) *Selector[T] {
   s.TableName = tableName
   return s
}
//那么Build就可以判断TableName是否为空来添加表名
func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	sb := s.sb
	sb.WriteString("SELECT * FROM ")
	//把表名加进去
	if s.TableName == "" {
		//通过反射获取T的名称，需先定义一个T
		var t T
		sb.WriteByte('`')
		//利用反射获得表名
		sb.WriteString(TransferName(reflect.TypeOf(t).Name()))
		sb.WriteByte('`')
	} else {
		sb.WriteString(s.TableName)
	}
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil
}
```

###### 问题：如果用户传入了带db的表名怎么办？

决策：如果用户指定了表名，就直接使用，不会使用反引号；否则使用反引号括起来。随着后面的演化决策可能会更改。

##### WHERE语句

构造完FROM之后，就可以开始着手构造WHERE。

简单设计：Where中直接传入列query和参数args，但这种设计对And，Or，Not难以支持。而且WHERE语句很复杂，这种设计肯定无法满足需求。

###### 设计：WHERE不再接收一个字符串，而是接受一个结构化的Predicate作为输入。

Predicate如何定义？用户如何创建Predicate？

Gorm设计：有一个顶级Expression接口，各种比较符都有一个实现，Not，and，or被认为是一个Expression的集合。

根据WHERE的语法，可以先简单的定义Predicate为(列，操作符，参数)

```go
//这种设计不好链式调用
type Predicate struct {
   Column string
   Opt   string
   Arg   any
}

//把操作符单独定义
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
//对比上面的设计，下面把Column单独拿出来可以链式调用
type Predicate struct {
   Column Column
   Opt   op
   Arg   any
}
type Column struct {
	Name string
}
func C(name string)Column{
    return Column{name:name}
}
//Eq 链式调用例如C("id").Eq(1)
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		c: c,
		Opt:  EQ,
        Arg:arg,
	}
}
```

基本的Predicate设计完后，那么And，Or就是定义在Predicate上的方法，Not左侧缺省，所以不定义在Predicate上。

###### 总结基本的Predicate

Predicate总的来说就是 left op right   And 就为  left and right；Or就为 left or right;Not就为 not right.

###### Expression抽象

WHERE语句就是Predicate op Predicate的二叉树,而到了二叉树底部，就为具体的操作了，如 id LT 12。

那么我们就可以设计一个标记接口，把Predicate和Column都标记为Expression。

```go
//Expression 可以把Where语句看做Expression Opt Expression   Expression可以是Predicate也可以是Column也可以是arg
//把where语句作为一个二叉树
//所以需要一个标记接口expression来把Predicate，Column，arg标记为expression
type Expression interface {
   expr()
}
```

那么Predicate的设计就要更改。

```go
//Predicate 完成标记后就要改造Predicate
type Predicate struct {
   left  Expression
   Opt   Op
   right Expression
}
```

当完成这种设计后，Column的Eq等方法中的Arg就需要新设计一个Arg结构体，让它实现expression标记接口。

```go
//Value 需要实现expr来标记arg所以arg需要改造成结构体
type Value struct {
   Arg any
}

func (v Value) expr() {}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		Opt:  EQ,
		right: Value{
			Arg: arg,
		},
	}
}
```

###### build实现

完成设计后，就要利用build来写WHERE语句了。

```go
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
```

##### 

#### 元数据解析

##### 元数据作用

ORM 框架需要解析模型以获得模型的元数据，这些元数据将被用于构建 SQL、执行校验，以及用于处理结果集。

##### 设计总结

设计总结：

• 模型：对应的表名、索引、主键、关联关系

• 列：列名、Go 类型、数据库类型、是否主键、是否外键

##### 元数据模型定义

###### model与field

元数据很复杂，但是都是一点点加进去的，，其实从我们已经出现的 From 和 Where 来看，我们就需要两个东西：

• 表名

• 列名先从最简定义开始：

```go
type model struct {
   tableName string
   fields    map[string]*field
}
//field 保存字段信息
type field struct {
   colName string
}
```

有了定义，学了反射就可以开始使用反射来解析结构体来获得元数据。

通过反射获得结构体在数据库中的表名和字段在数据库中的列名。

```go
// parseModel 解析模型
func parseModel(entity any) (*model, error) {
   typ := reflect.TypeOf(entity)
   //限制输入只能为一级指针
   if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
      return nil, errs.ErrPointerOnly
   }
   typ = typ.Elem()
   //获取字段数量
   numField := typ.NumField()
   fields := make(map[string]*field, numField)
   //解析字段名作为列名
   for i := 0; i < numField; i++ {
      fdType := typ.Field(i)
      fields[fdType.Name] = &field{
         colName: TransferName(fdType.Name),// TransferName是自己实现的字符串转换
      }
   }
   return &model{
      tableName: TransferName(typ.Name()),
      fields:    fields,
   }, nil
}
```

有了元数据就可以在selector中使用，在Column中校验列名是否在数据库中存在，用户若没有定义表名就可以使用元数据解析的表名。

```go
//处理expression为列的情况
case Column:
   //有了元数据后就可以校验列存不存在
   fd,ok:=s.model.fields[e.Name]
   if !ok{
      return errs.NewErrUnknownField(e.Name)
   }
   s.sb.WriteByte('`')
   s.sb.WriteString(fd.colName)
   s.sb.WriteByte('`')
```

##### 元数据注册中心

如果放在selector中，selector中每次都要解析一遍，所以我们可以把它的解析结果缓存住，那么存在哪，把注册中心交给DB维护。

###### 创建DB

DB在ORM中就相当于HTTPServer在Web框架中的地位，允许用户使用多个DB；DB实例可以单独配置，例如配置元数据中心；DB是天然的隔离和治理单位，所以使用DB来维护元数据。

暂时设计一个 NewDB 的方法，并且留下了 Option模式的口子，为将来留下扩展性的口子。

同样的，虽然目前的实现不会返回 error，但我依然在返回值里面加上了error。

记住，所有的公开方法尽量都要加上 error 作为返回值。

如果你不想返回 error，那么就需要同时提供两个版本的方法。

```go
type DB struct {
   r *registry
}

//DBOption 因为DB有多种，留下一个Option的口子
type DBOption func(*DB)

func NewDB(opts ...DBOption) (*DB, error) {
   db := &DB{
      r: NewRegistry(),
   }
   for _, opt := range opts {
      opt(db)
   }
   return db, nil
}
```

###### 元数据注册中心定义

先定义元数据注册中心registry,里面维护一个map[reflect.Type]*model，之所以要用reflect.Type是因为如果要用结构体名那么会有同结构体名不同表名无法处理;如果要使用表名，我们需要得到元数据但是我们现在在注册元数据;最后选择reflect.Type。把parseModel作为registry的方法把参数改为接受reflect.Type,因为我们希望用户使用get。

```go
type registry struct{
    models map[reflect.Type]*model
}
//get 得到相应的model
func(r *registry)get(val any)(*model,error){
   typ:=reflect.TypeOf(val)
   //判断是否已经缓存了此类型的元数据
   m,ok:=r.models[typ]
   if !ok{
      var err error
      m,err=r.parseModel(typ)
      if err!=nil{
         return nil, err
      }
   }
   r.models[typ]=m
   return m,nil
}
```

###### registry并发安全问题

使用普通map情况下，并发读写场景下肯定崩溃。

并发问题解决的思路有两种：

• 想办法去除掉并发读写的场景，但是可以保留并发读，因为只读永远都是并发安全的。这种就相当于web中的路由树，在服务启动之前就建立好了。

• 用并发工具保护起来

并发：double - check写法或者map使用sync.Map,使用 sync.Map，严格来说有个小问题，即同时解析的时候，会出现覆盖问题。但是我们假设解析的元数据是不会变的，所以问题不大。

###### selector改造

构造 SELECT 语句的时候从 registry 里面拿元数据。

新建selector，按照面向对象的思想，NewSelector 之类的东西应该是定义在 DB 之上的，但是因为泛型的限制，我们只能将db作为参数传入，在selector维护DB。

```go
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
```

##### 自定义表名和列名

有三种方案

+ 标签：直接和模型定义写在一起，非常内聚，但容易写错
+ 接口：直接定义在模型上，可以利用接口实现简单的分库分表功能。但比较隐晦，用户可能不知道实现什么接口。
+ 编程注册：部分享受编译期检查，运行到注册代码，就知道模型是否正确，但学习API难。

在考虑自定义字段名和列名的时候，要考虑模型的两种来源：

• 自己手写：那么三种都可以

• 代码生成：例如 protobuf，这种情况下只能考虑方法二和方法三，要想支持方法一，需要我们自己开发或者修改插件，例如修改 protobuf 编译 Go 语言的插件

###### 自定义列名

两个步骤：1.定义标签的语法。2.解析标签：利用反射提取到完整的标签，然后按照我们的需要进行切割。

```go
// parseTag 解析标签
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
   ormTag := tag.Get("orm")
   if ormTag == "" {
      // 返回一个空的 map，这样调用者就不需要判断 nil 了
      return map[string]string{}, nil
   }
   res := make(map[string]string, 1)
   pairs := strings.Split(ormTag, ",")
   for _, pair := range pairs {
      kv := strings.Split(pair, "=")
      //先限定只有一个tag
      if len(kv) != 2 {
         res[kv[0]] = kv[1]
      }
   }
   return res, nil
}
```

有了parseTag后就要在parseModel中在解析Field时解析tag。

```go
const (
	tagKeyColumn = "column"
)

//解析字段名作为列名
for i := 0; i < numField; i++ {
   fdType := typ.Field(i)
   //解析字段时检测标签
   tags, err := r.parseTag(fdType.Tag)
   if err != nil {
      return nil, err
   }
   colName := tags[tagKeyColumn]
   if colName == "" {
      colName = TransferName(fdType.Name)
   }
   fields[fdType.Name] = &field{
      colName: colName,
   }
}
```

###### 自定义表名

由于标签只用于字段结构体级别（或者说表级别），我们需要额外的手段。

类似于Gorm和Beego，我们需要让用户实现TableName接口来让用户指定表名。

那么在parseModel就要判断用户是否在结构体上实现了TableName接口，如果实现了那么TableName就直接使用。

```go
var tableName string
if tn, ok := val.(TableName); ok {
   tableName = tn.TableName()
}
if tableName == "" {
   tableName = TransferName(typ.Name())
}
return &model{
   tableName: tableName,
   fields:    fields,
}, nil
```

###### 编程方式自定义表名和列名

我们可以允许用户显式地注册模型，同时允许用户在注册的时候自定义一些信息。

为了达成这个目的，我们需要做几件事情：

• 改造设计，添加 Registry 接口

• 为 registry 添加 Register 方法

```go
// Registry 接口 元数据注册中心的抽象
type Registry interface {
   // 
   Get(val any) (*model, error)
   //Register 带Option，因为注册时可能带表名等
   Register(val any, opts ...ModelOpt) (*model, error)
}
```

让registry实现这个抽象。

```go
//ModelOpt option的变种，带error
type ModelOpt func(m *model) error

func (r *registry) Register(val any, opts ...ModelOpt) (*model, error) {
   typ := reflect.TypeOf(val)
   m, err := r.parseModel(val)
   if err != nil {
      return nil, err
   }
   for _, opt := range opts {
      err = opt(m)
      if err != nil {
         return nil, err
      }
   }
   r.models.Store(typ, m)
   return m, nil
}
func (r *registry) Get(val any) (*model, error) {
   typ := reflect.TypeOf(val)
   //判断是否已经缓存了此类型的元数据
   // m,ok:=r.models[typ]
   m, ok := r.models.Load(typ)
   if ok {
      return m.(*model), nil
   }
   return r.Register(val)
}
```

自定义表名：

```go
func ModelWithTableName(tableName string) ModelOpt {
  return func(model *Model)error{
  	model.tableName = tableName
  	return nil
  }
}
```

#### SQL编程

##### 连接数据库

Open：

• driver：也就是驱动的名字，例如 “ mysql” 、“ sqlite3”

• dsn：简单理解就是数据库链接信息

• 常见错误：忘记匿名引入 driver 包

OpenDB：一般用于接入一些自定义的驱动，例如说将分库分表做成一个驱动。

连接sqlite3，mysql：

```go
db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql")
```

##### 增删改查

增改删：

• Exec 或 ExecContext

• 可以用 ExecContext 来控制超时

• 同时检查 error 和 sql.Result

###### 建表

```go
_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
```

###### 插入数据

```go
// 使用 ？ 作为查询的参数的占位符
res, err := db.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)",
   1, "Tom", 18, "Jerry")
affected, err := res.RowsAffected()
log.Println("受影响行数", affected)
lastId, err := res.LastInsertId()
require.NoError(t, err)
log.Println("最后插入 ID", lastId)
```

查询：

• QueryRow 和 QueryRowContext：查询单行数据

• Query 和 QueryContext：查询多行数据

###### 查询单行数据

row必须有一行，如果没有调用row的scan时会返回sql.ErrNoRow

```go
row := db.QueryRowContext(ctx,
   "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
tm := TestModel{}
//查询出来后要用Scan注入
err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
```

###### 查询多行数据

```go
rows, err := db.QueryContext(ctx,
   "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
require.NoError(t, row.Err())
//使用rows.Next来逐行读取
for rows.Next() {
   tm = TestModel{}
   err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
   require.NoError(t, err)
   log.Println(tm)
}
```

要注意参数传递，一般的 SQL 都是使用 ? 作为参数占位符。

不要把参数拼接进去 SQL 本身，容易引起注入。

##### 自定义支持类型

• driver.Valuer：读取，实现该接口的类型可以作为查询参数使用(Go类型到数据库类型）

• sql.Scanner：写入，实现该接口的类型可以作为接收器用于 Scan 方法（数据库类型到Go 类型）

自定义类型一般是实现这两个接口。

如要支持json

```go
type JsonColumn[T any] struct {
   Val T
   // NULL 的问题
   Valid bool
}

func (j JsonColumn[T]) Value() (driver.Value, error) {
   // NULL
   if !j.Valid {
      return nil, nil
   }
   return json.Marshal(j.Val)
}

func (j *JsonColumn[T]) Scan(src any) error {
   //    int64
   //    float64
   //    bool
   //    []byte
   //    string
   //    time.Time
   //    nil - for NULL values
   var bs []byte
   switch data := src.(type) {
   case string:
      // 你可以考虑额外处理空字符串
      bs = []byte(data)
   case []byte:
      // 你也可以考虑额外处理 []byte{}
      bs = data
   case nil:
      // 说明数据库里面存的就是 NULL
      return nil
   default:
      return errors.New("不支持类型")
   }

   err := json.Unmarshal(bs, &j.Val)
   if err == nil {
      j.Valid = true
   }
   return err
}
```

##### sqlmock入门

在单元测试里面我们不希望依赖于真实的数据库，因为数据难以模拟，而且 error 更加难以模拟，所以我们采用 sqlmock 来做单元测试。

###### sqlmock 使用：

• 初始化：返回一个 mockDB，类型是*sql.DB，还有 mock 用于构造模拟的场景；

• 设置 mock：基本上是 ExpectXXXWillXXX，严格依赖于顺序。

```go
func TestSqlMock(t *testing.T) {
   _, mock, err := sqlmock.New()
   require.NoError(t, err)
   mock.ExpectBegin()
   //NewRows([]string{"id", "name"}).AddRow(1,"Tom") NewRows添加列 AddRow添加列中的数据
   mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(12, "Tom")
   mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)
   mockResult := sqlmock.NewResult(12, 1)
   mock.ExpectExec("UPDATE ,*").WillReturnResult(mockResult)
}
```

#### SELECT结果集处理

现在需要让ORM和sql包结合。

##### 发起查询

发起查询就是在selector的Get和GetMulti使用sql.DB发起查询。单行数据使用QueryRowContext,多行数据使用QueryContext。那么在实现时出现的问题就是db从哪来，*T如何转化。

```go
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
   var db sql.DB
   q, err := s.Build()
   if err != nil {
      return nil, err
   }
   row := db.QueryRowContext(ctx, q.SQL, q.Args...)

   return t, nil
}
func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
   var db sql.DB
   q, err := s.Build()
   if err != nil {
      return nil, err
   }
   rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
   for rows.Next() {
   }
   return nil, nil
}
```

###### 那么我们的db从何处来？

当然是与我们ORM层面上的DB绑定到一起，我们可以把DB当做sql.DB的一个封装。这样sql.DB就随着DB保存在selector中。

那么封装完成后如何连接数据库？

那当然是新建一个方法Open来让用户打开数据库。

```go
func Open(driver string, dst string, opts ...DBOption) (*DB,error){
   db, err := sql.Open(driver, dst)
   if err != nil {
      return nil,err
   }
   return &DB{
      r:&registry{},
      db:db,
   }, err
}
```

实际上用户可能自己就创建了sql.DB,我们要允许用户直接使用sql.DB创建我们的DB，所以增加一个OpenDB方法，同时把Open方法改造一下，让它使用OpenDB。且不再需要以前的NewDB。

```go
func Open(driver string, dst string, opts ...DBOption) (*DB, error) {
   db, err := sql.Open(driver, dst)
   if err != nil {
      return nil, err
   }
   return OpenDB(db,opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
   res := &DB{
      r:  &registry{},
      db: db,
   }
   for _, opt := range opts {
      opt(res)
   }
   return res, nil
}
```

那么之前使用sql.DB的地方都要转换成DB下的sql.DB。

##### 结果集构造(即转化为*T)

首先new一个T出来，然后给T中的每个字段赋值。

那么如何赋值？

就需要用到反射和元数据。

如果我们使用Scan取出数据之后用反射来set结构体的值。

```go
var vals []any
row.Scan(&vals[0], &vals[1], &vals[2])
//想办法把vals装进结构体
t := new(T)
tpValue := reflect.ValueOf(t)
//如何把vals与Name对应起来？
//两个问题
//类型要匹配
//顺序要匹配
tpValue.FieldByName("Name").Set(reflect.ValueOf(vals[0]))
```

有两个问题：顺序要匹配；类型要匹配；如果用户更改一下查询参数的顺序怎么办。

发现QueryContext的返回值有一个Column方法可以返回列的名字，而QueryRowContext没有，那么在Get中也要使用QueryContext方法来获取列，因为有了列名，就可以利用反射来设置数据了。

顺序匹配问题：使用Column返回的列名来匹配就可以。

类型匹配问题：如上述代码，我们怎么知道具体的类型，因为需要知道具体的类型才能赋值，query读出来的全是string，而不是使用any。方案：在field中维护一个reflect.Type，把类型记录下来，那么在列名匹配上之后就可以创建相应类型的数据。

值如何Set:使用FieldByName需要知道字段在结构体中的名字，我们在field中维护住它即可。

那么field最后如下：

```go
//field 保存字段信息
type field struct {
   //go中的名字
   goName string
   //列名
   colName string
   //代表字段的类型
   typ reflect.Type
}
```

selector的Get实现如下：

```go
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
   //把语句build出来
   q, err := s.Build()
   if err != nil {
      return nil, err
   }
   //对数据库发起查询
   row, err := s.db.db.QueryContext(ctx, q.SQL, q.Args...)
   if err != nil {
      return nil, err
   }
   //要确认有没有数据
   if !row.Next() {
      return nil, errs.ErrNoRows
   }

   t := new(T)
   //获得元数据
   meta, err := s.db.r.Get(t)
   if err != nil {
      return nil, err
   }
   //拿到列名后肯定要借助model元数据
   cols, err := row.Columns()
   if err != nil {
      return nil, err
   }
   //判断是否列过多
   if len(cols) > len(meta.fields) {
      return nil, errs.ErrMultiCols
   }
   //vals用来存储列中的值 
   vals := make([]any, 0, len(cols))
   //对每列创建实例 
   for _, c := range cols {
      //遍历model的fields字段来找到对应的列 
      for _, fd := range meta.fields {
         if fd.colName == c {
            //反射创建了新的实例
            //这里创建的时原本类型的指针 例如fd.typ=int那么val就是int的指针
            val := reflect.New(fd.typ)
            vals = append(vals, val.Interface())
         }
      }
   }
   //把列的值放入vals中 
   row.Scan(vals...)
   tpValue := reflect.ValueOf(t)
   //把值从vals中取出来给T赋值 
   for i, c := range cols {
      for _, fd := range meta.fields {
         if fd.colName == c {
             //由于new(T)返回的是*T所以要对tpValue取Elem 
            tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
         }   
      }
   }
   return t, nil
}
```

总结：使用reflect构造结果集就是，创建一串盒子，把值从字段中取出来放入盒子，再把值从盒子中取出来放入实例中。

###### 踩坑

使用sqlmock来测试一定要整个一起测试，而不能测试单个用例否则会出现query error。

###### 优化

1.因为我们在构造结果集时使用了循环，而循环效率是比较低的，可以在Model再维持一个columnMap列名到字段定义的映射，那么相应的就要在registry的parseModel中与fieldMap的构造相同，同样构造一个columnMap。有了columnMap在构造结果集时就不用循环去遍历字段看有没有了，只需要使用ok来判断。

2.我们对数据的操作，一会使用interface,一会使用reflect.ValueOf，可以直接把值缓存住就可以吧reflect.ValueOf(vals[i]).Elem()变更为缓存住的值。

```go
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
   q, err := s.Build()
   if err != nil {
      return nil, err
   }
   row, err := s.db.db.QueryContext(ctx, q.SQL, q.Args...)
   if err != nil {
      return nil, err
   }
   //要确认有没有数据
   if !row.Next() {
      return nil, errs.ErrNoRows
   }

   t := new(T)
   //获得元数据
   meta, err := s.db.r.Get(t)
   if err != nil {
      return nil, err
   }
   //拿到列名后肯定要借助model元数据
   cols, err := row.Columns()
   if err != nil {
      return nil, err
   }
   vals := make([]any, 0, len(cols))
   valElem := make([]reflect.Value, 0, len(cols))
   for _, c := range cols {
      fd, ok := s.model.columnMap[c]
      if !ok {
         //说明根本没有这个列，查错了
         return nil, errs.NewErrUnknownColumn(c)
      }
      //反射创建了新的实例
      //这里创建的时原本类型的指针 例如fd.typ=int那么val就是int的指针
      val := reflect.New(fd.typ)
      vals = append(vals, val.Interface())
      valElem = append(valElem, val.Elem())
   }
   //判断是否列过多
   if len(cols) > len(meta.fieldMap) {
      return nil, errs.ErrMultiCols
   }
   err = row.Scan(vals...)
   if err != nil {
      return nil, err
   }
   tpValue := reflect.ValueOf(t)
   for i, c := range cols {
      fd, ok := s.model.columnMap[c]
      if !ok {
         //说明根本没有这个列，查错了
         return nil, errs.NewErrUnknownColumn(c)
      }
      tpValue.Elem().FieldByName(fd.goName).Set(valElem[i])
   }
   return t, nil
}
```

##### 处理结果集API定义

reflect和unsafe设置值比较起来差不多，集成一个API，让用户可以自己选择使用Unsafe还是reflect。

有了valuer在select中把valuer创建一下直接set就可以了。

那么怎么把tp和valuer关联到一起。使用一个函数式的工厂模式。creator

那么creator又从哪来，和

```go
//Value 不在函数里面传entity，而是在创建Value时传入
//也可以使用在函数里传入entity的设计
type Value interface {
   SetColumns(row sql.Rows) error
}
type Creator func(entity any)Value
```

###### reflect

reflectValue实现Value，把selector中使用reflect处理结果集的代码复制到Value中。

之后就需要重构了，先把return更改一下，没有了selector，那么元数据就没有地方获取，那我就在reflectValue中维持model。而model是私有的，在internal中得不到，那么就将model改为Model作为公有的，其中的字段自然也要改为共有的，自然field也要改为共有的。

没有了t,tpValue := reflect.ValueOf(t)自然就不能用了，那么t自然就要在reflectValue维持。

#### Unsafe

要理解Unsafe就要理解go的内存布局。需要掌握

+ 计算地址
+ 计算偏移量
+ 直接操作内存

##### 输出偏移

接收结构体输出偏移

```go
func PrintFieldOffset(entity any) {
   tp := reflect.TypeOf(entity)
   for i := 0; i < tp.NumField(); i++ {
      val := tp.Field(i)
      fmt.Println(val.Name, "offset:", val.Offset)
   }
}
```

可以发现

```go
type User struct {
   name    string
   age     int32
   alias   []byte
   address string
}
```

偏移为：

name offset: 0
age offset: 16
alias offset: 24 为什么不是20而是24？
address offset: 48

因为go每一次访问都是按照字长的倍数来访问的，在32位机器上就是按照4个字节，在64位机器上就是按照8个字节，而age后的4个字节装不下alias所以在下一个字长开始装[]byte。所以我们在age后加字长少于四个字节的结构alias的offset不会变。

##### 使用Unsafe读写字段

注意点：unsafe 操作的是内存，本质上是对象的起始地址。

读：\*(\*T)(ptr)，T 是目标类型，如果类型不知道，只能拿到反射的 Type，那么可以用reflect.NewAt(typ, ptr).Elem()。

写：\*(\*T)(ptr) = T，T 是目标类型。

ptr 是字段偏移量：

ptr = 结构体起始地址 + 字段偏移量

###### 使用Unsafe读取结构体字段

为了契合orm框架UnsafeAccessor维持一个mapstring]fieldMeta和结构体的起始地址,fieldMeta维持字段的类型和偏移量。

```go
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
```

###### 为Accessor实现增删改查

读：Field方法，用来查询field。

reflect.NewAt(fd.typ, unsafe.Pointer(fdAddress)).Elem().Interface()：就是我知道地址在哪，把地址解释成了一个对象，这个对象其实是你真实类型的指针，所以要取个Elem，再用Interface转换成值。

```go
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
```

写：set方法。 

同样的不知道确切类型就用NewAt创建出来后给它Set。

###### Unsafe.Pointer和uintptr的区别

unsafe.Pointer：是go层面上的指针，GC会维护unsafe.Pointer的值

uintptr:直接就只一个数字，代表内存地址，在GC后会变化，我们在fieldMeta中使用这个是因为代表偏移量。使用uint也可以。

##### unsafe应用到结果集处理

首先可以预料到T本身还是需要reflect来构造，里面的字段可以用unsafe来操作。

现在我们就是将它和 Scan 方法集成在一起：

• 计算字段偏移量

• 计算对象起始地址

• 字段真实地址=对象起始地址 + 字段偏移量

• reflect.NewAt 在特定地址创建对象

• 调用 Scan。不再需要set步骤，因为scan本身就是拿到东西填入，有了偏移量就可以直接scan。

现在就是使用unsafe来把vals建好，scan就相当于Set操作。

在计算地址时，需要偏移量，那么只能在field结构体中加上偏移量，在创建时计算。

对象的起始地址就在创建T时计算。

```go
t := new(T)

if err != nil {
   return nil, err
}
//拿到列名后肯定要借助model元数据
cols, err := row.Columns()
if err != nil {
   return nil, err
}
vals := make([]any, 0, len(cols))
address := reflect.ValueOf(t).UnsafePointer()
for _, c := range cols {
   fd, ok := s.model.ColumnMap[c]
   if !ok {
      return nil, errs.NewErrUnknownColumn(c)
   }

   fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
   val := reflect.NewAt(fd.Typ, fdAddress)
   vals = append(vals, val.Interface())
}
row.Scan(vals)
```

现在就是T已经建好了，把数据都填进去。就是把箱子一开始就放在了正确的位置。



#### 事务API

##### Session抽象

核心就是允许用户创建事务，在事物内部进行增删改查，核心有三个API：

·Begin:开启一个事务

·Commit：提交一个事务

·Rollback：回滚一个事务

事务由DB开启，方法定义在DB上，Commit和Rollback由Tx来决定。而将Begin定义在DB上就限制了在一个事务无法开启一个新事务。

![image-20230328073532423](C:\Users\123456\AppData\Roaming\Typora\typora-user-images\image-20230328073532423.png)

Tx的使用：原本Selector接收的是DB做参数，现在使它也可以接收Tx，因为可以在事务中运行(Tx)也可以无事务运行(DB)，那么就需要一个共同的抽象，让DB和Tx来实现。

共同的抽象：session，在ORM语境下，一般代表一个上下文；也可以理解为一个分组机制，在此分组内所有的查询会共享一些基本配置。

Session接口的定义：想要进行抽象，就要把已经被使用的方法提取出来在接口中，在之前代码中，db的方法使用了*sql.DB的QueryContext和ExecContext那么在接口中就定义queryContext和execContext替换掉DB的调用。

core定义:在把session放入NewSelector之后，之前的db.dialect之类的都无法找到，为了得到在DB中我们需要的东西，定义一个core,把增删改查所需要的共同的东西放入core中，重点是DB中持有的，builder中需要什么就放入什么，最后让builder来组合这个core.为了得到core的内容，让DB持有core，在session中新定义一个getCore方法，在Tx中持有创建自己的DB来获得core。builder来使用core所以也要组合core。

##### 事务闭包API

用户传入方法，框架创建事务，事务执行方法然后根据方法的执行情况来判断是提交还是回滚。回滚的条件：出现error或者panic。

在DB上定义DoTx来做事务闭包API,用户传入上下文，业务代码和opts。注意在出错时，需要把err都包装在一起。

```go
func(db *DB)DoTx(ctx context.Context,
	fn func(ctx context.Context,tx *Tx)error,
	opts *sql.TxOptions)(err error){
	tx,err:=db.BeginTx(ctx,opts)
	if err!=nil{
		return err
	}
	panicked:=true
	defer func() {
		if panicked||err!=nil{
			e:=tx.Rollback()
			err=errs.NewErrFailedToRollbackTx(err,e,panicked)
		}else {
			err=tx.Commit()
		}
	}()
	fn(ctx,tx)
	panicked=false
	return err
}
```

由于go没有try-catch机制，虽然DoTx能解决大部分问题,但有时还要自己控制事务，如果事务没有提交就回滚，直接Rollback,返回的错误可以判断。

```go
func(t *Tx)RollbackIfNotCommit()error{
   t.done=true
   err:=t.tx.Rollback()
   //尝试回滚如果事务已经被提交或者回滚那么会返回ErrTxDone
   if err==sql.ErrTxDone{
      return nil
   }
   return err
}
```

###### 事务扩散方案

就是在调用链中，上游方法开启了事务，那么下游方法可以开一个新事务或无事务运行或报错。一般在其他语言中是thread-local，在go中就使用context。核心就是在创建事务时判断context中有没有未完成的事务,tx中定义done判断事务是否完成。

```go
type txKey struct {}
// ctx,tx,err:=db.BeginTxV2()
// doSomething(ctx,tx)
func(db *DB)BeginTxV2(ctx context.Context,opts *sql.TxOptions)(context.Context,*Tx,error){
   val:=ctx.Value(txKey{})
   tx,ok:=val.(*Tx)
   if ok&&!tx.done{
      return ctx,tx,nil
   }
   tx,err:=db.BeginTx(ctx,opts)
   if err!=nil{
      return nil,nil, err
   }
   ctx=context.WithValue(ctx,txKey{},tx)
   return ctx,tx,nil
}
```

#### AOP方案

基本上任何框架都要提供MiddleWare。设计基本照抄web框架的MiddleWare。

##### Beego

Beego设计为侵入式的设计，因为操作没有统一的接口，只有单独的Insert，Read等方法，而我们的ORM框架，对于Select出口只有Get和GetMuti。Insert,Update,Del只有Exec。

##### Gorm

Hook:跟时机强相关。

Create对应于插入，有四个分为两对，BeforeSave,BeforeCreate,AfterSave,AfterCreate。在自己的模型上定义这些方法就会自动执行。

Update也是四个，BeforeSave,BeforeUpdate,AfterSave,AfterUpdate。

Delete有两个，BeforeDelete，AfterDelete。

Query只有一个，AfterFind,没有Before就意味着没办法篡改语句。

##### Aop方案设计

###### 定义

而我们抄web的middleware,做一个函数式的。

```go
type Handler func(ctx context.Context,qc *QueryContext)*QueryResult
type Middleware func(next Handler)Handler
```

```go
//代表上下文
type QueryContext struct {
   // 查询类型，标记增删改查
   Type string

   //代表的是查询本身,大多数情况下需要转化到具体的类型才能篡改查询
   Builder QueryBuilder
   //一般都会暴露出来给用户做高级处理
   Model *model.Model
}
//代表查询结果
type QueryResult struct {
   //Result 在不同查询下类型不同
   //SELECT 可以是*T也可以是[]*T
   //其他就是类型Result
   Result any
   //查询本身出的问题
   Err error
}
```

Middleware用Builder模式

```go
type MiddlewareBuilder struct {
	logFunc func(query string,args []any)
}

func NewMiddlewareBuilder()*MiddlewareBuilder{
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("sql: %s ,args: %v \n",query,args)
		},
	}
}
func (m *MiddlewareBuilder)LogFunc(fn func(query string,args []any))*MiddlewareBuilder  {
	m.logFunc=fn
	return m
}
func(m MiddlewareBuilder)Build()orm.Middleware{
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			q,err:=qc.Builder.Build()
			if err!=nil{
				//要考虑记录下来吗？
				//log.Println("构造 SQL 出错",err)
				return &orm.QueryResult{
					Err: err,
				}
			}
			//log.Printf("sql: %s ,args: %v \n",q.SQL,q.Args)
			//交给用户输出
			m.logFunc(q.SQL,q.Args)
			res:=next(ctx,qc)
			return res
		}
	}
}
```

如何把middleware接入到orm中？

放在db中，而middleware用于所有的增删改查所以放到core中。在DB中再暴露一个Option给middleware。

###### selector改造

有了middleware之后就可以在select中改造，把get的功能放进getHandle中，get用来给getHandle添加middleware,Inserter的改造与selector相同。

```go
func (s *Selector[T])Get(ctx context.Context)(*T,error){
    root:=s.getHandler
	for i:=len(s.mdls)-1;i>=0;i--{
		root=s.mdls[i](root)
	}
    res:= root(ctx,&QueryContext{
        Type:"SELECT",
        Builder:s,
    })
    if res.Result!=nil{
        return res.Result.(*T),res.Err
    }
    return nil,res.Err
}

func (s *Selector[T])getHandler[T any](ctx context.Context,qc *QueryContext) *QueryResult{
   q,err:=s.Build()
   if err!=nil{
      return &QueryResult{
         Err: err,
      }
   }
   //在这里发起查询并处理结果集
   rows,err:=s.sess.queryContext(ctx,q.SQL,q.Args...)
   //这是查询错误，数据库返回的
   if err!=nil{
      return &QueryResult{
         Err: err,
      }
   }
   //将row 转化成*T
   //在这里处理结果集
   if !rows.Next(){
      //要不要返回error
      //返回error,和sql包语义保持一致 sql.ErrNoRows
      //return nil, ErrNoRows
      return &QueryResult{
         Err: ErrNoRows,
      }
   }
   tp:=new(T)
   creator:=c.creator
   val:=creator(c.model,tp)
   err=val.SetColumns(rows)
   return tp,err
}
```

###### middleware增强

我们希望m.Trace.Start(ctx,"","")的span name是select-table_name即类型和表名的结合，所以需要增强一下QueryContext，向其中添加一个model字段以获取表名。

```go
type QueryContext struct {
   // 查询类型，标记增删改查
   Type string
   //代表的是查询本身,大多数情况下需要转化到具体的类型才能篡改查询
   Builder QueryBuilder
   //一般都会暴露出来给用户做高级处理
   Model *model.Model
}
```

那么显然的，在Get构造QueryContext时要加上model,但是s中的model直到build时才会被赋值，那么我们可以考虑：

提前给model赋值，在Get中加上

```go
    var err error
    s.model,err=s.r.Get(new(T))
    if err!=nil{
       return nil, err
    }
```

或者专门给一个middleware给添加model。

#### 集成测试

orm框架：确保和数据库交互返回结果正确。

##### TestSuite

要使用不同的数据库，使用TestSuite:

1.它提供了一种分组机制效果

2.隔离：套件之间允许独立运行。

3.生命周期回调(钩子)：允许在套件前后执行一些动作。

4.参数控制：可用不同参数多次运行同一套件。

使用：

在一个结构体中集成suite.Suite，SetupSuite用来初始化db。

```go
type Suite struct {
   suite.Suite
   driver string
   dsn string
   db *orm.DB
}
// SetupSuite 所有suite执行前的钩子
func (s *Suite)SetupSuite(){
	db,err:=orm.Open(s.driver, s.dsn)
	require.NoError(s.T(), err)
	db.Wait()
	s.db=db
}
```

Wait是用来等待数据库启动并连接的，在DB上新增Wait方法：

```go
//Wait 主动等待数据库启动
func (d *DB) Wait()error{
   err:=d.db.Ping()
   //循环等待 
   for err==driver.ErrBadConn{
      log.Println("等待数据库启动...")
      err = d.db.Ping()
      time.Sleep(time.Second)
   }
   return err
}
```

想要进行什么测试，就用相应的结构体来集成Suite。

结构体上仍然可以再次定义SetupSuite，Suite中的用于在所有实例前运行，特定结构体的用于运行在特定实例上。

```go
type SelectSuite struct {
   Suite
}
//测试的进入方法
func TestMySQLTest(t *testing.T){
   suite.Run(t, &SelectSuite{
      Suite{
         driver: "mysql",
         dsn: "root:root@tcp(localhost:13306)/integration_test",
      },
   })
}
//TearDownSuite 所有都跑完清数据
func (s *InsertSuite)TearDownSuite(){
   orm.RawQuery[test.SimpleStruct](s.db,"TRUNCATE TABLE `simple_struct`").Exec(context.Background())
}
//Select的SetupSuite用来插入数据
func (s *SelectSuite)SetupSuite()  {
   s.Suite.SetupSuite()
   res:=orm.NewInserter[test.SimpleStruct](s.db).Values(test.NewSimpleStruct(100)).Exec(context.Background())
   require.NoError(s.T(),res.Err())
}

func(s *SelectSuite)TestSelect(){
   testCases:=[]struct{
      name string
      s *orm.Selector[test.SimpleStruct]

      wantRes *test.SimpleStruct
      wantErr error
   }{
      {
         name:"get data",
         s:orm.NewSelector[test.SimpleStruct](s.db).Where(orm.C("Id").Eq(100)),//数据从SetupSuite中插入
         wantRes: test.NewSimpleStruct(100),
      },
      {
         name:"no row",
         s:orm.NewSelector[test.SimpleStruct](s.db).Where(orm.C("Id").Eq(200)),//数据从SetupSuite中插入
         wantErr: orm.ErrNoRows,
      },
   }

   for _,tc:=range testCases{
       //t替换为s.T() 
      s.T().Run(tc.name, func(t *testing.T) {
         ctx,cancel:=context.WithTimeout(context.Background(),time.Second*10)
         defer cancel()
         res,err:=tc.s.Get(ctx)
         assert.Equal(t, tc.wantErr,err)
         if err!=nil{
            return
         }
         assert.Equal(t, tc.wantRes,res)
      })
   }
}
```

##### 数据的准备

select时，把数据准备好，测试全部完成后删除。

insert时，数据单独准备，每个用例完成后删除。

##### 标签

在头部添加//go:build tag    那么在go test -tags=tag ./...,如果不加tag那么就不会测试有标签的测试。

#### 原生查询

我们设计的orm select显然不能完全满足select的语法，那么就要给用户提供绕过orm框架写查询语句的机制，而结果集可通过orm框架也可通过sql.DB来封装。

##### 设计

显然的，我们需要原生的支持增删改查，那么就需要实现Querier(支持select)和Executor(支持增删改)和QueryBuilder(Build创建语句)。

```go
type RawQuerier[T any] struct {
   core
   sess Session
   //存储语句和参数
   sql string
   args []any
}
//需要一个构造函数来创建rawQuery，所以再实现QueryBuilder
func (r RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL: r.sql,
		Args: r.args,
	},nil
}
//需要一个构造函数来创建rawQuerier
func RawQuery[T any](sess Session,query string,args...any)*RawQuerier[T]{
	c:=sess.getCore()
	return &RawQuerier[T]{
		sql: query,
		args: args,
		sess: sess,
		core:c,
	}
}
func (i RawQuerier[T]) Exec(ctx context.Context) Result {
	if i.model==nil{
		var err error
		i.model,err=i.r.Get(new(T))
		if err!=nil{
			return Result{
				err: err,
			}
		}
	}

	res:=exec(ctx,i.sess,i.core,&QueryContext{
		Type: "RAW",
		Builder: i,
		Model: i.model,
	})

	var sqlRes sql.Result
	if res.Result!=nil{
		sqlRes = res.Result.(sql.Result)
	}
	return Result{
		err: res.Err,
		res:sqlRes,
	}
}
//实现的Querier的Get跟Selector中的差不多相同，但是getHandler定义在selector上，所以尝试把getHandler拆出来.
func (s RawQuerier[T]) Get(ctx context.Context) (*T, error) {
		var err error
    	//r从哪来？在RawQuerier中组合一个core
		s.model,err=s.r.Get(new(T))
		if err!=nil{
			return nil,err
		}
    //session从哪来？在RawQuerier中维护一个session
	res:=get[T](ctx,s.sess,s.core,&QueryContext{
			Type: "RAW",
			Builder: s,
        	//model从哪来？只能在get之前获取
			Model: s.model,
		})
	if res.Result!=nil{
		return res.Result.(*T),res.Err
	}
	return nil,res.Err
}

func (r RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("implement me")
}

//由于get getHandler exec execHandler是通用的，所以将这些方法放入core中。
//由于在selector和RawQuerier中的Get逻辑非常相似，所以把get也提取出来
func get[T any](ctx context.Context,sess Session,c core,qc  *QueryContext)*QueryResult{
    //不符合方法签名
	//var root Handler = getHandler[T](ctx,s.sess,s.core,&QueryContext{
	//	Type: "RAW",
	//	Builder: s,
	//	Model: s.model,
	//})
    //为了使用getHandler，所以我们get也需要传入sess,c,qc
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx,sess,c,qc)
	}
	for i:=len(c.mdls)-1;i>=0;i--{
		root=c.mdls[i](root)
	}
	//return root(ctx,&QueryContext{
	//	Type: "RAW",
	//	Builder: builder,
	//	//问题在于s.model在Build时才会赋值，1.在Get初始化s.model 2.专门设置一个middleware来设置model
	//	Model: c.model,
	//})
	return root(ctx,qc)
}
//拆除来的getHandler缺少了selector中的sess和core那么我们就给它传入sess和core，build就是qc中的Builder
func getHandler[T any](ctx context.Context,sess Session,c core,qc *QueryContext) *QueryResult{
	q,err:=qc.Builder.Build()
	if err!=nil{
		return &QueryResult{
			Err: err,
		}
	}
	//在这里发起查询并处理结果集
	rows,err:=sess.queryContext(ctx,q.SQL,q.Args...)
	//这是查询错误，数据库返回的
	if err!=nil{
		return &QueryResult{
			Err: err,
		}
	}
	if !rows.Next(){
		return &QueryResult{
			Err: ErrNoRows,
		}
	}
	tp:=new(T)
	creator:=c.creator
	val:=creator(c.model,tp)
	err=val.SetColumns(rows)
	return &QueryResult{
		Err: err,
		Result: tp,
	}
}
```

#### Join查询

Join查询有点像我们的 Expression，就是可以查询套查询无限套下去。

##### Join语法

JOIN 语法有两种形态：

+ JOIN ... ON
+ JOIN ... USING：USING 后面使用的是列名

JOIN本身有：

+ INNER JOIN、 JOIN
+ LEFT JOIN、RIGHT JOIN

SELECT

#### protobuf魔改

[protobuf-go](https://github.com/protocolbuffers/protobuf-go):下载了源码后，在proto-gen-go的main中找到了生成的函数GenerateFile，而我们要魔改的是生成的struct，在GenerateFile函数所在文件中找到了genMessageField，把

```go
tags := structTags{
   {"protobuf", fieldProtobufTagValue(field)},
   {"json", fieldJSONTagValue(field)},
}
```

改为：

```go
tags := structTags{
   {"protobuf", fieldProtobufTagValue(field)},
   {"json", fieldJSONTagValue(field)},
   {"orm", fieldORMTagValue(field)},
}
```

而我们自己定义的fieldORMTagValue的实现为：

```go
func fieldORMTagValue(field *protogen.Field) string {
   c:=field.Comments.Trailing.String()//Trailing就是跟在後面的，Leading是放在上面的
   c=strings.TrimSpace(c)
    //语法为 //@orm:column=xx 
   if strings.HasPrefix(c,"//@orm"){
      return c[7:]
   }
   return ""
}
```

总的流程为：

1.clone原来的protobuf-go代码库

2.修改protobuf-go代码

3.安装修改后的go插件  在protoc-gen-go文件夹下执行go install .

4.执行protoc命令