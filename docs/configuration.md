# 配置说明

gaea配置由两部分组成，本地配置为gaea_proxy直接使用的配置内容，一般不需要在运行时改变。gaea为多租户模式，每个租户称为一个namespace，namespace 的配置在运行时都可变，一般保存在etcd中。

## 本地配置说明

```ini
; 配置类型，目前支持file/etcd两种方式，file方式不支持热加载，但是可以快速体验功能
config_type=etcd

;file config path, 具体配置放到file_config_path的namespace目录下，该下级目录为固定目录
file_config_path=./etc/file

;配置中心地址，目前只支持etcd
coordinator_addr=http://127.0.0.1:2379

;远程配置(当前为etcd)根目录
coordinator_root=/gaea

;配置中心用户名和密码
username=test
password=test

;环境划分、test、online
environ=test

;组和服务名称，为服务化做准备
group_name=systech

;service name
service_name=gaea_proxy

;日志配置
log_path=./logs
log_level=Notice
log_filename=gaea
log_output=file

;管理地址
admin_addr=0.0.0.0:13307
;basic auth
admin_user=admin
admin_password=admin

;代理服务监听地址
proto_type=tcp4
proxy_addr=0.0.0.0:13306

; 默认编码
proxy_charset=utf8

;慢sql阈值，单位: 毫秒
slow_sql_time=100

;空闲会话超时时间,单位: 秒
session_timeout=3600

;打点统计配置
stats_enabled=true
stats_backend_type=prometheus
```

## namespace配置说明

namespace的配置格式为json，包含分表、非分表、实例等配置信息，都可在运行时改变。namespace的配置可以直接通过web平台进行操作，使用方不需要关心json里的内容，如果有兴趣参与到gaea的开发中，可以关注下字段含义，具体解释如下,格式为字段名称、类型、内容含义。

| 字段名称         | 字段类型   | 字段含义                                           |
| --------------- | ---------- | ----------------------------------------------- |
| name            | string     | namespace名称                                    |
| online          | bool       | 是否在线，逻辑上下线使用                            |
| read_only       | bool       | 是否只读，namespace级别                            |
| allowed_dbs     | map        | 数据库集合                                        |
| default_phy_dbs | map        | 默认数据库名, 与allowed_dbs一一对应                 |
| slow_sql_time   | string     | 慢sql时间，单位ms                                 |
| black_sql       | string数组 | 黑名单sql                                         |
| allowed_ip      | string数组 | 白名单IP                                          |
| slices          | map数组    | 一主多从的物理实例，slice里map的具体字段可参照slice配置 |
| shard_rules     | map数组    | 分库、分表、特殊表的配置内容，具体字段可参照shard配置    |
| users           | map数组    | 应用端连接gaea所需要的用户配置，具体字段可参照users配置 |

### slice配置

| 字段名称         | 字段类型   | 字段含义                                       |
| ---------------- | ---------- | ---------------------------------------------- |
| name             | string     | 分片名称，自动、有序生成                       |
| user_name        | string     | 连接后端mysql所需要的用户名称                  |
| password         | string     | 连接后端mysql所需要的用户密码                  |
| master           | string     | 主实例地址                                     |
| slaves           | string数组 | 从实例地址列表                                 |
| statistic_slaves | string数组 | 统计型从实例地址列表                           |
| capacity         | int        | gaea_proxy与每个实例的连接池大小               |
| max_capacity     | int        | gaea_proxy与每个实例的连接池最大大小           |
| idle_timeout     | int        | gaea_proxy与后端mysql空闲连接存活时间，单位:秒 |

### shard配置

这里列出了一些基本配置参数, 详细配置请参考[分片表配置](shard.md)

| 字段名称   | 字段类型 | 字段含义 |
| --------- | -------- | --------------------- |
| db        | string   | 分片表所在DB            |
| table     | string   | 分片表名                |
| type      | string   | 分片类型                |
| key       | string   | 分片列名                |
| locations | list     | 每个slice上分布的分片个数 |
| slices    | list     | slice列表              |
| databases | list     | mycat分片规则后端实际DB名 |

### users配置

| 字段名称       | 字段类型 | 字段含义                               |
| -------------- | -------- | -------------------------------------- |
| user_name      | string   | 用户名                                 |
| password       | string   | 用户密码                               |
| namespace      | string   | 对应的命名空间                         |
| rw_flag        | int      | 读写标识, 只读=1, 读写=2                |
| rw_split       | int      | 是否读写分离, 非读写分离=0, 读写分离=1     |
| other_property | int      | 目前用来标识是否走统计从实例, 普通用户=0, 统计用户=1 |
