package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithDeadline(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	cancel()
	fmt.Println(childCtx.Deadline()) // 返回定时的时间，和true，如果没有设置定时取消则为公元起点时间和false
	fmt.Println(childCtx.Err())      // 未cancel时为nil
	go func(c context.CancelFunc) {
		select {
		case <-childCtx.Done():
			fmt.Println(childCtx.Err())
		default:
			fmt.Println(childCtx.Err())
			fmt.Println("当走到这里，还没有被cancel时,会走默认分支")
		}
		cancel()
		fmt.Println(childCtx.Deadline()) // 返回定时的时间，和true，如果没有设置定时取消则为公元起点时间和false
	}(cancel)
	select {
	case <-childCtx.Done():
		fmt.Println(childCtx.Err()) // 只要cancel，则必定会走进来，手动cancel为context canceled，超时自动cancel为context deadline exceeded
		// 只要在运行到此处之前被cancel，则后续所有的<-childCtx.Done()都能拿到一个struct{},因为closedchan :=make(chan struct{}) init函数会close chan
		// 当运行cancel时，会把这个closedchan存入到ctx.done中，而从关闭的chan中取值会一直拿到空值(所有的childCtx.Done()，都是从ctx.done中找哪个closedchan)
	}
}

func TestWithTimeout(t *testing.T) {
	// WithTimeout 进源码会看到它直接调withdeadline，基本没区别
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept, 睡过头了")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}
}

func TestWithValue(t *testing.T) {
	valueCtx := context.WithValue(context.Background(), "a", 1)
	fmt.Println(valueCtx.Value("a"))
	childValueCtx := context.WithValue(valueCtx, "b", 2)
	fmt.Println(childValueCtx.Value("a")) // 子ctx可以拿到父ctx的value
	fmt.Println(childValueCtx.Value("b"))
	fmt.Println(valueCtx.Value("b")) // 父ctx拿不到子ctx的value
	fmt.Println(valueCtx.Deadline())
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	time.Sleep(3 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("已被cancel")
	default:
		fmt.Println("走到select时还没有cancel，则走默认分支")
	}
	select {
	case <-ctx.Done():
		fmt.Println("已被cancel")
	default:
		fmt.Println("走到select时还没有cancel，则走默认分支")
	}
}

func TestChildCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	childCtx, childCancel := context.WithCancel(ctx)
	childCancel()
	fmt.Println(ctx.Err()) // cancel子ctx，只会cancel子ctx以及子ctx的后代，而不会取消父ctx
	fmt.Println(childCtx.Err())
	cancel()
	//fmt.Println(ctx.Err()) // cancel父ctx，会cancel父ctx以及它的后代
	//fmt.Println(childCtx.Err())
}
