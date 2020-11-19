# ☕ kafo

[![License](_icon/license.svg)](https://opensource.org/licenses/MIT)

**Kafo** 是一个高性能的轻量级分布式缓存中间件，支持 tcp/http 调用。

[Read me in English](./README.en.md)

### 📃 功能特性

* 使用 Gossip 协议进行分布式通信
* 加入一致性哈希，集群每个节点负责独立的数据
* 提供 Get/Set/Delete/Status 几种调用接口
* 提供 HTTP / TCP 两种调用服务
* 支持获取缓存信息，比如 key 和 value 的占用空间
* 引入内存写满保护，使用 TTL 和 LRU 两种算法进行过期
* 引入 GC 机制，随机淘汰过期数据
* 基于内存快照实现持久化功能

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 📖 使用手册

```bash
$ go run main.go
```

客户端：[Github](https://github.com/avino-plan/kafo-client) / [码云](https://gitee.com/avino-plan/kafo-client)。

### 🔥 性能测试

> 场景：10000 个键值对的写入和读取的耗时

> 环境：R7-4700U CPU @ 2.0 GHZ，16 GB RAM

| type | write | read | rps |
|------|-------|------| ----- |
| http | 689.3ms | 5272.1ms | 1897 |
| tcp | 403.9ms | 387.1ms | 25833 |

测试详情参考文件 [_examples/performance_test.go](./_examples/performance_test.go)。

### 👤 贡献者

如果您觉得 **kafo** 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。

### 🔬 kafo 使用的技术

| 项目 | 作者 | 描述 | 链接 |
| -----------|--------|-------------|-------------------|
| logit | FishGoddess | 一个高性能、功能强大且极易上手的日志库 | [GitHub](https://github.com/FishGoddess/logit) / [码云](https://gitee.com/FishGoddess/logit) |
| vex | FishGoddess | 一个高性能、且极易上手的网络通信框架 | [GitHub](https://github.com/FishGoddess/vex) / [码云](https://gitee.com/FishGoddess/vex) |
