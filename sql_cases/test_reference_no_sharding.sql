drop table if exists noshard_t1
drop table if exists noshard_t2
drop table if exists noshard_t3
CREATE TABLE noshard_t1(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8
CREATE TABLE noshard_t2(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8
CREATE TABLE noshard_t3(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8
insert into noshard_t1 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6)
insert into noshard_t2 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1)
insert into noshard_t3 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6)
select * from noshard_t1 a,noshard_t2 b
select * from (select * from noshard_t1 where id<3) a,(select * from noshard_t2 where id>3) b
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b,(select * from noshard_t2 where id>3) c where a.pad=b.pad and c.pad=b.pad
    #
#join table
#
select * from noshard_t1 a join noshard_t2 as b order by a.id,b.id
select * from noshard_t1 a inner join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a cross join noshard_t2 b order by a.id,b.id
select a.id,a.name,a.pad,b.name from noshard_t1 a straight_join noshard_t2 b on a.pad=b.pad
    #
    #SELECT ... UNION [ALL | DISTINCT] SELECT ... [UNION [ALL | DISTINCT] SELECT ...]
    #
select * from noshard_t1 union all select * from noshard_t3 union all select * from noshard_t2
select * from noshard_t1 union distinct select * from noshard_t3 union distinct select * from noshard_t2
    (select name from noshard_t1 where pad=1 order by id limit 10) union all (select name from noshard_t2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
    (select name from noshard_t1 where pad=1 order by id limit 10) union distinct (select name from noshard_t2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
    (select * from noshard_t1 where pad=1) union (select * from noshard_t2 where pad=1) order by name limit 10
    (select name as sort_a from noshard_t1 where pad=1) union (select name from noshard_t2 where pad=1) order by sort_a limit 10
    (select name as sort_a,pad from noshard_t1 where pad=1) union (select name,pad from noshard_t2 where pad=1) order by sort_a,pad limit 10
    #
    #clear tables
    #
drop table if exists noshard_t1
drop table if exists noshard_t2
drop table if exists noshard_t3