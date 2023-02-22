package file

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func ReadFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(n)
	fmt.Println(string(buf))
	fmt.Println(string(buf[:n]))

}

func WriteFile(filePath string) {
	f, err := os.OpenFile(filePath, os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	n, err := f.Write([]byte("this is first line\n"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)
	n, err = f.WriteString("第二行内容\n")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)
	n, err = f.WriteAt([]byte("CHANGED"), io.SeekEnd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)
}
func TestFile(t *testing.T) {
	file := "../testdata/user"
	WriteFile(file)
	ReadFile(file)
}
