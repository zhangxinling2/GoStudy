package main

import (
	"GoStudy/dataStore/fatRank"
	"GoStudy/httpServer/httpPratice/frinterface"
	"fmt"
	"gorm.io/gorm"
	"log"
)

var _ frinterface.ServeInterface = &dbRank{}

type dbRank struct {
	conn      *gorm.DB
	embedRank frinterface.ServeInterface //为了使用已经实现的内存rank
}

func (d *dbRank) Init() error {
	output := make([]*fatRank.PersonalInformation, 0)
	resp := d.conn.Find(&output)
	if resp.Error != nil {
		fmt.Println("查找失败：", resp.Error)
		return resp.Error
	}
	for _, item := range output {
		if _, err := d.embedRank.UpdatePersonInformation(item); err != nil {
			log.Fatalf("初始化%s时失败：%v", item.Name, err)
		}
	}
	return nil
}

//NewDbRank 创建基于DB的RANK同时内嵌内存Rank
func NewDbRank(conn *gorm.DB, embedRank frinterface.ServeInterface) *dbRank {
	return &dbRank{
		conn:      conn,
		embedRank: embedRank,
	}
}

//RegisterPersonInformation 向DB中添加信息，同时向内存中添加
func (d *dbRank) RegisterPersonInformation(pi *fatRank.PersonalInformation) error {
	resp := d.conn.Create(pi)
	if err := resp.Error; err != nil {
		log.Println("注册失败")
		return err
	}
	log.Println("注册成功")
	//因为内嵌了内存的排行榜，数据库变动时内存排行榜也要进行相应的变动
	_ = d.embedRank.RegisterPersonInformation(pi)
	return nil
}

func (d *dbRank) UpdatePersonInformation(pi *fatRank.PersonalInformation) (*fatRank.PersonalInformationFatRate, error) {
	resp := d.conn.Updates(pi)
	if err := resp.Error; err != nil {
		log.Println("更新失败")
		return nil, err
	}
	log.Println("更新成功")
	bmi, err := BMI(pi.Weight, pi.Tall)
	if err != nil {
		log.Println("计算BMI出错")
		return nil, err
	}
	_, _ = d.embedRank.UpdatePersonInformation(pi)
	return &fatRank.PersonalInformationFatRate{
		Name:    pi.Name,
		Fatrate: CalcFatRate(float64(bmi), int(pi.Age), pi.Sex),
	}, nil
}

//GetFatrate 不从数据库取数据
func (d *dbRank) GetFatrate(name string) (*fatRank.PersonRank, error) {
	return d.embedRank.GetFatrate(name)
}

func (d *dbRank) GetTop() ([]*fatRank.PersonRank, error) {
	return d.embedRank.GetTop()
}
