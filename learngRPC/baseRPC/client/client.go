package main

import (
	"GoStudy/dataStore/fatRank"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	c := fatRank.NewRankServiceClient(conn)
	ret, err := c.Register(context.TODO(), &fatRank.PersonalInformation{
		Id:     1,
		Name:   "Tom",
		Sex:    "男",
		Tall:   1.77,
		Weight: 66,
		Age:    18,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("注册成功", ret)
	//得到的是注册的client
	regClient, err := c.RegisterPerson(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}
	if err := regClient.Send(&fatRank.PersonalInformation{
		Id:     1,
		Name:   fmt.Sprintf("tom-%d", time.Now().UnixNano()),
		Sex:    "男",
		Tall:   1.77,
		Weight: 66,
		Age:    18,
	}); err != nil {
		log.Println("注册失败")
	}
	if err := regClient.Send(&fatRank.PersonalInformation{
		Id:     1,
		Name:   fmt.Sprintf("tom-%d", time.Now().UnixNano()),
		Sex:    "男",
		Tall:   1.77,
		Weight: 66,
		Age:    18,
	}); err != nil {
		log.Println("注册失败")
	}
	if err := regClient.Send(&fatRank.PersonalInformation{
		Id:     1,
		Name:   fmt.Sprintf("tom-%d", time.Now().UnixNano()),
		Sex:    "男",
		Tall:   1.77,
		Weight: 66,
		Age:    18,
	}); err != nil {
		log.Println("注册失败")
	}
	resp, err := regClient.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("注册成功", resp.String())
}
