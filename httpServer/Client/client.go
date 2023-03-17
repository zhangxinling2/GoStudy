package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	var input Interface = &FakeInterface{
		name:       "小强",
		sex:        "男",
		baseWeight: 71.0,
		baseTall:   1.70,
		baseAge:    39,
	}
	for {
		//由于defer放在循环中容易内存泄漏所以写一个匿名方法存在方法中
		func() {
			conn, err := net.Dial("tcp", "localhost:8080")
			if err != nil {
				log.Fatalln("连接失败")
			}
			defer conn.Close()
			//读取输入
			fmt.Println("连接成功，开始发送数据")
			pi, err := input.ReadPersonalInformation()
			if err != nil {
				log.Println("WARNING: 读取失败请重新输入", err)
				return
			}
			data, err := json.Marshal(pi)
			if err != nil {
				log.Println("WARNING: 无法编码个人信息", err)
				return
			}
			// r := bufio.NewReader(os.Stdin)
			// input, _, _ := r.ReadLine()

			talk(conn, string(data))
		}()
		time.Sleep(1 * time.Second)
	}

}

func talk(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("通信失败")
	} else {
		data := make([]byte, 1024)
		valid, err := conn.Read(data)
		if err != nil {
			log.Println("warning:服务器返回数据失败")
		} else {
			content := data[:valid]
			log.Println("去：", message, "回：", string(content))
		}
	}
}
