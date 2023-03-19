package learnReflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_structReflect(t *testing.T) {
	type User struct {
		Name string
		age  int
	}
	testCases := []struct {
		name    string
		entity  any
		wantRes map[string]any
		wantErr error
	}{
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("不支持 nil"),
		},
		{
			name:    "*user nil",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持0值"),
		},
		{
			name: "user",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{"Name": "Tom", "age": 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := structReflect(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
