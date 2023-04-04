package frinterface

import (
	"GoStudy/dataStore/fatRank"
	"GoStudy/dataStore/types"
)

type ServeInterface interface {
	RegisterPersonInformation(pi *fatRank.PersonalInformation) error
	UpdatePersonInformation(pi *fatRank.PersonalInformation) (*types.PersonalInformationFatRate, error)
	GetFatrate(name string) (*types.PersonRank, error)
	GetTop() ([]*types.PersonRank, error)
}
type RankInitInterface interface {
	Init() error
}
