package main

import (
	"GoStudy/dataStore/fatRank"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
)

var _ fatRank.RankServiceServer = &rankServer{}

type rankServer struct {
	sync.Mutex
	persons map[string]*fatRank.PersonalInformation
	fatRank.UnimplementedRankServiceServer
	personCh chan *fatRank.PersonalInformation
}

func (r *rankServer) regPerson(p *fatRank.PersonalInformation) {
	r.Lock()
	defer r.Unlock()
	r.persons[p.Name] = p
	r.personCh <- p
}
func (r *rankServer) WatchPersons(null *fatRank.Null, server fatRank.RankService_WatchPersonsServer) error {
	for pi := range r.personCh {
		if err := server.Send(pi); err != nil {
			log.Println("发送失败", err)
			return err
		}
	}
	return nil
}

func (r *rankServer) RegisterPerson(server fatRank.RankService_RegisterPersonServer) error {
	pis := &fatRank.PersonalInformationList{}
	for {
		pi, err := server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}
		pis.Items = append(pis.Items, pi)
		//注册换成regRegistry
		r.regPerson(pi)
	}
	return server.SendAndClose(pis)
}

func (r *rankServer) Register(ctx context.Context, information *fatRank.PersonalInformation) (*fatRank.PersonalInformation, error) {
	r.regPerson(information)
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
		persons:  map[string]*fatRank.PersonalInformation{},
		personCh: make(chan *fatRank.PersonalInformation, 1),
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
