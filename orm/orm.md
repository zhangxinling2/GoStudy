##### **2022/3/20**

###### reflect设置值

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

###### reflect输出方法

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

##### **2022/3/21**

###### 元数据解析

元数据很复杂，但是都是一点点加进去的，先从最简定义开始：

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
   //限制输入
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

###### 元数据注册中心

selector中每次都要解析一遍，所以我们可以把它缓存住。

DB在ORM中就相当于HTTPServer在Web框架中的地位，允许用户使用多个DB，DB实例可以单独配置，例如配置元数据中心，DB是天然的隔离和治理单位，所以使用DB来维护元数据。

先定义元数据注册中心registry,里面维护一个map[reflect.Type]*model，之所以要用reflect.Type是因为如果要用结构体名那么会有同结构体名不同表名无法处理，如果要使用表名，我们需要得到元数据但是我们现在在注册元数据，最后选择reflect.Type。把parseModel作为registry的方法把参数改为接受reflect.Type,因为我们希望用户使用get。

```go
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

#### ORM：事务API

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