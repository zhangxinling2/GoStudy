package orm

import (
	"errors"
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
			wantErr: errors.New("只支持一级指针作为输入"),
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
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, m)
		})
	}
}
