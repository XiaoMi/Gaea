package containerTest

import (
	"encoding/json"
	"fmt"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd"
	"github.com/go-ini/ini"
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// >>>>> >>>>> >>>>> >>>>>> 载入 Containerd 初始配置
// >>>>> >>>>> >>>>> >>>>>> Load Containerd initial config

// ContainerIniConfig 为容器设定的初始配置
// ContainerIniConfig is the initial config for the container
type ContainerIniConfig struct {
	// config type
	ContainerTestEnable string `ini:"container_test_enable"`
}

// ParseContainerConfigFromFile 从档案获取容器初始配置
// ParseContainerConfigFromFile gets the container initial config from the file
func ParseContainerConfigFromFile(cfgFile string) (*ContainerIniConfig, error) {
	// 先决定配置文档的决对路径
	// determine the config file's absolute path
	absPath, err := absolutePath(cfgFile)
	if err != nil {
		return nil, err
	}

	// 读取配置文档的内容
	// read the config file's content
	cfg, err := ini.Load(absPath)
	if err != nil {
		return nil, err
	}

	// 创建一个容器 ini 设定配置对象
	// create a container ini config object
	var containerIniConfig = &ContainerIniConfig{}
	err = cfg.MapTo(containerIniConfig)
	if err != nil {
		return nil, err
	}

	// 回传结果
	// return the result
	return containerIniConfig, err
}

// >>>>> >>>>> >>>>> 获得所有的容器设定值
// >>>>> >>>>> >>>>> get all containerd configs

// Load 为用来执行容器服务
// Load is used to run containerd
type Load struct {
	// prefix 表示容器服务的设定目录的路径
	// prefix is the path of containerd configs
	prefix string
}

// ContainerdPath 返回容器服务设定目录的路径
// ContainerdPath returns the path of containerd configs.
func (r *Load) ContainerdPath() string {
	return filepath.Join(r.prefix, "containerd") // 在 prefix 下面的 containerd 目录 inside the containerd configs directory
}

// listContainerD 列出 ContainerConfig 的设置档列表, 有几个设置档就有几个名称
// listContainerD lists the containerd configs，one config file match one name.
func (r *Load) listContainerD() ([]string, error) {
	// 在 containerd 目录下面列出所有的设置档列表
	// list all configs in the containerd directory
	files, err := ioutil.ReadDir(r.ContainerdPath())
	if err != nil {
		return nil, err
	}

	// 收集所有的设置文件名称
	// collect all the configs names
	result := make([]string, 0)
	for _, f := range files {
		result = append(result, f.Name())
	}

	// 回传结果
	// return the result
	return result, nil
}

// loadContainerD 载入指定的设置档转成 ContainerConfig 设置值
// loadContainerD converts the specified config file to ContainerConfig config
func (r *Load) loadContainerD(file string) (containerd.ContainerConfig, error) {
	// 在 containerd 目录下面列出所有的设置档列表
	// list all configs in the containerd directory
	config := filepath.Join(r.ContainerdPath(), file)
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return containerd.ContainerConfig{}, err
	}

	// 开始把设置档转成 ContainerConfig 设置值
	// start converting the config to ContainerConfig config
	result := containerd.ContainerConfig{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return containerd.ContainerConfig{}, err
	}

	// 回传结果
	// return the result
	return result, nil
}

// loadAllContainerD 把所有的设置档转成 ContainerConfig 设置值
// loadAllContainerD converts all the config files to ContainerConfig config
func (r *Load) loadAllContainerD() (map[string]containerd.ContainerConfig, error) {
	// 在 containerd 目录下面列出所有的设置档列表
	// list all configs in the containerd directory
	files, err := r.listContainerD()
	if err != nil {
		return nil, err
	}

	// 创建一个 map 存放所有的设置档转成 ContainerConfig 设置值
	// create a map to store all the configs converted to ContainerConfig config
	result := make(map[string]containerd.ContainerConfig, len(files))
	for _, f := range files {
		config, err := r.loadContainerD(f)
		if err != nil {
			return nil, err
		}
		result[config.Name] = config
	}

	// 回传结果
	// return the result
	return result, nil
}

// 以下为群姐容器设定方式，这次 PR 不送出 !

// >>>>> >>>>> >>>>> 扩展容器服务的设定 extend the containerd config

// correctRange 是用来修正表示范围的字符串
// correctRange is used to correct the range string
func correctRange(name string) string {
	// 先修正字符串 fix the string
	// 使用 To 当成预期的格式 use To as the expected format
	name = strings.TrimSpace(name)
	name = strings.Replace(name, "TO", "To", -1)
	name = strings.Replace(name, "to", "To", -1)
	return name
}

// extendContainerName 是用来扩展容器名称的
// extendContainerName is used to extend the container name
func extendContainerName(name string) ([]string, error) {
	// 先修正字符串 fix the string
	name = correctRange(name)

	// 如没有 "To" "{" 或 "}" 如果有就返回错误 return error
	// if there is no "To" "{" or "}", return error
	if !strings.Contains(name, "To") ||
		!strings.Contains(name, "{") ||
		!strings.Contains(name, "}") {
		return nil, fmt.Errorf("the name %s is not valid", name)
	}

	// 先决定 { To 和 } 的位置 find the position of {, To, and }
	indexOpenCurly := strings.LastIndex(name, "{")
	indexCloseCurly := strings.LastIndex(name, "}")
	indexTo := strings.LastIndex(name, "To")

	// 决定前缀、中缀、后缀 prefix, middle, suffix
	prefix := name[:indexOpenCurly]
	middle := name[indexOpenCurly+1 : indexTo]
	suffix := name[indexTo+2 : indexCloseCurly]

	// 把 middleInt 和 middleFloat 转成 int 和 float convert middleInt and middleFloat to int and float
	middleInt, err := strconv.Atoi(middle)
	if err != nil {
		return nil, err
	}

	suffixInt, err := strconv.Atoi(suffix)
	if err != nil {
		return nil, err
	}

	// 收集所有的名称 collect all names
	var names = make([]string, suffixInt-middleInt+1)
	var j = 0

	// 开始收集 names。 collect names
	for i := middleInt; i <= suffixInt; i++ {
		names[j] = prefix + strconv.Itoa(i)
		j++
	}

	// 回传 names 回传 names。 return names
	return names, nil
}

// separateIPandPort 分离 IP 和 Port
// separateIPandPort is used to seperate IP and Port
func separateIPandPort(ip string) (string, string) {
	// 找到分离 IP 和 Port 的位置 find the position of seperating IP and Port.
	index := strings.LastIndex(ip, ":")

	// 如果 index 是 -1，那幺就没有 ":" 在 ip 中 if index is -1, it means that there is no ":" in the ip.
	if index == -1 {
		return ip, ""
	}

	// 回传 IP 和 Port 回传 IP 和 Port。 return IP and Port
	return ip[:index], ip[index+1:]
}

// extendContainerIP 是用来扩展容器网路位置的
// extendContainerIP is used to extend the container ip
func extendContainerIP(ipStr string) ([]string, error) {
	// 先修正字符串 fix the string
	// 分离 IP 和 Port。 separate IP and Port
	ipPortRange := strings.Split(correctRange(ipStr), "To")
	firstIP, firstPort := separateIPandPort(ipPortRange[0])
	endIP, endPort := separateIPandPort(ipPortRange[1])
	if firstPort != endPort {
		return nil, fmt.Errorf("the port %s is not match", endPort)
	}

	// ip 阵列。 ip array
	ips := make([]string, 0)

	// 下一个 IP 和 Port。 next IP and Port
	nextNetIP := net.ParseIP(firstIP)
	for {
		ips = append(ips, nextNetIP.String()+":"+firstPort)
		if nextNetIP.String() == endIP {
			// 如果 nextIP 字符串和 endIP 相同，那幺就退出退出循环 break the loop
			// if nextIP string is equal to endIP, then exit the loop
			break
		}
		// 下一个 IP 地址 nextIP
		nextNetIP = util.IncrementIP(nextNetIP)
	}

	// 回传 IP 阵列。 return IP array
	return ips, nil
}

func extendContainerConfig(config containerd.ContainerConfig) ([]containerd.ContainerConfig, error) {

	// names
	names, err := extendContainerName(config.Name)
	if err != nil {
		return nil, err
	}
	// ip
	ips, err := extendContainerIP(config.IP)

	// snapshot
	snapshots, err := extendContainerName(config.SnapShot)

	extendConfigs := make([]containerd.ContainerConfig, len(names))

	for i, name := range names {
		extendConfigs[i] = containerd.ContainerConfig{
			Sock:      config.Sock,
			Type:      config.Type,
			Name:      name,
			NameSpace: snapshots[i],
			Image:     config.Image,
			Task:      config.Task,
			NetworkNs: config.NetworkNs,
			IP:        ips[i],
			SnapShot:  config.SnapShot,
			Schema:    config.Schema,
			User:      config.User,
			Password:  config.Password,
		}
	}

	return extendConfigs, nil
}
