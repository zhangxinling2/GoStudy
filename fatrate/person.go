package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	pFilePath = "./testdata/person.json"
)

type Person struct {
	Name   string  `json:"name"`
	Sex    string  `json:"sex"`
	Tall   float64 `json:"tall"`
	Weight float64 `json:"weight"`
	Age    int     `json:"age"`

	Bmi     float64 `json:"bmi"`
	FatRate float64 `json:"fat_rate"`
}

func (p *Person) calcBmi() error {
	bmi, err := BMI(p.Weight, p.Tall)
	if err != nil {
		log.Printf("error when calculating BMP for Person[%s]: %v", p.Name, err)
		return err
	}
	p.Bmi = bmi
	return nil
}

func (p *Person) calcFatRate() {
	p.FatRate = CalcFatRate(p.Bmi, p.Age, p.Sex)
}

func (p Person) RegisterByJson() error {
	f, err := os.OpenFile(pFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return errors.New("文件打开失败")
	}
	defer f.Close()

	data, err := json.Marshal(p)
	fmt.Println(p)
	if err != nil {
		fmt.Println(err)
		return errors.New("解析失败")
	}
	data = append(data, '\n')
	if err != nil {
		return errors.New("json失败")
	}
	f.Write(data)
	return nil
}
