package orm

import (
	"fmt"
	"log"
	"testing"
)

func TestInsertPerson(t *testing.T) {
	conn, err := ConnectDataBase()
	if err != nil {
		log.Fatal("失败")
	}
	//err = InsertPerson(conn)
	//err = SelectPerson(conn)
	//err = UpdatePerson(conn)
	err = DeletePerson(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
}
