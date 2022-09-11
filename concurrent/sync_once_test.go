package concurrent

import (
	"fmt"
	"sync"
	"testing"
)

//once如果是在程序启动时就要运行，则可用init代替
type OnceClose struct {
	close sync.Once
}

func (o *OnceClose) Close() {
	o.close.Do(func() {
		fmt.Println("我只运行一次！")
	})
}

func TestOnceClose_Close(t *testing.T) {
	o := &OnceClose{}
	for i := 0; i < 100; i++ {
		o.Close()
	}
}

// 以下为Once的源码，其实是个Doubel Check机制（大佬写法：用atomic+Lock）
//func (o *Once) Do(f func()) {
//	// Note: Here is an incorrect implementation of Do:
//	//
//	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
//	//		f()
//	//	}
//	//
//	// Do guarantees that when it returns, f has finished.
//	// This implementation would not implement that guarantee:
//	// given two simultaneous calls, the winner of the cas would
//	// call f, and the second would return immediately, without
//	// waiting for the first's call to f to complete.
//	// This is why the slow path falls back to a mutex, and why
//	// the atomic.StoreUint32 must be delayed until after f returns.
//
//	if atomic.LoadUint32(&o.done) == 0 {
//		// Outlined slow-path to allow inlining of the fast-path.
//		o.doSlow(f)
//	}
//}
//
//func (o *Once) doSlow(f func()) {
//	o.m.Lock()
//	defer o.m.Unlock()
//	if o.done == 0 {
//		defer atomic.StoreUint32(&o.done, 1)
//		f()
//	}
//}
