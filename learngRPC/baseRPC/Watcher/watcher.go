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
	w, err := c.WatchPersons(context.TODO(), &fatRank.Null{})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		pi, err := w.Recv()
		if err != nil {
			log.Fatalln("接收异常")
			break
		}
		log.Println("收到变动", pi.String())
	}
}
