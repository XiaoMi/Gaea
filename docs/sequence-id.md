# 全局序列号说明

## 原理

参考mycat生成全局唯一序列号的设计，在数据库中建立一张表，存放sequence名称(name)，sequence当前值(current_value)，步长(increment int类型每次读取多少个sequence，假设为K)等信息；

Sequence获取步骤：

1).第一次使用该sequence时，根据传入的sequence名称，从数据库这张表中读取current_value，和increment到gaea中，并将数据库中的current_value设置为原current_value值+increment值（实现方式是基于后续的存储函数）

2).gaea将读取到current_value+increment作为本次要使用的sequence值，下次使用时，自动加1，当使用increment次后，执行步骤1)相同的操作.    
  
gaea只会修改和查询这张表，使用前需要按照下文中配置table的步骤在这张表中插入一条记录。   
若某次读取的sequence没有用完，系统就停掉了，则这次读取的sequence剩余值不会再使用。

## 如何使用

如：tbl_user_info的id列使用全局自增序列号
``` 
insert into gaea_test.tbl_user_info set name="zhangsan", age=15, id = nextval();
```

## 如何配置 
在指定的slice的master上操作，完成以下配置。
### 配置db
db: mycat
```
create database mycat;
```
### 配置table
table: mycat_sequence
``` 
DROP TABLE IF EXISTS MYCAT_SEQUENCE;
CREATE TABLE MYCAT_SEQUENCE (name VARCHAR(50) NOT NULL,current_value INT NOT NULL,increment INT NOT NULL DEFAULT 100, PRIMARY KEY(name)) ENGINE=InnoDB;
```

### 初始化table
``` 
INSERT INTO MYCAT_SEQUENCE(name,current_value,increment) VALUES ('GAEA_TEST.TBL_USER_INFO', -99, 100);
```
name由大写的db名称和table名称组成  
increment可根据参考业务qps峰值配置  
例：  
如果要在database gaea_test里的tbl_user_info表中使用全局唯一序列号  
name=>GAEA_TEST.TBL_USER_INFO  
``` 
mysql> select * from mycat_sequence;
+-------------------------+---------------+-----------+
| name                    | current_value | increment |
+-------------------------+---------------+-----------+
| GAEA_TEST.TBL_USER_INFO |           700 |       100 |
+-------------------------+---------------+-----------+
```

### 配置函数
- 获取当前sequence的值： 
``` 
DROP FUNCTION IF EXISTS mycat_seq_currval;
DELIMITER $
CREATE FUNCTION mycat_seq_currval(seq_name VARCHAR(50)) RETURNS varchar(64)     CHARSET utf8
DETERMINISTIC
BEGIN
DECLARE retval VARCHAR(64);
SET retval="-999999999,null";
SELECT concat(CAST(current_value AS CHAR),",",CAST(increment AS CHAR)) INTO retval FROM MYCAT_SEQUENCE WHERE name = seq_name;
RETURN retval;
END $
DELIMITER ;
```

- 设置sequence值： 
```

DROP FUNCTION IF EXISTS mycat_seq_setval; 
DELIMITER $ 
CREATE FUNCTION mycat_seq_setval(seq_name VARCHAR(50),value INTEGER) RETURNS     varchar(64) CHARSET utf8
DETERMINISTIC
BEGIN
UPDATE MYCAT_SEQUENCE
SET current_value = value
WHERE name = seq_name;
RETURN mycat_seq_currval(seq_name);
END $
DELIMITER ;
```

- 获取下一个sequence值：
 
```
DROP FUNCTION IF EXISTS mycat_seq_nextval;
DELIMITER $
CREATE FUNCTION mycat_seq_nextval(seq_name VARCHAR(50)) RETURNS varchar(64)     CHARSET utf8
DETERMINISTIC
BEGIN
UPDATE MYCAT_SEQUENCE
SET current_value = current_value + increment WHERE name = seq_name;
RETURN mycat_seq_currval(seq_name);
END $
DELIMITER ;
```

