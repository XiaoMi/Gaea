

DROP DATABASE IF EXISTS sbtest1_2;
CREATE DATABASE sbtest1_2;
USE sbtest1_2;

drop table if exists t,t1,t2,t3;

create table if not exists  t(id  int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
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
create table test5 (id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));


DROP DATABASE IF EXISTS sbtest1_3;
CREATE DATABASE sbtest1_3;
USE sbtest1_3;


drop table if exists t,t1,t2,t3;

create table if not exists  t(id  int(11) not null auto_increment,col1 varchar(20) default null, col2 int default null, primary key (`id`),KEY `idx1` (`col1`),KEY `idx2` (`col2`) ) ENGINE=Innodb DEFAULT CHARSET UTF8MB4;
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
create table test5 (id int(11),R_REGIONKEY int(11) primary key,R_NAME varchar(50),R_COMMENT varchar(50));
create table test6(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test7(id int(11),C_NAME varchar(20),C_NATIONKEY varchar(20),C_ORDERKEY varchar(20),C_CUSTKEY varchar(20) primary key);
create table test8 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date);
create table test9 (id int(11),O_ORDERKEY varchar(20) primary key,O_CUSTKEY varchar(20),O_TOTALPRICE int(20),MYDATE date, index ORDERS_FK1 (`O_CUSTKEY`));
