package model

import (
	"GoStudy/orm"
	"GoStudy/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseModel(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes *Model
		wantErr error
	}{
		{
			name:    "test model",
			entity:  orm.TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:   "test model ptr",
			entity: &orm.TestModel{},
			wantRes: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"Age": {
						ColName: "age",
					},
					"LastName": {
						ColName: "last_name",
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
	r := &orm.registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
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
