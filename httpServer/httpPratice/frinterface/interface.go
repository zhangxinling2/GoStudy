package frinterface

import "GoStudy/dataStore/fatRank"

type ServeInterface interface {
	RegisterPersonInformation(pi *fatRank.PersonalInformation) error
	UpdatePersonInformation(pi *fatRank.PersonalInformation) (*fatRank.PersonalInformationFatRate, error)
	GetFatrate(name string) (*fatRank.PersonRank, error)
	GetTop() ([]*fatRank.PersonRank, error)
}
