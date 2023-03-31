package moments

import "GoStudy/dataStore/fatRank"

type Moments interface {
	ReleaseMoment(text string) (fatRank.PersonalMoment, error)
	DeleteMoment(id int64) error
	GetAllMoment() (fatRank.PersonalMomentList, error)
}
