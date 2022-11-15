drop table if exists noshard_t1
drop table if exists noshard_t2
drop table if exists noshard_t3
CREATE TABLE noshard_t1(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8
CREATE TABLE noshard_t2(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8
CREATE TABLE noshard_t3(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8
insert into noshard_t1 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6)
insert into noshard_t2 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1)
insert into noshard_t3 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6)
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select all * from noshard_t2) b where a.t_id=b.o_id;
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select distinct * from noshard_t2) b where a.t_id=b.o_id;
select * from (select * from noshard_t2 a group by a.id) a;
select * from (select pad,count(*) from noshard_t2 a group by pad) a;
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 having pad>3) b where a.t_id=b.o_id;
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3 order by id) b where a.t_id=b.o_id;
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 order by id limit 3) b where a.t_id=b.o_id;
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 order by id limit 3) b where a.t_id=b.o_id limit 2;
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3) b where a.t_id=b.o_id;
select * from (select noshard_t2.pad from noshard_t1 left join noshard_t2 on noshard_t1.pad=noshard_t2.pad) a;
select * from (select * from noshard_t1 union select * from noshard_t2) a where a.id >3;
select id,pad from noshard_t1 where pad=(select min(id) from noshard_t2);
select id,pad,name from (select * from noshard_t1 where pad>2) a where id<5;
select pad,count(*) from (select * from noshard_t1 where pad>2) a group by pad;
select pad,count(*) from (select * from noshard_t1 where pad>2) a group by pad order by pad;
select count(*) from (select pad,count(*) a from noshard_t1 group by pad) a;
select * from noshard_t1 where pad<(select pad from noshard_t2 where id=3);
select * from noshard_t1 having pad<(select pad from noshard_t2 where id=3);
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3) b where a.t_id=b.o_id;
select id,name,(select count(*) from noshard_t2) count from noshard_t1;
select * from noshard_t1 where pad like (select pad from noshard_t2 where id=3);
select id,pad from noshard_t1 where pad>(select pad from noshard_t1 where id=2);
select id,pad from noshard_t1 where pad<(select pad from noshard_t1 where id=2);
select id,pad from noshard_t1 where pad=(select pad from noshard_t1 where id=2);
select id,pad from noshard_t1 where pad>=(select pad from noshard_t1 where id=2);
select id,pad from noshard_t1 where pad<=(select pad from noshard_t1 where id=2);
select id,pad from noshard_t1 where pad<>(select pad from noshard_t1 where id=2);
select id,pad from noshard_t1 where pad !=(select pad from noshard_t1 where id=2);
select * from noshard_t1 where exists(select * from noshard_t1 where pad>1);
select * from noshard_t1 where not exists(select * from noshard_t1 where pad>1);
select * from noshard_t1 where pad not in(select id from noshard_t1 where pad>1);
select * from noshard_t1 where pad in(select id from noshard_t1 where pad>1);
select * from noshard_t1 where pad=some(select id from noshard_t1 where pad>1);
select * from noshard_t1 where pad=any(select id from noshard_t1 where pad>1);
select * from noshard_t1 where pad !=any(select id from noshard_t1 where pad=3);
select a.id,b.id,b.pad,a.t_id from (select noshard_t1.id,noshard_t1.pad,noshard_t1.t_id from noshard_t1 join noshard_t2 where noshard_t1.pad=noshard_t2.pad ) a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b where a.pad=b.pad;
select * from noshard_t1 where pad>(select pad from noshard_t1 where pad=2);
select * from noshard_t1,(select * from noshard_t1 where id>3 union select * from noshard_t2 where id<2) a where a.id >3 and noshard_t1.pad=a.pad;
select count(*) from (select * from noshard_t1 where pad=(select pad from noshard_t2 where id=1)) a;
#
#Second supplement
#
select (select name from noshard_t1 limit 1)
select * from noshard_t1 where 'test_2'=(select name from noshard_t2 where id=2)
select * from noshard_t1 where 5=(select count(*) from noshard_t2)
select * from noshard_t1 where 'test_2' like(select name from noshard_t2 where id=2)
select * from noshard_t1 where 2 >any(select id from noshard_t1 where pad>1)
select * from noshard_t1 where 2 in(select id from noshard_t1 where pad>1)
select * from noshard_t1 where 2<>some(select id from noshard_t1 where pad>1)
select * from noshard_t1 where 2>all(select id from noshard_t1 where pad<1)
select * from noshard_t1 where (id,pad)=(select id,pad from noshard_t2 limit 1)
select * from noshard_t1 where row(id,pad)=(select id,pad from noshard_t2 limit 1)
select id,name,pad from noshard_t1 where (id,pad)in(select id,pad from noshard_t2)
select id,name,pad from noshard_t1 where (1,1)in(select id,pad from noshard_t2)
SELECT pad FROM noshard_t1 AS x WHERE x.id = (SELECT pad FROM noshard_t2 AS y WHERE x.id = (SELECT pad FROM noshard_t3 WHERE y.id = noshard_t3.id))
select co1,co2,co3 from (select id as co1,name as co2,pad as co3 from noshard_t1)as tb where co1>1
select avg(sum_column1) from (select sum(id) as sum_column1 from noshard_t1 group by pad) as t1
    #
#clear tables
#
drop table if exists noshard_t1
drop table if exists noshard_t2
drop table if exists noshard_t3