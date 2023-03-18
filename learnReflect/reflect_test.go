package learnReflect

import (
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
		user    User
		wantRes map[string]any
		wantErr error
	}{
		{
			name: "user",
			user: User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{"Name": "Tom", "age": 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := structReflect(User{
				Name: "Tom",
				age:  18,
			})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
