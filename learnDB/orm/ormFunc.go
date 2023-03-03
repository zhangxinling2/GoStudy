package orm

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDataBase() (*gorm.DB, error) {
	conn, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/mysql"))
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	return conn, nil
}
func DeletePerson(conn *gorm.DB) error {
	p := &PersonalInformation{Id: 4}
	resp := conn.Delete(p)
	if err := resp.Error; err != nil {
		fmt.Println("删除失败")
		return err
	}
	return nil
}

// func CloseConnection(conn *gorm.DB) error {

// }
func UpdatePerson(conn *gorm.DB) error {
	// resp := conn.Updates(&PersonalInformation{
	// 	Id:     5,
	// 	Name:   "王五",
	// 	Sex:    "男",
	// 	Tall:   1.8,
	// 	Weight: 91,
	// 	Age:    23,
	// })
	p := &PersonalInformation{
		Id:     5,
		Name:   "王五",
		Sex:    "男",
		Tall:   1.8,
		Weight: 91,
		Age:    23,
	}
	//只修改姓名身高，Model必须有实例
	resp := conn.Model(p).Select("Name", "Tall").Updates(p)
	if err := resp.Error; err != nil {
		fmt.Println("更新失败")
		return err
	}
	return nil
}
func SelectPerson(conn *gorm.DB) error {
	output := make([]*PersonalInformation, 0)
	// resp := conn.Where(&PersonalInformation{Name: "小强"}).Find(&output)
	// if err := resp.Error; err != nil {
	// 	fmt.Println("查询失败")
	// 	return err
	// }
	resp := conn.Where("age < ?", 80).Find(&output)
	if err := resp.Error; err != nil {
		fmt.Println("查询失败")
		return err
	}
	data, _ := json.Marshal(output)
	fmt.Printf("查询的结果为%+v\n", string(data))
	return nil
}
func InsertPerson(conn *gorm.DB) error {
	resp := conn.Create(&PersonalInformation{
		Name:   "王六",
		Sex:    "男",
		Tall:   1.8,
		Weight: 90,
		Age:    23,
	})
	if err := resp.Error; err != nil {
		fmt.Println("创建失败")
		return err
	}
	return nil
}

// func UpdatePerson(conn *gorm.DB)error{
// 	resp:=conn.Up
// }
