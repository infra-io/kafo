## ✒ 未来版本的新特性 (Features in future versions)

### v1.x
* 自定义 TCP 通信协议
* 加入分段锁，细化锁粒度
* 一致性哈希，分布式运行架构
* 引入 GC 遍历个数限制

### v0.1.0-alpha
* 提供 Get/Set/Delete 三种基本功能
* 提供 http 调用接口
* 支持获取缓存信息，比如 key 和 value 的占用空间
* 引入内存写满保护，使用 TTL 和 LRU 两种算法进行过期
* 引入 GC 机制，随机淘汰过期数据
