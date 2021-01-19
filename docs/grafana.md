# 监控配置

gaea proxy基于prometheus做统计数据的存储，使用grafana实现数据可视化展示。

## 监控说明

### gaea proxy监控

[proxy grafana配置模板](template/gaea_proxy.json)  

proxy监控概览这部分主要展示gaea proxy的整体运行情况，主要包含以下几个监控项:

- 集群QPS
- 业务流量
- 业务请求耗时
- SQL 错误数
- CPU 负载
- 内存 负载
- 流量负载
- 会话数
- 业务会话数
- 协程数量
- GC 停顿时间
- 堆对象数量
   

### 租户各指标监控

[namespace grafana配置模板](template/gaea_namespace.json)

导入模板之前需要把模板里的gaea_test_namespace 替换为实际使用的namespace

租户指标监控主要展示某个namespace的统计数据，主要包含以下几个监控项:

- QPS
- 流量
- SQL耗时
- SQL错误数
- 高耗时SQL指纹
- 错误SQL指纹
- 连接数
- 空闲连接数
- 连接等待队列


## prometheus配置说明

```
- job_name: 'gaea_proxy'
    metrics_path: '/api/metric/metrics'
    static_configs:
    - targets: ["admin_addr1"]
    - targets: ["admin_addr2"]
    - targets: ["admin_addr3"]
    basic_auth:
      username: admin_user
      password: admin_password
```
需要修改admin_addr,admin_user,admin_password与gaea.ini中的以下几项保持一致。
```
;管理地址
admin_addr=0.0.0.0:13307
;basic auth
admin_user=admin
admin_password=admin
```
增加Prometheus Recoding Rules
```
groups:
  - name: gaea_proxy_rule
    rules:
    - record: gaea_proxy_sql_timings_count_rate_each_namespace
      expr: sum(avg(rate(gaea_proxy_sql_timings_count[20s])) without (slave)) by (namespace)
    - record: gaea_proxy_sql_timings_count_rate_total
      expr: sum(sum(avg(rate(gaea_proxy_sql_timings_count[20s])) without (slave)) by (namespace))
    - record: gaea_proxy_flow_counts_rate_namespace_flowdirection
      expr: sum(avg(rate(gaea_proxy_flow_counts[20s])) without (slave)) by (namespace, flowdirection)
    - record: gaea_proxy_flow_counts_rate_namespace
      expr: sum(avg(rate(gaea_proxy_flow_counts[20s])) without (slave)) by (namespace)
    - record: gaea_proxy_flow_counts_rate_total
      expr: sum(sum(avg(rate(gaea_proxy_flow_counts[20s])) without (slave)) by (namespace))
    - record: gaea_proxy_sql_timings_rate_namespace_operation
      expr: sum(delta(gaea_proxy_sql_timings_sum[20s])) by (namespace,operation) / sum(delta(gaea_proxy_sql_timings_count[20s])) by (namespace,operation)
    - record: gaea_proxy_sql_timings_rate_namespace
      expr: sum(delta(gaea_proxy_sql_timings_sum[20s])) by (namespace) / sum(delta(gaea_proxy_sql_timings_count[20s])) by (namespace)
    - record: gaea_proxy_sql_error_counts_rate_namespace
      expr: sum(avg(rate(gaea_proxy_sql_error_counts[20s])) without (instance)) by (namespace)
``` 
##  

 
