# 基本概念

* namespace

命名空间，每一个业务系统对应一个namespace，一个namespace对应多个database，业务方可以自由切换。

* slice

分片，逻辑上的分组，一个分片包含一主多从。

* shard

分表规则，确定一个表如何分表，包括分表的类型、分布的分片位置。

* proxy

指代理本身，承接线上流量。

* gaea_cc

代理控制模块，主要负责配置下发、实例监控等

* gaea_agent

部署在mysql实例所在机器，负责实例部署、管理、执行插件等功能。
