# 基本概念

* cluster

集群, 按照业务重要程度划分集群, 一个集群内可包含多个gaea-proxy实例, 通过指定gaea-proxy启动时依赖的配置文件中的cluster_name
确定该proxy所属集群。集群内的proxy实例只为该集群内的namespace提供服务, 起到物理隔离的作用。
一套集群可为多个namespace提供服务。

* namespace

命名空间，每一个业务系统对应一个namespace，一个namespace对应多个database，业务方可以自由切换。
每个namespace理论上只会属于一个集群。
通过gaea-cc配置管理接口, 指定namespace所属集群。

* slice

分片，逻辑上的分组，一个分片包含mysql一主多从。

* shard

分表规则，确定一个表如何分表，包括分表的类型、分布的分片位置。

* proxy

指代理本身，承接线上流量。

* gaea_cc

代理控制模块，主要负责配置下发、实例监控等。

* gaea_agent

部署在mysql实例所在机器，负责实例部署、管理、执行插件等功能。
