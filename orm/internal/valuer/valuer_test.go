package valuer

import (
	"GoStudy/orm/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkValue(b *testing.B) {
	fn := func(b *testing.B, creator Creator) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(b, err)
		defer mockDB.Close()
		rows := sqlmock.NewRows([]string{"id", "last_name", "first_name", "age"})

		//需要跑N次，准备N行数据
		for i := 0; i < b.N; i++ {
			rows.AddRow("1", "Jerry", "Tom", "18")
		}
		mock.ExpectQuery("SELECT XX").WillReturnRows(rows)
		mockRows, err := mockDB.Query("SELECT XX")
		r := model.NewRegistry()
		m, err := r.Get(&TestModel{})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mockRows.Next()
			val := creator(m, &TestModel{})
			err = val.SetColumns(mockRows)
		}
		require.NoError(b, err)
	}
	b.Run("reflect", func(b *testing.B) {
		fn(b, NewReflectValue)
	})
	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValue)
	})
}
