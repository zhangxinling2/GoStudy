package main

import (
	"GoStudy/dataStore/fatRank"
	"crypto/rand"
	"math/big"
)

type Interface interface {
	ReadPersonalInformation() (fatRank.PersonalInformation, error)
}

var _ Interface = &FakeInterface{}

//FakeInterface传输假数据
type FakeInterface struct {
	baseWeight float64
	baseTall   float64
	baseAge    int
	name       string
	sex        string
}

func (f *FakeInterface) ReadPersonalInformation() (fatRank.PersonalInformation, error) {
	r, _ := rand.Int(rand.Reader, big.NewInt(1000))
	out := float64(r.Int64()) / 1000
	if r.Int64()%2 == 0 {
		out = 0 - out
	}
	pi := fatRank.PersonalInformation{
		Name:   f.name,
		Sex:    f.sex,
		Age:    int64(f.baseAge),
		Weight: float32(f.baseWeight),
		Tall:   float32(f.baseTall),
	}
	f.baseWeight += out
	return pi, nil
}
