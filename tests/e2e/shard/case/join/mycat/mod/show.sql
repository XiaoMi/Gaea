# databases
show databases like 'sbtest';

# schemas
show schemas;
show open tables;
show open tables from sbtest;
show open tables in sbtest;
show open tables from sbtest like 'aly_o%';

# tables
show table status like 'aly_o%'/*allow_diff*/;
show table status/*allow_diff*/;
show tables;
show full tables;
show tables like 'aly_o%';
show tables from sbtest_0;
show tables in sbtest_0;
show full tables from sbtest_0;
show full tables in sbtest_0;
show full tables from sbtest_0 like 'aly%';
show full tables in sbtest_0 where table_type like 'base%';


# columns
show columns in t1;
show columns from t1;
show full columns from t1;
show columns from t1 from sbtest_0;
show full columns from t1 from sbtest_0;
show full columns from t1 from sbtest_0 like 'n%';
show full columns from t1 from sbtest_0 where field like 's%';

# index
show index from t1;
show index in t1;
show index from t1 from sbtest_0;
show index in t1 in sbtest_0;
show index in t1 from sbtest_0;
show index from t1 in sbtest_0;

# keys
show keys from t1;
show keys in t1;
show keys from t1 from sbtest_0;
show keys from t1 in sbtest_0;
show keys in t1 in sbtest_0;
show keys in t1 from sbtest_0;

# create
show create database sbtest_0;
show create schema sbtest_0;
show create schema if not exists sbtest_0;
show create database if not exists sbtest_0;
