package main

import (
	"encoding/json"
	"errors"
	"os"
)

type RankItem struct {
	Name    string
	FatRate float64
}
type Rank struct {
	items []RankItem
}

var (
	rFilePath = "./testdata/rank.json"
)

func (r RankItem) RegisterByJson() error {
	f, err := os.OpenFile(rFilePath, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("文件打开失败")
	}
	defer f.Close()

	data, err := json.Marshal(r)
	data = append(data, '\n')
	if err != nil {
		return errors.New("json失败")
	}
	f.WriteString(string(data))
	return nil
}
func (r *Rank) Register(name string, fatRate float64) error {
	rankItem := RankItem{
		Name:    name,
		FatRate: fatRate,
	}
	err := rankItem.RegisterByJson()
	if err != nil {
		return err
	}
	r.items = append(r.items, rankItem)
	QuickSort(r, 0, len(r.items)-1)
	return nil
}

func BubleSort(r *Rank) {
	for i := 0; i < len(r.items); i++ {
		flag := false
		for j := len(r.items) - 1; j > i; j-- {
			if r.items[j].FatRate < r.items[j-1].FatRate {
				r.items[j], r.items[j-1] = r.items[j-1], r.items[j]
				flag = true
			}
		}
		if !flag {
			break
		}
	}
}
func QuickSort(r *Rank, left, right int) {
	if left >= right || left > len(r.items)-1 {
		return
	}
	pivot := r.items[0]
	for left < right {
		if left < right && r.items[right].FatRate > pivot.FatRate {
			right--
		}
		r.items[right] = r.items[left]
		if left < right && r.items[left].FatRate < pivot.FatRate {
			left++
		}
		r.items[left] = r.items[right]
	}
	r.items[left] = pivot
	QuickSort(r, 0, left-1)
	QuickSort(r, left+1, len(r.items)-1)
}
func (r Rank) GetRank(name string) (int, error) {
	i := 0
	for r.items[i].Name != name {
		i++
	}
	return i + 1, nil
}
