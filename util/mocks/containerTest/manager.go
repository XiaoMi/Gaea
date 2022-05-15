package containerTest

import (
	"errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run"
	"os"
	"runtime"
	"strings"
	"sync"
)

// 这次 PR 会把这里删除，不送出
var (
/*cmConfigFile = flag.String("cm", "./etc/containerd.ini", "containerd manager 配置")
  // DefaultConfigPath 缺省的容器设置路径
  DefaultConfigPath = "./etc/containerd/"*/
)

// ContainerManager 容器服务管理员
// ContainerManager is used to manage Containerd.
type ContainerManager struct {
	Enable        bool                      // 是否启用容器服务管理器 enable containerd manager
	ConfigPath    string                    // 容器服务配置路径 config path
	ContainerList map[string]*ContainerList // 容器服务列表 containerd list
}

// ContainerList 容器服务列表
// ContainerList is used to list containerd clients.
type ContainerList struct {
	NetworkLock sync.Locker                // 容器服务网络锁 containerd network lock
	Cfg         containerd.ContainerConfig // 容器服务配置 containerd config
	Builder     builder.Builder            // 容器服务构建器 containerd builder
	User        string                     // 用户名称 User name
	Status      int                        // 容器服务状态 containerd status
}

// NewContainderManager 新建容器服务管理员
// NewContainderManager is used to create a new container manager.
func NewContainderManager(path string) (*ContainerManager, error) {
	// 判断目录是否存在
	// check directory is existed
	if strings.TrimSpace(path) == "" {
		path = DefaultConfigPath
	}
	if err := checkDir(path); err != nil {
		log.Warn("check file config directory failed, %v", err)
		return nil, err
	}

	// 开始载入管理员的所有设定档
	// start to load all config files of manager
	r := Load{prefix: path}
	configs, err := r.loadAllContainerD()
	if err != nil {
		log.Warn("load containerd config failed, %v", err)
		return nil, err
	}

	// 产生管理员的列表
	// create manager's list
	containerList := make(map[string]*ContainerList)
	for container, config := range configs {
		newBuilder, err := containerd.NewBuilder(config)
		if err != nil {
			_ = log.Warn("make containerd client failed, %v", err)
			return nil, err
		}
		containerList[container] = &ContainerList{
			User:        "",                       // 获取函数名称 register's function (到时会改成执行函数名称 be changed by register's function)
			Status:      run.ContainerdStatusInit, // 容器服务状态为占用 containerd is occupied
			Cfg:         config,                   // 容器服务配置 containerd config
			Builder:     newBuilder,               // 容器服务构建器 containerd builder
			NetworkLock: &sync.Mutex{},            // 容器服务网络锁 containerd network lock
		}
	}

	// 回传管理员的对象
	// return manager's object
	return &ContainerManager{ConfigPath: path, ContainerList: containerList}, nil
}

// registerFunc 注册函数名称 register's function
type registerFunc func() string

// 预设的注册函数
// default register's function
func defaultRegFunction() string {
	// 以下会返回执行函数名称
	// return register's function
	counter, _, _, success := runtime.Caller(2)
	if !success {
		return "unknown"
	}
	return runtime.FuncForPC(counter).Name()
}

// AppendCurrentFunction 返回当前函数名
// AppendCurrentFunction returns the current function name.
func AppendCurrentFunction(layerNumber int, appendStr string) string {
	// appendStr 是用来判别各个协程
	// appendStr is to identify each goroutine's function name

	// 以下会返回执行函数名称
	// return register's function
	counter, _, _, success := runtime.Caller(layerNumber)
	if !success {
		return "unknown"
	}
	return runtime.FuncForPC(counter).Name() + appendStr
}

// checkDir 检查目录是否存在 checkDir is used to check directory is existed.
func checkDir(path string) error {
	// 先修正路径参数
	// fix parameter path
	if strings.TrimSpace(path) == "" {
		return errors.New("invalid path")
	}

	// 判断目录是否存在
	// check directory is existed
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("invalid path, should be a directory")
	}

	return nil
}

// isEnabled 获取容器是否被启用
// isEnabled is to check container is enabled.
func (cm *ContainerManager) isEnabled() bool {
	return cm.Enable
}

// getBuilder 获取容器服务构建器
// getBuilder is used to get containerd builder.
func (cm *ContainerManager) getBuilder(containerName string, regFunc registerFunc) (builder.Builder, error) {
	// 判断容器服务是否存在，不存在则使用默认的注册函数
	// check containerd is existed, if not, use default register's function
	if regFunc == nil {
		regFunc = defaultRegFunction
	}

	// 如果没有配置，则返回错误
	// if no config, return error
	if _, ok := cm.ContainerList[containerName]; !ok {
		return nil, errors.New("invalid config container name")
	}

	// 如果可以进行占用，则续续以下操作
	// if we can occupy, then continue to do the following operation.
	cm.ContainerList[containerName].NetworkLock.Lock()                    // 加锁 lock
	cm.ContainerList[containerName].User = regFunc()                      // 获取函数名称 register's function
	cm.ContainerList[containerName].Status = run.ContainerdStatusOccupied // 容器服务状态为占用 containerd is occupied

	// 记录日志 log
	_ = log.Notice(cm.ContainerList[containerName].User + " occupies " + containerName)

	// 正常返回容器服务构建器
	// return containerd builder
	return cm.ContainerList[containerName].Builder, nil
}

// returnBuilder 适放容器服务构建器
// returnBuilder is used to release containerd builder.
func (cm *ContainerManager) returnBuilder(containerName string, regFunc registerFunc) error {
	if regFunc == nil {
		regFunc = defaultRegFunction
	}

	// 如果没有配置，则返回错误
	if _, ok := cm.ContainerList[containerName]; !ok {
		return errors.New("invalid config container name")
	}

	// 如果可以进行占用，则续续以下操作
	// if we can occupy, then continue to do the following operation.
	cm.ContainerList[containerName].NetworkLock.Unlock()                  // 解锁 unlock
	cm.ContainerList[containerName].User = ""                             // 获取函数名称 register's function
	cm.ContainerList[containerName].Status = run.ContainerdStatusReturned // 容器服务状态为被适放 containerd is released

	// 记录日志 log
	_ = log.Notice(cm.ContainerList[containerName].User + " releases " + containerName)

	// 正常适放容器服务构建器
	// release containerd builder
	return nil
}

// getIPAddrPort 获取容器的网络位置
// getIPAddrPort is used to get containerd's status.
func (cm *ContainerManager) getIPAddrPort(containerName string) string {
	return cm.ContainerList[containerName].Cfg.IP
}
