# Gaea Changelog

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
