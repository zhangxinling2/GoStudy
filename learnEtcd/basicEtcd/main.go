package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)

func main() {
	etcdWatcherDemo()
}
func etcdGetDemo() {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"http://localhost:2379"}})
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := cli.Get(context.TODO(), "a")
	if err != nil {
		log.Fatalln(err)
	}
	for _, kv := range resp.Kvs {
		log.Printf("key: %s, value : %s", kv.Key, kv.Value)
	}
}
func etcdWatcherDemo() {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"http://localhost:2379"}})
	if err != nil {
		log.Fatalln(err)
	}
	//watch和put的通讯
	dataCh := make(chan int)
	go func() {
		//Watch拿到的是channel
		watcher := cli.Watch(context.TODO(), "b")
		for respData := range watcher {
			//etcd是基于批量的，所以会有一系列的事件
			evs := respData.Events
			whetherBreak := false
			for _, ev := range evs {
				i, err := strconv.Atoi(string(ev.Kv.Value))
				if err != nil {
					fmt.Println("不是数字，结束")
					whetherBreak = true
					break
				}
				dataCh <- i
			}
			if whetherBreak {
				break
			}
		}
	}()
	go func() {
		for i := range dataCh {
			_, err := cli.Put(context.TODO(), "b", fmt.Sprintf("%d", i+1))
			if err != nil {
				fmt.Println("更新失败")
			}
		}
	}()
	time.Sleep(2 * time.Second)
}
