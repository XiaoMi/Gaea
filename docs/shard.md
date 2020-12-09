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
分片方式说明：基于分表键的hash值计算子表下标。   
我们想将`db_example`库的`tbl_example`表配置为分片表, 共4个分片, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_example | tbl_example_0000 |
| slice-0 | db_example | tbl_example_0001 |
| slice-1 | db_example | tbl_example_0002 |
| slice-1 | db_example | tbl_example_0003 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_example",
    "table": "tbl_example",
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
配置说明：
-   该配置中的locations字段包含两个元素, locations[0]=2 代表slices字段数组slices[0]包含两个分片,即slice-0的master实例包含两个子表。locations[1]=2 代表slices字段数组slices[1]包含两个分片,即slice-1的master实例包含两个子表。
-   key字段代表用于分表的键。

##### mod
分片方式说明：基于分表键对子表数量的取模运算值计算子表下标。    
我们想将`db_example`库的`shard_mod`表配置为分片表, 共4个分片, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_example | shard_mod_0000 |
| slice-0 | db_example | shard_mod_0001 |
| slice-1 | db_example | shard_mod_0002 |
| slice-1 | db_example | shard_mod_0003 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_example",
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

// ]
```
配置说明：
-   该配置中的locations字段包含两个元素, locations[0]=2 代表slices字段数组slices[0]包含两个分片,即slice-0的master实例包含两个子表。locations[1]=2 代表slices字段数组slices[1]包含两个分片,即slice-1的master实例包含两个子表。
-   key字段代表用于分表的键。

##### range
分片方式说明：基于分表键的所在范围计算子表下标。  
该方式的优点：基于范围的查询或更新速度快，因为查询（或更新）的范围有可能落在同一张子表中。这样可以避免全部子表的查询（更新）。缺点：数据热点问题。因为在一段时间内整个集群的写压力都会落在一张子表上。此时整个mysql集群的写能力受限于单台mysql server的性能。并且，当正在集中写的mysql 节点如果宕机的话，整个mysql集群处于不可写状态。     
我们想将`db_example`库的`tbl_example`表配置为分片表, 共4个分片, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_example | tbl_example_0000 |
| slice-0 | db_example | tbl_example_0001 |
| slice-1 | db_example | tbl_example_0002 |
| slice-1 | db_example | tbl_example_0003 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_example",
    "table": "tbl_example",
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
    "table_row_limit": 100
}

// ]
```

配置说明：
-   该配置中的locations包含两个元素, locations[0]=2 代表slices字段数组slices[0]包含两个分片,即slice-0的master实例包含两个子表。locations[1]=2 代表slices字段数组slices[1]包含两个分片表,即slice-1的master实例包含两个子表。
-   key字段代表用于分表的键。
-   table_row_limit字段的值为100，代表每张子表的记录数。id字段的值为[0,100)在tbl_example_0000上，[100,200)在tbl_example_0001上,依此类推...

##### date_year
分片方式说明：基于分表键日期(年)计算子表下标。  
我们想将`db_example`库的`shard_year`表配置为分片表, 共4个分片, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_example | shard_year_2016 |
| slice-0 | db_example | shard_year_2017 |
| slice-1 | db_example | shard_year_2018 |
| slice-1 | db_example | shard_year_2019 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_example",
    "table": "shard_year",
    "type": "date_year",
    "key": "create_time",
    "slices": [
        "slice-0",
        "slice-1"
    ]
   "date_range": [
        "2016-2017",
        "2018-2019"
    ]
}

// ]
```

配置说明：
-   key:该配置表示shardding key是create_time 
-   data_range:表示shard_year_2016、shard_year_2017两个表在slice-0上，shard_year_2018、shard_year_2019在slice-1上, 左闭右闭。  

gaea 支持Mysql中三种格式的时间类型
-   date类型，格式：YYYY-MM-DD，例如:2016-03-04,注意：2016-3-04，2016-03-4，2016-3-4等格式都是不支持的。
-   datetime，格式：YYYY-MM-DD HH:MM:SS，例如:2016-03-04 13:23:43,注意：2016-3-04 13:23:43，2016-03-4 13:23:43，2016-3-4 13:23:43等格式都是不支持的。
-   timestamp，整数类型。   

注意：子表的命名格式必须是:shard_table_YYYY，shard_table是分表名，后面接具体的年。传入范围必须是有序递增，不能是[2018-2019,2016-2017]，且不能重叠，不能是[2017-2018,2018-2019]。

##### date_month 
分片方式说明：基于分表键日期(月)计算子表下标。   
我们想将`db_example`库的`shard_month`表配置为分片表, 共4个分片, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_example | shard_month_201405 |
| slice-0 | db_example | shard_month_201406 |
| slice-1 | db_example | shard_month_201408 |
| slice-1 | db_example | shard_month_201409 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_example",
    "table": "shard_month",
    "type": "date_month",
    "key": "create_time",
    "slices": [
        "slice-0",
        "slice-1"
    ]
    "date_range": [
         "201405-201406",
         "201408-201409"
    ]
}

// ]
```

配置说明：
-   key: sharding key 是create_time
-   type: 按月的分表类型是date_month
-   data_range: shard_month_201405、shard_month_201406两个子表在slice-0上，shard_month_201408、shard_month_201409在slice-1上,如果一个slice上只包含一张表，可以这样配置date_range[201609,201610-201611]

注意：子表的命名格式必须是:shard_table_YYYYMM,shard_table是分表名，后面接具体的年和月。传入范围必须是有序递增的，不能是[201609-201610,201501]。

##### date_day
分片方式说明：基于分表键日期(天)计算子表下标。     
我们想将`db_example`库的`shard_day`表配置为分片表, 共4个分片, 分布到2个slice上, 每个slice上有1个库, 每个库2张表, 即:

| slice | 后端数据库名 | 后端表名 |
| ----- | ---------- | ------- |
| slice-0 | db_example | shard_day_20201201 |
| slice-0 | db_example | shard_day_20201202 |
| slice-1 | db_example | shard_day_20201203 |
| slice-1 | db_example | shard_day_20201204 |

则namespace配置文件中的分片表规则可参考以下示例配置:

```
// namespace配置文件
// {
// ...
// "shard_rules": [

{
    "db": "db_example",
    "table": "shard_day",
    "type": "date_day",
    "key": "create_time",
    "slices": [
        "slice-0",
        "slice-1"
    ]
    "date_range": [
         "20201201-20201202",
         "20201203-20201204"
    ]
}

// ]
```

配置说明：
-   key: 分表键是create_time
-   type: 按天的分表类型是date_day
-   data_range: 表示shard_day_20201201、shard_day_20201202两个子表在slice-0上，shard_day_20201203、shard_day_20201204在slice-1上, 如果一个slice上只包含一张表，可以这样配置date_range[20160901,20161001-20161101]

注意：子表的命名格式必须是:shard_table_YYYYMMDD,shard_table是分表名，后面接具体的年、月和日。传入范围必须是有序递增的，不能是[20160901-20160902,20150901]。

### mycat分库配置

Gaea支持mycat的常用分库规则, 对应关系如下:

| mycat规则名称          | Gaea规则名称       |
| --------------------- | ----------------- |
| PartitionByMod        | mycat_mod         |
| PartitionByLong       | mycat_long        |
| PartitionByMurmurHash | mycat_murmur      |
| PartitionByString     | mycat_string      |

##### PartitionByMod
分片方式说明：基于分片键对子库数量的取模运算值计算子库下标。   
我们想将`db_mycat`库的`tbl_mycat`表配置为分片表, 共4个分片, 分布到2个slice上面, 每个slice上有2个库, 每个库1张表, 即:

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
    "db_mycat_3",
    "db_mycat_4"
]
```

##### PartitionByLong
分片方式说明：基于分片键固定分片hash算法计算子库下标。  
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
-   该配置中的locations字段包含两个元素, locations[0]=2 代表slices字段数组slices[0]包含两个分片,即slice-0的master实例包含两个子库。locations[1]=2 代表slices字段数组slices[1]包含两个分片,即slice-1的master实例包含两个子库。
-   partition_count标识分片个数，需要与设置的分片数量相等，由于定义了4个分片，因此这里只能是4
-   partition_length代表分片范围列表
    1. 配置"partition_count":"4"、"partition_length"："256"代表希望将数据水平分成4份，每份各占25%
    2. 配置"partition_count":"2,2"、"partition_length"："128,384"代表希望将数据水平分成4份，前两份占128/1024、后两份占384/1024 
    3. 配置"partition_count":"2,1"、"partition_length"："256,512"代表希望将数据水平分成3份，前两份各占25%，第三份占50%  
-   分区长度：默认为最大2^n=1024 ,即最大支持1024分区。
-   约束：1024 = sum((count[i]*length[i])). count和length两个向量的点积恒等于1024。如上述示例中，4 * 256=1024、128 * 2 + 384 * 2=1024、256 * 2 + 512=1024。

##### PartitionByMurmurHash
分片方式说明：基于分片键一致性hash算法计算子库下标。  
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
-   该配置中的locations字段包含两个元素, locations[0]=2 代表slices字段数组slices[0]包含两个分片,即slice-0的master实例包含两个子库。locations[1]=2 代表slices字段数组slices[1]包含两个分片,即slice-1的master实例包含两个子库。
-   其中`seed`, `virtual_bucket_times`配置项的含义与mycat `PartitionByMurmurHash`规则中的同名配置项的含义相同，代表一个实际的数据库节点被映射出该值对应的虚拟节点，这里设置160即虚拟节点数是物理节点数的160倍. 
-   目前Gaea中不支持配置weight, 所有bucket weight都是1.

##### PartitionByString
分片方式说明：基于分片键字串hash值计算子库下标。 
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
    "hash_slice": "-2:0"
}
```

其中`partition_count`, `partition_length`, `hash_slice`配置项的含义与mycat `PartitionByString`规则中的同名配置项的含义相同.
-   该配置中的locations字段包含两个元素, locations[0]=2 代表slices字段数组slices[0]包含两个分片,即slice-0的master实例包含两个子库。locations[1]=2 代表slices字段数组slices[1]包含两个分片,即slice-1的master实例包含两个子库。
-   partition_count代表分片数
-   partition_length代表字符串hash求模基数
-   hash_slice是hash运算位即根据子字符串hash运算.  
其中，hash_slice支持一下格式:  
    -   "2"代表(0,2)
    -   "1:2"代表(1,2)
    -   "1:"代表(1,0)
    -   "-1:"代表(-1,0)
    -   ":-1"代表(0,-1)  
例1：值“45abc”，hash运算位0:2 ，取其中45进行计算   
例2：值“aaaabbb2345”，hash预算位-4:0 ，取其中2345进行计算  

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
