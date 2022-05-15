package builder

import (
	"time"
)

// Builder 是操作一个新的测试容器环境的接口
// Builder is an interface for making a new container test environment.
type Builder interface {
	// Build 创建一个新的容器环境 Create a new container test environment.
	Build(t time.Duration) error
	// OnService 检查容器服务是否上线 CheckService is to check container service.
	OnService(t time.Duration) error
	// TearDown 删除容器环境 Delete the container test environment.
	TearDown(t time.Duration) error
}
