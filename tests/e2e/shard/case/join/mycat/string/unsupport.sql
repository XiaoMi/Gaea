select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
select co1,co2,co3 from (select id as co1,name as co2,pad as co3 from test1)as tb where co1>1;
