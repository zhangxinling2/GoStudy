package learnReflect

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type User struct {
	Name string
	age  int
}

func NewUser(name string, age int) User {
	return User{Name: name, age: age}
}
func NewUserPtr(name string, age int) *User {
	return &User{Name: name, age: age}
}
func (u User) GetName() string {
	return u.Name
}
func (u *User) GetAge() int {
	return u.age
}
func (u *User) ChangeName(name string) {
	u.Name = name
}
func Test_structReflect(t *testing.T) {

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

func TestSetColumn(t *testing.T) {

	tests := []struct {
		name       string
		entity     any
		field      string
		newVal     any
		wantErr    error
		wantEntity any
	}{
		{
			name:       "CanSet",
			entity:     &User{Name: "Tom", age: 18},
			field:      "Name",
			newVal:     "Jerry",
			wantEntity: &User{Name: "Jerry", age: 18},
		},
		{
			name:    "Can not Set",
			entity:  &User{Name: "Tom", age: 18},
			field:   "age",
			newVal:  19,
			wantErr: errors.New(fmt.Sprintf("age不能被设置")),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SetColumn(tc.entity, tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantEntity, tc.entity)
		})
	}
}

func TestIterateFunc(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		args    []any
		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			name:   "struct",
			entity: NewUser("Tom", 18),
			wantRes: map[string]FuncInfo{
				"GetName": {
					Name:   "GetName",
					Input:  []reflect.Type{reflect.TypeOf(User{})},
					Output: []reflect.Type{reflect.TypeOf("")},
					Result: []any{"Tom"},
				},
			},
		},
		{
			name:   "point",
			entity: NewUserPtr("Tom", 18),
			wantRes: map[string]FuncInfo{
				"ChangeName": {
					Name:   "ChangeName",
					Input:  []reflect.Type{reflect.TypeOf(&User{}), reflect.TypeOf("")},
					Output: []reflect.Type{},
					Result: []any{},
				},
				"GetAge": {
					Name:   "GetAge",
					Input:  []reflect.Type{reflect.TypeOf(&User{})},
					Output: []reflect.Type{reflect.TypeOf(0)},
					Result: []any{18},
				},
				"GetName": {
					Name:   "GetName",
					Input:  []reflect.Type{reflect.TypeOf(&User{})},
					Output: []reflect.Type{reflect.TypeOf("")},
					Result: []any{""},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFunc(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
