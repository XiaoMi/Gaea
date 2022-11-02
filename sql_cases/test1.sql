//Attention, You may need to name the test sql files with 'test' prefix and thus ensure test files will
//be executed after schema file

select * from sbtest1.t1 where id > 1 order by id desc limit 10;
select * from sbtest1.t2 where id < 0;
