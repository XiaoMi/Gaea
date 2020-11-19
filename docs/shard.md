# 分片表配置

Gaea支持kingshard分表规则和mycat分库规则, 用户可以在不迁移任何数据的情况下, 从kingshard和mycat切换到Gaea.

### kingshard分表配置

Gaea支持kingshard常用分表规则, 对应关系如下:

| kingshard规则名称 | Gaea规则名称 |
| ---------------- | ---------- |
| hash             | hash       |
| mod              | mod        |
| range            | range      |
| date_year        | date_year  |
| date_month       | date_month |
| date_day         | date_day   |

##### hash

我们想将`db_kingshard`库的`tbl_kingshard`表配置为分片表, 共4个分片表, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_kingshard | tbl_kingshard_0000 |
| slice-0 | db_kingshard | tbl_kingshard_0001 |
| slice-1 | db_kingshard | tbl_kingshard_0002 |
| slice-1 | db_kingshard | tbl_kingshard_0003 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_kingshard",
    "table": "tbl_kingshard",
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

// ]
```

### mycat分库配置

Gaea支持mycat的常用分库规则, 对应关系如下:

| mycat规则名称          | Gaea规则名称       |
| --------------------- | ----------------- |
| PartitionByMod        | mycat_mod         |
| PartitionByLong       | mycat_long        |
| PartitionByMurmurHash | mycat_murmur      |
| PartitionByString     | mycat_string      |

##### PartitionByMod

我们想将`db_mycat`库的`tbl_mycat`表配置为分片表, 共4个分片表, 分布到2个slice上面, 每个slice上有2个库, 每个库1张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_mycat_0 | tbl_mycat |
| slice-0 | db_mycat_1 | tbl_mycat |
| slice-1 | db_mycat_2 | tbl_mycat |
| slice-1 | db_mycat_3 | tbl_mycat |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

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

// ]
```

注: databases配置项, 指定了每个分片表的实际库名, 这里采用了简写的方式, 与以下配置等价:

```
"databases": [
    "db_mycat_0",
    "db_mycat_1",
    "db_mycat_2",
    "db_mycat_3"
]
```

如果指定的实际库名不是递增的, 也可以手动指定, 如:

```
"databases": [
    "db_mycat_0",
    "db_mycat_1",
    "db_mycat_0",
    "db_mycat_1"
]
```

##### PartitionByLong

mycat_long的配置规则如下:

```
{
    "db": "db_mycat",
    "table": "tbl_mycat_long",
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
```

其中`partition_count`, `partition_length`配置项的含义与mycat `PartitionByLong`规则中的同名配置项的含义相同.

##### PartitionByMurmurHash

mycat_murmur的配置规则如下:

```
{
    "db": "db_mycat",
    "table": "tbl_mycat_murmur",
    "type": "mycat_murmur",
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
        "db_mycat_0","db_mycat_1","db_mycat_2","db_mycat_3"
    ],
    "seed": "0",
    "virtual_bucket_times": "160"
}
```

其中`seed`, `virtual_bucket_times`配置项的含义与mycat `PartitionByMurmurHash`规则中的同名配置项的含义相同. 而在mycat中需要指定的`count`配置项, 在Gaea中通过locations自动判断, 不需要手动指定.

目前Gaea中不支持配置weight, 所有bucket weight都是1.

##### PartitionByString

mycat_string的配置规则如下:

```
{
    "db": "db_mycat",
    "table": "tbl_mycat_string",
    "type": "mycat_string",
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
    "partition_length": "256",
    "hash_slice": "20"
}
```

其中`partition_count`, `partition_length`, `hash_slice`配置项的含义与mycat `PartitionByString`规则中的同名配置项的含义相同.

### 关联表和全局表

Gaea分片SQL要求多个表具有关联关系 (一个分片表, 多个关联表), 或者只存在一个分片表, 其余均为全局表.

##### 关联表

关联表需要与同数据库的某个分片表关联. 例如, 将`db_mycat`库的`tbl_mycat_child`配置为关联表, 分片列名为`id`, 父表为`tbl_mycat`, 则需要在`shard_rules`中添加以下配置:

```
{
    "db": "db_mycat",
    "table": "tbl_mycat_child",
    "type": "linked",
    "parent_table": "tbl_mycat",
    "key": "id"
}
```

此时即可执行分片内的关联查询:

```
SELECT * FROM tbl_mycat, tbl_mycat_child WHERE tbl_mycat_child.id=5 AND tbl_mycat.user_name='hello';
```

##### 全局表

全局表是在各个slice上 (准确的说是各个slice的各个DB上) 数据完全一致的表, 方便执行一些跨分片查询, 配置如下:

```
{
    "db": "db_mycat",
    "table": "tbl_mycat_global",
    "type": "global",
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
```
