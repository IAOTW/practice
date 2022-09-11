package concurrent

// sync.Cond 经常用在多个 Goroutine 等待，一个 Goroutine 通知（事件发生）的场景。如果是一个通知，一个等待，使用互斥锁或 channel 就能搞定了。
// channel 可以实现，代码更加简洁，那么 sync.Cond 的存在还有必要吗？
//实际上 sync.Cond 与 Channel 是有区别的，channel 定位于通信，用于一发一收的场景，sync.Cond 定位于同步，用于一发多收的场景。虽然 channel 可以通过 close 操作来达到一发多收的效果，但是 closed 的 channel 已无法继续使用，而 sync.Cond 依旧可以继续使用。这可能就是“全能”与“专精”的区别。
// https://geektutu.com/post/hpg-sync-cond.html
