//Attention, You may need to name the schame sql files with 'schema' prefix and thus ensure schema files will
//be executed first


create database if not exists sbtest1;
drop table sbtest1.t1;
drop table sbtest1.t2;
create table if not exists sbtest1.t1(id int unsigned primary key auto_increment, name varchar(20));
insert into sbtest1.t1(name) values('x'),('y'),('z'),('zhangsan'),('lisi');

create table if not exists sbtest1.t2(id int unsigned primary key auto_increment, name varchar(20));
insert into sbtest1.t1(name) values('x'),('y'),('z'),('zhangsan'),('lisi');

