package orm

import (
	"GoStudy/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagKeyColumn = "column"
)

//ModelOpt option的变种，带error
type ModelOpt func(m *model) error

// Registry 接口 元数据注册中心的抽象
type Registry interface {
	Get(val any) (*model, error)
	//Register 带Option，因为注册时可能带表名等
	Register(val any, opts ...ModelOpt) (*model, error)
}

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
	//	m,ok:=r.models[typ]
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	return r.Register(val)
}

//GetV1 得到相应的model
//func (r *registry) GetV1(val any) (*model, error) {
//	typ := reflect.TypeOf(val)
//	//判断是否已经缓存了此类型的元数据
//	//	m,ok:=r.models[typ]
//	m, ok := r.models.Load(typ)
//	if !ok {
//		var err error
//		m, err = r.parseModel(typ)
//		if err != nil {
//			return nil, err
//		}
//	}
//	//r.models[typ]=m
//	r.models.Store(typ, m) //也可能同时执行到这里，引发覆盖问题，不过假设元数据解析的结果不变，影响不大
//	return m.(*model), nil
//}

//registry 元数据注册中心，里面维护元数据
type registry struct {
	models sync.Map //使用并发map来使线程安全
	//models map[reflect.Type]*model //结构体名存在同名不同表名的需求，表名则需要元数据，所以最后选择reflect.Type
}

func NewRegistry() *registry {
	return &registry{}
}

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
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}

// parseModel 解析模型
// 声明注册中心后把parseModel作为注册中心的私有方法，希望只用到get
func (r *registry) parseModel(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	//限制输入
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	//获取字段数量
	numField := typ.NumField()
	fields := make(map[string]*field, numField)
	columns := make(map[string]*field, numField)
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
		fd := &field{
			goName:  fdType.Name,
			colName: colName,
			typ:     fdType.Type,
			offset:  fdType.Offset,
		}
		//都需要相同的结果那么久提取出来
		fields[fdType.Name] = fd
		columns[colName] = fd
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}
	if tableName == "" {
		tableName = TransferName(typ.Name())
	}
	return &model{
		tableName: tableName,
		fieldMap:  fields,
		columnMap: columns,
	}, nil
}

//parseModel 不使用注册中心
//func parseModel(entity any) (*model, error) {
//	typ := reflect.TypeOf(entity)
//	//限制输入
//	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
//		return nil, errs.ErrPointerOnly
//	}
//	typ = typ.Elem()
//	//获取字段数量
//	numField := typ.NumField()
//	fieldMap := make(map[string]*field, numField)
//	//解析字段名作为列名
//	for i := 0; i < numField; i++ {
//		fdType := typ.Field(i)
//		fieldMap[fdType.Name] = &field{
//			colName: TransferName(fdType.Name),
//		}
//	}
//	return &model{
//		tableName: TransferName(typ.Name()),
//		fieldMap:    fieldMap,
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
