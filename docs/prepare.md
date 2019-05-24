# prepare的设计与实现

## 背景

应用端使用prepare主要考虑通过固定sql模板，在执行sql时只传输参数，减少数据包传输大小，提升sql执行效率。对于非分库分表的情况，我们可以直接通过转发execute(对应后端连接可能是prepare+execute+close)的方式进行支持。但是对于分库分表的情形，需要计算路由，重写sql，支持起来会非常麻烦。商城目前的分库分表中间件是mycat，而mycat是支持prepare的，而gaea的prepare方案也是参照mycat，即将prepare statements的执行转换为sql的执行，然后在应答阶段，根据文本的应答内容构造二进制应答内容，返回给客户端，从而统一了分库分表的处理逻辑。

## prepare

gaea在接到preprae请求后，首先计算参数个数、参数偏移位置和stmt-id。然后根据以上数据，构造statement对象并保存在SessionExecutor的stmts内，stmts为一个map，key为stmt-id，value即为构造的statement对象。

prepare阶段主要是计算、保存execute需要使用的变量信息，prepare应答数据内也会包含这些变量信息。

## execute

execute请求时会携带prepare应答返回的stmt-id，服务端根据stmt-id从SessionExecutor的stmts中查询对应的statement信息。根据statement信息的参数个数、偏移和execute上传的参数值，进行关联绑定，然后rewrite一条同等含义的sql。同时，为了安全性考虑，也会进行特殊字符过滤，防止比如sql注入的发生。

生成sql之后，无论是分表还是非分表，我们都可以调用handleQuery进行统一的处理，避免了因为要支持prepare，而存在两套计算分库、分表路由的逻辑。

处理完成之后，需要进行文本应答协议到二进制应答协议的转换，相关实现在BuildBinaryResultset内。

execute执行完成之后，执行ResetParams，重新初始化send_long_data对应的args字段，病返回应答。

## send_long_data

send_long_data不是必须的，但是如果execute有多个参数，且不止一个参数长度比较大，一次execute可能达到mysql max-payload-length，但是如果分多次，每次只发送一个，这样就绕过了max-payload-length，send_long_data就是基于这样的背景产生的。

客户端发送send_long_data报文，会携带stmt-id、param-id(参数位置)，我们根据stmt-id参数可以检索prepare阶段存储的stmt信息，根据param-id和对应上送的数据，可以建立一个k-v映射，存储在Stmt.args(interface slice)。在上述execute执行阶段，也会根据参数位置从Stmt.args查询对应位置的参数值，进行关联绑定。

send_long_data不需要应答。

## close

close的处理逻辑比较简单，服务端收到close请求后，删除prepare阶段stmt-id及其数据的对应关系。

## 总结

gaea对于prepare的处理初衷还是考虑协议的兼容和简化处理逻辑，对于client->proxy->mysql这样接口来说，client->proxy是prepare协议，proxy->mysql是文本协议，所以整体来看在gaea环境下使用prepare性能提升有限，还是建议直接使用sql。

## 参考资料

[mysql prepare statements 官方文档](https://dev.mysql.com/doc/internals/en/prepared-statements.html)