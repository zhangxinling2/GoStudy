package main

import (
	"GoStudy/dataStore/fatRank"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	//用来解析命令行参数
	var port string
	//用于解析一个字符串类型的命令行参数，并将其赋值给一个变量
	flag.StringVar(&port, "port", "8080", "配置启动端口")
	flag.Parse()
	//listen监听端口
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	rank := NewFatRateRank()
	fmt.Println("连接成功,开始接受数据")
	for {
		//Accept接收客户端请求
		conn, err := ln.Accept()
		if err != nil {
			log.Println("建立连接失败")
			continue
		}

		//talk(conn)
		//给多人提供服务
		//在talk时与rank做通信
		go talk(conn, rank)
	}
}
func talk(conn net.Conn, rank *FatRateRank) {
	defer fmt.Println("结束链接：", conn)
	defer conn.Close()
	// ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	// defer close()
	for {
		finalReceivedMeessage := make([]byte, 0, 1024)
		//如果读取数据失败则进行下一次读取
		enounterError := false
		for {
			buf := make([]byte, 1024)
			//类似文件读写，在没实现客户端之前，会一直循环 读数据失败
			valid, err := conn.Read(buf)
			if err != nil {
				enounterError = true
				log.Println("WARNING：读数据失败：", err)
				log.Println("重新读取", err)
				time.Sleep(1 * time.Second)
				break
			}
			if valid == 0 {
				break
			}
			//取出数据
			validContent := buf[:valid]
			finalReceivedMeessage = append(finalReceivedMeessage, validContent...)
			//如果没读满直接break
			if valid < len(buf) {
				break
			}
		}
		//如果读取数据失败则进行下一次读取
		if enounterError {
			//告诉客户端读取失败
			conn.Write([]byte("服务器获取数据失败，请重新输入"))
			continue
		}
		pi := &fatRank.PersonalInformation{}
		if err := json.Unmarshal(finalReceivedMeessage, pi); err != nil {
			conn.Write([]byte("数据无法解析，请重新输入"))
			continue
		}
		bmi, err := BMI(float64(pi.Weight), float64(pi.Tall))
		if err != nil {
			conn.Write([]byte("无法计算，请重新输入"))
			continue
		}
		fr := CalcFatRate(bmi, int(pi.Age), pi.Sex)
		rank.inputRecord(pi.Name, fr)
		rankId, _ := rank.getRank(pi.Name)
		log.Printf("姓名: %s,BMI: %v,体脂率: %v,排名: %d", pi.Name, bmi, fr, rankId)
		n, err := conn.Write([]byte(fmt.Sprintf("姓名: %s,BMI: %v,体脂率: %v,排名: %d", pi.Name, bmi, fr, rankId)))
		if n == 0 {
			conn.Write([]byte("服务器写入失败,写入0字节"))
			continue
		}
		if err != nil {
			conn.Write([]byte("服务器写入失败"))
			continue
		}
		//直接break 一次只服务一个人
		break
	}
}
func BMI(weightKG, heightM float64) (bmi float64, err error) {
	if weightKG < 0 {
		err = fmt.Errorf("weight cannot be negative")
		return
	}
	if heightM < 0 {
		err = fmt.Errorf("height cannot be negative")
		return
	}
	if heightM == 0 {
		err = fmt.Errorf("height cannot be 0")
		return
	}
	bmi = weightKG / (heightM * heightM)
	return
}
func CalcFatRate(bmi float64, age int, sex string) (fatRate float64) {
	sexWeight := 0
	if sex == "男" {
		sexWeight = 1
	} else {
		sexWeight = 0
	}
	fatRate = (1.2*bmi + getAgeWeight(age)*float64(age) - 5.4 - 10.8*float64(sexWeight)) / 100
	return
}

func getAgeWeight(age int) (ageWeight float64) {
	ageWeight = 0.23
	if age >= 30 && age <= 40 {
		ageWeight = 0.22
	}
	return
}
