package containerTest

import (
	"errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/log/xlog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// initContainerTestXLog 初始化容器测试日志
// initContainerTestXLog is init log of container test
func initContainerTestXLog() error {
	// 获取绝对路径
	// get absolute path
	absPath, err := absolutePath("util/mocks/containerTest/logs")
	if err != nil {
		return err
	}

	// 获得现在的时间字符串
	// get current time string
	timeStr := time.Now().Format("2006-01-02-15-04-05")

	// 容器测试日志设定
	// container test log setting
	cfg := make(map[string]string)
	cfg["path"] = absPath                       // 日志文档路径 log file path
	cfg["filename"] = "containerTest" + timeStr // 日志文档名 log file name
	cfg["level"] = "debug"                      // 日志级别 log level
	cfg["service"] = "containerTest"            // 日志服务名 log service name
	cfg["skip"] = "5"                           // 设置xlog打印方法堆栈需要跳过的层数, 5目前为调用log.Debug()等方法的方法名, 比xlog默认值多一层.

	// 创建容器测试日志对象
	// init container test log
	logger, err := xlog.CreateLogManager("file", cfg)
	if err != nil {
		return err
	}

	// 设置容器测试日志对象
	// set container test log
	log.SetGlobalLogger(logger)

	// 正常返回 return
	return nil
}

// absolutePath 获取绝对路径
// absolutePath get absolute path
func absolutePath(relativePath string) (string, error) {
	// 获得容器管理员服务配置
	// get the containerd manager config
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 决定容器管理员服务配置路径
	// decide the containerd manager config path
	absolutePath := ""
	if strings.Contains(path, "Gaea") {
		absolutePath = filepath.Join(strings.Split(path, "Gaea")[0], "Gaea", relativePath)
	} else {
		return "", errors.New("invalid config path")
	}

	// 成功返回路径
	// return path successfully
	return absolutePath, nil
}
