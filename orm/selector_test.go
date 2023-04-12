package orm

import (
	"GoStudy/internal/errs"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferName(t *testing.T) {
	t.Run("TableName", func(t *testing.T) {
		res := TransferName("TestModel")
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

func TestSelector_Where(t *testing.T) {
	db, err := NewDB()
	require.NoError(t, err)
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
	mock.ExpectExec("UPDATE ,*").WillReturnResult(mockResult)
}
