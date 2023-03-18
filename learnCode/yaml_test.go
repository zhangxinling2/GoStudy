package learnCode

import (
	"fmt"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

var (
	yamlPath = "../testdata/info.yaml"
)

func YamlMarshal(info PersonInfo, filePath string) {
	data, err := yaml.Marshal(info)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))

	data = append([]byte("---\n"), data...)
	writeFile(filePath, data)
}
func YamlUnMarshal(filePath string) {
	p := PersonInfo{}
	data := readFile(filePath)
	yaml.Unmarshal(data, &p)
	fmt.Println(p)
}
func TestYaml(t *testing.T) {
	p := PersonInfo{
		Name:   "小强",
		Sex:    "男",
		Tall:   1.76,
		Weight: 71,
		Age:    35,
	}
	YamlMarshal(p, yamlPath)
	YamlUnMarshal(filePath)
}
