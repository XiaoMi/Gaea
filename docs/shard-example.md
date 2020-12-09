# gaea 分片规则示例说明

## 导航
- [gaea kingshard hash分片示例](#gaea_kingshard_hash)
- [gaea kingshard mod分片示例](#gaea_kingshard_mod)
- [gaea kingshard range分片示例](#gaea_kingshard_range)
- [gaea kingshard date year分片示例](#gaea_kingshard_date_year)
- [gaea kingshard date month分片示例](#gaea_kingshard_date_month)
- [gaea kingshard date day分片示例](#gaea_kingshard_date_day)
- [gaea mycat mod分片示例](#gaea_mycat_mod)
- [gaea mycat_long(固定hash分片算法)分片示例](#gaea_mycat_long)
- [gaea mycat_murmur(一致性Hash)分片示例](#gaea_mycat_partitionByMurmurHash)
- [gaea mycat_string(字符串拆分hash)分片示例](#gaea_mycat_partitionByString)

<h2 id="gaea_kingshard_hash">gaea kingshard hash分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice定义一个库，每个库预定义2张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_kingshard_hash",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_kingshard": true
    },
    "default_phy_dbs": {
        "db_kingshard": "db_kingshard"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_kingshard",
            "table": "shard_hash",
            "type": "hash",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_kingshard_hash",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表shard_hash_0000、shard_hash_0001
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_hash_000"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表shard_hash_0002、shard_hash_0003
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_hash_000"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307示例，查询slice-0分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_hash_0000        |
| shard_hash_0001        |
+------------------------+
2 rows in set (0.01 sec)
#登录3308示例，查询slice-1分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_hash_0002        |
| shard_hash_0003        |
+------------------------+
2 rows in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 10`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_hash (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> select * from shard_hash;
+----+--------+
| id | col1   |
+----+--------+
|  4 | test4  |
|  8 | test8  |
|  1 | test1  |
|  5 | test5  |
|  9 | test9  |
|  2 | test2  |
|  6 | test6  |
| 10 | test10 |
|  3 | test3  |
|  7 | test7  |
+----+--------+
10 rows in set (0.03 sec)
#连接3307数据库实例，对slice-0分表数据进行查询：
mysql> select * from shard_hash_0000;
+----+-------+
| id | col1  |
+----+-------+
|  4 | test4 |
|  8 | test8 |
+----+-------+
2 rows in set (0.00 sec)
mysql> select * from shard_hash_0001;
+----+-------+
| id | col1  |
+----+-------+
|  1 | test1 |
|  5 | test5 |
|  9 | test9 |
+----+-------+
3 rows in set (0.01 sec)
#连接3308数据库实例，对slice-1分表数据进行查询：
mysql>  select * from shard_hash_0002;
+----+--------+
| id | col1   |
+----+--------+
|  2 | test2  |
|  6 | test6  |
| 10 | test10 |
+----+--------+
3 rows in set (0.01 sec)
mysql>  select * from shard_hash_0003;
+----+-------+
| id | col1  |
+----+-------+
|  3 | test3 |
|  7 | test7 |
+----+-------+
2 rows in set (0.01 sec)
```

<h2 id="gaea_kingshard_mod">gaea kingshard mod分片示例</h2>

我们预定义两个分片slice-0、slice-1，每个slice定义一个库，每个库预定义2张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_kingshard_mod",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_kingshard": true
    },
    "default_phy_dbs": {
        "db_kingshard": "db_kingshard"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_kingshard",
            "table": "shard_mod",
            "type": "mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_kingshard_mod",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_mod_0000、shard_mod_0001
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_mod_000"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_mod_0002、shard_mod_0003
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_mod_000"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_mod_0000         |
| shard_mod_0001         |
+------------------------+
2 rows in set (0.01 sec)
#登录3308示例，查询slice-1分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_mod_0002         |
| shard_mod_0003         |
+------------------------+
2 rows in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 10`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_mod (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> select * from shard_mod;
+----+--------+
| id | col1   |
+----+--------+
|  4 | test4  |
|  8 | test8  |
|  1 | test1  |
|  5 | test5  |
|  9 | test9  |
|  2 | test2  |
|  6 | test6  |
| 10 | test10 |
|  3 | test3  |
|  7 | test7  |
+----+--------+
10 rows in set (0.03 sec)

#连接3307数据库实例，对slice-0分表数据进行查询：
mysql> select * from shard_mod_0000;
+----+-------+
| id | col1  |
+----+-------+
|  4 | test4 |
|  8 | test8 |
+----+-------+
2 rows in set (0.00 sec)
mysql> select * from shard_mod_0001;
+----+-------+
| id | col1  |
+----+-------+
|  1 | test1 |
|  5 | test5 |
|  9 | test9 |
+----+-------+
3 rows in set (0.01 sec)

#连接3308数据库实例，对slice-1分表数据进行查询：
mysql> select * from shard_mod_0002;
+----+--------+
| id | col1   |
+----+--------+
|  2 | test2  |
|  6 | test6  |
| 10 | test10 |
+----+--------+
3 rows in set (0.00 sec)
mysql> select * from shard_mod_0003;
+----+-------+
| id | col1  |
+----+-------+
|  3 | test3 |
|  7 | test7 |
+----+-------+
2 rows in set (0.01 sec)
```

<h2 id="gaea_kingshard_range">gaea kingshard range分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice定义一个库，每个库预定义2张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_kingshard_range",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_kingshard": true
    },
    "default_phy_dbs": {
        "db_kingshard": "db_kingshard"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_kingshard",
            "table": "shard_range",
            "type": "range",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "table_row_limit": 3
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_kingshard_range",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```
其中，"table_row_limit:3"配置含义为:每张子表的记录数，分表字段位于区间[0,3)在shard_range_0000上，分表字段位于区间[3,6)在子表shard_range_0001上，依此类推...

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_range_0000、shard_range_0001
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_range_000"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_range_0002、shard_range_0003
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_range_000"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_range_0000       |
| shard_range_0001       |
+------------------------+
2 rows in set (0.00 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_range_0002       |
| shard_range_0003       |
+------------------------+
2 rows in set (0.01 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 10`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_range (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> select * from shard_range;
+----+--------+
| id | col1   |
+----+--------+
|  1 | test1  |
|  2 | test2  |
|  3 | test3  |
|  4 | test4  |
|  5 | test5  |
|  6 | test6  |
|  7 | test7  |
|  8 | test8  |
|  9 | test9  |
| 10 | test10 |
+----+--------+
10 rows in set (0.03 sec)
#连接3307数据库实例，对slice-0分表数据进行查询：
mysql> select * from shard_range_0000;
+----+-------+
| id | col1  |
+----+-------+
|  1 | test1 |
|  2 | test2 |
+----+-------+
2 rows in set (0.01 sec)
mysql> select * from shard_range_0001;
+----+-------+
| id | col1  |
+----+-------+
|  3 | test3 |
|  4 | test4 |
|  5 | test5 |
+----+-------+
3 rows in set (0.01 sec)
#连接3308数据库实例，对slice-1分表数据进行查询：
mysql> select * from shard_range_0002;
+----+-------+
| id | col1  |
+----+-------+
|  6 | test6 |
|  7 | test7 |
|  8 | test8 |
+----+-------+
3 rows in set (0.01 sec)
mysql> select * from shard_range_0003;
+----+--------+
| id | col1   |
+----+--------+
|  9 | test9  |
| 10 | test10 |
+----+--------+
2 rows in set (0.00 sec)
```

<h2 id="gaea_kingshard_date_year">gaea kingshard date year分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice定义一个库，每个库预定义2张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_kingshard_date_year",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_kingshard": true
    },
    "default_phy_dbs": {
        "db_kingshard": "db_kingshard"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_kingshard",
            "table": "shard_year",
            "type": "date_year",
            "key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "date_range": [
                "2016-2017",
                "2018-2019"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_kingshard_date_year",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_year_2016、shard_year_2017
for i in `seq 6 7`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_year_201"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),create_time datetime DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_year_2018、shard_year_2019
for i in `seq 8 9`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_year_201"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),create_time datetime DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_year_2016        |
| shard_year_2017        |
+------------------------+
2 rows in set (0.00 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_year_2018        |
| shard_year_2019        |
+------------------------+
2 rows in set (0.01 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 6 9`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_year (id, col1,create_time) values(${i}, 'test$i','201$i-07-01')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> select * from shard_year;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  6 | test6 | 2016-07-01 00:00:00 |
|  7 | test7 | 2017-07-01 00:00:00 |
|  8 | test8 | 2018-07-01 00:00:00 |
|  9 | test9 | 2019-07-01 00:00:00 |
+----+-------+---------------------+
4 rows in set (0.03 sec)

#连接3307数据库实例，对slice-0分表数据进行查询：
mysql> select * from shard_year_2016;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  6 | test6 | 2016-07-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)
mysql> select * from shard_year_2017;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  7 | test7 | 2017-07-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)
#连接3308数据库实例，对slice-1分表数据进行查询：
mysql> select * from shard_year_2018;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  8 | test8 | 2018-07-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)
mysql> select * from shard_year_2019;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  9 | test9 | 2019-07-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.00 sec)
```


<h2 id="gaea_kingshard_date_month">gaea kingshard date month分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice定义一个库，每个库预定义2张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_kingshard_date_month",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_kingshard": true
    },
    "default_phy_dbs": {
        "db_kingshard": "db_kingshard"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_kingshard",
            "table": "shard_month",
            "type": "date_month",
            "key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "date_range": [
                "201405-201406",
                "201408-201409"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_kingshard_date_month",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_month_201405、shard_month_201406
for i in `seq 5 6`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_month_20140"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),create_time datetime DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_month_201408、shard_month_201409
for i in `seq 8 9`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_month_20140"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),create_time datetime DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_month_201405     |
| shard_month_201406     |
+------------------------+
2 rows in set (0.01 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_month_201408     |
| shard_month_201409     |
+------------------------+
2 rows in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 5 6`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_month (id, col1,create_time) values(${i}, 'test$i','2014-0$i-01')";done
for i in `seq 8 9`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_month (id, col1,create_time) values(${i}, 'test$i','2014-0$i-01')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> select * from shard_month;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  5 | test5 | 2014-05-01 00:00:00 |
|  6 | test6 | 2014-06-01 00:00:00 |
|  8 | test8 | 2014-08-01 00:00:00 |
|  9 | test9 | 2014-09-01 00:00:00 |
+----+-------+---------------------+
4 rows in set (0.03 sec)

#连接3307数据库实例，对slice-0分表数据进行查询：
mysql> select * from shard_month_201405;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  5 | test5 | 2014-05-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)

mysql> select * from shard_month_201406;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  6 | test6 | 2014-06-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)
#连接3308数据库实例，对slice-1分表数据进行查询：
mysql> select * from shard_month_201408;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  8 | test8 | 2014-08-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.00 sec)

mysql> select * from shard_month_201409;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  9 | test9 | 2014-09-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.00 sec)

```

<h2 id="gaea_kingshard_date_day">gaea kingshard date day分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice定义一个库，每个库预定义2张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_kingshard_date_day",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_kingshard": true
    },
    "default_phy_dbs": {
        "db_kingshard": "db_kingshard"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_kingshard",
            "table": "shard_day",
            "type": "date_day",
            "key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "date_range": [
                "20201201-20201202",
                "20201203-20201204"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_kingshard_date_day",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_month_201405、shard_day_20201201、shard_day_20201202
for i in `seq 1 2`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_day_2020120"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),create_time datetime DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_kingshard;
#在命令行执行以下命令，创建分表,shard_day_20201203、shard_day_20201204
for i in `seq 3 4`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_kingshard -e "CREATE TABLE IF NOT EXISTS shard_day_2020120"${i}" ( id INT(64) NOT NULL, col1 VARCHAR(256),create_time datetime DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_day_20201201     |
| shard_day_20201202     |
+------------------------+
2 rows in set (0.00 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> show tables;
+------------------------+
| Tables_in_db_kingshard |
+------------------------+
| shard_day_20201203     |
| shard_day_20201204     |
+------------------------+
2 rows in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 4`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_kingshard -e "insert into shard_day (id, col1,create_time) values(${i}, 'test$i','2020-12-0$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> select * from shard_day;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  1 | test1 | 2020-12-01 00:00:00 |
|  2 | test2 | 2020-12-02 00:00:00 |
|  3 | test3 | 2020-12-03 00:00:00 |
|  4 | test4 | 2020-12-04 00:00:00 |
+----+-------+---------------------+
4 rows in set (0.03 sec)

#连接3307数据库实例，对slice-0分表数据进行查询：
mysql> select * from shard_day_20201201;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  1 | test1 | 2020-12-01 00:00:00 |
+----+-------+---------------------+
1 row in set (0.00 sec)

mysql> select * from shard_day_20201202;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  2 | test2 | 2020-12-02 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)
#连接3308数据库实例，对slice-1分表数据进行查询：
mysql> select * from shard_day_20201203;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  3 | test3 | 2020-12-03 00:00:00 |
+----+-------+---------------------+
1 row in set (0.00 sec)

mysql> select * from shard_day_20201204;
+----+-------+---------------------+
| id | col1  | create_time         |
+----+-------+---------------------+
|  4 | test4 | 2020-12-04 00:00:00 |
+----+-------+---------------------+
1 row in set (0.01 sec)

```

<h2 id="gaea_mycat_mod">gaea mycat mod分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice预定义2个库，每个库一张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_mycat_mod",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_mycat": "db_mycat"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_mycat_mod",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_mycat_0;
create database db_mycat_1;
#在命令行执行以下命令，创建分表
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat  ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_mycat_2;
create database db_mycat_3;
#在命令行执行以下命令，创建分表
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_0 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_1
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_1 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.01 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_2 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_3 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 10`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_mycat -e "insert into tbl_mycat (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> use db_mycat
Database changed
mysql> select * from tbl_mycat;
+----+--------+
| id | col1   |
+----+--------+
|  4 | test4  |
|  8 | test8  |
|  1 | test1  |
|  5 | test5  |
|  9 | test9  |
|  3 | test3  |
|  7 | test7  |
|  2 | test2  |
|  6 | test6  |
| 10 | test10 |
+----+--------+
10 rows in set (0.04 sec)

#连接3307数据库实例，对slice-0分片数据进行查询：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+----+-------+
| id | col1  |
+----+-------+
|  4 | test4 |
|  8 | test8 |
+----+-------+
2 rows in set (0.01 sec)

mysql> use db_mycat_1;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+----+-------+
| id | col1  |
+----+-------+
|  1 | test1 |
|  5 | test5 |
|  9 | test9 |
+----+-------+
3 rows in set (0.00 sec)
#连接3308数据库实例，对slice-1分片数据进行查询：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+----+--------+
| id | col1   |
+----+--------+
|  2 | test2  |
|  6 | test6  |
| 10 | test10 |
+----+--------+
3 rows in set (0.01 sec)

mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+----+-------+
| id | col1  |
+----+-------+
|  3 | test3 |
|  7 | test7 |
+----+-------+
2 rows in set (0.01 sec)
```

<h2 id="gaea_mycat_long">gaea mycat_long(固定hash分片算法)分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice预定义2个库，每个库一张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_mycat_long",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_mycat": "db_mycat"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_long",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
            "partition_count": "4",
            "partition_length": "256"
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_mycat_long",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_mycat_0;
create database db_mycat_1;
#在命令行执行以下命令，创建分表
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat  ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_mycat_2;
create database db_mycat_3;
#在命令行执行以下命令，创建分表
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_0 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_1
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_1 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.01 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_2 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_3 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入,插入2000行记录。
for i in `seq 1 2000`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_mycat -e "insert into tbl_mycat (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> use db_mycat
Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    1 | test1    |
|    2 | test2    |
|    3 | test3    |
...................
| 1786 | test1786 |
| 1787 | test1787 |
| 1788 | test1788 |
| 1789 | test1789 |
| 1790 | test1790 |
| 1791 | test1791 |
+------+----------+
2000 rows in set (0.61 sec)
#连接3307数据库实例，对slice-0分片数据进行查询：
mysql> use db_mycat_0;
Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    1 | test1    |
|    2 | test2    |
|    3 | test3    |
..................
|  249 | test249  |
|  250 | test250  |
|  251 | test251  |
|  252 | test252  |
|  253 | test253  |
|  254 | test254  |
|  255 | test255  |
| 1024 | test1024 |
| 1025 | test1025 |
| 1026 | test1026 |
| 1027 | test1027 |
| 1028 | test1028 |
..................
| 1277 | test1277 |
| 1278 | test1278 |
| 1279 | test1279 |
+------+----------+
511 rows in set (0.01 sec)

mysql> use db_mycat_1;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|  256 | test256  |
|  257 | test257  |
|  258 | test258  |
|  259 | test259  |
|  260 | test260  |
...................
|  509 | test509  |
|  510 | test510  |
|  511 | test511  |
| 1280 | test1280 |
| 1281 | test1281 |
| 1282 | test1282 |
| 1283 | test1283 |
| 1284 | test1284 |
| 1285 | test1285 |
...................
| 1532 | test1532 |
| 1533 | test1533 |
| 1534 | test1534 |
| 1535 | test1535 |
+------+----------+
512 rows in set (0.00 sec)
#连接3308数据库实例，对slice-1分片数据进行查询：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|  512 | test512  |
|  513 | test513  |
|  514 | test514  |
|  515 | test515  |
...................
|  765 | test765  |
|  766 | test766  |
|  767 | test767  |
| 1536 | test1536 |
| 1537 | test1537 |
| 1538 | test1538 |
| 1539 | test1539 |
| 1540 | test1540 |
...................
| 1786 | test1786 |
| 1787 | test1787 |
| 1788 | test1788 |
| 1789 | test1789 |
| 1790 | test1790 |
| 1791 | test1791 |
+------+----------+
512 rows in set (0.01 sec)

mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|  768 | test768  |
|  769 | test769  |
|  770 | test770  |
|  771 | test771  |
|  772 | test772  |
|  773 | test773  |
...................
|  996 | test996  |
|  997 | test997  |
|  998 | test998  |
|  999 | test999  |
| 1000 | test1000 |
| 1001 | test1001 |
| 1002 | test1002 |
| 1003 | test1003 |
| 1004 | test1004 |
...................
| 1993 | test1993 |
| 1994 | test1994 |
| 1995 | test1995 |
| 1996 | test1996 |
| 1997 | test1997 |
| 1998 | test1998 |
| 1999 | test1999 |
| 2000 | test2000 |
+------+----------+
465 rows in set (0.00 sec)
```

<h2 id="gaea_mycat_partitionByMurmurHash">gaea mycat_murmur(一致性Hash)分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice预定义2个库，每个库一张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_mycat_string",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_mycat": "db_mycat"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_string",
            "key": "col1",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
            "partition_count": "4",
            "partition_length": "256",
            "hash_slice": "-3:0"
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_mycat_string",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_mycat_0;
create database db_mycat_1;
#在命令行执行以下命令，创建分表
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat  ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_mycat_2;
create database db_mycat_3;
#在命令行执行以下命令，创建分表
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_0 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_1
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_1 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.01 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_2 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_3 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 2000`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_mycat -e "insert into tbl_mycat (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> use db_mycat
Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    1 | test1    |
|    2 | test2    |
|    3 | test3    |
|    9 | test9    |
|   10 | test10   |
|   15 | test15   |
..................
| 1982 | test1982 |
| 1987 | test1987 |
| 1988 | test1988 |
| 1990 | test1990 |
| 1997 | test1997 |
| 1999 | test1999 |
+------+----------+
2000 rows in set (0.05 sec)

#连接3307数据库实例，对slice-0分片数据进行查询：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    5 | test5    |
|    6 | test6    |
|    8 | test8    |
|   14 | test14   |
|   16 | test16   |
...................
| 1984 | test1984 |
| 1989 | test1989 |
| 1992 | test1992 |
| 1998 | test1998 |
| 2000 | test2000 |
+------+----------+
522 rows in set (0.01 sec)

mysql> use db_mycat_1;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    1 | test1    |
|    2 | test2    |
|    3 | test3    |
|    9 | test9    |
|   10 | test10   |
|   15 | test15   |
...................
| 1973 | test1973 |
| 1976 | test1976 |
| 1979 | test1979 |
| 1983 | test1983 |
| 1985 | test1985 |
| 1993 | test1993 |
+------+----------+
502 rows in set (0.01 sec)
#连接3308数据库实例，对slice-1分片数据进行查询：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    4 | test4    |
|    7 | test7    |
|   11 | test11   |
|   24 | test24   |
|   30 | test30   |
...................
| 1974 | test1974 |
| 1986 | test1986 |
| 1991 | test1991 |
| 1994 | test1994 |
| 1995 | test1995 |
| 1996 | test1996 |
+------+----------+
457 rows in set (0.01 sec)

mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|   12 | test12   |
|   13 | test13   |
|   20 | test20   |
|   22 | test22   |
...................
| 1982 | test1982 |
| 1987 | test1987 |
| 1988 | test1988 |
| 1990 | test1990 |
| 1997 | test1997 |
| 1999 | test1999 |
+------+----------+
519 rows in set (0.01 sec)
```

<h2 id="gaea_mycat_partitionByString">gaea mycat_string(字符串拆分hash)分片示例</h2>

我们预定义两个slice slice-0、slice-1，每个slice预定义2个库，每个库一张表，其中slice-0的主库地址为127.0.0.1:3307，slice-1的主库地址为127.0.0.1:3308。

Gaea启动地址为127.0.0.1:13307

### namespace配置
```json
{
    "name": "test_mycat_mod",
    "online": true,
    "read_only": false,
    "allowed_dbs": {
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_mycat": "db_mycat"
    },
    "slow_sql_time": "1000",
    "black_sql": [
        ""
    ],
    "allowed_ip": null,
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3307",
            "slaves": [],
            "statistic_slaves": null,
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "1234",
            "master": "127.0.0.1:3308",
            "slaves": [],
            "statistic_slaves": [],
            "capacity": 12,
            "max_capacity": 24,
            "idle_timeout": 60
        }
    ],
    "shard_rules": [
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        }
    ],
    "users": [
        {
            "user_name": "test",
            "password": "1234",
            "namespace": "test_mycat_mod",
            "rw_flag": 2,
            "rw_split": 1,
            "other_property": 0
        }
    ],
    "default_slice": "slice-1",
    "global_sequences": null
}
```

### 创建数据库表
```shell script
#连接3307数据库实例
mysql -h127.0.0.1 -P3307 -uroot -p1234
#创建数据库
create database db_mycat_0;
create database db_mycat_1;
#在命令行执行以下命令，创建分表
for i in `seq 0 1`;do  mysql -h127.0.0.1 -P3307 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat  ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done

#连接3306数据库实例
mysql -h127.0.0.1 -P3308 -uroot -p1234
#创建数据库
create database db_mycat_2;
create database db_mycat_3;
#在命令行执行以下命令，创建分表
for i in `seq 2 3`;do  mysql -h127.0.0.1 -P3308 -uroot -p1234  db_mycat_$i -e "CREATE TABLE IF NOT EXISTS tbl_mycat ( id INT(64) NOT NULL, col1 VARCHAR(256),PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";done
#登录3307实例，查询slice-0分片表展示：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_0 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_1
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_1 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.01 sec)

#登录3308示例，查询slice-1分片表展示：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_2 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
Database changed
mysql> show tables;
+----------------------+
| Tables_in_db_mycat_3 |
+----------------------+
| tbl_mycat            |
+----------------------+
1 row in set (0.00 sec)
```

### 插入数据
```shell script
#命令行执行,该命令连接Gaea执行插入：
for i in `seq 1 2000`;do mysql -h127.0.0.1 -P13306 -utest -p1234  db_mycat -e "insert into tbl_mycat (id, col1) values(${i}, 'test$i')";done
```

### 查看数据
```shell script
#连接gaea，进行数据查询：
mysql> use db_mycat
Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|   50 | test50   |
|   51 | test51   |
|   52 | test52   |
...................
| 1996 | test1996 |
| 1997 | test1997 |
| 1998 | test1998 |
| 1999 | test1999 |
+------+----------+
2000 rows in set (0.03 sec)

#连接3307数据库实例，对slice-0分片数据进行查询：
mysql> use db_mycat_0;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|  500 | test500  |
|  501 | test501  |
|  502 | test502  |
|  503 | test503  |
...................
| 1985 | test1985 |
| 1986 | test1986 |
| 1987 | test1987 |
| 1988 | test1988 |
| 1989 | test1989 |
+------+----------+
486 rows in set (0.01 sec)
mysql> use db_mycat_1;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|    1 | test1    |
|    2 | test2    |
|    3 | test3    |
...................
| 1995 | test1995 |
| 1996 | test1996 |
| 1997 | test1997 |
| 1998 | test1998 |
| 1999 | test1999 |
+------+----------+
849 rows in set (0.01 sec)
#连接3308数据库实例，对slice-1分片数据进行查询：
mysql> use db_mycat_2;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|   50 | test50   |
|   51 | test51   |
|   52 | test52   |
...................
| 1594 | test1594 |
| 1595 | test1595 |
| 1596 | test1596 |
| 1597 | test1597 |
| 1598 | test1598 |
| 1599 | test1599 |
| 2000 | test2000 |
+------+----------+
601 rows in set (0.01 sec)

mysql> use db_mycat_3;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> select * from tbl_mycat;
+------+----------+
| id   | col1     |
+------+----------+
|  190 | test190  |
|  191 | test191  |
|  192 | test192  |
|  193 | test193  |
|  194 | test194  |
...................
| 1900 | test1900 |
| 1901 | test1901 |
| 1902 | test1902 |
| 1903 | test1903 |
| 1904 | test1904 |
| 1905 | test1905 |
| 1906 | test1906 |
+------+----------+
64 rows in set (0.00 sec)
```
