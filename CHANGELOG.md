# Gaea Changelog

## Gaea 2.3 Release Notes

### 新功能

- 支持 group_concat_max_len 变量设置
- 支持可配置的 session 级别变量
- 支持分库分表情况下，显示详细的 sql explain 信息
- 支持客户端使用 hash 后的密码连接 Gaea
- 支持非分片情况下使用 MySQL 8.0
- 支持非分片情况下大结果集流式返回，避免较大内存占用问题
- 支持日志配置热加载
- 支持日志文件保留数量的配置
- 新增 keep session 功能，keep session 模式下，会一对一保持前后端数据库连接会话
- 新增 分片情况下 group by ... order by max(col) 语句 ORDER BY 使用聚合函数
- 新增 SQL 支持：子查询中的多括号支持
- 新增 SQL 支持：create table t (... on update current_timestamp(NUMBER))
- 新增 SQL 行为支持：select @@read_only 及 show variables like "%read_only%" 发送至主库，避免应用框架通过从库获取数据库只读状态导致业务报错
- 新增 SQL 支持：默认支持 multistatement
- 新增 API 用于获取当前 Gaea 版本，便于后续自动化升级降级操作
- 新增了大量只读 SQL 测试 Case，提高后续版本稳定性

### 优化提升

- 优化 Gaea session timeout 功能，在高 qps 情况下，cpu 使用率较高问题
- 优化日志功能，更换日志库，新增异步日志写入能力，提高日志性能
- 优化请求处理逻辑，提高单条 SQL 处理效率
- 优化 result 结果集对象复用，降低内存使用
- 优化 MySQL 协议写入数据包逻辑，提升网络写入效率，降低 P99延迟
- 优化计算 p99/p95 cpu 资源占用率过高问题
- 优化 prepare execute 参数过多时 cpu 占用率过高问题
- 优化 SQL 黑名单和连接数过多的报错信息
- 优化获取 buffer pool 性能
- 优化 gaea_proxy_sql_error_counts 监控项，只记录 Gaea 处理过程中的错误，不记录后端 MySQL 执行的错误
- 优化 general 日志展示，增加用户名及执行从库显示
- 忽略 `lock table ...` 语句，解决避免后端未解锁表导致复用失败的问题
- 调整 gaea_proxy_cpu_busy 监控项名为 gaea_proxy_cpu_busy_by_core，区分新旧版本值
- 修改 CheckSelectLock 默认值为 True，避免旧版兼容性。并优化部分语句的主从路由


### Bug修复

- 修复高版本 JDBC 连接低版本 MySQL 时报错 `Unknown system variable 'transaction_isolation'`
- 修复修改 namespace 旧的探活连接未关闭的问题，并限制探活连接为独立的连接
- 修复分库场景下 explain 语句报错问题
- 修复 2.3.7 版本引入的 prepare 结果集 > 16M gorm panic 问题
- 修复偶发 Cannot execute statement in a READ ONLY transaction 问题
- 修复设置 sql_mode 使用 replace 函数报错问题
- 修复忽略 /*!40103 xxx */ sql 存在语法错误情况下，执行的报错问题
- 修复keep session 事务场景和配置变更场景连接回收 bug
- 修复查询 8.0 mysql 偶发的空指针 panic 问题
- 修复 set character_set_client=binary 语法不支持问题
- 修复 mysql 8.0 字符集排序规则 utf8mb4_0900_ai_ci 设置不符合预期问题
- 修复没有流量时，p99/p95 值不更新问题
- 修复 sql 被 kill 时，客户端数据库连接未被断开问题
- 修复 executeInSlice 慢查询 panic 问题
- 修复 mysql > 8.0.3 版本 tx_read_only 报错问题
- 修复 net_buffer_size 配置不生效问题
- 修复 select last_insert_id() 返回值类型错误 string->int
- 修复 warning 日志未被清理的问题
- 修复非分片情况下多语句默认全部发送到后端 MySQL，避免多语句拆分高并发情况下负载高问题
- 修复默认支持多语句导致 QPS 高的情况下 CPU 负载升高的问题
- 修复部分 SQL 语句类型识别错误的问题
- 修复 2.0 版本部分 union 语句无法执行的问题


## Gaea 2.2 Release Notes

### 新功能

- 【SQL 支持】分片情况下支持跨分片 batch insert
- 【SQL 支持】支持 Multi Statement 多语句
- 【SQL 支持】 支持 char 函数
- 【SQL 支持】支持分片情况下 group_concat 函数
- 【SQL 支持】支持分片情况下`!mycat:sql` hint 将 SQL 打到指定分片
- 【SQL 支持】 分片情况下支持插入时忽略全局自增 ID 列或设置 null 情况下自动生成全局自增 ID
- 【SQL 支持】非分片情况下 Explain 返回 MySQL 的 Explain 信息
- 【功能】支持优先访问本地从库
- 【配置】支持日志配置热加载
- 【配置】支持分片环境下限制全局自增 ID 上限
- 【监控】新增 P99/P95 响应时间 Metric 信息
  
### 优化提升

- 【SQL 行为】分片情况下全局表随机查询分片
- 【SQL 行为】优化分片情况下默认返回一致的结果集顺序
- 【测试】优化 e2e 测试使用独立的 Docker 环境，避免环境导致 CI 失败
- 【测试】e2e 测试增加分片环境下的基础 SQL 测试
- JDBC 高版本支持：支持 JDBC 8.0.31 以上版本
- 优化获取后端 MySQL 连接超时时间，减少从库宕机时业务响应时间 Gaea集群下MySQL从库宕机流量切换测试

### Bug 修复

- 修复探活引起的负载均衡问题
- 修复分片情况下 sum 函数聚合 decimal 字段结果不准确问题
- 修复 namespace 关闭时探活未关闭的问题
- 修复非分片情况下返回 info 信息
- 修复全局表执行 order by/ group by 返回多一列的问题

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
