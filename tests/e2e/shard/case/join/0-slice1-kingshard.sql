CREATE DATABASE  sbtest;
USE sbtest;

drop table if exists t_0002,t1_0002,t2_0002,t3_0002;
create table if not exists  t_0002(id  int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
create table if not exists t1_0002 like t_0002;
create table if not exists t2_0002 like t_0002;
create table if not exists t3_0002 like t_0002;

drop table if exists test1_0002;
drop table if exists test2_0002;
drop table if exists test3_0002;
CREATE TABLE test1_0002(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test2_0002(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test3_0002(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;

drop table if exists t4_0002;
drop table if exists t5_0002;
drop table if exists t6_0002;
CREATE TABLE t4_0002(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t5_0002(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t6_0002(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;

drop table if exists test4_0002,test5_0002,test6_0002,test7_0002,test8_0002,test9_0002;
CREATE TABLE test4_0002(ID INT NOT NULL,FirstName VARCHAR(20),LastName VARCHAR(20),Department VARCHAR(20),Salary INT,key `ID_index` (`ID`));
create table test5_0002 (id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6_0002(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7_0002(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8_0002 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9_0002 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));


drop table if exists t_0003,t1_0003,t2_0003,t3_0003;
create table if not exists  t_0003(id  int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
create table if not exists t1_0003 like t_0003;
create table if not exists t2_0003 like t_0003;
create table if not exists t3_0003 like t_0003;


drop table if exists test1_0003;
drop table if exists test2_0003;
drop table if exists test3_0003;
CREATE TABLE test1_0003(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test2_0003(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test3_0003(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;


drop table if exists t4_0003;
drop table if exists t5_0003;
drop table if exists t6_0003;
CREATE TABLE t4_0003(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t5_0003(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t6_0003(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;

drop table if exists test4_0003,test5_0003,test6_0003,test7_0003,test8_0003,test9_0003;
CREATE TABLE test4_0003(ID INT NOT NULL,FirstName VARCHAR(20),LastName VARCHAR(20),Department VARCHAR(20),Salary INT,key `ID_index` (`ID`));
create table test5_0003 (id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6_0003(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7_0003(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8_0003 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9_0003 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));
