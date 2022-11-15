use sbtest1;
drop table if exists test1;
drop table if exists sbtest1.test2
drop table if exists sbtest2.test3
CREATE TABLE test1(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8
CREATE TABLE sbtest1.test2(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8
CREATE TABLE sbtest2.test3(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8
insert into test1 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6)
insert into sbtest1.test2 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1)
insert into sbtest2.test3 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6)
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1,sbtest1.test2 where test1.pad=sbtest1.test2.pad
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a,sbtest1.test2 b where a.pad=b.pad
select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id
select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id
select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest2.test3.id,sbtest2.test3.pad from test1 join sbtest2.test3 where test1.pad=sbtest2.test3.pad) b where a.pad=b.pad
select test1.id,test1.name,a.name from test1,(select name from sbtest1.test2) a
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 inner join sbtest1.test2 order by test1.id,sbtest1.test2.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 cross join sbtest1.test2 order by test1.id,sbtest1.test2.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 join sbtest1.test2 order by test1.id,sbtest1.test2.id
select a.id,a.name,a.pad,b.name from test1 a inner join sbtest1.test2 b order by a.id,b.id
select a.id,a.name,a.pad,b.name from test1 a cross join sbtest1.test2 b order by a.id,b.id
select a.id,a.name,a.pad,b.name from test1 a join sbtest1.test2 b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a inner join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a cross join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a inner join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a cross join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b  using(pad) order by a.id,b.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 straight_join sbtest1.test2 order by test1.id,sbtest1.test2.id
select a.id,a.name,a.pad,b.name from test1 a straight_join sbtest1.test2 b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a straight_join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a straight_join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a straight_join (select * from sbtest1.test2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 using(pad) order by test1.id,sbtest1.test2.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left outer join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right outer join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join sbtest1.test2 b using(pad) order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left outer join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right outer join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join (select * from sbtest1.test2 where pad>2) b using(pad) order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a right join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left outer join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a right outer join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left join (select * from sbtest1.test2 where pad>3) b using(pad) order by a.id,b.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural left join sbtest1.test2
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural right join sbtest1.test2
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural left outer join sbtest1.test2
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural right outer join sbtest1.test2
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left join sbtest1.test2 b order by a.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right join sbtest1.test2 b order by b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left outer join sbtest1.test2 b order by a.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right outer join sbtest1.test2 b order by b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left outer join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right outer join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural left join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural right join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural left outer join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural right outer join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad and test1.id>3 order by test1.id,sbtest1.test2.id
(select pad from test1) union distinct (select pad from sbtest1.test2)
(select test1.id,test1.t_id,test1.name,test1.pad from test1 where id=2) union distinct (select sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from sbtest1.test2 where id=2)
select distinct a.pad from test1 a,sbtest1.test2 b where a.pad=b.pad
select distinct b.pad,a.pad from test1 a,(select * from sbtest1.test2 where pad=1) b where a.t_id=b.o_id
select count(distinct pad,name),avg(distinct t_id) from test1
select count(distinct id),sum(distinct name) from test1 where id=3 or id=7



drop table if exists test1
drop table if exists sbtest1.test2
drop table if exists sbtest2.test3