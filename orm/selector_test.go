package orm

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
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
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name:      "eq where",
			builder:   (&Selector[TestModel]{}).Where(C("Tom").Eq(18)),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model` WHERE `Tom` = ?;", Args: []any{18}},
		},
		{
			name:      "and where",
			builder:   (&Selector[TestModel]{}).Where(C("Tom").Eq(18).And(C("Jerry").Lt(19))),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model` WHERE (`Tom` = ?) AND (`Jerry` < ?);", Args: []any{18, 19}},
		},
		{
			name:      "not where",
			builder:   (&Selector[TestModel]{}).Where(Not(C("Tom").Eq(18).And(C("Jerry").Lt(19)))),
			wantQuery: &Query{SQL: "SELECT * FROM `test_model` WHERE  NOT ((`Tom` = ?) AND (`Jerry` < ?));", Args: []any{18, 19}},
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
}
