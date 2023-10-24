sql:select * from noshard_t1,noshard_t2 where noshard_t1.pad=noshard_t2.pad
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select * from noshard_t1 a,noshard_t2 b where a.pad=b.pad
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3) b where a.t_id=b.o_id
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from noshard_t1) a,(select * from noshard_t2) b where a.t_id=b.o_id
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]

sql:select a.id,b.id,b.pad,a.t_id from (select noshard_t1.id,noshard_t1.pad,noshard_t1.t_id from noshard_t1 join noshard_t2 where noshard_t1.pad=noshard_t2.pad ) a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b where a.pad=b.pad
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]

sql:select noshard_t1.id,noshard_t1.name,a.name from noshard_t1,(select name from noshard_t2) a
mysqlRes:[[1 test中id为1 order中id为1] [1 test中id为1 test_2] [1 test中id为1 order中id为3] [1 test中id为1 $order$4] [1 test中id为1 order...5] [2 test_2 order中id为1] [2 test_2 test_2] [2 test_2 order中id为3] [2 test_2 $order$4] [2 test_2 order...5] [3 test中id为3 order中id为1] [3 test中id为3 test_2] [3 test中id为3 order中id为3] [3 test中id为3 $order$4] [3 test中id为3 order...5] [4 $test$4 order中id为1] [4 $test$4 test_2] [4 $test$4 order中id为3] [4 $test$4 $order$4] [4 $test$4 order...5] [5 test...5 order中id为1] [5 test...5 test_2] [5 test...5 order中id为3] [5 test...5 $order$4] [5 test...5 order...5] [6 test6 order中id为1] [6 test6 test_2] [6 test6 order中id为3] [6 test6 $order$4] [6 test6 order...5]]
gaeaRes:[[1 test中id为1 order中id为1] [1 test中id为1 test_2] [1 test中id为1 order中id为3] [1 test中id为1 $order$4] [1 test中id为1 order...5] [2 test_2 order中id为1] [2 test_2 test_2] [2 test_2 order中id为3] [2 test_2 $order$4] [2 test_2 order...5] [3 test中id为3 order中id为1] [3 test中id为3 test_2] [3 test中id为3 order中id为3] [3 test中id为3 $order$4] [3 test中id为3 order...5] [4 $test$4 order中id为1] [4 $test$4 test_2] [4 $test$4 order中id为3] [4 $test$4 $order$4] [4 $test$4 order...5] [5 test...5 order中id为1] [5 test...5 test_2] [5 test...5 order中id为3] [5 test...5 $order$4] [5 test...5 order...5] [6 test6 order中id为1] [6 test6 test_2] [6 test6 order中id为3] [6 test6 $order$4] [6 test6 order...5]]

sql:select * from noshard_t1 inner join noshard_t2 order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 cross join noshard_t2 order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 join noshard_t2 order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.name,a.pad,b.name from noshard_t1 a inner join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select a.id,a.name,a.pad,b.name from noshard_t1 a cross join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select a.id,a.name,a.pad,b.name from noshard_t1 a join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select * from noshard_t1 a inner join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 a cross join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 a join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from (select * from noshard_t1 where pad>0) a inner join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from (select * from noshard_t1 where pad>0) a cross join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from (select * from noshard_t1 where pad>0) a join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 a join (select * from noshard_t2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]

sql:select * from noshard_t1 a join (select * from noshard_t2 where pad>0) b  using(pad) order by a.id,b.id
mysqlRes:[[1 1 1 test中id为1 1 1 order中id为1] [1 1 1 test中id为1 5 5 order...5] [2 2 2 test_2 2 2 test_2] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 1 1 order中id为1] [1 5 5 test...5 5 5 order...5]]
gaeaRes:[[1 1 1 test中id为1 1 1 order中id为1] [1 1 1 test中id为1 5 5 order...5] [2 2 2 test_2 2 2 test_2] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 1 1 order中id为1] [1 5 5 test...5 5 5 order...5]]

sql:select * from noshard_t1 straight_join noshard_t2 order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.name,a.pad,b.name from noshard_t1 a straight_join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select * from noshard_t1 a straight_join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from (select * from noshard_t1 where pad>0) a straight_join (select * from noshard_t2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 a straight_join (select * from noshard_t2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]

sql:select * from noshard_t1 left join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from noshard_t1 right join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select * from noshard_t1 left outer join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from noshard_t1 right outer join noshard_t2 on noshard_t1.pad=noshard_t2.pad order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select * from noshard_t1 left join noshard_t2 using(pad) order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 1 test中id为1 1 1 order中id为1] [1 1 1 test中id为1 5 5 order...5] [2 2 2 test_2 2 2 test_2] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 1 1 order中id为1] [1 5 5 test...5 5 5 order...5] [6 6 6 test6 NULL NULL NULL]]
gaeaRes:[[1 1 1 test中id为1 1 1 order中id为1] [1 1 1 test中id为1 5 5 order...5] [2 2 2 test_2 2 2 test_2] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 1 1 order中id为1] [1 5 5 test...5 5 5 order...5] [6 6 6 test6 NULL NULL NULL]]

sql:select * from noshard_t1 a left join noshard_t2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from noshard_t1 a right join noshard_t2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select * from noshard_t1 a left outer join noshard_t2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from noshard_t1 a right outer join noshard_t2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select * from noshard_t1 a left join noshard_t2 b using(pad) order by a.id,b.id
mysqlRes:[[1 1 1 test中id为1 1 1 order中id为1] [1 1 1 test中id为1 5 5 order...5] [2 2 2 test_2 2 2 test_2] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 1 1 order中id为1] [1 5 5 test...5 5 5 order...5] [6 6 6 test6 NULL NULL NULL]]
gaeaRes:[[1 1 1 test中id为1 1 1 order中id为1] [1 1 1 test中id为1 5 5 order...5] [2 2 2 test_2 2 2 test_2] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 1 1 order中id为1] [1 5 5 test...5 5 5 order...5] [6 6 6 test6 NULL NULL NULL]]

sql:select * from noshard_t1 a left join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from noshard_t1 a right join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]

sql:select * from noshard_t1 a left outer join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from noshard_t1 a right outer join (select * from noshard_t2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]

sql:select * from noshard_t1 a left join (select * from noshard_t2 where pad>2) b using(pad) order by a.id,b.id
mysqlRes:[[1 1 1 test中id为1 NULL NULL NULL] [2 2 2 test_2 NULL NULL NULL] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 NULL NULL NULL] [6 6 6 test6 NULL NULL NULL]]
gaeaRes:[[1 1 1 test中id为1 NULL NULL NULL] [2 2 2 test_2 NULL NULL NULL] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 3 3 order中id为3] [1 5 5 test...5 NULL NULL NULL] [6 6 6 test6 NULL NULL NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a left join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a right join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4]]

sql:select * from (select * from noshard_t1 where pad>1) a left outer join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a right outer join (select * from noshard_t2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4]]

sql:select * from (select * from noshard_t1 where pad>1) a left join (select * from noshard_t2 where pad>3) b using(pad) order by a.id,b.id
mysqlRes:[[2 2 2 test_2 NULL NULL NULL] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 NULL NULL NULL] [6 6 6 test6 NULL NULL NULL]]
gaeaRes:[[2 2 2 test_2 NULL NULL NULL] [4 3 3 test中id为3 4 4 $order$4] [3 4 4 $test$4 NULL NULL NULL] [6 6 6 test6 NULL NULL NULL]]

sql:select * from noshard_t1 natural left join noshard_t2
mysqlRes:[[2 test_2 2 2 2] [1 test中id为1 1 1 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[2 test_2 2 2 2] [1 test中id为1 1 1 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]

sql:select * from noshard_t1 natural right join noshard_t2
mysqlRes:[[1 order中id为1 1 1 NULL] [2 test_2 2 2 2] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL]]
gaeaRes:[[1 order中id为1 1 1 NULL] [2 test_2 2 2 2] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL]]

sql:select * from noshard_t1 natural left outer join noshard_t2
mysqlRes:[[2 test_2 2 2 2] [1 test中id为1 1 1 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[2 test_2 2 2 2] [1 test中id为1 1 1 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]

sql:select * from noshard_t1 natural right outer join noshard_t2
mysqlRes:[[1 order中id为1 1 1 NULL] [2 test_2 2 2 2] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL]]
gaeaRes:[[1 order中id为1 1 1 NULL] [2 test_2 2 2 2] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL]]

sql:select * from noshard_t1 a natural left join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 2] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 2] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]

sql:select * from noshard_t1 a natural right join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 order中id为1 1 1 NULL] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL] [2 test_2 2 2 2]]
gaeaRes:[[1 order中id为1 1 1 NULL] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL] [2 test_2 2 2 2]]

sql:select * from noshard_t1 a natural left outer join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 2] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 2] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]

sql:select * from noshard_t1 a natural right outer join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 order中id为1 1 1 NULL] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL] [2 test_2 2 2 2]]
gaeaRes:[[1 order中id为1 1 1 NULL] [3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL] [5 order...5 1 5 NULL] [2 test_2 2 2 2]]

sql:select * from noshard_t1 a natural left join (select * from noshard_t2 where pad>2) b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]

sql:select * from noshard_t1 a natural right join (select * from noshard_t2 where pad>2) b order by a.id,b.id
mysqlRes:[[3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL]]
gaeaRes:[[3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL]]

sql:select * from noshard_t1 a natural left outer join (select * from noshard_t2 where pad>2) b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[1 test中id为1 1 1 NULL] [2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [5 test...5 1 5 NULL] [6 test6 6 6 NULL]]

sql:select * from noshard_t1 a natural right outer join (select * from noshard_t2 where pad>2) b order by a.id,b.id
mysqlRes:[[3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL]]
gaeaRes:[[3 order中id为3 3 3 NULL] [4 $order$4 4 4 NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a natural left join (select * from noshard_t2 where pad>3) b order by a.id,b.id
mysqlRes:[[2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [6 test6 6 6 NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a natural right join (select * from noshard_t2 where pad>3) b order by a.id,b.id
mysqlRes:[[4 $order$4 4 4 NULL]]
gaeaRes:[[4 $order$4 4 4 NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a natural left outer join (select * from noshard_t2 where pad>3) b order by a.id,b.id
mysqlRes:[[2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [6 test6 6 6 NULL]]
gaeaRes:[[2 test_2 2 2 NULL] [3 test中id为3 4 3 NULL] [4 $test$4 3 4 NULL] [6 test6 6 6 NULL]]

sql:select * from (select * from noshard_t1 where pad>1) a natural right outer join (select * from noshard_t2 where pad>3) b order by a.id,b.id
mysqlRes:[[4 $order$4 4 4 NULL]]
gaeaRes:[[4 $order$4 4 4 NULL]]

sql:select * from noshard_t1 left join noshard_t2 on noshard_t1.pad=noshard_t2.pad and noshard_t1.id>3 order by noshard_t1.id,noshard_t2.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:(select pad from noshard_t1) union distinct (select pad from noshard_t2)
mysqlRes:[[1] [2] [4] [3] [6]]
gaeaRes:[[1] [2] [4] [3] [6]]

sql:    (select * from noshard_t1 where id=2) union distinct (select * from noshard_t2 where id=2)
mysqlRes:[[2 2 test_2 2]]
gaeaRes:[[2 2 test_2 2]]

sql:select distinct a.pad from noshard_t1 a,noshard_t2 b where a.pad=b.pad
mysqlRes:[[1] [2] [4] [3]]
gaeaRes:[[1] [2] [4] [3]]

sql:select distinct b.pad,a.pad from noshard_t1 a,(select * from noshard_t2 where pad=1) b where a.t_id=b.o_id
mysqlRes:[[1 1]]
gaeaRes:[[1 1]]

sql:select count(distinct pad,name),avg(distinct t_id) from noshard_t1
mysqlRes:[[6 3.5000]]
gaeaRes:[[6 3.5000]]

sql:select count(distinct id),sum(distinct name) from noshard_t1 where id=3 or id=7
mysqlRes:[[1 0]]
gaeaRes:[[1 0]]

sql:select * from (select a.id, t_id, a.name, a.pad from noshard_t1 a join noshard_t3 b on a.id = b.id union select * from  noshard_t2) as c order by id,name;
mysqlRes:[[1 1 order中id为1 1] [1 1 test中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [3 3 test中id为3 4] [4 4 $order$4 4] [4 4 $test$4 3] [5 5 order...5 1] [5 5 test...5 1]]
gaeaRes:[[1 1 order中id为1 1] [1 1 test中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [3 3 test中id为3 4] [4 4 $order$4 4] [4 4 $test$4 3] [5 5 order...5 1] [5 5 test...5 1]]

sql:select id, t_id, name,pad from (select a.id,t_id, b.name,b.pad from noshard_t1 a join noshard_t3 b on a.id = b.id) as c union select * from noshard_t2 order by id, name;
mysqlRes:[[1 1 manager中id为1 1] [1 1 order中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [3 3 order中id为3 3] [4 4 $manager$4 4] [4 4 $order$4 4] [5 5 manager...5 6] [5 5 order...5 1]]
gaeaRes:[[1 1 manager中id为1 1] [1 1 order中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [3 3 order中id为3 3] [4 4 $manager$4 4] [4 4 $order$4 4] [5 5 manager...5 6] [5 5 order...5 1]]

sql:select * from (select * from ( select a.id, t_id, a.name,a.pad from noshard_t1 a join noshard_t2 b on a.id=b.id union (select * from  noshard_t2 where id in (select id from noshard_t3))) as d UNION select * from noshard_t3) as c order by c.id,c.name;
mysqlRes:[[1 1 manager中id为1 1] [1 1 order中id为1 1] [1 1 test中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [3 3 order中id为3 3] [3 3 test中id为3 4] [4 4 $manager$4 4] [4 4 $order$4 4] [4 4 $test$4 3] [5 5 manager...5 6] [5 5 order...5 1] [5 5 test...5 1]]
gaeaRes:[[1 1 manager中id为1 1] [1 1 order中id为1 1] [1 1 test中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [3 3 order中id为3 3] [3 3 test中id为3 4] [4 4 $manager$4 4] [4 4 $order$4 4] [4 4 $test$4 3] [5 5 manager...5 6] [5 5 order...5 1] [5 5 test...5 1]]

sql:select * from noshard_t1 a,noshard_t2 b
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from (select * from noshard_t1 where id<3) a,(select * from noshard_t2 where id>3) b
mysqlRes:[[1 1 test中id为1 1 4 4 $order$4 4] [2 2 test_2 2 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 4 4 $order$4 4] [2 2 test_2 2 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 5 5 order...5 1]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b,(select * from noshard_t2 where id>3) c where a.pad=b.pad and c.pad=b.pad
mysqlRes:[[1 1 1 1] [5 1 1 5] [3 4 4 3] [1 1 1 1] [5 1 1 5]]
gaeaRes:[[1 1 1 1] [5 1 1 5] [3 4 4 3] [1 1 1 1] [5 1 1 5]]

sql:select * from noshard_t1 a join noshard_t2 as b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 a inner join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select * from noshard_t1 a cross join noshard_t2 b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.name,a.pad,b.name from noshard_t1 a straight_join noshard_t2 b on a.pad=b.pad
mysqlRes:[[1 test中id为1 1 order中id为1] [5 test...5 1 order中id为1] [2 test_2 2 test_2] [4 $test$4 3 order中id为3] [3 test中id为3 4 $order$4] [1 test中id为1 1 order...5] [5 test...5 1 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [5 test...5 1 order中id为1] [2 test_2 2 test_2] [4 $test$4 3 order中id为3] [3 test中id为3 4 $order$4] [1 test中id为1 1 order...5] [5 test...5 1 order...5]]

sql:select * from noshard_t1 union all select * from noshard_t3 union all select * from noshard_t2
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]

sql:select * from noshard_t1 union distinct select * from noshard_t3 union distinct select * from noshard_t2
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]

sql:    (select name from noshard_t1 where pad=1 order by id limit 10) union all (select name from noshard_t2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
mysqlRes:[[test中id为1] [test...5] [order中id为1] [order...5]]
gaeaRes:[[test中id为1] [test...5] [order中id为1] [order...5]]

sql:    (select name from noshard_t1 where pad=1 order by id limit 10) union distinct (select name from noshard_t2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
mysqlRes:[[test中id为1] [test...5] [order中id为1] [order...5]]
gaeaRes:[[test中id为1] [test...5] [order中id为1] [order...5]]

sql:    (select * from noshard_t1 where pad=1) union (select * from noshard_t2 where pad=1) order by name limit 10
mysqlRes:[[5 5 order...5 1] [1 1 order中id为1 1] [5 5 test...5 1] [1 1 test中id为1 1]]
gaeaRes:[[5 5 order...5 1] [1 1 order中id为1 1] [5 5 test...5 1] [1 1 test中id为1 1]]

sql:    (select name as sort_a from noshard_t1 where pad=1) union (select name from noshard_t2 where pad=1) order by sort_a limit 10
mysqlRes:[[order...5] [order中id为1] [test...5] [test中id为1]]
gaeaRes:[[order...5] [order中id为1] [test...5] [test中id为1]]

sql:    (select name as sort_a,pad from noshard_t1 where pad=1) union (select name,pad from noshard_t2 where pad=1) order by sort_a,pad limit 10
mysqlRes:[[order...5 1] [order中id为1 1] [test...5 1] [test中id为1 1]]
gaeaRes:[[order...5 1] [order中id为1 1] [test...5 1] [test中id为1 1]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select all * from noshard_t2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select distinct * from noshard_t2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]

sql:select * from (select * from noshard_t2 a group by a.id) a;
mysqlRes:[[1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]

sql:select * from (select pad,count(*) from noshard_t2 a group by pad) a;
mysqlRes:[[1 2] [2 1] [3 1] [4 1]]
gaeaRes:[[1 2] [2 1] [3 1] [4 1]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 having pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3 order by id) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 order by id limit 3) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 order by id limit 3) b where a.t_id=b.o_id limit 2;
mysqlRes:[[1 1 1 1] [2 2 2 2]]
gaeaRes:[[1 1 1 1] [2 2 2 2]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select * from (select noshard_t2.pad from noshard_t1 left join noshard_t2 on noshard_t1.pad=noshard_t2.pad) a;
mysqlRes:[[1] [1] [2] [3] [4] [1] [1] [NULL]]
gaeaRes:[[1] [1] [2] [3] [4] [1] [1] [NULL]]

sql:select * from (select * from noshard_t1 union select * from noshard_t2) a where a.id >3;
mysqlRes:[[4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [4 4 $order$4 4] [5 5 order...5 1]]

sql:select id,pad from noshard_t1 where pad=(select min(id) from noshard_t2);
mysqlRes:[[1 1] [5 1]]
gaeaRes:[[1 1] [5 1]]

sql:select id,pad,name from (select * from noshard_t1 where pad>2) a where id<5;
mysqlRes:[[3 4 test中id为3] [4 3 $test$4]]
gaeaRes:[[3 4 test中id为3] [4 3 $test$4]]

sql:select pad,count(*) from (select * from noshard_t1 where pad>2) a group by pad;
mysqlRes:[[3 1] [4 1] [6 1]]
gaeaRes:[[3 1] [4 1] [6 1]]

sql:select pad,count(*) from (select * from noshard_t1 where pad>2) a group by pad order by pad;
mysqlRes:[[3 1] [4 1] [6 1]]
gaeaRes:[[3 1] [4 1] [6 1]]

sql:select count(*) from (select pad,count(*) a from noshard_t1 group by pad) a;
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select * from noshard_t1 where pad<(select pad from noshard_t2 where id=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]

sql:select * from noshard_t1 having pad<(select pad from noshard_t2 where id=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]

sql:select a.id,b.id,b.pad,a.t_id from noshard_t1 a,(select * from noshard_t2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select id,name,(select count(*) from noshard_t2) count from noshard_t1;
mysqlRes:[[1 test中id为1 5] [2 test_2 5] [3 test中id为3 5] [4 $test$4 5] [5 test...5 5] [6 test6 5]]
gaeaRes:[[1 test中id为1 5] [2 test_2 5] [3 test中id为3 5] [4 $test$4 5] [5 test...5 5] [6 test6 5]]

sql:select * from noshard_t1 where pad like (select pad from noshard_t2 where id=3);
mysqlRes:[[4 4 $test$4 3]]
gaeaRes:[[4 4 $test$4 3]]

sql:select id,pad from noshard_t1 where pad>(select pad from noshard_t1 where id=2);
mysqlRes:[[3 4] [4 3] [6 6]]
gaeaRes:[[3 4] [4 3] [6 6]]

sql:select id,pad from noshard_t1 where pad<(select pad from noshard_t1 where id=2);
mysqlRes:[[1 1] [5 1]]
gaeaRes:[[1 1] [5 1]]

sql:select id,pad from noshard_t1 where pad=(select pad from noshard_t1 where id=2);
mysqlRes:[[2 2]]
gaeaRes:[[2 2]]

sql:select id,pad from noshard_t1 where pad>=(select pad from noshard_t1 where id=2);
mysqlRes:[[2 2] [3 4] [4 3] [6 6]]
gaeaRes:[[2 2] [3 4] [4 3] [6 6]]

sql:select id,pad from noshard_t1 where pad<=(select pad from noshard_t1 where id=2);
mysqlRes:[[1 1] [2 2] [5 1]]
gaeaRes:[[1 1] [2 2] [5 1]]

sql:select id,pad from noshard_t1 where pad<>(select pad from noshard_t1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]

sql:select id,pad from noshard_t1 where pad !=(select pad from noshard_t1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]

sql:select * from noshard_t1 where exists(select * from noshard_t1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where not exists(select * from noshard_t1 where pad>1);
mysqlRes:[]
gaeaRes:[]

sql:select * from noshard_t1 where pad not in(select id from noshard_t1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [5 5 test...5 1]]
gaeaRes:[[1 1 test中id为1 1] [5 5 test...5 1]]

sql:select * from noshard_t1 where pad in(select id from noshard_t1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select * from noshard_t1 where pad=some(select id from noshard_t1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select * from noshard_t1 where pad=any(select id from noshard_t1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select * from noshard_t1 where pad !=any(select id from noshard_t1 where pad=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select a.id,b.id,b.pad,a.t_id from (select noshard_t1.id,noshard_t1.pad,noshard_t1.t_id from noshard_t1 join noshard_t2 where noshard_t1.pad=noshard_t2.pad ) a,(select noshard_t3.id,noshard_t3.pad from noshard_t1 join noshard_t3 where noshard_t1.pad=noshard_t3.pad) b where a.pad=b.pad;
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]

sql:select * from noshard_t1 where pad>(select pad from noshard_t1 where pad=2);
mysqlRes:[[3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select * from noshard_t1,(select * from noshard_t1 where id>3 union select * from noshard_t2 where id<2) a where a.id >3 and noshard_t1.pad=a.pad;
mysqlRes:[[1 1 test中id为1 1 5 5 test...5 1] [4 4 $test$4 3 4 4 $test$4 3] [5 5 test...5 1 5 5 test...5 1] [6 6 test6 6 6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1 5 5 test...5 1] [4 4 $test$4 3 4 4 $test$4 3] [5 5 test...5 1 5 5 test...5 1] [6 6 test6 6 6 6 test6 6]]

sql:select count(*) from (select * from noshard_t1 where pad=(select pad from noshard_t2 where id=1)) a;
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:select (select name from noshard_t1 limit 1)
mysqlRes:[[test中id为1]]
gaeaRes:[[test中id为1]]

sql:select * from noshard_t1 where 'test_2'=(select name from noshard_t2 where id=2)
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where 5=(select count(*) from noshard_t2)
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where 'test_2' like(select name from noshard_t2 where id=2)
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where 2 >any(select id from noshard_t1 where pad>1)
mysqlRes:[]
gaeaRes:[]

sql:select * from noshard_t1 where 2 in(select id from noshard_t1 where pad>1)
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where 2<>some(select id from noshard_t1 where pad>1)
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where 2>all(select id from noshard_t1 where pad<1)
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select * from noshard_t1 where (id,pad)=(select id,pad from noshard_t2 limit 1)
mysqlRes:[[1 1 test中id为1 1]]
gaeaRes:[[1 1 test中id为1 1]]

sql:select * from noshard_t1 where row(id,pad)=(select id,pad from noshard_t2 limit 1)
mysqlRes:[[1 1 test中id为1 1]]
gaeaRes:[[1 1 test中id为1 1]]

sql:select id,name,pad from noshard_t1 where (id,pad)in(select id,pad from noshard_t2)
mysqlRes:[[1 test中id为1 1] [2 test_2 2] [5 test...5 1]]
gaeaRes:[[1 test中id为1 1] [2 test_2 2] [5 test...5 1]]

sql:select id,name,pad from noshard_t1 where (1,1)in(select id,pad from noshard_t2)
mysqlRes:[[1 test中id为1 1] [2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]
gaeaRes:[[1 test中id为1 1] [2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]

sql:SELECT pad FROM noshard_t1 AS x WHERE x.id = (SELECT pad FROM noshard_t2 AS y WHERE x.id = (SELECT pad FROM noshard_t3 WHERE y.id = noshard_t3.id))
mysqlRes:[[1] [2] [4] [3]]
gaeaRes:[[1] [2] [4] [3]]

sql:select co1,co2,co3 from (select id as co1,name as co2,pad as co3 from noshard_t1)as tb where co1>1
mysqlRes:[[2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]
gaeaRes:[[2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]

sql:select avg(sum_column1) from (select sum(id) as sum_column1 from noshard_t1 group by pad) as t1
mysqlRes:[[4.2000]]
gaeaRes:[[4.2000]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1,sbtest1.test2 where test1.pad=sbtest1.test2.pad
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a,sbtest1.test2 b where a.pad=b.pad
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select a.id,b.id,b.pad,a.t_id from (select id,t_id from test1) a,(select * from sbtest1.test2) b where a.t_id=b.o_id
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]

sql:select a.id,b.id,b.pad,a.t_id from (select test1.id,test1.pad,test1.t_id from test1 join sbtest1.test2 where test1.pad=sbtest1.test2.pad ) a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b where a.pad=b.pad
mysqlRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]
gaeaRes:[[1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 1 1 1] [5 1 1 5] [5 1 1 5]]

sql:select test1.id,test1.name,a.name from test1,(select name from sbtest1.test2) a
mysqlRes:[[1 test中id为1 order中id为1] [1 test中id为1 test_2] [1 test中id为1 order中id为3] [1 test中id为1 $order$4] [1 test中id为1 order...5] [2 test_2 order中id为1] [2 test_2 test_2] [2 test_2 order中id为3] [2 test_2 $order$4] [2 test_2 order...5] [3 test中id为3 order中id为1] [3 test中id为3 test_2] [3 test中id为3 order中id为3] [3 test中id为3 $order$4] [3 test中id为3 order...5] [4 $test$4 order中id为1] [4 $test$4 test_2] [4 $test$4 order中id为3] [4 $test$4 $order$4] [4 $test$4 order...5] [5 test...5 order中id为1] [5 test...5 test_2] [5 test...5 order中id为3] [5 test...5 $order$4] [5 test...5 order...5] [6 test6 order中id为1] [6 test6 test_2] [6 test6 order中id为3] [6 test6 $order$4] [6 test6 order...5]]
gaeaRes:[[1 test中id为1 order中id为1] [1 test中id为1 test_2] [1 test中id为1 order中id为3] [1 test中id为1 $order$4] [1 test中id为1 order...5] [2 test_2 order中id为1] [2 test_2 test_2] [2 test_2 order中id为3] [2 test_2 $order$4] [2 test_2 order...5] [3 test中id为3 order中id为1] [3 test中id为3 test_2] [3 test中id为3 order中id为3] [3 test中id为3 $order$4] [3 test中id为3 order...5] [4 $test$4 order中id为1] [4 $test$4 test_2] [4 $test$4 order中id为3] [4 $test$4 $order$4] [4 $test$4 order...5] [5 test...5 order中id为1] [5 test...5 test_2] [5 test...5 order中id为3] [5 test...5 $order$4] [5 test...5 order...5] [6 test6 order中id为1] [6 test6 test_2] [6 test6 order中id为3] [6 test6 $order$4] [6 test6 order...5]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 inner join sbtest1.test2 order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 cross join sbtest1.test2 order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 join sbtest1.test2 order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.name,a.pad,b.name from test1 a inner join sbtest1.test2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select a.id,a.name,a.pad,b.name from test1 a cross join sbtest1.test2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select a.id,a.name,a.pad,b.name from test1 a join sbtest1.test2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a inner join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a cross join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a inner join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a cross join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join (select * from sbtest1.test2 where pad>0) b  using(pad) order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 straight_join sbtest1.test2 order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.name,a.pad,b.name from test1 a straight_join sbtest1.test2 b order by a.id,b.id
mysqlRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [1 test中id为1 1 test_2] [1 test中id为1 1 order中id为3] [1 test中id为1 1 $order$4] [1 test中id为1 1 order...5] [2 test_2 2 order中id为1] [2 test_2 2 test_2] [2 test_2 2 order中id为3] [2 test_2 2 $order$4] [2 test_2 2 order...5] [3 test中id为3 4 order中id为1] [3 test中id为3 4 test_2] [3 test中id为3 4 order中id为3] [3 test中id为3 4 $order$4] [3 test中id为3 4 order...5] [4 $test$4 3 order中id为1] [4 $test$4 3 test_2] [4 $test$4 3 order中id为3] [4 $test$4 3 $order$4] [4 $test$4 3 order...5] [5 test...5 1 order中id为1] [5 test...5 1 test_2] [5 test...5 1 order中id为3] [5 test...5 1 $order$4] [5 test...5 1 order...5] [6 test6 6 order中id为1] [6 test6 6 test_2] [6 test6 6 order中id为3] [6 test6 6 $order$4] [6 test6 6 order...5]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a straight_join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>0) a straight_join (select * from sbtest1.test2 where pad>0) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a straight_join (select * from sbtest1.test2 where pad>0) b on a.id<b.id and a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[1 1 test中id为1 1 5 5 order...5 1] [3 3 test中id为3 4 4 4 $order$4 4]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 right outer join sbtest1.test2 on test1.pad=sbtest1.test2.pad order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 using(pad) order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left outer join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right outer join sbtest1.test2 b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join sbtest1.test2 b using(pad) order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left outer join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a right outer join (select * from sbtest1.test2 where pad>2) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a left join (select * from sbtest1.test2 where pad>2) b using(pad) order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a right join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left outer join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a right outer join (select * from sbtest1.test2 where pad>3) b on a.pad=b.pad order by a.id,b.id
mysqlRes:[[3 3 test中id为3 4 4 4 $order$4 4]]
gaeaRes:[[3 3 test中id为3 4 4 4 $order$4 4]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a left join (select * from sbtest1.test2 where pad>3) b using(pad) order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 4 4 $order$4 4] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural left join sbtest1.test2
mysqlRes:[[2 2 test_2 2 2 2 test_2 2] [1 1 test中id为1 1 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 2 2 test_2 2] [1 1 test中id为1 1 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural right join sbtest1.test2
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural left outer join sbtest1.test2
mysqlRes:[[2 2 test_2 2 2 2 test_2 2] [1 1 test中id为1 1 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 2 2 test_2 2] [1 1 test中id为1 1 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 natural right outer join sbtest1.test2
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left join sbtest1.test2 b order by a.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right join sbtest1.test2 b order by b.id
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left outer join sbtest1.test2 b order by a.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 2 2 test_2 2] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right outer join sbtest1.test2 b order by b.id
mysqlRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]
gaeaRes:[[NULL NULL NULL NULL 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4] [NULL NULL NULL NULL 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
mysqlRes:[[NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaRes:[[NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural left outer join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [5 5 test...5 1 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a natural right outer join (select * from sbtest1.test2 where pad>2) b order by a.id,b.id
mysqlRes:[[NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaRes:[[NULL NULL NULL NULL 3 3 order中id为3 3] [NULL NULL NULL NULL 4 4 $order$4 4]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural left join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural right join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
mysqlRes:[[NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaRes:[[NULL NULL NULL NULL 4 4 $order$4 4]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural left outer join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
mysqlRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 NULL NULL NULL NULL] [6 6 test6 6 NULL NULL NULL NULL]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from (select * from test1 where pad>1) a natural right outer join (select * from sbtest1.test2 where pad>3) b order by a.id,b.id
mysqlRes:[[NULL NULL NULL NULL 4 4 $order$4 4]]
gaeaRes:[[NULL NULL NULL NULL 4 4 $order$4 4]]

sql:select test1.id,test1.t_id,test1.name,test1.pad,sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from test1 left join sbtest1.test2 on test1.pad=sbtest1.test2.pad and test1.id>3 order by test1.id,sbtest1.test2.id
mysqlRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]
gaeaRes:[[1 1 test中id为1 1 NULL NULL NULL NULL] [2 2 test_2 2 NULL NULL NULL NULL] [3 3 test中id为3 4 NULL NULL NULL NULL] [4 4 $test$4 3 3 3 order中id为3 3] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 NULL NULL NULL NULL]]

sql:(select pad from test1) union distinct (select pad from sbtest1.test2)
mysqlRes:[[1] [2] [4] [3] [6]]
gaeaRes:[[1] [2] [4] [3] [6]]

sql:(select test1.id,test1.t_id,test1.name,test1.pad from test1 where id=2) union distinct (select sbtest1.test2.id,sbtest1.test2.o_id,sbtest1.test2.name,sbtest1.test2.pad from sbtest1.test2 where id=2)
mysqlRes:[[2 2 test_2 2]]
gaeaRes:[[2 2 test_2 2]]

sql:select distinct a.pad from test1 a,sbtest1.test2 b where a.pad=b.pad
mysqlRes:[[1] [2] [4] [3]]
gaeaRes:[[1] [2] [4] [3]]

sql:select distinct b.pad,a.pad from test1 a,(select * from sbtest1.test2 where pad=1) b where a.t_id=b.o_id
mysqlRes:[[1 1]]
gaeaRes:[[1 1]]

sql:select count(distinct pad,name),avg(distinct t_id) from test1
mysqlRes:[[6 3.5000]]
gaeaRes:[[6 3.5000]]

sql:select count(distinct id),sum(distinct name) from test1 where id=3 or id=7
mysqlRes:[[1 0]]
gaeaRes:[[1 0]]

sql:select a.id,a.t_id,b.o_id,b.name from (select * from test1 where id<3) a,(select * from sbtest1.test2 where id>3) b
mysqlRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]
gaeaRes:[[1 1 4 $order$4] [1 1 5 order...5] [2 2 4 $order$4] [2 2 5 order...5]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select sbtest1.test3.id,sbtest1.test3.pad from test1 join sbtest1.test3 where test1.pad=sbtest1.test3.pad) b,(select * from sbtest1.test2 where id>3) c where a.pad=b.pad and c.pad=b.pad
mysqlRes:[[1 1 1 1] [5 1 1 5] [3 4 4 3] [1 1 1 1] [5 1 1 5]]
gaeaRes:[[1 1 1 1] [5 1 1 5] [3 4 4 3] [1 1 1 1] [5 1 1 5]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a join sbtest1.test2 as b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a inner join sbtest1.test2 b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad,b.id,b.o_id,b.name,b.pad from test1 a cross join sbtest1.test2 b order by a.id,b.id
mysqlRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1 1 1 order中id为1 1] [1 1 test中id为1 1 2 2 test_2 2] [1 1 test中id为1 1 3 3 order中id为3 3] [1 1 test中id为1 1 4 4 $order$4 4] [1 1 test中id为1 1 5 5 order...5 1] [2 2 test_2 2 1 1 order中id为1 1] [2 2 test_2 2 2 2 test_2 2] [2 2 test_2 2 3 3 order中id为3 3] [2 2 test_2 2 4 4 $order$4 4] [2 2 test_2 2 5 5 order...5 1] [3 3 test中id为3 4 1 1 order中id为1 1] [3 3 test中id为3 4 2 2 test_2 2] [3 3 test中id为3 4 3 3 order中id为3 3] [3 3 test中id为3 4 4 4 $order$4 4] [3 3 test中id为3 4 5 5 order...5 1] [4 4 $test$4 3 1 1 order中id为1 1] [4 4 $test$4 3 2 2 test_2 2] [4 4 $test$4 3 3 3 order中id为3 3] [4 4 $test$4 3 4 4 $order$4 4] [4 4 $test$4 3 5 5 order...5 1] [5 5 test...5 1 1 1 order中id为1 1] [5 5 test...5 1 2 2 test_2 2] [5 5 test...5 1 3 3 order中id为3 3] [5 5 test...5 1 4 4 $order$4 4] [5 5 test...5 1 5 5 order...5 1] [6 6 test6 6 1 1 order中id为1 1] [6 6 test6 6 2 2 test_2 2] [6 6 test6 6 3 3 order中id为3 3] [6 6 test6 6 4 4 $order$4 4] [6 6 test6 6 5 5 order...5 1]]

sql:select a.id,a.name,a.pad,b.name from test1 a straight_join sbtest1.test2 b on a.pad=b.pad
mysqlRes:[[1 test中id为1 1 order中id为1] [5 test...5 1 order中id为1] [2 test_2 2 test_2] [4 $test$4 3 order中id为3] [3 test中id为3 4 $order$4] [1 test中id为1 1 order...5] [5 test...5 1 order...5]]
gaeaRes:[[1 test中id为1 1 order中id为1] [5 test...5 1 order中id为1] [2 test_2 2 test_2] [4 $test$4 3 order中id为3] [3 test中id为3 4 $order$4] [1 test中id为1 1 order...5] [5 test...5 1 order...5]]

sql:select a.id,a.t_id,a.name,a.pad from test1 a union all select b.id,b.m_id,b.name,b.pad from sbtest1.test3 b union all select c.id,c.o_id,c.name,c.pad from sbtest1.test2 c
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [2 2 test_2 2] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]

sql:select a.id,a.t_id,a.name,a.pad from test1 a union distinct select b.id,b.m_id,b.name,b.pad from sbtest1.test3 b union distinct select c.id,c.o_id,c.name,c.pad from sbtest1.test2 c
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [1 1 manager中id为1 1] [3 3 manager中id为3 3] [4 4 $manager$4 4] [5 5 manager...5 6] [1 1 order中id为1 1] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]

sql:(select name from test1 where pad=1 order by id limit 10) union all (select name from sbtest1.test2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
mysqlRes:[[test中id为1] [test...5] [order中id为1] [order...5]]
gaeaRes:[[test中id为1] [test...5] [order中id为1] [order...5]]

sql:(select name from test1 where pad=1 order by id limit 10) union distinct (select name from sbtest1.test2 where pad=1 order by id limit 10)/*allow_diff_sequence*/
mysqlRes:[[test中id为1] [test...5] [order中id为1] [order...5]]
gaeaRes:[[test中id为1] [test...5] [order中id为1] [order...5]]

sql:(select a.id,a.t_id,a.name,a.pad from test1 a where a.pad=1) union (select c.id,c.o_id,c.name,c.pad from sbtest1.test2 c where c.pad=1) order by id limit 10/*allow_diff_sequence*/
mysqlRes:[[1 1 order中id为1 1] [1 1 test中id为1 1] [5 5 order...5 1] [5 5 test...5 1]]
gaeaRes:[[1 1 order中id为1 1] [1 1 test中id为1 1] [5 5 order...5 1] [5 5 test...5 1]]

sql:(select name as sort_a from test1 where pad=1) union (select name from sbtest1.test2 where pad=1) order by sort_a limit 10/*allow_diff_sequence*/
mysqlRes:[[order...5] [order中id为1] [test...5] [test中id为1]]
gaeaRes:[[order...5] [order中id为1] [test...5] [test中id为1]]

sql:(select name as sort_a,pad from test1 where pad=1) union (select name,pad from sbtest1.test2 where pad=1) order by sort_a,pad limit 10/*allow_diff_sequence*/
mysqlRes:[[order...5 1] [order中id为1 1] [test...5 1] [test中id为1 1]]
gaeaRes:[[order...5 1] [order中id为1 1] [test...5 1] [test中id为1 1]]

sql:(select * from test1 where id=2) union (select * from sbtest1.test2 where id=2);
mysqlRes:[[2 2 test_2 2]]
gaeaRes:[[2 2 test_2 2]]

sql:select id,FirstName,lastname,department,salary from test4 use index(ID_index) where  Department ='Finance'
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]

sql:select id,FirstName,lastname,department,salary from test4 force index(ID_index) where ID= 205
mysqlRes:[[205 Jssme Bdnaa Development 75000]]
gaeaRes:[[205 Jssme Bdnaa Development 75000]]

sql:SELECT FirstName, LastName,Department = CASE Department WHEN 'F' THEN 'Financial' WHEN 'D' THEN 'Development'  ELSE 'Other' END FROM test4
mysqlRes:[[Mazojys Fxoj 0] [Jozzh Lnanyo 0] [Syllauu Dfaafk 0] [Gecrrcc Srlkrt 0] [Jssme Bdnaa 0] [Dnnaao Errllov 0] [Tyoysww Osk 0]]
gaeaRes:[[Mazojys Fxoj 0] [Jozzh Lnanyo 0] [Syllauu Dfaafk 0] [Gecrrcc Srlkrt 0] [Jssme Bdnaa 0] [Dnnaao Errllov 0] [Tyoysww Osk 0]]

sql:select id,FirstName,lastname,department,salary from test4 where Salary >'40000' and Salary <'70000' and Department ='Finance'
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]

sql:SELECT count(*), Department  FROM (SELECT * FROM test4 ORDER BY FirstName DESC) AS Actions GROUP BY Department ORDER BY ID DESC
mysqlRes:[[3 Development] [4 Finance]]
gaeaRes:[[3 Development] [4 Finance]]

sql:SELECT id,FirstName,lastname,department,salary FROM test4 ORDER BY FIELD( ID, 203, 207,206)
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [203 Syllauu Dfaafk Finance 57000] [207 Tyoysww Osk Development 49000] [206 Dnnaao Errllov Development 55000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [203 Syllauu Dfaafk Finance 57000] [207 Tyoysww Osk Development 49000] [206 Dnnaao Errllov Development 55000]]

sql:SELECT Department, COUNT(ID) FROM test4 GROUP BY Department HAVING COUNT(ID)>3
mysqlRes:[[Finance 4]]
gaeaRes:[[Finance 4]]

sql:select id,FirstName,lastname,department,salary from test4 order by ID ASC
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 group by Department
mysqlRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]
gaeaRes:[[205 Jssme Bdnaa Development 75000] [201 Mazojys Fxoj Finance 7800]]

sql:SELECT Department, MIN(Salary) FROM test4  GROUP BY Department HAVING MIN(Salary)>46000
mysqlRes:[[Development 49000]]
gaeaRes:[[Development 49000]]

sql:select Department,max(Salary) from test4 group by Department order by Department asc
mysqlRes:[[Development 75000] [Finance 62000]]
gaeaRes:[[Development 75000] [Finance 62000]]

sql:select Department,min(Salary) from test4 group by Department order by Department asc
mysqlRes:[[Development 49000] [Finance 7800]]
gaeaRes:[[Development 49000] [Finance 7800]]

sql:select Department,avg(Salary) from test4 group by Department order by Department asc
mysqlRes:[[Development 59666.6667] [Finance 43150.0000]]
gaeaRes:[[Development 59666.6667] [Finance 43150.0000]]

sql:select Department,sum(Salary) from test4 group by Department order by Department asc
mysqlRes:[[Development 179000] [Finance 172600]]
gaeaRes:[[Development 179000] [Finance 172600]]

sql:select ID,Department,Salary from test4 order by 2,3
mysqlRes:[[207 Development 49000] [206 Development 55000] [205 Development 75000] [201 Finance 7800] [202 Finance 45800] [203 Finance 57000] [204 Finance 62000]]
gaeaRes:[[207 Development 49000] [206 Development 55000] [205 Development 75000] [201 Finance 7800] [202 Finance 45800] [203 Finance 57000] [204 Finance 62000]]

sql:select id,FirstName,lastname,department,salary from test4 order by Department ,ID desc
mysqlRes:[[207 Tyoysww Osk Development 49000] [206 Dnnaao Errllov Development 55000] [205 Jssme Bdnaa Development 75000] [204 Gecrrcc Srlkrt Finance 62000] [203 Syllauu Dfaafk Finance 57000] [202 Jozzh Lnanyo Finance 45800] [201 Mazojys Fxoj Finance 7800]]
gaeaRes:[[207 Tyoysww Osk Development 49000] [206 Dnnaao Errllov Development 55000] [205 Jssme Bdnaa Development 75000] [204 Gecrrcc Srlkrt Finance 62000] [203 Syllauu Dfaafk Finance 57000] [202 Jozzh Lnanyo Finance 45800] [201 Mazojys Fxoj Finance 7800]]

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department
mysqlRes:[[Development 3] [Finance 4]]
gaeaRes:[[Development 3] [Finance 4]]

sql:SELECT Department, COUNT(Salary) FROM test4 GROUP BY Department DESC
mysqlRes:[[Finance 4] [Development 3]]
gaeaRes:[[Finance 4] [Development 3]]

sql:select Department,count(Salary) as a from test4 group by Department having a=3
mysqlRes:[[Development 3]]
gaeaRes:[[Development 3]]

sql:select Department,count(Salary) as a from test4 group by Department having a>0
mysqlRes:[[Development 3] [Finance 4]]
gaeaRes:[[Development 3] [Finance 4]]

sql:select Department,count(Salary) from test4 group by Department having count(ID) >2
mysqlRes:[[Development 3] [Finance 4]]
gaeaRes:[[Development 3] [Finance 4]]

sql:select Department,count(*) as num from test4 group by Department having count(*) >1
mysqlRes:[[Development 3] [Finance 4]]
gaeaRes:[[Development 3] [Finance 4]]

sql:select Department,count(*) as num from test4 group by Department having count(*) <=3
mysqlRes:[[Development 3]]
gaeaRes:[[Development 3]]

sql:select Department from test4 having Department >3
mysqlRes:[]
gaeaRes:[]

sql:select Department from test4 where Department >0
mysqlRes:[]
gaeaRes:[]

sql:select Department,max(salary) from test4 group by Department having max(salary) >10
mysqlRes:[[Development 75000] [Finance 62000]]
gaeaRes:[[Development 75000] [Finance 62000]]

sql:select 12 as Department, Department from test4 group by Department
mysqlRes:[[12 Development] [12 Finance]]
gaeaRes:[[12 Development] [12 Finance]]

sql:select id,FirstName,lastname,department,salary from test4 limit 2,10
mysqlRes:[[203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 order by FirstName in ('Syllauu','Dnnaao') desc
mysqlRes:[[203 Syllauu Dfaafk Finance 57000] [206 Dnnaao Errllov Development 55000] [201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[203 Syllauu Dfaafk Finance 57000] [206 Dnnaao Errllov Development 55000] [201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [207 Tyoysww Osk Development 49000]]

sql:select max(salary) from test4 group by department order by department asc
mysqlRes:[[75000] [62000]]
gaeaRes:[[75000] [62000]]

sql:select min(salary) from test4 group by department order by department asc
mysqlRes:[[49000] [7800]]
gaeaRes:[[49000] [7800]]

sql:select avg(salary) from test4 group by department order by department asc
mysqlRes:[[59666.6667] [43150.0000]]
gaeaRes:[[59666.6667] [43150.0000]]

sql:select sum(salary) from test4 group by department order by department asc
mysqlRes:[[179000] [172600]]
gaeaRes:[[179000] [172600]]

sql:select count(salary) from test4 group by department order by department asc
mysqlRes:[[3] [4]]
gaeaRes:[[3] [4]]

sql:select Department,sum(Salary) a from test4 group by Department having a >=1 order by Department DESC
mysqlRes:[[Finance 172600] [Development 179000]]
gaeaRes:[[Finance 172600] [Development 179000]]

sql:select Department,count(*) as num from test4 group by Department having count(*) >=4 order by Department ASC
mysqlRes:[[Finance 4]]
gaeaRes:[[Finance 4]]

sql:select FirstName,LastName,Department,ABS(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]

sql:select FirstName,LastName,Department,ACOS(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development NULL] [Dnnaao Errllov Development NULL] [Tyoysww Osk Development NULL] [Mazojys Fxoj Finance NULL] [Jozzh Lnanyo Finance NULL] [Syllauu Dfaafk Finance NULL] [Gecrrcc Srlkrt Finance NULL]]
gaeaRes:[[Jssme Bdnaa Development NULL] [Dnnaao Errllov Development NULL] [Tyoysww Osk Development NULL] [Mazojys Fxoj Finance NULL] [Jozzh Lnanyo Finance NULL] [Syllauu Dfaafk Finance NULL] [Gecrrcc Srlkrt Finance NULL]]

sql:select FirstName,LastName,Department,ASIN(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development NULL] [Dnnaao Errllov Development NULL] [Tyoysww Osk Development NULL] [Mazojys Fxoj Finance NULL] [Jozzh Lnanyo Finance NULL] [Syllauu Dfaafk Finance NULL] [Gecrrcc Srlkrt Finance NULL]]
gaeaRes:[[Jssme Bdnaa Development NULL] [Dnnaao Errllov Development NULL] [Tyoysww Osk Development NULL] [Mazojys Fxoj Finance NULL] [Jozzh Lnanyo Finance NULL] [Syllauu Dfaafk Finance NULL] [Gecrrcc Srlkrt Finance NULL]]

sql:select FirstName,LastName,Department,ATAN(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1.570782993461564] [Dnnaao Errllov Development 1.5707781449767169] [Tyoysww Osk Development 1.5707759186316341] [Mazojys Fxoj Finance 1.570668121667394] [Jozzh Lnanyo Finance 1.5707744927337648] [Syllauu Dfaafk Finance 1.5707787829352493] [Gecrrcc Srlkrt Finance 1.5707801977626399]]
gaeaRes:[[Jssme Bdnaa Development 1.570782993461564] [Dnnaao Errllov Development 1.5707781449767169] [Tyoysww Osk Development 1.5707759186316341] [Mazojys Fxoj Finance 1.570668121667394] [Jozzh Lnanyo Finance 1.5707744927337648] [Syllauu Dfaafk Finance 1.5707787829352493] [Gecrrcc Srlkrt Finance 1.5707801977626399]]

sql:select FirstName,LastName,Department,ATAN(salary,100) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1.569462994251686] [Dnnaao Errllov Development 1.5689781469802169] [Tyoysww Osk Development 1.5687555133016455] [Mazojys Fxoj Finance 1.557976516321996] [Jozzh Lnanyo Finance 1.5686129241509728] [Syllauu Dfaafk Finance 1.5690419426299052] [Gecrrcc Srlkrt Finance 1.5691834249677208]]
gaeaRes:[[Jssme Bdnaa Development 1.569462994251686] [Dnnaao Errllov Development 1.5689781469802169] [Tyoysww Osk Development 1.5687555133016455] [Mazojys Fxoj Finance 1.557976516321996] [Jozzh Lnanyo Finance 1.5686129241509728] [Syllauu Dfaafk Finance 1.5690419426299052] [Gecrrcc Srlkrt Finance 1.5691834249677208]]

sql:select FirstName,LastName,Department,ATAN2(salary,100) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1.569462994251686] [Dnnaao Errllov Development 1.5689781469802169] [Tyoysww Osk Development 1.5687555133016455] [Mazojys Fxoj Finance 1.557976516321996] [Jozzh Lnanyo Finance 1.5686129241509728] [Syllauu Dfaafk Finance 1.5690419426299052] [Gecrrcc Srlkrt Finance 1.5691834249677208]]
gaeaRes:[[Jssme Bdnaa Development 1.569462994251686] [Dnnaao Errllov Development 1.5689781469802169] [Tyoysww Osk Development 1.5687555133016455] [Mazojys Fxoj Finance 1.557976516321996] [Jozzh Lnanyo Finance 1.5686129241509728] [Syllauu Dfaafk Finance 1.5690419426299052] [Gecrrcc Srlkrt Finance 1.5691834249677208]]

sql:select FirstName,LastName,Department,CEIL(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]

sql:select FirstName,LastName,Department,CEILING(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]

sql:select FirstName,LastName,Department,COT(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1.055126176601761] [Dnnaao Errllov Development 7.231437579987603] [Tyoysww Osk Development 1.5283848380509242] [Mazojys Fxoj Finance -1.5445941093101545] [Jozzh Lnanyo Finance -0.300046703526689] [Syllauu Dfaafk Finance -0.5642127571169799] [Gecrrcc Srlkrt Finance 1.2648660040338875]]
gaeaRes:[[Jssme Bdnaa Development 1.055126176601761] [Dnnaao Errllov Development 7.231437579987603] [Tyoysww Osk Development 1.5283848380509242] [Mazojys Fxoj Finance -1.5445941093101545] [Jozzh Lnanyo Finance -0.300046703526689] [Syllauu Dfaafk Finance -0.5642127571169799] [Gecrrcc Srlkrt Finance 1.2648660040338875]]

sql:select FirstName,LastName,Department,CRC32(Department) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1179299581] [Dnnaao Errllov Development 1179299581] [Tyoysww Osk Development 1179299581] [Mazojys Fxoj Finance 26596220] [Jozzh Lnanyo Finance 26596220] [Syllauu Dfaafk Finance 26596220] [Gecrrcc Srlkrt Finance 26596220]]
gaeaRes:[[Jssme Bdnaa Development 1179299581] [Dnnaao Errllov Development 1179299581] [Tyoysww Osk Development 1179299581] [Mazojys Fxoj Finance 26596220] [Jozzh Lnanyo Finance 26596220] [Syllauu Dfaafk Finance 26596220] [Gecrrcc Srlkrt Finance 26596220]]

sql:select FirstName,LastName,Department,FLOOR(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]

sql:select FirstName,LN(FirstName),LastName,Department from test4 order by Department
mysqlRes:[[Jssme NULL Bdnaa Development] [Dnnaao NULL Errllov Development] [Tyoysww NULL Osk Development] [Mazojys NULL Fxoj Finance] [Jozzh NULL Lnanyo Finance] [Syllauu NULL Dfaafk Finance] [Gecrrcc NULL Srlkrt Finance]]
gaeaRes:[[Jssme NULL Bdnaa Development] [Dnnaao NULL Errllov Development] [Tyoysww NULL Osk Development] [Mazojys NULL Fxoj Finance] [Jozzh NULL Lnanyo Finance] [Syllauu NULL Dfaafk Finance] [Gecrrcc NULL Srlkrt Finance]]

sql:select FirstName,LastName,Department,LOG(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 11.225243392518447] [Dnnaao Errllov Development 10.915088464214607] [Tyoysww Osk Development 10.799575577092764] [Mazojys Fxoj Finance 8.961879012677683] [Jozzh Lnanyo Finance 10.732039370102276] [Syllauu Dfaafk Finance 10.950806546816688] [Gecrrcc Srlkrt Finance 11.03488966402723]]
gaeaRes:[[Jssme Bdnaa Development 11.225243392518447] [Dnnaao Errllov Development 10.915088464214607] [Tyoysww Osk Development 10.799575577092764] [Mazojys Fxoj Finance 8.961879012677683] [Jozzh Lnanyo Finance 10.732039370102276] [Syllauu Dfaafk Finance 10.950806546816688] [Gecrrcc Srlkrt Finance 11.03488966402723]]

sql:select FirstName,LastName,Department,LOG2(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 16.194602975157967] [Dnnaao Errllov Development 15.747143998186745] [Tyoysww Osk Development 15.580494128777296] [Mazojys Fxoj Finance 12.929258408636972] [Jozzh Lnanyo Finance 15.483059977871669] [Syllauu Dfaafk Finance 15.79867429882683] [Gecrrcc Srlkrt Finance 15.919980595048964]]
gaeaRes:[[Jssme Bdnaa Development 16.194602975157967] [Dnnaao Errllov Development 15.747143998186745] [Tyoysww Osk Development 15.580494128777296] [Mazojys Fxoj Finance 12.929258408636972] [Jozzh Lnanyo Finance 15.483059977871669] [Syllauu Dfaafk Finance 15.79867429882683] [Gecrrcc Srlkrt Finance 15.919980595048964]]

sql:select FirstName,LastName,Department,LOG10(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 4.8750612633917] [Dnnaao Errllov Development 4.740362689494244] [Tyoysww Osk Development 4.690196080028514] [Mazojys Fxoj Finance 3.8920946026904804] [Jozzh Lnanyo Finance 4.660865478003869] [Syllauu Dfaafk Finance 4.7558748556724915] [Gecrrcc Srlkrt Finance 4.7923916894982534]]
gaeaRes:[[Jssme Bdnaa Development 4.8750612633917] [Dnnaao Errllov Development 4.740362689494244] [Tyoysww Osk Development 4.690196080028514] [Mazojys Fxoj Finance 3.8920946026904804] [Jozzh Lnanyo Finance 4.660865478003869] [Syllauu Dfaafk Finance 4.7558748556724915] [Gecrrcc Srlkrt Finance 4.7923916894982534]]

sql:select FirstName,LastName,Department,MOD(salary,2) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 0] [Dnnaao Errllov Development 0] [Tyoysww Osk Development 0] [Mazojys Fxoj Finance 0] [Jozzh Lnanyo Finance 0] [Syllauu Dfaafk Finance 0] [Gecrrcc Srlkrt Finance 0]]
gaeaRes:[[Jssme Bdnaa Development 0] [Dnnaao Errllov Development 0] [Tyoysww Osk Development 0] [Mazojys Fxoj Finance 0] [Jozzh Lnanyo Finance 0] [Syllauu Dfaafk Finance 0] [Gecrrcc Srlkrt Finance 0]]

sql:select FirstName,LastName,Department,RADIANS(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1308.996938995747] [Dnnaao Errllov Development 959.9310885968813] [Tyoysww Osk Development 855.2113334772215] [Mazojys Fxoj Finance 136.1356816555577] [Jozzh Lnanyo Finance 799.3607974134029] [Syllauu Dfaafk Finance 994.8376736367678] [Gecrrcc Srlkrt Finance 1082.1041362364842]]
gaeaRes:[[Jssme Bdnaa Development 1308.996938995747] [Dnnaao Errllov Development 959.9310885968813] [Tyoysww Osk Development 855.2113334772215] [Mazojys Fxoj Finance 136.1356816555577] [Jozzh Lnanyo Finance 799.3607974134029] [Syllauu Dfaafk Finance 994.8376736367678] [Gecrrcc Srlkrt Finance 1082.1041362364842]]

sql:select FirstName,LastName,Department,ROUND(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]

sql:select FirstName,LastName,Department,SIGN(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 1] [Dnnaao Errllov Development 1] [Tyoysww Osk Development 1] [Mazojys Fxoj Finance 1] [Jozzh Lnanyo Finance 1] [Syllauu Dfaafk Finance 1] [Gecrrcc Srlkrt Finance 1]]
gaeaRes:[[Jssme Bdnaa Development 1] [Dnnaao Errllov Development 1] [Tyoysww Osk Development 1] [Mazojys Fxoj Finance 1] [Jozzh Lnanyo Finance 1] [Syllauu Dfaafk Finance 1] [Gecrrcc Srlkrt Finance 1]]

sql:select FirstName,LastName,Department,SIN(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development -0.6878921013205519] [Dnnaao Errllov Development -0.13698155956977381] [Tyoysww Osk Development -0.5475068647790385] [Mazojys Fxoj Finance 0.5434645393886025] [Jozzh Lnanyo Finance 0.957813972427147] [Syllauu Dfaafk Finance -0.8709373957204783] [Gecrrcc Srlkrt Finance -0.6201872685349616]]
gaeaRes:[[Jssme Bdnaa Development -0.6878921013205519] [Dnnaao Errllov Development -0.13698155956977381] [Tyoysww Osk Development -0.5475068647790385] [Mazojys Fxoj Finance 0.5434645393886025] [Jozzh Lnanyo Finance 0.957813972427147] [Syllauu Dfaafk Finance -0.8709373957204783] [Gecrrcc Srlkrt Finance -0.6201872685349616]]

sql:select FirstName,LastName,Department,SQRT(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 273.8612787525831] [Dnnaao Errllov Development 234.5207879911715] [Tyoysww Osk Development 221.35943621178654] [Mazojys Fxoj Finance 88.31760866327846] [Jozzh Lnanyo Finance 214.00934559032697] [Syllauu Dfaafk Finance 238.74672772626644] [Gecrrcc Srlkrt Finance 248.99799195977465]]
gaeaRes:[[Jssme Bdnaa Development 273.8612787525831] [Dnnaao Errllov Development 234.5207879911715] [Tyoysww Osk Development 221.35943621178654] [Mazojys Fxoj Finance 88.31760866327846] [Jozzh Lnanyo Finance 214.00934559032697] [Syllauu Dfaafk Finance 238.74672772626644] [Gecrrcc Srlkrt Finance 248.99799195977465]]

sql:select FirstName,LastName,Department,TAN(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 0.9477539484620641] [Dnnaao Errllov Development 0.13828509047321602] [Tyoysww Osk Development 0.6542854751655689] [Mazojys Fxoj Finance -0.6474192760236663] [Jozzh Lnanyo Finance -3.3328144860323405] [Syllauu Dfaafk Finance -1.7723810519808347] [Gecrrcc Srlkrt Finance 0.7905975785662815]]
gaeaRes:[[Jssme Bdnaa Development 0.9477539484620641] [Dnnaao Errllov Development 0.13828509047321602] [Tyoysww Osk Development 0.6542854751655689] [Mazojys Fxoj Finance -0.6474192760236663] [Jozzh Lnanyo Finance -3.3328144860323405] [Syllauu Dfaafk Finance -1.7723810519808347] [Gecrrcc Srlkrt Finance 0.7905975785662815]]

sql:select FirstName,LastName,Department,TRUNCATE(salary,1) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[Jssme Bdnaa Development 75000] [Dnnaao Errllov Development 55000] [Tyoysww Osk Development 49000] [Mazojys Fxoj Finance 7800] [Jozzh Lnanyo Finance 45800] [Syllauu Dfaafk Finance 57000] [Gecrrcc Srlkrt Finance 62000]]

sql:select FirstName,LastName,Department,TRUNCATE(salary*100,0) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development 7500000] [Dnnaao Errllov Development 5500000] [Tyoysww Osk Development 4900000] [Mazojys Fxoj Finance 780000] [Jozzh Lnanyo Finance 4580000] [Syllauu Dfaafk Finance 5700000] [Gecrrcc Srlkrt Finance 6200000]]
gaeaRes:[[Jssme Bdnaa Development 7500000] [Dnnaao Errllov Development 5500000] [Tyoysww Osk Development 4900000] [Mazojys Fxoj Finance 780000] [Jozzh Lnanyo Finance 4580000] [Syllauu Dfaafk Finance 5700000] [Gecrrcc Srlkrt Finance 6200000]]

sql:select FirstName,LastName,Department,SIN(salary) from test4 order by Department
mysqlRes:[[Jssme Bdnaa Development -0.6878921013205519] [Dnnaao Errllov Development -0.13698155956977381] [Tyoysww Osk Development -0.5475068647790385] [Mazojys Fxoj Finance 0.5434645393886025] [Jozzh Lnanyo Finance 0.957813972427147] [Syllauu Dfaafk Finance -0.8709373957204783] [Gecrrcc Srlkrt Finance -0.6201872685349616]]
gaeaRes:[[Jssme Bdnaa Development -0.6878921013205519] [Dnnaao Errllov Development -0.13698155956977381] [Tyoysww Osk Development -0.5475068647790385] [Mazojys Fxoj Finance 0.5434645393886025] [Jozzh Lnanyo Finance 0.957813972427147] [Syllauu Dfaafk Finance -0.8709373957204783] [Gecrrcc Srlkrt Finance -0.6201872685349616]]

sql:select id,FirstName,lastname,department,salary from test4 where Department is Null
mysqlRes:[]
gaeaRes:[]

sql:select id,FirstName,lastname,department,salary from test4 where Department is not Null
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where NOT (ID < 200)
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID <300
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID <1
mysqlRes:[]
gaeaRes:[]

sql:select id,FirstName,lastname,department,salary from test4 where ID <> 0
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID <> 0 and ID <=1
mysqlRes:[]
gaeaRes:[]

sql:select id,FirstName,lastname,department,salary from test4 where ID >=205
mysqlRes:[[205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID <=205
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID >=205 and ID <=205
mysqlRes:[[205 Jssme Bdnaa Development 75000]]
gaeaRes:[[205 Jssme Bdnaa Development 75000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID >1 and ID <=203
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID >=1 and ID=205
mysqlRes:[[205 Jssme Bdnaa Development 75000]]
gaeaRes:[[205 Jssme Bdnaa Development 75000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID=(ID>>1)<<1
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [206 Dnnaao Errllov Development 55000]]
gaeaRes:[[202 Jozzh Lnanyo Finance 45800] [204 Gecrrcc Srlkrt Finance 62000] [206 Dnnaao Errllov Development 55000]]

sql:select id,FirstName,lastname,department,salary from test4 where ID&1
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [203 Syllauu Dfaafk Finance 57000] [205 Jssme Bdnaa Development 75000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [203 Syllauu Dfaafk Finance 57000] [205 Jssme Bdnaa Development 75000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where Salary >'40000' and Salary <'70000' and Department ='Finance' order by Salary ASC
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]

sql:select id,FirstName,lastname,department,salary from test4 where (Salary >'50000' and Salary <'70000') or Department ='Finance' order by Salary ASC
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [206 Dnnaao Errllov Development 55000] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [206 Dnnaao Errllov Development 55000] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000]]

sql:select id,FirstName,lastname,department,salary from test4 where FirstName like 'J%'
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [205 Jssme Bdnaa Development 75000]]
gaeaRes:[[202 Jozzh Lnanyo Finance 45800] [205 Jssme Bdnaa Development 75000]]

sql:select count(*) FROM test4 WHERE Salary is null or FirstName not like '%M%'
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:SELECT id,FirstName,lastname,department,salary FROM test4  WHERE ID IN (SELECT ID FROM test4 WHERE ID >0)
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [202 Jozzh Lnanyo Finance 45800] [203 Syllauu Dfaafk Finance 57000] [204 Gecrrcc Srlkrt Finance 62000] [205 Jssme Bdnaa Development 75000] [206 Dnnaao Errllov Development 55000] [207 Tyoysww Osk Development 49000]]

sql:SELECT distinct salary,id,FirstName,lastname,department FROM test4 WHERE ID IN ( SELECT ID FROM test4 WHERE ID >0)
mysqlRes:[[7800 201 Mazojys Fxoj Finance] [45800 202 Jozzh Lnanyo Finance] [57000 203 Syllauu Dfaafk Finance] [62000 204 Gecrrcc Srlkrt Finance] [75000 205 Jssme Bdnaa Development] [55000 206 Dnnaao Errllov Development] [49000 207 Tyoysww Osk Development]]
gaeaRes:[[7800 201 Mazojys Fxoj Finance] [45800 202 Jozzh Lnanyo Finance] [57000 203 Syllauu Dfaafk Finance] [62000 204 Gecrrcc Srlkrt Finance] [75000 205 Jssme Bdnaa Development] [55000 206 Dnnaao Errllov Development] [49000 207 Tyoysww Osk Development]]

sql:select id,FirstName,lastname,department,salary from test4 where FirstName in ('Mazojys','Syllauu','Tyoysww')
mysqlRes:[[201 Mazojys Fxoj Finance 7800] [203 Syllauu Dfaafk Finance 57000] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[201 Mazojys Fxoj Finance 7800] [203 Syllauu Dfaafk Finance 57000] [207 Tyoysww Osk Development 49000]]

sql:select id,FirstName,lastname,department,salary from test4 where Salary between 40000 and 50000
mysqlRes:[[202 Jozzh Lnanyo Finance 45800] [207 Tyoysww Osk Development 49000]]
gaeaRes:[[202 Jozzh Lnanyo Finance 45800] [207 Tyoysww Osk Development 49000]]

sql:select sum(salary) from test4 where department = 'Finance'
mysqlRes:[[172600]]
gaeaRes:[[172600]]

sql:select max(salary) from test4 where department = 'Finance'
mysqlRes:[[62000]]
gaeaRes:[[62000]]

sql:select min(salary) from test4 where department = 'Finance'
mysqlRes:[[7800]]
gaeaRes:[[7800]]

sql:select avg(salary) from test4 where department = 'Finance'
mysqlRes:[[43150.0000]]
gaeaRes:[[43150.0000]]

sql:select CURRENT_USER FROM test5
mysqlRes:[[superroot@%] [superroot@%] [superroot@%] [superroot@%]]
gaeaRes:[[superroot@%] [superroot@%] [superroot@%] [superroot@%]]

sql:select sum(distinct id) from test5
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select sum(all id) from test5
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select id, R_REGIONKEY from test5
mysqlRes:[[1 1] [2 2] [3 3] [4 4]]
gaeaRes:[[1 1] [2 2] [3 3] [4 4]]

sql:select id,'user is user' from test5
mysqlRes:[[1 user is user] [2 user is user] [3 user is user] [4 user is user]]
gaeaRes:[[1 user is user] [2 user is user] [3 user is user] [4 user is user]]

sql:select id*5,'user is user',10 from test5
mysqlRes:[[5 user is user 10] [10 user is user 10] [15 user is user 10] [20 user is user 10]]
gaeaRes:[[5 user is user 10] [10 user is user 10] [15 user is user 10] [20 user is user 10]]

sql:select ALL id, R_REGIONKEY, R_NAME, R_COMMENT from test5
mysqlRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]
gaeaRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]

sql:select DISTINCT id, R_REGIONKEY, R_NAME, R_COMMENT from test5
mysqlRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]
gaeaRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]

sql:select DISTINCTROW id, R_REGIONKEY, R_NAME, R_COMMENT from test5
mysqlRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]
gaeaRes:[[1 1 Eastern test001] [2 2 Western test002] [3 3 Northern test003] [4 4 Southern test004]]

sql:select ALL HIGH_PRIORITY id,'ID' as detail  from test5
mysqlRes:[[1 ID] [2 ID] [3 ID] [4 ID]]
gaeaRes:[[1 ID] [2 ID] [3 ID] [4 ID]]

sql:select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1
mysqlRes:[[1 ORDERKEY_001 200000 2014-10-22]]
gaeaRes:[[1 ORDERKEY_001 200000 2014-10-22]]

sql:select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 where id=1 and not id=1
mysqlRes:[]
gaeaRes:[]

sql:select id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE from test8 where id=1
mysqlRes:[[1 ORDERKEY_001 CUSTKEY_003 200000 2014-10-22]]
gaeaRes:[[1 ORDERKEY_001 CUSTKEY_003 200000 2014-10-22]]

sql:select count(*) counts from test8 a where MYDATE is null
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select count(*) counts from test8 a where id is null
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select count(*) counts from test8 a where id is not null
mysqlRes:[[8]]
gaeaRes:[[8]]

sql:select count(*) counts from test8 a where not (id is null)
mysqlRes:[[8]]
gaeaRes:[[8]]

sql:select count(O_ORDERKEY) counts from test8 a where O_ORDERKEY like 'ORDERKEY_00%'
mysqlRes:[[6]]
gaeaRes:[[6]]

sql:select count(O_ORDERKEY) counts from test8 a where O_ORDERKEY not like '%00%'
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_ORDERKEY<'ORDERKEY_010' and O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_ORDERKEY<'ORDERKEY_010' or O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY in ('CUSTKEY_003','CUSTKEY_420','CUSTKEY_980') group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]

sql:select sum(O_TOTALPRICE) as sums,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY not in ('CUSTKEY_003','CUSTKEY_420','CUSTKEY_980')
mysqlRes:[[89212944 4]]
gaeaRes:[[89212944 4]]

sql:select O_CUSTKEY,case when sum(O_TOTALPRICE)<100000 then 'D' when sum(O_TOTALPRICE)>100000 and sum(O_TOTALPRICE)<1000000 then 'C' when sum(O_TOTALPRICE)>1000000 and sum(O_TOTALPRICE)<5000000 then 'B' else 'A' end as jibie  from test8 a group by O_CUSTKEY order by jibie, O_CUSTKEY limit 10
mysqlRes:[[CUSTKEY_333 A] [CUSTKEY_003 C] [CUSTKEY_012 C] [CUSTKEY_111 D] [CUSTKEY_132 D] [CUSTKEY_420 D] [CUSTKEY_980 D]]
gaeaRes:[[CUSTKEY_333 A] [CUSTKEY_003 C] [CUSTKEY_012 C] [CUSTKEY_111 D] [CUSTKEY_132 D] [CUSTKEY_420 D] [CUSTKEY_980 D]]

sql:select count(*) from test8 where MYDATE between concat(date_format('1992-05-01','%Y-%m'),'-00') and concat(date_format(date_add('1992-05-01',interval 2 month),'%Y-%m'),'-00')
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select id,sum(O_TOTALPRICE) from test8 where id>1 and id<50 group by  id
mysqlRes:[[2 100000] [4 500] [5 100] [6 231] [7 12000] [10 88888888] [11 323456]]
gaeaRes:[[2 100000] [4 500] [5 100] [6 231] [7 12000] [10 88888888] [11 323456]]

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by date_format(id,'%Y-%m')
mysqlRes:[[7 1992-05]]
gaeaRes:[[7 1992-05]]

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth from test8 where id>1 and id<50 group by 2 asc
mysqlRes:[[1 1992-05] [1 1992-06] [1 1992-07] [1 1992-08] [1 1992-09] [1 1992-11] [1 2008-01]]
gaeaRes:[[1 1992-05] [1 1992-06] [1 1992-07] [1 1992-08] [1 1992-09] [1 1992-11] [1 2008-01]]

sql:select id, O_ORDERKEY, O_TOTALPRICE,MYDATE from test8 group by id,O_ORDERKEY,MYDATE
mysqlRes:[[1 ORDERKEY_001 200000 2014-10-22] [2 ORDERKEY_002 100000 1992-05-01] [4 ORDERKEY_004 500 2008-01-05] [5 ORDERKEY_005 100 1992-06-28] [6 ORDERKEY_006 231 1992-11-11] [7 ORDERKEY_007 12000 1992-09-10] [10 ORDERKEY_010 88888888 1992-07-20] [11 ORDERKEY_011 323456 1992-08-22]]
gaeaRes:[[1 ORDERKEY_001 200000 2014-10-22] [2 ORDERKEY_002 100000 1992-05-01] [4 ORDERKEY_004 500 2008-01-05] [5 ORDERKEY_005 100 1992-06-28] [6 ORDERKEY_006 231 1992-11-11] [7 ORDERKEY_007 12000 1992-09-10] [10 ORDERKEY_010 88888888 1992-07-20] [11 ORDERKEY_011 323456 1992-08-22]]

sql:select count(id) as counts,date_format(MYDATE,'%Y-%m') as mouth,id from test8 where id>1 and id<50 group by 2 asc ,id desc
mysqlRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]
gaeaRes:[[1 1992-05 2] [1 1992-06 5] [1 1992-07 10] [1 1992-08 11] [1 1992-09 7] [1 1992-11 6] [1 2008-01 4]]

sql:select sum(O_TOTALPRICE) as sums,id from test8 where id>1 and id<50 group by 2 asc
mysqlRes:[[100000 2] [500 4] [100 5] [231 6] [12000 7] [88888888 10] [323456 11]]
gaeaRes:[[100000 2] [500 4] [100 5] [231 6] [12000 7] [88888888 10] [323456 11]]

sql:select sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc having sums>2000000
mysqlRes:[[88888888 10]]
gaeaRes:[[88888888 10]]

sql:select sum(O_TOTALPRICE ) as sums,id from test8 where id>1 and id<50 group by 2 asc   having sums>2000000
mysqlRes:[[88888888 10]]
gaeaRes:[[88888888 10]]

sql:select sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having count(O_ORDERKEY)>2
mysqlRes:[]
gaeaRes:[]

sql:select sum(O_TOTALPRICE ) as sums,id,count(O_ORDERKEY) counts from test8 where id>1 and id<50 group by 2 asc   having min(O_ORDERKEY)>10 and max(O_ORDERKEY)<10000000
mysqlRes:[]
gaeaRes:[]

sql:select sum(O_TOTALPRICE ) from test8 where id>1 and id<50 having min(O_ORDERKEY)<10000
mysqlRes:[[89325175]]
gaeaRes:[[89325175]]

sql:select id,O_ORDERKEY,O_TOTALPRICE from test8 where id>36900 and id<36902 group by O_ORDERKEY  having O_ORDERKEY in (select O_ORDERKEY from test8 group by id having sum(id)>10000)
mysqlRes:[]
gaeaRes:[]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where O_CUSTKEY not between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaRes:[[88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test8 a where not (O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300') group by o_custkey
mysqlRes:[[88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]
gaeaRes:[[88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1] [12000 CUSTKEY_980 1]]

sql:select C_NAME from sbtest1.test6 where C_NAME like 'A__A'
mysqlRes:[]
gaeaRes:[]

sql:select C_NAME from sbtest1.test6 where C_NAME like 'm___i'
mysqlRes:[[marui]]
gaeaRes:[[marui]]

sql:select C_NAME from sbtest1.test6 where C_NAME like 'ch___i%%'
mysqlRes:[[chenxiao] [chenqi]]
gaeaRes:[[chenxiao] [chenqi]]

sql:select C_NAME from sbtest1.test6 where C_NAME like 'ch___i%%' ESCAPE 'i'
mysqlRes:[]
gaeaRes:[]

sql:select count(*) from sbtest1.test6 where C_NAME not like 'chen%'
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select count(*) from sbtest1.test6 where not (C_NAME like 'chen%')
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select C_NAME from sbtest1.test6 where C_NAME like binary 'chen%'
mysqlRes:[[chenxiao] [chenqi]]
gaeaRes:[[chenxiao] [chenqi]]

sql:select C_NAME from sbtest1.test6 where C_NAME like 'chen%'
mysqlRes:[[chenxiao] [chenqi]]
gaeaRes:[[chenxiao] [chenqi]]

sql:select C_NAME from sbtest1.test6 where  C_NAME  like concat('%','AM','%')
mysqlRes:[]
gaeaRes:[]

sql:select C_NAME from sbtest1.test6 where  C_NAME  like concat('%','en','%')
mysqlRes:[[chenxiao] [chenqi] [huachen]]
gaeaRes:[[chenxiao] [chenqi] [huachen]]

sql:select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132'
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:    (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' order by O_ORDERKEY) union (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by O_ORDERKEY)
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by O_ORDERKEY
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by 1, O_ORDERKEY desc/*allow_diff_sequence*/
mysqlRes:[[CUSTKEY_003 ORDERKEY_002] [CUSTKEY_003 ORDERKEY_001] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_002] [CUSTKEY_003 ORDERKEY_001] [CUSTKEY_132 ORDERKEY_005]]

sql:select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' union select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132' order by 2,1
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:    (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_330') union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_420') union (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY<'CUSTKEY_200') order by 2,1/*allow_diff_sequence*/;
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_111 ORDERKEY_004] [CUSTKEY_132 ORDERKEY_005] [CUSTKEY_420 ORDERKEY_006] [CUSTKEY_012 ORDERKEY_011]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_111 ORDERKEY_004] [CUSTKEY_132 ORDERKEY_005] [CUSTKEY_420 ORDERKEY_006] [CUSTKEY_012 ORDERKEY_011]]

sql:(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') union (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') order by O_ORDERKEY;
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') order by O_ORDERKEY limit 5;
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003' order by O_ORDERKEY limit 5) union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_980' order by O_ORDERKEY limit 5);
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_980 ORDERKEY_007]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_980 ORDERKEY_007]]

sql:(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') union all (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') order by O_ORDERKEY limit 5
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:(select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_132') union DISTINCT (select O_CUSTKEY,O_ORDERKEY from test8 a where O_CUSTKEY ='CUSTKEY_003') order by O_ORDERKEY limit 5
mysqlRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]
gaeaRes:[[CUSTKEY_003 ORDERKEY_001] [CUSTKEY_003 ORDERKEY_002] [CUSTKEY_132 ORDERKEY_005]]

sql:select C_ORDERKEY,C_CUSTKEY,C_NAME from test1,sbtest1.test6 where C_CUSTKEY=c_CUSTKEY and C_ORDERKEY<'ORDERKEY_006'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]

sql:select count(*) from test6 as a where a.id >= any(select id from sbtest1.test6)
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from test6 as a where a.id >= any(select avg(id) from sbtest1.test6 group by id)
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from test6 as a where a.id in(select id from sbtest1.test6)
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from test6 as a where a.id not in(select id from sbtest1.test6)
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select count(*) from test6 as a where not (a.id in(select id from sbtest1.test6))
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select count(*) from test6 as a where a.id not in(select id from sbtest1.test6 where C_ORDERKEY in(1,2))
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from test6 as a where not (a.id in(select id from sbtest1.test6 where C_ORDERKEY in(1,2)))
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from test6 as a where a.id =some(select id from sbtest1.test6)
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from test6 as a where a.id != any(select id from sbtest1.test6)
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select count(*) from (select O_CUSTKEY,count(O_CUSTKEY) as counts from test8 group by O_CUSTKEY) as a where counts<10 group by counts
mysqlRes:[[6] [1]]
gaeaRes:[[6] [1]]

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 INNER join sbtest1.test6 where O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 INNER join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_006'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from test8 CROSS join sbtest1.test6 on O_CUSTKEY=c_CUSTKEY and O_ORDERKEY<'ORDERKEY_006'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan]]

sql:select test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME from test6 CROSS join sbtest1.test8 using(id) order by test8.O_ORDERKEY,test8.O_CUSTKEY,C_NAME
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 xiaojuan] [ORDERKEY_004 CUSTKEY_111 chenqi] [ORDERKEY_005 CUSTKEY_132 marui] [ORDERKEY_007 CUSTKEY_980 yanglu]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 xiaojuan] [ORDERKEY_004 CUSTKEY_111 chenqi] [ORDERKEY_005 CUSTKEY_132 marui] [ORDERKEY_007 CUSTKEY_980 yanglu]]

sql:select O_ORDERKEY,O_CUSTKEY,C_NAME from sbtest1.test6 as a STRAIGHT_JOIN test8 b where b.O_CUSTKEY=a.c_CUSTKEY and O_ORDERKEY<'ORDERKEY_007'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]

sql:select b.O_ORDERKEY,b.O_CUSTKEY,a.C_NAME from sbtest1.test6 a STRAIGHT_JOIN test8 b on b.O_CUSTKEY=a.c_CUSTKEY and b.O_ORDERKEY<'ORDERKEY_007'
mysqlRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003 chenxiao] [ORDERKEY_002 CUSTKEY_003 chenxiao] [ORDERKEY_004 CUSTKEY_111 wangye] [ORDERKEY_005 CUSTKEY_132 xiaojuan] [ORDERKEY_006 CUSTKEY_420 yanglu]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left join test8 b on b.O_CUSTKEY=a.C_CUSTKEY and a.C_CUSTKEY<'CUSTKEY_300'
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left join test8 b using(id) where a.C_CUSTKEY<'CUSTKEY_300'
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_111 wangye NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_111 wangye NULL]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test6 as a left OUTER join test8 b using(id)
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_333 chenqi ORDERKEY_004] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_420 yanglu ORDERKEY_007] [CUSTKEY_111 wangye NULL] [CUSTKEY_980 huachen NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_333 chenqi ORDERKEY_004] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_420 yanglu ORDERKEY_007] [CUSTKEY_111 wangye NULL] [CUSTKEY_980 huachen NULL]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a using(id) where a.c_CUSTKEY<'CUSTKEY_300'
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_111 wangye NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_111 wangye NULL]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a  using(id)
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_333 chenqi ORDERKEY_004] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_420 yanglu ORDERKEY_007] [CUSTKEY_111 wangye NULL] [CUSTKEY_980 huachen NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_132 xiaojuan ORDERKEY_002] [CUSTKEY_333 chenqi ORDERKEY_004] [CUSTKEY_012 marui ORDERKEY_005] [CUSTKEY_420 yanglu ORDERKEY_007] [CUSTKEY_111 wangye NULL] [CUSTKEY_980 huachen NULL]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from test8 b right join sbtest1.test6 as a on b.O_CUSTKEY=a.C_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300'
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by 2
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY), sums
mysqlRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [88888888 CUSTKEY_333 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [88888888 CUSTKEY_333 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by (select C_ORDERKEY from sbtest1.test7 where c_custkey=o_custkey) asc
mysqlRes:[[300000 CUSTKEY_003 2] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1] [323456 CUSTKEY_012 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [500 CUSTKEY_111 1] [88888888 CUSTKEY_333 1] [323456 CUSTKEY_012 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by counts asc,2 desc
mysqlRes:[[88888888 CUSTKEY_333 1] [500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[88888888 CUSTKEY_333 1] [500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_013' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by o_custkey
mysqlRes:[]
gaeaRes:[]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by 2
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY)
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by (select C_ORDERKEY from sbtest1.test7 where c_custkey=o_custkey) asc
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by counts asc,2 desc
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,2 desc limit 2
mysqlRes:[[88888888 CUSTKEY_333 1] [500 CUSTKEY_111 1]]
gaeaRes:[[88888888 CUSTKEY_333 1] [500 CUSTKEY_111 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,3
mysqlRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_003' and 'CUSTKEY_500' group by 2 asc having sums>400 order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1
mysqlRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[500 CUSTKEY_111 1] [323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,2 desc limit 10
mysqlRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]
gaeaRes:[[323456 CUSTKEY_012 1] [300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 1,10
mysqlRes:[[300000 CUSTKEY_003 2]]
gaeaRes:[[300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_050' group by 2 asc having O_CUSTKEY<(select max(c_custkey) from sbtest1.test7) order by count(O_ORDERKEY) asc,O_CUSTKEY desc limit 10 offset 1
mysqlRes:[[300000 CUSTKEY_003 2]]
gaeaRes:[[300000 CUSTKEY_003 2]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts,MYDATE from test9 use index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]
gaeaRes:[[300000 CUSTKEY_003 2 2014-10-22] [323456 CUSTKEY_012 1 1992-08-22] [500 CUSTKEY_111 1 2008-01-05] [100 CUSTKEY_132 1 1992-06-28]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use key (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_800' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use key (ORDERS_FK1) ignore index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_800' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1] [88888888 CUSTKEY_333 1] [231 CUSTKEY_420 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 ignore key (ORDERS_FK1) force index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 ignore key (ORDERS_FK1) force index (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 use index for order by (ORDERS_FK1) ignore index for group by (ORDERS_FK1) where O_CUSTKEY between 1 and 50 group by o_custkey
mysqlRes:[]
gaeaRes:[]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 as a use index for group by (ORDERS_FK1) ignore index for join (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]

sql:select sum(O_TOTALPRICE) as sums,O_CUSTKEY,count(O_ORDERKEY) counts from test9 a use index for join (ORDERS_FK1) ignore index for join (ORDERS_FK1) where O_CUSTKEY between 'CUSTKEY_002' and 'CUSTKEY_300' group by o_custkey
mysqlRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]
gaeaRes:[[300000 CUSTKEY_003 2] [323456 CUSTKEY_012 1] [500 CUSTKEY_111 1] [100 CUSTKEY_132 1]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test7 as a left join test9 b ignore index for join(ORDERS_FK1) on b.O_CUSTKEY=a.c_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300'
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]

sql:select a.c_CUSTKEY,a.C_NAME,b.O_ORDERKEY from sbtest1.test7 as a force index for join(primary) left join test9 b ignore index for join(ORDERS_FK1) on b.O_CUSTKEY=a.c_CUSTKEY and a.c_CUSTKEY<'CUSTKEY_300'
mysqlRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]
gaeaRes:[[CUSTKEY_003 chenxiao ORDERKEY_001] [CUSTKEY_003 chenxiao ORDERKEY_002] [CUSTKEY_111 wangye ORDERKEY_004] [CUSTKEY_132 xiaojuan ORDERKEY_005] [CUSTKEY_012 marui ORDERKEY_011] [CUSTKEY_333 chenqi NULL] [CUSTKEY_420 yanglu NULL] [CUSTKEY_980 huachen NULL]]

sql:select UPPER((select C_NAME FROM sbtest1.test7 limit 1)) FROM sbtest1.test7 limit 1
mysqlRes:[[CHENXIAO]]
gaeaRes:[[CHENXIAO]]

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY=(select min(C_CUSTKEY) from sbtest1.test7)
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY) from sbtest1.test7)
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003]]

sql:select O_ORDERKEY,O_CUSTKEY from test9 as a where a.O_CUSTKEY<=(select min(C_CUSTKEY)+1 from sbtest1.test7)
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003] [ORDERKEY_011 CUSTKEY_012] [ORDERKEY_004 CUSTKEY_111] [ORDERKEY_005 CUSTKEY_132] [ORDERKEY_010 CUSTKEY_333] [ORDERKEY_006 CUSTKEY_420] [ORDERKEY_007 CUSTKEY_980]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003] [ORDERKEY_011 CUSTKEY_012] [ORDERKEY_004 CUSTKEY_111] [ORDERKEY_005 CUSTKEY_132] [ORDERKEY_010 CUSTKEY_333] [ORDERKEY_006 CUSTKEY_420] [ORDERKEY_007 CUSTKEY_980]]

sql:select count(*) from sbtest1.test7 as a where a.c_CUSTKEY=(select max(C_CUSTKEY) from test9 where C_CUSTKEY=a.C_CUSTKEY)
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select C_CUSTKEY  from sbtest1.test7 as a where (select count(*) from test9 where O_CUSTKEY=a.C_CUSTKEY)=2
mysqlRes:[[CUSTKEY_003]]
gaeaRes:[[CUSTKEY_003]]

sql:select count(*) from test9 as a where a.id <> all(select id from sbtest1.test7)
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:select count(*) from test9 as a where 56000< all(select id from sbtest1.test7)
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select count(*) from sbtest1.test7 as a where 2>all(select count(*) from test9 where O_CUSTKEY=C_CUSTKEY)
mysqlRes:[[6]]
gaeaRes:[[6]]

sql:select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where (a.O_ORDERKEY,O_CUSTKEY)=(select c_ORDERKEY,c_CUSTKEY from sbtest1.test7 where c_name='yanglu')
mysqlRes:[[6 CUSTKEY_420 ORDERKEY_006 231]]
gaeaRes:[[6 CUSTKEY_420 ORDERKEY_006 231]]

sql:select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where (a.O_ORDERKEY,a.O_CUSTKEY) in (select c_ORDERKEY,c_CUSTKEY from sbtest1.test7) order by id,O_ORDERKEY,O_CUSTKEY,O_TOTALPRICE,MYDATE
mysqlRes:[[1 CUSTKEY_003 ORDERKEY_001 200000] [4 CUSTKEY_111 ORDERKEY_004 500] [5 CUSTKEY_132 ORDERKEY_005 100] [6 CUSTKEY_420 ORDERKEY_006 231] [7 CUSTKEY_980 ORDERKEY_007 12000] [10 CUSTKEY_333 ORDERKEY_010 88888888] [11 CUSTKEY_012 ORDERKEY_011 323456]]
gaeaRes:[[1 CUSTKEY_003 ORDERKEY_001 200000] [4 CUSTKEY_111 ORDERKEY_004 500] [5 CUSTKEY_132 ORDERKEY_005 100] [6 CUSTKEY_420 ORDERKEY_006 231] [7 CUSTKEY_980 ORDERKEY_007 12000] [10 CUSTKEY_333 ORDERKEY_010 88888888] [11 CUSTKEY_012 ORDERKEY_011 323456]]

sql:select id,O_CUSTKEY,O_ORDERKEY,O_TOTALPRICE from test9 a where exists(select * from sbtest1.test7 where a.O_CUSTKEY=sbtest1.test7.C_CUSTKEY)
mysqlRes:[[1 CUSTKEY_003 ORDERKEY_001 200000] [2 CUSTKEY_003 ORDERKEY_002 100000] [4 CUSTKEY_111 ORDERKEY_004 500] [5 CUSTKEY_132 ORDERKEY_005 100] [6 CUSTKEY_420 ORDERKEY_006 231] [7 CUSTKEY_980 ORDERKEY_007 12000] [10 CUSTKEY_333 ORDERKEY_010 88888888] [11 CUSTKEY_012 ORDERKEY_011 323456]]
gaeaRes:[[1 CUSTKEY_003 ORDERKEY_001 200000] [2 CUSTKEY_003 ORDERKEY_002 100000] [4 CUSTKEY_111 ORDERKEY_004 500] [5 CUSTKEY_132 ORDERKEY_005 100] [6 CUSTKEY_420 ORDERKEY_006 231] [7 CUSTKEY_980 ORDERKEY_007 12000] [10 CUSTKEY_333 ORDERKEY_010 88888888] [11 CUSTKEY_012 ORDERKEY_011 323456]]

sql:select distinct O_ORDERKEY,O_CUSTKEY from test9 a where exists(select * from sbtest1.test7 where a.O_CUSTKEY=sbtest1.test7.C_CUSTKEY)
mysqlRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003] [ORDERKEY_011 CUSTKEY_012] [ORDERKEY_004 CUSTKEY_111] [ORDERKEY_005 CUSTKEY_132] [ORDERKEY_010 CUSTKEY_333] [ORDERKEY_006 CUSTKEY_420] [ORDERKEY_007 CUSTKEY_980]]
gaeaRes:[[ORDERKEY_001 CUSTKEY_003] [ORDERKEY_002 CUSTKEY_003] [ORDERKEY_011 CUSTKEY_012] [ORDERKEY_004 CUSTKEY_111] [ORDERKEY_005 CUSTKEY_132] [ORDERKEY_010 CUSTKEY_333] [ORDERKEY_006 CUSTKEY_420] [ORDERKEY_007 CUSTKEY_980]]

sql:select count(*) from test9 a where not exists(select * from sbtest1.test7 where a.O_CUSTKEY=sbtest1.test7.C_CUSTKEY)
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:use  sbtest1;
mysqlRes:[]
gaeaRes:[]

sql:/* some comments */ SELECT CONVERT('111', SIGNED);
mysqlRes:[[111]]
gaeaRes:[[111]]

sql:/* some comments */ /*comment*/ SELECT CONVERT('111', SIGNED);
mysqlRes:[[111]]
gaeaRes:[[111]]

sql:SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1 INNER JOIN t2 where t1.col1 = t2.col1;
mysqlRes:[[1 aa 5 1 aa 5] [6 aa 5 1 aa 5] [2 bb 10 2 bb 10] [7 bb 10 2 bb 10] [3 cc 15 3 cc 15] [8 cc 15 3 cc 15] [4 dd 20 4 dd 20] [9 dd 20 4 dd 20] [5 ee 30 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [6 aa 5 6 aa 5] [2 bb 10 7 bb 10] [7 bb 10 7 bb 10] [3 cc 15 8 cc 15] [8 cc 15 8 cc 15] [4 dd 20 9 dd 20] [9 dd 20 9 dd 20] [5 ee 30 10 ee 30] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [6 aa 5 1 aa 5] [2 bb 10 2 bb 10] [7 bb 10 2 bb 10] [3 cc 15 3 cc 15] [8 cc 15 3 cc 15] [4 dd 20 4 dd 20] [9 dd 20 4 dd 20] [5 ee 30 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [6 aa 5 6 aa 5] [2 bb 10 7 bb 10] [7 bb 10 7 bb 10] [3 cc 15 8 cc 15] [8 cc 15 8 cc 15] [4 dd 20 9 dd 20] [9 dd 20 9 dd 20] [5 ee 30 10 ee 30] [10 ee 30 10 ee 30]]

sql:SELECT /*comment*/ CONVERT('111', SIGNED) ;
mysqlRes:[[111]]
gaeaRes:[[111]]

sql:SELECT CONVERT('111', /*comment*/ SIGNED) ;
mysqlRes:[[111]]
gaeaRes:[[111]]

sql:SELECT CONVERT('111', SIGNED) /*comment*/;
mysqlRes:[[111]]
gaeaRes:[[111]]

sql:/*comment*/ /*comment*/ select col1 /* this is a comment */ from t;
mysqlRes:[[aa] [aa] [bb] [bb] [cc] [cc] [dd] [dd] [ee] [ee]]
gaeaRes:[[aa] [aa] [bb] [bb] [cc] [cc] [dd] [dd] [ee] [ee]]

sql:SELECT /*!40001 SQL_NO_CACHE */ * FROM t WHERE 1 limit 0, 2000;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select b'';
mysqlRes:[[]]
gaeaRes:[[]]

sql:select B'';
mysqlRes:[[]]
gaeaRes:[[]]

sql:SELECT * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:SELECT * FROM t AS u;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:SELECT * FROM t1, t2;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]

sql:SELECT * FROM t1 AS u, t2;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]

sql:SELECT * FROM t1, t2 AS u;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]

sql:SELECT * FROM t1 AS u, t2 AS v;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]

sql:SELECT * FROM t, t1, t2;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]

sql:SELECT DISTINCT * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:SELECT DISTINCTROW * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:SELECT ALL * FROM t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:SELECT id FROM t;
mysqlRes:[[1] [6] [2] [7] [3] [8] [4] [9] [5] [10]]
gaeaRes:[[1] [6] [2] [7] [3] [8] [4] [9] [5] [10]]

sql:SELECT * FROM t WHERE 1 = 1;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select 1 from t;
mysqlRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]
gaeaRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]

sql:select 1 from t limit 1;
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select 1 from t where not exists (select 2);
mysqlRes:[]
gaeaRes:[]

sql:select * from sbtest1.t1 where id > 4 and id <=8 order by col1 desc;
mysqlRes:[[5 ee 30] [8 cc 15] [7 bb 10] [6 aa 5]]
gaeaRes:[[5 ee 30] [8 cc 15] [7 bb 10] [6 aa 5]]

sql:select 1 as a from t order by a;
mysqlRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]
gaeaRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]

sql:select * from sbtest1.t1 where id > 1 order by id desc limit 10;
mysqlRes:[[10 ee 30] [9 dd 20] [8 cc 15] [7 bb 10] [6 aa 5] [5 ee 30] [4 dd 20] [3 cc 15] [2 bb 10]]
gaeaRes:[[10 ee 30] [9 dd 20] [8 cc 15] [7 bb 10] [6 aa 5] [5 ee 30] [4 dd 20] [3 cc 15] [2 bb 10]]

sql:select * from sbtest1.t2 where id < 0;
mysqlRes:[]
gaeaRes:[]

sql:select 1 as a from t where 1 < any (select 2) order by a;
mysqlRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]
gaeaRes:[[1] [1] [1] [1] [1] [1] [1] [1] [1] [1]]

sql:select 1 order by 1;
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT * from t for update;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:SELECT * from t lock in share mode;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:DO 1;
mysqlRes:[]
gaeaRes:[]

sql:DO 1, sleep(1);
mysqlRes:[]
gaeaRes:[]

sql:SELECT * from t1, t2, t3;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]

sql:select * from t1 join t2 left join t3 on t2.id = t3.id;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]

sql:select * from t1 right join t2 on t1.id = t2.id left join t3 on t3.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 10 ee 30 10 ee 30]]

sql:select * from t1 right join t2 using (id) left join t3 using (id);
mysqlRes:[[1 aa 5 aa 5 aa 5] [2 bb 10 bb 10 bb 10] [3 cc 15 cc 15 cc 15] [4 dd 20 dd 20 dd 20] [5 ee 30 ee 30 ee 30] [6 aa 5 aa 5 aa 5] [7 bb 10 bb 10 bb 10] [8 cc 15 cc 15 cc 15] [9 dd 20 dd 20 dd 20] [10 ee 30 ee 30 ee 30]]
gaeaRes:[[1 aa 5 aa 5 aa 5] [2 bb 10 bb 10 bb 10] [3 cc 15 cc 15 cc 15] [4 dd 20 dd 20 dd 20] [5 ee 30 ee 30 ee 30] [6 aa 5 aa 5 aa 5] [7 bb 10 bb 10 bb 10] [8 cc 15 cc 15 cc 15] [9 dd 20 dd 20 dd 20] [10 ee 30 ee 30 ee 30]]

sql:select * from t1 natural join t2;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t1 natural right join t2;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t1 natural left outer join t2;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t1 straight_join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]

sql:select straight_join * from t1 join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]

sql:select straight_join * from t1 left join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]

sql:select straight_join * from t1 right join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]

sql:select straight_join * from t1 straight_join t2 on t1.id = t2.id;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 2 bb 10] [3 cc 15 3 cc 15] [4 dd 20 4 dd 20] [5 ee 30 5 ee 30] [6 aa 5 6 aa 5] [7 bb 10 7 bb 10] [8 cc 15 8 cc 15] [9 dd 20 9 dd 20] [10 ee 30 10 ee 30]]

sql:(select 1);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:ANALYZE TABLE t;
mysqlRes:[[sbtest1.t analyze status OK]]
gaeaRes:[[sbtest1.t analyze status OK]]

sql:SHOW VARIABLES LIKE 'character_set_results';
mysqlRes:[[character_set_results utf8mb4]]
gaeaRes:[[character_set_results utf8mb4]]

sql:SHOW GLOBAL VARIABLES LIKE 'character_set_results';
mysqlRes:[[character_set_results latin1]]
gaeaRes:[[character_set_results latin1]]

sql:SHOW SESSION VARIABLES LIKE 'character_set_results';
mysqlRes:[[character_set_results utf8mb4]]
gaeaRes:[[character_set_results utf8mb4]]

sql:SHOW GLOBAL VARIABLES;
mysqlRes:[[auto_increment_increment 1] [auto_increment_offset 1] [autocommit ON] [automatic_sp_privileges ON] [avoid_temporal_upgrade OFF] [back_log 80] [basedir /usr/] [big_tables OFF] [bind_address *] [binlog_cache_size 32768] [binlog_checksum CRC32] [binlog_direct_non_transactional_updates OFF] [binlog_error_action ABORT_SERVER] [binlog_format ROW] [binlog_group_commit_sync_delay 0] [binlog_group_commit_sync_no_delay_count 0] [binlog_gtid_simple_recovery ON] [binlog_max_flush_queue_time 0] [binlog_order_commits ON] [binlog_row_image FULL] [binlog_rows_query_log_events OFF] [binlog_stmt_cache_size 32768] [binlog_transaction_dependency_history_size 25000] [binlog_transaction_dependency_tracking COMMIT_ORDER] [block_encryption_mode aes-128-ecb] [bulk_insert_buffer_size 8388608] [character_set_client latin1] [character_set_connection latin1] [character_set_database latin1] [character_set_filesystem binary] [character_set_results latin1] [character_set_server latin1] [character_set_system utf8] [character_sets_dir /usr/share/mysql/charsets/] [check_proxy_users OFF] [collation_connection latin1_swedish_ci] [collation_database latin1_swedish_ci] [collation_server latin1_swedish_ci] [completion_type NO_CHAIN] [concurrent_insert AUTO] [connect_timeout 10] [core_file OFF] [datadir /var/lib/mysql/] [date_format %Y-%m-%d] [datetime_format %Y-%m-%d %H:%i:%s] [default_authentication_plugin mysql_native_password] [default_password_lifetime 0] [default_storage_engine InnoDB] [default_tmp_storage_engine InnoDB] [default_week_format 0] [delay_key_write ON] [delayed_insert_limit 100] [delayed_insert_timeout 300] [delayed_queue_size 1000] [disabled_storage_engines ] [disconnect_on_expired_password ON] [div_precision_increment 4] [end_markers_in_json OFF] [enforce_gtid_consistency OFF] [eq_range_index_dive_limit 200] [event_scheduler OFF] [expire_logs_days 0] [explicit_defaults_for_timestamp OFF] [flush OFF] [flush_time 0] [foreign_key_checks ON] [ft_boolean_syntax + -><()~*:""&|] [ft_max_word_len 84] [ft_min_word_len 4] [ft_query_expansion_limit 20] [ft_stopword_file (built-in)] [general_log ON] [general_log_file /var/lib/mysql/a8244020552d.log] [group_concat_max_len 1024] [gtid_executed ] [gtid_executed_compression_period 1000] [gtid_mode OFF] [gtid_owned ] [gtid_purged ] [have_compress YES] [have_crypt YES] [have_dynamic_loading YES] [have_geometry YES] [have_openssl DISABLED] [have_profiling YES] [have_query_cache YES] [have_rtree_keys YES] [have_ssl DISABLED] [have_statement_timeout YES] [have_symlink YES] [host_cache_size 279] [hostname a8244020552d] [ignore_builtin_innodb OFF] [ignore_db_dirs ] [init_connect ] [init_file ] [init_slave ] [innodb_adaptive_flushing ON] [innodb_adaptive_flushing_lwm 10] [innodb_adaptive_hash_index ON] [innodb_adaptive_hash_index_parts 8] [innodb_adaptive_max_sleep_delay 150000] [innodb_api_bk_commit_interval 5] [innodb_api_disable_rowlock OFF] [innodb_api_enable_binlog OFF] [innodb_api_enable_mdl OFF] [innodb_api_trx_level 0] [innodb_autoextend_increment 64] [innodb_autoinc_lock_mode 1] [innodb_buffer_pool_chunk_size 134217728] [innodb_buffer_pool_dump_at_shutdown ON] [innodb_buffer_pool_dump_now OFF] [innodb_buffer_pool_dump_pct 25] [innodb_buffer_pool_filename ib_buffer_pool] [innodb_buffer_pool_instances 1] [innodb_buffer_pool_load_abort OFF] [innodb_buffer_pool_load_at_startup ON] [innodb_buffer_pool_load_now OFF] [innodb_buffer_pool_size 134217728] [innodb_change_buffer_max_size 25] [innodb_change_buffering all] [innodb_checksum_algorithm crc32] [innodb_checksums ON] [innodb_cmp_per_index_enabled OFF] [innodb_commit_concurrency 0] [innodb_compression_failure_threshold_pct 5] [innodb_compression_level 6] [innodb_compression_pad_pct_max 50] [innodb_concurrency_tickets 5000] [innodb_data_file_path ibdata1:12M:autoextend] [innodb_data_home_dir ] [innodb_deadlock_detect ON] [innodb_default_row_format dynamic] [innodb_disable_sort_file_cache OFF] [innodb_doublewrite ON] [innodb_fast_shutdown 1] [innodb_file_format Barracuda] [innodb_file_format_check ON] [innodb_file_format_max Barracuda] [innodb_file_per_table ON] [innodb_fill_factor 100] [innodb_flush_log_at_timeout 1] [innodb_flush_log_at_trx_commit 1] [innodb_flush_method ] [innodb_flush_neighbors 1] [innodb_flush_sync ON] [innodb_flushing_avg_loops 30] [innodb_force_load_corrupted OFF] [innodb_force_recovery 0] [innodb_ft_aux_table ] [innodb_ft_cache_size 8000000] [innodb_ft_enable_diag_print OFF] [innodb_ft_enable_stopword ON] [innodb_ft_max_token_size 84] [innodb_ft_min_token_size 3] [innodb_ft_num_word_optimize 2000] [innodb_ft_result_cache_limit 2000000000] [innodb_ft_server_stopword_table ] [innodb_ft_sort_pll_degree 2] [innodb_ft_total_cache_size 640000000] [innodb_ft_user_stopword_table ] [innodb_io_capacity 200] [innodb_io_capacity_max 2000] [innodb_large_prefix ON] [innodb_lock_wait_timeout 50] [innodb_locks_unsafe_for_binlog OFF] [innodb_log_buffer_size 16777216] [innodb_log_checksums ON] [innodb_log_compressed_pages ON] [innodb_log_file_size 50331648] [innodb_log_files_in_group 2] [innodb_log_group_home_dir ./] [innodb_log_write_ahead_size 8192] [innodb_lru_scan_depth 1024] [innodb_max_dirty_pages_pct 75.000000] [innodb_max_dirty_pages_pct_lwm 0.000000] [innodb_max_purge_lag 0] [innodb_max_purge_lag_delay 0] [innodb_max_undo_log_size 1073741824] [innodb_monitor_disable ] [innodb_monitor_enable ] [innodb_monitor_reset ] [innodb_monitor_reset_all ] [innodb_old_blocks_pct 37] [innodb_old_blocks_time 1000] [innodb_online_alter_log_max_size 134217728] [innodb_open_files 2000] [innodb_optimize_fulltext_only OFF] [innodb_page_cleaners 1] [innodb_page_size 16384] [innodb_print_all_deadlocks OFF] [innodb_purge_batch_size 300] [innodb_purge_rseg_truncate_frequency 128] [innodb_purge_threads 4] [innodb_random_read_ahead OFF] [innodb_read_ahead_threshold 56] [innodb_read_io_threads 4] [innodb_read_only OFF] [innodb_replication_delay 0] [innodb_rollback_on_timeout OFF] [innodb_rollback_segments 128] [innodb_sort_buffer_size 1048576] [innodb_spin_wait_delay 6] [innodb_stats_auto_recalc ON] [innodb_stats_include_delete_marked OFF] [innodb_stats_method nulls_equal] [innodb_stats_on_metadata OFF] [innodb_stats_persistent ON] [innodb_stats_persistent_sample_pages 20] [innodb_stats_sample_pages 8] [innodb_stats_transient_sample_pages 8] [innodb_status_output OFF] [innodb_status_output_locks OFF] [innodb_strict_mode ON] [innodb_support_xa ON] [innodb_sync_array_size 1] [innodb_sync_spin_loops 30] [innodb_table_locks ON] [innodb_temp_data_file_path ibtmp1:12M:autoextend] [innodb_thread_concurrency 0] [innodb_thread_sleep_delay 10000] [innodb_tmpdir ] [innodb_undo_directory ./] [innodb_undo_log_truncate OFF] [innodb_undo_logs 128] [innodb_undo_tablespaces 0] [innodb_use_native_aio OFF] [innodb_version 5.7.26] [innodb_write_io_threads 4] [interactive_timeout 2880000] [internal_tmp_disk_storage_engine InnoDB] [join_buffer_size 262144] [keep_files_on_create OFF] [key_buffer_size 8388608] [key_cache_age_threshold 300] [key_cache_block_size 1024] [key_cache_division_limit 100] [keyring_operations ON] [large_files_support ON] [large_page_size 0] [large_pages OFF] [lc_messages en_US] [lc_messages_dir /usr/share/mysql/] [lc_time_names en_US] [license GPL] [local_infile ON] [lock_wait_timeout 31536000] [locked_in_memory OFF] [log_bin OFF] [log_bin_basename ] [log_bin_index ] [log_bin_trust_function_creators OFF] [log_bin_use_v1_row_events OFF] [log_builtin_as_identified_by_password OFF] [log_error stderr] [log_error_verbosity 3] [log_output TABLE] [log_queries_not_using_indexes OFF] [log_slave_updates OFF] [log_slow_admin_statements OFF] [log_slow_slave_statements OFF] [log_statements_unsafe_for_binlog ON] [log_syslog OFF] [log_syslog_facility daemon] [log_syslog_include_pid ON] [log_syslog_tag ] [log_throttle_queries_not_using_indexes 0] [log_timestamps UTC] [log_warnings 2] [long_query_time 10.000000] [low_priority_updates OFF] [lower_case_file_system OFF] [lower_case_table_names 0] [master_info_repository FILE] [master_verify_checksum OFF] [max_allowed_packet 4194304] [max_binlog_cache_size 18446744073709547520] [max_binlog_size 1073741824] [max_binlog_stmt_cache_size 18446744073709547520] [max_connect_errors 100] [max_connections 151] [max_delayed_threads 20] [max_digest_length 1024] [max_error_count 64] [max_execution_time 0] [max_heap_table_size 16777216] [max_insert_delayed_threads 20] [max_join_size 18446744073709551615] [max_length_for_sort_data 1024] [max_points_in_geometry 65536] [max_prepared_stmt_count 16382] [max_relay_log_size 0] [max_seeks_for_key 4294967295] [max_sort_length 1024] [max_sp_recursion_depth 0] [max_tmp_tables 32] [max_user_connections 0] [max_write_lock_count 4294967295] [metadata_locks_cache_size 1024] [metadata_locks_hash_instances 8] [min_examined_row_limit 0] [multi_range_count 256] [myisam_data_pointer_size 6] [myisam_max_sort_file_size 2146435072] [myisam_mmap_size 4294967295] [myisam_recover_options OFF] [myisam_repair_threads 1] [myisam_sort_buffer_size 8388608] [myisam_stats_method nulls_unequal] [myisam_use_mmap OFF] [mysql_native_password_proxy_users OFF] [net_buffer_length 16384] [net_read_timeout 30] [net_retry_count 10] [net_write_timeout 60] [new OFF] [ngram_token_size 2] [offline_mode OFF] [old OFF] [old_alter_table OFF] [old_passwords 0] [open_files_limit 1048576] [optimizer_prune_level 1] [optimizer_search_depth 62] [optimizer_switch index_merge=on,index_merge_union=on,index_merge_sort_union=on,index_merge_intersection=on,engine_condition_pushdown=on,index_condition_pushdown=on,mrr=on,mrr_cost_based=on,block_nested_loop=on,batched_key_access=off,materialization=on,semijoin=on,loosescan=on,firstmatch=on,duplicateweedout=on,subquery_materialization_cost_based=on,use_index_extensions=on,condition_fanout_filter=on,derived_merge=on] [optimizer_trace enabled=off,one_line=off] [optimizer_trace_features greedy_search=on,range_optimizer=on,dynamic_range=on,repeated_subselect=on] [optimizer_trace_limit 1] [optimizer_trace_max_mem_size 16384] [optimizer_trace_offset -1] [parser_max_mem_size 4294967295] [performance_schema ON] [performance_schema_accounts_size -1] [performance_schema_digests_size 10000] [performance_schema_events_stages_history_long_size 10000] [performance_schema_events_stages_history_size 10] [performance_schema_events_statements_history_long_size 10000] [performance_schema_events_statements_history_size 10] [performance_schema_events_transactions_history_long_size 10000] [performance_schema_events_transactions_history_size 10] [performance_schema_events_waits_history_long_size 10000] [performance_schema_events_waits_history_size 10] [performance_schema_hosts_size -1] [performance_schema_max_cond_classes 80] [performance_schema_max_cond_instances -1] [performance_schema_max_digest_length 1024] [performance_schema_max_file_classes 80] [performance_schema_max_file_handles 32768] [performance_schema_max_file_instances -1] [performance_schema_max_index_stat -1] [performance_schema_max_memory_classes 320] [performance_schema_max_metadata_locks -1] [performance_schema_max_mutex_classes 210] [performance_schema_max_mutex_instances -1] [performance_schema_max_prepared_statements_instances -1] [performance_schema_max_program_instances -1] [performance_schema_max_rwlock_classes 50] [performance_schema_max_rwlock_instances -1] [performance_schema_max_socket_classes 10] [performance_schema_max_socket_instances -1] [performance_schema_max_sql_text_length 1024] [performance_schema_max_stage_classes 150] [performance_schema_max_statement_classes 193] [performance_schema_max_statement_stack 10] [performance_schema_max_table_handles -1] [performance_schema_max_table_instances -1] [performance_schema_max_table_lock_stat -1] [performance_schema_max_thread_classes 50] [performance_schema_max_thread_instances -1] [performance_schema_session_connect_attrs_size 512] [performance_schema_setup_actors_size -1] [performance_schema_setup_objects_size -1] [performance_schema_users_size -1] [pid_file /var/lib/mysql/a8244020552d.pid] [plugin_dir /usr/lib/mysql/plugin/] [port 3306] [preload_buffer_size 32768] [profiling OFF] [profiling_history_size 15] [protocol_version 10] [query_alloc_block_size 8192] [query_cache_limit 1048576] [query_cache_min_res_unit 4096] [query_cache_size 1048576] [query_cache_type OFF] [query_cache_wlock_invalidate OFF] [query_prealloc_size 8192] [range_alloc_block_size 4096] [range_optimizer_max_mem_size 8388608] [rbr_exec_mode STRICT] [read_buffer_size 131072] [read_only OFF] [read_rnd_buffer_size 262144] [relay_log ] [relay_log_basename /var/lib/mysql/a8244020552d-relay-bin] [relay_log_index /var/lib/mysql/a8244020552d-relay-bin.index] [relay_log_info_file relay-log.info] [relay_log_info_repository FILE] [relay_log_purge ON] [relay_log_recovery OFF] [relay_log_space_limit 0] [report_host ] [report_password ] [report_port 3306] [report_user ] [require_secure_transport OFF] [rpl_stop_slave_timeout 31536000] [secure_auth ON] [secure_file_priv /var/lib/mysql-files/] [server_id 3379] [server_id_bits 32] [server_uuid 6b47d1d0-6bd7-11ee-8d59-0242c0a8ed04] [session_track_gtids OFF] [session_track_schema ON] [session_track_state_change OFF] [session_track_system_variables time_zone,autocommit,character_set_client,character_set_results,character_set_connection] [session_track_transaction_info OFF] [sha256_password_proxy_users OFF] [show_compatibility_56 OFF] [show_create_table_verbosity OFF] [show_old_temporals OFF] [skip_external_locking ON] [skip_name_resolve ON] [skip_networking OFF] [skip_show_database OFF] [slave_allow_batching OFF] [slave_checkpoint_group 512] [slave_checkpoint_period 300] [slave_compressed_protocol OFF] [slave_exec_mode STRICT] [slave_load_tmpdir /tmp] [slave_max_allowed_packet 1073741824] [slave_net_timeout 60] [slave_parallel_type DATABASE] [slave_parallel_workers 0] [slave_pending_jobs_size_max 16777216] [slave_preserve_commit_order OFF] [slave_rows_search_algorithms TABLE_SCAN,INDEX_SCAN] [slave_skip_errors OFF] [slave_sql_verify_checksum ON] [slave_transaction_retries 10] [slave_type_conversions ] [slow_launch_time 2] [slow_query_log OFF] [slow_query_log_file /var/lib/mysql/a8244020552d-slow.log] [socket /var/lib/mysql/mysql.sock] [sort_buffer_size 262144] [sql_auto_is_null OFF] [sql_big_selects ON] [sql_buffer_result OFF] [sql_log_off OFF] [sql_mode STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION] [sql_notes ON] [sql_quote_show_create ON] [sql_safe_updates OFF] [sql_select_limit 18446744073709551615] [sql_slave_skip_counter 0] [sql_warnings OFF] [ssl_ca ] [ssl_capath ] [ssl_cert ] [ssl_cipher ] [ssl_crl ] [ssl_crlpath ] [ssl_key ] [stored_program_cache 256] [super_read_only OFF] [sync_binlog 1] [sync_frm ON] [sync_master_info 10000] [sync_relay_log 10000] [sync_relay_log_info 10000] [system_time_zone UTC] [table_definition_cache 1400] [table_open_cache 2000] [table_open_cache_instances 16] [thread_cache_size 9] [thread_handling one-thread-per-connection] [thread_stack 196608] [time_format %H:%i:%s] [time_zone SYSTEM] [tls_version TLSv1,TLSv1.1] [tmp_table_size 16777216] [tmpdir /tmp] [transaction_alloc_block_size 8192] [transaction_isolation REPEATABLE-READ] [transaction_prealloc_size 4096] [transaction_read_only OFF] [transaction_write_set_extraction OFF] [tx_isolation REPEATABLE-READ] [tx_read_only OFF] [unique_checks ON] [updatable_views_with_limit YES] [version 5.7.26-1-log] [version_comment (Debian)] [version_compile_machine armv7l] [version_compile_os Linux] [wait_timeout 2880000]]
gaeaRes:[[auto_increment_increment 1] [auto_increment_offset 1] [autocommit ON] [automatic_sp_privileges ON] [avoid_temporal_upgrade OFF] [back_log 80] [basedir /usr/] [big_tables OFF] [bind_address *] [binlog_cache_size 32768] [binlog_checksum CRC32] [binlog_direct_non_transactional_updates OFF] [binlog_error_action ABORT_SERVER] [binlog_format ROW] [binlog_group_commit_sync_delay 0] [binlog_group_commit_sync_no_delay_count 0] [binlog_gtid_simple_recovery ON] [binlog_max_flush_queue_time 0] [binlog_order_commits ON] [binlog_row_image FULL] [binlog_rows_query_log_events OFF] [binlog_stmt_cache_size 32768] [binlog_transaction_dependency_history_size 25000] [binlog_transaction_dependency_tracking COMMIT_ORDER] [block_encryption_mode aes-128-ecb] [bulk_insert_buffer_size 8388608] [character_set_client latin1] [character_set_connection latin1] [character_set_database latin1] [character_set_filesystem binary] [character_set_results latin1] [character_set_server latin1] [character_set_system utf8] [character_sets_dir /usr/share/mysql/charsets/] [check_proxy_users OFF] [collation_connection latin1_swedish_ci] [collation_database latin1_swedish_ci] [collation_server latin1_swedish_ci] [completion_type NO_CHAIN] [concurrent_insert AUTO] [connect_timeout 10] [core_file OFF] [datadir /var/lib/mysql/] [date_format %Y-%m-%d] [datetime_format %Y-%m-%d %H:%i:%s] [default_authentication_plugin mysql_native_password] [default_password_lifetime 0] [default_storage_engine InnoDB] [default_tmp_storage_engine InnoDB] [default_week_format 0] [delay_key_write ON] [delayed_insert_limit 100] [delayed_insert_timeout 300] [delayed_queue_size 1000] [disabled_storage_engines ] [disconnect_on_expired_password ON] [div_precision_increment 4] [end_markers_in_json OFF] [enforce_gtid_consistency OFF] [eq_range_index_dive_limit 200] [event_scheduler OFF] [expire_logs_days 0] [explicit_defaults_for_timestamp OFF] [flush OFF] [flush_time 0] [foreign_key_checks ON] [ft_boolean_syntax + -><()~*:""&|] [ft_max_word_len 84] [ft_min_word_len 4] [ft_query_expansion_limit 20] [ft_stopword_file (built-in)] [general_log ON] [general_log_file /var/lib/mysql/a8244020552d.log] [group_concat_max_len 1024] [gtid_executed ] [gtid_executed_compression_period 1000] [gtid_mode OFF] [gtid_owned ] [gtid_purged ] [have_compress YES] [have_crypt YES] [have_dynamic_loading YES] [have_geometry YES] [have_openssl DISABLED] [have_profiling YES] [have_query_cache YES] [have_rtree_keys YES] [have_ssl DISABLED] [have_statement_timeout YES] [have_symlink YES] [host_cache_size 279] [hostname a8244020552d] [ignore_builtin_innodb OFF] [ignore_db_dirs ] [init_connect ] [init_file ] [init_slave ] [innodb_adaptive_flushing ON] [innodb_adaptive_flushing_lwm 10] [innodb_adaptive_hash_index ON] [innodb_adaptive_hash_index_parts 8] [innodb_adaptive_max_sleep_delay 150000] [innodb_api_bk_commit_interval 5] [innodb_api_disable_rowlock OFF] [innodb_api_enable_binlog OFF] [innodb_api_enable_mdl OFF] [innodb_api_trx_level 0] [innodb_autoextend_increment 64] [innodb_autoinc_lock_mode 1] [innodb_buffer_pool_chunk_size 134217728] [innodb_buffer_pool_dump_at_shutdown ON] [innodb_buffer_pool_dump_now OFF] [innodb_buffer_pool_dump_pct 25] [innodb_buffer_pool_filename ib_buffer_pool] [innodb_buffer_pool_instances 1] [innodb_buffer_pool_load_abort OFF] [innodb_buffer_pool_load_at_startup ON] [innodb_buffer_pool_load_now OFF] [innodb_buffer_pool_size 134217728] [innodb_change_buffer_max_size 25] [innodb_change_buffering all] [innodb_checksum_algorithm crc32] [innodb_checksums ON] [innodb_cmp_per_index_enabled OFF] [innodb_commit_concurrency 0] [innodb_compression_failure_threshold_pct 5] [innodb_compression_level 6] [innodb_compression_pad_pct_max 50] [innodb_concurrency_tickets 5000] [innodb_data_file_path ibdata1:12M:autoextend] [innodb_data_home_dir ] [innodb_deadlock_detect ON] [innodb_default_row_format dynamic] [innodb_disable_sort_file_cache OFF] [innodb_doublewrite ON] [innodb_fast_shutdown 1] [innodb_file_format Barracuda] [innodb_file_format_check ON] [innodb_file_format_max Barracuda] [innodb_file_per_table ON] [innodb_fill_factor 100] [innodb_flush_log_at_timeout 1] [innodb_flush_log_at_trx_commit 1] [innodb_flush_method ] [innodb_flush_neighbors 1] [innodb_flush_sync ON] [innodb_flushing_avg_loops 30] [innodb_force_load_corrupted OFF] [innodb_force_recovery 0] [innodb_ft_aux_table ] [innodb_ft_cache_size 8000000] [innodb_ft_enable_diag_print OFF] [innodb_ft_enable_stopword ON] [innodb_ft_max_token_size 84] [innodb_ft_min_token_size 3] [innodb_ft_num_word_optimize 2000] [innodb_ft_result_cache_limit 2000000000] [innodb_ft_server_stopword_table ] [innodb_ft_sort_pll_degree 2] [innodb_ft_total_cache_size 640000000] [innodb_ft_user_stopword_table ] [innodb_io_capacity 200] [innodb_io_capacity_max 2000] [innodb_large_prefix ON] [innodb_lock_wait_timeout 50] [innodb_locks_unsafe_for_binlog OFF] [innodb_log_buffer_size 16777216] [innodb_log_checksums ON] [innodb_log_compressed_pages ON] [innodb_log_file_size 50331648] [innodb_log_files_in_group 2] [innodb_log_group_home_dir ./] [innodb_log_write_ahead_size 8192] [innodb_lru_scan_depth 1024] [innodb_max_dirty_pages_pct 75.000000] [innodb_max_dirty_pages_pct_lwm 0.000000] [innodb_max_purge_lag 0] [innodb_max_purge_lag_delay 0] [innodb_max_undo_log_size 1073741824] [innodb_monitor_disable ] [innodb_monitor_enable ] [innodb_monitor_reset ] [innodb_monitor_reset_all ] [innodb_old_blocks_pct 37] [innodb_old_blocks_time 1000] [innodb_online_alter_log_max_size 134217728] [innodb_open_files 2000] [innodb_optimize_fulltext_only OFF] [innodb_page_cleaners 1] [innodb_page_size 16384] [innodb_print_all_deadlocks OFF] [innodb_purge_batch_size 300] [innodb_purge_rseg_truncate_frequency 128] [innodb_purge_threads 4] [innodb_random_read_ahead OFF] [innodb_read_ahead_threshold 56] [innodb_read_io_threads 4] [innodb_read_only OFF] [innodb_replication_delay 0] [innodb_rollback_on_timeout OFF] [innodb_rollback_segments 128] [innodb_sort_buffer_size 1048576] [innodb_spin_wait_delay 6] [innodb_stats_auto_recalc ON] [innodb_stats_include_delete_marked OFF] [innodb_stats_method nulls_equal] [innodb_stats_on_metadata OFF] [innodb_stats_persistent ON] [innodb_stats_persistent_sample_pages 20] [innodb_stats_sample_pages 8] [innodb_stats_transient_sample_pages 8] [innodb_status_output OFF] [innodb_status_output_locks OFF] [innodb_strict_mode ON] [innodb_support_xa ON] [innodb_sync_array_size 1] [innodb_sync_spin_loops 30] [innodb_table_locks ON] [innodb_temp_data_file_path ibtmp1:12M:autoextend] [innodb_thread_concurrency 0] [innodb_thread_sleep_delay 10000] [innodb_tmpdir ] [innodb_undo_directory ./] [innodb_undo_log_truncate OFF] [innodb_undo_logs 128] [innodb_undo_tablespaces 0] [innodb_use_native_aio OFF] [innodb_version 5.7.26] [innodb_write_io_threads 4] [interactive_timeout 2880000] [internal_tmp_disk_storage_engine InnoDB] [join_buffer_size 262144] [keep_files_on_create OFF] [key_buffer_size 8388608] [key_cache_age_threshold 300] [key_cache_block_size 1024] [key_cache_division_limit 100] [keyring_operations ON] [large_files_support ON] [large_page_size 0] [large_pages OFF] [lc_messages en_US] [lc_messages_dir /usr/share/mysql/] [lc_time_names en_US] [license GPL] [local_infile ON] [lock_wait_timeout 31536000] [locked_in_memory OFF] [log_bin OFF] [log_bin_basename ] [log_bin_index ] [log_bin_trust_function_creators OFF] [log_bin_use_v1_row_events OFF] [log_builtin_as_identified_by_password OFF] [log_error stderr] [log_error_verbosity 3] [log_output TABLE] [log_queries_not_using_indexes OFF] [log_slave_updates OFF] [log_slow_admin_statements OFF] [log_slow_slave_statements OFF] [log_statements_unsafe_for_binlog ON] [log_syslog OFF] [log_syslog_facility daemon] [log_syslog_include_pid ON] [log_syslog_tag ] [log_throttle_queries_not_using_indexes 0] [log_timestamps UTC] [log_warnings 2] [long_query_time 10.000000] [low_priority_updates OFF] [lower_case_file_system OFF] [lower_case_table_names 0] [master_info_repository FILE] [master_verify_checksum OFF] [max_allowed_packet 4194304] [max_binlog_cache_size 18446744073709547520] [max_binlog_size 1073741824] [max_binlog_stmt_cache_size 18446744073709547520] [max_connect_errors 100] [max_connections 151] [max_delayed_threads 20] [max_digest_length 1024] [max_error_count 64] [max_execution_time 0] [max_heap_table_size 16777216] [max_insert_delayed_threads 20] [max_join_size 18446744073709551615] [max_length_for_sort_data 1024] [max_points_in_geometry 65536] [max_prepared_stmt_count 16382] [max_relay_log_size 0] [max_seeks_for_key 4294967295] [max_sort_length 1024] [max_sp_recursion_depth 0] [max_tmp_tables 32] [max_user_connections 0] [max_write_lock_count 4294967295] [metadata_locks_cache_size 1024] [metadata_locks_hash_instances 8] [min_examined_row_limit 0] [multi_range_count 256] [myisam_data_pointer_size 6] [myisam_max_sort_file_size 2146435072] [myisam_mmap_size 4294967295] [myisam_recover_options OFF] [myisam_repair_threads 1] [myisam_sort_buffer_size 8388608] [myisam_stats_method nulls_unequal] [myisam_use_mmap OFF] [mysql_native_password_proxy_users OFF] [net_buffer_length 16384] [net_read_timeout 30] [net_retry_count 10] [net_write_timeout 60] [new OFF] [ngram_token_size 2] [offline_mode OFF] [old OFF] [old_alter_table OFF] [old_passwords 0] [open_files_limit 1048576] [optimizer_prune_level 1] [optimizer_search_depth 62] [optimizer_switch index_merge=on,index_merge_union=on,index_merge_sort_union=on,index_merge_intersection=on,engine_condition_pushdown=on,index_condition_pushdown=on,mrr=on,mrr_cost_based=on,block_nested_loop=on,batched_key_access=off,materialization=on,semijoin=on,loosescan=on,firstmatch=on,duplicateweedout=on,subquery_materialization_cost_based=on,use_index_extensions=on,condition_fanout_filter=on,derived_merge=on] [optimizer_trace enabled=off,one_line=off] [optimizer_trace_features greedy_search=on,range_optimizer=on,dynamic_range=on,repeated_subselect=on] [optimizer_trace_limit 1] [optimizer_trace_max_mem_size 16384] [optimizer_trace_offset -1] [parser_max_mem_size 4294967295] [performance_schema ON] [performance_schema_accounts_size -1] [performance_schema_digests_size 10000] [performance_schema_events_stages_history_long_size 10000] [performance_schema_events_stages_history_size 10] [performance_schema_events_statements_history_long_size 10000] [performance_schema_events_statements_history_size 10] [performance_schema_events_transactions_history_long_size 10000] [performance_schema_events_transactions_history_size 10] [performance_schema_events_waits_history_long_size 10000] [performance_schema_events_waits_history_size 10] [performance_schema_hosts_size -1] [performance_schema_max_cond_classes 80] [performance_schema_max_cond_instances -1] [performance_schema_max_digest_length 1024] [performance_schema_max_file_classes 80] [performance_schema_max_file_handles 32768] [performance_schema_max_file_instances -1] [performance_schema_max_index_stat -1] [performance_schema_max_memory_classes 320] [performance_schema_max_metadata_locks -1] [performance_schema_max_mutex_classes 210] [performance_schema_max_mutex_instances -1] [performance_schema_max_prepared_statements_instances -1] [performance_schema_max_program_instances -1] [performance_schema_max_rwlock_classes 50] [performance_schema_max_rwlock_instances -1] [performance_schema_max_socket_classes 10] [performance_schema_max_socket_instances -1] [performance_schema_max_sql_text_length 1024] [performance_schema_max_stage_classes 150] [performance_schema_max_statement_classes 193] [performance_schema_max_statement_stack 10] [performance_schema_max_table_handles -1] [performance_schema_max_table_instances -1] [performance_schema_max_table_lock_stat -1] [performance_schema_max_thread_classes 50] [performance_schema_max_thread_instances -1] [performance_schema_session_connect_attrs_size 512] [performance_schema_setup_actors_size -1] [performance_schema_setup_objects_size -1] [performance_schema_users_size -1] [pid_file /var/lib/mysql/a8244020552d.pid] [plugin_dir /usr/lib/mysql/plugin/] [port 3306] [preload_buffer_size 32768] [profiling OFF] [profiling_history_size 15] [protocol_version 10] [query_alloc_block_size 8192] [query_cache_limit 1048576] [query_cache_min_res_unit 4096] [query_cache_size 1048576] [query_cache_type OFF] [query_cache_wlock_invalidate OFF] [query_prealloc_size 8192] [range_alloc_block_size 4096] [range_optimizer_max_mem_size 8388608] [rbr_exec_mode STRICT] [read_buffer_size 131072] [read_only OFF] [read_rnd_buffer_size 262144] [relay_log ] [relay_log_basename /var/lib/mysql/a8244020552d-relay-bin] [relay_log_index /var/lib/mysql/a8244020552d-relay-bin.index] [relay_log_info_file relay-log.info] [relay_log_info_repository FILE] [relay_log_purge ON] [relay_log_recovery OFF] [relay_log_space_limit 0] [report_host ] [report_password ] [report_port 3306] [report_user ] [require_secure_transport OFF] [rpl_stop_slave_timeout 31536000] [secure_auth ON] [secure_file_priv /var/lib/mysql-files/] [server_id 3379] [server_id_bits 32] [server_uuid 6b47d1d0-6bd7-11ee-8d59-0242c0a8ed04] [session_track_gtids OFF] [session_track_schema ON] [session_track_state_change OFF] [session_track_system_variables time_zone,autocommit,character_set_client,character_set_results,character_set_connection] [session_track_transaction_info OFF] [sha256_password_proxy_users OFF] [show_compatibility_56 OFF] [show_create_table_verbosity OFF] [show_old_temporals OFF] [skip_external_locking ON] [skip_name_resolve ON] [skip_networking OFF] [skip_show_database OFF] [slave_allow_batching OFF] [slave_checkpoint_group 512] [slave_checkpoint_period 300] [slave_compressed_protocol OFF] [slave_exec_mode STRICT] [slave_load_tmpdir /tmp] [slave_max_allowed_packet 1073741824] [slave_net_timeout 60] [slave_parallel_type DATABASE] [slave_parallel_workers 0] [slave_pending_jobs_size_max 16777216] [slave_preserve_commit_order OFF] [slave_rows_search_algorithms TABLE_SCAN,INDEX_SCAN] [slave_skip_errors OFF] [slave_sql_verify_checksum ON] [slave_transaction_retries 10] [slave_type_conversions ] [slow_launch_time 2] [slow_query_log OFF] [slow_query_log_file /var/lib/mysql/a8244020552d-slow.log] [socket /var/lib/mysql/mysql.sock] [sort_buffer_size 262144] [sql_auto_is_null OFF] [sql_big_selects ON] [sql_buffer_result OFF] [sql_log_off OFF] [sql_mode STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION] [sql_notes ON] [sql_quote_show_create ON] [sql_safe_updates OFF] [sql_select_limit 18446744073709551615] [sql_slave_skip_counter 0] [sql_warnings OFF] [ssl_ca ] [ssl_capath ] [ssl_cert ] [ssl_cipher ] [ssl_crl ] [ssl_crlpath ] [ssl_key ] [stored_program_cache 256] [super_read_only OFF] [sync_binlog 1] [sync_frm ON] [sync_master_info 10000] [sync_relay_log 10000] [sync_relay_log_info 10000] [system_time_zone UTC] [table_definition_cache 1400] [table_open_cache 2000] [table_open_cache_instances 16] [thread_cache_size 9] [thread_handling one-thread-per-connection] [thread_stack 196608] [time_format %H:%i:%s] [time_zone SYSTEM] [tls_version TLSv1,TLSv1.1] [tmp_table_size 16777216] [tmpdir /tmp] [transaction_alloc_block_size 8192] [transaction_isolation REPEATABLE-READ] [transaction_prealloc_size 4096] [transaction_read_only OFF] [transaction_write_set_extraction OFF] [tx_isolation REPEATABLE-READ] [tx_read_only OFF] [unique_checks ON] [updatable_views_with_limit YES] [version 5.7.26-1-log] [version_comment (Debian)] [version_compile_machine armv7l] [version_compile_os Linux] [wait_timeout 2880000]]

sql:SHOW GLOBAL VARIABLES WHERE Variable_name = 'autocommit';
mysqlRes:[[autocommit ON]]
gaeaRes:[[autocommit ON]]

sql:SHOW STATUS LIKE 'Up%';
mysqlRes:[[Uptime 433328] [Uptime_since_flush_status 433328]]
gaeaRes:[[Uptime 433328] [Uptime_since_flush_status 433328]]

sql:SHOW STATUS WHERE Variable_name;
mysqlRes:[]
gaeaRes:[]

sql:SHOW STATUS WHERE Variable_name LIKE 'Up%';
mysqlRes:[[Uptime 433328] [Uptime_since_flush_status 433328]]
gaeaRes:[[Uptime 433328] [Uptime_since_flush_status 433328]]

sql:SHOW FULL TABLES FROM sbtest1 LIKE 't%';
mysqlRes:[[t BASE TABLE] [t1 BASE TABLE] [t2 BASE TABLE] [t3 BASE TABLE] [test1 BASE TABLE] [test2 BASE TABLE] [test3 BASE TABLE] [test4 BASE TABLE] [test5 BASE TABLE] [test6 BASE TABLE] [test7 BASE TABLE] [test8 BASE TABLE] [test9 BASE TABLE]]
gaeaRes:[[t BASE TABLE] [t1 BASE TABLE] [t2 BASE TABLE] [t3 BASE TABLE] [test1 BASE TABLE] [test2 BASE TABLE] [test3 BASE TABLE] [test4 BASE TABLE] [test5 BASE TABLE] [test6 BASE TABLE] [test7 BASE TABLE] [test8 BASE TABLE] [test9 BASE TABLE]]

sql:SHOW FULL TABLES WHERE Table_type != 'VIEW';
mysqlRes:[[noshard_t1 BASE TABLE] [noshard_t2 BASE TABLE] [noshard_t3 BASE TABLE] [t BASE TABLE] [t1 BASE TABLE] [t2 BASE TABLE] [t3 BASE TABLE] [test1 BASE TABLE] [test2 BASE TABLE] [test3 BASE TABLE] [test4 BASE TABLE] [test5 BASE TABLE] [test6 BASE TABLE] [test7 BASE TABLE] [test8 BASE TABLE] [test9 BASE TABLE]]
gaeaRes:[[noshard_t1 BASE TABLE] [noshard_t2 BASE TABLE] [noshard_t3 BASE TABLE] [t BASE TABLE] [t1 BASE TABLE] [t2 BASE TABLE] [t3 BASE TABLE] [test1 BASE TABLE] [test2 BASE TABLE] [test3 BASE TABLE] [test4 BASE TABLE] [test5 BASE TABLE] [test6 BASE TABLE] [test7 BASE TABLE] [test8 BASE TABLE] [test9 BASE TABLE]]

sql:SHOW GRANTS;
mysqlRes:[[GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION]]
gaeaRes:[[GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION]]

sql:SHOW GRANTS FOR current_user();
mysqlRes:[[GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION]]
gaeaRes:[[GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION]]

sql:SHOW GRANTS FOR current_user;
mysqlRes:[[GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION]]
gaeaRes:[[GRANT ALL PRIVILEGES ON *.* TO 'superroot'@'%' WITH GRANT OPTION]]

sql:SHOW COLUMNS FROM t;
mysqlRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]
gaeaRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]

sql:SHOW COLUMNS FROM sbtest1.t;
mysqlRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]
gaeaRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]

sql:SHOW FIELDS FROM t;
mysqlRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]
gaeaRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]

sql:SHOW TRIGGERS LIKE 't';
mysqlRes:[]
gaeaRes:[]

sql:SHOW PROCEDURE STATUS WHERE Db='test';
mysqlRes:[]
gaeaRes:[]

sql:SHOW FUNCTION STATUS WHERE Db='test';
mysqlRes:[]
gaeaRes:[]

sql:SHOW INDEX FROM t;
mysqlRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]
gaeaRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]

sql:SHOW KEYS FROM t;
mysqlRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]
gaeaRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]

sql:SHOW INDEX IN t;
mysqlRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]
gaeaRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]

sql:SHOW KEYS IN t;
mysqlRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]
gaeaRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]

sql:SHOW INDEXES IN t where true;
mysqlRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]
gaeaRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]

sql:SHOW KEYS FROM t FROM sbtest1 where true;
mysqlRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]
gaeaRes:[[t 0 PRIMARY 1 id A 10 NULL NULL  BTREE  ] [t 1 idx1 1 col1 A 5 NULL NULL YES BTREE  ] [t 1 idx2 1 col2 A 5 NULL NULL YES BTREE  ]]

sql:SHOW EVENTS FROM t WHERE definer = 'current_user';
mysqlRes:[]
gaeaRes:[]

sql:SHOW PLUGINS;
mysqlRes:[[binlog ACTIVE STORAGE ENGINE NULL GPL] [mysql_native_password ACTIVE AUTHENTICATION NULL GPL] [sha256_password ACTIVE AUTHENTICATION NULL GPL] [CSV ACTIVE STORAGE ENGINE NULL GPL] [MEMORY ACTIVE STORAGE ENGINE NULL GPL] [InnoDB ACTIVE STORAGE ENGINE NULL GPL] [INNODB_TRX ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_LOCKS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_LOCK_WAITS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP_RESET ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMPMEM ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMPMEM_RESET ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP_PER_INDEX ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP_PER_INDEX_RESET ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_BUFFER_PAGE ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_BUFFER_PAGE_LRU ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_BUFFER_POOL_STATS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_TEMP_TABLE_INFO ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_METRICS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_DEFAULT_STOPWORD ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_DELETED ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_BEING_DELETED ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_CONFIG ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_INDEX_CACHE ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_INDEX_TABLE ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_TABLES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_TABLESTATS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_INDEXES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_COLUMNS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_FIELDS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_FOREIGN ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_FOREIGN_COLS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_TABLESPACES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_DATAFILES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_VIRTUAL ACTIVE INFORMATION SCHEMA NULL GPL] [MyISAM ACTIVE STORAGE ENGINE NULL GPL] [MRG_MYISAM ACTIVE STORAGE ENGINE NULL GPL] [PERFORMANCE_SCHEMA ACTIVE STORAGE ENGINE NULL GPL] [ARCHIVE ACTIVE STORAGE ENGINE NULL GPL] [BLACKHOLE ACTIVE STORAGE ENGINE NULL GPL] [FEDERATED DISABLED STORAGE ENGINE NULL GPL] [partition ACTIVE STORAGE ENGINE NULL GPL] [ngram ACTIVE FTPARSER NULL GPL]]
gaeaRes:[[binlog ACTIVE STORAGE ENGINE NULL GPL] [mysql_native_password ACTIVE AUTHENTICATION NULL GPL] [sha256_password ACTIVE AUTHENTICATION NULL GPL] [CSV ACTIVE STORAGE ENGINE NULL GPL] [MEMORY ACTIVE STORAGE ENGINE NULL GPL] [InnoDB ACTIVE STORAGE ENGINE NULL GPL] [INNODB_TRX ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_LOCKS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_LOCK_WAITS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP_RESET ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMPMEM ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMPMEM_RESET ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP_PER_INDEX ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_CMP_PER_INDEX_RESET ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_BUFFER_PAGE ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_BUFFER_PAGE_LRU ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_BUFFER_POOL_STATS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_TEMP_TABLE_INFO ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_METRICS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_DEFAULT_STOPWORD ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_DELETED ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_BEING_DELETED ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_CONFIG ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_INDEX_CACHE ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_FT_INDEX_TABLE ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_TABLES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_TABLESTATS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_INDEXES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_COLUMNS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_FIELDS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_FOREIGN ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_FOREIGN_COLS ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_TABLESPACES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_DATAFILES ACTIVE INFORMATION SCHEMA NULL GPL] [INNODB_SYS_VIRTUAL ACTIVE INFORMATION SCHEMA NULL GPL] [MyISAM ACTIVE STORAGE ENGINE NULL GPL] [MRG_MYISAM ACTIVE STORAGE ENGINE NULL GPL] [PERFORMANCE_SCHEMA ACTIVE STORAGE ENGINE NULL GPL] [ARCHIVE ACTIVE STORAGE ENGINE NULL GPL] [BLACKHOLE ACTIVE STORAGE ENGINE NULL GPL] [FEDERATED DISABLED STORAGE ENGINE NULL GPL] [partition ACTIVE STORAGE ENGINE NULL GPL] [ngram ACTIVE FTPARSER NULL GPL]]

sql:SHOW PROFILES;
mysqlRes:[]
gaeaRes:[]

sql:SHOW PRIVILEGES;
mysqlRes:[[Alter Tables To alter the table] [Alter routine Functions,Procedures To alter or drop stored functions/procedures] [Create Databases,Tables,Indexes To create new databases and tables] [Create routine Databases To use CREATE FUNCTION/PROCEDURE] [Create temporary tables Databases To use CREATE TEMPORARY TABLE] [Create view Tables To create new views] [Create user Server Admin To create new users] [Delete Tables To delete existing rows] [Drop Databases,Tables To drop databases, tables, and views] [Event Server Admin To create, alter, drop and execute events] [Execute Functions,Procedures To execute stored routines] [File File access on server To read and write files on the server] [Grant option Databases,Tables,Functions,Procedures To give to other users those privileges you possess] [Index Tables To create or drop indexes] [Insert Tables To insert data into tables] [Lock tables Databases To use LOCK TABLES (together with SELECT privilege)] [Process Server Admin To view the plain text of currently executing queries] [Proxy Server Admin To make proxy user possible] [References Databases,Tables To have references on tables] [Reload Server Admin To reload or refresh tables, logs and privileges] [Replication client Server Admin To ask where the slave or master servers are] [Replication slave Server Admin To read binary log events from the master] [Select Tables To retrieve rows from table] [Show databases Server Admin To see all databases with SHOW DATABASES] [Show view Tables To see views with SHOW CREATE VIEW] [Shutdown Server Admin To shut down the server] [Super Server Admin To use KILL thread, SET GLOBAL, CHANGE MASTER, etc.] [Trigger Tables To use triggers] [Create tablespace Server Admin To create/alter/drop tablespaces] [Update Tables To update existing rows] [Usage Server Admin No privileges - allow connect only]]
gaeaRes:[[Alter Tables To alter the table] [Alter routine Functions,Procedures To alter or drop stored functions/procedures] [Create Databases,Tables,Indexes To create new databases and tables] [Create routine Databases To use CREATE FUNCTION/PROCEDURE] [Create temporary tables Databases To use CREATE TEMPORARY TABLE] [Create view Tables To create new views] [Create user Server Admin To create new users] [Delete Tables To delete existing rows] [Drop Databases,Tables To drop databases, tables, and views] [Event Server Admin To create, alter, drop and execute events] [Execute Functions,Procedures To execute stored routines] [File File access on server To read and write files on the server] [Grant option Databases,Tables,Functions,Procedures To give to other users those privileges you possess] [Index Tables To create or drop indexes] [Insert Tables To insert data into tables] [Lock tables Databases To use LOCK TABLES (together with SELECT privilege)] [Process Server Admin To view the plain text of currently executing queries] [Proxy Server Admin To make proxy user possible] [References Databases,Tables To have references on tables] [Reload Server Admin To reload or refresh tables, logs and privileges] [Replication client Server Admin To ask where the slave or master servers are] [Replication slave Server Admin To read binary log events from the master] [Select Tables To retrieve rows from table] [Show databases Server Admin To see all databases with SHOW DATABASES] [Show view Tables To see views with SHOW CREATE VIEW] [Shutdown Server Admin To shut down the server] [Super Server Admin To use KILL thread, SET GLOBAL, CHANGE MASTER, etc.] [Trigger Tables To use triggers] [Create tablespace Server Admin To create/alter/drop tablespaces] [Update Tables To update existing rows] [Usage Server Admin No privileges - allow connect only]]

sql:show character set;
mysqlRes:[[big5 Big5 Traditional Chinese big5_chinese_ci 2] [dec8 DEC West European dec8_swedish_ci 1] [cp850 DOS West European cp850_general_ci 1] [hp8 HP West European hp8_english_ci 1] [koi8r KOI8-R Relcom Russian koi8r_general_ci 1] [latin1 cp1252 West European latin1_swedish_ci 1] [latin2 ISO 8859-2 Central European latin2_general_ci 1] [swe7 7bit Swedish swe7_swedish_ci 1] [ascii US ASCII ascii_general_ci 1] [ujis EUC-JP Japanese ujis_japanese_ci 3] [sjis Shift-JIS Japanese sjis_japanese_ci 2] [hebrew ISO 8859-8 Hebrew hebrew_general_ci 1] [tis620 TIS620 Thai tis620_thai_ci 1] [euckr EUC-KR Korean euckr_korean_ci 2] [koi8u KOI8-U Ukrainian koi8u_general_ci 1] [gb2312 GB2312 Simplified Chinese gb2312_chinese_ci 2] [greek ISO 8859-7 Greek greek_general_ci 1] [cp1250 Windows Central European cp1250_general_ci 1] [gbk GBK Simplified Chinese gbk_chinese_ci 2] [latin5 ISO 8859-9 Turkish latin5_turkish_ci 1] [armscii8 ARMSCII-8 Armenian armscii8_general_ci 1] [utf8 UTF-8 Unicode utf8_general_ci 3] [ucs2 UCS-2 Unicode ucs2_general_ci 2] [cp866 DOS Russian cp866_general_ci 1] [keybcs2 DOS Kamenicky Czech-Slovak keybcs2_general_ci 1] [macce Mac Central European macce_general_ci 1] [macroman Mac West European macroman_general_ci 1] [cp852 DOS Central European cp852_general_ci 1] [latin7 ISO 8859-13 Baltic latin7_general_ci 1] [utf8mb4 UTF-8 Unicode utf8mb4_general_ci 4] [cp1251 Windows Cyrillic cp1251_general_ci 1] [utf16 UTF-16 Unicode utf16_general_ci 4] [utf16le UTF-16LE Unicode utf16le_general_ci 4] [cp1256 Windows Arabic cp1256_general_ci 1] [cp1257 Windows Baltic cp1257_general_ci 1] [utf32 UTF-32 Unicode utf32_general_ci 4] [binary Binary pseudo charset binary 1] [geostd8 GEOSTD8 Georgian geostd8_general_ci 1] [cp932 SJIS for Windows Japanese cp932_japanese_ci 2] [eucjpms UJIS for Windows Japanese eucjpms_japanese_ci 3] [gb18030 China National Standard GB18030 gb18030_chinese_ci 4]]
gaeaRes:[[big5 Big5 Traditional Chinese big5_chinese_ci 2] [dec8 DEC West European dec8_swedish_ci 1] [cp850 DOS West European cp850_general_ci 1] [hp8 HP West European hp8_english_ci 1] [koi8r KOI8-R Relcom Russian koi8r_general_ci 1] [latin1 cp1252 West European latin1_swedish_ci 1] [latin2 ISO 8859-2 Central European latin2_general_ci 1] [swe7 7bit Swedish swe7_swedish_ci 1] [ascii US ASCII ascii_general_ci 1] [ujis EUC-JP Japanese ujis_japanese_ci 3] [sjis Shift-JIS Japanese sjis_japanese_ci 2] [hebrew ISO 8859-8 Hebrew hebrew_general_ci 1] [tis620 TIS620 Thai tis620_thai_ci 1] [euckr EUC-KR Korean euckr_korean_ci 2] [koi8u KOI8-U Ukrainian koi8u_general_ci 1] [gb2312 GB2312 Simplified Chinese gb2312_chinese_ci 2] [greek ISO 8859-7 Greek greek_general_ci 1] [cp1250 Windows Central European cp1250_general_ci 1] [gbk GBK Simplified Chinese gbk_chinese_ci 2] [latin5 ISO 8859-9 Turkish latin5_turkish_ci 1] [armscii8 ARMSCII-8 Armenian armscii8_general_ci 1] [utf8 UTF-8 Unicode utf8_general_ci 3] [ucs2 UCS-2 Unicode ucs2_general_ci 2] [cp866 DOS Russian cp866_general_ci 1] [keybcs2 DOS Kamenicky Czech-Slovak keybcs2_general_ci 1] [macce Mac Central European macce_general_ci 1] [macroman Mac West European macroman_general_ci 1] [cp852 DOS Central European cp852_general_ci 1] [latin7 ISO 8859-13 Baltic latin7_general_ci 1] [utf8mb4 UTF-8 Unicode utf8mb4_general_ci 4] [cp1251 Windows Cyrillic cp1251_general_ci 1] [utf16 UTF-16 Unicode utf16_general_ci 4] [utf16le UTF-16LE Unicode utf16le_general_ci 4] [cp1256 Windows Arabic cp1256_general_ci 1] [cp1257 Windows Baltic cp1257_general_ci 1] [utf32 UTF-32 Unicode utf32_general_ci 4] [binary Binary pseudo charset binary 1] [geostd8 GEOSTD8 Georgian geostd8_general_ci 1] [cp932 SJIS for Windows Japanese cp932_japanese_ci 2] [eucjpms UJIS for Windows Japanese eucjpms_japanese_ci 3] [gb18030 China National Standard GB18030 gb18030_chinese_ci 4]]

sql:show charset;
mysqlRes:[[big5 Big5 Traditional Chinese big5_chinese_ci 2] [dec8 DEC West European dec8_swedish_ci 1] [cp850 DOS West European cp850_general_ci 1] [hp8 HP West European hp8_english_ci 1] [koi8r KOI8-R Relcom Russian koi8r_general_ci 1] [latin1 cp1252 West European latin1_swedish_ci 1] [latin2 ISO 8859-2 Central European latin2_general_ci 1] [swe7 7bit Swedish swe7_swedish_ci 1] [ascii US ASCII ascii_general_ci 1] [ujis EUC-JP Japanese ujis_japanese_ci 3] [sjis Shift-JIS Japanese sjis_japanese_ci 2] [hebrew ISO 8859-8 Hebrew hebrew_general_ci 1] [tis620 TIS620 Thai tis620_thai_ci 1] [euckr EUC-KR Korean euckr_korean_ci 2] [koi8u KOI8-U Ukrainian koi8u_general_ci 1] [gb2312 GB2312 Simplified Chinese gb2312_chinese_ci 2] [greek ISO 8859-7 Greek greek_general_ci 1] [cp1250 Windows Central European cp1250_general_ci 1] [gbk GBK Simplified Chinese gbk_chinese_ci 2] [latin5 ISO 8859-9 Turkish latin5_turkish_ci 1] [armscii8 ARMSCII-8 Armenian armscii8_general_ci 1] [utf8 UTF-8 Unicode utf8_general_ci 3] [ucs2 UCS-2 Unicode ucs2_general_ci 2] [cp866 DOS Russian cp866_general_ci 1] [keybcs2 DOS Kamenicky Czech-Slovak keybcs2_general_ci 1] [macce Mac Central European macce_general_ci 1] [macroman Mac West European macroman_general_ci 1] [cp852 DOS Central European cp852_general_ci 1] [latin7 ISO 8859-13 Baltic latin7_general_ci 1] [utf8mb4 UTF-8 Unicode utf8mb4_general_ci 4] [cp1251 Windows Cyrillic cp1251_general_ci 1] [utf16 UTF-16 Unicode utf16_general_ci 4] [utf16le UTF-16LE Unicode utf16le_general_ci 4] [cp1256 Windows Arabic cp1256_general_ci 1] [cp1257 Windows Baltic cp1257_general_ci 1] [utf32 UTF-32 Unicode utf32_general_ci 4] [binary Binary pseudo charset binary 1] [geostd8 GEOSTD8 Georgian geostd8_general_ci 1] [cp932 SJIS for Windows Japanese cp932_japanese_ci 2] [eucjpms UJIS for Windows Japanese eucjpms_japanese_ci 3] [gb18030 China National Standard GB18030 gb18030_chinese_ci 4]]
gaeaRes:[[big5 Big5 Traditional Chinese big5_chinese_ci 2] [dec8 DEC West European dec8_swedish_ci 1] [cp850 DOS West European cp850_general_ci 1] [hp8 HP West European hp8_english_ci 1] [koi8r KOI8-R Relcom Russian koi8r_general_ci 1] [latin1 cp1252 West European latin1_swedish_ci 1] [latin2 ISO 8859-2 Central European latin2_general_ci 1] [swe7 7bit Swedish swe7_swedish_ci 1] [ascii US ASCII ascii_general_ci 1] [ujis EUC-JP Japanese ujis_japanese_ci 3] [sjis Shift-JIS Japanese sjis_japanese_ci 2] [hebrew ISO 8859-8 Hebrew hebrew_general_ci 1] [tis620 TIS620 Thai tis620_thai_ci 1] [euckr EUC-KR Korean euckr_korean_ci 2] [koi8u KOI8-U Ukrainian koi8u_general_ci 1] [gb2312 GB2312 Simplified Chinese gb2312_chinese_ci 2] [greek ISO 8859-7 Greek greek_general_ci 1] [cp1250 Windows Central European cp1250_general_ci 1] [gbk GBK Simplified Chinese gbk_chinese_ci 2] [latin5 ISO 8859-9 Turkish latin5_turkish_ci 1] [armscii8 ARMSCII-8 Armenian armscii8_general_ci 1] [utf8 UTF-8 Unicode utf8_general_ci 3] [ucs2 UCS-2 Unicode ucs2_general_ci 2] [cp866 DOS Russian cp866_general_ci 1] [keybcs2 DOS Kamenicky Czech-Slovak keybcs2_general_ci 1] [macce Mac Central European macce_general_ci 1] [macroman Mac West European macroman_general_ci 1] [cp852 DOS Central European cp852_general_ci 1] [latin7 ISO 8859-13 Baltic latin7_general_ci 1] [utf8mb4 UTF-8 Unicode utf8mb4_general_ci 4] [cp1251 Windows Cyrillic cp1251_general_ci 1] [utf16 UTF-16 Unicode utf16_general_ci 4] [utf16le UTF-16LE Unicode utf16le_general_ci 4] [cp1256 Windows Arabic cp1256_general_ci 1] [cp1257 Windows Baltic cp1257_general_ci 1] [utf32 UTF-32 Unicode utf32_general_ci 4] [binary Binary pseudo charset binary 1] [geostd8 GEOSTD8 Georgian geostd8_general_ci 1] [cp932 SJIS for Windows Japanese cp932_japanese_ci 2] [eucjpms UJIS for Windows Japanese eucjpms_japanese_ci 3] [gb18030 China National Standard GB18030 gb18030_chinese_ci 4]]

sql:show collation;
mysqlRes:[[big5_chinese_ci big5 1 Yes Yes 1] [big5_bin big5 84  Yes 1] [dec8_swedish_ci dec8 3 Yes Yes 1] [dec8_bin dec8 69  Yes 1] [cp850_general_ci cp850 4 Yes Yes 1] [cp850_bin cp850 80  Yes 1] [hp8_english_ci hp8 6 Yes Yes 1] [hp8_bin hp8 72  Yes 1] [koi8r_general_ci koi8r 7 Yes Yes 1] [koi8r_bin koi8r 74  Yes 1] [latin1_german1_ci latin1 5  Yes 1] [latin1_swedish_ci latin1 8 Yes Yes 1] [latin1_danish_ci latin1 15  Yes 1] [latin1_german2_ci latin1 31  Yes 2] [latin1_bin latin1 47  Yes 1] [latin1_general_ci latin1 48  Yes 1] [latin1_general_cs latin1 49  Yes 1] [latin1_spanish_ci latin1 94  Yes 1] [latin2_czech_cs latin2 2  Yes 4] [latin2_general_ci latin2 9 Yes Yes 1] [latin2_hungarian_ci latin2 21  Yes 1] [latin2_croatian_ci latin2 27  Yes 1] [latin2_bin latin2 77  Yes 1] [swe7_swedish_ci swe7 10 Yes Yes 1] [swe7_bin swe7 82  Yes 1] [ascii_general_ci ascii 11 Yes Yes 1] [ascii_bin ascii 65  Yes 1] [ujis_japanese_ci ujis 12 Yes Yes 1] [ujis_bin ujis 91  Yes 1] [sjis_japanese_ci sjis 13 Yes Yes 1] [sjis_bin sjis 88  Yes 1] [hebrew_general_ci hebrew 16 Yes Yes 1] [hebrew_bin hebrew 71  Yes 1] [tis620_thai_ci tis620 18 Yes Yes 4] [tis620_bin tis620 89  Yes 1] [euckr_korean_ci euckr 19 Yes Yes 1] [euckr_bin euckr 85  Yes 1] [koi8u_general_ci koi8u 22 Yes Yes 1] [koi8u_bin koi8u 75  Yes 1] [gb2312_chinese_ci gb2312 24 Yes Yes 1] [gb2312_bin gb2312 86  Yes 1] [greek_general_ci greek 25 Yes Yes 1] [greek_bin greek 70  Yes 1] [cp1250_general_ci cp1250 26 Yes Yes 1] [cp1250_czech_cs cp1250 34  Yes 2] [cp1250_croatian_ci cp1250 44  Yes 1] [cp1250_bin cp1250 66  Yes 1] [cp1250_polish_ci cp1250 99  Yes 1] [gbk_chinese_ci gbk 28 Yes Yes 1] [gbk_bin gbk 87  Yes 1] [latin5_turkish_ci latin5 30 Yes Yes 1] [latin5_bin latin5 78  Yes 1] [armscii8_general_ci armscii8 32 Yes Yes 1] [armscii8_bin armscii8 64  Yes 1] [utf8_general_ci utf8 33 Yes Yes 1] [utf8_bin utf8 83  Yes 1] [utf8_unicode_ci utf8 192  Yes 8] [utf8_icelandic_ci utf8 193  Yes 8] [utf8_latvian_ci utf8 194  Yes 8] [utf8_romanian_ci utf8 195  Yes 8] [utf8_slovenian_ci utf8 196  Yes 8] [utf8_polish_ci utf8 197  Yes 8] [utf8_estonian_ci utf8 198  Yes 8] [utf8_spanish_ci utf8 199  Yes 8] [utf8_swedish_ci utf8 200  Yes 8] [utf8_turkish_ci utf8 201  Yes 8] [utf8_czech_ci utf8 202  Yes 8] [utf8_danish_ci utf8 203  Yes 8] [utf8_lithuanian_ci utf8 204  Yes 8] [utf8_slovak_ci utf8 205  Yes 8] [utf8_spanish2_ci utf8 206  Yes 8] [utf8_roman_ci utf8 207  Yes 8] [utf8_persian_ci utf8 208  Yes 8] [utf8_esperanto_ci utf8 209  Yes 8] [utf8_hungarian_ci utf8 210  Yes 8] [utf8_sinhala_ci utf8 211  Yes 8] [utf8_german2_ci utf8 212  Yes 8] [utf8_croatian_ci utf8 213  Yes 8] [utf8_unicode_520_ci utf8 214  Yes 8] [utf8_vietnamese_ci utf8 215  Yes 8] [utf8_general_mysql500_ci utf8 223  Yes 1] [ucs2_general_ci ucs2 35 Yes Yes 1] [ucs2_bin ucs2 90  Yes 1] [ucs2_unicode_ci ucs2 128  Yes 8] [ucs2_icelandic_ci ucs2 129  Yes 8] [ucs2_latvian_ci ucs2 130  Yes 8] [ucs2_romanian_ci ucs2 131  Yes 8] [ucs2_slovenian_ci ucs2 132  Yes 8] [ucs2_polish_ci ucs2 133  Yes 8] [ucs2_estonian_ci ucs2 134  Yes 8] [ucs2_spanish_ci ucs2 135  Yes 8] [ucs2_swedish_ci ucs2 136  Yes 8] [ucs2_turkish_ci ucs2 137  Yes 8] [ucs2_czech_ci ucs2 138  Yes 8] [ucs2_danish_ci ucs2 139  Yes 8] [ucs2_lithuanian_ci ucs2 140  Yes 8] [ucs2_slovak_ci ucs2 141  Yes 8] [ucs2_spanish2_ci ucs2 142  Yes 8] [ucs2_roman_ci ucs2 143  Yes 8] [ucs2_persian_ci ucs2 144  Yes 8] [ucs2_esperanto_ci ucs2 145  Yes 8] [ucs2_hungarian_ci ucs2 146  Yes 8] [ucs2_sinhala_ci ucs2 147  Yes 8] [ucs2_german2_ci ucs2 148  Yes 8] [ucs2_croatian_ci ucs2 149  Yes 8] [ucs2_unicode_520_ci ucs2 150  Yes 8] [ucs2_vietnamese_ci ucs2 151  Yes 8] [ucs2_general_mysql500_ci ucs2 159  Yes 1] [cp866_general_ci cp866 36 Yes Yes 1] [cp866_bin cp866 68  Yes 1] [keybcs2_general_ci keybcs2 37 Yes Yes 1] [keybcs2_bin keybcs2 73  Yes 1] [macce_general_ci macce 38 Yes Yes 1] [macce_bin macce 43  Yes 1] [macroman_general_ci macroman 39 Yes Yes 1] [macroman_bin macroman 53  Yes 1] [cp852_general_ci cp852 40 Yes Yes 1] [cp852_bin cp852 81  Yes 1] [latin7_estonian_cs latin7 20  Yes 1] [latin7_general_ci latin7 41 Yes Yes 1] [latin7_general_cs latin7 42  Yes 1] [latin7_bin latin7 79  Yes 1] [utf8mb4_general_ci utf8mb4 45 Yes Yes 1] [utf8mb4_bin utf8mb4 46  Yes 1] [utf8mb4_unicode_ci utf8mb4 224  Yes 8] [utf8mb4_icelandic_ci utf8mb4 225  Yes 8] [utf8mb4_latvian_ci utf8mb4 226  Yes 8] [utf8mb4_romanian_ci utf8mb4 227  Yes 8] [utf8mb4_slovenian_ci utf8mb4 228  Yes 8] [utf8mb4_polish_ci utf8mb4 229  Yes 8] [utf8mb4_estonian_ci utf8mb4 230  Yes 8] [utf8mb4_spanish_ci utf8mb4 231  Yes 8] [utf8mb4_swedish_ci utf8mb4 232  Yes 8] [utf8mb4_turkish_ci utf8mb4 233  Yes 8] [utf8mb4_czech_ci utf8mb4 234  Yes 8] [utf8mb4_danish_ci utf8mb4 235  Yes 8] [utf8mb4_lithuanian_ci utf8mb4 236  Yes 8] [utf8mb4_slovak_ci utf8mb4 237  Yes 8] [utf8mb4_spanish2_ci utf8mb4 238  Yes 8] [utf8mb4_roman_ci utf8mb4 239  Yes 8] [utf8mb4_persian_ci utf8mb4 240  Yes 8] [utf8mb4_esperanto_ci utf8mb4 241  Yes 8] [utf8mb4_hungarian_ci utf8mb4 242  Yes 8] [utf8mb4_sinhala_ci utf8mb4 243  Yes 8] [utf8mb4_german2_ci utf8mb4 244  Yes 8] [utf8mb4_croatian_ci utf8mb4 245  Yes 8] [utf8mb4_unicode_520_ci utf8mb4 246  Yes 8] [utf8mb4_vietnamese_ci utf8mb4 247  Yes 8] [cp1251_bulgarian_ci cp1251 14  Yes 1] [cp1251_ukrainian_ci cp1251 23  Yes 1] [cp1251_bin cp1251 50  Yes 1] [cp1251_general_ci cp1251 51 Yes Yes 1] [cp1251_general_cs cp1251 52  Yes 1] [utf16_general_ci utf16 54 Yes Yes 1] [utf16_bin utf16 55  Yes 1] [utf16_unicode_ci utf16 101  Yes 8] [utf16_icelandic_ci utf16 102  Yes 8] [utf16_latvian_ci utf16 103  Yes 8] [utf16_romanian_ci utf16 104  Yes 8] [utf16_slovenian_ci utf16 105  Yes 8] [utf16_polish_ci utf16 106  Yes 8] [utf16_estonian_ci utf16 107  Yes 8] [utf16_spanish_ci utf16 108  Yes 8] [utf16_swedish_ci utf16 109  Yes 8] [utf16_turkish_ci utf16 110  Yes 8] [utf16_czech_ci utf16 111  Yes 8] [utf16_danish_ci utf16 112  Yes 8] [utf16_lithuanian_ci utf16 113  Yes 8] [utf16_slovak_ci utf16 114  Yes 8] [utf16_spanish2_ci utf16 115  Yes 8] [utf16_roman_ci utf16 116  Yes 8] [utf16_persian_ci utf16 117  Yes 8] [utf16_esperanto_ci utf16 118  Yes 8] [utf16_hungarian_ci utf16 119  Yes 8] [utf16_sinhala_ci utf16 120  Yes 8] [utf16_german2_ci utf16 121  Yes 8] [utf16_croatian_ci utf16 122  Yes 8] [utf16_unicode_520_ci utf16 123  Yes 8] [utf16_vietnamese_ci utf16 124  Yes 8] [utf16le_general_ci utf16le 56 Yes Yes 1] [utf16le_bin utf16le 62  Yes 1] [cp1256_general_ci cp1256 57 Yes Yes 1] [cp1256_bin cp1256 67  Yes 1] [cp1257_lithuanian_ci cp1257 29  Yes 1] [cp1257_bin cp1257 58  Yes 1] [cp1257_general_ci cp1257 59 Yes Yes 1] [utf32_general_ci utf32 60 Yes Yes 1] [utf32_bin utf32 61  Yes 1] [utf32_unicode_ci utf32 160  Yes 8] [utf32_icelandic_ci utf32 161  Yes 8] [utf32_latvian_ci utf32 162  Yes 8] [utf32_romanian_ci utf32 163  Yes 8] [utf32_slovenian_ci utf32 164  Yes 8] [utf32_polish_ci utf32 165  Yes 8] [utf32_estonian_ci utf32 166  Yes 8] [utf32_spanish_ci utf32 167  Yes 8] [utf32_swedish_ci utf32 168  Yes 8] [utf32_turkish_ci utf32 169  Yes 8] [utf32_czech_ci utf32 170  Yes 8] [utf32_danish_ci utf32 171  Yes 8] [utf32_lithuanian_ci utf32 172  Yes 8] [utf32_slovak_ci utf32 173  Yes 8] [utf32_spanish2_ci utf32 174  Yes 8] [utf32_roman_ci utf32 175  Yes 8] [utf32_persian_ci utf32 176  Yes 8] [utf32_esperanto_ci utf32 177  Yes 8] [utf32_hungarian_ci utf32 178  Yes 8] [utf32_sinhala_ci utf32 179  Yes 8] [utf32_german2_ci utf32 180  Yes 8] [utf32_croatian_ci utf32 181  Yes 8] [utf32_unicode_520_ci utf32 182  Yes 8] [utf32_vietnamese_ci utf32 183  Yes 8] [binary binary 63 Yes Yes 1] [geostd8_general_ci geostd8 92 Yes Yes 1] [geostd8_bin geostd8 93  Yes 1] [cp932_japanese_ci cp932 95 Yes Yes 1] [cp932_bin cp932 96  Yes 1] [eucjpms_japanese_ci eucjpms 97 Yes Yes 1] [eucjpms_bin eucjpms 98  Yes 1] [gb18030_chinese_ci gb18030 248 Yes Yes 2] [gb18030_bin gb18030 249  Yes 1] [gb18030_unicode_520_ci gb18030 250  Yes 8]]
gaeaRes:[[big5_chinese_ci big5 1 Yes Yes 1] [big5_bin big5 84  Yes 1] [dec8_swedish_ci dec8 3 Yes Yes 1] [dec8_bin dec8 69  Yes 1] [cp850_general_ci cp850 4 Yes Yes 1] [cp850_bin cp850 80  Yes 1] [hp8_english_ci hp8 6 Yes Yes 1] [hp8_bin hp8 72  Yes 1] [koi8r_general_ci koi8r 7 Yes Yes 1] [koi8r_bin koi8r 74  Yes 1] [latin1_german1_ci latin1 5  Yes 1] [latin1_swedish_ci latin1 8 Yes Yes 1] [latin1_danish_ci latin1 15  Yes 1] [latin1_german2_ci latin1 31  Yes 2] [latin1_bin latin1 47  Yes 1] [latin1_general_ci latin1 48  Yes 1] [latin1_general_cs latin1 49  Yes 1] [latin1_spanish_ci latin1 94  Yes 1] [latin2_czech_cs latin2 2  Yes 4] [latin2_general_ci latin2 9 Yes Yes 1] [latin2_hungarian_ci latin2 21  Yes 1] [latin2_croatian_ci latin2 27  Yes 1] [latin2_bin latin2 77  Yes 1] [swe7_swedish_ci swe7 10 Yes Yes 1] [swe7_bin swe7 82  Yes 1] [ascii_general_ci ascii 11 Yes Yes 1] [ascii_bin ascii 65  Yes 1] [ujis_japanese_ci ujis 12 Yes Yes 1] [ujis_bin ujis 91  Yes 1] [sjis_japanese_ci sjis 13 Yes Yes 1] [sjis_bin sjis 88  Yes 1] [hebrew_general_ci hebrew 16 Yes Yes 1] [hebrew_bin hebrew 71  Yes 1] [tis620_thai_ci tis620 18 Yes Yes 4] [tis620_bin tis620 89  Yes 1] [euckr_korean_ci euckr 19 Yes Yes 1] [euckr_bin euckr 85  Yes 1] [koi8u_general_ci koi8u 22 Yes Yes 1] [koi8u_bin koi8u 75  Yes 1] [gb2312_chinese_ci gb2312 24 Yes Yes 1] [gb2312_bin gb2312 86  Yes 1] [greek_general_ci greek 25 Yes Yes 1] [greek_bin greek 70  Yes 1] [cp1250_general_ci cp1250 26 Yes Yes 1] [cp1250_czech_cs cp1250 34  Yes 2] [cp1250_croatian_ci cp1250 44  Yes 1] [cp1250_bin cp1250 66  Yes 1] [cp1250_polish_ci cp1250 99  Yes 1] [gbk_chinese_ci gbk 28 Yes Yes 1] [gbk_bin gbk 87  Yes 1] [latin5_turkish_ci latin5 30 Yes Yes 1] [latin5_bin latin5 78  Yes 1] [armscii8_general_ci armscii8 32 Yes Yes 1] [armscii8_bin armscii8 64  Yes 1] [utf8_general_ci utf8 33 Yes Yes 1] [utf8_bin utf8 83  Yes 1] [utf8_unicode_ci utf8 192  Yes 8] [utf8_icelandic_ci utf8 193  Yes 8] [utf8_latvian_ci utf8 194  Yes 8] [utf8_romanian_ci utf8 195  Yes 8] [utf8_slovenian_ci utf8 196  Yes 8] [utf8_polish_ci utf8 197  Yes 8] [utf8_estonian_ci utf8 198  Yes 8] [utf8_spanish_ci utf8 199  Yes 8] [utf8_swedish_ci utf8 200  Yes 8] [utf8_turkish_ci utf8 201  Yes 8] [utf8_czech_ci utf8 202  Yes 8] [utf8_danish_ci utf8 203  Yes 8] [utf8_lithuanian_ci utf8 204  Yes 8] [utf8_slovak_ci utf8 205  Yes 8] [utf8_spanish2_ci utf8 206  Yes 8] [utf8_roman_ci utf8 207  Yes 8] [utf8_persian_ci utf8 208  Yes 8] [utf8_esperanto_ci utf8 209  Yes 8] [utf8_hungarian_ci utf8 210  Yes 8] [utf8_sinhala_ci utf8 211  Yes 8] [utf8_german2_ci utf8 212  Yes 8] [utf8_croatian_ci utf8 213  Yes 8] [utf8_unicode_520_ci utf8 214  Yes 8] [utf8_vietnamese_ci utf8 215  Yes 8] [utf8_general_mysql500_ci utf8 223  Yes 1] [ucs2_general_ci ucs2 35 Yes Yes 1] [ucs2_bin ucs2 90  Yes 1] [ucs2_unicode_ci ucs2 128  Yes 8] [ucs2_icelandic_ci ucs2 129  Yes 8] [ucs2_latvian_ci ucs2 130  Yes 8] [ucs2_romanian_ci ucs2 131  Yes 8] [ucs2_slovenian_ci ucs2 132  Yes 8] [ucs2_polish_ci ucs2 133  Yes 8] [ucs2_estonian_ci ucs2 134  Yes 8] [ucs2_spanish_ci ucs2 135  Yes 8] [ucs2_swedish_ci ucs2 136  Yes 8] [ucs2_turkish_ci ucs2 137  Yes 8] [ucs2_czech_ci ucs2 138  Yes 8] [ucs2_danish_ci ucs2 139  Yes 8] [ucs2_lithuanian_ci ucs2 140  Yes 8] [ucs2_slovak_ci ucs2 141  Yes 8] [ucs2_spanish2_ci ucs2 142  Yes 8] [ucs2_roman_ci ucs2 143  Yes 8] [ucs2_persian_ci ucs2 144  Yes 8] [ucs2_esperanto_ci ucs2 145  Yes 8] [ucs2_hungarian_ci ucs2 146  Yes 8] [ucs2_sinhala_ci ucs2 147  Yes 8] [ucs2_german2_ci ucs2 148  Yes 8] [ucs2_croatian_ci ucs2 149  Yes 8] [ucs2_unicode_520_ci ucs2 150  Yes 8] [ucs2_vietnamese_ci ucs2 151  Yes 8] [ucs2_general_mysql500_ci ucs2 159  Yes 1] [cp866_general_ci cp866 36 Yes Yes 1] [cp866_bin cp866 68  Yes 1] [keybcs2_general_ci keybcs2 37 Yes Yes 1] [keybcs2_bin keybcs2 73  Yes 1] [macce_general_ci macce 38 Yes Yes 1] [macce_bin macce 43  Yes 1] [macroman_general_ci macroman 39 Yes Yes 1] [macroman_bin macroman 53  Yes 1] [cp852_general_ci cp852 40 Yes Yes 1] [cp852_bin cp852 81  Yes 1] [latin7_estonian_cs latin7 20  Yes 1] [latin7_general_ci latin7 41 Yes Yes 1] [latin7_general_cs latin7 42  Yes 1] [latin7_bin latin7 79  Yes 1] [utf8mb4_general_ci utf8mb4 45 Yes Yes 1] [utf8mb4_bin utf8mb4 46  Yes 1] [utf8mb4_unicode_ci utf8mb4 224  Yes 8] [utf8mb4_icelandic_ci utf8mb4 225  Yes 8] [utf8mb4_latvian_ci utf8mb4 226  Yes 8] [utf8mb4_romanian_ci utf8mb4 227  Yes 8] [utf8mb4_slovenian_ci utf8mb4 228  Yes 8] [utf8mb4_polish_ci utf8mb4 229  Yes 8] [utf8mb4_estonian_ci utf8mb4 230  Yes 8] [utf8mb4_spanish_ci utf8mb4 231  Yes 8] [utf8mb4_swedish_ci utf8mb4 232  Yes 8] [utf8mb4_turkish_ci utf8mb4 233  Yes 8] [utf8mb4_czech_ci utf8mb4 234  Yes 8] [utf8mb4_danish_ci utf8mb4 235  Yes 8] [utf8mb4_lithuanian_ci utf8mb4 236  Yes 8] [utf8mb4_slovak_ci utf8mb4 237  Yes 8] [utf8mb4_spanish2_ci utf8mb4 238  Yes 8] [utf8mb4_roman_ci utf8mb4 239  Yes 8] [utf8mb4_persian_ci utf8mb4 240  Yes 8] [utf8mb4_esperanto_ci utf8mb4 241  Yes 8] [utf8mb4_hungarian_ci utf8mb4 242  Yes 8] [utf8mb4_sinhala_ci utf8mb4 243  Yes 8] [utf8mb4_german2_ci utf8mb4 244  Yes 8] [utf8mb4_croatian_ci utf8mb4 245  Yes 8] [utf8mb4_unicode_520_ci utf8mb4 246  Yes 8] [utf8mb4_vietnamese_ci utf8mb4 247  Yes 8] [cp1251_bulgarian_ci cp1251 14  Yes 1] [cp1251_ukrainian_ci cp1251 23  Yes 1] [cp1251_bin cp1251 50  Yes 1] [cp1251_general_ci cp1251 51 Yes Yes 1] [cp1251_general_cs cp1251 52  Yes 1] [utf16_general_ci utf16 54 Yes Yes 1] [utf16_bin utf16 55  Yes 1] [utf16_unicode_ci utf16 101  Yes 8] [utf16_icelandic_ci utf16 102  Yes 8] [utf16_latvian_ci utf16 103  Yes 8] [utf16_romanian_ci utf16 104  Yes 8] [utf16_slovenian_ci utf16 105  Yes 8] [utf16_polish_ci utf16 106  Yes 8] [utf16_estonian_ci utf16 107  Yes 8] [utf16_spanish_ci utf16 108  Yes 8] [utf16_swedish_ci utf16 109  Yes 8] [utf16_turkish_ci utf16 110  Yes 8] [utf16_czech_ci utf16 111  Yes 8] [utf16_danish_ci utf16 112  Yes 8] [utf16_lithuanian_ci utf16 113  Yes 8] [utf16_slovak_ci utf16 114  Yes 8] [utf16_spanish2_ci utf16 115  Yes 8] [utf16_roman_ci utf16 116  Yes 8] [utf16_persian_ci utf16 117  Yes 8] [utf16_esperanto_ci utf16 118  Yes 8] [utf16_hungarian_ci utf16 119  Yes 8] [utf16_sinhala_ci utf16 120  Yes 8] [utf16_german2_ci utf16 121  Yes 8] [utf16_croatian_ci utf16 122  Yes 8] [utf16_unicode_520_ci utf16 123  Yes 8] [utf16_vietnamese_ci utf16 124  Yes 8] [utf16le_general_ci utf16le 56 Yes Yes 1] [utf16le_bin utf16le 62  Yes 1] [cp1256_general_ci cp1256 57 Yes Yes 1] [cp1256_bin cp1256 67  Yes 1] [cp1257_lithuanian_ci cp1257 29  Yes 1] [cp1257_bin cp1257 58  Yes 1] [cp1257_general_ci cp1257 59 Yes Yes 1] [utf32_general_ci utf32 60 Yes Yes 1] [utf32_bin utf32 61  Yes 1] [utf32_unicode_ci utf32 160  Yes 8] [utf32_icelandic_ci utf32 161  Yes 8] [utf32_latvian_ci utf32 162  Yes 8] [utf32_romanian_ci utf32 163  Yes 8] [utf32_slovenian_ci utf32 164  Yes 8] [utf32_polish_ci utf32 165  Yes 8] [utf32_estonian_ci utf32 166  Yes 8] [utf32_spanish_ci utf32 167  Yes 8] [utf32_swedish_ci utf32 168  Yes 8] [utf32_turkish_ci utf32 169  Yes 8] [utf32_czech_ci utf32 170  Yes 8] [utf32_danish_ci utf32 171  Yes 8] [utf32_lithuanian_ci utf32 172  Yes 8] [utf32_slovak_ci utf32 173  Yes 8] [utf32_spanish2_ci utf32 174  Yes 8] [utf32_roman_ci utf32 175  Yes 8] [utf32_persian_ci utf32 176  Yes 8] [utf32_esperanto_ci utf32 177  Yes 8] [utf32_hungarian_ci utf32 178  Yes 8] [utf32_sinhala_ci utf32 179  Yes 8] [utf32_german2_ci utf32 180  Yes 8] [utf32_croatian_ci utf32 181  Yes 8] [utf32_unicode_520_ci utf32 182  Yes 8] [utf32_vietnamese_ci utf32 183  Yes 8] [binary binary 63 Yes Yes 1] [geostd8_general_ci geostd8 92 Yes Yes 1] [geostd8_bin geostd8 93  Yes 1] [cp932_japanese_ci cp932 95 Yes Yes 1] [cp932_bin cp932 96  Yes 1] [eucjpms_japanese_ci eucjpms 97 Yes Yes 1] [eucjpms_bin eucjpms 98  Yes 1] [gb18030_chinese_ci gb18030 248 Yes Yes 2] [gb18030_bin gb18030 249  Yes 1] [gb18030_unicode_520_ci gb18030 250  Yes 8]]

sql:show collation like 'utf8%';
mysqlRes:[[utf8_general_ci utf8 33 Yes Yes 1] [utf8_bin utf8 83  Yes 1] [utf8_unicode_ci utf8 192  Yes 8] [utf8_icelandic_ci utf8 193  Yes 8] [utf8_latvian_ci utf8 194  Yes 8] [utf8_romanian_ci utf8 195  Yes 8] [utf8_slovenian_ci utf8 196  Yes 8] [utf8_polish_ci utf8 197  Yes 8] [utf8_estonian_ci utf8 198  Yes 8] [utf8_spanish_ci utf8 199  Yes 8] [utf8_swedish_ci utf8 200  Yes 8] [utf8_turkish_ci utf8 201  Yes 8] [utf8_czech_ci utf8 202  Yes 8] [utf8_danish_ci utf8 203  Yes 8] [utf8_lithuanian_ci utf8 204  Yes 8] [utf8_slovak_ci utf8 205  Yes 8] [utf8_spanish2_ci utf8 206  Yes 8] [utf8_roman_ci utf8 207  Yes 8] [utf8_persian_ci utf8 208  Yes 8] [utf8_esperanto_ci utf8 209  Yes 8] [utf8_hungarian_ci utf8 210  Yes 8] [utf8_sinhala_ci utf8 211  Yes 8] [utf8_german2_ci utf8 212  Yes 8] [utf8_croatian_ci utf8 213  Yes 8] [utf8_unicode_520_ci utf8 214  Yes 8] [utf8_vietnamese_ci utf8 215  Yes 8] [utf8_general_mysql500_ci utf8 223  Yes 1] [utf8mb4_general_ci utf8mb4 45 Yes Yes 1] [utf8mb4_bin utf8mb4 46  Yes 1] [utf8mb4_unicode_ci utf8mb4 224  Yes 8] [utf8mb4_icelandic_ci utf8mb4 225  Yes 8] [utf8mb4_latvian_ci utf8mb4 226  Yes 8] [utf8mb4_romanian_ci utf8mb4 227  Yes 8] [utf8mb4_slovenian_ci utf8mb4 228  Yes 8] [utf8mb4_polish_ci utf8mb4 229  Yes 8] [utf8mb4_estonian_ci utf8mb4 230  Yes 8] [utf8mb4_spanish_ci utf8mb4 231  Yes 8] [utf8mb4_swedish_ci utf8mb4 232  Yes 8] [utf8mb4_turkish_ci utf8mb4 233  Yes 8] [utf8mb4_czech_ci utf8mb4 234  Yes 8] [utf8mb4_danish_ci utf8mb4 235  Yes 8] [utf8mb4_lithuanian_ci utf8mb4 236  Yes 8] [utf8mb4_slovak_ci utf8mb4 237  Yes 8] [utf8mb4_spanish2_ci utf8mb4 238  Yes 8] [utf8mb4_roman_ci utf8mb4 239  Yes 8] [utf8mb4_persian_ci utf8mb4 240  Yes 8] [utf8mb4_esperanto_ci utf8mb4 241  Yes 8] [utf8mb4_hungarian_ci utf8mb4 242  Yes 8] [utf8mb4_sinhala_ci utf8mb4 243  Yes 8] [utf8mb4_german2_ci utf8mb4 244  Yes 8] [utf8mb4_croatian_ci utf8mb4 245  Yes 8] [utf8mb4_unicode_520_ci utf8mb4 246  Yes 8] [utf8mb4_vietnamese_ci utf8mb4 247  Yes 8]]
gaeaRes:[[utf8_general_ci utf8 33 Yes Yes 1] [utf8_bin utf8 83  Yes 1] [utf8_unicode_ci utf8 192  Yes 8] [utf8_icelandic_ci utf8 193  Yes 8] [utf8_latvian_ci utf8 194  Yes 8] [utf8_romanian_ci utf8 195  Yes 8] [utf8_slovenian_ci utf8 196  Yes 8] [utf8_polish_ci utf8 197  Yes 8] [utf8_estonian_ci utf8 198  Yes 8] [utf8_spanish_ci utf8 199  Yes 8] [utf8_swedish_ci utf8 200  Yes 8] [utf8_turkish_ci utf8 201  Yes 8] [utf8_czech_ci utf8 202  Yes 8] [utf8_danish_ci utf8 203  Yes 8] [utf8_lithuanian_ci utf8 204  Yes 8] [utf8_slovak_ci utf8 205  Yes 8] [utf8_spanish2_ci utf8 206  Yes 8] [utf8_roman_ci utf8 207  Yes 8] [utf8_persian_ci utf8 208  Yes 8] [utf8_esperanto_ci utf8 209  Yes 8] [utf8_hungarian_ci utf8 210  Yes 8] [utf8_sinhala_ci utf8 211  Yes 8] [utf8_german2_ci utf8 212  Yes 8] [utf8_croatian_ci utf8 213  Yes 8] [utf8_unicode_520_ci utf8 214  Yes 8] [utf8_vietnamese_ci utf8 215  Yes 8] [utf8_general_mysql500_ci utf8 223  Yes 1] [utf8mb4_general_ci utf8mb4 45 Yes Yes 1] [utf8mb4_bin utf8mb4 46  Yes 1] [utf8mb4_unicode_ci utf8mb4 224  Yes 8] [utf8mb4_icelandic_ci utf8mb4 225  Yes 8] [utf8mb4_latvian_ci utf8mb4 226  Yes 8] [utf8mb4_romanian_ci utf8mb4 227  Yes 8] [utf8mb4_slovenian_ci utf8mb4 228  Yes 8] [utf8mb4_polish_ci utf8mb4 229  Yes 8] [utf8mb4_estonian_ci utf8mb4 230  Yes 8] [utf8mb4_spanish_ci utf8mb4 231  Yes 8] [utf8mb4_swedish_ci utf8mb4 232  Yes 8] [utf8mb4_turkish_ci utf8mb4 233  Yes 8] [utf8mb4_czech_ci utf8mb4 234  Yes 8] [utf8mb4_danish_ci utf8mb4 235  Yes 8] [utf8mb4_lithuanian_ci utf8mb4 236  Yes 8] [utf8mb4_slovak_ci utf8mb4 237  Yes 8] [utf8mb4_spanish2_ci utf8mb4 238  Yes 8] [utf8mb4_roman_ci utf8mb4 239  Yes 8] [utf8mb4_persian_ci utf8mb4 240  Yes 8] [utf8mb4_esperanto_ci utf8mb4 241  Yes 8] [utf8mb4_hungarian_ci utf8mb4 242  Yes 8] [utf8mb4_sinhala_ci utf8mb4 243  Yes 8] [utf8mb4_german2_ci utf8mb4 244  Yes 8] [utf8mb4_croatian_ci utf8mb4 245  Yes 8] [utf8mb4_unicode_520_ci utf8mb4 246  Yes 8] [utf8mb4_vietnamese_ci utf8mb4 247  Yes 8]]

sql:show collation where Charset = 'utf8' and Collation = 'utf8_bin';
mysqlRes:[[utf8_bin utf8 83  Yes 1]]
gaeaRes:[[utf8_bin utf8 83  Yes 1]]

sql:show columns in t;
mysqlRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]
gaeaRes:[[id int(11) NO PRI NULL auto_increment] [col1 varchar(20) YES MUL NULL ] [col2 int(11) YES MUL NULL ]]

sql:show full columns in t;
mysqlRes:[[id int(11) NULL NO PRI NULL auto_increment select,insert,update,references ] [col1 varchar(20) utf8mb4_general_ci YES MUL NULL  select,insert,update,references ] [col2 int(11) NULL YES MUL NULL  select,insert,update,references ]]
gaeaRes:[[id int(11) NULL NO PRI NULL auto_increment select,insert,update,references ] [col1 varchar(20) utf8mb4_general_ci YES MUL NULL  select,insert,update,references ] [col2 int(11) NULL YES MUL NULL  select,insert,update,references ]]

sql:show create table sbtest1.t;
mysqlRes:[[t CREATE TABLE `t` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `col1` varchar(20) DEFAULT NULL,
  `col2` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx1` (`col1`),
  KEY `idx2` (`col2`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4]]
gaeaRes:[[t CREATE TABLE `t` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `col1` varchar(20) DEFAULT NULL,
  `col2` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx1` (`col1`),
  KEY `idx2` (`col2`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4]]

sql:show create table t;
mysqlRes:[[t CREATE TABLE `t` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `col1` varchar(20) DEFAULT NULL,
  `col2` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx1` (`col1`),
  KEY `idx2` (`col2`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4]]
gaeaRes:[[t CREATE TABLE `t` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `col1` varchar(20) DEFAULT NULL,
  `col2` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx1` (`col1`),
  KEY `idx2` (`col2`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4]]

sql:show create database sbtest1;
mysqlRes:[[sbtest1 CREATE DATABASE `sbtest1` /*!40100 DEFAULT CHARACTER SET latin1 */]]
gaeaRes:[[sbtest1 CREATE DATABASE `sbtest1` /*!40100 DEFAULT CHARACTER SET latin1 */]]

sql:show create database if not exists sbtest1;
mysqlRes:[[sbtest1 CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sbtest1` /*!40100 DEFAULT CHARACTER SET latin1 */]]
gaeaRes:[[sbtest1 CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sbtest1` /*!40100 DEFAULT CHARACTER SET latin1 */]]

sql:SELECT ++1;
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT -+1;
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT -1;
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT --1;
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select '''a''', """a""";
mysqlRes:[['a' "a"]]
gaeaRes:[['a' "a"]]

sql:select '''a''';
mysqlRes:[['a']]
gaeaRes:[['a']]

sql:select '\'a\'';
mysqlRes:[['a']]
gaeaRes:[['a']]

sql:select "\"a\"";
mysqlRes:[["a"]]
gaeaRes:[["a"]]

sql:select """a""";
mysqlRes:[["a"]]
gaeaRes:[["a"]]

sql:select _utf8"string";
mysqlRes:[[string]]
gaeaRes:[[string]]

sql:select _binary"string";
mysqlRes:[[string]]
gaeaRes:[[string]]

sql:select N'string';
mysqlRes:[[string]]
gaeaRes:[[string]]

sql:select n'string';
mysqlRes:[[string]]
gaeaRes:[[string]]

sql:select 1 <=> 0, 1 <=> null, 1 = null;
mysqlRes:[[0 0 NULL]]
gaeaRes:[[0 0 NULL]]

sql:select {ts '1989-09-10 11:11:11'};
mysqlRes:[[1989-09-10 11:11:11]]
gaeaRes:[[1989-09-10 11:11:11]]

sql:select {d '1989-09-10'};
mysqlRes:[[1989-09-10]]
gaeaRes:[[1989-09-10]]

sql:select {t '00:00:00.111'};
mysqlRes:[[00:00:00.111]]
gaeaRes:[[00:00:00.111]]

sql:select {ts123 '1989-09-10 11:11:11'};
mysqlRes:[[1989-09-10 11:11:11]]
gaeaRes:[[1989-09-10 11:11:11]]

sql:select {ts123 123};
mysqlRes:[[123]]
gaeaRes:[[123]]

sql:select {ts123 1 xor 1};
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT POW(1, 2);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT POW(1, 0.5);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT POW(1, -1);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT POW(-1, 1);
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT MOD(10, 2);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT ROUND(-1.23);
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT ROUND(1.23, 1);
mysqlRes:[[1.2]]
gaeaRes:[[1.2]]

sql:SELECT CEIL(-1.23);
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT CEILING(1.23);
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:SELECT FLOOR(-1.23);
mysqlRes:[[-2]]
gaeaRes:[[-2]]

sql:SELECT LN(1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT LOG(-2);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT LOG(2, 65536);
mysqlRes:[[16]]
gaeaRes:[[16]]

sql:SELECT LOG2(2);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT LOG10(10);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT ABS(10);
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:SELECT CRC32('MySQL');
mysqlRes:[[3259397556]]
gaeaRes:[[3259397556]]

sql:SELECT SIGN(0);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT SQRT(0);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT ACOS(1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT ASIN(1);
mysqlRes:[[1.5707963267948966]]
gaeaRes:[[1.5707963267948966]]

sql:SELECT ATAN(0), ATAN(1), ATAN(1, 2);
mysqlRes:[[0 0.7853981633974483 0.4636476090008061]]
gaeaRes:[[0 0.7853981633974483 0.4636476090008061]]

sql:SELECT COS(0);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT COS(1);
mysqlRes:[[0.5403023058681398]]
gaeaRes:[[0.5403023058681398]]

sql:SELECT COT(1);
mysqlRes:[[0.6420926159343306]]
gaeaRes:[[0.6420926159343306]]

sql:SELECT DEGREES(0);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT EXP(1);
mysqlRes:[[2.718281828459045]]
gaeaRes:[[2.718281828459045]]

sql:SELECT PI();
mysqlRes:[[3.141593]]
gaeaRes:[[3.141593]]

sql:SELECT RADIANS(1);
mysqlRes:[[0.017453292519943295]]
gaeaRes:[[0.017453292519943295]]

sql:SELECT SIN(1);
mysqlRes:[[0.8414709848078965]]
gaeaRes:[[0.8414709848078965]]

sql:SELECT TAN(1);
mysqlRes:[[1.5574077246549023]]
gaeaRes:[[1.5574077246549023]]

sql:SELECT TRUNCATE(1.223,1);
mysqlRes:[[1.2]]
gaeaRes:[[1.2]]

sql:SELECT SUBSTR('Quadratically',5);
mysqlRes:[[ratically]]
gaeaRes:[[ratically]]

sql:SELECT SUBSTR('Quadratically',5, 3);
mysqlRes:[[rat]]
gaeaRes:[[rat]]

sql:SELECT SUBSTR('Quadratically' FROM 5);
mysqlRes:[[ratically]]
gaeaRes:[[ratically]]

sql:SELECT SUBSTR('Quadratically' FROM 5 FOR 3);
mysqlRes:[[rat]]
gaeaRes:[[rat]]

sql:SELECT SUBSTRING('Quadratically',5);
mysqlRes:[[ratically]]
gaeaRes:[[ratically]]

sql:SELECT SUBSTRING('Quadratically',5, 3);
mysqlRes:[[rat]]
gaeaRes:[[rat]]

sql:SELECT SUBSTRING('Quadratically' FROM 5);
mysqlRes:[[ratically]]
gaeaRes:[[ratically]]

sql:SELECT SUBSTRING('Quadratically' FROM 5 FOR 3);
mysqlRes:[[rat]]
gaeaRes:[[rat]]

sql:SELECT CONVERT('111', SIGNED);
mysqlRes:[[111]]
gaeaRes:[[111]]

sql:SELECT LEAST(1, 2, 3);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT INTERVAL(1, 0, 1, 2);
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:SELECT DATE_ADD('2008-01-02', INTERVAL INTERVAL(1, 0, 1) DAY);
mysqlRes:[[2008-01-04]]
gaeaRes:[[2008-01-04]]

sql:SELECT DATABASE();
mysqlRes:[[sbtest1]]
gaeaRes:[[sbtest1]]

sql:SELECT SCHEMA();
mysqlRes:[[sbtest1]]
gaeaRes:[[sbtest1]]

sql:SELECT USER();
mysqlRes:[[superroot@192.168.237.1]]
gaeaRes:[[superroot@192.168.237.1]]

sql:SELECT CURRENT_USER();
mysqlRes:[[superroot@%]]
gaeaRes:[[superroot@%]]

sql:SELECT CURRENT_USER;
mysqlRes:[[superroot@%]]
gaeaRes:[[superroot@%]]

sql:SELECT VERSION();
mysqlRes:[[5.7.26-1-log]]
gaeaRes:[[5.7.26-1-log]]

sql:SELECT BENCHMARK(1000000, AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3')));
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT CHARSET('abc');
mysqlRes:[[utf8mb4]]
gaeaRes:[[utf8mb4]]

sql:SELECT COERCIBILITY('abc');
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:SELECT COLLATION('abc');
mysqlRes:[[utf8mb4_general_ci]]
gaeaRes:[[utf8mb4_general_ci]]

sql:SELECT ROW_COUNT();
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT SESSION_USER();
mysqlRes:[[superroot@192.168.237.1]]
gaeaRes:[[superroot@192.168.237.1]]

sql:SELECT SYSTEM_USER();
mysqlRes:[[superroot@192.168.237.1]]
gaeaRes:[[superroot@192.168.237.1]]

sql:SELECT SUBSTRING_INDEX('www.mysql.com', '.', 2);
mysqlRes:[[www.mysql]]
gaeaRes:[[www.mysql]]

sql:SELECT SUBSTRING_INDEX('www.mysql.com', '.', -2);
mysqlRes:[[mysql.com]]
gaeaRes:[[mysql.com]]

sql:SELECT LOWER("A"), UPPER("a");
mysqlRes:[[a A]]
gaeaRes:[[a A]]

sql:SELECT LCASE("A"), UCASE("a");
mysqlRes:[[a A]]
gaeaRes:[[a A]]

sql:SELECT REPLACE('www.mysql.com', 'w', 'Ww');
mysqlRes:[[WwWwWw.mysql.com]]
gaeaRes:[[WwWwWw.mysql.com]]

sql:SELECT LOCATE('bar', 'foobarbar');
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:SELECT LOCATE('bar', 'foobarbar', 5);
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:select row(1, 1) > row(1, 1), row(1, 1, 1) > row(1, 1, 1);
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:Select (1, 1) > (1, 1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT *, CAST(col1 AS CHAR CHARACTER SET utf8) FROM t;
mysqlRes:[[1 aa 5 aa] [2 bb 10 bb] [3 cc 15 cc] [4 dd 20 dd] [5 ee 30 ee] [6 aa 5 aa] [7 bb 10 bb] [8 cc 15 cc] [9 dd 20 dd] [10 ee 30 ee]]
gaeaRes:[[1 aa 5 aa] [2 bb 10 bb] [3 cc 15 cc] [4 dd 20 dd] [5 ee 30 ee] [6 aa 5 aa] [7 bb 10 bb] [8 cc 15 cc] [9 dd 20 dd] [10 ee 30 ee]]

sql:select cast(1 as signed int);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT last_insert_id();
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT binary 'a';
mysqlRes:[[a]]
gaeaRes:[[a]]

sql:SELECT BIT_COUNT(1);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT time('01:02:03');
mysqlRes:[[01:02:03]]
gaeaRes:[[01:02:03]]

sql:SELECT time('01:02:03.1');
mysqlRes:[[01:02:03.1]]
gaeaRes:[[01:02:03.1]]

sql:SELECT time('20.1');
mysqlRes:[[00:00:20.1]]
gaeaRes:[[00:00:20.1]]

sql:SELECT TIMEDIFF('2000:01:01 00:00:00', '2000:01:01 00:00:00.000001');
mysqlRes:[[-00:00:00.000001]]
gaeaRes:[[-00:00:00.000001]]

sql:SELECT TIMESTAMPDIFF(MONTH,'2003-02-01','2003-05-01');
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:SELECT TIMESTAMPDIFF(YEAR,'2002-05-01','2001-01-01');
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT TIMESTAMPDIFF(MINUTE,'2003-02-01','2003-05-01 12:05:55');
mysqlRes:[[128885]]
gaeaRes:[[128885]]

sql:SELECT MICROSECOND('2009-12-31 23:59:59.000010');
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:SELECT SECOND('10:05:03');
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:SELECT MINUTE('2008-02-03 10:05:03');
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE();
mysqlRes:[[2023-10-24 2023-10-24 2023-10-24]]
gaeaRes:[[2023-10-24 2023-10-24 2023-10-24]]

sql:SELECT DATEDIFF('2003-12-31', '2003-12-30');
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT DATE('2003-12-31 01:02:03');
mysqlRes:[[2003-12-31]]
gaeaRes:[[2003-12-31]]

sql:SELECT DATE_FORMAT('2003-12-31 01:02:03', '%W %M %Y');
mysqlRes:[[Wednesday December 2003]]
gaeaRes:[[Wednesday December 2003]]

sql:SELECT DAY('2007-02-03');
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:SELECT DAYOFMONTH('2007-02-03');
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:SELECT DAYOFWEEK('2007-02-03');
mysqlRes:[[7]]
gaeaRes:[[7]]

sql:SELECT DAYOFYEAR('2007-02-03');
mysqlRes:[[34]]
gaeaRes:[[34]]

sql:SELECT DAYNAME('2007-02-03');
mysqlRes:[[Saturday]]
gaeaRes:[[Saturday]]

sql:SELECT FROM_DAYS(1423);
mysqlRes:[[0003-11-24]]
gaeaRes:[[0003-11-24]]

sql:SELECT WEEKDAY('2007-02-03');
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:SELECT UTC_DATE, UTC_DATE();
mysqlRes:[[2023-10-24 2023-10-24]]
gaeaRes:[[2023-10-24 2023-10-24]]

sql:SELECT UTC_DATE(), UTC_DATE()+0;
mysqlRes:[[2023-10-24 20231024]]
gaeaRes:[[2023-10-24 20231024]]

sql:SELECT WEEK('2007-02-03');
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:SELECT WEEK('2007-02-03', 0);
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:SELECT WEEKOFYEAR('2007-02-03');
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:SELECT MONTH('2007-02-03');
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:SELECT MONTHNAME('2007-02-03');
mysqlRes:[[February]]
gaeaRes:[[February]]

sql:SELECT YEAR('2007-02-03');
mysqlRes:[[2007]]
gaeaRes:[[2007]]

sql:SELECT YEARWEEK('2007-02-03');
mysqlRes:[[200704]]
gaeaRes:[[200704]]

sql:SELECT YEARWEEK('2007-02-03', 0);
mysqlRes:[[200704]]
gaeaRes:[[200704]]

sql:SELECT ADDTIME('01:00:00.999999', '02:00:00.999998');
mysqlRes:[[03:00:01.999997]]
gaeaRes:[[03:00:01.999997]]

sql:SELECT SUBTIME('01:00:00.999999', '02:00:00.999998');
mysqlRes:[[-00:59:59.999999]]
gaeaRes:[[-00:59:59.999999]]

sql:SELECT CONVERT_TZ('2004-01-01 12:00:00','+00:00','+10:00');
mysqlRes:[[2004-01-01 22:00:00]]
gaeaRes:[[2004-01-01 22:00:00]]

sql:SELECT GET_FORMAT(DATE, 'USA');
mysqlRes:[[%m.%d.%Y]]
gaeaRes:[[%m.%d.%Y]]

sql:SELECT GET_FORMAT(DATETIME, 'USA');
mysqlRes:[[%Y-%m-%d %H.%i.%s]]
gaeaRes:[[%Y-%m-%d %H.%i.%s]]

sql:SELECT GET_FORMAT(TIME, 'USA');
mysqlRes:[[%h:%i:%s %p]]
gaeaRes:[[%h:%i:%s %p]]

sql:SELECT GET_FORMAT(TIMESTAMP, 'USA');
mysqlRes:[[%Y-%m-%d %H.%i.%s]]
gaeaRes:[[%Y-%m-%d %H.%i.%s]]

sql:SELECT MAKEDATE(2011,31);
mysqlRes:[[2011-01-31]]
gaeaRes:[[2011-01-31]]

sql:SELECT MAKETIME(12,15,30);
mysqlRes:[[12:15:30]]
gaeaRes:[[12:15:30]]

sql:SELECT PERIOD_ADD(200801,2);
mysqlRes:[[200803]]
gaeaRes:[[200803]]

sql:SELECT PERIOD_DIFF(200802,200703);
mysqlRes:[[11]]
gaeaRes:[[11]]

sql:SELECT QUARTER('2008-04-01');
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:SELECT SEC_TO_TIME(2378);
mysqlRes:[[00:39:38]]
gaeaRes:[[00:39:38]]

sql:SELECT TIME_FORMAT('100:00:00', '%H %k %h %I %l');
mysqlRes:[[100 100 04 04 4]]
gaeaRes:[[100 100 04 04 4]]

sql:SELECT TIME_TO_SEC('22:23:00');
mysqlRes:[[80580]]
gaeaRes:[[80580]]

sql:SELECT TIMESTAMPADD(WEEK,1,'2003-01-02');
mysqlRes:[[2003-01-09]]
gaeaRes:[[2003-01-09]]

sql:SELECT TO_DAYS('2007-10-07');
mysqlRes:[[733321]]
gaeaRes:[[733321]]

sql:SELECT TO_SECONDS('2009-11-29');
mysqlRes:[[63426672000]]
gaeaRes:[[63426672000]]

sql:SELECT LAST_DAY('2003-02-05');
mysqlRes:[[2003-02-28]]
gaeaRes:[[2003-02-28]]

sql:select extract(microsecond from "2011-11-11 10:10:10.123456");
mysqlRes:[[123456]]
gaeaRes:[[123456]]

sql:select extract(second from "2011-11-11 10:10:10.123456");
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select extract(minute from "2011-11-11 10:10:10.123456");
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select extract(hour from "2011-11-11 10:10:10.123456");
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select extract(day from "2011-11-11 10:10:10.123456");
mysqlRes:[[11]]
gaeaRes:[[11]]

sql:select extract(week from "2011-11-11 10:10:10.123456");
mysqlRes:[[45]]
gaeaRes:[[45]]

sql:select extract(month from "2011-11-11 10:10:10.123456");
mysqlRes:[[11]]
gaeaRes:[[11]]

sql:select extract(quarter from "2011-11-11 10:10:10.123456");
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:select extract(year from "2011-11-11 10:10:10.123456");
mysqlRes:[[2011]]
gaeaRes:[[2011]]

sql:select extract(second_microsecond from "2011-11-11 10:10:10.123456");
mysqlRes:[[10123456]]
gaeaRes:[[10123456]]

sql:select extract(minute_microsecond from "2011-11-11 10:10:10.123456");
mysqlRes:[[1010123456]]
gaeaRes:[[1010123456]]

sql:select extract(minute_second from "2011-11-11 10:10:10.123456");
mysqlRes:[[1010]]
gaeaRes:[[1010]]

sql:select extract(hour_microsecond from "2011-11-11 10:10:10.123456");
mysqlRes:[[101010123456]]
gaeaRes:[[101010123456]]

sql:select extract(hour_second from "2011-11-11 10:10:10.123456");
mysqlRes:[[101010]]
gaeaRes:[[101010]]

sql:select extract(hour_minute from "2011-11-11 10:10:10.123456");
mysqlRes:[[1010]]
gaeaRes:[[1010]]

sql:select extract(day_microsecond from "2011-11-11 10:10:10.123456");
mysqlRes:[[11101010123456]]
gaeaRes:[[11101010123456]]

sql:select extract(day_second from "2011-11-11 10:10:10.123456");
mysqlRes:[[11101010]]
gaeaRes:[[11101010]]

sql:select extract(day_minute from "2011-11-11 10:10:10.123456");
mysqlRes:[[111010]]
gaeaRes:[[111010]]

sql:select extract(day_hour from "2011-11-11 10:10:10.123456");
mysqlRes:[[1110]]
gaeaRes:[[1110]]

sql:select extract(year_month from "2011-11-11 10:10:10.123456");
mysqlRes:[[201111]]
gaeaRes:[[201111]]

sql:select from_unixtime(1447430881);
mysqlRes:[[2015-11-13 16:08:01]]
gaeaRes:[[2015-11-13 16:08:01]]

sql:select from_unixtime(1447430881.123456);
mysqlRes:[[2015-11-13 16:08:01.123456]]
gaeaRes:[[2015-11-13 16:08:01.123456]]

sql:select from_unixtime(1447430881.1234567);
mysqlRes:[[2015-11-13 16:08:01.123457]]
gaeaRes:[[2015-11-13 16:08:01.123457]]

sql:select from_unixtime(1447430881.9999999);
mysqlRes:[[2015-11-13 16:08:02.000000]]
gaeaRes:[[2015-11-13 16:08:02.000000]]

sql:select from_unixtime(1447430881, "%Y %D %M %h:%i:%s %x");
mysqlRes:[[2015 13th November 04:08:01 2015]]
gaeaRes:[[2015 13th November 04:08:01 2015]]

sql:select from_unixtime(1447430881.123456, "%Y %D %M %h:%i:%s %x");
mysqlRes:[[2015 13th November 04:08:01 2015]]
gaeaRes:[[2015 13th November 04:08:01 2015]]

sql:select from_unixtime(1447430881.1234567, "%Y %D %M %h:%i:%s %x");
mysqlRes:[[2015 13th November 04:08:01 2015]]
gaeaRes:[[2015 13th November 04:08:01 2015]]

sql:SELECT CAST('test collated returns' AS CHAR CHARACTER SET utf8) COLLATE utf8_bin;
mysqlRes:[[test collated returns]]
gaeaRes:[[test collated returns]]

sql:SELECT TRIM('  bar   ');
mysqlRes:[[bar]]
gaeaRes:[[bar]]

sql:SELECT TRIM(LEADING 'x' FROM 'xxxbarxxx');
mysqlRes:[[barxxx]]
gaeaRes:[[barxxx]]

sql:SELECT TRIM(BOTH 'x' FROM 'xxxbarxxx');
mysqlRes:[[bar]]
gaeaRes:[[bar]]

sql:SELECT TRIM(TRAILING 'xyz' FROM 'barxxyz');
mysqlRes:[[barx]]
gaeaRes:[[barx]]

sql:SELECT LTRIM(' foo ');
mysqlRes:[[foo ]]
gaeaRes:[[foo ]]

sql:SELECT RTRIM(' bar ');
mysqlRes:[[ bar]]
gaeaRes:[[ bar]]

sql:SELECT RPAD('hi', 6, 'c');
mysqlRes:[[hicccc]]
gaeaRes:[[hicccc]]

sql:SELECT BIT_LENGTH('hi');
mysqlRes:[[16]]
gaeaRes:[[16]]

sql:SELECT CHAR_LENGTH('abc');
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:SELECT CHARACTER_LENGTH('abc');
mysqlRes:[[3]]
gaeaRes:[[3]]

sql:SELECT FIELD('ej', 'Hej', 'ej', 'Heja', 'hej', 'foo');
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:SELECT FIND_IN_SET('foo', 'foo,bar');
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT MAKE_SET(1,'a'), MAKE_SET(1,'a','b','c');
mysqlRes:[[a a]]
gaeaRes:[[a a]]

sql:SELECT MID('Sakila', -5, 3);
mysqlRes:[[aki]]
gaeaRes:[[aki]]

sql:SELECT OCT(12);
mysqlRes:[[14]]
gaeaRes:[[14]]

sql:SELECT OCTET_LENGTH('text');
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:SELECT ORD('2');
mysqlRes:[[50]]
gaeaRes:[[50]]

sql:SELECT POSITION('bar' IN 'foobarbar');
mysqlRes:[[4]]
gaeaRes:[[4]]

sql:SELECT BIN(12);
mysqlRes:[[1100]]
gaeaRes:[[1100]]

sql:SELECT ELT(1, 'ej', 'Heja', 'hej', 'foo');
mysqlRes:[[ej]]
gaeaRes:[[ej]]

sql:SELECT EXPORT_SET(5,'Y','N'), EXPORT_SET(5,'Y','N',','), EXPORT_SET(5,'Y','N',',',4);
mysqlRes:[[Y,N,Y,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N Y,N,Y,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N Y,N,Y,N]]
gaeaRes:[[Y,N,Y,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N Y,N,Y,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N,N Y,N,Y,N]]

sql:SELECT FROM_BASE64('abc');
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT TO_BASE64('abc');
mysqlRes:[[YWJj]]
gaeaRes:[[YWJj]]

sql:SELECT LOAD_FILE('/tmp/picture');
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT LPAD('hi',4,'??');
mysqlRes:[[??hi]]
gaeaRes:[[??hi]]

sql:SELECT LEFT("foobar", 3);
mysqlRes:[[foo]]
gaeaRes:[[foo]]

sql:SELECT RIGHT("foobar", 3);
mysqlRes:[[bar]]
gaeaRes:[[bar]]

sql:SELECT REPEAT("a", 10);
mysqlRes:[[aaaaaaaaaa]]
gaeaRes:[[aaaaaaaaaa]]

sql:SELECT SLEEP(1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT ANY_VALUE(@arg);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT INET_ATON('10.0.5.9');
mysqlRes:[[167773449]]
gaeaRes:[[167773449]]

sql:SELECT INET_NTOA(167773449);
mysqlRes:[[10.0.5.9]]
gaeaRes:[[10.0.5.9]]

sql:SELECT INET6_ATON('fdfe::5a55:caff:fefa:9089');
mysqlRes:[[��      ZU������]]
gaeaRes:[[��      ZU������]]

sql:SELECT INET6_NTOA(INET_NTOA(167773449));
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT IS_IPV4('10.0.5.9');
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT IS_IPV4_COMPAT(INET6_ATON('::10.0.5.9'));
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT IS_IPV4_MAPPED(INET6_ATON('::10.0.5.9'));
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT IS_IPV6('10.0.5.9');
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT MASTER_POS_WAIT(@log_name, @log_pos), MASTER_POS_WAIT(@log_name, @log_pos, @timeout), MASTER_POS_WAIT(@log_name, @log_pos, @timeout, @channel_name);
mysqlRes:[[NULL NULL NULL]]
gaeaRes:[[NULL NULL NULL]]

sql:SELECT NAME_CONST('myname', 14);
mysqlRes:[[14]]
gaeaRes:[[14]]

sql:SELECT RELEASE_ALL_LOCKS();
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select "2011-11-11 10:10:10.123456" + interval 10 second;
mysqlRes:[[2011-11-11 10:10:20.123456]]
gaeaRes:[[2011-11-11 10:10:20.123456]]

sql:select "2011-11-11 10:10:10.123456" - interval 10 second;
mysqlRes:[[2011-11-11 10:10:00.123456]]
gaeaRes:[[2011-11-11 10:10:00.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 10 microsecond);
mysqlRes:[[2011-11-11 10:10:10.123466]]
gaeaRes:[[2011-11-11 10:10:10.123466]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 10 second);
mysqlRes:[[2011-11-11 10:10:20.123456]]
gaeaRes:[[2011-11-11 10:10:20.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 10 minute);
mysqlRes:[[2011-11-11 10:20:10.123456]]
gaeaRes:[[2011-11-11 10:20:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 10 hour);
mysqlRes:[[2011-11-11 20:10:10.123456]]
gaeaRes:[[2011-11-11 20:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 10 day);
mysqlRes:[[2011-11-21 10:10:10.123456]]
gaeaRes:[[2011-11-21 10:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 1 week);
mysqlRes:[[2011-11-18 10:10:10.123456]]
gaeaRes:[[2011-11-18 10:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 1 month);
mysqlRes:[[2011-12-11 10:10:10.123456]]
gaeaRes:[[2011-12-11 10:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 1 quarter);
mysqlRes:[[2012-02-11 10:10:10.123456]]
gaeaRes:[[2012-02-11 10:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 1 year);
mysqlRes:[[2012-11-11 10:10:10.123456]]
gaeaRes:[[2012-11-11 10:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
mysqlRes:[[2011-11-11 10:10:20.223456]]
gaeaRes:[[2011-11-11 10:10:20.223456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
mysqlRes:[[2011-11-11 10:20:20.223456]]
gaeaRes:[[2011-11-11 10:20:20.223456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
mysqlRes:[[2011-11-11 10:20:20.123456]]
gaeaRes:[[2011-11-11 10:20:20.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
mysqlRes:[[2011-11-11 20:20:20.223456]]
gaeaRes:[[2011-11-11 20:20:20.223456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
mysqlRes:[[2011-11-11 20:20:20.123456]]
gaeaRes:[[2011-11-11 20:20:20.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
mysqlRes:[[2011-11-11 20:20:10.123456]]
gaeaRes:[[2011-11-11 20:20:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
mysqlRes:[[2011-11-11 20:20:10.123456]]
gaeaRes:[[2011-11-11 20:20:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
mysqlRes:[[2011-11-22 20:20:20.223456]]
gaeaRes:[[2011-11-22 20:20:20.223456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
mysqlRes:[[2011-11-22 20:20:20.123456]]
gaeaRes:[[2011-11-22 20:20:20.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
mysqlRes:[[2011-11-22 20:20:10.123456]]
gaeaRes:[[2011-11-22 20:20:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
mysqlRes:[[2011-11-22 20:10:10.123456]]
gaeaRes:[[2011-11-22 20:10:10.123456]]

sql:select date_add("2011-11-11 10:10:10.123456", interval "11-11" year_month);
mysqlRes:[[2023-10-11 10:10:10.123456]]
gaeaRes:[[2023-10-11 10:10:10.123456]]

sql:select strcmp('abc', 'def');
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 10 microsecond);
mysqlRes:[[2011-11-11 10:10:10.123466]]
gaeaRes:[[2011-11-11 10:10:10.123466]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 10 second);
mysqlRes:[[2011-11-11 10:10:20.123456]]
gaeaRes:[[2011-11-11 10:10:20.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 10 minute);
mysqlRes:[[2011-11-11 10:20:10.123456]]
gaeaRes:[[2011-11-11 10:20:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 10 hour);
mysqlRes:[[2011-11-11 20:10:10.123456]]
gaeaRes:[[2011-11-11 20:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 10 day);
mysqlRes:[[2011-11-21 10:10:10.123456]]
gaeaRes:[[2011-11-21 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 1 week);
mysqlRes:[[2011-11-18 10:10:10.123456]]
gaeaRes:[[2011-11-18 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 1 month);
mysqlRes:[[2011-12-11 10:10:10.123456]]
gaeaRes:[[2011-12-11 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 1 quarter);
mysqlRes:[[2012-02-11 10:10:10.123456]]
gaeaRes:[[2012-02-11 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 1 year);
mysqlRes:[[2012-11-11 10:10:10.123456]]
gaeaRes:[[2012-11-11 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
mysqlRes:[[2011-11-11 10:10:20.223456]]
gaeaRes:[[2011-11-11 10:10:20.223456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
mysqlRes:[[2011-11-11 10:20:20.223456]]
gaeaRes:[[2011-11-11 10:20:20.223456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
mysqlRes:[[2011-11-11 10:20:20.123456]]
gaeaRes:[[2011-11-11 10:20:20.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
mysqlRes:[[2011-11-11 20:20:20.223456]]
gaeaRes:[[2011-11-11 20:20:20.223456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
mysqlRes:[[2011-11-11 20:20:20.123456]]
gaeaRes:[[2011-11-11 20:20:20.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
mysqlRes:[[2011-11-11 20:20:10.123456]]
gaeaRes:[[2011-11-11 20:20:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
mysqlRes:[[2011-11-11 20:20:10.123456]]
gaeaRes:[[2011-11-11 20:20:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
mysqlRes:[[2011-11-22 20:20:20.223456]]
gaeaRes:[[2011-11-22 20:20:20.223456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
mysqlRes:[[2011-11-22 20:20:20.123456]]
gaeaRes:[[2011-11-22 20:20:20.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
mysqlRes:[[2011-11-22 20:20:10.123456]]
gaeaRes:[[2011-11-22 20:20:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
mysqlRes:[[2011-11-22 20:10:10.123456]]
gaeaRes:[[2011-11-22 20:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", interval "11-11" year_month);
mysqlRes:[[2023-10-11 10:10:10.123456]]
gaeaRes:[[2023-10-11 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", 10);
mysqlRes:[[2011-11-21 10:10:10.123456]]
gaeaRes:[[2011-11-21 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", 0.10);
mysqlRes:[[2011-11-11 10:10:10.123456]]
gaeaRes:[[2011-11-11 10:10:10.123456]]

sql:select adddate("2011-11-11 10:10:10.123456", "11,11");
mysqlRes:[[2011-11-22 10:10:10.123456]]
gaeaRes:[[2011-11-22 10:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 10 microsecond);
mysqlRes:[[2011-11-11 10:10:10.123446]]
gaeaRes:[[2011-11-11 10:10:10.123446]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 10 second);
mysqlRes:[[2011-11-11 10:10:00.123456]]
gaeaRes:[[2011-11-11 10:10:00.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 10 minute);
mysqlRes:[[2011-11-11 10:00:10.123456]]
gaeaRes:[[2011-11-11 10:00:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 10 hour);
mysqlRes:[[2011-11-11 00:10:10.123456]]
gaeaRes:[[2011-11-11 00:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 10 day);
mysqlRes:[[2011-11-01 10:10:10.123456]]
gaeaRes:[[2011-11-01 10:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 1 week);
mysqlRes:[[2011-11-04 10:10:10.123456]]
gaeaRes:[[2011-11-04 10:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 1 month);
mysqlRes:[[2011-10-11 10:10:10.123456]]
gaeaRes:[[2011-10-11 10:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 1 quarter);
mysqlRes:[[2011-08-11 10:10:10.123456]]
gaeaRes:[[2011-08-11 10:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 1 year);
mysqlRes:[[2010-11-11 10:10:10.123456]]
gaeaRes:[[2010-11-11 10:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
mysqlRes:[[2011-11-11 10:10:00.023456]]
gaeaRes:[[2011-11-11 10:10:00.023456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
mysqlRes:[[2011-11-11 10:00:00.023456]]
gaeaRes:[[2011-11-11 10:00:00.023456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
mysqlRes:[[2011-11-11 10:00:00.123456]]
gaeaRes:[[2011-11-11 10:00:00.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
mysqlRes:[[2011-11-11 00:00:00.023456]]
gaeaRes:[[2011-11-11 00:00:00.023456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
mysqlRes:[[2011-11-11 00:00:00.123456]]
gaeaRes:[[2011-11-11 00:00:00.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
mysqlRes:[[2011-11-11 00:00:10.123456]]
gaeaRes:[[2011-11-11 00:00:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
mysqlRes:[[2011-11-11 00:00:10.123456]]
gaeaRes:[[2011-11-11 00:00:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
mysqlRes:[[2011-10-31 00:00:00.023456]]
gaeaRes:[[2011-10-31 00:00:00.023456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
mysqlRes:[[2011-10-31 00:00:00.123456]]
gaeaRes:[[2011-10-31 00:00:00.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
mysqlRes:[[2011-10-31 00:00:10.123456]]
gaeaRes:[[2011-10-31 00:00:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
mysqlRes:[[2011-10-31 00:10:10.123456]]
gaeaRes:[[2011-10-31 00:10:10.123456]]

sql:select date_sub("2011-11-11 10:10:10.123456", interval "11-11" year_month);
mysqlRes:[[1999-12-11 10:10:10.123456]]
gaeaRes:[[1999-12-11 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 10 microsecond);
mysqlRes:[[2011-11-11 10:10:10.123446]]
gaeaRes:[[2011-11-11 10:10:10.123446]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 10 second);
mysqlRes:[[2011-11-11 10:10:00.123456]]
gaeaRes:[[2011-11-11 10:10:00.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 10 minute);
mysqlRes:[[2011-11-11 10:00:10.123456]]
gaeaRes:[[2011-11-11 10:00:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 10 hour);
mysqlRes:[[2011-11-11 00:10:10.123456]]
gaeaRes:[[2011-11-11 00:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 10 day);
mysqlRes:[[2011-11-01 10:10:10.123456]]
gaeaRes:[[2011-11-01 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 1 week);
mysqlRes:[[2011-11-04 10:10:10.123456]]
gaeaRes:[[2011-11-04 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 1 month);
mysqlRes:[[2011-10-11 10:10:10.123456]]
gaeaRes:[[2011-10-11 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 1 quarter);
mysqlRes:[[2011-08-11 10:10:10.123456]]
gaeaRes:[[2011-08-11 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 1 year);
mysqlRes:[[2010-11-11 10:10:10.123456]]
gaeaRes:[[2010-11-11 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
mysqlRes:[[2011-11-11 10:10:00.023456]]
gaeaRes:[[2011-11-11 10:10:00.023456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
mysqlRes:[[2011-11-11 10:00:00.023456]]
gaeaRes:[[2011-11-11 10:00:00.023456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
mysqlRes:[[2011-11-11 10:00:00.123456]]
gaeaRes:[[2011-11-11 10:00:00.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
mysqlRes:[[2011-11-11 00:00:00.023456]]
gaeaRes:[[2011-11-11 00:00:00.023456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
mysqlRes:[[2011-11-11 00:00:00.123456]]
gaeaRes:[[2011-11-11 00:00:00.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
mysqlRes:[[2011-11-11 00:00:10.123456]]
gaeaRes:[[2011-11-11 00:00:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
mysqlRes:[[2011-11-11 00:00:10.123456]]
gaeaRes:[[2011-11-11 00:00:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
mysqlRes:[[2011-10-31 00:00:00.023456]]
gaeaRes:[[2011-10-31 00:00:00.023456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
mysqlRes:[[2011-10-31 00:00:00.123456]]
gaeaRes:[[2011-10-31 00:00:00.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
mysqlRes:[[2011-10-31 00:00:10.123456]]
gaeaRes:[[2011-10-31 00:00:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
mysqlRes:[[2011-10-31 00:10:10.123456]]
gaeaRes:[[2011-10-31 00:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", interval "11-11" year_month);
mysqlRes:[[1999-12-11 10:10:10.123456]]
gaeaRes:[[1999-12-11 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", 10);
mysqlRes:[[2011-11-01 10:10:10.123456]]
gaeaRes:[[2011-11-01 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", 0.10);
mysqlRes:[[2011-11-11 10:10:10.123456]]
gaeaRes:[[2011-11-11 10:10:10.123456]]

sql:select subdate("2011-11-11 10:10:10.123456", "11,11");
mysqlRes:[[2011-10-31 10:10:10.123456]]
gaeaRes:[[2011-10-31 10:10:10.123456]]

sql:select unix_timestamp();
mysqlRes:[[1698133401]]
gaeaRes:[[1698133401]]

sql:select unix_timestamp('2015-11-13 10:20:19.012');
mysqlRes:[[1447410019.012]]
gaeaRes:[[1447410019.012]]

sql:select avg(distinct col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select avg(distinctrow col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select avg(distinct all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select avg(distinctrow all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select avg(col2) from t;
mysqlRes:[[16.0000]]
gaeaRes:[[16.0000]]

sql:select bit_and(col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select bit_and(all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select bit_or(col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select bit_or(all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select bit_xor(col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select bit_xor(all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select max(distinct col1) from t;
mysqlRes:[[ee]]
gaeaRes:[[ee]]

sql:select max(distinctrow col1) from t;
mysqlRes:[[ee]]
gaeaRes:[[ee]]

sql:select max(distinct all col1) from t;
mysqlRes:[[ee]]
gaeaRes:[[ee]]

sql:select max(distinctrow all col1) from t;
mysqlRes:[[ee]]
gaeaRes:[[ee]]

sql:select max(col2) from t;
mysqlRes:[[30]]
gaeaRes:[[30]]

sql:select min(distinct col1) from t;
mysqlRes:[[aa]]
gaeaRes:[[aa]]

sql:select min(distinctrow col1) from t;
mysqlRes:[[aa]]
gaeaRes:[[aa]]

sql:select min(distinct all col1) from t;
mysqlRes:[[aa]]
gaeaRes:[[aa]]

sql:select min(distinctrow all col1) from t;
mysqlRes:[[aa]]
gaeaRes:[[aa]]

sql:select min(col2) from t;
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select sum(distinct col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select sum(distinctrow col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select sum(distinct all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select sum(distinctrow all col1) from t;
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select sum(col2) from t;
mysqlRes:[[160]]
gaeaRes:[[160]]

sql:select count(col1) from t;
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select count(*) from t;
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select count(distinct col1, col2) from t;
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select count(distinctrow col1, col2) from t;
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select count(all col1) from t;
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select group_concat(col2,col1) from t group by col1;
mysqlRes:[[5aa,5aa] [10bb,10bb] [15cc,15cc] [20dd,20dd] [30ee,30ee]]
gaeaRes:[[5aa,5aa] [10bb,10bb] [15cc,15cc] [20dd,20dd] [30ee,30ee]]

sql:select group_concat(col2,col1 SEPARATOR ';') from t group by col1;
mysqlRes:[[5aa;5aa] [10bb;10bb] [15cc;15cc] [20dd;20dd] [30ee;30ee]]
gaeaRes:[[5aa;5aa] [10bb;10bb] [15cc;15cc] [20dd;20dd] [30ee;30ee]]

sql:select group_concat(distinct col2,col1) from t group by col1;
mysqlRes:[[5aa] [10bb] [15cc] [20dd] [30ee]]
gaeaRes:[[5aa] [10bb] [15cc] [20dd] [30ee]]

sql:select group_concat(distinctrow col2,col1) from t group by col1;
mysqlRes:[[5aa] [10bb] [15cc] [20dd] [30ee]]
gaeaRes:[[5aa] [10bb] [15cc] [20dd] [30ee]]

sql:SELECT col1, GROUP_CONCAT(DISTINCT col2 ORDER BY col2 DESC SEPARATOR ' ') FROM t GROUP BY col1;
mysqlRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30]]
gaeaRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30]]

sql:select std(col1), std(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select stddev(col1), stddev(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select stddev_pop(col1), stddev_pop(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select stddev_samp(col1), stddev_samp(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select variance(col1), variance(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select var_pop(col1), var_pop(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select var_samp(col1), var_samp(all col1) from t;
mysqlRes:[[0 0]]
gaeaRes:[[0 0]]

sql:select AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3'));
mysqlRes:[[�lS�1��GsH?���]]
gaeaRes:[[�lS�1��GsH?���]]

sql:select AES_DECRYPT(@crypt_str,@key_str);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:select AES_DECRYPT(@crypt_str,@key_str,@init_vector);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT COMPRESS('');
mysqlRes:[[]]
gaeaRes:[[]]

sql:SELECT MD5('testing');
mysqlRes:[[ae2b1fca515949e5d54fb22b8ed95575]]
gaeaRes:[[ae2b1fca515949e5d54fb22b8ed95575]]

sql:SELECT RANDOM_BYTES(@len);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT SHA1('abc');
mysqlRes:[[a9993e364706816aba3e25717850c26c9cd0d89d]]
gaeaRes:[[a9993e364706816aba3e25717850c26c9cd0d89d]]

sql:SELECT SHA('abc');
mysqlRes:[[a9993e364706816aba3e25717850c26c9cd0d89d]]
gaeaRes:[[a9993e364706816aba3e25717850c26c9cd0d89d]]

sql:SELECT SHA2('abc', 224);
mysqlRes:[[23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7]]
gaeaRes:[[23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7]]

sql:SELECT UNCOMPRESS('any string');
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT UNCOMPRESSED_LENGTH(@compressed_string);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT VALIDATE_PASSWORD_STRENGTH(@str);
mysqlRes:[[NULL]]
gaeaRes:[[NULL]]

sql:SELECT CHAR(65);
mysqlRes:[[A]]
gaeaRes:[[A]]

sql:SELECT CHAR(65, 66, 67);
mysqlRes:[[ABC]]
gaeaRes:[[ABC]]

sql:SELECT HEX(CHAR(1, 0)),HEX(CHAR(256)),HEX(CHAR(1, 1)),HEX(CHAR(257));
mysqlRes:[[0100 0100 0101 0101]]
gaeaRes:[[0100 0100 0101 0101]]

sql:SELECT CHAR(0x027FA USING ucs2);
mysqlRes:[[⟺]]
gaeaRes:[[⟺]]

sql:SELECT CHAR(0xc2a7 USING utf8);
mysqlRes:[[§]]
gaeaRes:[[§]]

sql:select 1 as a, 1 as `a`, 1 as 'a';
mysqlRes:[[1 1 1]]
gaeaRes:[[1 1 1]]

sql:select 1 as a, 1 as "a", 1 as 'a';
mysqlRes:[[1 1 1]]
gaeaRes:[[1 1 1]]

sql:select 1 a, 1 "a", 1 'a';
mysqlRes:[[1 1 1]]
gaeaRes:[[1 1 1]]

sql:select * from t a;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t as a;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select 1 full, 1 `row`, 1 abs;
mysqlRes:[[1 1 1]]
gaeaRes:[[1 1 1]]

sql:select * from t full, t1 `row`, t2 abs;
mysqlRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5 1 aa 5] [2 bb 10 1 aa 5 1 aa 5] [3 cc 15 1 aa 5 1 aa 5] [4 dd 20 1 aa 5 1 aa 5] [5 ee 30 1 aa 5 1 aa 5] [6 aa 5 1 aa 5 1 aa 5] [7 bb 10 1 aa 5 1 aa 5] [8 cc 15 1 aa 5 1 aa 5] [9 dd 20 1 aa 5 1 aa 5] [10 ee 30 1 aa 5 1 aa 5] [1 aa 5 2 bb 10 1 aa 5] [2 bb 10 2 bb 10 1 aa 5] [3 cc 15 2 bb 10 1 aa 5] [4 dd 20 2 bb 10 1 aa 5] [5 ee 30 2 bb 10 1 aa 5] [6 aa 5 2 bb 10 1 aa 5] [7 bb 10 2 bb 10 1 aa 5] [8 cc 15 2 bb 10 1 aa 5] [9 dd 20 2 bb 10 1 aa 5] [10 ee 30 2 bb 10 1 aa 5] [1 aa 5 3 cc 15 1 aa 5] [2 bb 10 3 cc 15 1 aa 5] [3 cc 15 3 cc 15 1 aa 5] [4 dd 20 3 cc 15 1 aa 5] [5 ee 30 3 cc 15 1 aa 5] [6 aa 5 3 cc 15 1 aa 5] [7 bb 10 3 cc 15 1 aa 5] [8 cc 15 3 cc 15 1 aa 5] [9 dd 20 3 cc 15 1 aa 5] [10 ee 30 3 cc 15 1 aa 5] [1 aa 5 4 dd 20 1 aa 5] [2 bb 10 4 dd 20 1 aa 5] [3 cc 15 4 dd 20 1 aa 5] [4 dd 20 4 dd 20 1 aa 5] [5 ee 30 4 dd 20 1 aa 5] [6 aa 5 4 dd 20 1 aa 5] [7 bb 10 4 dd 20 1 aa 5] [8 cc 15 4 dd 20 1 aa 5] [9 dd 20 4 dd 20 1 aa 5] [10 ee 30 4 dd 20 1 aa 5] [1 aa 5 5 ee 30 1 aa 5] [2 bb 10 5 ee 30 1 aa 5] [3 cc 15 5 ee 30 1 aa 5] [4 dd 20 5 ee 30 1 aa 5] [5 ee 30 5 ee 30 1 aa 5] [6 aa 5 5 ee 30 1 aa 5] [7 bb 10 5 ee 30 1 aa 5] [8 cc 15 5 ee 30 1 aa 5] [9 dd 20 5 ee 30 1 aa 5] [10 ee 30 5 ee 30 1 aa 5] [1 aa 5 6 aa 5 1 aa 5] [2 bb 10 6 aa 5 1 aa 5] [3 cc 15 6 aa 5 1 aa 5] [4 dd 20 6 aa 5 1 aa 5] [5 ee 30 6 aa 5 1 aa 5] [6 aa 5 6 aa 5 1 aa 5] [7 bb 10 6 aa 5 1 aa 5] [8 cc 15 6 aa 5 1 aa 5] [9 dd 20 6 aa 5 1 aa 5] [10 ee 30 6 aa 5 1 aa 5] [1 aa 5 7 bb 10 1 aa 5] [2 bb 10 7 bb 10 1 aa 5] [3 cc 15 7 bb 10 1 aa 5] [4 dd 20 7 bb 10 1 aa 5] [5 ee 30 7 bb 10 1 aa 5] [6 aa 5 7 bb 10 1 aa 5] [7 bb 10 7 bb 10 1 aa 5] [8 cc 15 7 bb 10 1 aa 5] [9 dd 20 7 bb 10 1 aa 5] [10 ee 30 7 bb 10 1 aa 5] [1 aa 5 8 cc 15 1 aa 5] [2 bb 10 8 cc 15 1 aa 5] [3 cc 15 8 cc 15 1 aa 5] [4 dd 20 8 cc 15 1 aa 5] [5 ee 30 8 cc 15 1 aa 5] [6 aa 5 8 cc 15 1 aa 5] [7 bb 10 8 cc 15 1 aa 5] [8 cc 15 8 cc 15 1 aa 5] [9 dd 20 8 cc 15 1 aa 5] [10 ee 30 8 cc 15 1 aa 5] [1 aa 5 9 dd 20 1 aa 5] [2 bb 10 9 dd 20 1 aa 5] [3 cc 15 9 dd 20 1 aa 5] [4 dd 20 9 dd 20 1 aa 5] [5 ee 30 9 dd 20 1 aa 5] [6 aa 5 9 dd 20 1 aa 5] [7 bb 10 9 dd 20 1 aa 5] [8 cc 15 9 dd 20 1 aa 5] [9 dd 20 9 dd 20 1 aa 5] [10 ee 30 9 dd 20 1 aa 5] [1 aa 5 10 ee 30 1 aa 5] [2 bb 10 10 ee 30 1 aa 5] [3 cc 15 10 ee 30 1 aa 5] [4 dd 20 10 ee 30 1 aa 5] [5 ee 30 10 ee 30 1 aa 5] [6 aa 5 10 ee 30 1 aa 5] [7 bb 10 10 ee 30 1 aa 5] [8 cc 15 10 ee 30 1 aa 5] [9 dd 20 10 ee 30 1 aa 5] [10 ee 30 10 ee 30 1 aa 5] [1 aa 5 1 aa 5 2 bb 10] [2 bb 10 1 aa 5 2 bb 10] [3 cc 15 1 aa 5 2 bb 10] [4 dd 20 1 aa 5 2 bb 10] [5 ee 30 1 aa 5 2 bb 10] [6 aa 5 1 aa 5 2 bb 10] [7 bb 10 1 aa 5 2 bb 10] [8 cc 15 1 aa 5 2 bb 10] [9 dd 20 1 aa 5 2 bb 10] [10 ee 30 1 aa 5 2 bb 10] [1 aa 5 2 bb 10 2 bb 10] [2 bb 10 2 bb 10 2 bb 10] [3 cc 15 2 bb 10 2 bb 10] [4 dd 20 2 bb 10 2 bb 10] [5 ee 30 2 bb 10 2 bb 10] [6 aa 5 2 bb 10 2 bb 10] [7 bb 10 2 bb 10 2 bb 10] [8 cc 15 2 bb 10 2 bb 10] [9 dd 20 2 bb 10 2 bb 10] [10 ee 30 2 bb 10 2 bb 10] [1 aa 5 3 cc 15 2 bb 10] [2 bb 10 3 cc 15 2 bb 10] [3 cc 15 3 cc 15 2 bb 10] [4 dd 20 3 cc 15 2 bb 10] [5 ee 30 3 cc 15 2 bb 10] [6 aa 5 3 cc 15 2 bb 10] [7 bb 10 3 cc 15 2 bb 10] [8 cc 15 3 cc 15 2 bb 10] [9 dd 20 3 cc 15 2 bb 10] [10 ee 30 3 cc 15 2 bb 10] [1 aa 5 4 dd 20 2 bb 10] [2 bb 10 4 dd 20 2 bb 10] [3 cc 15 4 dd 20 2 bb 10] [4 dd 20 4 dd 20 2 bb 10] [5 ee 30 4 dd 20 2 bb 10] [6 aa 5 4 dd 20 2 bb 10] [7 bb 10 4 dd 20 2 bb 10] [8 cc 15 4 dd 20 2 bb 10] [9 dd 20 4 dd 20 2 bb 10] [10 ee 30 4 dd 20 2 bb 10] [1 aa 5 5 ee 30 2 bb 10] [2 bb 10 5 ee 30 2 bb 10] [3 cc 15 5 ee 30 2 bb 10] [4 dd 20 5 ee 30 2 bb 10] [5 ee 30 5 ee 30 2 bb 10] [6 aa 5 5 ee 30 2 bb 10] [7 bb 10 5 ee 30 2 bb 10] [8 cc 15 5 ee 30 2 bb 10] [9 dd 20 5 ee 30 2 bb 10] [10 ee 30 5 ee 30 2 bb 10] [1 aa 5 6 aa 5 2 bb 10] [2 bb 10 6 aa 5 2 bb 10] [3 cc 15 6 aa 5 2 bb 10] [4 dd 20 6 aa 5 2 bb 10] [5 ee 30 6 aa 5 2 bb 10] [6 aa 5 6 aa 5 2 bb 10] [7 bb 10 6 aa 5 2 bb 10] [8 cc 15 6 aa 5 2 bb 10] [9 dd 20 6 aa 5 2 bb 10] [10 ee 30 6 aa 5 2 bb 10] [1 aa 5 7 bb 10 2 bb 10] [2 bb 10 7 bb 10 2 bb 10] [3 cc 15 7 bb 10 2 bb 10] [4 dd 20 7 bb 10 2 bb 10] [5 ee 30 7 bb 10 2 bb 10] [6 aa 5 7 bb 10 2 bb 10] [7 bb 10 7 bb 10 2 bb 10] [8 cc 15 7 bb 10 2 bb 10] [9 dd 20 7 bb 10 2 bb 10] [10 ee 30 7 bb 10 2 bb 10] [1 aa 5 8 cc 15 2 bb 10] [2 bb 10 8 cc 15 2 bb 10] [3 cc 15 8 cc 15 2 bb 10] [4 dd 20 8 cc 15 2 bb 10] [5 ee 30 8 cc 15 2 bb 10] [6 aa 5 8 cc 15 2 bb 10] [7 bb 10 8 cc 15 2 bb 10] [8 cc 15 8 cc 15 2 bb 10] [9 dd 20 8 cc 15 2 bb 10] [10 ee 30 8 cc 15 2 bb 10] [1 aa 5 9 dd 20 2 bb 10] [2 bb 10 9 dd 20 2 bb 10] [3 cc 15 9 dd 20 2 bb 10] [4 dd 20 9 dd 20 2 bb 10] [5 ee 30 9 dd 20 2 bb 10] [6 aa 5 9 dd 20 2 bb 10] [7 bb 10 9 dd 20 2 bb 10] [8 cc 15 9 dd 20 2 bb 10] [9 dd 20 9 dd 20 2 bb 10] [10 ee 30 9 dd 20 2 bb 10] [1 aa 5 10 ee 30 2 bb 10] [2 bb 10 10 ee 30 2 bb 10] [3 cc 15 10 ee 30 2 bb 10] [4 dd 20 10 ee 30 2 bb 10] [5 ee 30 10 ee 30 2 bb 10] [6 aa 5 10 ee 30 2 bb 10] [7 bb 10 10 ee 30 2 bb 10] [8 cc 15 10 ee 30 2 bb 10] [9 dd 20 10 ee 30 2 bb 10] [10 ee 30 10 ee 30 2 bb 10] [1 aa 5 1 aa 5 3 cc 15] [2 bb 10 1 aa 5 3 cc 15] [3 cc 15 1 aa 5 3 cc 15] [4 dd 20 1 aa 5 3 cc 15] [5 ee 30 1 aa 5 3 cc 15] [6 aa 5 1 aa 5 3 cc 15] [7 bb 10 1 aa 5 3 cc 15] [8 cc 15 1 aa 5 3 cc 15] [9 dd 20 1 aa 5 3 cc 15] [10 ee 30 1 aa 5 3 cc 15] [1 aa 5 2 bb 10 3 cc 15] [2 bb 10 2 bb 10 3 cc 15] [3 cc 15 2 bb 10 3 cc 15] [4 dd 20 2 bb 10 3 cc 15] [5 ee 30 2 bb 10 3 cc 15] [6 aa 5 2 bb 10 3 cc 15] [7 bb 10 2 bb 10 3 cc 15] [8 cc 15 2 bb 10 3 cc 15] [9 dd 20 2 bb 10 3 cc 15] [10 ee 30 2 bb 10 3 cc 15] [1 aa 5 3 cc 15 3 cc 15] [2 bb 10 3 cc 15 3 cc 15] [3 cc 15 3 cc 15 3 cc 15] [4 dd 20 3 cc 15 3 cc 15] [5 ee 30 3 cc 15 3 cc 15] [6 aa 5 3 cc 15 3 cc 15] [7 bb 10 3 cc 15 3 cc 15] [8 cc 15 3 cc 15 3 cc 15] [9 dd 20 3 cc 15 3 cc 15] [10 ee 30 3 cc 15 3 cc 15] [1 aa 5 4 dd 20 3 cc 15] [2 bb 10 4 dd 20 3 cc 15] [3 cc 15 4 dd 20 3 cc 15] [4 dd 20 4 dd 20 3 cc 15] [5 ee 30 4 dd 20 3 cc 15] [6 aa 5 4 dd 20 3 cc 15] [7 bb 10 4 dd 20 3 cc 15] [8 cc 15 4 dd 20 3 cc 15] [9 dd 20 4 dd 20 3 cc 15] [10 ee 30 4 dd 20 3 cc 15] [1 aa 5 5 ee 30 3 cc 15] [2 bb 10 5 ee 30 3 cc 15] [3 cc 15 5 ee 30 3 cc 15] [4 dd 20 5 ee 30 3 cc 15] [5 ee 30 5 ee 30 3 cc 15] [6 aa 5 5 ee 30 3 cc 15] [7 bb 10 5 ee 30 3 cc 15] [8 cc 15 5 ee 30 3 cc 15] [9 dd 20 5 ee 30 3 cc 15] [10 ee 30 5 ee 30 3 cc 15] [1 aa 5 6 aa 5 3 cc 15] [2 bb 10 6 aa 5 3 cc 15] [3 cc 15 6 aa 5 3 cc 15] [4 dd 20 6 aa 5 3 cc 15] [5 ee 30 6 aa 5 3 cc 15] [6 aa 5 6 aa 5 3 cc 15] [7 bb 10 6 aa 5 3 cc 15] [8 cc 15 6 aa 5 3 cc 15] [9 dd 20 6 aa 5 3 cc 15] [10 ee 30 6 aa 5 3 cc 15] [1 aa 5 7 bb 10 3 cc 15] [2 bb 10 7 bb 10 3 cc 15] [3 cc 15 7 bb 10 3 cc 15] [4 dd 20 7 bb 10 3 cc 15] [5 ee 30 7 bb 10 3 cc 15] [6 aa 5 7 bb 10 3 cc 15] [7 bb 10 7 bb 10 3 cc 15] [8 cc 15 7 bb 10 3 cc 15] [9 dd 20 7 bb 10 3 cc 15] [10 ee 30 7 bb 10 3 cc 15] [1 aa 5 8 cc 15 3 cc 15] [2 bb 10 8 cc 15 3 cc 15] [3 cc 15 8 cc 15 3 cc 15] [4 dd 20 8 cc 15 3 cc 15] [5 ee 30 8 cc 15 3 cc 15] [6 aa 5 8 cc 15 3 cc 15] [7 bb 10 8 cc 15 3 cc 15] [8 cc 15 8 cc 15 3 cc 15] [9 dd 20 8 cc 15 3 cc 15] [10 ee 30 8 cc 15 3 cc 15] [1 aa 5 9 dd 20 3 cc 15] [2 bb 10 9 dd 20 3 cc 15] [3 cc 15 9 dd 20 3 cc 15] [4 dd 20 9 dd 20 3 cc 15] [5 ee 30 9 dd 20 3 cc 15] [6 aa 5 9 dd 20 3 cc 15] [7 bb 10 9 dd 20 3 cc 15] [8 cc 15 9 dd 20 3 cc 15] [9 dd 20 9 dd 20 3 cc 15] [10 ee 30 9 dd 20 3 cc 15] [1 aa 5 10 ee 30 3 cc 15] [2 bb 10 10 ee 30 3 cc 15] [3 cc 15 10 ee 30 3 cc 15] [4 dd 20 10 ee 30 3 cc 15] [5 ee 30 10 ee 30 3 cc 15] [6 aa 5 10 ee 30 3 cc 15] [7 bb 10 10 ee 30 3 cc 15] [8 cc 15 10 ee 30 3 cc 15] [9 dd 20 10 ee 30 3 cc 15] [10 ee 30 10 ee 30 3 cc 15] [1 aa 5 1 aa 5 4 dd 20] [2 bb 10 1 aa 5 4 dd 20] [3 cc 15 1 aa 5 4 dd 20] [4 dd 20 1 aa 5 4 dd 20] [5 ee 30 1 aa 5 4 dd 20] [6 aa 5 1 aa 5 4 dd 20] [7 bb 10 1 aa 5 4 dd 20] [8 cc 15 1 aa 5 4 dd 20] [9 dd 20 1 aa 5 4 dd 20] [10 ee 30 1 aa 5 4 dd 20] [1 aa 5 2 bb 10 4 dd 20] [2 bb 10 2 bb 10 4 dd 20] [3 cc 15 2 bb 10 4 dd 20] [4 dd 20 2 bb 10 4 dd 20] [5 ee 30 2 bb 10 4 dd 20] [6 aa 5 2 bb 10 4 dd 20] [7 bb 10 2 bb 10 4 dd 20] [8 cc 15 2 bb 10 4 dd 20] [9 dd 20 2 bb 10 4 dd 20] [10 ee 30 2 bb 10 4 dd 20] [1 aa 5 3 cc 15 4 dd 20] [2 bb 10 3 cc 15 4 dd 20] [3 cc 15 3 cc 15 4 dd 20] [4 dd 20 3 cc 15 4 dd 20] [5 ee 30 3 cc 15 4 dd 20] [6 aa 5 3 cc 15 4 dd 20] [7 bb 10 3 cc 15 4 dd 20] [8 cc 15 3 cc 15 4 dd 20] [9 dd 20 3 cc 15 4 dd 20] [10 ee 30 3 cc 15 4 dd 20] [1 aa 5 4 dd 20 4 dd 20] [2 bb 10 4 dd 20 4 dd 20] [3 cc 15 4 dd 20 4 dd 20] [4 dd 20 4 dd 20 4 dd 20] [5 ee 30 4 dd 20 4 dd 20] [6 aa 5 4 dd 20 4 dd 20] [7 bb 10 4 dd 20 4 dd 20] [8 cc 15 4 dd 20 4 dd 20] [9 dd 20 4 dd 20 4 dd 20] [10 ee 30 4 dd 20 4 dd 20] [1 aa 5 5 ee 30 4 dd 20] [2 bb 10 5 ee 30 4 dd 20] [3 cc 15 5 ee 30 4 dd 20] [4 dd 20 5 ee 30 4 dd 20] [5 ee 30 5 ee 30 4 dd 20] [6 aa 5 5 ee 30 4 dd 20] [7 bb 10 5 ee 30 4 dd 20] [8 cc 15 5 ee 30 4 dd 20] [9 dd 20 5 ee 30 4 dd 20] [10 ee 30 5 ee 30 4 dd 20] [1 aa 5 6 aa 5 4 dd 20] [2 bb 10 6 aa 5 4 dd 20] [3 cc 15 6 aa 5 4 dd 20] [4 dd 20 6 aa 5 4 dd 20] [5 ee 30 6 aa 5 4 dd 20] [6 aa 5 6 aa 5 4 dd 20] [7 bb 10 6 aa 5 4 dd 20] [8 cc 15 6 aa 5 4 dd 20] [9 dd 20 6 aa 5 4 dd 20] [10 ee 30 6 aa 5 4 dd 20] [1 aa 5 7 bb 10 4 dd 20] [2 bb 10 7 bb 10 4 dd 20] [3 cc 15 7 bb 10 4 dd 20] [4 dd 20 7 bb 10 4 dd 20] [5 ee 30 7 bb 10 4 dd 20] [6 aa 5 7 bb 10 4 dd 20] [7 bb 10 7 bb 10 4 dd 20] [8 cc 15 7 bb 10 4 dd 20] [9 dd 20 7 bb 10 4 dd 20] [10 ee 30 7 bb 10 4 dd 20] [1 aa 5 8 cc 15 4 dd 20] [2 bb 10 8 cc 15 4 dd 20] [3 cc 15 8 cc 15 4 dd 20] [4 dd 20 8 cc 15 4 dd 20] [5 ee 30 8 cc 15 4 dd 20] [6 aa 5 8 cc 15 4 dd 20] [7 bb 10 8 cc 15 4 dd 20] [8 cc 15 8 cc 15 4 dd 20] [9 dd 20 8 cc 15 4 dd 20] [10 ee 30 8 cc 15 4 dd 20] [1 aa 5 9 dd 20 4 dd 20] [2 bb 10 9 dd 20 4 dd 20] [3 cc 15 9 dd 20 4 dd 20] [4 dd 20 9 dd 20 4 dd 20] [5 ee 30 9 dd 20 4 dd 20] [6 aa 5 9 dd 20 4 dd 20] [7 bb 10 9 dd 20 4 dd 20] [8 cc 15 9 dd 20 4 dd 20] [9 dd 20 9 dd 20 4 dd 20] [10 ee 30 9 dd 20 4 dd 20] [1 aa 5 10 ee 30 4 dd 20] [2 bb 10 10 ee 30 4 dd 20] [3 cc 15 10 ee 30 4 dd 20] [4 dd 20 10 ee 30 4 dd 20] [5 ee 30 10 ee 30 4 dd 20] [6 aa 5 10 ee 30 4 dd 20] [7 bb 10 10 ee 30 4 dd 20] [8 cc 15 10 ee 30 4 dd 20] [9 dd 20 10 ee 30 4 dd 20] [10 ee 30 10 ee 30 4 dd 20] [1 aa 5 1 aa 5 5 ee 30] [2 bb 10 1 aa 5 5 ee 30] [3 cc 15 1 aa 5 5 ee 30] [4 dd 20 1 aa 5 5 ee 30] [5 ee 30 1 aa 5 5 ee 30] [6 aa 5 1 aa 5 5 ee 30] [7 bb 10 1 aa 5 5 ee 30] [8 cc 15 1 aa 5 5 ee 30] [9 dd 20 1 aa 5 5 ee 30] [10 ee 30 1 aa 5 5 ee 30] [1 aa 5 2 bb 10 5 ee 30] [2 bb 10 2 bb 10 5 ee 30] [3 cc 15 2 bb 10 5 ee 30] [4 dd 20 2 bb 10 5 ee 30] [5 ee 30 2 bb 10 5 ee 30] [6 aa 5 2 bb 10 5 ee 30] [7 bb 10 2 bb 10 5 ee 30] [8 cc 15 2 bb 10 5 ee 30] [9 dd 20 2 bb 10 5 ee 30] [10 ee 30 2 bb 10 5 ee 30] [1 aa 5 3 cc 15 5 ee 30] [2 bb 10 3 cc 15 5 ee 30] [3 cc 15 3 cc 15 5 ee 30] [4 dd 20 3 cc 15 5 ee 30] [5 ee 30 3 cc 15 5 ee 30] [6 aa 5 3 cc 15 5 ee 30] [7 bb 10 3 cc 15 5 ee 30] [8 cc 15 3 cc 15 5 ee 30] [9 dd 20 3 cc 15 5 ee 30] [10 ee 30 3 cc 15 5 ee 30] [1 aa 5 4 dd 20 5 ee 30] [2 bb 10 4 dd 20 5 ee 30] [3 cc 15 4 dd 20 5 ee 30] [4 dd 20 4 dd 20 5 ee 30] [5 ee 30 4 dd 20 5 ee 30] [6 aa 5 4 dd 20 5 ee 30] [7 bb 10 4 dd 20 5 ee 30] [8 cc 15 4 dd 20 5 ee 30] [9 dd 20 4 dd 20 5 ee 30] [10 ee 30 4 dd 20 5 ee 30] [1 aa 5 5 ee 30 5 ee 30] [2 bb 10 5 ee 30 5 ee 30] [3 cc 15 5 ee 30 5 ee 30] [4 dd 20 5 ee 30 5 ee 30] [5 ee 30 5 ee 30 5 ee 30] [6 aa 5 5 ee 30 5 ee 30] [7 bb 10 5 ee 30 5 ee 30] [8 cc 15 5 ee 30 5 ee 30] [9 dd 20 5 ee 30 5 ee 30] [10 ee 30 5 ee 30 5 ee 30] [1 aa 5 6 aa 5 5 ee 30] [2 bb 10 6 aa 5 5 ee 30] [3 cc 15 6 aa 5 5 ee 30] [4 dd 20 6 aa 5 5 ee 30] [5 ee 30 6 aa 5 5 ee 30] [6 aa 5 6 aa 5 5 ee 30] [7 bb 10 6 aa 5 5 ee 30] [8 cc 15 6 aa 5 5 ee 30] [9 dd 20 6 aa 5 5 ee 30] [10 ee 30 6 aa 5 5 ee 30] [1 aa 5 7 bb 10 5 ee 30] [2 bb 10 7 bb 10 5 ee 30] [3 cc 15 7 bb 10 5 ee 30] [4 dd 20 7 bb 10 5 ee 30] [5 ee 30 7 bb 10 5 ee 30] [6 aa 5 7 bb 10 5 ee 30] [7 bb 10 7 bb 10 5 ee 30] [8 cc 15 7 bb 10 5 ee 30] [9 dd 20 7 bb 10 5 ee 30] [10 ee 30 7 bb 10 5 ee 30] [1 aa 5 8 cc 15 5 ee 30] [2 bb 10 8 cc 15 5 ee 30] [3 cc 15 8 cc 15 5 ee 30] [4 dd 20 8 cc 15 5 ee 30] [5 ee 30 8 cc 15 5 ee 30] [6 aa 5 8 cc 15 5 ee 30] [7 bb 10 8 cc 15 5 ee 30] [8 cc 15 8 cc 15 5 ee 30] [9 dd 20 8 cc 15 5 ee 30] [10 ee 30 8 cc 15 5 ee 30] [1 aa 5 9 dd 20 5 ee 30] [2 bb 10 9 dd 20 5 ee 30] [3 cc 15 9 dd 20 5 ee 30] [4 dd 20 9 dd 20 5 ee 30] [5 ee 30 9 dd 20 5 ee 30] [6 aa 5 9 dd 20 5 ee 30] [7 bb 10 9 dd 20 5 ee 30] [8 cc 15 9 dd 20 5 ee 30] [9 dd 20 9 dd 20 5 ee 30] [10 ee 30 9 dd 20 5 ee 30] [1 aa 5 10 ee 30 5 ee 30] [2 bb 10 10 ee 30 5 ee 30] [3 cc 15 10 ee 30 5 ee 30] [4 dd 20 10 ee 30 5 ee 30] [5 ee 30 10 ee 30 5 ee 30] [6 aa 5 10 ee 30 5 ee 30] [7 bb 10 10 ee 30 5 ee 30] [8 cc 15 10 ee 30 5 ee 30] [9 dd 20 10 ee 30 5 ee 30] [10 ee 30 10 ee 30 5 ee 30] [1 aa 5 1 aa 5 6 aa 5] [2 bb 10 1 aa 5 6 aa 5] [3 cc 15 1 aa 5 6 aa 5] [4 dd 20 1 aa 5 6 aa 5] [5 ee 30 1 aa 5 6 aa 5] [6 aa 5 1 aa 5 6 aa 5] [7 bb 10 1 aa 5 6 aa 5] [8 cc 15 1 aa 5 6 aa 5] [9 dd 20 1 aa 5 6 aa 5] [10 ee 30 1 aa 5 6 aa 5] [1 aa 5 2 bb 10 6 aa 5] [2 bb 10 2 bb 10 6 aa 5] [3 cc 15 2 bb 10 6 aa 5] [4 dd 20 2 bb 10 6 aa 5] [5 ee 30 2 bb 10 6 aa 5] [6 aa 5 2 bb 10 6 aa 5] [7 bb 10 2 bb 10 6 aa 5] [8 cc 15 2 bb 10 6 aa 5] [9 dd 20 2 bb 10 6 aa 5] [10 ee 30 2 bb 10 6 aa 5] [1 aa 5 3 cc 15 6 aa 5] [2 bb 10 3 cc 15 6 aa 5] [3 cc 15 3 cc 15 6 aa 5] [4 dd 20 3 cc 15 6 aa 5] [5 ee 30 3 cc 15 6 aa 5] [6 aa 5 3 cc 15 6 aa 5] [7 bb 10 3 cc 15 6 aa 5] [8 cc 15 3 cc 15 6 aa 5] [9 dd 20 3 cc 15 6 aa 5] [10 ee 30 3 cc 15 6 aa 5] [1 aa 5 4 dd 20 6 aa 5] [2 bb 10 4 dd 20 6 aa 5] [3 cc 15 4 dd 20 6 aa 5] [4 dd 20 4 dd 20 6 aa 5] [5 ee 30 4 dd 20 6 aa 5] [6 aa 5 4 dd 20 6 aa 5] [7 bb 10 4 dd 20 6 aa 5] [8 cc 15 4 dd 20 6 aa 5] [9 dd 20 4 dd 20 6 aa 5] [10 ee 30 4 dd 20 6 aa 5] [1 aa 5 5 ee 30 6 aa 5] [2 bb 10 5 ee 30 6 aa 5] [3 cc 15 5 ee 30 6 aa 5] [4 dd 20 5 ee 30 6 aa 5] [5 ee 30 5 ee 30 6 aa 5] [6 aa 5 5 ee 30 6 aa 5] [7 bb 10 5 ee 30 6 aa 5] [8 cc 15 5 ee 30 6 aa 5] [9 dd 20 5 ee 30 6 aa 5] [10 ee 30 5 ee 30 6 aa 5] [1 aa 5 6 aa 5 6 aa 5] [2 bb 10 6 aa 5 6 aa 5] [3 cc 15 6 aa 5 6 aa 5] [4 dd 20 6 aa 5 6 aa 5] [5 ee 30 6 aa 5 6 aa 5] [6 aa 5 6 aa 5 6 aa 5] [7 bb 10 6 aa 5 6 aa 5] [8 cc 15 6 aa 5 6 aa 5] [9 dd 20 6 aa 5 6 aa 5] [10 ee 30 6 aa 5 6 aa 5] [1 aa 5 7 bb 10 6 aa 5] [2 bb 10 7 bb 10 6 aa 5] [3 cc 15 7 bb 10 6 aa 5] [4 dd 20 7 bb 10 6 aa 5] [5 ee 30 7 bb 10 6 aa 5] [6 aa 5 7 bb 10 6 aa 5] [7 bb 10 7 bb 10 6 aa 5] [8 cc 15 7 bb 10 6 aa 5] [9 dd 20 7 bb 10 6 aa 5] [10 ee 30 7 bb 10 6 aa 5] [1 aa 5 8 cc 15 6 aa 5] [2 bb 10 8 cc 15 6 aa 5] [3 cc 15 8 cc 15 6 aa 5] [4 dd 20 8 cc 15 6 aa 5] [5 ee 30 8 cc 15 6 aa 5] [6 aa 5 8 cc 15 6 aa 5] [7 bb 10 8 cc 15 6 aa 5] [8 cc 15 8 cc 15 6 aa 5] [9 dd 20 8 cc 15 6 aa 5] [10 ee 30 8 cc 15 6 aa 5] [1 aa 5 9 dd 20 6 aa 5] [2 bb 10 9 dd 20 6 aa 5] [3 cc 15 9 dd 20 6 aa 5] [4 dd 20 9 dd 20 6 aa 5] [5 ee 30 9 dd 20 6 aa 5] [6 aa 5 9 dd 20 6 aa 5] [7 bb 10 9 dd 20 6 aa 5] [8 cc 15 9 dd 20 6 aa 5] [9 dd 20 9 dd 20 6 aa 5] [10 ee 30 9 dd 20 6 aa 5] [1 aa 5 10 ee 30 6 aa 5] [2 bb 10 10 ee 30 6 aa 5] [3 cc 15 10 ee 30 6 aa 5] [4 dd 20 10 ee 30 6 aa 5] [5 ee 30 10 ee 30 6 aa 5] [6 aa 5 10 ee 30 6 aa 5] [7 bb 10 10 ee 30 6 aa 5] [8 cc 15 10 ee 30 6 aa 5] [9 dd 20 10 ee 30 6 aa 5] [10 ee 30 10 ee 30 6 aa 5] [1 aa 5 1 aa 5 7 bb 10] [2 bb 10 1 aa 5 7 bb 10] [3 cc 15 1 aa 5 7 bb 10] [4 dd 20 1 aa 5 7 bb 10] [5 ee 30 1 aa 5 7 bb 10] [6 aa 5 1 aa 5 7 bb 10] [7 bb 10 1 aa 5 7 bb 10] [8 cc 15 1 aa 5 7 bb 10] [9 dd 20 1 aa 5 7 bb 10] [10 ee 30 1 aa 5 7 bb 10] [1 aa 5 2 bb 10 7 bb 10] [2 bb 10 2 bb 10 7 bb 10] [3 cc 15 2 bb 10 7 bb 10] [4 dd 20 2 bb 10 7 bb 10] [5 ee 30 2 bb 10 7 bb 10] [6 aa 5 2 bb 10 7 bb 10] [7 bb 10 2 bb 10 7 bb 10] [8 cc 15 2 bb 10 7 bb 10] [9 dd 20 2 bb 10 7 bb 10] [10 ee 30 2 bb 10 7 bb 10] [1 aa 5 3 cc 15 7 bb 10] [2 bb 10 3 cc 15 7 bb 10] [3 cc 15 3 cc 15 7 bb 10] [4 dd 20 3 cc 15 7 bb 10] [5 ee 30 3 cc 15 7 bb 10] [6 aa 5 3 cc 15 7 bb 10] [7 bb 10 3 cc 15 7 bb 10] [8 cc 15 3 cc 15 7 bb 10] [9 dd 20 3 cc 15 7 bb 10] [10 ee 30 3 cc 15 7 bb 10] [1 aa 5 4 dd 20 7 bb 10] [2 bb 10 4 dd 20 7 bb 10] [3 cc 15 4 dd 20 7 bb 10] [4 dd 20 4 dd 20 7 bb 10] [5 ee 30 4 dd 20 7 bb 10] [6 aa 5 4 dd 20 7 bb 10] [7 bb 10 4 dd 20 7 bb 10] [8 cc 15 4 dd 20 7 bb 10] [9 dd 20 4 dd 20 7 bb 10] [10 ee 30 4 dd 20 7 bb 10] [1 aa 5 5 ee 30 7 bb 10] [2 bb 10 5 ee 30 7 bb 10] [3 cc 15 5 ee 30 7 bb 10] [4 dd 20 5 ee 30 7 bb 10] [5 ee 30 5 ee 30 7 bb 10] [6 aa 5 5 ee 30 7 bb 10] [7 bb 10 5 ee 30 7 bb 10] [8 cc 15 5 ee 30 7 bb 10] [9 dd 20 5 ee 30 7 bb 10] [10 ee 30 5 ee 30 7 bb 10] [1 aa 5 6 aa 5 7 bb 10] [2 bb 10 6 aa 5 7 bb 10] [3 cc 15 6 aa 5 7 bb 10] [4 dd 20 6 aa 5 7 bb 10] [5 ee 30 6 aa 5 7 bb 10] [6 aa 5 6 aa 5 7 bb 10] [7 bb 10 6 aa 5 7 bb 10] [8 cc 15 6 aa 5 7 bb 10] [9 dd 20 6 aa 5 7 bb 10] [10 ee 30 6 aa 5 7 bb 10] [1 aa 5 7 bb 10 7 bb 10] [2 bb 10 7 bb 10 7 bb 10] [3 cc 15 7 bb 10 7 bb 10] [4 dd 20 7 bb 10 7 bb 10] [5 ee 30 7 bb 10 7 bb 10] [6 aa 5 7 bb 10 7 bb 10] [7 bb 10 7 bb 10 7 bb 10] [8 cc 15 7 bb 10 7 bb 10] [9 dd 20 7 bb 10 7 bb 10] [10 ee 30 7 bb 10 7 bb 10] [1 aa 5 8 cc 15 7 bb 10] [2 bb 10 8 cc 15 7 bb 10] [3 cc 15 8 cc 15 7 bb 10] [4 dd 20 8 cc 15 7 bb 10] [5 ee 30 8 cc 15 7 bb 10] [6 aa 5 8 cc 15 7 bb 10] [7 bb 10 8 cc 15 7 bb 10] [8 cc 15 8 cc 15 7 bb 10] [9 dd 20 8 cc 15 7 bb 10] [10 ee 30 8 cc 15 7 bb 10] [1 aa 5 9 dd 20 7 bb 10] [2 bb 10 9 dd 20 7 bb 10] [3 cc 15 9 dd 20 7 bb 10] [4 dd 20 9 dd 20 7 bb 10] [5 ee 30 9 dd 20 7 bb 10] [6 aa 5 9 dd 20 7 bb 10] [7 bb 10 9 dd 20 7 bb 10] [8 cc 15 9 dd 20 7 bb 10] [9 dd 20 9 dd 20 7 bb 10] [10 ee 30 9 dd 20 7 bb 10] [1 aa 5 10 ee 30 7 bb 10] [2 bb 10 10 ee 30 7 bb 10] [3 cc 15 10 ee 30 7 bb 10] [4 dd 20 10 ee 30 7 bb 10] [5 ee 30 10 ee 30 7 bb 10] [6 aa 5 10 ee 30 7 bb 10] [7 bb 10 10 ee 30 7 bb 10] [8 cc 15 10 ee 30 7 bb 10] [9 dd 20 10 ee 30 7 bb 10] [10 ee 30 10 ee 30 7 bb 10] [1 aa 5 1 aa 5 8 cc 15] [2 bb 10 1 aa 5 8 cc 15] [3 cc 15 1 aa 5 8 cc 15] [4 dd 20 1 aa 5 8 cc 15] [5 ee 30 1 aa 5 8 cc 15] [6 aa 5 1 aa 5 8 cc 15] [7 bb 10 1 aa 5 8 cc 15] [8 cc 15 1 aa 5 8 cc 15] [9 dd 20 1 aa 5 8 cc 15] [10 ee 30 1 aa 5 8 cc 15] [1 aa 5 2 bb 10 8 cc 15] [2 bb 10 2 bb 10 8 cc 15] [3 cc 15 2 bb 10 8 cc 15] [4 dd 20 2 bb 10 8 cc 15] [5 ee 30 2 bb 10 8 cc 15] [6 aa 5 2 bb 10 8 cc 15] [7 bb 10 2 bb 10 8 cc 15] [8 cc 15 2 bb 10 8 cc 15] [9 dd 20 2 bb 10 8 cc 15] [10 ee 30 2 bb 10 8 cc 15] [1 aa 5 3 cc 15 8 cc 15] [2 bb 10 3 cc 15 8 cc 15] [3 cc 15 3 cc 15 8 cc 15] [4 dd 20 3 cc 15 8 cc 15] [5 ee 30 3 cc 15 8 cc 15] [6 aa 5 3 cc 15 8 cc 15] [7 bb 10 3 cc 15 8 cc 15] [8 cc 15 3 cc 15 8 cc 15] [9 dd 20 3 cc 15 8 cc 15] [10 ee 30 3 cc 15 8 cc 15] [1 aa 5 4 dd 20 8 cc 15] [2 bb 10 4 dd 20 8 cc 15] [3 cc 15 4 dd 20 8 cc 15] [4 dd 20 4 dd 20 8 cc 15] [5 ee 30 4 dd 20 8 cc 15] [6 aa 5 4 dd 20 8 cc 15] [7 bb 10 4 dd 20 8 cc 15] [8 cc 15 4 dd 20 8 cc 15] [9 dd 20 4 dd 20 8 cc 15] [10 ee 30 4 dd 20 8 cc 15] [1 aa 5 5 ee 30 8 cc 15] [2 bb 10 5 ee 30 8 cc 15] [3 cc 15 5 ee 30 8 cc 15] [4 dd 20 5 ee 30 8 cc 15] [5 ee 30 5 ee 30 8 cc 15] [6 aa 5 5 ee 30 8 cc 15] [7 bb 10 5 ee 30 8 cc 15] [8 cc 15 5 ee 30 8 cc 15] [9 dd 20 5 ee 30 8 cc 15] [10 ee 30 5 ee 30 8 cc 15] [1 aa 5 6 aa 5 8 cc 15] [2 bb 10 6 aa 5 8 cc 15] [3 cc 15 6 aa 5 8 cc 15] [4 dd 20 6 aa 5 8 cc 15] [5 ee 30 6 aa 5 8 cc 15] [6 aa 5 6 aa 5 8 cc 15] [7 bb 10 6 aa 5 8 cc 15] [8 cc 15 6 aa 5 8 cc 15] [9 dd 20 6 aa 5 8 cc 15] [10 ee 30 6 aa 5 8 cc 15] [1 aa 5 7 bb 10 8 cc 15] [2 bb 10 7 bb 10 8 cc 15] [3 cc 15 7 bb 10 8 cc 15] [4 dd 20 7 bb 10 8 cc 15] [5 ee 30 7 bb 10 8 cc 15] [6 aa 5 7 bb 10 8 cc 15] [7 bb 10 7 bb 10 8 cc 15] [8 cc 15 7 bb 10 8 cc 15] [9 dd 20 7 bb 10 8 cc 15] [10 ee 30 7 bb 10 8 cc 15] [1 aa 5 8 cc 15 8 cc 15] [2 bb 10 8 cc 15 8 cc 15] [3 cc 15 8 cc 15 8 cc 15] [4 dd 20 8 cc 15 8 cc 15] [5 ee 30 8 cc 15 8 cc 15] [6 aa 5 8 cc 15 8 cc 15] [7 bb 10 8 cc 15 8 cc 15] [8 cc 15 8 cc 15 8 cc 15] [9 dd 20 8 cc 15 8 cc 15] [10 ee 30 8 cc 15 8 cc 15] [1 aa 5 9 dd 20 8 cc 15] [2 bb 10 9 dd 20 8 cc 15] [3 cc 15 9 dd 20 8 cc 15] [4 dd 20 9 dd 20 8 cc 15] [5 ee 30 9 dd 20 8 cc 15] [6 aa 5 9 dd 20 8 cc 15] [7 bb 10 9 dd 20 8 cc 15] [8 cc 15 9 dd 20 8 cc 15] [9 dd 20 9 dd 20 8 cc 15] [10 ee 30 9 dd 20 8 cc 15] [1 aa 5 10 ee 30 8 cc 15] [2 bb 10 10 ee 30 8 cc 15] [3 cc 15 10 ee 30 8 cc 15] [4 dd 20 10 ee 30 8 cc 15] [5 ee 30 10 ee 30 8 cc 15] [6 aa 5 10 ee 30 8 cc 15] [7 bb 10 10 ee 30 8 cc 15] [8 cc 15 10 ee 30 8 cc 15] [9 dd 20 10 ee 30 8 cc 15] [10 ee 30 10 ee 30 8 cc 15] [1 aa 5 1 aa 5 9 dd 20] [2 bb 10 1 aa 5 9 dd 20] [3 cc 15 1 aa 5 9 dd 20] [4 dd 20 1 aa 5 9 dd 20] [5 ee 30 1 aa 5 9 dd 20] [6 aa 5 1 aa 5 9 dd 20] [7 bb 10 1 aa 5 9 dd 20] [8 cc 15 1 aa 5 9 dd 20] [9 dd 20 1 aa 5 9 dd 20] [10 ee 30 1 aa 5 9 dd 20] [1 aa 5 2 bb 10 9 dd 20] [2 bb 10 2 bb 10 9 dd 20] [3 cc 15 2 bb 10 9 dd 20] [4 dd 20 2 bb 10 9 dd 20] [5 ee 30 2 bb 10 9 dd 20] [6 aa 5 2 bb 10 9 dd 20] [7 bb 10 2 bb 10 9 dd 20] [8 cc 15 2 bb 10 9 dd 20] [9 dd 20 2 bb 10 9 dd 20] [10 ee 30 2 bb 10 9 dd 20] [1 aa 5 3 cc 15 9 dd 20] [2 bb 10 3 cc 15 9 dd 20] [3 cc 15 3 cc 15 9 dd 20] [4 dd 20 3 cc 15 9 dd 20] [5 ee 30 3 cc 15 9 dd 20] [6 aa 5 3 cc 15 9 dd 20] [7 bb 10 3 cc 15 9 dd 20] [8 cc 15 3 cc 15 9 dd 20] [9 dd 20 3 cc 15 9 dd 20] [10 ee 30 3 cc 15 9 dd 20] [1 aa 5 4 dd 20 9 dd 20] [2 bb 10 4 dd 20 9 dd 20] [3 cc 15 4 dd 20 9 dd 20] [4 dd 20 4 dd 20 9 dd 20] [5 ee 30 4 dd 20 9 dd 20] [6 aa 5 4 dd 20 9 dd 20] [7 bb 10 4 dd 20 9 dd 20] [8 cc 15 4 dd 20 9 dd 20] [9 dd 20 4 dd 20 9 dd 20] [10 ee 30 4 dd 20 9 dd 20] [1 aa 5 5 ee 30 9 dd 20] [2 bb 10 5 ee 30 9 dd 20] [3 cc 15 5 ee 30 9 dd 20] [4 dd 20 5 ee 30 9 dd 20] [5 ee 30 5 ee 30 9 dd 20] [6 aa 5 5 ee 30 9 dd 20] [7 bb 10 5 ee 30 9 dd 20] [8 cc 15 5 ee 30 9 dd 20] [9 dd 20 5 ee 30 9 dd 20] [10 ee 30 5 ee 30 9 dd 20] [1 aa 5 6 aa 5 9 dd 20] [2 bb 10 6 aa 5 9 dd 20] [3 cc 15 6 aa 5 9 dd 20] [4 dd 20 6 aa 5 9 dd 20] [5 ee 30 6 aa 5 9 dd 20] [6 aa 5 6 aa 5 9 dd 20] [7 bb 10 6 aa 5 9 dd 20] [8 cc 15 6 aa 5 9 dd 20] [9 dd 20 6 aa 5 9 dd 20] [10 ee 30 6 aa 5 9 dd 20] [1 aa 5 7 bb 10 9 dd 20] [2 bb 10 7 bb 10 9 dd 20] [3 cc 15 7 bb 10 9 dd 20] [4 dd 20 7 bb 10 9 dd 20] [5 ee 30 7 bb 10 9 dd 20] [6 aa 5 7 bb 10 9 dd 20] [7 bb 10 7 bb 10 9 dd 20] [8 cc 15 7 bb 10 9 dd 20] [9 dd 20 7 bb 10 9 dd 20] [10 ee 30 7 bb 10 9 dd 20] [1 aa 5 8 cc 15 9 dd 20] [2 bb 10 8 cc 15 9 dd 20] [3 cc 15 8 cc 15 9 dd 20] [4 dd 20 8 cc 15 9 dd 20] [5 ee 30 8 cc 15 9 dd 20] [6 aa 5 8 cc 15 9 dd 20] [7 bb 10 8 cc 15 9 dd 20] [8 cc 15 8 cc 15 9 dd 20] [9 dd 20 8 cc 15 9 dd 20] [10 ee 30 8 cc 15 9 dd 20] [1 aa 5 9 dd 20 9 dd 20] [2 bb 10 9 dd 20 9 dd 20] [3 cc 15 9 dd 20 9 dd 20] [4 dd 20 9 dd 20 9 dd 20] [5 ee 30 9 dd 20 9 dd 20] [6 aa 5 9 dd 20 9 dd 20] [7 bb 10 9 dd 20 9 dd 20] [8 cc 15 9 dd 20 9 dd 20] [9 dd 20 9 dd 20 9 dd 20] [10 ee 30 9 dd 20 9 dd 20] [1 aa 5 10 ee 30 9 dd 20] [2 bb 10 10 ee 30 9 dd 20] [3 cc 15 10 ee 30 9 dd 20] [4 dd 20 10 ee 30 9 dd 20] [5 ee 30 10 ee 30 9 dd 20] [6 aa 5 10 ee 30 9 dd 20] [7 bb 10 10 ee 30 9 dd 20] [8 cc 15 10 ee 30 9 dd 20] [9 dd 20 10 ee 30 9 dd 20] [10 ee 30 10 ee 30 9 dd 20] [1 aa 5 1 aa 5 10 ee 30] [2 bb 10 1 aa 5 10 ee 30] [3 cc 15 1 aa 5 10 ee 30] [4 dd 20 1 aa 5 10 ee 30] [5 ee 30 1 aa 5 10 ee 30] [6 aa 5 1 aa 5 10 ee 30] [7 bb 10 1 aa 5 10 ee 30] [8 cc 15 1 aa 5 10 ee 30] [9 dd 20 1 aa 5 10 ee 30] [10 ee 30 1 aa 5 10 ee 30] [1 aa 5 2 bb 10 10 ee 30] [2 bb 10 2 bb 10 10 ee 30] [3 cc 15 2 bb 10 10 ee 30] [4 dd 20 2 bb 10 10 ee 30] [5 ee 30 2 bb 10 10 ee 30] [6 aa 5 2 bb 10 10 ee 30] [7 bb 10 2 bb 10 10 ee 30] [8 cc 15 2 bb 10 10 ee 30] [9 dd 20 2 bb 10 10 ee 30] [10 ee 30 2 bb 10 10 ee 30] [1 aa 5 3 cc 15 10 ee 30] [2 bb 10 3 cc 15 10 ee 30] [3 cc 15 3 cc 15 10 ee 30] [4 dd 20 3 cc 15 10 ee 30] [5 ee 30 3 cc 15 10 ee 30] [6 aa 5 3 cc 15 10 ee 30] [7 bb 10 3 cc 15 10 ee 30] [8 cc 15 3 cc 15 10 ee 30] [9 dd 20 3 cc 15 10 ee 30] [10 ee 30 3 cc 15 10 ee 30] [1 aa 5 4 dd 20 10 ee 30] [2 bb 10 4 dd 20 10 ee 30] [3 cc 15 4 dd 20 10 ee 30] [4 dd 20 4 dd 20 10 ee 30] [5 ee 30 4 dd 20 10 ee 30] [6 aa 5 4 dd 20 10 ee 30] [7 bb 10 4 dd 20 10 ee 30] [8 cc 15 4 dd 20 10 ee 30] [9 dd 20 4 dd 20 10 ee 30] [10 ee 30 4 dd 20 10 ee 30] [1 aa 5 5 ee 30 10 ee 30] [2 bb 10 5 ee 30 10 ee 30] [3 cc 15 5 ee 30 10 ee 30] [4 dd 20 5 ee 30 10 ee 30] [5 ee 30 5 ee 30 10 ee 30] [6 aa 5 5 ee 30 10 ee 30] [7 bb 10 5 ee 30 10 ee 30] [8 cc 15 5 ee 30 10 ee 30] [9 dd 20 5 ee 30 10 ee 30] [10 ee 30 5 ee 30 10 ee 30] [1 aa 5 6 aa 5 10 ee 30] [2 bb 10 6 aa 5 10 ee 30] [3 cc 15 6 aa 5 10 ee 30] [4 dd 20 6 aa 5 10 ee 30] [5 ee 30 6 aa 5 10 ee 30] [6 aa 5 6 aa 5 10 ee 30] [7 bb 10 6 aa 5 10 ee 30] [8 cc 15 6 aa 5 10 ee 30] [9 dd 20 6 aa 5 10 ee 30] [10 ee 30 6 aa 5 10 ee 30] [1 aa 5 7 bb 10 10 ee 30] [2 bb 10 7 bb 10 10 ee 30] [3 cc 15 7 bb 10 10 ee 30] [4 dd 20 7 bb 10 10 ee 30] [5 ee 30 7 bb 10 10 ee 30] [6 aa 5 7 bb 10 10 ee 30] [7 bb 10 7 bb 10 10 ee 30] [8 cc 15 7 bb 10 10 ee 30] [9 dd 20 7 bb 10 10 ee 30] [10 ee 30 7 bb 10 10 ee 30] [1 aa 5 8 cc 15 10 ee 30] [2 bb 10 8 cc 15 10 ee 30] [3 cc 15 8 cc 15 10 ee 30] [4 dd 20 8 cc 15 10 ee 30] [5 ee 30 8 cc 15 10 ee 30] [6 aa 5 8 cc 15 10 ee 30] [7 bb 10 8 cc 15 10 ee 30] [8 cc 15 8 cc 15 10 ee 30] [9 dd 20 8 cc 15 10 ee 30] [10 ee 30 8 cc 15 10 ee 30] [1 aa 5 9 dd 20 10 ee 30] [2 bb 10 9 dd 20 10 ee 30] [3 cc 15 9 dd 20 10 ee 30] [4 dd 20 9 dd 20 10 ee 30] [5 ee 30 9 dd 20 10 ee 30] [6 aa 5 9 dd 20 10 ee 30] [7 bb 10 9 dd 20 10 ee 30] [8 cc 15 9 dd 20 10 ee 30] [9 dd 20 9 dd 20 10 ee 30] [10 ee 30 9 dd 20 10 ee 30] [1 aa 5 10 ee 30 10 ee 30] [2 bb 10 10 ee 30 10 ee 30] [3 cc 15 10 ee 30 10 ee 30] [4 dd 20 10 ee 30 10 ee 30] [5 ee 30 10 ee 30 10 ee 30] [6 aa 5 10 ee 30 10 ee 30] [7 bb 10 10 ee 30 10 ee 30] [8 cc 15 10 ee 30 10 ee 30] [9 dd 20 10 ee 30 10 ee 30] [10 ee 30 10 ee 30 10 ee 30]]

sql:select .78+123;
mysqlRes:[[123.78]]
gaeaRes:[[123.78]]

sql:select .78+.21;
mysqlRes:[[0.99]]
gaeaRes:[[0.99]]

sql:select .78-123;
mysqlRes:[[-122.22]]
gaeaRes:[[-122.22]]

sql:select .78-.21;
mysqlRes:[[0.57]]
gaeaRes:[[0.57]]

sql:select .78--123;
mysqlRes:[[123.78]]
gaeaRes:[[123.78]]

sql:select .78*123;
mysqlRes:[[95.94]]
gaeaRes:[[95.94]]

sql:select .78*.21;
mysqlRes:[[0.1638]]
gaeaRes:[[0.1638]]

sql:select .78/123;
mysqlRes:[[0.006341]]
gaeaRes:[[0.006341]]

sql:select .78/.21;
mysqlRes:[[3.714286]]
gaeaRes:[[3.714286]]

sql:select .78,123;
mysqlRes:[[0.78 123]]
gaeaRes:[[0.78 123]]

sql:select .78,.21;
mysqlRes:[[0.78 0.21]]
gaeaRes:[[0.78 0.21]]

sql:select .78 , 123;
mysqlRes:[[0.78 123]]
gaeaRes:[[0.78 123]]

sql:select .78'123';
mysqlRes:[[0.78]]
gaeaRes:[[0.78]]

sql:select .78`123`;
mysqlRes:[[0.78]]
gaeaRes:[[0.78]]

sql:select .78"123";
mysqlRes:[[0.78]]
gaeaRes:[[0.78]]

sql:select x'0a', X'11', 0x11;
mysqlRes:[[
  ]]
gaeaRes:[[
  ]]

sql:select x'13181C76734725455A';
mysqlRes:[[vsG%EZ]]
gaeaRes:[[vsG%EZ]]

sql:select 0x4920616D2061206C6F6E672068657820737472696E67;
mysqlRes:[[I am a long hex string]]
gaeaRes:[[I am a long hex string]]

sql:select 0b01, 0b0, b'11', B'11';
mysqlRes:[[    ]]
gaeaRes:[[    ]]

sql:SELECT 1 > (select 1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT 1 > ANY (select 1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT 1 > ALL (select 1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT 1 > SOME (select 1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:SELECT EXISTS (select 1);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT + EXISTS (select 1);
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:SELECT - EXISTS (select 1);
mysqlRes:[[-1]]
gaeaRes:[[-1]]

sql:SELECT NOT EXISTS (select 1);
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select * from t1 where not exists (select * from t2 where t1.col1 = t2.col1);
mysqlRes:[]
gaeaRes:[]

sql:select col1 from t1 union select col2 from t2;
mysqlRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]
gaeaRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]

sql:select col1 from t1 union (select col2 from t2);
mysqlRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]
gaeaRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]

sql:select col1 from t1 union (select col2 from t2) order by col1;
mysqlRes:[[10] [15] [20] [30] [5] [aa] [bb] [cc] [dd] [ee]]
gaeaRes:[[10] [15] [20] [30] [5] [aa] [bb] [cc] [dd] [ee]]

sql:select col1 from t1 union (select col2 from t2) limit 1;
mysqlRes:[[aa]]
gaeaRes:[[aa]]

sql:select col1 from t1 union (select col2 from t2) limit 1, 1;
mysqlRes:[[bb]]
gaeaRes:[[bb]]

sql:select col1 from t1 union (select col2 from t2) order by col1 limit 1;
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:(select col1 from t1) union distinct select col2 from t2;
mysqlRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]
gaeaRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]

sql:(select col1 from t1) union distinctrow select col2 from t2;
mysqlRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]
gaeaRes:[[aa] [bb] [cc] [dd] [ee] [5] [10] [15] [20] [30]]

sql:(select col1 from t1) union all select col2 from t2;
mysqlRes:[[aa] [aa] [bb] [bb] [cc] [cc] [dd] [dd] [ee] [ee] [5] [5] [10] [10] [15] [15] [20] [20] [30] [30]]
gaeaRes:[[aa] [aa] [bb] [bb] [cc] [cc] [dd] [dd] [ee] [ee] [5] [5] [10] [10] [15] [15] [20] [20] [30] [30]]

sql:(select col1 from t1) union select col2 from t2 union (select col2 from t3) order by col1 limit 1;
mysqlRes:[[10]]
gaeaRes:[[10]]

sql:select (select 1 union select 1) as a;
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select * from (select 1 union select 2) as a;
mysqlRes:[[1] [2]]
gaeaRes:[[1] [2]]

sql:select 2 as a from dual union select 1 as b from dual order by a;
mysqlRes:[[1] [2]]
gaeaRes:[[1] [2]]

sql:select 2 as a from t1 union select 1 as b from t1 order by a;
mysqlRes:[[1] [2]]
gaeaRes:[[1] [2]]

sql:(select 2 as a from t order by a) union select 1 as b from t order by a;
mysqlRes:[[1] [2]]
gaeaRes:[[1] [2]]

sql:select 1 a, 2 b from dual order by a;
mysqlRes:[[1 2]]
gaeaRes:[[1 2]]

sql:select 1 a, 2 b from dual;
mysqlRes:[[1 2]]
gaeaRes:[[1 2]]

sql:select "abc_" like "abc\\_" escape '';
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select "abc_" like "abc\\_" escape '\\';
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select "abc" like "escape" escape '+';
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select '''_' like '''_' escape '''';
mysqlRes:[[0]]
gaeaRes:[[0]]

sql:select * from t use index (primary);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t use index (`primary`);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t use index ();
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t use index (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t use index (idx1, idx2);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t ignore key (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t force index for join (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t use index for order by (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t force index for group by (idx1);
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t use index for group by (idx1) use index for order by (idx2), t2;
mysqlRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]
gaeaRes:[[1 aa 5 1 aa 5] [2 bb 10 1 aa 5] [3 cc 15 1 aa 5] [4 dd 20 1 aa 5] [5 ee 30 1 aa 5] [6 aa 5 1 aa 5] [7 bb 10 1 aa 5] [8 cc 15 1 aa 5] [9 dd 20 1 aa 5] [10 ee 30 1 aa 5] [1 aa 5 2 bb 10] [2 bb 10 2 bb 10] [3 cc 15 2 bb 10] [4 dd 20 2 bb 10] [5 ee 30 2 bb 10] [6 aa 5 2 bb 10] [7 bb 10 2 bb 10] [8 cc 15 2 bb 10] [9 dd 20 2 bb 10] [10 ee 30 2 bb 10] [1 aa 5 3 cc 15] [2 bb 10 3 cc 15] [3 cc 15 3 cc 15] [4 dd 20 3 cc 15] [5 ee 30 3 cc 15] [6 aa 5 3 cc 15] [7 bb 10 3 cc 15] [8 cc 15 3 cc 15] [9 dd 20 3 cc 15] [10 ee 30 3 cc 15] [1 aa 5 4 dd 20] [2 bb 10 4 dd 20] [3 cc 15 4 dd 20] [4 dd 20 4 dd 20] [5 ee 30 4 dd 20] [6 aa 5 4 dd 20] [7 bb 10 4 dd 20] [8 cc 15 4 dd 20] [9 dd 20 4 dd 20] [10 ee 30 4 dd 20] [1 aa 5 5 ee 30] [2 bb 10 5 ee 30] [3 cc 15 5 ee 30] [4 dd 20 5 ee 30] [5 ee 30 5 ee 30] [6 aa 5 5 ee 30] [7 bb 10 5 ee 30] [8 cc 15 5 ee 30] [9 dd 20 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [2 bb 10 6 aa 5] [3 cc 15 6 aa 5] [4 dd 20 6 aa 5] [5 ee 30 6 aa 5] [6 aa 5 6 aa 5] [7 bb 10 6 aa 5] [8 cc 15 6 aa 5] [9 dd 20 6 aa 5] [10 ee 30 6 aa 5] [1 aa 5 7 bb 10] [2 bb 10 7 bb 10] [3 cc 15 7 bb 10] [4 dd 20 7 bb 10] [5 ee 30 7 bb 10] [6 aa 5 7 bb 10] [7 bb 10 7 bb 10] [8 cc 15 7 bb 10] [9 dd 20 7 bb 10] [10 ee 30 7 bb 10] [1 aa 5 8 cc 15] [2 bb 10 8 cc 15] [3 cc 15 8 cc 15] [4 dd 20 8 cc 15] [5 ee 30 8 cc 15] [6 aa 5 8 cc 15] [7 bb 10 8 cc 15] [8 cc 15 8 cc 15] [9 dd 20 8 cc 15] [10 ee 30 8 cc 15] [1 aa 5 9 dd 20] [2 bb 10 9 dd 20] [3 cc 15 9 dd 20] [4 dd 20 9 dd 20] [5 ee 30 9 dd 20] [6 aa 5 9 dd 20] [7 bb 10 9 dd 20] [8 cc 15 9 dd 20] [9 dd 20 9 dd 20] [10 ee 30 9 dd 20] [1 aa 5 10 ee 30] [2 bb 10 10 ee 30] [3 cc 15 10 ee 30] [4 dd 20 10 ee 30] [5 ee 30 10 ee 30] [6 aa 5 10 ee 30] [7 bb 10 10 ee 30] [8 cc 15 10 ee 30] [9 dd 20 10 ee 30] [10 ee 30 10 ee 30]]

sql:select high_priority * from t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select * from t;
mysqlRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]
gaeaRes:[[1 aa 5] [2 bb 10] [3 cc 15] [4 dd 20] [5 ee 30] [6 aa 5] [7 bb 10] [8 cc 15] [9 dd 20] [10 ee 30]]

sql:select """";
mysqlRes:[["]]
gaeaRes:[["]]

sql:select "汉字";
mysqlRes:[[汉字]]
gaeaRes:[[汉字]]

sql:select 'abc"def';
mysqlRes:[[abc"def]]
gaeaRes:[[abc"def]]

sql:select 'a\r\n';
mysqlRes:[[a
]]
gaeaRes:[[a
]]

sql:select "\a\r\n";
mysqlRes:[[a
]]
gaeaRes:[[a
]]

sql:select "\xFF";
mysqlRes:[[xFF]]
gaeaRes:[[xFF]]

sql:explain select col1 from t1;
mysqlRes:[[1 SIMPLE t1 NULL index NULL idx1 83 NULL 10 100.00 Using index]]
gaeaRes:[[1 SIMPLE t1 NULL index NULL idx1 83 NULL 10 100.00 Using index]]

sql:explain delete t1, t2 from t1 inner join t2 inner join t3 where t1.id=t2.id and t2.id=t3.id;
mysqlRes:[[1 DELETE t1 NULL ALL PRIMARY NULL NULL NULL 10 100.00 NULL] [1 DELETE t2 NULL eq_ref PRIMARY PRIMARY 4 sbtest1.t1.id 1 100.00 NULL] [1 SIMPLE t3 NULL eq_ref PRIMARY PRIMARY 4 sbtest1.t1.id 1 100.00 Using index]]
gaeaRes:[[1 DELETE t1 NULL ALL PRIMARY NULL NULL NULL 10 100.00 NULL] [1 DELETE t2 NULL eq_ref PRIMARY PRIMARY 4 sbtest1.t1.id 1 100.00 NULL] [1 SIMPLE t3 NULL eq_ref PRIMARY PRIMARY 4 sbtest1.t1.id 1 100.00 Using index]]

sql:explain update t set id = id + 1 order by id desc;
mysqlRes:[[1 UPDATE t NULL index NULL PRIMARY 4 NULL 10 100.00 Using filesort]]
gaeaRes:[[1 UPDATE t NULL index NULL PRIMARY 4 NULL 10 100.00 Using filesort]]

sql:EXPLAIN select col1 from t1 union (select col2 from t2) limit 1, 1;
mysqlRes:[[1 PRIMARY t1 NULL index NULL idx1 83 NULL 10 100.00 Using index] [2 UNION t2 NULL index NULL idx2 5 NULL 10 100.00 Using index] [NULL UNION RESULT <union1,2> NULL ALL NULL NULL NULL NULL NULL NULL Using temporary]]
gaeaRes:[[1 PRIMARY t1 NULL index NULL idx1 83 NULL 10 100.00 Using index] [2 UNION t2 NULL index NULL idx2 5 NULL 10 100.00 Using index] [NULL UNION RESULT <union1,2> NULL ALL NULL NULL NULL NULL NULL NULL Using temporary]]

sql:EXPLAIN SELECT 1;
mysqlRes:[[1 SIMPLE NULL NULL NULL NULL NULL NULL NULL NULL NULL No tables used]]
gaeaRes:[[1 SIMPLE NULL NULL NULL NULL NULL NULL NULL NULL NULL No tables used]]

sql:PREPARE pname FROM 'SELECT ?';
mysqlRes:[]
gaeaRes:[]

sql:SELECT col1, MAX(col2) AS max_col2 FROM t GROUP BY col1 WITH ROLLUP;
mysqlRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30] [NULL 30]]
gaeaRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30] [NULL 30]]

sql:SELECT COALESCE(col1,'ALL') AS coalesced_col1, MAX(col2) AS max_col2 FROM t GROUP BY col1 WITH ROLLUP;
mysqlRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30] [ALL 30]]
gaeaRes:[[aa 5] [bb 10] [cc 15] [dd 20] [ee 30] [ALL 30]]

sql:select col1 from t1 group by col1 order by null;
mysqlRes:[[aa] [bb] [cc] [dd] [ee]]
gaeaRes:[[aa] [bb] [cc] [dd] [ee]]

sql:select col1 from t1 group by col1 order by 1;
mysqlRes:[[aa] [bb] [cc] [dd] [ee]]
gaeaRes:[[aa] [bb] [cc] [dd] [ee]]

sql:SELECT CHAR(65);
mysqlRes:[[A]]
gaeaRes:[[A]]

sql:use  sbtest1;
mysqlRes:[]
gaeaRes:[]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select all * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select distinct * from sbtest1.test2) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3] [4 4 4 4] [5 5 1 5]]

sql:select id,o_id,name,pad from (select * from sbtest1.test2 a group by a.id) a;
mysqlRes:[[1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]
gaeaRes:[[1 1 order中id为1 1] [2 2 test_2 2] [3 3 order中id为3 3] [4 4 $order$4 4] [5 5 order...5 1]]

sql:select * from (select pad,count(*) from sbtest1.test2 a group by pad) a;
mysqlRes:[[1 2] [2 1] [3 1] [4 1]]
gaeaRes:[[1 2] [2 1] [3 1] [4 1]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 having pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3 order by id) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 order by id limit 3) b where a.t_id=b.o_id;
mysqlRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3]]
gaeaRes:[[1 1 1 1] [2 2 2 2] [3 3 3 3]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 order by id limit 3) b where a.t_id=b.o_id limit 2;
mysqlRes:[[1 1 1 1] [2 2 2 2]]
gaeaRes:[[1 1 1 1] [2 2 2 2]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select * from (select sbtest1.test3.pad from test1 left join sbtest1.test3 on test1.pad=sbtest1.test3.pad) a;
mysqlRes:[[1] [1] [2] [3] [4] [6]]
gaeaRes:[[1] [1] [2] [3] [4] [6]]

sql:select id,t_id,name,pad from (select * from test1 union select * from sbtest1.test3) a where a.id >3;
mysqlRes:[[4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [4 4 $manager$4 4] [5 5 manager...5 6]]
gaeaRes:[[4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6] [4 4 $manager$4 4] [5 5 manager...5 6]]

sql:select id,pad from test1 where pad=(select min(id) from sbtest1.test3);
mysqlRes:[[1 1] [5 1]]
gaeaRes:[[1 1] [5 1]]

sql:select id,pad,name from (select * from test1 where pad>2) a where id<5;
mysqlRes:[[3 4 test中id为3] [4 3 $test$4]]
gaeaRes:[[3 4 test中id为3] [4 3 $test$4]]

sql:select pad,count(*) from (select * from test1 where pad>2) a group by pad;
mysqlRes:[[3 1] [4 1] [6 1]]
gaeaRes:[[3 1] [4 1] [6 1]]

sql:select pad,count(*) from (select * from test1 where pad>2) a group by pad order by pad;
mysqlRes:[[3 1] [4 1] [6 1]]
gaeaRes:[[3 1] [4 1] [6 1]]

sql:select count(*) from (select pad,count(*) a from test1 group by pad) a;
mysqlRes:[[5]]
gaeaRes:[[5]]

sql:select id,t_id,name,pad from test1 where pad<(select pad from sbtest1.test3 where id=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]

sql:select id,t_id,name,pad from test1 having pad<(select pad from sbtest1.test3 where id=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [5 5 test...5 1]]

sql:select a.id,b.id,b.pad,a.t_id from test1 a,(select * from sbtest1.test2 where pad>3) b where a.t_id=b.o_id;
mysqlRes:[[4 4 4 4]]
gaeaRes:[[4 4 4 4]]

sql:select id,name,(select count(*) from sbtest1.test3) count from test1;
mysqlRes:[[1 test中id为1 5] [2 test_2 5] [3 test中id为3 5] [4 $test$4 5] [5 test...5 5] [6 test6 5]]
gaeaRes:[[1 test中id为1 5] [2 test_2 5] [3 test中id为3 5] [4 $test$4 5] [5 test...5 5] [6 test6 5]]

sql:select id,t_id,name,pad from test1 where pad like (select pad from sbtest1.test3 where id=3);
mysqlRes:[[4 4 $test$4 3]]
gaeaRes:[[4 4 $test$4 3]]

sql:select id,pad from test1 where pad>(select pad from test1 where id=2);
mysqlRes:[[3 4] [4 3] [6 6]]
gaeaRes:[[3 4] [4 3] [6 6]]

sql:select id,pad from test1 where pad<(select pad from test1 where id=2);
mysqlRes:[[1 1] [5 1]]
gaeaRes:[[1 1] [5 1]]

sql:select id,pad from test1 where pad=(select pad from test1 where id=2);
mysqlRes:[[2 2]]
gaeaRes:[[2 2]]

sql:select id,pad from test1 where pad>=(select pad from test1 where id=2);
mysqlRes:[[2 2] [3 4] [4 3] [6 6]]
gaeaRes:[[2 2] [3 4] [4 3] [6 6]]

sql:select id,pad from test1 where pad<=(select pad from test1 where id=2);
mysqlRes:[[1 1] [2 2] [5 1]]
gaeaRes:[[1 1] [2 2] [5 1]]

sql:select id,pad from test1 where pad<>(select pad from test1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]

sql:select id,pad from test1 where pad !=(select pad from test1 where id=2);
mysqlRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]
gaeaRes:[[1 1] [3 4] [4 3] [5 1] [6 6]]

sql:select id,t_id,name,pad from test1 where exists(select * from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where not exists(select * from test1 where pad>1);
mysqlRes:[]
gaeaRes:[]

sql:select id,t_id,name,pad from test1 where pad not in(select id from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [5 5 test...5 1]]
gaeaRes:[[1 1 test中id为1 1] [5 5 test...5 1]]

sql:select id,t_id,name,pad from test1 where pad in(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where pad=some(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where pad=any(select id from test1 where pad>1);
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where pad !=any(select id from test1 where pad=3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:SELECT a.id AS a_id, b.id AS b_id, b.pad AS b_pad, a.t_id AS a_t_id FROM (SELECT t1.id, t1.pad, t1.t_id FROM test1 AS t1 JOIN sbtest1.test3 AS t3 ON t1.pad = t3.pad) AS a JOIN (SELECT t2.id, t2.pad FROM test1 AS t1_2 JOIN sbtest1.test2 AS t2 ON t1_2.pad = t2.pad) AS b ON a.pad = b.pad;
mysqlRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]
gaeaRes:[[1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5] [2 2 2 2] [3 4 4 3] [4 3 3 4] [1 1 1 1] [1 5 1 1] [5 1 1 5] [5 5 1 5]]

sql:select id,t_id,name,pad from test1 where pad>(select pad from test1 where pad=2);
mysqlRes:[[3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]
gaeaRes:[[3 3 test中id为3 4] [4 4 $test$4 3] [6 6 test6 6]]

sql:select b.id,b.t_id,b.name,b.pad,a.id,a.id,a.pad,a.t_id from test1 b,(select * from test1 where id>3 union select * from sbtest1.test3 where id<2) a where a.id >3 and b.pad=a.pad;
mysqlRes:[[1 1 test中id为1 1 5 5 1 5] [4 4 $test$4 3 4 4 3 4] [5 5 test...5 1 5 5 1 5] [6 6 test6 6 6 6 6 6]]
gaeaRes:[[1 1 test中id为1 1 5 5 1 5] [4 4 $test$4 3 4 4 3 4] [5 5 test...5 1 5 5 1 5] [6 6 test6 6 6 6 6 6]]

sql:select count(*) from (select * from test1 where pad=(select pad from sbtest1.test3 where id=1)) a;
mysqlRes:[[2]]
gaeaRes:[[2]]

sql:select (select name from test1 limit 1);
mysqlRes:[[test中id为1]]
gaeaRes:[[test中id为1]]

sql:select id,t_id,name,pad from test1 where 'test_2'=(select name from sbtest1.test3 where id=2);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where 5=(select count(*) from sbtest1.test3);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where 'test_2' like(select name from sbtest1.test3 where id=2);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where 2 >any(select id from test1 where pad>1);
mysqlRes:[]
gaeaRes:[]

sql:select id,t_id,name,pad from test1 where 2 in(select id from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where 2<>some(select id from test1 where pad>1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where 2>all(select id from test1 where pad<1);
mysqlRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[1 1 test中id为1 1] [2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

sql:select id,t_id,name,pad from test1 where (id,pad)=(select id,pad from sbtest1.test3 limit 1);
mysqlRes:[[1 1 test中id为1 1]]
gaeaRes:[[1 1 test中id为1 1]]

sql:select id,t_id,name,pad from test1 where row(id,pad)=(select id,pad from sbtest1.test3 limit 1);
mysqlRes:[[1 1 test中id为1 1]]
gaeaRes:[[1 1 test中id为1 1]]

sql:select id,name,pad from test1 where (id,pad)in(select id,pad from sbtest1.test3);
mysqlRes:[[1 test中id为1 1] [2 test_2 2]]
gaeaRes:[[1 test中id为1 1] [2 test_2 2]]

sql:select id,name,pad from test1 where (1,1)in(select id,pad from sbtest1.test3);
mysqlRes:[[1 test中id为1 1] [2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]
gaeaRes:[[1 test中id为1 1] [2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]

sql:SELECT x.pad FROM test1 AS x WHERE x.id = ( SELECT y.pad FROM sbtest1.test3 AS y WHERE y.id = ( SELECT z.pad FROM sbtest1.test2 AS z WHERE y.id = z.id LIMIT 1 ) LIMIT 1 );
mysqlRes:[[1]]
gaeaRes:[[1]]

sql:select co1,co2,co3 from (select id as co1,name as co2,pad as co3 from test1)as tb where co1>1;
mysqlRes:[[2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]
gaeaRes:[[2 test_2 2] [3 test中id为3 4] [4 $test$4 3] [5 test...5 1] [6 test6 6]]

sql:select avg(sum_column1) from (select sum(id) as sum_column1 from test1 group by pad) as t1;
mysqlRes:[[4.2000]]
gaeaRes:[[4.2000]]

sql:select * from (select m.id,n.pad from test1 m,sbtest1.test2 n where m.id=n.id AND m.name='test中id为1' and m.pad>7 and m.pad<10)a;
mysqlRes:[]
gaeaRes:[]

sql:select id,t_id,name,pad from test1 where id in (select id from test1 where id in ( select id from test1 where id in (select id from test1 where id =0) or id not in(select id from test1 where id =1)));
mysqlRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]
gaeaRes:[[2 2 test_2 2] [3 3 test中id为3 4] [4 4 $test$4 3] [5 5 test...5 1] [6 6 test6 6]]

