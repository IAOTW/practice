package concurrent

import (
	"fmt"
	"testing"
)

// 实现消费队列-发布订阅
// 方式1：每一个订阅者都维护一个ch，发布者发布时，会将消息给到每一个订阅者
type Consumer struct {
	ch chan string
}

type Broker struct {
	consumers []*Consumer
}

// 订阅
func (b *Broker) Subscribe(c *Consumer) {
	b.consumers = append(b.consumers, c)
}

func (b *Broker) Publish(msg string) {
	for _, c := range b.consumers {
		c.ch <- msg
	}
}

func TestBroker(t *testing.T) {
	b := &Broker{
		consumers: make([]*Consumer, 0, 10),
	}
	c1 := &Consumer{
		ch: make(chan string, 1),
	}
	c2 := &Consumer{
		ch: make(chan string, 1),
	}
	b.Subscribe(c1)
	b.Subscribe(c2)
	b.Publish("Hello，everyone!")
	fmt.Println(<-c1.ch)
	fmt.Println(<-c2.ch)
}

// 方式2：
// 订阅推送方法
type consumeFunc func(s string)

type BrokerV2 struct {
	ch       chan string
	consumes []consumeFunc
}

func (b *BrokerV2) SubscribeV2(c consumeFunc) {
	b.consumes = append(b.consumes, c)
}

func (b *BrokerV2) PublishV2(msg string) {
	b.ch <- msg
}

func (b *BrokerV2) Start() {
	go func() {
		msg := <-b.ch
		for _, consumeFunc := range b.consumes {
			consumeFunc(msg)
		}
	}()
}

func NewBrokerV2() *BrokerV2 {
	b := &BrokerV2{
		ch:       make(chan string, 10),
		consumes: make([]consumeFunc, 0, 10),
	}
	go func() {
		msg := <-b.ch
		for _, consumeFunc := range b.consumes {
			consumeFunc(msg)
		}
	}()
	return b
}
