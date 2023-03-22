package orm

type DB struct {
	r *registry
}

//DBOption 因为DB有多种，留下一个Option的口子
type DBOption func(*DB)

func NewDB(opts ...DBOption) (*DB, error) {
	db := &DB{
		r: &registry{},
	}
	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}
