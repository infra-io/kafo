# ☕ kafo

[![License](_icon/license.svg)](https://opensource.org/licenses/MIT)

**Kafo** 是一个高性能的轻量级分布式缓存中间件，支持 tcp/http 调用。

[Read me in English](./README.en.md)

### 📃 功能特性

* 提供 Get/Set/Delete/Status 几种调用接口
* 提供 HTTP 服务
* 支持获取缓存信息，比如 key 和 value 的占用空间
* 引入内存写满保护，使用 TTL 和 LRU 两种算法进行过期
* 引入 GC 机制，随机淘汰过期数据

_历史版本的特性请查看 [HISTORY.md](./HISTORY.md)。未来版本的新特性和计划请查看 [FUTURE.md](./FUTURE.md)。_

### 🔧 安装方式

（正在开发。。。）

### 📖 使用手册

（正在开发。。。）

### 🔥 性能测试

（正在开发。。。）

> 测试环境：R7-4700U CPU @ 2.0 GHZ，16 GB RAM

### 👤 贡献者

如果您觉得 **kafo** 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。

### 🔬 kafo 使用的技术

| 项目 | 作者 | 描述 | 链接 |
| -----------|--------|-------------|-------------------|
| logit | FishGoddess | 一个高性能、功能强大且极易上手的日志库 | [GitHub](https://github.com/FishGoddess/logit) / [码云](https://gitee.com/FishGoddess/logit) |
