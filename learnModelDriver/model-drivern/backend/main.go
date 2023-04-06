package main

import (
	"GoStudy/learnDocker/dockertool"
	"GoStudy/learnModelDriver/model-drivern/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

func main() {
	backendName := "/rankservice/backend"
	//etcd客户端
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2380"}})
	if err != nil {
		log.Fatalln(err)
	}
	//docker客户端
	dockerCli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	go func() {
		watcher := cli.Watch(ctx, backendName)
		for respData := range watcher {
			evs := respData.Events
			for _, ev := range evs {
				rawBackend := ev.Kv.Value
				backend := &models.RankServiceBackend{}
				json.Unmarshal(rawBackend, backend)
				//拿到后比较
				ids, _ := dockertool.List(ctx, dockerCli, map[string]string{"backend": backendName})
				backend.Status.RunningCount = len(ids)
				if backend.Expected.Count != backend.Status.RunningCount {
					fmt.Println("开始创建实例")
					//todo 调用docker创建新的容器实例，或 删除一些，并更新 Status
					ip, err := dockertool.Run(ctx, dockerCli, map[string]string{"backend": backendName}, "nginx:stable-alpine", nil)
					if err != nil {
						//retry
					} else {
						backend.Status.RunningCount++
						backend.Status.InstanceIPs = append(backend.Status.InstanceIPs, ip)
						rawData, _ := json.Marshal(backend)
						cli.Put(ctx, backendName, string(rawData))
					}
				} else {
					fmt.Println("已满足预期", backend.Expected.Count)
				}
			}
		}
	}()
	<-ctx.Done()
}
