package orm

import "database/sql"

type DB struct {
	r  *registry
	db *sql.DB
}

//DBOption 因为DB有多种，留下一个Option的口子
type DBOption func(*DB)

func Open(driver string, dst string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dst)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
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
