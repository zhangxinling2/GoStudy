package orm

import (
	"GoStudy/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

//registry 元数据注册中心，里面维护元数据
type registry struct {
	models sync.Map //使用并发map来使线程安全
	//models map[reflect.Type]*model //结构体名存在同名不同表名的需求，表名则需要元数据，所以最后选择reflect.Type
}

func NewRegistry() *registry {
	return &registry{}
}

//get 得到相应的model
func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	//判断是否已经缓存了此类型的元数据
	//	m,ok:=r.models[typ]
	m, ok := r.models.Load(typ)
	if !ok {
		var err error
		m, err = r.parseModel(typ)
		if err != nil {
			return nil, err
		}
	}
	//r.models[typ]=m
	r.models.Store(typ, m) //也可能同时执行到这里，引发覆盖问题，不过假设元数据解析的结果不变，影响不大
	return m.(*model), nil
}

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

// parseModel 解析模型
// 声明注册中心后把parseModel作为注册中心的私有方法，希望只用到get
func (r *registry) parseModel(typ reflect.Type) (*model, error) {
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
			colName: TransferName(fdType.Name),
		}
	}
	return &model{
		tableName: TransferName(typ.Name()),
		fields:    fields,
	}, nil
}

//func parseModel(entity any) (*model, error) {
//	typ := reflect.TypeOf(entity)
//	//限制输入
//	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
//		return nil, errs.ErrPointerOnly
//	}
//	typ = typ.Elem()
//	//获取字段数量
//	numField := typ.NumField()
//	fields := make(map[string]*field, numField)
//	//解析字段名作为列名
//	for i := 0; i < numField; i++ {
//		fdType := typ.Field(i)
//		fields[fdType.Name] = &field{
//			colName: TransferName(fdType.Name),
//		}
//	}
//	return &model{
//		tableName: TransferName(typ.Name()),
//		fields:    fields,
//	}, nil
//}

func TransferName(name string) string {
	var s strings.Builder
	n := []rune(name)
	for i := 0; i < len(name); i++ {
		//判断是否是大写
		if unicode.IsUpper(n[i]) {
			//如果是开头的大写那么只转换成小写，如果不是则在前面加个_
			if i == 0 {
				s.WriteRune(unicode.ToLower(n[i]))
			} else {
				s.WriteByte('_')
				s.WriteRune(unicode.ToLower(n[i]))
			}
		} else {
			s.WriteRune(n[i])
		}
	}
	return s.String()
}
