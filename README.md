## 简介

Gaea是小米商城/系统组研发的基于mysql协议的数据库中间件，目前在小米商城大陆和海外得到广泛使用，包括订单、社区、活动等各个业务。Gaea支持分库分表、sql路由、读写分离等基本特性，更多详细功能可以参照功能列表，其中分库分表方案兼容了mycat和kingshard两个项目的路由方式。Gaea在设计、实现阶段参照了mycat、kingshard和vitess，并使用tidb parser作为sql parser，在此表达诚挚感谢。为了方面使用和学习Gaea，我们也提供了详细的使用和设计文档，详情参见Gaea Wiki文档。

## 自有开发模块

backend  
cmd  
log  
models  
proxy/plan  
proxy/router(按照代码行数，自有占到70%，标识为外部版权的文件很多实现都是
自己开发的与外部并不相同，开放时会加入gaea版权)  
proxy/sequence
server(session.go有函数重名和几个函数实现拷贝，最早也是move文件修改的，
目前已经非常不同，但是对应的文件也保留了对方版权，感觉应该可以去掉)  

## 外部模块

mysql(google vitess、tidb、kingshard都有引入)  
parser(tidb)  
stats(google vitess，打点统计)  
util(公用函数，也有几个自己实现的，大部分是引入的google vitess)  

