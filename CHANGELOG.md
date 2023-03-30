# Gaea Changelog

## Gaea 2.1 Release Notes

### 新功能
- 新增 SQL 支持：分片情况下支持在 `group by ... order by max(col)` 语句 ORDER BY 使用聚合函数
- 新增 SQL 支持：子查询中的多括号支持
- 新增 SQL 支持：`create table t (... on update current_timestamp(NUMBER))` 
- 新增 SQL 行为支持：`select @@read_only` 及 `show variables like "%read_only%"` 发送至主库，避免应用框架通过从库获取数据库只读状态导致业务报错

### Bug 修复
- 修复 2.0 版本部分 union 语句无法执行的问题
- 回退 2.0 版本 Gaea 启动时初始化后端数据库连接，避免影响多租户下的集群启动时间

### 优化提升
- 新增 API 用于获取当前 Gaea 版本，便于后续自动化升级降级操作
- 新增了大量只读 SQL 测试 Case，提高后续版本稳定性，测试 Case 会在后续版本继续完善

## Gaea 2.0-GA Release Notes

### 新功能

- 支持限制 Gaea 的 CPU 核数，便于将 Gaea 作为单租户使用

### Bug修复

- 修复不兼容 SQL：支持 `password` 函数
- 修复不兼容 SQL：支持 `group by ... with rollup` 语句
- 修复不兼容 SQL：支持分片情况下 `order by null` 语句
- 修复 LVS 探活引起的 Broken Pipe 问题
- 修复后端从库探活偶尔失败的问题

### 优化提升

- 优化了 general_log 配置项，直接修改 namespace 配置即可在所有 Gaea 节点生效
- 增加监控项，监控 Gaea CPU 使用情况
- 修复本地 Mac 环境下无法生成 parser 文件的问题

### 其他说明

- 进行了 Gaea 与 Miproxy 的做了基本的性能对比，两者各有较小优劣，整体对比差距不大。详见：https://xiaomi.f.mioffice.cn/wiki/wikk4HMGGB0A56lU1NyCj5likQb


## Gaea 2.0-Beta Release Notes

### 新功能

- 支持配置日志保留天数，默认保留 3 天
- 支持将 `select ... for update` or `select ... in share mode` 语句优先发给主库

### Bug修复

- 修复了 Alpha 版本测试发现的旧版兼容性问题，如配置后端实例为空字符串（默认为空数组）

### 优化提升

- 优化了日志展示，去除了冗余字段及冗余日志
- 新增二进制文件自动更新，用户可以自行拉取新版本。具体使用见用户文档
- 新增监控项，用于 Gaea 中间件存活状态及后端实例存活告警

### 其他说明

- 本次新增了部署、迁移、用户等文档，便于后续使用。详见：https://xiaomi.f.mioffice.cn/wiki/wikk4RF3oZ2XMT5hOJm8ISb7Jqh


## Gaea 2.0-Alpha Release Notes

### 新功能

- 兼容 MiProxy 强制读主库语句：`select /*master*/ from ...`
- 支持配置连接后端 MySQL 时自定义 [Capability Flags](https://dev.mysql.com/doc/dev/mysql-server/latest/group__group__cs__capabilities__flags.html) 
- 后端实例探活：支持主从库探活并自动剔除，降低单实例宕机对业务的影响
- 后端从库状态探测：支持从库复制状态及延迟探测并自动剔除，降低主从数据不一致对业务的影响
- 支持限制前端连接数配置
- 支持配置后端连接初始化时执行特定 SQL

### Bug修复

- 修复 show 语句会默认发到主库的问题
- 修复 gaea-cc 修改配置失败会阻塞的问题
- 修复 gaea-cc 修改配置失败未回滚 etcd 配置的问题

### 优化提升

- 连接池优化：初始化时创建初始连接
- 日志：支持输出 SQL 日志和慢日志
- 日志：支持 SQL 日志输出到单独的文件
- 配置：增加配置检查及配置错误提示

### 测试

- 新增代码静态检测，统一代码风格并减少低级错误
- 新增集成测试框架，支持测试各种 SQL 在 Gaea 及实际 MySQL 结果对比
- 新增 e2e 测试框架，支持端对端验证 Gaea 从启动到停止过程中的所有行为
- 新增 CI，包含代码静态检查、单元测试、集成测试、e2e测试

### 其他说明

- 本版本基于公司开源的 Gaea(https://github.com/XiaoMi/Gaea) 最新的 [commit](https://github.com/XiaoMi/Gaea/commit/9e16060c4afcabe8ca910313e2abde4113e9af79) 开发。
