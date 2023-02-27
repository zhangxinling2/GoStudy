package code

import (
	"encoding/base64"
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
)

var (
	protoPath = "../testdata/proto"
)

func ProtoMarshal(info PersonInfo, filePath string) {
	data, err := proto.Marshal(&info)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
	output := base64.StdEncoding.EncodeToString(data)
	fmt.Println(output)
	writeFile(filePath, []byte(output))
}
func ProtoUnMarshal(filePath string) {
	p := PersonInfo{}
	data := readFile(filePath)
	output, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	proto.Unmarshal(output, &p)
	fmt.Println(p)
}
func TestProtoBuf(t *testing.T) {
	p := PersonInfo{
		Name:   "小强",
		Sex:    "男",
		Tall:   1.76,
		Weight: 71,
		Age:    35,
	}
	ProtoMarshal(p, protoPath)
	ProtoUnMarshal(protoPath)
}
