CREATE DATABASE sbtest;
USE sbtest;

drop table if exists  t_0000,t1_0000,t2_0000,t3_0000;
create table if not exists  t_0000(id  int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
create table if not exists t1_0000 like t_0000;
create table if not exists t2_0000 like t_0000;
create table if not exists t3_0000 like t_0000;


drop table if exists test1_0000;
drop table if exists test2_0000;
drop table if exists test3_0000;
CREATE TABLE test1_0000(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test2_0000(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test3_0000(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;


drop table if exists t4_0000;
drop table if exists t5_0000;
drop table if exists t6_0000;
CREATE TABLE t4_0000(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t5_0000(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t6_0000(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;

drop table if exists test4_0000,test5_0000,test6_0000,test7_0000,test8_0000,test9_0000;
CREATE TABLE test4_0000(ID INT NOT NULL,FirstName VARCHAR(20),LastName VARCHAR(20),Department VARCHAR(20),Salary INT,key `ID_index` (`ID`));
create table test5_0000 (id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6_0000(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7_0000(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8_0000 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9_0000 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));



drop table if exists t_0001,t1_0001,t2_0001,t3_0001;
create table if not exists  t_0001(id  int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
create table if not exists t1_0001 like t_0001;
create table if not exists t2_0001 like t_0001;
create table if not exists t3_0001 like t_0001;


drop table if exists test1_0001;
drop table if exists test2_0001;
drop table if exists test3_0001;
CREATE TABLE test1_0001(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test2_0001(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE test3_0001(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;


drop table if exists t4_0001;
drop table if exists t5_0001;
drop table if exists t6_0001;
CREATE TABLE t4_0001(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t5_0001(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE t6_0001(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;

drop table if exists test4_0001,test5_0001,test6_0001,test7_0001,test8_0001,test9_0001;
CREATE TABLE test4_0001(ID INT NOT NULL,FirstName VARCHAR(20),LastName VARCHAR(20),Department VARCHAR(20),Salary INT,key `ID_index` (`ID`));
create table test5_0001 (id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6_0001(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7_0001(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8_0001 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9_0001 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));
