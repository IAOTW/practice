package concurrent

import (
	"sync"
	"time"
)

type SafeMap[K comparable, V any] struct {
	values map[K]V
	lock   sync.RWMutex
}

// LoadOrStore: Double Check机制
// 已经有key，返回对应的值，然后loaded=true
// 没有，则放入，返回 loaded false
func (s *SafeMap[K, V]) LoadOrStore(key K, newValue V) (val V, loaded bool) {
	s.lock.RLock() // 加读锁可以一直多个goroutine都加，都能读到
	oldVal, ok := s.values[key]
	s.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	//time.Sleep(500 * time.Millisecond) // 此处放大复现几率
	s.lock.Lock()
	defer s.lock.Unlock()
	oldVal, ok = s.values[key]
	if ok {
		return oldVal, true
	}
	s.values[key] = newValue
	return newValue, false
}

// 无Double Check示例（错误示例）
// goroutine1： ("key", 1)
// goroutine2： ("key", 2)
// 两个goroutine同时读，都没读到，此时进到写的环节
func (s *SafeMap[K, V]) UnSafeLoadOrStore(key K, newValue V) (val V, loaded bool) {
	s.lock.RLock() // 加读锁可以一直多个goroutine都加，都能读到
	oldVal, ok := s.values[key]
	s.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	time.Sleep(500 * time.Millisecond) // 此处放大复现几率
	s.lock.Lock()
	defer s.lock.Unlock()
	// goroutine1先进来，那么 "key" = 1
	// goroutine2后进来，按照原设计，后进来的应该读到1才对，可是他却覆盖了，故不安全
	s.values[key] = newValue
	return newValue, false
}

// 死锁示例
func (s *SafeMap[K, V]) DeferLoadOrStore(key K, newValue V) (val V, loaded bool) {
	s.lock.RLock() // 加读锁可以一直多个goroutine都加，都能读到
	oldVal, ok := s.values[key]
	defer s.lock.RUnlock() // 延后执行
	if ok {
		return oldVal, true
	}
	s.lock.Lock() // 当运行到此处，由于上面的RUnlock还未执行，即拿着写锁，再去拿写锁，会死锁
	defer s.lock.Unlock()
	oldVal, ok = s.values[key]
	if ok {
		return oldVal, true
	}
	s.values[key] = newValue
	return newValue, false
}
