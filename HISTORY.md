## ✒ 历史版本的特性介绍 (Features in old versions)

### v0.1.0-alpha
* 提供 Get/Set/Delete 三种基本功能
* 提供 http 调用接口
* 支持获取缓存信息，比如 key 和 value 的占用空间
* 引入内存写满保护，使用 TTL 和 LRU 两种算法进行过期
* 引入 GC 机制，随机淘汰过期数据
