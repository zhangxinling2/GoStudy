package code

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func writeFile(filePath string, data []byte) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	f.Write(data)
}
func utilReadFile(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	return data
}
func readFile(filePath string) []byte {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()
	data := make([]byte, 4096)
	n, err := f.Read(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data[:n]
}

// type PersonInfo struct {
// 	Name   string  `json:"name"`
// 	Sex    string  `json:"sex"`
// 	Tall   float64 `json:"tall"`
// 	Weight int     `json:"weight"`
// 	Age    int     `json:"age"`
// }

var (
	filePath = "../testdata/info.json"
)

func TestJson(t *testing.T) {

	p := PersonInfo{
		Name:   "小强",
		Sex:    "男",
		Tall:   1.76,
		Weight: 71,
		Age:    35,
	}
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("marshal原生结果为", p)
	fmt.Println("marshal结果为", string(data))
	writeFile(filePath, data)
	data = readFile(filePath)
	tmp := PersonInfo{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("使用readFile", tmp)

	data = utilReadFile(filePath)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("使用ioutil.readFile", tmp)
}
