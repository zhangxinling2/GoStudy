package orm

import (
	"GoStudy/orm/internal/errs"
	"GoStudy/orm/model"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferName(t *testing.T) {
	t.Run("TableName", func(t *testing.T) {
		res := model.TransferName("TestModel")
		assert.Equal(t, "test_model", res)
	})
}
func TestSelector_Build(t *testing.T) {
	testcases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery string
		wantErr   error
	}{
		{
			name:      "table name",
			builder:   &Selector[TestModel]{},
			wantQuery: "SELECT * FROM `test_model`;",
		},
		{
			name:      "custom table name",
			builder:   (&Selector[TestModel]{}).From("`test_model_test`"),
			wantQuery: "SELECT * FROM `test_model_test`;",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			require.NoError(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q.SQL)
		})
	}
}

//memoryDB 仅用于单元测试
func memoryDB(t *testing.T) *DB {
	db, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	return db
}
func TestSelector_Where(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:      "eq where",
			builder:   (NewSelector[TestModel](db)).Where(C("Age").Eq(18)),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model` WHERE `age` = ?;", Args: []any{18}},
		},
		{
			name:      "and where",
			builder:   (NewSelector[TestModel](db)).Where(C("Age").Eq(18).And(C("Age").Lt(19))),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model` WHERE (`age` = ?) AND (`age` < ?);", Args: []any{18, 19}},
		},
		{
			name:      "not where",
			builder:   (NewSelector[TestModel](db)).Where(Not(C("Age").Eq(18).And(C("Age").Lt(19)))),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model` WHERE  NOT ((`age` = ?) AND (`age` < ?));", Args: []any{18, 19}},
		},
		{
			name:    "error",
			builder: (NewSelector[TestModel](db)).Where(C("Invalid").Eq(18)),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSqlMock(t *testing.T) {
	_, mock, err := sqlmock.New()
	require.NoError(t, err)
	mock.ExpectBegin()
	//NewRows([]string{"id", "name"}).AddRow(1,"Tom") NewRows添加列 AddRow添加列中的数据
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(12, "Tom")
	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)
	mockResult := sqlmock.NewResult(12, 1)
	mock.ExpectExec("UPDATE .*").WillReturnResult(mockResult)
}

func TestGet(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	db, err := OpenDB(mockDB)
	require.NoError(t, err)
	//对应于query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))
	//对应于no rows
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"}).AddRow("1", "Tom", "Jerry", "18")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
	testCases := []struct {
		name    string
		s       *Selector[TestModel]
		wantRes *TestModel
		wantErr error
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(C("XXX").Eq(1)),
			wantErr: errs.NewErrUnknownField("XXX"),
		},
		{
			name:    "Query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errs.ErrNoRows,
		},
		{
			name: "data",
			s:    NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantRes: &TestModel{
				Id:        1,
				FirstName: "Tom",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Jerry"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
