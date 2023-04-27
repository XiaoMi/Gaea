# FAQ

## 权限说明

部分语句默认的请求到后端从库或主库的情况：
- 普通 Select：
    - 读写分离用户（RWFlag=2，RWSplit=1）：从库
    - 只写用户（RWFlag=2，RWSplit=0）：主库
    - 只读用户（RWFlag=1，RWSplit=1）：从库
- Select 带 Master Hint：
    - 读写分离用户（RWFlag=2，RWSplit=1）：主库
    - 只写用户（RWFlag=2，RWSplit=0）：主库
    - 只读用户（RWFlag=1，RWSplit=1）：从库（V2.0 以下版本会请求到主库，MiProxy 会打到从库）
- Select... For Update(lock in share mode):（同上）
    - 读写分离用户（RWFlag=2，RWSplit=1）：主库
    - 只写用户（RWFlag=2，RWSplit=0）：主库
    - 只读用户（RWFlag=1，RWSplit=1）：从库
- Update/Insert/Delete:
    - 读写分离用户（RWFlag=2，RWSplit=1）：主库
    - 只写用户（RWFlag=2，RWSplit=0）：主库
    - 只读用户（RWFlag=1，RWSplit=1）：报错
- 开启事务
    - 读写分离用户（RWFlag=2，RWSplit=1）：主库（与旧版相同）
    - 只写用户（RWFlag=2，RWSplit=0）：主库（与旧版相同）
    - 只读用户（RWFlag=1，RWSplit=1）：主库（与旧版相同）

- 当从库宕机时：会重新从主库获取连接，并请求到主库（无论权限如何）
- 当 Select 处于显式开启的事务中，也都会请求到主库（包括只读用户，此处与 MiProxy 不一致，MiProxy 只读用户开启事务也会请求到从库）

## 配置热加载
目前 Gaea namespace 配置支持热加载，Gaea 本身的配置文件只支持日志的热加载，执行热加载的方式为
```bash
# 使用 signal 重新加载
kill -SIGUSER1 ${gaea_pid}
# 使用 API 热加载
curl -X PUT 'http://127.0.0.1:13307/api/proxy/proxyconfig/reload' \
-H 'Authorization: Basic YWRtaW46YWRtaW4='
```