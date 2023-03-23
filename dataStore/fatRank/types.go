package fatRank

type PersonalInformationFatRate struct {
	Name    string
	Fatrate float64
}

type PersonRank struct {
	RankNumber int
	Sex        string
	Fr         PersonalInformationFatRate
}
