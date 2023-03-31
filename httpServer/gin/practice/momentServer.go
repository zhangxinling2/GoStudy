package main

import (
	"GoStudy/dataStore/fatRank"
	"GoStudy/httpServer/gin/practice/moments"
	"gorm.io/gorm"
	"log"
	"time"
)

var _ moments.Moments = &momentServer{}

type momentServer struct {
	conn   *gorm.DB
	person fatRank.PersonalInformation
}

func NewMomentServer(conn *gorm.DB, person fatRank.PersonalInformation) *momentServer {
	return &momentServer{
		conn:   conn,
		person: person,
	}
}
func (m momentServer) ReleaseMoment(text string) (fatRank.PersonalMoment, error) {
	bmi, err := BMI(m.person.Weight, m.person.Tall)
	if err != nil {
		log.Println("计算BMI出错")
		return fatRank.PersonalMoment{}, err
	}
	fatrate := CalcFatRate(float64(bmi), int(m.person.Age), m.person.Sex)
	mi := fatRank.PersonalMoment{
		PersonId:    m.person.Id,
		CreatedTime: time.Now(),
		Content:     text,
		Fatrate:     float32(fatrate),
		Visible:     true,
	}
	res := m.conn.Create(mi)
	if res.Error != nil {
		log.Println("插入数据出错")
		return fatRank.PersonalMoment{}, err
	}
	return mi, nil
}

func (m momentServer) DeleteMoment(id int64) error {
	mi := fatRank.PersonalMoment{}
	res := m.conn.Model(&mi).Update("visible", false)
	if res.Error != nil {
		log.Println("删除数据出错")
		return res.Error
	}
	return nil
}

func (m momentServer) GetAllMoment() (fatRank.PersonalMomentList, error) {
	var list fatRank.PersonalMomentList
	res := m.conn.Model(&fatRank.PersonalMoment{}).Where("visible", true).Scan(&list.Items)
	if res.Error != nil {
		log.Println("查找数据出错")
		return fatRank.PersonalMomentList{}, res.Error
	}
	return list, nil
}
