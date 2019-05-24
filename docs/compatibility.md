# Gaea兼容范围

## 协议兼容性

Gaea支持text协议和binary协议. 

## SQL兼容性

Gaea对分表和非分表的兼容性有所不同. 非分表理论上支持所有DML语句, 部分ADMIN语句.

对分表情况, Gaea本身的定位是**轻量级, 高性能**, 因此采用轻量的分表实现方式, 对一条SQL的执行, 只做字段改写和结果聚合, 不做SQL语义上的改写和多条SQL结果集的拼接计算. 

**以下支持/不支持操作均指分表情况.**

### SELECT

明确支持以下操作:

- JOIN操作支持一个父表和多个关联子表, 以及全局表.
- 聚合函数支持SUM, MAX, MIN, COUNT, 且必须出现在最外层.
- WHERE语句的条件支持AND, OR, 操作符支持=, >, >=, <, <=, <=>, IN, NOT IN, LIKE, NOT LIKE.
- 支持GROUP BY.

明确不支持以下操作:

- 不支持跨分片JOIN. JOIN中非分片键相关的条件, 只改写表名, 不计算路由, 走默认的广播路由.
- JOIN USING不支持指定表名或DB名.
- 表别名不允许与表名重复.
  - select animals.id from animals, test1.xm_order_extend as animals;
  - 这句SQL在MySQL中被认为是正确的, 但是gaea会明确拒绝这种操作.

### INSERT

明确不支持以下操作:

- 不明确指定列名的INSERT
- 跨分片批量INSERT
- INSERT INTO SELECT
 
### UPDATE

明确不支持以下操作:

- UPDATE多个表


## 事务兼容性

- Gaea目前未实现分布式事务, 只支持单分片事务, 使用跨分片事务会报错.
- 不支持SAVEPOINT, RELEASE SAVEPOINT, ROLLBACK TO SAVEPOINT **TODO**
