
select id,FirstName,lastname,department,salary from test4 use index(ID_index) where  Department ='Finance'
select id,FirstName,lastname,department,salary from test4 force index(ID_index) where ID= 205
SELECT FirstName, LastName,Department = CASE Department WHEN 'F' THEN 'Financial' WHEN 'D' THEN 'Development'  ELSE 'Other' END FROM test4
select id,FirstName,lastname,department,salary from test4 where Salary >'40000' and Salary <'70000' and Department ='Finance'
SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC
SELECT id,FirstName,lastname,department,salary FROM test4 ORDER BY FIELD( ID, 203, 207,206)
SELECT Department, COUNT(ID) FROM test4 GROUP BY Department HAVING COUNT(ID)>3
select id,FirstName,lastname,department,salary from test4 order by ID ASC
select id,FirstName,lastname,department,salary from test4 group by Department
SELECT Department, MIN(Salary) FROM test4  GROUP BY Department HAVING MIN(Salary)>46000
select Department,max(Salary) from test4 group by Department order by Department asc
select Department,min(Salary) from test4 group by Department order by Department asc
select Department,avg(Salary) from test4 group by Department order by Department asc
select Department,sum(Salary) from test4 group by Department order by Department asc
select ID,Department,Salary from test4 order by 2,3
select id,FirstName,lastname,department,salary from test4 order by Department ,ID desc
SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department
SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC
select Department,count(Salary) as a from test4 group by Department having a=3
select Department,count(Salary) as a from test4 group by Department having a>0
select Department,count(Salary) from test4 group by Department having count(ID) >2
select Department,count(*) as num from test4 group by Department having count(*) >1
select Department,count(*) as num from test4 group by Department having count(*) <=3
select Department from test4 having Department >3
select Department from test4 where Department >0
select Department,max(salary) from test4 group by Department having max(salary) >10
select 12 as Department, Department from test4 group by Department
select id,FirstName,lastname,department,salary from test4 limit 2,10
select id,FirstName,lastname,department,salary from test4 order by FirstName in ('Syllauu','Dnnaao') desc
select max(salary) from test4 group by department order by department asc
select min(salary) from test4 group by department order by department asc
select avg(salary) from test4 group by department order by department asc
select sum(salary) from test4 group by department order by department asc
select count(salary) from test4 group by department order by department asc
select Department,sum(Salary) a from test4 group by Department having a >=1 order by Department DESC
select Department,count(*) as num from test4 group by Department having count(*) >=4 order by Department ASC
select FirstName,LastName,Department,ABS(salary) from test4 order by Department
select FirstName,LastName,Department,ACOS(salary) from test4 order by Department
select FirstName,LastName,Department,ASIN(salary) from test4 order by Department
select FirstName,LastName,Department,ATAN(salary) from test4 order by Department
select FirstName,LastName,Department,ATAN(salary,100) from test4 order by Department
select FirstName,LastName,Department,ATAN2(salary,100) from test4 order by Department
select FirstName,LastName,Department,CEIL(salary) from test4 order by Department
select FirstName,LastName,Department,CEILING(salary) from test4 order by Department
select FirstName,LastName,Department,COT(salary) from test4 order by Department
select FirstName,LastName,Department,CRC32(Department) from test4 order by Department
select FirstName,LastName,Department,FLOOR(salary) from test4 order by Department
select FirstName,LN(FirstName),LastName,Department from test4 order by Department
select FirstName,LastName,Department,LOG(salary) from test4 order by Department
select FirstName,LastName,Department,LOG2(salary) from test4 order by Department
select FirstName,LastName,Department,LOG10(salary) from test4 order by Department
select FirstName,LastName,Department,MOD(salary,2) from test4 order by Department
select FirstName,LastName,Department,RADIANS(salary) from test4 order by Department
select FirstName,LastName,Department,ROUND(salary) from test4 order by Department
select FirstName,LastName,Department,SIGN(salary) from test4 order by Department
select FirstName,LastName,Department,SIN(salary) from test4 order by Department
select FirstName,LastName,Department,SQRT(salary) from test4 order by Department
select FirstName,LastName,Department,TAN(salary) from test4 order by Department
select FirstName,LastName,Department,TRUNCATE(salary,1) from test4 order by Department
select FirstName,LastName,Department,TRUNCATE(salary*100,0) from test4 order by Department
select FirstName,LastName,Department,SIN(salary) from test4 order by Department
select id,FirstName,lastname,department,salary from test4 where Department is Null
select id,FirstName,lastname,department,salary from test4 where Department is not Null
select id,FirstName,lastname,department,salary from test4 where NOT (ID < 200)
select id,FirstName,lastname,department,salary from test4 where ID <300
select id,FirstName,lastname,department,salary from test4 where ID <1
select id,FirstName,lastname,department,salary from test4 where ID <> 0
select id,FirstName,lastname,department,salary from test4 where ID <> 0 and ID <=1
select id,FirstName,lastname,department,salary from test4 where ID >=205
select id,FirstName,lastname,department,salary from test4 where ID <=205
select id,FirstName,lastname,department,salary from test4 where ID >=205 and ID <=205
select id,FirstName,lastname,department,salary from test4 where ID >1 and ID <=203
select id,FirstName,lastname,department,salary from test4 where ID >=1 and ID=205
select id,FirstName,lastname,department,salary from test4 where ID=(ID>>1)<<1
select id,FirstName,lastname,department,salary from test4 where ID&1
select id,FirstName,lastname,department,salary from test4 where Salary >'40000' and Salary <'70000' and Department ='Finance' order by Salary ASC
select id,FirstName,lastname,department,salary from test4 where (Salary >'50000' and Salary <'70000') or Department ='Finance' order by Salary ASC
select id,FirstName,lastname,department,salary from test4 where FirstName like 'J%'
select count(*) FROM test4 WHERE Salary is null or FirstName not like '%M%'
SELECT id,FirstName,lastname,department,salary FROM test4  WHERE ID IN (SELECT ID FROM test4 WHERE ID >0)
SELECT distinct salary,id,FirstName,lastname,department FROM test4 WHERE ID IN ( SELECT ID FROM test4 WHERE ID >0)
select id,FirstName,lastname,department,salary from test4 where FirstName in ('Mazojys','Syllauu','Tyoysww')
select id,FirstName,lastname,department,salary from test4 where Salary between 40000 and 50000
select sum(salary) from test4 where department = 'Finance'
select max(salary) from test4 where department = 'Finance'
select min(salary) from test4 where department = 'Finance'
select avg(salary) from test4 where department = 'Finance'

select sum(distinct id) from test5
select sum(all id) from test5
select id, R_REGIONKEY from test5
select id,'user is user' from test5
select id*5,'user is user',10 from test5
select ALL id, R_REGIONKEY, R_NAME, R_COMMENT from test5
select DISTINCT id, R_REGIONKEY, R_NAME, R_COMMENT from test5
select DISTINCTROW id, R_REGIONKEY, R_NAME, R_COMMENT from test5
select ALL HIGH_PRIORITY id,'ID' as detail  from test5

select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1
    ## TODO BUG EXISTS here
##select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1 or not id=1
select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1 and not id=1
select id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE from test8 where id=1
select count(*) counts from test8 a where MYDATE is null
select count(*) counts from test8 a where id is null
select count(*) counts from test8 a where id is not null
select count(*) counts from test8 a where not (id is null)
select count(O_ORDERKEY) counts from test8 a where O_ORDERKEY like 'ORDERKEY_00%'
select count(O_ORDERKEY) counts from test8 a where O_ORDERKEY not like '%00%'
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_ORDERKEY<'ORDERKEY_010' and O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_ORDERKEY<'ORDERKEY_010' or O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY in ('CUSTKEY_003','CUSTKEY_420','CUSTKEY_980') group by o_custkey
select sum(O_TOTALPRICE) as sums,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY not in ('CUSTKEY_003','CUSTKEY_420','CUSTKEY_980')
    #select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE from test8 where id=1
select O_CUSTKEY,case when sum(O_TOTALPRICE)<100000 then 'D' when sum(O_TOTALPRICE)>100000 and sum(O_TOTALPRICE)<1000000 then 'C' when sum(O_TOTALPRICE)>1000000 and sum(O_TOTALPRICE)<5000000 then 'B' else 'A' end as jibie  from test8 a group by O_CUSTKEY order by jibie, O_CUSTKEY limit 10

## TODO BUG EXISTS here
#select sum(O_TOTALPRICE) as sums,count(O_ORDERKEY) counts from test8 a where not O_CUSTKEY ='CUSTKEY_003'

    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE from test8 where id=1
select count(*) from test8 where MYDATE between concat(date_format('1992-05-01','%Y-%m'),'-00') and concat(date_format(date_add('1992-05-01',interval 2 month),'%Y-%m'),'-00')
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT count(*) from test8 where mydate between concat(date_format('1992-05-01','%Y-%m'),'-00') and concat(date_format(date_add('1992-05-01',interval 2 month),'%Y-%m'),'-00')
select id,sum(O_TOTALPRICE) from test8 where id>1 and id<50 group by  id
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE id,sum(O_TOTALPRICE) from test8 where id>1 and id<50 group by  id
select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by date_format(id,'%Y-%m')
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by date_format(MYDATE,'%Y-%m')
select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by 2 asc
select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 group by id,O_ORDERKEY,MYDATE
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by 2 asc
select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc
##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc
select sum(O_TOTALPRICE) as sums,id from test8 where id>1 and id<50 group by 2 asc
##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,id from test8 where id>1 and id<50 group by 2 asc
select sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc having sums>2000000
select sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc   having sums>2000000
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc   having sums>2000000
select sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having count(O_ORDERKEY)>2
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having count(O_ORDERKEY)>2
select sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having min(O_ORDERKEY)>10 and max(O_ORDERKEY)<10000000
select sum(O_TOTALPRICE ) from test8 where id>1 and id<50 having min(O_ORDERKEY)<10000
select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000)
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY not between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where not (O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300') group by o_custkey

select C_NAME from sbtest1.test6 where C_NAME like 'A__A'
select C_NAME from sbtest1.test6 where C_NAME like 'm___i'
select C_NAME from sbtest1.test6 where C_NAME like 'ch___i%%'
select C_NAME from sbtest1.test6 where C_NAME like 'ch___i%%' ESCAPE 'i'
select count(*) from sbtest1.test6 where C_NAME not like 'chen%'
select count(*) from sbtest1.test6 where not (C_NAME like 'chen%')
select C_NAME from sbtest1.test6 where C_NAME like binary 'chen%'
select C_NAME from sbtest1.test6 where C_NAME like 'chen%'
select C_NAME from sbtest1.test6 where  C_NAME  like concat('%','AM','%')
select C_NAME from sbtest1.test6 where  C_NAME  like concat('%','en','%')
select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132'
    (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' order by O_ORDERKEY) union (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by O_ORDERKEY)
select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by O_ORDERKEY
select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by 1, O_ORDERKEY desc/*allow_diff_sequence*/
select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by 2,1
    (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_330') union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_420') union (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY<'CUSTKEY_200') order by 2,1/*allow_diff_sequence*/;
(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') union (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') order by O_ORDERKEY;
(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') order by O_ORDERKEY limit 5;
(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' order by O_ORDERKEY limit 5) union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_980' order by O_ORDERKEY limit 5);
(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') order by O_ORDERKEY limit 5
(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') union DISTINCT (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') order by O_ORDERKEY limit 5
select C_ORDERKEY,C_CUSTKEY,C_NAME from test1,sbtest1.test6 where C_CUSTKEY=c_CUSTKEY and C_ORDERKEY<'ORDERKEY_006'
select count(*) from test6 as a where a.id >= any(select id from sbtest1.test6)
select count(*) from test6 as a where a.id >= any(select avg(id) from sbtest1.test6 group by id)
select count(*) from test6 as a where a.id in(select id from sbtest1.test6)
select count(*) from test6 as a where a.id not in(select id from sbtest1.test6)
select count(*) from test6 as a where not (a.id in(select id from sbtest1.test6))
select count(*) from test6 as a where a.id not in(select id from sbtest1.test6 where C_ORDERKEY in(1,2))
select count(*) from test6 as a where not (a.id in(select id from sbtest1.test6 where C_ORDERKEY in(1,2)))
select count(*) from test6 as a where a.id =some(select id from sbtest1.test6)
select count(*) from test6 as a where a.id != any(select id from sbtest1.test6)
select count(*) from (select O_CUSTKEY,count(O_CUSTKEY) as counts from test8 group by O_CUSTKEY) as a where counts<10 group by counts
select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007'
select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 INNER join sbtest1.test6 where O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007'
select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 INNER join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_006'
select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 CROSS join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_006'
select test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME from test6 CROSS join sbtest1.test8 using(id) order by test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME
select O_ORDERKEY,O_CUSTKEY,C_NAME from sbtest1.test6 as a STRAIGHT_JOIN test8 b where b.O_CUSTKEY=a.c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007'
select b.O_ORDERKEY,b.O_CUSTKEY,a.C_NAME from sbtest1.test6 a STRAIGHT_JOIN test8 b on b.O_CUSTKEY=a.c_CUSTKEY and b.O_ORDERKEY<'ORDERKEY_007'
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left join test8 b on b.O_CUSTKEY=a.C_CUSTKEY and a.C_CUSTKEY<'CUSTKEY_300'
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left join test8 b using(id) where a.C_CUSTKEY<'CUSTKEY_300'
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left OUTER join test8 b using(id)
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a using(id) where a.c_CUSTKEY<'CUSTKEY_300'
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a  using(id)
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a on b.O_CUSTKEY=a.C_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300'


select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by o_custkey
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9  where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by 2
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by 2
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY), sums
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY),sums
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by (select C_ORDERKEY from sbtest1.test7 where c_custkey=o_custkey) asc
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by (select C_ORDERKEY from sbtest1.test7 where c_custkey=o_custkey) asc
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by counts asc,2 desc
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by counts asc,2 desc
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_013' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by o_custkey
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_012' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by 2
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by 2
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY)
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY)
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by (select C_ORDERKEY from sbtest1.test7 where c_custkey=o_custkey) asc
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by (select C_ORDERKEY from sbtest1.test7 where c_custkey=o_custkey) asc
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by counts asc,2 desc
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by counts asc,2 desc
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,2 desc limit 2
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,2 desc limit 2
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,3
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,3
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,2 desc limit 10
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,2 desc limit 10
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,10
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,10
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1
    ##select ALL HIGH_PRIORITY STRAIGHT_JOIN SQL_SMALL_RESULT SQL_BIG_RESULT SQL_BUFFER_RESULT SQL_CACHE SQL_CALC_FOUND_ROWS sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use key (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_800' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use key (ORDERS_FK1) ignore index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_800' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 ignore key (ORDERS_FK1) force index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 ignore key (ORDERS_FK1) force index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use index for order by (ORDERS_FK1) ignore index for group by (ORDERS_FK1) where O_CUSTKEY between 1 and 50 group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 as a use index for group by (ORDERS_FK1) ignore index for join (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 a use index for join (ORDERS_FK1) ignore index for join (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test7 as a left join test9 b ignore index for join(ORDERS_FK1) on b.O_CUSTKEY=a.c_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300'
select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test7 as a force index for join(primary) left join test9 b ignore index for join(ORDERS_FK1) on b.O_CUSTKEY=a.c_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300'
select UPPER((select C_NAME FROM sbtest1.test7 limit 1)) FROM sbtest1.test7 limit 1
select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY=(select min(C_CUSTKEY) from sbtest1.test7)
select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY) from sbtest1.test7)
select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY)+1 from sbtest1.test7)
select count(*) from sbtest1.test7 as a where a.c_CUSTKEY=(select max(C_CUSTKEY) from test9 where C_CUSTKEY=a.C_CUSTKEY)
select C_CUSTKEY  from sbtest1.test7 as a where (select count(*) from test9 where O_CUSTKEY=a.C_CUSTKEY)=2
select count(*) from test9 as a where a.id <> all(select id from sbtest1.test7)
select count(*) from test9 as a where 56000< all(select id from sbtest1.test7)
select count(*) from sbtest1.test7 as a where 2>all(select count(*) from test9 where O_CUSTKEY=C_CUSTKEY)
select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where (a.O_ORDERKEY,O_CUSTKEY)=(select c_ORDERKEY,c_CUSTKEY from sbtest1.test7 where c_name='yanglu')
select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where (a.O_ORDERKEY,a.O_CUSTKEY) in (select c_ORDERKEY,c_CUSTKEY from sbtest1.test7) order by id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE
select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where exists(select * from sbtest1.test7 where a.O_CUSTKEY=sbtest1.test7.C_CUSTKEY)
select distinct O_ORDERKEY,O_CUSTKEY from test9 a where exists(select * from sbtest1.test7 where a.O_CUSTKEY=sbtest1.test7.C_CUSTKEY)
select count(*) from test9 a where not exists(select * from sbtest1.test7 where a.O_CUSTKEY=sbtest1.test7.C_CUSTKEY)
