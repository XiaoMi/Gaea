# gaea proxy后端连接池的设计与实现

## 理想的连接池

基于go实现连接池的方式有很多种，比如通过chan、通过map+锁的方式，但是从使用者的角度来看，一个优秀的连接池我认为有以下几个特性: 1.有最大连接数和初始连接数限制 2.实际连接数可以在上述范围内动态伸缩 3.有限时间内获取连接且连接可用。在这三个基础之上，可能还包含一些其他非必须特性比如在运行时改变连接数最大限制、暴露连接池的一些状态信息等等。gaea的连接池是基于vitess的resource pool进行封装并添加了连接重试的功能，相关代码在backend和util目录下。

## 连接池的创建、使用

### 定义

ConnectionPool定义

```golang
// ConnectionPool means connection pool with specific addr
type ConnectionPool struct {
    mu          sync.RWMutex
    connections *util.ResourcePool

    addr     string
    user     string
    password string
    db       string

    charset     string
    collationID mysql.CollationID

    capacity    int // capacity of pool
    maxCapacity int // max capacity of pool
    idleTimeout time.Duration
}
```

NewConnectionPool定义

```golang
// NewConnectionPool create connection pool
func NewConnectionPool(addr, user, password, db string, capacity, maxCapacity int, idleTimeout time.Duration, charset string, collationID mysql.CollationID) *ConnectionPool {
    cp := &ConnectionPool{addr: addr, user: user, password: password, db: db, capacity: capacity, maxCapacity: maxCapacity, idleTimeout: idleTimeout,   charset: charset, collationID: collationID}
    return cp
}
```

Open定义

```golang
// Open open connection pool without error, should be called before use the pool
func (cp *ConnectionPool) Open() {
    if cp.capacity == 0 {
        cp.capacity = DefaultCapacity
    }

    if cp.maxCapacity == 0 {
        cp.maxCapacity = cp.capacity
    }
    cp.mu.Lock()
    defer cp.mu.Unlock()
    cp.connections = util.NewResourcePool(cp.connect, cp.capacity, cp.maxCapacity, cp.idleTimeout)
    return
}
```

每一个连接作为一个资源单位，connections存放所有的连接资源，类型为ResourcePool。capacity定义连接池的初始容量、maxCapacity定义连接池的最大容量，实际后端连接数量会根据使用情况在[0,maxCapacity]之间浮动。idleTimeout为连接空闲关闭时间，当连接不活跃时间达到该值时，连接池将会与后端mysql断开该连接并回收资源。

### 创建

外部调用者通过NewConnectionPool函数创建一个连接池对象，然后通过Open函数，初始化连接池并与后端mysql建立实际的连接(有的连接池或长连接会在第一次请求时才去建立连接)。连接池的connect方法作为资源池初始化的工厂方法，所以如果你要基于该resource pool实现其他池子时，需要实现该工厂方法。

### 使用

通过连接池的Get方法可以获取一个连接，Get的入口参数包含context，该context最初设计用来传入一些全局上下文，比如超时上下文，目前gaea的获取连接超时时间是固定的，所以超时上下文也是基于该context在内部构造的。为防止上层一直阻塞在获取连接处，超过getConnTimeout未获取到连接会报超时错误，从而避免发生更严重的状况。如果发生超时次数过多，可以通过配置平台调整最大连接数大小，不同namespace的最大连接数其实是一个经验值，根据不同的业务状态有所不同，支持动态实时调整。

拿到连接后，连接池会主动调用tryReuse，用来保证连接是自动提交的状态。应用层需要主动调用initBackendConn，初始化与mysql有关的状态信息，包括use db、charset、session variables、sql mode等。为了提升效率会检测后端连接与前端会话的对应变量值，不一致才会进行设置，在设置的时候也是通过批量发送sql的方式，尽最大可能减少与后端mysql的网络交互。

连接使用完成后，需要手动调用recycleBackendConn回收连接，注意: 事务相关的连接是在commit或者rollback的时候进行统一释放。

## 动态维护连接

动态维护连接其实包含两部分，一部分维护连接池容量，不活跃的连接要主动关闭。连接池的Open函数会调用NewResourcePool函数，NewResourcePool函数会启动一个Timer定时器，定时器通过定期检测连接活跃时间与空闲时间的差值，决定是否关闭连接、回收资源从而实现动态调整连接池的容量。另一部分是保证获取到的连接是有效连接，这里在通过writeEphemeralPacket向后端mysql连接写数据时，如果报错含有"broken pipe"即连接可能无效，则会进行重试，直到连接成功或者重试次数达到三次报错，通过主动监测和连接重试，我们不需要进行定期ping后端连接，就可以保证后端连接是有效的。

## 总结

gaea的连接池实现还是相对简单可依赖的，在使用gaea的过程中，最好将后端mysql的wait_timeout值设置为比gaea idleTimeout长，可以减少不必要的连接重试。
