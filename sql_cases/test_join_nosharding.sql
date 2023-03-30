# test join nosharding
select * from noshard_t1,noshard_t2 where noshard_t1.pad=noshard_t2.pad
select * from noshard_t1 a,noshard_t2 b where a.pad=b.pad
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3) b where a.t_id=b.o_id
select a.id,b.id,b.pad,a.t_id from (select id,t_id from noshard_t1) a,(select * from noshard_t2) b where a.t_id=b.o_id
select a.id,b.id,b.pad,a.t_id from (select noshard_t1.id,noshard_t1.pad,noshard_t1.t_id from noshard_t1 join noshard_t2 where noshard_t1.pad=noshard_t2.pad ) a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b where a.pad=b.pad
select noshard_t1.id,noshard_t1.name,a.name from noshard_t1,(select name from noshard_t2) a
select * from noshard_t1 inner join noshard_t2 order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 cross join noshard_t2 order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 join noshard_t2 order by noshard_t1.id,noshard_t2.id
select a.id,a.name,a.pad,b.name from noshard_t1 a inner join noshard_t2 b order by a.id,b.id
select a.id,a.name,a.pad,b.name from noshard_t1 a cross join noshard_t2 b order by a.id,b.id
select a.id,a.name,a.pad,b.name from noshard_t1 a join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a inner join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from noshard_t1 a cross join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from noshard_t1 a join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>0) a inner join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>0) a cross join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>0) a join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from noshard_t1 a join (select * from noshard_t2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a join (select * from noshard_t2 where pad>0) b  using(pad) order by a.id,b.id
select * from noshard_t1 straight_join noshard_t2 order by noshard_t1.id,noshard_t2.id
select a.id,a.name,a.pad,b.name from noshard_t1 a straight_join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a straight_join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>0) a straight_join (select * from noshard_t2 where pad>0) b order by a.id,b.id
select * from noshard_t1 a straight_join (select * from noshard_t2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
select * from noshard_t1 left join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 right join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 left outer join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 right outer join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 left join noshard_t2 using(pad) order by noshard_t1.id,noshard_t2.id
select * from noshard_t1 a left join noshard_t2 b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a right join noshard_t2 b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a left outer join noshard_t2 b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a right outer join noshard_t2 b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a left join noshard_t2 b using(pad) order by a.id,b.id
select * from noshard_t1 a left join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a right join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a left outer join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a right outer join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
select * from noshard_t1 a left join (select * from noshard_t2 where pad>2) b using(pad) order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a left join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a right join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a left outer join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a right outer join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a left join (select * from noshard_t2 where pad>3) b using(pad) order by a.id,b.id
select * from noshard_t1 natural left join noshard_t2
select * from noshard_t1 natural right join noshard_t2
select * from noshard_t1 natural left outer join noshard_t2
select * from noshard_t1 natural right outer join noshard_t2
select * from noshard_t1 a natural left join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a natural right join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a natural left outer join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a natural right outer join noshard_t2 b order by a.id,b.id
select * from noshard_t1 a natural left join (select * from noshard_t2 where pad>2) b order by a.id,b.id
select * from noshard_t1 a natural right join (select * from noshard_t2 where pad>2) b order by a.id,b.id
select * from noshard_t1 a natural left outer join (select * from noshard_t2 where pad>2) b order by a.id,b.id
select * from noshard_t1 a natural right outer join (select * from noshard_t2 where pad>2) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a natural left join (select * from noshard_t2 where pad>3) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a natural right join (select * from noshard_t2 where pad>3) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a natural left outer join (select * from noshard_t2 where pad>3) b order by a.id,b.id
select * from (select * from noshard_t1 where pad>1) a natural right outer join (select * from noshard_t2 where pad>3) b order by a.id,b.id
select * from noshard_t1 left join noshard_t2 on noshard_t1.pad=noshard_t2.pad and noshard_t1.id>3 order by noshard_t1.id,noshard_t2.id
    #
#distinct(special_scene)
#
(select pad from noshard_t1) union distinct (select pad from noshard_t2)
    (select * from noshard_t1 where id=2) union distinct (select * from noshard_t2 where id=2)
select distinct a.pad from noshard_t1 a,noshard_t2 b where a.pad=b.pad
select distinct b.pad,a.pad from noshard_t1 a,(select * from noshard_t2 where pad=1) b where a.t_id=b.o_id
select count(distinct pad,name),avg(distinct t_id) from noshard_t1
select count(distinct id),sum(distinct name) from noshard_t1 where id=3 or id=7


# test_join_subquery_mixed
select * from (select a.id, t_id, a.name, a.pad from noshard_t1 a join noshard_t3 b on a.id = b.id union select * from  noshard_t2) as c order by id,name;
select id, t_id, name,pad from (select a.id,t_id, b.name,b.pad from noshard_t1 a join noshard_t3 b on a.id = b.id) as c union select * from noshard_t2 order by id, name;
select * from (select * from ( select a.id, t_id, a.name,a.pad from noshard_t1 a join noshard_t2 b on a.id=b.id union (select * from  noshard_t2 where id in (select id from noshard_t3))) as d UNION select * from noshard_t3) as c order by c.id,c.name;


# test_reference_no_sharding
select * from noshard_t1 a,noshard_t2 b
select * from (select * from noshard_t1 where id<3) a,(select * from noshard_t2 where id>3) b
select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b,(select * from noshard_t2 where id>3) c where a.pad=b.pad and c.pad=b.pad
#
#join table
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


# test_subquery_no_sharding
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

#Second supplement
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
