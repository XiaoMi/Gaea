## 简介

Gaea是小米商城/系统组研发的基于mysql协议的数据库中间件，目前在小米商城大陆和海外得到广泛使用，包括订单、社区、活动等各个业务。Gaea支持分库分表、sql路由、读写分离等基本特性，更多详细功能可以参照功能列表，其中分库分表方案兼容了mycat和kingshard两个项目的路由方式。Gaea在设计、实现阶段参照了mycat、kingshard和vitess，并使用tidb parser作为sql parser，在此表达诚挚感谢。为了方面使用和学习Gaea，我们也提供了详细的使用和设计文档，详情参见Gaea Wiki文档。

## 自有开发模块

backend  
cmd  
log  
models  
proxy/plan  
proxy/router(kingshard路由方式源自kingshard项目本身)
proxy/sequence
server

## 外部模块

mysql(google vitess、tidb、kingshard都有引入)  
parser(tidb)  
stats(google vitess，打点统计)  
util(公用函数，也有几个自己实现的，大部分是引入的google vitess)  

