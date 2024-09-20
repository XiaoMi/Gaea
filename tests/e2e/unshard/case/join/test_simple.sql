
use  sbtest1;
/* some comments */ SELECT CONVERT('111', SIGNED);
/* some comments */ /*comment*/ SELECT CONVERT('111', SIGNED);

SELECT /*comment*/ CONVERT('111', SIGNED) ;
SELECT CONVERT('111', /*comment*/ SIGNED) ;
SELECT CONVERT('111', SIGNED) /*comment*/;
/*comment*/ /*comment*/ select col1 /* this is a comment */ from t;
SELECT /*!40001 SQL_NO_CACHE */ * FROM t WHERE 1 limit 0, 2000;


select b'';
select B'';
SELECT * FROM t;
SELECT * FROM t AS u;
SELECT * FROM t1, t2;
SELECT * FROM t1 AS u, t2;
SELECT * FROM t1, t2 AS u;
SELECT * FROM t1 AS u, t2 AS v;
SELECT * FROM t, t1, t2;
SELECT DISTINCT * FROM t;
SELECT DISTINCTROW * FROM t;
SELECT ALL * FROM t;
SELECT id FROM t;
SELECT * FROM t WHERE 1 = 1;
select 1 from t;
select 1 from t limit 1;
select 1 from t where not exists (select 2);
select * from sbtest1.t1 where id > 4 and id <=8 order by col1 desc;
select 1 as a from t order by a;
select * from sbtest1.t1 where id > 1 order by id desc limit 10;
select * from sbtest1.t2 where id < 0;
select 1 as a from t where 1 < any (select 2) order by a;
select 1 order by 1;


SELECT * from t for update;
SELECT * from t lock in share mode;

DO 1;
DO 1, sleep(1);

SELECT * from t1, t2, t3;
select * from t1 join t2 left join t3 on t2.id = t3.id;
select * from t1 right join t2 on t1.id = t2.id left join t3 on t3.id = t2.id;
select * from t1 right join t2 using (id) left join t3 using (id);
select * from t1 natural join t2;
select * from t1 natural right join t2;
select * from t1 natural left outer join t2;


select * from t1 straight_join t2 on t1.id = t2.id;
select straight_join * from t1 join t2 on t1.id = t2.id;
select straight_join * from t1 left join t2 on t1.id = t2.id;
select straight_join * from t1 right join t2 on t1.id = t2.id;
select straight_join * from t1 straight_join t2 on t1.id = t2.id;


(select 1);
ANALYZE TABLE t;


SHOW VARIABLES LIKE 'character_set_results';
SHOW GLOBAL VARIABLES LIKE 'character_set_results';
SHOW SESSION VARIABLES LIKE 'character_set_results';
SHOW GLOBAL VARIABLES WHERE Variable_name = 'autocommit';


SHOW STATUS WHERE Variable_name;
SHOW FULL TABLES FROM sbtest1 LIKE 't%';
SHOW FULL TABLES WHERE Table_type != 'VIEW';

SHOW COLUMNS FROM t;
SHOW COLUMNS FROM sbtest1.t;
SHOW FIELDS FROM t;
SHOW TRIGGERS LIKE 't';

SHOW DATABASES LIKE 'test2';
SHOW PROCEDURE STATUS WHERE Db='test';
SHOW FUNCTION STATUS WHERE Db='test';
SHOW INDEX FROM t;
SHOW KEYS FROM t;
SHOW INDEX IN t;
SHOW KEYS IN t;
SHOW INDEXES IN t where true;
SHOW KEYS FROM t FROM sbtest1 where true;
SHOW PLUGINS;
SHOW PROFILES;
SHOW PRIVILEGES;
show character set;
show charset;


show collation;
show collation like 'utf8%';
show collation where Charset = 'utf8' and Collation = 'utf8_bin';

show columns in t;

show create table sbtest1.t;
show create table t;


show create database sbtest1;
show create database if not exists sbtest1;


SELECT ++1;
SELECT -+1;
SELECT -1;
SELECT --1;


select '''a''', """a""";
select '''a''';
select '\'a\'';
select "\"a\"";
select """a""";
select _utf8"string";
select _binary"string";
select N'string';
select n'string';

select 1 <=> 0, 1 <=> null, 1 = null;



select {ts '1989-09-10 11:11:11'};
select {d '1989-09-10'};
select {t '00:00:00.111'};


select {ts123 '1989-09-10 11:11:11'};
select {ts123 123};
select {ts123 1 xor 1};

SELECT POW(1, 2);
SELECT POW(1, 0.5);
SELECT POW(1, -1);
SELECT POW(-1, 1);

SELECT MOD(10, 2);
SELECT ROUND(-1.23);
SELECT ROUND(1.23, 1);
SELECT CEIL(-1.23);
SELECT CEILING(1.23);
SELECT FLOOR(-1.23);
SELECT LN(1);
SELECT LOG(-2);
SELECT LOG(2, 65536);
SELECT LOG2(2);
SELECT LOG10(10);
SELECT ABS(10);
SELECT CRC32('MySQL');
SELECT SIGN(0);
SELECT SQRT(0);
SELECT ACOS(1);
SELECT ASIN(1);
SELECT ATAN(0), ATAN(1), ATAN(1, 2);
SELECT COS(0);
SELECT COS(1);
SELECT COT(1);
SELECT DEGREES(0);
SELECT EXP(1);
SELECT PI();
SELECT RADIANS(1);
SELECT SIN(1);
SELECT TAN(1);
SELECT TRUNCATE(1.223,1);
SELECT SUBSTR('Quadratically',5);
SELECT SUBSTR('Quadratically',5, 3);
SELECT SUBSTR('Quadratically' FROM 5);
SELECT SUBSTR('Quadratically' FROM 5 FOR 3);
SELECT SUBSTRING('Quadratically',5);
SELECT SUBSTRING('Quadratically',5, 3);
SELECT SUBSTRING('Quadratically' FROM 5);
SELECT SUBSTRING('Quadratically' FROM 5 FOR 3);
SELECT CONVERT('111', SIGNED);
SELECT LEAST(1, 2, 3);
SELECT INTERVAL(1, 0, 1, 2);
SELECT DATE_ADD('2008-01-02', INTERVAL INTERVAL(1, 0, 1) DAY);


SELECT DATABASE();
SELECT SCHEMA();


SELECT VERSION();
SELECT BENCHMARK(1000000, AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3')));
SELECT CHARSET('abc');
SELECT COERCIBILITY('abc');
SELECT COLLATION('abc');

SELECT SUBSTRING_INDEX('www.mysql.com', '.', 2);
SELECT SUBSTRING_INDEX('www.mysql.com', '.', -2);
SELECT LOWER("A"), UPPER("a");
SELECT LCASE("A"), UCASE("a");
SELECT REPLACE('www.mysql.com', 'w', 'Ww');
SELECT LOCATE('bar', 'foobarbar');
SELECT LOCATE('bar', 'foobarbar', 5);


select row(1, 1) > row(1, 1), row(1, 1, 1) > row(1, 1, 1);
Select (1, 1) > (1, 1);

SELECT *, CAST(col1 AS CHAR CHARACTER SET utf8) FROM t;

select cast(1 as signed int);

SELECT last_insert_id();
SELECT last_insert_id(1);
SELECT binary 'a';
SELECT BIT_COUNT(1);



SELECT time('01:02:03');
SELECT time('01:02:03.1');
SELECT time('20.1');
SELECT TIMEDIFF('2000:01:01 00:00:00', '2000:01:01 00:00:00.000001');
SELECT TIMESTAMPDIFF(MONTH,'2003-02-01','2003-05-01');
SELECT TIMESTAMPDIFF(YEAR,'2002-05-01','2001-01-01');
SELECT TIMESTAMPDIFF(MINUTE,'2003-02-01','2003-05-01 12:05:55');


SELECT MICROSECOND('2009-12-31 23:59:59.000010');
SELECT SECOND('10:05:03');
SELECT MINUTE('2008-02-03 10:05:03');

SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE();
SELECT DATEDIFF('2003-12-31', '2003-12-30');
SELECT DATE('2003-12-31 01:02:03');
SELECT DATE_FORMAT('2003-12-31 01:02:03', '%W %M %Y');
SELECT DAY('2007-02-03');
SELECT DAYOFMONTH('2007-02-03');
SELECT DAYOFWEEK('2007-02-03');
SELECT DAYOFYEAR('2007-02-03');
SELECT DAYNAME('2007-02-03');
SELECT FROM_DAYS(1423);
SELECT WEEKDAY('2007-02-03');


SELECT UTC_DATE, UTC_DATE();
SELECT UTC_DATE(), UTC_DATE()+0;


SELECT WEEK('2007-02-03');
SELECT WEEK('2007-02-03', 0);
SELECT WEEKOFYEAR('2007-02-03');
SELECT MONTH('2007-02-03');
SELECT MONTHNAME('2007-02-03');
SELECT YEAR('2007-02-03');
SELECT YEARWEEK('2007-02-03');
SELECT YEARWEEK('2007-02-03', 0);
SELECT ADDTIME('01:00:00.999999', '02:00:00.999998');
SELECT SUBTIME('01:00:00.999999', '02:00:00.999998');


SELECT CONVERT_TZ('2004-01-01 12:00:00','+00:00','+10:00');
SELECT GET_FORMAT(DATE, 'USA');
SELECT GET_FORMAT(DATETIME, 'USA');
SELECT GET_FORMAT(TIME, 'USA');
SELECT GET_FORMAT(TIMESTAMP, 'USA');




SELECT MAKEDATE(2011,31);
SELECT MAKETIME(12,15,30);

SELECT PERIOD_ADD(200801,2);
SELECT PERIOD_DIFF(200802,200703);


SELECT QUARTER('2008-04-01');
SELECT SEC_TO_TIME(2378);
SELECT TIME_FORMAT('100:00:00', '%H %k %h %I %l');
SELECT TIME_TO_SEC('22:23:00');
SELECT TIMESTAMPADD(WEEK,1,'2003-01-02');
SELECT TO_DAYS('2007-10-07');
SELECT TO_SECONDS('2009-11-29');
SELECT LAST_DAY('2003-02-05');


select extract(microsecond from "2011-11-11 10:10:10.123456");
select extract(second from "2011-11-11 10:10:10.123456");
select extract(minute from "2011-11-11 10:10:10.123456");
select extract(hour from "2011-11-11 10:10:10.123456");
select extract(day from "2011-11-11 10:10:10.123456");
select extract(week from "2011-11-11 10:10:10.123456");
select extract(month from "2011-11-11 10:10:10.123456");
select extract(quarter from "2011-11-11 10:10:10.123456");
select extract(year from "2011-11-11 10:10:10.123456");
select extract(second_microsecond from "2011-11-11 10:10:10.123456");
select extract(minute_microsecond from "2011-11-11 10:10:10.123456");
select extract(minute_second from "2011-11-11 10:10:10.123456");
select extract(hour_microsecond from "2011-11-11 10:10:10.123456");
select extract(hour_second from "2011-11-11 10:10:10.123456");
select extract(hour_minute from "2011-11-11 10:10:10.123456");
select extract(day_microsecond from "2011-11-11 10:10:10.123456");
select extract(day_second from "2011-11-11 10:10:10.123456");
select extract(day_minute from "2011-11-11 10:10:10.123456");
select extract(day_hour from "2011-11-11 10:10:10.123456");
select extract(year_month from "2011-11-11 10:10:10.123456");


select from_unixtime(1447430881);
select from_unixtime(1447430881.123456);
select from_unixtime(1447430881.1234567);
select from_unixtime(1447430881.9999999);
select from_unixtime(1447430881, "%Y %D %M %h:%i:%s %x");
select from_unixtime(1447430881.123456, "%Y %D %M %h:%i:%s %x");
select from_unixtime(1447430881.1234567, "%Y %D %M %h:%i:%s %x");


SELECT CAST('test collated returns' AS CHAR CHARACTER SET utf8) COLLATE utf8_bin;
SELECT TRIM('  bar   ');
SELECT TRIM(LEADING 'x' FROM 'xxxbarxxx');
SELECT TRIM(BOTH 'x' FROM 'xxxbarxxx');
SELECT TRIM(TRAILING 'xyz' FROM 'barxxyz');
SELECT LTRIM(' foo ');
SELECT RTRIM(' bar ');
SELECT RPAD('hi', 6, 'c');
SELECT BIT_LENGTH('hi');

SELECT CHAR_LENGTH('abc');
SELECT CHARACTER_LENGTH('abc');
SELECT FIELD('ej', 'Hej', 'ej', 'Heja', 'hej', 'foo');
SELECT FIND_IN_SET('foo', 'foo,bar');
SELECT MAKE_SET(1,'a'), MAKE_SET(1,'a','b','c');
SELECT MID('Sakila', -5, 3);
SELECT OCT(12);
SELECT OCTET_LENGTH('text');
SELECT ORD('2');
SELECT POSITION('bar' IN 'foobarbar');

SELECT BIN(12);
SELECT ELT(1, 'ej', 'Heja', 'hej', 'foo');
SELECT EXPORT_SET(5,'Y','N'), EXPORT_SET(5,'Y','N',','), EXPORT_SET(5,'Y','N',',',4);
SELECT FROM_BASE64('abc');
SELECT TO_BASE64('abc');

SELECT LOAD_FILE('/tmp/picture');
SELECT LPAD('hi',4,'??');
SELECT LEFT("foobar", 3);
SELECT RIGHT("foobar", 3);


SELECT REPEAT("a", 10);


SELECT SLEEP(1);
SELECT ANY_VALUE(@arg);
SELECT INET_ATON('10.0.5.9');
SELECT INET_NTOA(167773449);
SELECT INET6_ATON('fdfe::5a55:caff:fefa:9089');
SELECT INET6_NTOA(INET_NTOA(167773449));
SELECT IS_IPV4('10.0.5.9');
SELECT IS_IPV4_COMPAT(INET6_ATON('::10.0.5.9'));
SELECT IS_IPV4_MAPPED(INET6_ATON('::10.0.5.9'));
SELECT IS_IPV6('10.0.5.9');
SELECT MASTER_POS_WAIT(@log_name, @log_pos), MASTER_POS_WAIT(@log_name, @log_pos, @timeout), MASTER_POS_WAIT(@log_name, @log_pos, @timeout, @channel_name);
SELECT NAME_CONST('myname', 14);
SELECT RELEASE_ALL_LOCKS();



select "2011-11-11 10:10:10.123456" + interval 10 second;
select "2011-11-11 10:10:10.123456" - interval 10 second;


select date_add("2011-11-11 10:10:10.123456", interval 10 microsecond);
select date_add("2011-11-11 10:10:10.123456", interval 10 second);
select date_add("2011-11-11 10:10:10.123456", interval 10 minute);
select date_add("2011-11-11 10:10:10.123456", interval 10 hour);
select date_add("2011-11-11 10:10:10.123456", interval 10 day);
select date_add("2011-11-11 10:10:10.123456", interval 1 week);
select date_add("2011-11-11 10:10:10.123456", interval 1 month);
select date_add("2011-11-11 10:10:10.123456", interval 1 quarter);
select date_add("2011-11-11 10:10:10.123456", interval 1 year);
select date_add("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
select date_add("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
select date_add("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
select date_add("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
select date_add("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
select date_add("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
select date_add("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
select date_add("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
select date_add("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
select date_add("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
select date_add("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
select date_add("2011-11-11 10:10:10.123456", interval "11-11" year_month);

select strcmp('abc', 'def');


select adddate("2011-11-11 10:10:10.123456", interval 10 microsecond);
select adddate("2011-11-11 10:10:10.123456", interval 10 second);
select adddate("2011-11-11 10:10:10.123456", interval 10 minute);
select adddate("2011-11-11 10:10:10.123456", interval 10 hour);
select adddate("2011-11-11 10:10:10.123456", interval 10 day);
select adddate("2011-11-11 10:10:10.123456", interval 1 week);
select adddate("2011-11-11 10:10:10.123456", interval 1 month);
select adddate("2011-11-11 10:10:10.123456", interval 1 quarter);
select adddate("2011-11-11 10:10:10.123456", interval 1 year);
select adddate("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
select adddate("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
select adddate("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
select adddate("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
select adddate("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
select adddate("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
select adddate("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
select adddate("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
select adddate("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
select adddate("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
select adddate("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
select adddate("2011-11-11 10:10:10.123456", interval "11-11" year_month);
select adddate("2011-11-11 10:10:10.123456", 10);
select adddate("2011-11-11 10:10:10.123456", 0.10);
select adddate("2011-11-11 10:10:10.123456", "11,11");


select date_sub("2011-11-11 10:10:10.123456", interval 10 microsecond);
select date_sub("2011-11-11 10:10:10.123456", interval 10 second);
select date_sub("2011-11-11 10:10:10.123456", interval 10 minute);
select date_sub("2011-11-11 10:10:10.123456", interval 10 hour);
select date_sub("2011-11-11 10:10:10.123456", interval 10 day);
select date_sub("2011-11-11 10:10:10.123456", interval 1 week);
select date_sub("2011-11-11 10:10:10.123456", interval 1 month);
select date_sub("2011-11-11 10:10:10.123456", interval 1 quarter);
select date_sub("2011-11-11 10:10:10.123456", interval 1 year);
select date_sub("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
select date_sub("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
select date_sub("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
select date_sub("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
select date_sub("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
select date_sub("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
select date_sub("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
select date_sub("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
select date_sub("2011-11-11 10:10:10.123456", interval "11-11" year_month);


select subdate("2011-11-11 10:10:10.123456", interval 10 microsecond);
select subdate("2011-11-11 10:10:10.123456", interval 10 second);
select subdate("2011-11-11 10:10:10.123456", interval 10 minute);
select subdate("2011-11-11 10:10:10.123456", interval 10 hour);
select subdate("2011-11-11 10:10:10.123456", interval 10 day);
select subdate("2011-11-11 10:10:10.123456", interval 1 week);
select subdate("2011-11-11 10:10:10.123456", interval 1 month);
select subdate("2011-11-11 10:10:10.123456", interval 1 quarter);
select subdate("2011-11-11 10:10:10.123456", interval 1 year);
select subdate("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond);
select subdate("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond);
select subdate("2011-11-11 10:10:10.123456", interval "10:10" minute_second);
select subdate("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond);
select subdate("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second);
select subdate("2011-11-11 10:10:10.123456", interval "10:10" hour_minute);
select subdate("2011-11-11 10:10:10.123456", interval 10.10 hour_minute);
select subdate("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond);
select subdate("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second);
select subdate("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute);
select subdate("2011-11-11 10:10:10.123456", interval "11 10" day_hour);
select subdate("2011-11-11 10:10:10.123456", interval "11-11" year_month);
select subdate("2011-11-11 10:10:10.123456", 10);
select subdate("2011-11-11 10:10:10.123456", 0.10);
select subdate("2011-11-11 10:10:10.123456", "11,11");


select unix_timestamp();
select unix_timestamp('2015-11-13 10:20:19.012');




select avg(distinct col1) from t;
select avg(distinctrow col1) from t;
select avg(distinct all col1) from t;
select avg(distinctrow all col1) from t;
select avg(col2) from t;

select bit_and(col1) from t;
select bit_and(all col1) from t;
select bit_or(col1) from t;
select bit_or(all col1) from t;
select bit_xor(col1) from t;
select bit_xor(all col1) from t;
select max(distinct col1) from t;
select max(distinctrow col1) from t;
select max(distinct all col1) from t;
select max(distinctrow all col1) from t;
select max(col2) from t;
select min(distinct col1) from t;
select min(distinctrow col1) from t;
select min(distinct all col1) from t;
select min(distinctrow all col1) from t;

select min(col2) from t;
select sum(distinct col1) from t;
select sum(distinctrow col1) from t;
select sum(distinct all col1) from t;
select sum(distinctrow all col1) from t;
select sum(col2) from t;
select count(col1) from t;
select count(*) from t;
select count(distinct col1, col2) from t;
select count(distinctrow col1, col2) from t;
select count(all col1) from t;

select group_concat(col2,col1) from t group by col1;
select group_concat(col2,col1 SEPARATOR ';') from t group by col1;
select group_concat(distinct col2,col1) from t group by col1;
select group_concat(distinctrow col2,col1) from t group by col1;
SELECT col1, GROUP_CONCAT(DISTINCT col2 ORDER BY col2 DESC SEPARATOR ' ') FROM t GROUP BY col1;
select std(col1), std(all col1) from t;
select stddev(col1), stddev(all col1) from t;
select stddev_pop(col1), stddev_pop(all col1) from t;
select stddev_samp(col1), stddev_samp(all col1) from t;
select variance(col1), variance(all col1) from t;
select var_pop(col1), var_pop(all col1) from t;
select var_samp(col1), var_samp(all col1) from t;


select AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3'));
select AES_DECRYPT(@crypt_str,@key_str);
select AES_DECRYPT(@crypt_str,@key_str,@init_vector);
SELECT COMPRESS('');


SELECT MD5('testing');

SELECT RANDOM_BYTES(@len);
SELECT SHA1('abc');
SELECT SHA('abc');
SELECT SHA2('abc', 224);
SELECT UNCOMPRESS('any string');
SELECT UNCOMPRESSED_LENGTH(@compressed_string);
SELECT VALIDATE_PASSWORD_STRENGTH(@str);


SELECT CHAR(65);
SELECT CHAR(65, 66, 67);
SELECT HEX(CHAR(1, 0)),HEX(CHAR(256)),HEX(CHAR(1, 1)),HEX(CHAR(257));
SELECT CHAR(0x027FA USING ucs2);
SELECT CHAR(0xc2a7 USING utf8);



select 1 as a, 1 as `a`, 1 as 'a';
select 1 as a, 1 as "a", 1 as 'a';
select 1 a, 1 "a", 1 'a';
select * from t a;
select * from t as a;
select 1 full, 1 `row`, 1 abs;
select * from t full, t1 `row`, t2 abs;


select .78+123;
select .78+.21;
select .78-123;
select .78-.21;
select .78--123;
select .78*123;
select .78*.21;
select .78/123;
select .78/.21;
select .78,123;
select .78,.21;
select .78 , 123;
select .78'123';
select .78`123`;
select .78"123";


select x'0a', X'11', 0x11;
select x'13181C76734725455A';
select 0x4920616D2061206C6F6E672068657820737472696E67;


select 0b01, 0b0, b'11', B'11';

SELECT 1 > (select 1);
SELECT 1 > ANY (select 1);
SELECT 1 > ALL (select 1);
SELECT 1 > SOME (select 1);


SELECT EXISTS (select 1);
SELECT + EXISTS (select 1);
SELECT - EXISTS (select 1);
SELECT NOT EXISTS (select 1);

select * from t1 where not exists (select * from t2 where t1.col1 = t2.col1);


select col1 from t1 union select col2 from t2;
select col1 from t1 union (select col2 from t2);
select col1 from t1 union (select col2 from t2) order by col1;
select col1 from t1 union (select col2 from t2) limit 1;
select col1 from t1 union (select col2 from t2) limit 1, 1;
select col1 from t1 union (select col2 from t2) order by col1 limit 1;
(select col1 from t1) union distinct select col2 from t2;
(select col1 from t1) union distinctrow select col2 from t2;
(select col1 from t1) union all select col2 from t2;
(select col1 from t1) union select col2 from t2 union (select col2 from t3) order by col1 limit 1;
select (select 1 union select 1) as a;
select * from (select 1 union select 2) as a;
select 2 as a from dual union select 1 as b from dual order by a;


select 2 as a from t1 union select 1 as b from t1 order by a;
(select 2 as a from t order by a) union select 1 as b from t order by a;
select 1 a, 2 b from dual order by a;
select 1 a, 2 b from dual;


select "abc_" like "abc\\_" escape '';
select "abc_" like "abc\\_" escape '\\';
select "abc" like "escape" escape '+';
select '''_' like '''_' escape '''';


select * from t use index (primary);
select * from t use index (`primary`);
select * from t use index ();
select * from t use index (idx1);
select * from t use index (idx1, idx2);
select * from t ignore key (idx1);
select * from t force index for join (idx1);
select * from t use index for order by (idx1);
select * from t force index for group by (idx1);
select * from t use index for group by (idx1) use index for order by (idx2), t2;


select high_priority * from t;

select * from t;


select """";
select "汉字";
select 'abc"def';
select 'a\r\n';
select "\a\r\n";
select "\xFF";


explain select col1 from t1;
explain delete t1, t2 from t1 inner join t2 inner join t3 where t1.id=t2.id and t2.id=t3.id;
explain update t set id = id + 1 order by id desc;
EXPLAIN select col1 from t1 union (select col2 from t2) limit 1, 1;
EXPLAIN SELECT 1;

PREPARE pname FROM 'SELECT ?';


SELECT col1, MAX(col2) AS max_col2 FROM t GROUP BY col1 WITH ROLLUP;
SELECT COALESCE(col1,'ALL') AS coalesced_col1, MAX(col2) AS max_col2 FROM t GROUP BY col1 WITH ROLLUP;

select col1 from t1 group by col1 order by null;
select col1 from t1 group by col1 order by 1;

# simple function
## DIV
SELECT 5 DIV 2, -5 DIV 2, 5 DIV -2, -5 DIV -2;

## other
SELECT INSERT('Quadratic', 3, 4, 'What');
SELECT CHAR(65);
