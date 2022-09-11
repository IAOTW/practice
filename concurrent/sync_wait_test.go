package concurrent

import (
	"fmt"
	"sync"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	//runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	for i := 10; i < 20; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("i: ", i)
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("i: ", i)
		}(i)
	}
	wg.Wait()
}
