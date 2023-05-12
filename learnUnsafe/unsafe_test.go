package learnUnsafe

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"reflect"
	"testing"
)

type User struct {
	name    string
	age     int32
	ageV1   int32
	alias   []byte
	address string
}

func TestPrintFieldOffset(t *testing.T) {
	PrintFieldOffset(User{
		name:    "",
		age:     0,
		ageV1:   0,
		alias:   nil,
		address: "",
	})
}
func PrintFieldOffset(entity any) {
	tp := reflect.TypeOf(entity)
	for i := 0; i < tp.NumField(); i++ {
		val := tp.Field(i)
		fmt.Println(val.Name, "offset:", val.Offset)
	}
}
func TestUnsafeAccessor_Field(t *testing.T) {
	user := &User{
		name:    "tom",
		age:     0,
		ageV1:   0,
		alias:   nil,
		address: "",
	}
	accessor := NewUnsafeAccessor(user)
	val, err := accessor.Field("age")
	require.NoError(t, err)
	log.Println(val)
}
