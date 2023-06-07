package orm

import (
	"GoStudy/orm/internal/valuer"
	"GoStudy/orm/model"
	"database/sql"
)

type DB struct {
	creator valuer.Creator
	r       model.Registry
	db      *sql.DB
}

//DBOption 因为DB有多种，留下一个Option的口子
type DBOption func(*DB)

func DBWithUnsafe() DBOption {
	return func(db *DB) {
		db.creator = valuer.NewUnsafeValue
	}
}
func DBWithReflect() DBOption {
	return func(db *DB) {
		db.creator = valuer.NewReflectValue
	}
}
func Open(driver string, dst string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dst)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:       model.NewRegistry(),
		db:      db,
		creator: valuer.NewReflectValue,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}
