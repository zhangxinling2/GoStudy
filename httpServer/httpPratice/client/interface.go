package main

import (
	"GoStudy/dataStore/fatRank"
	"crypto/rand"
	"math/big"
)

type Interface interface {
	ReadPersonalInformation() (fatRank.PersonalInformation, error)
}

var _ Interface = &fakeInterface{}

type fakeInterface struct {
	name       string
	sex        string
	baseWeight float64
	baseTall   float64
	baseAge    int
}

func (f *fakeInterface) ReadPersonalInformation() (fatRank.PersonalInformation, error) {
	r, _ := rand.Int(rand.Reader, big.NewInt(1000))
	out := float64(r.Int64()) / 1000
	if r.Int64()%2 == 0 {
		out = 0 - out
	}
	pi := fatRank.PersonalInformation{
		Name:   f.name,
		Sex:    f.sex,
		Tall:   float32(f.baseTall),
		Weight: float32(f.baseWeight),
		Age:    int64(f.baseAge),
	}
	f.baseWeight += out
	return pi, nil
}
