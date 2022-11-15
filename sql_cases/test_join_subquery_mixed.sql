drop table if exists noshard_t1;
drop table if exists noshard_t2;
drop table if exists noshard_t3;
CREATE TABLE noshard_t1(`id` int(10) unsigned NOT NULL,`t_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`t_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE noshard_t2(`id` int(10) unsigned NOT NULL,`o_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`o_id`))DEFAULT CHARSET=UTF8;
CREATE TABLE noshard_t3(`id` int(10) unsigned NOT NULL,`m_id` int(10) unsigned NOT NULL DEFAULT '0',`name` char(120) NOT NULL DEFAULT '',`pad` int(11) NOT NULL,PRIMARY KEY (`id`),KEY `k_1` (`m_id`))DEFAULT CHARSET=UTF8;
insert into noshard_t1 values(1,1,'test中id为1',1),(2,2,'test_2',2),(3,3,'test中id为3',4),(4,4,'$test$4',3),(5,5,'test...5',1),(6,6,'test6',6);
insert into noshard_t2 values(1,1,'order中id为1',1),(2,2,'test_2',2),(3,3,'order中id为3',3),(4,4,'$order$4',4),(5,5,'order...5',1);
insert into noshard_t3 values(1,1,'manager中id为1',1),(2,2,'test_2',2),(3,3,'manager中id为3',3),(4,4,'$manager$4',4),(5,5,'manager...5',6);

select * from (select a.id, t_id, a.name, a.pad from noshard_t1 a join noshard_t3 b on a.id = b.id union select * from  noshard_t2) as c order by id,name;
select id, t_id, name,pad from (select a.id,t_id, b.name,b.pad from noshard_t1 a join noshard_t3 b on a.id = b.id) as c union select * from noshard_t2 order by id, name;
select * from (select * from ( select a.id, t_id, a.name,a.pad from noshard_t1 a join noshard_t2 b on a.id=b.id union (select * from  noshard_t2 where id in (select id from noshard_t3))) as d UNION select * from noshard_t3) as c order by c.id,c.name;