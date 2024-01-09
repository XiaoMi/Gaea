SELECT last_insert_id();

# databases
show databases;
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
show tables from sbtest;
show tables in sbtest;
show full tables from sbtest;
show full tables in sbtest;
show full tables from sbtest like 'aly%';
show full tables in sbtest where table_type like 'base%';


# columns
show columns in t1;
show columns from t1;
show full columns from t1;
show columns from t1 from sbtest;
show full columns from t1 from sbtest;
show full columns from t1 from sbtest like 'n%';
show full columns from t1 from sbtest where field like 's%';

# index
show index from t1;
show index in t1;
show index from t1 from sbtest;
show index in t1 in sbtest;
show index in t1 from sbtest;
show index from t1 in sbtest;

# keys
show keys from t1;
show keys in t1;
show keys from t1 from sbtest;
show keys from t1 in sbtest;
show keys in t1 in sbtest;
show keys in t1 from sbtest;

# create
show create database sbtest;
show create schema sbtest;
show create schema if not exists sbtest;
show create database if not exists sbtest;
