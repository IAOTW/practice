package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestCloseChannel(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		time.Sleep(1 * time.Second)
		close(ch)
		c, ok := <-ch // 当关闭的管道中有数据时，会先取完：1 true
		fmt.Println(c, ok)

		c, ok = <-ch // 当关闭的管道为空：类型默认值 false
		fmt.Println(c, ok)

		ch <- 2
		fmt.Println("往关闭的channel中塞元素会报错：panic: send on closed channel")
	}()
	ch <- 1
	time.Sleep(3 * time.Second)
}

// 透支channel
func TestOverdraftChannel(t *testing.T) {
	ch := make(chan int, 3)
	go func() {
		time.Sleep(1 * time.Second)
		//close(ch)
		c, ok := <-ch // 1 true
		fmt.Println(c, ok)

		c, ok = <-ch // 当未关闭管道为空，阻塞住，直到塞入一个元素
		fmt.Println(c, ok)

		fmt.Println("当未关闭管道为空，且没有元素进去，一直阻塞住")
		c, ok = <-ch
		fmt.Println(c, ok) // 不会走到这一步
	}()
	ch <- 1
	time.Sleep(3 * time.Second)
	ch <- 2
	time.Sleep(10 * time.Second)
}

func TestUnBufferChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		<-ch
		fmt.Println("准备取数据")
	}()
	ch <- 1
	fmt.Println("无缓冲chan，阻塞，当其他goroutine准备取数据的那一刻才能往ch里面塞")
}

func TestBufferChannel(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		time.Sleep(3 * time.Second)
		c := <-ch
		fmt.Println("准备取数据: ", c)
	}()
	ch <- 1
	fmt.Println("有缓冲，放入可以不等待，但取出的goroutine必须要等待！！！")
	time.Sleep(5 * time.Second)
}
