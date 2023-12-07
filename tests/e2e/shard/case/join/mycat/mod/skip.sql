SELECT last_insert_id();

# show
show table status/*allow_diff*/;
show index from t1;
show index in t1;
show keys from t1;
show keys in t1;
show databases;
show schemas;
show open tables;
show open tables from sbtest1;
show open tables in sbtest1;

show tables from sbtest1;
show tables in sbtest1;
show full tables from sbtest1;
show full tables in sbtest1;
show full tables from sbtest1 like 'aly%';
show full tables in sbtest1 where table_type like 'base%';
show columns from t1 from sbtest1;
show full columns from t1 from sbtest1;
show full columns from t1 from sbtest1 like 'n%';
show full columns from t1 from sbtest1 where field like 's%';
show index from t1 from sbtest1;
show index in t1 in sbtest1;
show index in t1 from sbtest1;
show index from t1 in sbtest1;
show keys from t1 from sbtest1;
show keys from t1 in sbtest1;
show keys in t1 in sbtest1;
show keys in t1 from sbtest1;
show create database sbtest1;
show create schema sbtest1;
show create schema if not exists sbtest1;
show create database if not exists sbtest1;
