package main

import (
	"GoStudy/dataStore/fatRank"
	"context"
	"google.golang.org/grpc"
	"log"
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
}
