package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// ctx, close := context.WithTimeout(context.TODO(), 5*time.Second)
	// defer close()
	// go directGet()
	// go httpGetWithContext(ctx)
	postMethod()
	time.Sleep(5 * time.Second)
}
func postMethod() {
	r := strings.NewReader("http://localhost:8088/?name=xiaoqiang&&sex=男")
	resp, err := http.Post("http://localhost:8088", "*/*", r)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("resp: " + string(data))
}
func directGet() {
	//Get
	resp, err := http.Get("http://localhost:8088")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))
}
func httpGetWithContext(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	//生成请求
	req, err := http.NewRequest("get", "http://localhost:8088", nil)
	if err != nil {
		log.Fatal("无法生成请求", err)
	}
	//将刚刚创建的上下文 ctx 传递给 HTTP 请求的 WithContext 方法，这样在进行请求时就会将上下文与请求关联起来
	//一定要新赋值，因为源码中新建了一个req
	req = req.WithContext(ctx)
	//执行请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("无法发送请求", err)
	}
	//读取请求的结果
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("无法读取返回内容", err)
	}
	fmt.Println(string(data))
}
