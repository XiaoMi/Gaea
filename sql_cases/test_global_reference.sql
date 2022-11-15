drop table if exists test1
drop table if exists sbtest1.test2
drop table if exists sbtest2.test3
CREATE TABLE test1(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8
CREATE TABLE sbtest1.test2(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8
CREATE TABLE sbtest2.test3(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8
insert into test1 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6)
insert into sbtest1.test2 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1)
insert into sbtest2.test3 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6)

select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b
select a.id,b.id,b.pad,a.t_id from test1 a,(select sbtest2.test3.id,sbtest2.test3.pad from test1 join sbtest2.test3 where test1.pad=sbtest2.test3.pad) b,(select * from sbtest1.test2 where id>3) c where a.pad=b.pad and c.pad=b.pad
    #
#join table
#
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join sbtest1.test2 as b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a inner join sbtest1.test2 b order by a.id,b.id
select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a cross join sbtest1.test2 b order by a.id,b.id
select a.id,a.name,a.pad,b.name from test1 a straight_join sbtest1.test2 b on a.pad=b.pad
    #
    #SELECT ... UNION [ALL | DISTINCT] SELECT ... [UNION [ALL | DISTINCT] SELECT ...]
    #
select a.id,a.t_id,a.name,a.pad from test1 a union all select b.id,b.m_id,b.name,b.pad from sbtest2.test3 b union all select c.id,c.o_id,c.name,c.pad from sbtest1.test2 c
select a.id,a.t_id,a.name,a.pad from test1 a union distinct select b.id,b.m_id,b.name,b.pad from sbtest2.test3 b union distinct select c.id,c.o_id,c.name,c.pad from sbtest1.test2 c
    (select name from test1 where pad=1 order by id limit 10) union all (select name from sbtest1.test2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
    (select name from test1 where pad=1 order by id limit 10) union distinct (select name from sbtest1.test2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
    (select a.id,a.t_id,a.name,a.pad from test1 a where a.pad=1) union (select c.id,c.o_id,c.name,c.pad from sbtest1.test2 c where c.pad=1) order by id limit 10/*allow_diff_sequence*/
    (select name as sort_a from test1 where pad=1) union (select name from sbtest1.test2 where pad=1) order by sort_a limit 10/*allow_diff_sequence*/
    (select name as sort_a,pad from test1 where pad=1) union (select name,pad from sbtest1.test2 where pad=1) order by sort_a,pad limit 10/*allow_diff_sequence*/
    #-- case union,from issue #275
    (select * from test1 where id=2) union (select * from sbtest1.test2 where id=2);
#clear tables
#
drop table if exists test1
drop table if exists sbtest1.test2
drop table if exists sbtest2.test3