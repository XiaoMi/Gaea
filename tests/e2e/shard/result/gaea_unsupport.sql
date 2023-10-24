sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1`.`test2_0000`)) AS `b` WHERE `a`.`t_id`=`b`.`o_i' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1`.`test3_0000`.`id`,`sbtest1`.`test3_0000`.`pad` FROM (`t' at line 1

sql:select test1.id,test1.t_id,test1.name,test1.pad,test2.id,test2.o_id,test2.name,test2.pad from test1 inner join test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,test2.id,test2.o_id,test2.name,test2.pad from test1 cross join test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,test2.id,test2.o_id,test2.name,test2.pad from test1 join test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 straight_join sbtest1.test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 using(pad) order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad and test1.id>3 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test1.id' in 'field list'

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b;
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1`.`test2_0000`) WHERE `id`>3) AS `b`' at line 1

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC;
mysqlRes:[[3 Development] [4 Finance]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #3 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'Actions.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select id,FirstName,lastname,department,salary from test4 group by Department;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1.test4_0000.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC;
mysqlRes:[[Finance 4] [Development 3]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc;
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test8' doesn't exist

sql:select test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME from test6 CROSS join sbtest1.test8 using(id) order by test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME;
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 xiaojuan] [ORDERKEY_004 CUSTKEY_111 chenqi] [ORDERKEY_005 CUSTKEY_132 marui] [ORDERKEY_007 CUSTKEY_980 yanglu]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1054 (42S22): Unknown column 'test8.O_ORDERKEY' in 'field list'

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_013' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by o_custkey;   
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by 2;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY);
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by counts asc,2 desc;
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY) asc,2 desc limit 10;
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,10;
mysqlRes:[[300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1;
mysqlRes:[[300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #4 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1.test9_0000.MYDATE' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select UPPER((select C_NAME FROM test7 limit 1)) FROM test7 limit 1;
mysqlRes:[[CHENXIAO]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY=(select min(C_CUSTKEY) from test7);
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY) from test7);
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY)+1 from test7);
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003] [ORDERKEY_011 CUSTKEY_012] [ORDERKEY_004 CUSTKEY_111] [ORDERKEY_005 CUSTKEY_132] [ORDERKEY_010 CUSTKEY_333] [ORDERKEY_006 CUSTKEY_420] [ORDERKEY_007 CUSTKEY_980]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select count(*) from sbtest1.test7 as a where a.c_CUSTKEY=(select max(C_CUSTKEY) from test9 where C_CUSTKEY=a.C_CUSTKEY);
mysqlRes:[[1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test9' doesn't exist

sql:select C_CUSTKEY  from sbtest1.test7 as a where (select count(*) from test9 where O_CUSTKEY=a.C_CUSTKEY)=2;
mysqlRes:[[CUSTKEY_003]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test9' doesn't exist

sql:select count(*) from test9 as a where a.id <> all(select id from test7);
mysqlRes:[[3]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select count(*) from test9 as a where 56000< all(select id from test7);
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select count(*) from sbtest1.test7 as a where 2>all(select count(*) from test9 where O_CUSTKEY=C_CUSTKEY);
mysqlRes:[[6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test9' doesn't exist

sql:select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where (a.O_ORDERKEY,O_CUSTKEY)=(select c_ORDERKEY,c_CUSTKEY from test7 where c_name='yanglu');
mysqlRes:[[6 CUSTKEY_420 ORDERKEY_006 231]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test7' doesn't exist

sql:select col2,col1 from t group by col1 with rollup;
mysqlRes:[[5 aa] [10 bb] [15 cc] [20 dd] [30 ee] [30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1.t_0000.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select coalesce(col1,'ALL'),col2,col1 from t group by col1 with rollup;
mysqlRes:[[aa 5 aa] [bb 10 bb] [cc 15 cc] [dd 20 dd] [ee 30 ee] [ALL 30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1.t_0000.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select id,pad from test1 where pad>(select pad from test1 where id=2);
mysqlRes:[[3 4] [4 3] [6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,pad from test1 where pad<(select pad from test1 where id=2);
mysqlRes:[[1 1] [5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,pad from test1 where pad=(select pad from test1 where id=2);
mysqlRes:[[2 2]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,pad from test1 where pad>=(select pad from test1 where id=2);
mysqlRes:[[2 2] [3 4] [4 3] [6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,pad from test1 where pad<=(select pad from test1 where id=2);
mysqlRes:[[1 1] [2 2] [5 1]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,pad from test1 where pad<>(select pad from test1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,pad from test1 where pad !=(select pad from test1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where exists(select * from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where not exists(select * from test1 where pad>1);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where pad=some(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where pad=any(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where pad !=any(select id from test1 where pad=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad ) a,(select sbtest1.test2.id,sbtest1.test2.pad from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1`.`test2_0000`.`id`,`sbtest1`.`test2_0000`.`pad` FROM (`t' at line 1

sql:select id,t_id,name,pad from test1 where pad>(select pad from test1 where pad=2);
mysqlRes:[[3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where 2 >any(select id from test1 where pad>1);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where 2<>some(select id from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select id,t_id,name,pad from test1 where 2>all(select id from test1 where pad<1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1146 (42S02): Table 'sbtest1.test1' doesn't exist

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`)) AS `b` WHERE `a`.`t_id`=`b`.`o_id`' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test3`.`id`,`sbtest1_0`.`test3`.`pad` FROM (`test1` ' at line 1

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b;
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`) WHERE `id`>3) AS `b`' at line 1

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC;
mysqlRes:[[3 Development] [4 Finance]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #3 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'Actions.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select id,FirstName,lastname,department,salary from test4 group by Department;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test4.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC;
mysqlRes:[[Finance 4] [Development 3]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc;
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test8.O_ORDERKEY' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #4 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test9.MYDATE' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select col2,col1 from t group by col1 with rollup;
mysqlRes:[[5 aa] [10 bb] [15 cc] [20 dd] [30 ee] [30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select coalesce(col1,'ALL'),col2,col1 from t group by col1 with rollup;
mysqlRes:[[aa 5 aa] [bb 10 bb] [cc 15 cc] [dd 20 dd] [ee 30 ee] [ALL 30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad ) a,(select sbtest1.test2.id,sbtest1.test2.pad from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test2`.`id`,`sbtest1_0`.`test2`.`pad` FROM (`test1` ' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`)) AS `b` WHERE `a`.`t_id`=`b`.`o_id`' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test3`.`id`,`sbtest1_0`.`test3`.`pad` FROM (`test1` ' at line 1

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b;
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`) WHERE `id`>3) AS `b`' at line 1

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC;
mysqlRes:[[3 Development] [4 Finance]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #3 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'Actions.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select id,FirstName,lastname,department,salary from test4 group by Department;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test4.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC;
mysqlRes:[[Finance 4] [Development 3]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc;
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test8.O_ORDERKEY' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #4 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test9.MYDATE' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select col2,col1 from t group by col1 with rollup;
mysqlRes:[[5 aa] [10 bb] [15 cc] [20 dd] [30 ee] [30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select coalesce(col1,'ALL'),col2,col1 from t group by col1 with rollup;
mysqlRes:[[aa 5 aa] [bb 10 bb] [cc 15 cc] [dd 20 dd] [ee 30 ee] [ALL 30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad ) a,(select sbtest1.test2.id,sbtest1.test2.pad from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test2`.`id`,`sbtest1_0`.`test2`.`pad` FROM (`test1` ' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`)) AS `b` WHERE `a`.`t_id`=`b`.`o_id`' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test3`.`id`,`sbtest1_0`.`test3`.`pad` FROM (`test1` ' at line 1

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b;
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`) WHERE `id`>3) AS `b`' at line 1

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC;
mysqlRes:[[3 Development] [4 Finance]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #3 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'Actions.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select id,FirstName,lastname,department,salary from test4 group by Department;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test4.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC;
mysqlRes:[[Finance 4] [Development 3]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc;
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test8.O_ORDERKEY' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #4 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test9.MYDATE' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select col2,col1 from t group by col1 with rollup;
mysqlRes:[[5 aa] [10 bb] [15 cc] [20 dd] [30 ee] [30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select coalesce(col1,'ALL'),col2,col1 from t group by col1 with rollup;
mysqlRes:[[aa 5 aa] [bb 10 bb] [cc 15 cc] [dd 20 dd] [ee 30 ee] [ALL 30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad ) a,(select sbtest1.test2.id,sbtest1.test2.pad from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test2`.`id`,`sbtest1_0`.`test2`.`pad` FROM (`test1` ' at line 1

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1,sbtest1.test2 where test1.pad=sbtest1.test2.pad;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a,sbtest1.test2 b where a.pad=b.pad;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.name,a.name from test1,(select name from sbtest1.test2) a;
mysqlRes:[[1 test中id为1 order中id为1] [1 test中id为1 test_2] [1 test中id为1 order中id为3] [1 test中id为1 $order$4] [1 test中id为1 order...5] [2 test_2 order中id为1] [2 test_2 test_2] [2 test_2 order中id为3] [2 test_2 $order$4] [2 test_2 order...5] [3 test中id为3 order中id为1] [3 test中id为3 test_2] [3 test中id为3 order中id为3] [3 test中id为3 $order$4] [3 test中id为3 order...5] [4 $test$4 order中id为1] [4 $test$4 test_2] [4 $test$4 order中id为3] [4 $test$4 $order$4] [4 $test$4 order...5] [5 test...5 order中id为1] [5 test...5 test_2] [5 test...5 order中id为3] [5 test...5 $order$4] [5 test...5 order...5] [6 test6 order中id为1] [6 test6 test_2] [6 test6 order中id为3] [6 test6 $order$4] [6 test6 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,test2.id,test2.o_id,test2.name,test2.pad from test1 inner join test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,test2.id,test2.o_id,test2.name,test2.pad from test1 cross join test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,test2.id,test2.o_id,test2.name,test2.pad from test1 join test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.name,a.pad,b.name from test1 a inner join sbtest1.test2 b order by a.id,b.id;
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.name,a.pad,b.name from test1 a cross join sbtest1.test2 b order by a.id,b.id;
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.name,a.pad,b.name from test1 a join sbtest1.test2 b order by a.id,b.id;
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a inner join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a cross join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a inner join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a cross join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b  using(pad) order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 straight_join sbtest1.test2 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.name,a.pad,b.name from test1 a straight_join sbtest1.test2 b order by a.id,b.id;
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a straight_join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a straight_join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a straight_join (select * from sbtest1.test2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 using(pad) order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left outer join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right outer join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join sbtest1.test2 b using(pad) order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left outer join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right outer join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join (select * from sbtest1.test2 where pad>2) b using(pad) order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a right join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left outer join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a right outer join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id;
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left join (select * from sbtest1.test2 where pad>3) b using(pad) order by a.id,b.id;
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural left join sbtest1.test2;
mysqlRes:[[2 2 test_2 2 2 2 test_2 2] [1 1 test中id为1 1 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural right join sbtest1.test2;
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural left outer join sbtest1.test2;
mysqlRes:[[2 2 test_2 2 2 2 test_2 2] [1 1 test中id为1 1 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural right outer join sbtest1.test2;
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left join sbtest1.test2 b order by a.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right join sbtest1.test2 b order by b.id;
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left outer join sbtest1.test2 b order by a.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right outer join sbtest1.test2 b order by b.id;
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id;
mysqlRes:[[NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left outer join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right outer join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id;
mysqlRes:[[NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural left join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id;
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural right join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id;
mysqlRes:[[NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural left outer join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id;
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural right outer join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id;
mysqlRes:[[NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad and test1.id>3 order by test1.id,test2.id;
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select distinct a.pad from test1 a,sbtest1.test2 b where a.pad=b.pad;
mysqlRes:[[1] [2] [4] [3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select distinct b.pad,a.pad from test1 a,(select * from sbtest1.test2 where pad=1) b where a.t_id=b.o_id;
mysqlRes:[[1 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(distinct id),sum(distinct name) from test1 where id=3 or id=7;
mysqlRes:[[1 0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b;
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b,(select * from sbtest1.test2 where id>3) c where a.pad=b.pad and c.pad=b.pad;
mysqlRes:[[1 1 1 1] [5 1 1 5] [3 4 4 3] [1 1 1 1] [5 1 1 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join sbtest1.test2 as b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a inner join sbtest1.test2 b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a cross join sbtest1.test2 b order by a.id,b.id;
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,a.name,a.pad,b.name from test1 a straight_join sbtest1.test2 b on a.pad=b.pad;
mysqlRes:[[1 test中id为1 1 order中id为1] [5 test...5 1 order中id为1] [2 test_2 2 test_2] [4 $test$4 3 order中id为3] [3 test中id为3 4 $order$4] [1 test中id为1 1 order...5] [5 test...5 1 order...5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 use index(ID_index) where  Department ='Finance';
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 force index(ID_index) where ID= 205;
mysqlRes:[[205 Jssme Bdnaa Development 75000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT FirstName, LastName,Department = CASE Department WHEN 'F' THEN 'Financial' WHEN 'D' THEN 'Development'  ELSE 'Other' END FROM test4;
mysqlRes:[[Mazojys Fxoj 0] [Jozzh Lnanyo 0] [Syllauu Dfaafk 0] [Gecrrcc Srlkrt 0] [Jssme Bdnaa 0] [Dnnaao Errllov 0] [Tyoysww Osk 0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where Salary >'40000' and Salary <'70000' and Department ='Finance';
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC;
mysqlRes:[[3 Development] [4 Finance]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT Department, COUNT(ID) FROM test4 GROUP BY Department HAVING COUNT(ID)>3;
mysqlRes:[[Finance 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 order by ID ASC;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 group by Department;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT Department, MIN(Salary) FROM test4  GROUP BY Department HAVING MIN(Salary)>46000;
mysqlRes:[[Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,max(Salary) from test4 group by Department order by Department asc;
mysqlRes:[[Development 75000] [Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,min(Salary) from test4 group by Department order by Department asc;
mysqlRes:[[Development 49000] [Finance 7800]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,sum(Salary) from test4 group by Department order by Department asc;
mysqlRes:[[Development 179000] [Finance 172600]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select ID,Department,Salary from test4 order by 2,3;
mysqlRes:[[207 Development 49000] [206 Development 55000] [205 Development 75000] [201 Finance 7800] [202 Finance 45800] [203 Finance 57000] [204 Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 order by Department ,ID desc;
mysqlRes:[[207 Tyoysww Osk Development 49000] [206 Dnnaao Errllov Development 55000] [205 Jssme Bdnaa Development 75000] [204 Gecrrcc Srlkrt Finance 62000] [203 Syllauu Dfaafk Finance 57000] [202 Jozzh Lnanyo Finance 45800] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department;
mysqlRes:[[Development 3] [Finance 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC;
mysqlRes:[[Finance 4] [Development 3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,count(Salary) as a from test4 group by Department having a=3;
mysqlRes:[[Development 3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,count(Salary) as a from test4 group by Department having a>0;
mysqlRes:[[Development 3] [Finance 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,count(Salary) from test4 group by Department having count(ID) >2;
mysqlRes:[[Development 3] [Finance 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,count(*) as num from test4 group by Department having count(*) >1;
mysqlRes:[[Development 3] [Finance 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,count(*) as num from test4 group by Department having count(*) <=3;
mysqlRes:[[Development 3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department from test4 having Department >3;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department from test4 where Department >0;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,max(salary) from test4 group by Department having max(salary) >10;
mysqlRes:[[Development 75000] [Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select 12 as Department, Department from test4 group by Department;
mysqlRes:[[12 Development] [12 Finance]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 limit 2,10;
mysqlRes:[[203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(salary) from test4 group by department order by department asc;
mysqlRes:[[75000] [62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select min(salary) from test4 group by department order by department asc;
mysqlRes:[[49000] [7800]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(salary) from test4 group by department order by department asc;
mysqlRes:[[179000] [172600]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(salary) from test4 group by department order by department asc;
mysqlRes:[[3] [4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,sum(Salary) a from test4 group by Department having a >=1 order by Department DESC;
mysqlRes:[[Finance 172600] [Development 179000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select Department,count(*) as num from test4 group by Department having count(*) >=4 order by Department ASC;
mysqlRes:[[Finance 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ABS(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ACOS(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development NULL] [Dnnaao Errllov Development NULL] [Tyoysww Osk Development NULL] [Mazojys Fxoj Finance NULL] [Jozzh Lnanyo Finance NULL] [Syllauu Dfaafk Finance NULL] [Gecrrcc Srlkrt Finance NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ASIN(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development NULL] [Dnnaao Errllov Development NULL] [Tyoysww Osk Development NULL] [Mazojys Fxoj Finance NULL] [Jozzh Lnanyo Finance NULL] [Syllauu Dfaafk Finance NULL] [Gecrrcc Srlkrt Finance NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ATAN(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1.570782993461564] [Dnnaao Errllov Development 1.5707781449767169] [Tyoysww Osk Development 1.5707759186316341] [Mazojys Fxoj Finance 1.570668121667394] [Jozzh Lnanyo Finance 1.5707744927337648] [Syllauu Dfaafk Finance 1.5707787829352493] [Gecrrcc Srlkrt Finance 1.5707801977626399]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ATAN(salary,100) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1.569462994251686] [Dnnaao Errllov Development 1.5689781469802169] [Tyoysww Osk Development 1.5687555133016455] [Mazojys Fxoj Finance 1.557976516321996] [Jozzh Lnanyo Finance 1.5686129241509728] [Syllauu Dfaafk Finance 1.5690419426299052] [Gecrrcc Srlkrt Finance 1.5691834249677208]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ATAN2(salary,100) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1.569462994251686] [Dnnaao Errllov Development 1.5689781469802169] [Tyoysww Osk Development 1.5687555133016455] [Mazojys Fxoj Finance 1.557976516321996] [Jozzh Lnanyo Finance 1.5686129241509728] [Syllauu Dfaafk Finance 1.5690419426299052] [Gecrrcc Srlkrt Finance 1.5691834249677208]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,CEIL(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,CEILING(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,COT(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1.055126176601761] [Dnnaao Errllov Development 7.231437579987603] [Tyoysww Osk Development 1.5283848380509242] [Mazojys Fxoj Finance -1.5445941093101545] [Jozzh Lnanyo Finance -0.300046703526689] [Syllauu Dfaafk Finance -0.5642127571169799] [Gecrrcc Srlkrt Finance 1.2648660040338875]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,CRC32(Department) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1179299581] [Dnnaao Errllov Development 1179299581] [Tyoysww Osk Development 1179299581] [Mazojys Fxoj Finance 26596220] [Jozzh Lnanyo Finance 26596220] [Syllauu Dfaafk Finance 26596220] [Gecrrcc Srlkrt Finance 26596220]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,FLOOR(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LN(FirstName),LastName,Department from test4 order by Department;
mysqlRes:[[Jssme NULL Bdnaa Development] [Dnnaao NULL Errllov Development] [Tyoysww NULL Osk Development] [Mazojys NULL Fxoj Finance] [Jozzh NULL Lnanyo Finance] [Syllauu NULL Dfaafk Finance] [Gecrrcc NULL Srlkrt Finance]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,LOG(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 11.225243392518447] [Dnnaao Errllov Development 10.915088464214607] [Tyoysww Osk Development 10.799575577092764] [Mazojys Fxoj Finance 8.961879012677683] [Jozzh Lnanyo Finance 10.732039370102276] [Syllauu Dfaafk Finance 10.950806546816688] [Gecrrcc Srlkrt Finance 11.03488966402723]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,LOG2(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 16.194602975157967] [Dnnaao Errllov Development 15.747143998186745] [Tyoysww Osk Development 15.580494128777296] [Mazojys Fxoj Finance 12.929258408636972] [Jozzh Lnanyo Finance 15.483059977871669] [Syllauu Dfaafk Finance 15.79867429882683] [Gecrrcc Srlkrt Finance 15.919980595048964]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,LOG10(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 4.8750612633917] [Dnnaao Errllov Development 4.740362689494244] [Tyoysww Osk Development 4.690196080028514] [Mazojys Fxoj Finance 3.8920946026904804] [Jozzh Lnanyo Finance 4.660865478003869] [Syllauu Dfaafk Finance 4.7558748556724915] [Gecrrcc Srlkrt Finance 4.7923916894982534]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,MOD(salary,2) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 0] [Dnnaao Errllov Development 0] [Tyoysww Osk Development 0] [Mazojys Fxoj Finance 0] [Jozzh Lnanyo Finance 0] [Syllauu Dfaafk Finance 0] [Gecrrcc Srlkrt Finance 0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,RADIANS(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1308.996938995747] [Dnnaao Errllov Development 959.9310885968813] [Tyoysww Osk Development 855.2113334772215] [Mazojys Fxoj Finance 136.1356816555577] [Jozzh Lnanyo Finance 799.3607974134029] [Syllauu Dfaafk Finance 994.8376736367678] [Gecrrcc Srlkrt Finance 1082.1041362364842]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,ROUND(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,SIGN(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 1] [Dnnaao Errllov Development 1] [Tyoysww Osk Development 1] [Mazojys Fxoj Finance 1] [Jozzh Lnanyo Finance 1] [Syllauu Dfaafk Finance 1] [Gecrrcc Srlkrt Finance 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,SIN(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development -0.6878921013205519] [Dnnaao Errllov Development -0.13698155956977381] [Tyoysww Osk Development -0.5475068647790385] [Mazojys Fxoj Finance 0.5434645393886025] [Jozzh Lnanyo Finance 0.957813972427147] [Syllauu Dfaafk Finance -0.8709373957204783] [Gecrrcc Srlkrt Finance -0.6201872685349616]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,SQRT(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 273.8612787525831] [Dnnaao Errllov Development 234.5207879911715] [Tyoysww Osk Development 221.35943621178654] [Mazojys Fxoj Finance 88.31760866327846] [Jozzh Lnanyo Finance 214.00934559032697] [Syllauu Dfaafk Finance 238.74672772626644] [Gecrrcc Srlkrt Finance 248.99799195977465]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,TAN(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 0.9477539484620641] [Dnnaao Errllov Development 0.13828509047321602] [Tyoysww Osk Development 0.6542854751655689] [Mazojys Fxoj Finance -0.6474192760236663] [Jozzh Lnanyo Finance -3.3328144860323405] [Syllauu Dfaafk Finance -1.7723810519808347] [Gecrrcc Srlkrt Finance 0.7905975785662815]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,TRUNCATE(salary,1) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,TRUNCATE(salary*100,0) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development 7500000] [Dnnaao Errllov Development 5500000] [Tyoysww Osk Development 4900000] [Mazojys Fxoj Finance 780000] [Jozzh Lnanyo Finance 4580000] [Syllauu Dfaafk Finance 5700000] [Gecrrcc Srlkrt Finance 6200000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select FirstName,LastName,Department,SIN(salary) from test4 order by Department;
mysqlRes:[[Jssme Bdnaa Development -0.6878921013205519] [Dnnaao Errllov Development -0.13698155956977381] [Tyoysww Osk Development -0.5475068647790385] [Mazojys Fxoj Finance 0.5434645393886025] [Jozzh Lnanyo Finance 0.957813972427147] [Syllauu Dfaafk Finance -0.8709373957204783] [Gecrrcc Srlkrt Finance -0.6201872685349616]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where Department is Null;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where Department is not Null;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where NOT (ID < 200);
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID <300;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID <1;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID <> 0;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID <> 0 and ID <=1;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID >=205;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID <=205;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID >=205 and ID <=205;
mysqlRes:[[205 Jssme Bdnaa Development 75000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID >1 and ID <=203;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID >=1 and ID=205;
mysqlRes:[[205 Jssme Bdnaa Development 75000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID=(ID>>1)<<1;
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [206 Dnnaao Errllov Development 55000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where ID&1;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [203 Syllauu Dfaafk Finance 57000] [205 Jssme Bdnaa Development 75000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where Salary >'40000' and Salary <'70000' and Department ='Finance' order by Salary ASC;
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where (Salary >'50000' and Salary <'70000') or Department ='Finance' order by Salary ASC;
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [206 Dnnaao Errllov Development 55000] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where FirstName like 'J%';
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [205 Jssme Bdnaa Development 75000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) FROM test4 WHERE Salary is null or FirstName not like '%M%';
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where FirstName in ('Mazojys','Syllauu','Tyoysww');
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [203 Syllauu Dfaafk Finance 57000] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,FirstName,lastname,department,salary from test4 where Salary between 40000 and 50000;
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [207 Tyoysww Osk Development 49000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(salary) from test4 where department = 'Finance';
mysqlRes:[[172600]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(salary) from test4 where department = 'Finance';
mysqlRes:[[62000]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select min(salary) from test4 where department = 'Finance';
mysqlRes:[[7800]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select CURRENT_USER FROM test5;
mysqlRes:[[superroot@%] [superroot@%] [superroot@%] [superroot@%]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(distinct id) from test5;
mysqlRes:[[10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(all id) from test5;
mysqlRes:[[10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id, R_REGIONKEY from test5;
mysqlRes:[[1 1] [2 2] [3 3] [4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,'user is user' from test5;
mysqlRes:[[1 user is user] [2 user is user] [3 user is user] [4 user is user]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id*5,'user is user',10 from test5;
mysqlRes:[[5 user is user 10] [10 user is user 10] [15 user is user 10] [20 user is user 10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select ALL id, R_REGIONKEY, R_NAME, R_COMMENT from test5;
mysqlRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select DISTINCT id, R_REGIONKEY, R_NAME, R_COMMENT from test5;
mysqlRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select DISTINCTROW id, R_REGIONKEY, R_NAME, R_COMMENT from test5;
mysqlRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select ALL HIGH_PRIORITY id,'ID' as detail  from test5;
mysqlRes:[[1 ID] [2 ID] [3 ID] [4 ID]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1;
mysqlRes:[[1 ORDERKEY_001 200000 2014-10-22]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1 and not id=1;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE from test8 where id=1;
mysqlRes:[[1 ORDERKEY_001 CUSTKEY_003 200000 2014-10-22]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) counts from test8 a where MYDATE is null;
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) counts from test8 a where id is null;
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) counts from test8 a where id is not null;
mysqlRes:[[8]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) counts from test8 a where not (id is null);
mysqlRes:[[8]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(O_ORDERKEY) counts from test8 a where O_ORDERKEY like 'ORDERKEY_00%';
mysqlRes:[[6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(O_ORDERKEY) counts from test8 a where O_ORDERKEY not like '%00%';
mysqlRes:[[2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_ORDERKEY<'ORDERKEY_010' and O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_ORDERKEY<'ORDERKEY_010' or O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY in ('CUSTKEY_003','CUSTKEY_420','CUSTKEY_980') group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY not in ('CUSTKEY_003','CUSTKEY_420','CUSTKEY_980');
mysqlRes:[[89212944 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_CUSTKEY,case when sum(O_TOTALPRICE)<100000 then 'D' when sum(O_TOTALPRICE)>100000 and sum(O_TOTALPRICE)<1000000 then 'C' when sum(O_TOTALPRICE)>1000000 and sum(O_TOTALPRICE)<5000000 then 'B' else 'A' end as jibie  from test8 a group by O_CUSTKEY order by jibie, O_CUSTKEY limit 10;
mysqlRes:[[CUSTKEY_333 A] [CUSTKEY_003 C] [CUSTKEY_012 C] [CUSTKEY_111 D] [CUSTKEY_132 D] [CUSTKEY_420 D] [CUSTKEY_980 D]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from test8 where MYDATE between concat(date_format('1992-05-01','%Y-%m'),'-00') and concat(date_format(date_add('1992-05-01',interval 2 month),'%Y-%m'),'-00');
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,sum(O_TOTALPRICE) from test8 where id>1 and id<50 group by  id;
mysqlRes:[[2 100000] [4 500] [5 100] [6 231] [7 12000] [10 88888888] [11 323456]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by 2 asc;
mysqlRes:[[1 1992-05] [1 1992-06] [1 1992-07] [1 1992-08] [1 1992-09] [1 1992-11] [1 2008-01]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 group by id,O_ORDERKEY,MYDATE;
mysqlRes:[[1 ORDERKEY_001 200000 2014-10-22] [2 ORDERKEY_002 100000 1992-05-01] [4 ORDERKEY_004 500 2008-01-05] [5 ORDERKEY_005 100 1992-06-28] [6 ORDERKEY_006 231 1992-11-11] [7 ORDERKEY_007 12000 1992-09-10] [10 ORDERKEY_010 88888888 1992-07-20] [11 ORDERKEY_011 323456 1992-08-22]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc;
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,id from test8 where id>1 and id<50 group by 2 asc;
mysqlRes:[[100000 2] [500 4] [100 5] [231 6] [12000 7] [88888888 10] [323456 11]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc having sums>2000000;
mysqlRes:[[88888888 10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc   having sums>2000000;
mysqlRes:[[88888888 10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having count(O_ORDERKEY)>2;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having min(O_ORDERKEY)>10 and max(O_ORDERKEY)<10000000;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE ) from test8 where id>1 and id<50 having min(O_ORDERKEY)<10000;
mysqlRes:[[89325175]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000);
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY not between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where not (O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300') group by o_custkey;
mysqlRes:[[88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where C_NAME like 'A__A';
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where C_NAME like 'm___i';
mysqlRes:[[marui]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where C_NAME like 'ch___i%%';
mysqlRes:[[chenxiao] [chenqi]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where C_NAME like 'ch___i%%' ESCAPE 'i';
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from sbtest1.test6 where C_NAME not like 'chen%';
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from sbtest1.test6 where not (C_NAME like 'chen%');
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where C_NAME like binary 'chen%';
mysqlRes:[[chenxiao] [chenqi]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where C_NAME like 'chen%';
mysqlRes:[[chenxiao] [chenqi]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where  C_NAME  like concat('%','AM','%');
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_NAME from sbtest1.test6 where  C_NAME  like concat('%','en','%');
mysqlRes:[[chenxiao] [chenqi] [huachen]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_ORDERKEY,C_CUSTKEY,C_NAME from test1,sbtest1.test6 where C_CUSTKEY=c_CUSTKEY and C_ORDERKEY<'ORDERKEY_006';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from (select O_CUSTKEY,count(O_CUSTKEY) as counts from test8 group by O_CUSTKEY) as a where counts<10 group by counts;
mysqlRes:[[6] [1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 INNER join sbtest1.test6 where O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 INNER join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_006';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 CROSS join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_006';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME from test6 CROSS join sbtest1.test8 using(id) order by test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME;
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 xiaojuan] [ORDERKEY_004 CUSTKEY_111 chenqi] [ORDERKEY_005 CUSTKEY_132 marui] [ORDERKEY_007 CUSTKEY_980 yanglu]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from sbtest1.test6 as a STRAIGHT_JOIN test8 b where b.O_CUSTKEY=a.c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select b.O_ORDERKEY,b.O_CUSTKEY,a.C_NAME from sbtest1.test6 a STRAIGHT_JOIN test8 b on b.O_CUSTKEY=a.c_CUSTKEY and b.O_ORDERKEY<'ORDERKEY_007';
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left join test8 b on b.O_CUSTKEY=a.C_CUSTKEY and a.C_CUSTKEY<'CUSTKEY_300';
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left join test8 b using(id) where a.C_CUSTKEY<'CUSTKEY_300';
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_111 wangye NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left OUTER join test8 b using(id);
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_333 chenqi ORDERKEY_004] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_420 yanglu ORDERKEY_007] [CUSTKEY_111 wangye NULL] [CUSTKEY_980 huachen NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a using(id) where a.c_CUSTKEY<'CUSTKEY_300';
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_111 wangye NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a  using(id);
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_333 chenqi ORDERKEY_004] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_420 yanglu ORDERKEY_007] [CUSTKEY_111 wangye NULL] [CUSTKEY_980 huachen NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a on b.O_CUSTKEY=a.C_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300';
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by 2;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY), sums;
mysqlRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [88888888 CUSTKEY_333 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by counts asc,2 desc;
mysqlRes:[[88888888 CUSTKEY_333 1] [500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_013' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by o_custkey;   
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by 2;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY);
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by counts asc,2 desc;
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,2 desc limit 2;
mysqlRes:[[88888888 CUSTKEY_333 1] [500 CUSTKEY_111 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,3;
mysqlRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1;
mysqlRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY) asc,2 desc limit 10;
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,10;
mysqlRes:[[300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1;
mysqlRes:[[300000 CUSTKEY_003 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use key (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_800' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use key (ORDERS_FK1) ignore index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_800' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 ignore key (ORDERS_FK1) force index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 ignore key (ORDERS_FK1) force index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use index for order by (ORDERS_FK1) ignore index for group by (ORDERS_FK1) where O_CUSTKEY between 1 and 50 group by o_custkey;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 as a use index for group by (ORDERS_FK1) ignore index for join (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 a use index for join (ORDERS_FK1) ignore index for join (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test7 as a left join test9 b ignore index for join(ORDERS_FK1) on b.O_CUSTKEY=a.c_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300';
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test7 as a force index for join(primary) left join test9 b ignore index for join(ORDERS_FK1) on b.O_CUSTKEY=a.c_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300';
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select UPPER((select C_NAME FROM test7 limit 1)) FROM test7 limit 1;
mysqlRes:[[CHENXIAO]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY=(select min(C_CUSTKEY) from test7);
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY) from test7);
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY)+1 from test7);
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003] [ORDERKEY_011 CUSTKEY_012] [ORDERKEY_004 CUSTKEY_111] [ORDERKEY_005 CUSTKEY_132] [ORDERKEY_010 CUSTKEY_333] [ORDERKEY_006 CUSTKEY_420] [ORDERKEY_007 CUSTKEY_980]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from sbtest1.test7 as a where a.c_CUSTKEY=(select max(C_CUSTKEY) from test9 where C_CUSTKEY=a.C_CUSTKEY);
mysqlRes:[[1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select C_CUSTKEY  from sbtest1.test7 as a where (select count(*) from test9 where O_CUSTKEY=a.C_CUSTKEY)=2;
mysqlRes:[[CUSTKEY_003]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from test9 as a where a.id <> all(select id from test7);
mysqlRes:[[3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from test9 as a where 56000< all(select id from test7);
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from sbtest1.test7 as a where 2>all(select count(*) from test9 where O_CUSTKEY=C_CUSTKEY);
mysqlRes:[[6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where (a.O_ORDERKEY,O_CUSTKEY)=(select c_ORDERKEY,c_CUSTKEY from test7 where c_name='yanglu');
mysqlRes:[[6 CUSTKEY_420 ORDERKEY_006 231]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1 INNER JOIN t2 where t1.col1 = t2.col1;
mysqlRes:[[1 aa 5 1 aa 5] [6 aa 5 1 aa 5] [2 bb 10 2 bb 10] [7 bb 10 2 bb 10] [3 cc 15 3 cc 15] [8 cc 15 3 cc 15] [4 dd 20 4 dd 20] [9 dd 20 4 dd 20] [5 ee 30 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [6 aa 5 6 aa 5] [2 bb 10 7 bb 10] [7 bb 10 7 bb 10] [3 cc 15 8 cc 15] [8 cc 15 8 cc 15] [4 dd 20 9 dd 20] [9 dd 20 9 dd 20] [5 ee 30 10 ee 30] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:/*comment*/ /*comment*/ select col1 /* this is a comment */ from t;
mysqlRes:[[aa] [aa] [bb] [bb] [cc] [cc] [dd] [dd] [ee] [ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT /*!40001 SQL_NO_CACHE */ * FROM t WHERE 1 limit 0, 2000;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t AS u;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t1, t2;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t1 AS u, t2;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t1, t2 AS u;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t1 AS u, t2 AS v;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t, t1, t2;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT DISTINCT * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT DISTINCTROW * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT ALL * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT id FROM t;
mysqlRes:[[1] [6] [2] [7] [3] [8] [4] [9] [5] [10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * FROM t WHERE 1 = 1;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select 1 from t;
mysqlRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select 1 from t limit 1;
mysqlRes:[[1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select 1 from t where not exists (select 2);
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from sbtest1.t1 where id > 4 and id <=8 order by col1 desc;
mysqlRes:[[5 ee 30] [8 cc 15] [7 bb 10] [6 aa 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select 1 as a from t order by a;
mysqlRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from sbtest1.t1 where id > 1 order by id desc limit 10;
mysqlRes:[[10 ee 30] [9 dd 20] [8 cc 15] [7 bb 10] [6 aa 5] [5 ee 30] [4 dd 20] [3 cc 15] [2 bb 10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from sbtest1.t2 where id < 0;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select 1 as a from t where 1 < any (select 2) order by a;
mysqlRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * from t for update;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * from t lock in share mode;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT * from t1, t2, t3;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 join t2 left join t3 on t2.id = t3.id;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 right join t2 on t1.id = t2.id left join t3 on t3.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 right join t2 using (id) left join t3 using (id);
mysqlRes:[[1 aa 5 aa 5 aa 5] [2 bb 10 bb 10 bb 10] [3 cc 15 cc 15 cc 15] [4 dd 20 dd 20 dd 20] [5 ee 30 ee 30 ee 30] [6 aa 5 aa 5 aa 5] [7 bb 10 bb 10 bb 10] [8 cc 15 cc 15 cc 15] [9 dd 20 dd 20 dd 20] [10 ee 30 ee 30 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 natural join t2;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 natural right join t2;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 natural left outer join t2;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t1 straight_join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select straight_join * from t1 join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select straight_join * from t1 left join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select straight_join * from t1 right join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select straight_join * from t1 straight_join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT *, CAST(col1 AS CHAR CHARACTER SET utf8) FROM t;
mysqlRes:[[1 aa 5 aa] [2 bb 10 bb] [3 cc 15 cc] [4 dd 20 dd] [5 ee 30 ee] [6 aa 5 aa] [7 bb 10 bb] [8 cc 15 cc] [9 dd 20 dd] [10 ee 30 ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT last_insert_id();
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(distinct col1) from t;
mysqlRes:[[ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(distinctrow col1) from t;
mysqlRes:[[ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(distinct all col1) from t;
mysqlRes:[[ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(distinctrow all col1) from t;
mysqlRes:[[ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select max(col2) from t;
mysqlRes:[[30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select min(distinct col1) from t;
mysqlRes:[[aa]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select min(distinctrow col1) from t;
mysqlRes:[[aa]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select min(distinct all col1) from t;
mysqlRes:[[aa]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select min(col2) from t;
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(distinct col1) from t;
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(distinctrow col1) from t;
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(distinct all col1) from t;
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(distinctrow all col1) from t;
mysqlRes:[[0]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select sum(col2) from t;
mysqlRes:[[160]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(col1) from t;
mysqlRes:[[10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from t;
mysqlRes:[[10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(distinct col1, col2) from t;
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(distinctrow col1, col2) from t;
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(all col1) from t;
mysqlRes:[[10]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select group_concat(col2,col1) from t group by col1;
mysqlRes:[[5aa,5aa] [10bb,10bb] [15cc,15cc] [20dd,20dd] [30ee,30ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select group_concat(col2,col1 SEPARATOR ';') from t group by col1;
mysqlRes:[[5aa;5aa] [10bb;10bb] [15cc;15cc] [20dd;20dd] [30ee;30ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select group_concat(distinct col2,col1) from t group by col1;
mysqlRes:[[5aa] [10bb] [15cc] [20dd] [30ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select group_concat(distinctrow col2,col1) from t group by col1;
mysqlRes:[[5aa] [10bb] [15cc] [20dd] [30ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:SELECT col1, GROUP_CONCAT(DISTINCT col2 ORDER BY col2 DESC SEPARATOR ' ') FROM t GROUP BY col1;
mysqlRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t a;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t as a;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t full, t1 `row`, t2 abs;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index (primary);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index (`primary`);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index ();
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index (idx1, idx2);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t ignore key (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t force index for join (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index for order by (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t force index for group by (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t use index for group by (idx1) use index for order by (idx2), t2;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select high_priority * from t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select SQL_CACHE * from t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select col2,col1 from t group by col1 with rollup;
mysqlRes:[[5 aa] [10 bb] [15 cc] [20 dd] [30 ee] [30 NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select coalesce(col1,'ALL'),col2,col1 from t group by col1 with rollup;
mysqlRes:[[aa 5 aa] [bb 10 bb] [cc 15 cc] [dd 20 dd] [ee 30 ee] [ALL 30 NULL]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select col1 from t1 group by col1 order by null;
mysqlRes:[[aa] [bb] [cc] [dd] [ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select col1 from t1 group by col1 order by 1;
mysqlRes:[[aa] [bb] [cc] [dd] [ee]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select all * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select distinct * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,o_id,name,pad from (select * from sbtest1.test2 a group by a.id) a;
mysqlRes:[[1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from (select pad,count(*) from sbtest1.test2 a group by pad) a;
mysqlRes:[[1 2] [2 1] [3 1] [4 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 having pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3 order by id) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 order by id limit 3) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 order by id limit 3) b where a.t_id=b.o_id limit 2;
mysqlRes:[[1 1 1 1] [2 2 2 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from (select sbtest1.test3.pad from test1 left join sbtest1.test3 on test1.pad=sbtest1.test3.pad) a;
mysqlRes:[[1] [1] [2] [3] [4] [6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select pad,count(*) from (select * from test1 where pad>2) a group by pad;
mysqlRes:[[3 1] [4 1] [6 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select pad,count(*) from (select * from test1 where pad>2) a group by pad order by pad;
mysqlRes:[[3 1] [4 1] [6 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select count(*) from (select pad,count(*) a from test1 group by pad) a;
mysqlRes:[[5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad>(select pad from test1 where id=2);
mysqlRes:[[3 4] [4 3] [6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad<(select pad from test1 where id=2);
mysqlRes:[[1 1] [5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad=(select pad from test1 where id=2);
mysqlRes:[[2 2]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad>=(select pad from test1 where id=2);
mysqlRes:[[2 2] [3 4] [4 3] [6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad<=(select pad from test1 where id=2);
mysqlRes:[[1 1] [2 2] [5 1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad<>(select pad from test1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,pad from test1 where pad !=(select pad from test1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where exists(select * from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where not exists(select * from test1 where pad>1);
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where pad=some(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where pad=any(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where pad !=any(select id from test1 where pad=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad ) a,(select sbtest1.test2.id,sbtest1.test2.pad from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where pad>(select pad from test1 where pad=2);
mysqlRes:[[3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select (select name from test1 limit 1);
mysqlRes:[[test中id为1]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where 2 >any(select id from test1 where pad>1);
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where 2<>some(select id from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select id,t_id,name,pad from test1 where 2>all(select id from test1 where pad<1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select co1,co2,co3 from (select id as co1,name as co2,pad as co3 from test1)as tb where co1>1;
mysqlRes:[[2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select * from (select m.id,n.pad from test1 m,sbtest1.test2 n where m.id=n.id AND m.name='test中id为1' and m.pad>7 and m.pad<10)a;
mysqlRes:[]
gaeaError:Error 1105: unknown error: get plan error, db: sbtest1, origin sql: Query, trimmedSql: Query, err: parse sql error, sql: Query, err: line 1 column 5 near "Query" 

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`)) AS `b` WHERE `a`.`t_id`=`b`.`o_id`' at line 1

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test3`.`id`,`sbtest1_0`.`test3`.`pad` FROM (`test1` ' at line 1

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b;
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT * FROM (`sbtest1_0`.`test2`) WHERE `id`>3) AS `b`' at line 1

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC;
mysqlRes:[[3 Development] [4 Finance]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #3 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'Actions.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select id,FirstName,lastname,department,salary from test4 group by Department;
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test4.ID' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC;
mysqlRes:[[Finance 4] [Development 3]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc;
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near 'DESC' at line 1

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000);
mysqlRes:[]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test8.O_ORDERKEY' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey;
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #4 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.test9.MYDATE' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select col2,col1 from t group by col1 with rollup;
mysqlRes:[[5 aa] [10 bb] [15 cc] [20 dd] [30 ee] [30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select coalesce(col1,'ALL'),col2,col1 from t group by col1 with rollup;
mysqlRes:[[aa 5 aa] [bb 10 bb] [cc 15 cc] [dd 20 dd] [ee 30 ee] [ALL 30 NULL]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'sbtest1_0.t.col2' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad ) a,(select sbtest1.test2.id,sbtest1.test2.pad from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaError:Error 1105: unknown error: execute in SelectPlan error: ERROR 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ') JOIN (SELECT `sbtest1_0`.`test2`.`id`,`sbtest1_0`.`test2`.`pad` FROM (`test1` ' at line 1

