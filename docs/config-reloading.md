# gaea配置热加载设计与实现

## 背景

我们在开始设计gaea的时候列了几个关键词: 配置热加载、多租户、读写分离、分库分表、路由sql，配置热加载就是其中一个非常重要的目标。以往基于配置文件的方式，存在以下几个问题: DBA
使用起来比较容易出错、多个代理的配置文件维护比较麻烦、配置更新到生效的过程比较漫长、通过接口更新配置也存在文件和内存中配置可能不一致的情况。而gaea配置热加载就是希望能解决以上弊端并通过一个统一的平台进行配置管理。

设计之初，我们就把配置分为静态配置和动态配置。静态配置即在运行时不需要、也不可以进行修改的配置内容，比如监听端口、etcd地址、log路径等。动态配置即为运行时需要不断进行更改的配置，比如增删namespace、namespace内实例、用户、配置项的变更等等，动态配置是配置热加载的主角。

## 配置结构

动态配置最开始加载为models里对应的结构，比如models.Namespace，但是程序运行使用的配置为proxy/server包下的Manager。

动态配置对应结构如下:

```golang
// Manager contains namespace manager and user manager
type Manager struct {
    switchIndex util.BoolIndex
    namespaces  [2]*NamespaceManager
    users       [2]*UserManager
    statistics  *StatisticManager
}
```

Manager是一个全局变量，负责动态配置的管理，按照配置的适用范围分为用户配置、namespace配置和打点统计配置，而switchIndex则作为滚动数组进行配置滚动的标识。

### Namespace配置

主要包含名称、数据库白名单、用户属性、sql指纹、执行计划、分片配置等。

```golang
// NamespaceManager namespace manager
type NamespaceManager struct {
    namespaces map[string]*Namespace
}
// Namespace is struct driected used by server
type Namespace struct {
    name               string
    allowedDBs         map[string]bool
    defaultPhyDBs      map[string]string // logicDBName-phyDBName
    sqls               map[string]string //key: sql fingerprint
    slowSQLTime        int64             // session slow sql time, millisecond, default 1000
    allowips           []util.IPInfo
    router             *router.Router
    slices             map[string]*backend.Slice // key: slice name
    userProperties     map[string]*UserProperty  // key: user name ,value: user's properties
    defaultCharset     string
    defaultCollationID mysql.CollationID

    slowSQLCache         *cache.LRUCache
    errorSQLCache        *cache.LRUCache
    backendSlowSQLCache  *cache.LRUCache
    backendErrorSQLCache *cache.LRUCache
    planCache            *cache.LRUCache
}
```

### UserManager配置

主要包含auth阶段所需要的结构化信息

```golang
// UserManager means user for auth
// username+password是全局唯一的, 而username可以对应多个namespace
type UserManager struct {
    users          map[string][]string // key: user name, value: user password, same user may have different password, so array of passwords is needed
    userNamespaces map[string]string   // key: UserName+Password, value: name of namespace
}
```

### StatisticManager配置

作为Manager里的一个子模块，用以打点统计各项指标，并非配置热加载的配置项，在此不进行过多讨论。

## 配置初始化

通过InitManager，加载所有的namespace，初始化users、namespaces、statistics，构造全局manager(globalManager)

## 配置变更接口

gaea的配置变更没有针对每个配置项分别进行接口封装，而是基于配置完整替换+被替换配置动态资源延迟关闭的策略。无论是某一个还是多个配置项发生变化或者新增、删除namespace，对于gaea来说，都会构建一份新的配置并构造后端的动态资源。当新配置滚动为当前使用的配置之后，旧版本的动态资源在延迟60秒之后主动释放，进行比如关闭套接字、文件描述符等的操作。所以gaea的配置变更可以理解为，只要配置项实现了构造、延迟关闭功能，都可以纳入配置热加载模块内，进行统一管理。

## 滚动数组实现无锁化

滚动数组是配置热加载过程中经常使用的技巧，通过用空间换时间的方式，规避配置获取和配置变更之间的竟态条件。Manager中的switchIndex即为当前生效配置的下标，switchIndex的类型为BoolIndex，其对应的Get、Set方法均为原子操作，这样就保证了二元数组对应配置的切换为原子过程，而配置复制和赋值永远都是变更当前未在使用的另一元素，即通过写时复制+原子切换实现了配置的无锁化。为了防止多次滚动，导致丢失配置，在进行配置更改时，需要持有全局锁，保证同一时间，只有一个namespace进行配置变更。

## 延迟关闭回收动态资源

配置在提交之后，新配置生效，老配置需要进行资源回收。通过一个单独的goroutine，在sleep一段时间之后(尽最大努力保证请求得到应答)，调用各个单独项的Close，回收资源。

## 两阶段提交保证一致性

一个集群会包含多台gaea-proxy，为了保证多台gaea-proxy快速生效相同的配置，故而引入了两阶段提交的配置变更方式，其中协调者为gaea-cc。第一阶段: gaea-cc调用各个gaea-proxy的prepare接口，gaea-proxy在prepare阶段首先复制一份当前的全量配置，然后从etcd加载对应namespace的最新的配置，最后更新对应的全量配置；第二阶段: gaea-cc如果在prepare阶段发生错误(任何一个gaea-proxy报错)则直接报错，prepare成功后则调用gaea-proxy的commit接口，gaea-proxy在commit接口只进行一次简单的配置切换，这样prepare工作重、commit工作非常轻量，可以很大程度上提升配置变更成功的几率。如果commit失败，则gaea-cc也是直接报错，对应的web平台上看到错误后可以决定是否停止变更或者重新发起一次变更(多次发送相同配置幂等)。

## 集群配置一致性校验

通过两阶段提交配置后，当前所有gaea-proxy的生效配置是相同的。为了方便验证: 1.配置是否发生变化 2.是否所有gaea-proxy的最新配置已经生效，gaea-proxy提供了获取当前配置签名的接口。通过该接口，DBA可以直接通过管理平台查看到各个gaea-proxy前后及当前配置的md5签名，保证配置变更的执行效果符合预期。