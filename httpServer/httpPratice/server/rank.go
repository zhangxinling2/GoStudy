package main

import (
	"GoStudy/dataStore/fatRank"
	"GoStudy/dataStore/types"
	"GoStudy/httpServer/httpPratice/frinterface"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
)

var _ frinterface.ServeInterface = &FatRateRank{}

type RandItem struct {
	Name    string
	Sex     string
	FatRate float64
}

type FatRateRank struct {
	items    []RandItem
	itemLock sync.Mutex
}

func BMI(weightKG, heightM float32) (bmi float32, err error) {
	if weightKG < 0 {
		err = fmt.Errorf("weight cannot be negative")
		return
	}
	if heightM < 0 {
		err = fmt.Errorf("height cannot be negative")
		return
	}
	if heightM == 0 {
		err = fmt.Errorf("height cannot be 0")
		return
	}
	bmi = weightKG / (heightM * heightM)
	return
}
func CalcFatRate(bmi float64, age int, sex string) (fatRate float64) {
	sexWeight := 0
	if sex == "男" {
		sexWeight = 1
	} else {
		sexWeight = 0
	}
	fatRate = (1.2*bmi + getAgeWeight(age)*float64(age) - 5.4 - 10.8*float64(sexWeight)) / 100
	return
}

func getAgeWeight(age int) (ageWeight float64) {
	ageWeight = 0.23
	if age >= 30 && age <= 40 {
		ageWeight = 0.22
	}
	return
}
func (r *FatRateRank) RegisterPersonInformation(pi *fatRank.PersonalInformation) error {
	gobmi, err := BMI(pi.Weight, pi.Tall)
	if err != nil {
		log.Println("计算BMI失败")
		return err
	}
	fr := CalcFatRate(float64(gobmi), int(pi.Age), pi.Sex)
	r.inputRecord(pi.Name, pi.Sex, fr)
	return nil
}

func (r *FatRateRank) UpdatePersonInformation(pi *fatRank.PersonalInformation) (*types.PersonalInformationFatRate, error) {
	gobmi, err := BMI(pi.Weight, pi.Tall)
	if err != nil {
		log.Println("计算BMI失败")
		return nil, err
	}
	fr := CalcFatRate(float64(gobmi), int(pi.Age), pi.Sex)
	r.inputRecord(pi.Name, pi.Sex, fr)
	return &types.PersonalInformationFatRate{
		Name:    pi.Name,
		Fatrate: fr,
	}, nil
}

func (r *FatRateRank) GetFatrate(name string) (*types.PersonRank, error) {
	rankId, sex, fr := r.getRank(name)
	Fr := types.PersonalInformationFatRate{
		Name:    name,
		Fatrate: fr,
	}
	return &types.PersonRank{
		RankNumber: rankId,
		Sex:        sex,
		Fr:         Fr,
	}, nil
}

func (r *FatRateRank) GetTop() ([]*types.PersonRank, error) {
	return r.getRankTop(), nil
}

func NewFatRateRank() *FatRateRank {
	return &FatRateRank{
		items: make([]RandItem, 0, 100),
	}
}

func (r *FatRateRank) inputRecord(name, sex string, fatRate ...float64) {
	r.itemLock.Lock()
	defer r.itemLock.Unlock()
	minFatRate := math.MaxFloat64
	for _, item := range fatRate {
		if minFatRate > item {
			minFatRate = item
		}
	}

	found := false
	for i, item := range r.items {
		if item.Name == name {
			if item.FatRate >= minFatRate {
				item.FatRate = minFatRate
			}
			r.items[i] = item
			found = true
			break
		}
	}
	if !found {
		r.items = append(r.items, RandItem{
			Name:    name,
			Sex:     sex,
			FatRate: minFatRate,
		})
	}
}

func (r *FatRateRank) getRank(name string) (rank int, sex string, fatRate float64) {
	r.itemLock.Lock()
	defer r.itemLock.Unlock()
	sort.Slice(r.items, func(i, j int) bool {
		return r.items[i].FatRate < r.items[j].FatRate
	})
	frs := map[float64]struct{}{}
	for _, item := range r.items {
		frs[item.FatRate] = struct{}{}
		if item.Name == name {
			fatRate = item.FatRate
		}
	}
	rankArr := make([]float64, 0, len(frs))
	for k := range frs {
		rankArr = append(rankArr, k)
	}
	sort.Float64s(rankArr)
	for i, frItem := range rankArr {
		if frItem == fatRate {
			rank = i + 1
			sex = r.items[i].Sex
			break
		}
	}
	return
}
func (r *FatRateRank) getRankTop() []*types.PersonRank {
	r.itemLock.Lock()
	defer r.itemLock.Unlock()
	sort.Slice(r.items, func(i, j int) bool {
		return r.items[i].FatRate < r.items[j].FatRate
	})
	out := make([]*types.PersonRank, 0, 10)
	for i := 0; i < 10 && i < len(r.items); i++ {
		out = append(out, &types.PersonRank{
			RankNumber: i,
			Sex:        r.items[i].Sex,
			Fr: types.PersonalInformationFatRate{
				Name:    r.items[i].Name,
				Fatrate: r.items[i].FatRate,
			},
		})
	}

	return out
}
