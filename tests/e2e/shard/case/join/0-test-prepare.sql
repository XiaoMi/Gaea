

DROP DATABASE IF EXISTS sbtest1;
CREATE DATABASE sbtest1;
USE sbtest1;

drop table if exists t,t1,t2,t3;
create table if not exists  t(id   int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
create table if not exists t1 like t;
create table if not exists t2 like t;
create table if not exists t3 like t;


drop table if exists test1;
drop table if exists test2;
drop table if exists test3;
CREATE TABLE test1(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test2(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test3(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;

drop table if exists t4;
drop table if exists t5;
drop table if exists t6;
CREATE TABLE t4(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t5(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t6(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;


drop table if exists test4,test5,test6,test7,test8,test9;
CREATE TABLE test4(ID INT NOT NULL,FirstName VARCHAR(20),LastName VARCHAR(20),Department VARCHAR(20),Salary INT,key `ID_index` (`ID`));
create table test5(id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8(id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9(id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));


insert into test1 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6);
insert into test2 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1);
insert into test3 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6);
insert into t4 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6);
insert into t5 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1);
insert into t6 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6);
INSERT INTO test4 VALUES(201,'Mazojys','Fxoj','Finance',7800),(202,'Jozzh','Lnanyo','Finance',45800),(203,'Syllauu','Dfaafk','Finance',57000),(204,'Gecrrcc','Srlkrt','Finance',62000),(205,'Jssme','Bdnaa','Development',75000),(206,'Dnnaao','Errllov','Development',55000),(207,'Tyoysww','Osk','Development',49000);
insert into test5 (id,R_REGIONKEY,R_NAME,R_COMMENT) values (1,1, 'Eastern','test001'),(3,3, 'Northern','test003'),(2,2, 'Western','test002'),(4,4, 'Southern','test004');
insert into test6 (id,C_NAME,C_NATIONKEY,C_ORDERKEY,C_CUSTKEY) values (1,'chenxiao','NATIONKEY_001','ORDERKEY_001','CUSTKEY_003'),(3,'wangye','NATIONKEY_001','ORDERKEY_004','CUSTKEY_111'),(2,'xiaojuan','NATIONKEY_001','ORDERKEY_005','CUSTKEY_132'),(4,'chenqi','NATIONKEY_051','ORDERKEY_010','CUSTKEY_333'),(5,'marui','NATIONKEY_002','ORDERKEY_011','CUSTKEY_012'),(8,'huachen','NATIONKEY_002','ORDERKEY_007','CUSTKEY_980'),(7,'yanglu','NATIONKEY_132','ORDERKEY_006','CUSTKEY_420')
insert into test7 (id,C_NAME,C_NATIONKEY,C_ORDERKEY,C_CUSTKEY) values (1,'chenxiao','NATIONKEY_001','ORDERKEY_001','CUSTKEY_003'),(3,'wangye','NATIONKEY_001','ORDERKEY_004','CUSTKEY_111'),(2,'xiaojuan','NATIONKEY_001','ORDERKEY_005','CUSTKEY_132'),(4,'chenqi','NATIONKEY_051','ORDERKEY_010','CUSTKEY_333'),(5,'marui','NATIONKEY_002','ORDERKEY_011','CUSTKEY_012'),(8,'huachen','NATIONKEY_002','ORDERKEY_007','CUSTKEY_980'),(7,'yanglu','NATIONKEY_132','ORDERKEY_006','CUSTKEY_420');
insert into test8 (id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE) values (1,'ORDERKEY_001','CUSTKEY_003',200000,'20141022'),(2,'ORDERKEY_002','CUSTKEY_003',100000,'19920501'),(4,'ORDERKEY_004','CUSTKEY_111',500,'20080105'),(5,'ORDERKEY_005','CUSTKEY_132',100,'19920628'),(10,'ORDERKEY_010','CUSTKEY_333',88888888,'19920720'),(11,'ORDERKEY_011','CUSTKEY_012',323456,'19920822'),(7,'ORDERKEY_007','CUSTKEY_980',12000,'19920910'),(6,'ORDERKEY_006','CUSTKEY_420',231,'19921111');
insert into test9 (id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE) values (1,'ORDERKEY_001','CUSTKEY_003',200000,'20141022'),(2,'ORDERKEY_002','CUSTKEY_003',100000,'19920501'),(4,'ORDERKEY_004','CUSTKEY_111',500,'20080105'),(5,'ORDERKEY_005','CUSTKEY_132',100,'19920628'),(10,'ORDERKEY_010','CUSTKEY_333',88888888,'19920720'),(11,'ORDERKEY_011','CUSTKEY_012',323456,'19920822'),(7,'ORDERKEY_007','CUSTKEY_980',12000,'19920910'),(6,'ORDERKEY_006','CUSTKEY_420',231,'19921111');
insert into t  (id, col1, col2) values  (1, 'aa', 5),(2, 'bb', 10),(3, 'cc', 15),(4, 'dd', 20),(5, 'ee', 30),(6, 'aa', 5),(7, 'bb', 10),(8, 'cc', 15),(9, 'dd', 20),(10, 'ee', 30);
insert into t1 (id, col1, col2) values  (1, 'aa', 5),(2, 'bb', 10),(3, 'cc', 15),(4, 'dd', 20),(5, 'ee', 30),(6, 'aa', 5),(7, 'bb', 10),(8, 'cc', 15),(9, 'dd', 20),(10, 'ee', 30);
insert into t2 (id, col1, col2) values  (1, 'aa', 5),(2, 'bb', 10),(3, 'cc', 15),(4, 'dd', 20),(5, 'ee', 30),(6, 'aa', 5),(7, 'bb', 10),(8, 'cc', 15),(9, 'dd', 20),(10, 'ee', 30);
insert into t3 (id, col1, col2) values  (1, 'aa', 5),(2, 'bb', 10),(3, 'cc', 15),(4, 'dd', 20),(5, 'ee', 30),(6, 'aa', 5),(7, 'bb', 10),(8, 'cc', 15),(9, 'dd', 20),(10, 'ee', 30);
