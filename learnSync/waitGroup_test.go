package learnSync

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type runner struct {
	name string
}

func (r runner) run(c *sync.WaitGroup, wg *sync.WaitGroup) {
	fmt.Println("准备跑")
	c.Wait()
	defer wg.Done()
	start := time.Now()
	fmt.Println(r.name + "开始跑" + start.String())
	end := time.Now()
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	dur := end.Sub(start)
	fmt.Println(r.name + "完成跑,总耗时" + dur.String())
}

func TestRun(t *testing.T) {
	runCount := 10
	//定义一个裁判用来使runner一起执行
	c := sync.WaitGroup{}
	c.Add(1)

	r := sync.WaitGroup{}
	r.Add(10)
	runners := []runner{}
	for i := 0; i < runCount; i++ {
		runners = append(runners, runner{name: fmt.Sprintf("%d", i)})
	}
	for i := 0; i < runCount; i++ {
		go runners[i].run(&c, &r)
	}
	time.Sleep(1 * time.Second)
	c.Done()
	fmt.Println("裁判吹哨")
	r.Wait()
	fmt.Println("结束")
}
