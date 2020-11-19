# ☕ kafo

[![License](_icon/license.svg)](https://opensource.org/licenses/MIT)

**Kafo** is a high-performance and distributed cache middleware with persistence and http/tcp mixed interface.

[阅读中文版的 Read me](./README.md)

### 📃 Features

* Use Gossip to communicate with each node in cluster
* Consistent hash supports, each node has independent data
* Get/Set/Delete/Status interface supports
* HTTP / TCP server supports
* Status in control, such as memory size
* Memory exceed protection, eliminating entries with ttl and lru
* Automatically gc, using random strategy to clean up the dead entries.
* Persistence supports, based on memory snapshot.

_Check [HISTORY.md](./HISTORY.md) and [FUTURE.md](./FUTURE.md) to know about more information._

### 📖 Guides

```bash
$ go run main.go
```

Client: [Github](https://github.com/avino-plan/kafo-client) / [码云](https://gitee.com/avino-plan/kafo-client).

### 🔥 Benchmarks

> Case: 10000 entries, writing and reading

> Environment: R7-4700U CPU @ 2.0 GHZ，16 GB RAM

| type | Write | Read |
|------|-------|------|
| http | 689.3ms | 5272.1ms | 1897 |
| tcp | 403.9ms | 387.1ms | 25833 |

More detail in [_examples/performance_test.go](./_examples/performance_test.go).

### 👤 Contributing

If you find that something is not working as expected please open an _**issue**_.

### 🔬 Projects kafo used

| Project | Author | Description | link |
| -----------|--------|-------------|------------------|
| logit | FishGoddess | A high-performance and easy-to-use logging foundation | [GitHub](https://github.com/FishGoddess/logit) / [Gitee](https://gitee.com/FishGoddess/logit) |
| vex | FishGoddess | A high-performance and easy-to-use net foundation | [GitHub](https://github.com/FishGoddess/vex) / [Gitee](https://gitee.com/FishGoddess/vex) |
