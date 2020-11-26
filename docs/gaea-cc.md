## 1. listNamespace

- 方法描述：查询该集群下的namespace 列表
- URL地址：/api/cc/namespace/list
- 请求方式：get
- 请求参数：

| 字段    | 类型   | 说明     | 是否必传 |
| :------ | :----- | :------- | :------- |
| cluster | string | 集群名称 | Y        |



- 返回参数

| 字段                    | 类型      | 说明     | json key    |
| :---------------------- | :-------- | :------- | :---------- |
| RetHeader               | RetHeader | 返回头   | ret_header  |
| Data                    | string[]  | 返回值   | data        |
| 此后为RetHeader对应字段 |           |          |             |
| RetCode                 | int       | 返回码   | ret_code    |
| RetMessage              | string    | 返回信息 | ret_message |



## 2.queryNamespace (弃用)

## 3.detailNamespacce

- 方法描述：查询namespace的详细配置
- URL地址：/api/cc/namespace/detail/:name
- 请求方式：get
- 请求参数：

| 字段    | 类型   | 说明                  | 是否必传 |
| :------ | :----- | :-------------------- | :------- |
| name    | string | 一次查询一个namespace | Y        |
| cluster | string |                       | Y        |



- 返回参数

| 字段                    | 类型        | 说明     | json key    |
| :---------------------- | :---------- | :------- | :---------- |
| RetHeader               | RetHeader   | 返回头   | ret_header  |
| Data                    | Namespace[] | 返回值   | data        |
| 此后为RetHeader对应字段 |             |          |             |
| RetCode                 | int         | 返回码   | ret_code    |
| RetMessage              | string      | 返回信息 | ret_message |

Namespace结构参考：https://github.com/XiaoMi/Gaea/blob/master/docs/configuration.md

## 4.modifyNamespace

- 方法描述：修改namespace配置
- URL地址：/api/cc/namespace/modify
- 请求方式：put
- 请求参数

| 字段      | 类型      | 说明                        | 是否必传 |
| :-------- | :-------- | :-------------------------- | :------- |
| cluster   | string    |                             | Y        |
| namespace | Namespace | 在body中传递namespace的json | Y        |

Namespace结构参考：https://github.com/XiaoMi/Gaea/blob/master/docs/configuration.md

- 返回参数

| 字段       | 类型   | 说明     | json key    |
| :--------- | :----- | :------- | :---------- |
| RetCode    | int    | 返回码   | ret_code    |
| RetMessage | string | 返回信息 | ret_message |



## 5.delNamespace

- 方法描述：删除namespace 
- URL地址：/api/cc/namespace/delete/:name
- 请求方式：put
- 请求参数

| 字段    | 类型   | 说明                  | 是否必传 |
| :------ | :----- | :-------------------- | :------- |
| name    | string | 一次删除一个namespace | Y        |
| cluster | string |                       | Y        |



- 返回参数

| 字段       | 类型   | 说明     | json key    |
| :--------- | :----- | :------- | :---------- |
| RetCode    | int    | 返回码   | ret_code    |
| RetMessage | string | 返回信息 | ret_message |



## 6.sqlFingerprint

- 方法描述：获取慢sql , 错误sql 指纹
- URL地址：/api/cc/namespace/sqlfingerprint/:name
- 请求方式：get
- 请求参数

| 字段    | 类型   | 说明          | 是否必传 |
| :------ | :----- | :------------ | :------- |
| name    | string | namespace名称 | Y        |
| cluster | string | 集群名称      | Y        |



- 返回参数

| 字段                    | 类型              | 说明     | json key    |
| :---------------------- | :---------------- | :------- | :---------- |
| RetHeader               | RetHeader         | 返回头   | ret_header  |
| ErrSQLs                 | map[string]string |          | err_sqls    |
| SlowSQLs                | map[string]string |          | slow_sqls   |
| 此后为RetHeader对应字段 |                   |          |             |
| RetCode                 | int               | 返回码   | ret_code    |
| RetMessage              | string            | 返回信息 | ret_message |



## 7.proxyConfigFingerprint

- 方法描述：获取该集群下所有proxy配置的md5值
- URL地址：/api/cc/proxy/config/fingerprint
- 请求方式：get
- 请求参数

| 字段    | 类型   | 说明 | 是否必传 |
| :------ | :----- | :--- | :------- |
| cluster | string |      | Y        |



- 返回参数

| 字段                    | 类型              | 说明                                  | json key    |
| :---------------------- | :---------------- | :------------------------------------ | :---------- |
| RetHeader               | RetHeader         | 返回头                                | ret_header  |
| Data                    | map[string]string | key: proxy-ip:portvalue:md5 of config | data        |
| 此后为RetHeader对应字段 |                   |                                       |             |
| RetCode                 | int               | 返回码                                | ret_code    |
| RetMessage              | string            | 返回信息                              | ret_message |