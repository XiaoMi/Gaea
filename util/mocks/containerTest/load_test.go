package containerTest

import (
	"fmt"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// 测试容器设定文档的位罝 Test the container config file position.
func TestContainerdPath(t *testing.T) {
	r := Load{prefix: "./test/"}
	path := r.ContainerdPath()
	require.Equal(t, path, "test/containerd")
}

// 测试容器设定文档的数量和名称 Test the container config files' number and name.
func TestListContainerD(t *testing.T) {
	r := Load{prefix: "./example/"}
	files, err := r.listContainerD()
	require.Nil(t, err)
	require.Equal(t, len(files), 2)
	require.Contains(t, files, "defaults.json")
	require.Contains(t, files, "mariadb.json")
}

// 测试容器设定文档的载入 Test the container config file loading.
func TestLoadContainerD(t *testing.T) {
	r := Load{prefix: "./example/"}
	mariadb, err := r.loadContainerD("mariadb.json")
	require.Nil(t, err)
	require.Equal(t, mariadb.Name, "mariadb-server")
}

// 测试所有容器设定文档的载入 Test all container config files loading.
func TestLoadAllContainerD(t *testing.T) {
	r := Load{prefix: "./example/"}
	configs, err := r.loadAllContainerD()
	require.Nil(t, err)
	for key, value := range configs {
		require.Equal(t, key, value.Name)
	}
}

// 扩展容器设定文档的载入 extend the container config file loading
func TestLoadExtendContainerD(t *testing.T) {
	// 先载入所有的容器设定文档 Load all container config files
	r := Load{prefix: "./example/"}
	configs, err := r.loadAllContainerD()
	require.Nil(t, err)

	//
	tmp := make([]containerd.ContainerConfig, 0)
	for key, value := range configs {
		if strings.Contains(key, "{") {
			extendConfig, err := extendContainerConfig(value)
			require.Nil(t, err)
			tmp = append(tmp, extendConfig...)
			delete(configs, key)
		}
	}

	for _, value := range tmp {
		configs[value.Name] = value
	}

	fmt.Println(configs)
}

// >>>>> >>>>> >>>>> 测试扩展容器服务的设定 test extending the containerd config

// TestExtendContainerName 测试扩展容器名称 test extending the container name
func TestExtendContainerName(t *testing.T) {
	// 容器名称由 mariadb-server-12 至 mariadb-server-100
	names, err := extendContainerName("mariadb-server-{12To100}")
	require.Nil(t, err)
	// 所扩展的容器名称数量要为 89
	// extend number of container names must be 89
	require.Equal(t, len(names), 89)
	// 检查最后一个容器名称是否正确
	// check the last container name is correct
	require.Equal(t, names[len(names)-1], "mariadb-server-100")
}

// TestSeparateIPandPort 测试分离IP和端口 test separating the ip and port
func TestSeparateIPandPort(t *testing.T) {
	ipPort := "192.168.122.1:3306"
	ip, port := separateIPandPort(ipPort)
	require.Equal(t, ip, "192.168.122.1")
	require.Equal(t, port, "3306")
}
