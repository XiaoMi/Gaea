package defaults

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
	"time"
)

// Defaults 为预设容器配置项
// default container configuration.
type Defaults struct {
	debianVersion string // Debian 版本 Debian version.
}

// 约定的接口
// defined interface

// >>>>> >>>>> >>>>> 创建部分 create part

// Pull 为容器拉取镜像
// Pull is to pull image from registry.
func (d *Defaults) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	// 从注册中心拉取镜像 pull image from registry.
	download := func(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
		image, err := client.Pull(ctx, imageUrl, containerd.WithPullUnpack)
		if err != nil {
			return nil, err
		}
		return image, nil
	}

	// 之后会用协程去拉取镜像，创建通信信道.
	// Then, use goroutine to pull image. Create communication channel.

	// message. 消息
	type message struct {
		image containerd.Image
		err   error
	}

	// channel. 信道
	chMessage := make(chan message)

RETRY:
	// 开启一个goroutine从注册中心拉取镜像
	// goroutine to pull image from registry.
	go func(client *containerd.Client, ctx context.Context, imageUrl string) {
		image, err := download(client, ctx, imageUrl)
		chMessage <- message{image: image, err: err} // return the message. 回传消息
	}(client, ctx, imageUrl)

	// 等待镜像拉取
	// wait for image pull.
	for {
		select {
		case <-ctx.Done():
			// 逾时停止下载镜像
			// stop downloading the image because the context is canceled.
			return nil, ctx.Err()
		case downloadMsg := <-chMessage:
			// 下载镜像失败
			// download the image failed.
			if downloadMsg.err != nil {
				goto RETRY
			}
			// 下载镜像成功
			// download the image successfully.
			return downloadMsg.image, downloadMsg.err
		default:
			// 等待一段时间
			// wait for a while.
			time.Sleep(1 * time.Second)
		}
	}
}

// Create 为容器创建
// Create is to create container.
func (d *Defaults) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	// 连接到网路环境
	// connect to the network environment.
	defaultNS := specs.LinuxNamespace{
		Type: specs.NetworkNamespace,
		Path: networkNS}

	// 创建一个新的预设容器
	// create a default container.
	return client.NewContainer(
		ctx,                               // 序文 context.
		containerName,                     // 容器名称 container name.
		containerd.WithImage(imagePulled), // 容器镜像 image.
		containerd.WithNewSnapshot(snapShot, imagePulled),                                           // 容器快照 snapshot.
		containerd.WithNewSpec(oci.WithImageConfig(imagePulled), oci.WithLinuxNamespace(defaultNS)), // 容器配置 spec.
	)
}

// Task 为容器任务创建
// Task is to create task.
func (d *Defaults) Task(container containerd.Container, ctx context.Context) (containerd.Task, error) {
	// 创建新的容器工作
	// create a task from the container.
	return container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
}

// Start 为容器任务启动
// Start is to start task.
func (d *Defaults) Start(task containerd.Task, ctx context.Context) error {
	// 开始执行容器工作
	// start the task.
	if err := task.Start(ctx); err != nil {
		return err
	}

	// 检查容器任务状态
	// check the status of the task.
LOOP:
	for {
		select {
		case <-ctx.Done():
			// 逾时停止容器工作
			// stop the task because the context is canceled.
			return ctx.Err()
		default:
			// 开始监听容器工作
			// start listening for the container work.
			status, err := task.Status(ctx)
			if err != nil {
				// 容器工作监听失败
				// container monitor failed.
				return err
			}
			if status.Status == containerd.Running {
				// 容器工作正在运行
				// container work is running.
				break LOOP
			}

			// 等待一秒
			// wait for one second.
			time.Sleep(1 * time.Second)
		}
	}

	// 容器工作执行成功
	// container work executed successfully.
	return nil
}

// >>>>> >>>>> >>>>> 检查部分 check part

// CheckService 为检查容器服务是否上线
// CheckService is to check container service.
func (d *Defaults) CheckService(ctx context.Context, ipAddrPort string) error {
	// 缺省容器没有服务，所以不需要检查.
	// Default container has no service, so no need to check.

	// 等待一秒
	// wait for one second.
	time.Sleep(1 * time.Second)

	// 正常返回.
	// return successfully.
	return nil
}

// CheckSchema 为检查容器资料是否存在 CheckService is to check container data exists.
func (d *Defaults) CheckSchema(ctx context.Context, ipAddrPort string) error {
	// 缺省容器没有服务，所以不需要检查.
	// Default container has no service, so no need to check.

	// 等待一秒
	// wait for one second.
	time.Sleep(1 * time.Second)

	// 正常返回.
	// return successfully.
	return nil
}

// >>>>> >>>>> >>>>> 删除部分 delete part

// Interrupt 为立刻停止容器任务
// Interrupt is to stop task immediately.
func (d *Defaults) Interrupt(task containerd.Task, ctx context.Context) error {
	// stop task immediately. 停止任务
LOOP:
	for {
		// 停止任务
		// stop the task.
		_ = task.Kill(ctx, syscall.SIGKILL)

		// 开始监听容器状态
		// start listening for the container status.
		status, err := task.Status(ctx)
		if err != nil {
			// container monitor failed. 容器工作监听失败
			return err
		}
		if status.Status != containerd.Running {
			// 容器工作不为正在运行
			// container is not running.
			break LOOP
		}
		// 等待一段时间
		// wait for a while.
		time.Sleep(1 * time.Second)
	}

	// 容器工作中止成功
	// stop task successfully.
	return nil
}

// Delete 为容器任务停止
// Delete is to delete task.
func (d *Defaults) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {

	// 删除容器工作
	// delete the task.
	if task != nil { // 确认容器工作存在 check task exist.
		if _, err := task.Delete(ctx); err != nil {
			// 删除容器工作失败
			// delete task failed.
			return err
		}
	}

	// 删除容器
	// delete the container.
	if container != nil { // 确认容器存在 check container exist.
		if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
			// 删除容器快照失败
			// delete container snapshot failed.
			return err
		}
	}

	// 删除成功
	// delete success.
	return nil
}

// 非约定的函数 non-defined interface.

// Set 为获取当前版本
// Set is to set current version.
func (d *Defaults) Set(version string) {
	d.debianVersion = version
}

// Version 获取当前版本
// Version is to get current version.
func (d *Defaults) Version() string {
	return d.debianVersion
}
