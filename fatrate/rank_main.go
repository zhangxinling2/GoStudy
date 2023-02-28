package main

import (
	"fmt"
)

func main() {
	input := &inputFromStd{}
	rank := Rank{}
	for {
		pi := input.GetInput()
		pi.calcBmi()
		pi.calcFatRate()
		err := pi.RegisterByJson()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = rank.Register(pi.Name, pi.FatRate)
		if err != nil {
			return
		}
		res, err := rank.GetRank(pi.Name)
		if err != nil {
			return
		}
		fmt.Println("排名结果为", res)
	}
}

// func TestRank(t *testing.T) {
// 	input := &inputFromStd{}
// 	rank := Rank{}
// 	for {
// 		pi := input.GetInput()
// 		err := pi.RegisterByJson()
// 		if err != nil {
// 			return
// 		}
// 		pi.calcBmi()
// 		pi.calcFatRate()
// 		err = rank.Register(pi.name, pi.fatRate)
// 		if err != nil {
// 			return
// 		}
// 		res, err := rank.GetRank(pi.name)
// 		if err != nil {
// 			return
// 		}
// 		fmt.Println("排名结果为", res)
// 	}
// }
