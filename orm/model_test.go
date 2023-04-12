package orm

import (
	"GoStudy/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseModel(t *testing.T) {
	testCases := []struct {
		name    string
		entity  any
		wantRes *model
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
			wantRes: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
			},
		},
		{
			name:   "test tag",
			entity: &TestTag{},
			wantRes: &model{
				tableName: "TestTag---",
				fields: map[string]*field{
					"Id": {
						colName: "idTest",
					},
				},
			},
		},
	}
	r := &registry{}
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
