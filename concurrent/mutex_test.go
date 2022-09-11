package concurrent

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestLoadOrStore(t *testing.T) {
	sm := SafeMap[string, int]{
		values: make(map[string]int, 3),
	}
	go func() {
		fmt.Println(sm.LoadOrStore("key", 1))
	}()
	go func() {
		fmt.Println(sm.LoadOrStore("key", 2))
	}()
	go func() {
		fmt.Println(sm.LoadOrStore("key", 3))
	}()
	time.Sleep(1 * time.Second)
}

func TestUnSafeLoadOrStore(t *testing.T) {
	sm := SafeMap[string, int]{
		values: make(map[string]int, 3),
	}
	go func() {
		fmt.Println(sm.UnSafeLoadOrStore("key", 1))
	}()
	go func() {
		fmt.Println(sm.UnSafeLoadOrStore("key", 2))
	}()
	go func() {
		fmt.Println(sm.UnSafeLoadOrStore("key", 3))
	}()
	time.Sleep(1 * time.Second)
}

func TestDeferLoadOrStore(t *testing.T) {
	sm := SafeMap[string, int]{
		values: make(map[string]int, 3),
	}
	sm.DeferLoadOrStore("key", 1) // deadlock!
}

const Count = 1000000000 // 1亿次

// sync.Mutex和sync.RWMutex的区别
// Mutex: 互斥锁，任意两个锁都是互斥的，适用于读写不确定场景，即读写次数没有明显的区别，并且只允许只有一个读或者写的场景
// RWMutex：读写互斥锁，写操作都是互斥的、读和写是互斥的、读和读不互斥，该锁可以加多个读锁或者一个写锁，其经常用于读次数远远多于写次数的场景

type LockTest struct {
	sync sync.Mutex
}

func TestLock(t *testing.T) {
	runtime.GOMAXPROCS(2)
	start := time.Now()
	s := &LockTest{}
	for i := 0; i < Count; i++ {
		s.sync.Lock()
		s.sync.Unlock()
	}
	fmt.Println(time.Since(start)) // 16s
}

type RWLockTest struct {
	sync sync.RWMutex
}

func TestRWLock(t *testing.T) {
	runtime.GOMAXPROCS(2)
	start := time.Now()
	s := &RWLockTest{}
	for i := 0; i < Count; i++ {
		s.sync.Lock()
		s.sync.Unlock()
	}
	fmt.Println(time.Since(start)) // 36s
}

func TestAdd1(t *testing.T) {
	runtime.GOMAXPROCS(2)
	start := time.Now()
	a := 0
	for i := 0; i < Count; i++ {
		a += i
	}
	fmt.Println(time.Since(start)) // 普通+1操作需要0.5s
}
