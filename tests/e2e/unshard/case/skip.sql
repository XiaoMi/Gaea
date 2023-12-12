SHOW VARIABLES;
SHOW STATUS;
SHOW GLOBAL STATUS;
SHOW SESSION STATUS;
SELECT CONNECTION_ID();
SELECT ROW_COUNT();
SELECT GET_LOCK('lock1',10);
SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t1 INNER JOIN t2 where t1.col1 = t2.col1;
/*
res1 [[1 aa 5 1 aa 5] [1 aa 5 6 aa 5] [2 bb 10 2 bb 10] [2 bb 10 7 bb 10] [3 cc 15 3 cc 15] [3 cc 15 8 cc 15] [4 dd 20 4 dd 20] [4 dd 20 9 dd 20] [5 ee 30 5 ee 30] [5 ee 30 10 ee 30] [6 aa 5 1 aa 5] [6 aa 5 6 aa 5] [7 bb 10 2 bb 10] [7 bb 10 7 bb 10] [8 cc 15 3 cc 15] [8 cc 15 8 cc 15] [9 dd 20 4 dd 20] [9 dd 20 9 dd 20] [10 ee 30 5 ee 30] [10 ee 30 10 ee 30]]
res2 [[1 aa 5 1 aa 5] [6 aa 5 1 aa 5] [2 bb 10 2 bb 10] [7 bb 10 2 bb 10] [3 cc 15 3 cc 15] [8 cc 15 3 cc 15] [4 dd 20 4 dd 20] [9 dd 20 4 dd 20] [5 ee 30 5 ee 30] [10 ee 30 5 ee 30] [1 aa 5 6 aa 5] [6 aa 5 6 aa 5] [2 bb 10 7 bb 10] [7 bb 10 7 bb 10] [3 cc 15 8 cc 15] [8 cc 15 8 cc 15] [4 dd 20 9 dd 20] [9 dd 20 9 dd 20] [5 ee 30 10 ee 30] [10 ee 30 10 ee 30]] 
*/
