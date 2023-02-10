package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Store struct {
	init  sync.Once
	max   int
	store chan int
}
type cproducer struct{}
type cconsumer struct{}

func (s *Store) initStore() {
	s.init.Do(func() {
		s.store = make(chan int, s.max)
	})
}

func (c cproducer) produce(s *Store) {
	fmt.Println("开始生产，库存+1")
	s.store <- 1
}

func (c cconsumer) consume(s *Store) {
	fmt.Println("开始购买，库存-1")
	<-s.store
}
func TestStore(t *testing.T) {
	s := &Store{max: 10}
	s.initStore()
	cCount, pCount := 50, 50
	for i := 0; i < pCount; i++ {
		go func(i int) {
			cproducer{}.produce(s)
		}(i)
	}
	for i := 0; i < cCount; i++ {
		go func(i int) {
			cconsumer{}.consume(s)
		}(i)
	}
	time.Sleep(5 * time.Second)
	//fmt.Printf("库存余量:%d", s.count)
	close(s.store)
}
