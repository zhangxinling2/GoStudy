package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//一个仓库可以让生产者放东西，可以让消费者拿东西，sync.Cond可以用来通知生产者，消费者
type store struct {
	max   int
	count int
	mutex sync.Mutex
	pCon  *sync.Cond
	cCon  *sync.Cond
}

//定义生产者和消费者
type producer struct{}
type consumer struct{}

func (p producer) Produce(s *store) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.count == s.max {
		fmt.Println("仓库已满")
		s.pCon.Wait()
	}
	fmt.Println("开始生产+1")
	s.count++
	s.pCon.Signal()
}
func (c consumer) Comsume(s *store) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.count == 0 {
		fmt.Println("仓库已空")
		s.pCon.Wait()
	}
	fmt.Println("开始购买-1")
	s.count--
	s.pCon.Signal()
}
func TestPC(t *testing.T) {
	s := &store{max: 10}
	s.cCon = sync.NewCond(&s.mutex)
	s.pCon = sync.NewCond(&s.mutex)
	cCount, pCount := 50, 50
	for i := 0; i < pCount; i++ {
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				producer{}.Produce(s)
			}
		}()
	}
	for i := 0; i < cCount; i++ {
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				consumer{}.Comsume(s)
			}
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("库存还剩%d", s.count)
}
