package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

var (
	GOOS   string
	GOARCH string
	File   string
)

func main() {
	initFlag()
	buildFile()
}
func initFlag() {
	flag.StringVar(&GOOS, "GOOS", "linux", "目标平台操作系统")
	flag.StringVar(&GOARCH, "GOARCH", "amd64", "目标平台体系架构")
	flag.StringVar(&File, "File", "main.go", "编译文件")
	flag.Parse()
}
func buildFile() {
	//set 环境变量
	os.Setenv("CGO_ENABLE", "0")
	os.Setenv("GOOS", GOOS)
	os.Setenv("GOARCH", GOARCH)
	//写命令
	cmd := exec.Command("go", "build", File)
	//设置命令执行的环境
	cmd.Env = os.Environ()
	//执行命令
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
