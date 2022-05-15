package containerd

import (
	"context"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run/defaults"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run/etcd"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run/mariadb"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"time"
)

// >>>>> >>>>> >>>>> contianerd 客户端的设置结构体
// ContainerdClient's config struct.

// ContainerConfig 为 Containerd 容器服务所要 Parse 对应的内容
type ContainerConfig struct {
	Sock      string `json:"sock"`      // 容器服务在 Linux 的连接位置 the connection location of the container service.
	Type      string `json:"type"`      // 容器服务类型 the type of the container service.
	Name      string `json:"name"`      // 容器服务名称 the name of the container service.
	NameSpace string `json:"namespace"` // 容器服务命名空间 the namespace of the container service.
	Image     string `json:"image"`     // 容器服务镜像 the image of the container service.
	Task      string `json:"task"`      // 容器服务任务名称 the task name of the container service.
	NetworkNs string `json:"networkNs"` // 容器服务网络 namespace of the container service.
	IP        string `json:"ip"`        // 容器服务 IP。 the IP of the container service.
	SnapShot  string `json:"snapshot"`  // 容器服务快照 the snapshot of the container service.
	Schema    string `json:"schema"`    // 容器服务 Schema，用于数据库设置 the schema of the container service.
	User      string `json:"user"`      // 容器服务用户名 the user name of the container service.
	Password  string `json:"password"`  // 容器服务密码 the password of the container service.
}

// NewBuilder 为创建一个新的 Builder 接口
// NewBuilder is a function to create a new Builder interface.
func NewBuilder(cfg ContainerConfig) (builder.Builder, error) {
	return NewContainerdClient(cfg)
}

// NewContainerdClient 为新建容器服务的客户端
// NewContainerdClient is a function to create a new containerd client.
func NewContainerdClient(cfg ContainerConfig) (*ContainerdClient, error) {
	// >>>>> >>>>> >>>>> 决定容器服务的连接 sock 对象
	// >>>>> >>>>> >>>>> decide the sock.

	// 创建容器服务的客户端 Socket 对象
	// create a new sock for the containerd client.
	currentSock := new(containerd.Client) // 新的连接 a new connection to containerd.
	var err error                         // 报错信息 error message
	var usedSock = ""                     // 客户端的 sock 指定路径 usedSock is the user defined sock path.

	// 如果没有配置，则使用默认的路径
	// if socketPath is empty, use default path.
	if cfg.Sock != "" { // 使用指定的路径 use the specified path.
		usedSock = cfg.Sock
	} else { // 使用默认的路径 use default path.
		usedSock = run.DefaultSock
	}

	// 创建容器服务的客户端连接
	// create a new containerd connection.
	currentSock, err = containerd.New(usedSock)
	if err != nil {
		return nil, err
	}

	// >>>>> >>>>> >>>>> 创建容器服务的 ContainerdClient 对象 create a new ContainerdClient object

	// 创建容器服务的客户端
	// create a new containerd client.
	client := &ContainerdClient{
		// 在 NewContainerdClient 中创建
		// create in NewContainerdClient.

		// 客户端的配置 client's config.
		Status: run.ContainerdStatusInit, // 现在为初始化状态 init status.
		Conn:   currentSock,              // 容器服务的连接 a connection to containerd.
		Type:   cfg.Type,                 // 容器服务的类型 a containerd type.
		IP:     cfg.IP,                   // 容器服务的网路位置 a containerd IP.
		// 容器的配置 container's config.
		Container: ClientContainerd{
			Name:      cfg.Name,      // 容器的名称 a container's name.
			NameSpace: cfg.NameSpace, // 容器服务的名称 a container name.
			Image:     cfg.Image,     // 容器服务的镜像 a container image.
			NetworkNS: cfg.NetworkNs, // 容器服务的网络命名空间 a container network namespace.
			SnapShot:  cfg.SnapShot,  // 容器服务的快照 a container snapshot.
		},
		// Schema的配置 Schema's config.
		Schema: ClientSchema{
			User: cfg.User, // 容器服务的用户 containerd user.
		},
	}

	// >>>>> >>>>> >>>>> 实现 Run 接口并回传 implement the Run interface and return

	// 实现在 Distinguish 中
	// implement in Distinguish.
	err = client.Distinguish()
	if err != nil {
		return nil, err
	}

	// 返回新的容器服务的客户端
	// return the new containerd client.
	return client, nil
}

// >>>>> >>>>> >>>>> containerd client 的内核结构体
// >>>>> >>>>> >>>>> ContainerdClient is a struct to represent a containerd client.

// ContainerdClient 为容器服务的内核客户端
// ContainerdClient is core component of Containerd client.
type ContainerdClient struct {
	// 在 NewContainerdClient 中创建 create in NewContainerdClient.
	Status    int                // 容器服务的状态 containerd client's status.
	Conn      *containerd.Client // 容器服务的连接 containerd client's connection.
	Type      string             // 容器类型，是 etcd 或者 mariaDB 等等 containerd client's type.
	IP        string             // 容器服务的 IP containerd client's IP.
	Container ClientContainerd   // 容器服务的容器 containerd client's container.
	Schema    ClientSchema       // 容器服务的 Schema containerd client's schema.
	Running   *ClientRunning     // 容器服务运行中的暂时资料 containerd client's running data.

	// 在容器区分时候创建 create in Distinguish.
	Run run.Run // 容器服务的运行的接口 interface for Containerd.
}

// ClientContainerd 为客户端的容器服务设置
// ClientContainerd is configured for Containerd.
type ClientContainerd struct {
	Name      string // 容器服务的名称 Name
	NameSpace string // 容器服务的命名空间 NameSpace
	Image     string // 容器服务的镜像 Image
	SnapShot  string // 容器服务的快照 SnapShot
	NetworkNS string // 容器服务的网络命名空间 NetworkNS
	Container string // 容器服务的容器 Container
	Task      string // 容器服务的任务 Task
}

// ClientSchema 客户端的 Schema 设置
// ClientSchema is the schema of the containerd client.
type ClientSchema struct {
	User     string // 容器服务的用户名 user name.
	Password string // 容器服务的密码 password.
	Schema   string // 容器服务的 schema.
}

// ClientRunning 客户端的运行时的对象
// ClientRunning is the running object of the containerd client.
type ClientRunning struct {
	ctx    context.Context      // 容器服务的上下文 context.
	img    containerd.Image     // 容器服务的镜像 image.
	c      containerd.Container // 容器服务的容器 container.
	tsk    containerd.Task      // 容器服务的任务 task.
	cancel context.CancelFunc   // 容器服务的取消函数 cancel function.
}

// Distinguish 对容器服务进行区分，判断容器服务的类型，给容器服务所需要的功能
// Distinguish is a function to distinguish the containerd client.
func (cc *ContainerdClient) Distinguish() error {
	// distinguish the containerd client. 对容器服务进行区分
	switch cc.Type {
	case "etcd":
		cc.Run = new(etcd.Etcd) // use etcd. 容器服务为 etcd
		return nil              // return nil. 返回 nil
	case "mariadb": // use mariaDB. 容器服务为 mariaDB
		cc.Run = new(mariadb.MariaDB) // return mariaDB. 返回 mariaDB
		return nil                    // return nil. 返回 nil
	default:
		cc.Run = new(defaults.Defaults) // use defaults. 容器服务为 defaults
		return nil                      // return nil. 返回 nil
	}
}

// Build 创建容器测试环境
// Build create a new container environment for test.
func (cc *ContainerdClient) Build(t time.Duration) error {
	// 错误信息 error message.
	var err error

	// 设置容器管理器为创建状态
	// container manager's status is building.
	cc.Status = run.ContainerdStatusBuilding

	// 创建 Running 对象
	// create Running object.
	cc.Running = new(ClientRunning)

	// 决定之后的决行时间
	// decide the time duration.
	ctx := context.Background() // 创建一个上下文对象 create a context object.
	if t > 0 {
		ctx, cc.Running.cancel = context.WithTimeout(ctx, t)
	}

	// 测立一个新的命名空间
	// create a new context with a namespace
	cc.Running.ctx = namespaces.WithNamespace(ctx, cc.Container.NameSpace)

	// 设置容器管理器为下载镜像状态
	// container manager's status is pulling image.
	cc.Status = run.ContainerdStatusBuildPullingImage

	// 拉取预设的测试印象档
	// pull the default test image from DockerHub
	// example: "docker.io/panhongrainbow/mariadb:testing" OR "localhost/mariadb:latest"
	cc.Running.img, err = cc.Run.Pull(cc.Conn, cc.Running.ctx, cc.Container.Image)
	if err != nil {
		return err
	}

	// 设置容器管理器为创建容器状态
	// container manager's status is creating container.
	cc.Status = run.ContainerdStatusBuildCreateContainer

	// 创建一个新的容器
	// create a new container
	cc.Running.c, err = cc.Run.Create(cc.Conn, cc.Running.ctx, cc.Container.Name, cc.Container.NetworkNS, cc.Running.img, cc.Container.SnapShot)
	if err != nil {
		return err
	}

	// 设置容器管理器为创建容器任务状态
	// container manager's status is creating container task.
	cc.Status = run.ContainerdStatusBuildCreateTask

	// 创建新的容器工作
	// create a task from the container
	cc.Running.tsk, err = cc.Run.Task(cc.Running.c, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 设置容器管理器为启动容器任务状态
	// container manager's status is starting container task.
	cc.Status = run.ContainerdStatusBuildStartTask

	// 开始执行容器工作
	// start the task.
	err = cc.Run.Start(cc.Running.tsk, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 设置容器管理器为容器任务执行状态
	// container manager's status is running container task.
	cc.Status = run.ContainerdStatusBuildRunning

	// 创建容器环境成功
	// build the container environment successfully.
	return nil
}

// OnService 确认容器是否处放服务状态
// OnService is making sure the container is on service.
func (cc *ContainerdClient) OnService(t time.Duration) error {
	// 设置容器管理器为检查连线状态
	// container manager's status is checking on service.
	cc.Status = run.ContainerdStatusChecking

	// 决定之后的决行时间
	// decide the time duration.
	ctx := context.Background() // 创建一个上下文对象 create a context object.
	if t > 0 {
		ctx, cc.Running.cancel = context.WithTimeout(ctx, t)
	}

	// 测立一个新的命名空间
	// create a new context with a "mariadb" namespace
	cc.Running.ctx = namespaces.WithNamespace(ctx, cc.Container.NameSpace)

	// 设置容器管理器为检查连线服务
	// container manager's status is checking on service.
	cc.Status = run.ContainerdStatusCheckOnService

	// 检查服务是否正常
	// check if the container is on service.
	err := cc.Run.CheckService(cc.Running.ctx, cc.IP)
	if err != nil {
		return err
	}

	// 设置容器管理器为检查数据库资料完整中
	// container manager's status is checking database data.
	cc.Status = run.ContainerdStatusCheckSchema

	// 检查数据库资料是否完整
	// check if the database is complete.
	err = cc.Run.CheckSchema(cc.Running.ctx, cc.IP)
	if err != nil {
		return err
	}

	// 设置容器管理器为服务正常中
	// container manager's status is service ready.
	cc.Status = run.ContainerdStatusReady

	// 容器检查完成
	// Container check completed.
	return nil
}

// TearDown 拆除容器测试环境
// TearDown is a function to tear down the container environment.
func (cc *ContainerdClient) TearDown(t time.Duration) error {
	// 设置容器管理器为拆除容器状态
	// container manager's status is tearing down container.
	cc.Status = run.ContainerdStatusTearingDown

	// 决定之后的决行时间
	// decide the time duration.
	ctx := context.Background() // 创建一个上下文对象 create a context object.
	if t > 0 {
		ctx, cc.Running.cancel = context.WithTimeout(ctx, t)
	}

	// 测立一个新的命名空间
	// create a new context with a namespace
	cc.Running.ctx = namespaces.WithNamespace(ctx, cc.Container.NameSpace)

	// 设置容器管理器为中断容器状态
	// container manager's status is interrupting container.
	cc.Status = run.ContainerdStatusTearDownInterrupted

	// 强制中断容器工作
	// interrupt the task.
	err := cc.Run.Interrupt(cc.Running.tsk, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 设置容器管理器为删除容器
	// container manager's status is containerd killed.
	cc.Status = run.ContainerdStatusTearDownKilled

	// 删除容器和获得离开讯息
	// kill the process and get the exit status
	err = cc.Run.Delete(cc.Running.tsk, cc.Running.c, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 取消上下文
	// cancel the context.
	defer cc.Running.cancel()

	// 清空 Running 对象
	// clear the Running object.
	cc.Running = nil

	// 删除容器环境成功
	// delete the container environment successfully.
	return nil
}
