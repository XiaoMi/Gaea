# 多租户的设计与实现

## 背景

gaea多租户是为了实现一套gaea集群，可以接入多个业务系统的不同数据库，方便部署、运维。gaea多租户为软多租户，一个租户称为一个namespace，多个namespace之间存在于一套gaea proxy集群内，所以是一种软隔离。我们也可以为一些重要等级业务系统单独部署一套gaea集群，甚至一套业务系统对应一套gaea集群实现物理隔离。

## 接入方式

mysql的授权方式为用户名+密码+ip+数据库，接入gaea的情况下，授权方式为用户名+密码确定唯一一个namespace，ip则在白名单IP/IP段起作用，如果未配置白名单IP，则默认对所有IP生效。所以，不同的业务系统，用户名可以有相同的，但是用户名+密码要保证是唯一的，密码我们内部是根据一定的规则随机生成的，并且会校验是否重复。

## 实现原理

### 主要结构

在授权阶段还未能确定对应的namespace，所以user和namespace的配置是分别加载的，授权的实现主要依赖于UserManager结构体，其定义如下

```golang
type UserManager struct {
    users          map[string][]string // key: user name, value: user password, same user may have different password, so array of passwords is needed
    userNamespaces map[string]string   // key: UserName+Password, value: name of namespace
}
```

### 配置加载过程

在系统初始化阶段，会依次加载对应user、namespace配置到一个全局Manager内，其中user部分用以授权检查，整个配置通过滚动数组的方式实现了无锁热加载，具体实现可以参照[gaea配置热加载实现原理](config-reloading.md)一章。

### 校验过程

其中users同一个用户名对应一个string数组，用以处理同一用户名不同密码的情形，而userNamespaces则用以通过用户名+密码快速获取对应的namespace名称。在验证阶段，首先通过CheckUser检查用户名是否存在，不存在则直接授权失败。然后，通过CheckPassword，依次对比确定是否可以找到对应密码，如果找不到，则最终授权失败；如果找到，则授权检查通过并记录对应的会话信息。

## 结语

gaea的多租户确实为部署、运维带来了不少方便，后续也会考虑支持kubernetes的部署、调度多租户等，但是当下的多租户结构不会发生太大变化。