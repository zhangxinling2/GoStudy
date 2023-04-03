package main

import (
	"GoStudy/dataStore/fatRank"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var _ fatRank.RankServiceServer = &rankServer{}

type rankServer struct {
	sync.Mutex
	persons map[string]*fatRank.PersonalInformation
	fatRank.UnimplementedRankServiceServer
}

func (r *rankServer) Register(ctx context.Context, information *fatRank.PersonalInformation) (*fatRank.PersonalInformation, error) {
	r.Lock()
	defer r.Unlock()
	r.persons[information.Name] = information
	log.Printf("收到新注册人%s\n", information.String())
	return information, nil
}

func startGRPCServer(ctx context.Context) {
	//监听端口
	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalln(err)
	}
	//新建一个server
	s := grpc.NewServer([]grpc.ServerOption{}...)
	//向server中注册服务
	fatRank.RegisterRankServiceServer(s, &rankServer{
		persons: map[string]*fatRank.PersonalInformation{},
	})
	go func() {
		select {
		case <-ctx.Done():
			s.Stop()
		}
	}()
	//启动server
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	startGRPCServer(ctx)
}
