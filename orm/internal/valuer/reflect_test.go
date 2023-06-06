package valuer

import (
	"GoStudy/orm/model"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestReflectValue_SetColumns(t *testing.T) {
	testCases := []struct {
		name    string
		wantErr error

		//一定是指针
		entity     any
		wantEntity any
		row        func() *sqlmock.Rows
	}{
		{
			name: "setColumn",

			entity: &TestModel{},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
			row: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow("1", "Tom", "18", "Jerry")
				return rows
			},
		},
		{
			// 测试列的不同顺序
			name:   "order",
			entity: &TestModel{},
			row: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "age"})
				rows.AddRow("1", "Jerry", "Tom", "18")
				return rows
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},

		{
			// 测试列的不同顺序
			name:   "partial columns",
			entity: &TestModel{},
			row: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "last_name"})
				rows.AddRow("1", "Jerry")
				return rows
			},
			wantEntity: &TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}
	r := model.NewRegistry()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRows := tc.row()
			//随便写Select 只是为了转化Row
			mock.ExpectQuery("SELECT XX").WillReturnRows(mockRows)
			rows, err := mockDB.Query("SELECT XX")
			require.NoError(t, err)
			rows.Next()

			//得到元数据
			m, err := r.Get(tc.entity)
			require.NoError(t, err)
			if err != nil {
				return
			}
			val := NewReflectValue(m, tc.entity)

			err = val.SetColumns(rows)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}
