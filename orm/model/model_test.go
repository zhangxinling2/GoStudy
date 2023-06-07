package model

import (
	"GoStudy/orm/internal/errs"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestParseModel(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes *Model
		wantErr error
	}{
		{
			name:    "test model",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "test model ptr",
			entity: &TestModel{},
			wantRes: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						GoName:  "Id",
						ColName: "id",
						Typ:     reflect.TypeOf(int64(0)),
						Offset:  0,
					},
					"FirstName": {
						GoName:  "FirstName",
						ColName: "first_name",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					"Age": {
						GoName:  "Age",
						ColName: "age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
					"LastName": {
						GoName:  "LastName",
						ColName: "last_name",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
				},
			},
		},
		{
			name:   "test tag",
			entity: &TestTag{},
			wantRes: &Model{
				TableName: "TestTag---",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "idTest",
					},
				},
			},
		},
	}
	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			columnMap := make(map[string]*Field)
			for _, f := range tc.wantRes.FieldMap {
				columnMap[f.ColName] = f
			}
			tc.wantRes.ColumnMap = columnMap
			assert.Equal(t, tc.wantRes, m)
		})
	}
}

type TestTag struct {
	Id int `orm:"column=idTest"`
}

func (t TestTag) TableName() string {
	return "TestTag---"
}
