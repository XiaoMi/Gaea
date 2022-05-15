package etcd

import "github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run/defaults"

// Etcd 为 etcd 容器配置项
// Etcd container configuration.
type Etcd struct {
	defaults.Defaults // 续承 Defaults. inherit Defaults.
}

// 约定的接口 defined interface

// >>>>> >>>>> >>>>> 创建部分 create part

// Pull 为容器拉取镜像
// Pull is to pull image from registry.
/*func (e *Etcd) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	return nil, nil
}*/

// Create 为容器创建
// Create is to create container.
/*func (e *Etcd) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	return nil, nil
}*/

// Task 为容器任务创建
// Task is to create task.
/*func (e *Etcd) Task(container containerd.Container, ctx context.Context) (containerd.Task, error) {
	return nil, nil
}*/

// Start 为容器任务启动
// Start is to start task.
/*func (e *Etcd) Start(task containerd.Task, ctx context.Context) error {
	return nil
}*/

// >>>>> >>>>> >>>>> 检查部分 check part

// CheckService 为检查容器服务是否上线
// CheckService is to check container service.
/*func (e *Etcd) CheckService(task containerd.Task, ctx context.Context) error {
	return nil
}*/

// CheckSchema 为检查容器资料是否存在
// CheckService is to check container data exists.
/*func (e *Etcd) CheckSchema(ctx context.Context, ipAddrPort string) error {
	return nil
}*/

// >>>>> >>>>> >>>>> 删除部分 delete part

// Interrupt 为立刻停止容器任务
// Interrupt is to stop task immediately.
/*func (e *Etcd) Interrupt(task containerd.Task, ctx context.Context) error {
	// kill the process work. 删除容器工作
	return task.Kill(ctx, syscall.SIGKILL)
}*/

// Delete 为容器任务停止
// Delete is to delete task.
/*func (e *Etcd) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}*/
