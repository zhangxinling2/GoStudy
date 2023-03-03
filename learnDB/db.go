package learnDB

import (
	"GoStudy/code"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func connectMySql() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		fmt.Println(err)
		return
	}
	InsertInfo(db)
	QueryAllInfo(db)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

}

func InsertInfo(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO personinfo(name,sex,tall,weight,age) VALUES('%s','%s',%f,%d,%d)", "阿珍", "女", 1.76, 73, 23))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func QueryAllInfo(db *sql.DB) {
	res, err := db.Query("SELECT name,age FROM personinfo")
	if err != nil {
		fmt.Println(err)
		return
	}

	persons := []*code.PersonInfo{}
	for res.Next() {
		p := &code.PersonInfo{}
		res.Scan(&p.Name, &p.Age) //Scan中参数顺序一定与select中的严格匹配
		persons = append(persons, p)
	}
	data, _ := json.Marshal(persons)
	fmt.Println(string(data))
}
