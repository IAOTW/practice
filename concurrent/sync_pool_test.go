package concurrent

import (
	"sync"
	"unsafe"
)

type MyCache struct {
	pool sync.Pool
}

// sync.Pool会准备一个资源池，比如一共有50个资源坑位
// 先查看pool中是否现有该资源，有则直接返回该资源，没有则创建一个
// 当gc的时候，会释放该资源，资源会放到到资源池
// 目的：复用内存，减少内存分配，不用一直开辟内存，减轻了gc的压力
// 对cpu的消耗较少，内存分配和gc都是cpu密集操作
type MyPool struct {
	p        sync.Pool
	maxCount int32
	count    int32
}

func (p *MyPool) Get() any {
	return p.p.Get()
}

func (p *MyPool) Put(val any) {
	// 大对象不放回去
	if unsafe.Sizeof(val) > 1024 {
		return
	}
	// 超过数量限制,以下代码由于gc的原因，效果不行
	//cnt := atomic.AddInt32(&p.count, 1)
	//if cnt >= p.maxCount {
	//	atomic.AddInt32(&p.count, -1)
	//	return
	//}
	p.p.Put(val)
}

// github有bytebufferpool的实现
