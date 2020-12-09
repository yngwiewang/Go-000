学习笔记

# 开 goroutine 的原则：
1. 由调用者决定是否异步执行
2. 调用者要知道 goroutine 什么时候结束
3. 要管控 goroutine 的生命周期，控制它的退出

# 使用锁的原则：
最晚加锁，最早释放，临界区代码尽量少

# sync.pool 使用注意：
1. Trace 对象可以放在 pool 里
2. 只能放可被随时释放的对象，因为 gc 有可能回收 pool 中的对象

# context 使用注意：
1. 不要放在 struct 中，而是作为首参数传递给函数

# chan 使用注意：
一定要由 sender close chan

# atomic.Value 适用于配置中心这类读多写少的场景

# copy on write

# 内存模型：
1. 如果不用同步原语，不保证可见性
2. happens before