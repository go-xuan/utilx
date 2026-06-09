# utilx

Go 通用工具集，包含 22 个子包，覆盖加解密、文件操作、切片泛型、并发任务、时间处理等高频场景。

## 安装

```bash
go get github.com/go-xuan/utilx
```

## 子包速查

| 子包 | 用途 |
|------|------|
| **bytex** | 字节单位转换（Byte/KB/MB/GB/TB） |
| **contextx** | Context 增强工具 |
| **cryptox** | AES 加解密，子包 `cryptox/rsa` 提供 RSA 密钥与加解密 |
| **encodingx** | 编码转换 |
| **errorx** | 带堆栈的错误处理（New/Wrap/Wrapf/Is/As） |
| **excelx** | Excel 读写（基于 xlsx） |
| **execx** | 命令执行 |
| **filex** | 文件操作 |
| **funcx** | 函数包装执行（Execute/Merge/Wrap/Duration/Recover） |
| **httpx** | HTTP 客户端增强 |
| **idx** | ID 生成（UUID 等） |
| **marshalx** | 序列化/反序列化 |
| **maskx** | 数据脱敏 |
| **mathx** | 数学计算 |
| **osx** | 操作系统相关 |
| **randx** | 随机数/字符串生成 |
| **reflectx** | 反射工具 |
| **slicex** | 泛型切片操作（Contains/Distinct/Union/Intersect/Map/Filter/Reduce/Conv2Map） |
| **stringx** | 字符串处理（Parse/SnakeCase/CamelCase/Similarity/MatchUrl） |
| **taskx** | 并发任务（Concurrency/Retry/CronJob/ResultHook） |
| **timex** | 时间处理工具 |
