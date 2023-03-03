package orm

import (
	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
)

type PersonalInformation struct {
	Id int `json:"id"		gorm:"column:id"`
	//Id     int     `json:"id"		gorm:"primaryKey;column:id"`
	Name   string  `json:"name"  	gorm:"column:name"`
	Sex    string  `json:"sex		gorm:"column:sex""`
	Tall   float64 `json:"tall"		gorm:"column:tall"`
	Weight int     `json:"weight"	gorm:"column:weight"`
	Age    int     `json:"age"		gorm:"column:age"`
}

func (*PersonalInformation) TableName() string {
	return "person_info"
}
