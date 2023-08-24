package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/XiaoMi/Gaea/cc"
	"github.com/XiaoMi/Gaea/models"
)

const (
	// 注意这里的端口要和gaea_cc.ini 中的addr=0.0.0.0:23306
	gaeaCCBaseRouter    = "http://localhost:23306/api/cc/namespace/"
	gaeaCCAdminUsername = "admin"
	gaeaCCAdminPassword = "admin"
)

var once sync.Once
var defaultNamespaceRegisterManger *NamespaceRegisterManger

// GetDefaultNamespaceRegisterManger returns the default NamespaceRegisterManger instance,
// and initializes it the first time this function is called.
func GetNamespaceRegisterManager() *NamespaceRegisterManger {
	once.Do(func() {
		defaultNamespaceRegisterManger = NewNamespaceRegisterManger()
	})
	return defaultNamespaceRegisterManger
}

type NamespaceRegisterManger struct {
	Namespaces    map[string]*models.Namespace
	GaeaCCManager *GaeaCCManager
}

func NewNamespaceRegisterManger() *NamespaceRegisterManger {
	return &NamespaceRegisterManger{
		Namespaces: map[string]*models.Namespace{},
		GaeaCCManager: &GaeaCCManager{
			GaeaCCAdminUser:     gaeaCCAdminUsername,
			GaeaCCAdminPassword: gaeaCCAdminPassword,
			GaeaCCBaseRouter:    gaeaCCBaseRouter,
		},
	}
}

func (nr *NamespaceRegisterManger) AddNamespace(m *models.Namespace) error {
	// 存在就替换
	err := m.Verify()
	if err != nil {
		return err
	}
	nr.Namespaces[m.Name] = m
	return nil
}

func (nr *NamespaceRegisterManger) UnRegisterNamespaceByKey(key string) error {

	delete(nr.Namespaces, key)
	// 在 ETCD 中删除 key
	err := nr.GaeaCCManager.deleteNamespace(key)
	if err != nil {
		return fmt.Errorf("error deleting from ETCD: %v", err)
	}
	return nil
}

// Register namespaces by gaeacc
func (nr *NamespaceRegisterManger) RegisterNamespaces() error {
	for _, namespace := range nr.Namespaces {
		err := nr.GaeaCCManager.modifyNamespace(namespace)
		if err != nil {
			return fmt.Errorf("error writing to ETCD: %v", err)
		}
	}
	return nil
}

// Unregister namespaces by gaeacc
func (nr *NamespaceRegisterManger) UnRegisterNamespaces() {
	for _, namespace := range nr.Namespaces {
		nr.UnRegisterNamespaceByKey(namespace.Name)
	}
}

func (nr *NamespaceRegisterManger) Get(key string) ([]*models.Namespace, error) {
	return nr.GaeaCCManager.getNamespace(key)
}

func (nr *NamespaceRegisterManger) List() ([]*models.Namespace, error) {
	return nr.GaeaCCManager.listNamespace()
}

func (nr *NamespaceRegisterManger) GetNsAllowedDBS(name string) map[string]bool {
	v, ok := nr.Namespaces[name]
	if ok {
		return v.AllowedDBS
	}
	return map[string]bool{}
}

func (nr *NamespaceRegisterManger) GetNsDefaultPhyDBS(name string) map[string]string {
	v, ok := nr.Namespaces[name]
	if ok {
		return v.DefaultPhyDBS
	}
	return map[string]string{}
}

func (nr *NamespaceRegisterManger) GetNsSlices(name string) []*models.Slice {

	v, ok := nr.Namespaces[name]
	if ok {
		return v.Slices
	}
	log.Fatal("key not exist:")
	return []*models.Slice{}
}

func (nr *NamespaceRegisterManger) GetNsShardRules(name string) []*models.Shard {
	v, ok := nr.Namespaces[name]
	if ok {
		return v.ShardRules
	}
	return []*models.Shard{}
}

func (nr *NamespaceRegisterManger) GetNsUsers(name string) []*models.User {
	v, ok := nr.Namespaces[name]
	if ok {
		return v.Users
	}
	return []*models.User{}
}

func (nr *NamespaceRegisterManger) GetNsDefaultSlice(name string) string {
	v, ok := nr.Namespaces[name]
	if ok {
		return v.DefaultSlice
	}
	return ""
}

func (nr *NamespaceRegisterManger) GetGaeaCCManager() *GaeaCCManager {
	return nr.GaeaCCManager
}

type GaeaCCManager struct {
	GaeaCCAdminUser     string
	GaeaCCAdminPassword string
	GaeaCCBaseRouter    string
}

func (g *GaeaCCManager) modifyNamespace(n *models.Namespace) error {
	// 将命名空间对象转换为 JSON
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequest("PUT", g.GaeaCCBaseRouter+"modify", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *GaeaCCManager) deleteNamespace(key string) error {
	req, err := http.NewRequest("PUT", g.GaeaCCBaseRouter+"delete/"+key, nil)
	if err != nil {
		log.Fatal("create http request ", g.GaeaCCBaseRouter+"delete/"+key, " error:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("send request to ", g.GaeaCCBaseRouter+"delete/"+key, " error:", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("resp StatusCode error ", err)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *GaeaCCManager) getNamespace(key string) ([]*models.Namespace, error) {
	//  req.SetBasicAuth("username", "password")
	req, err := http.NewRequest("GET", g.GaeaCCBaseRouter+"detail/"+key, nil)
	if err != nil {
		log.Fatal("create http request ", g.GaeaCCBaseRouter+"detail/"+key, " error:", err)
		return []*models.Namespace{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("send request to ", g.GaeaCCBaseRouter+"detail/"+key, " error:", err)
		return []*models.Namespace{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read response failed:", err)
		return nil, err
	}
	var res cc.QueryNamespaceResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal("unmarshal response failed::", err)
		return nil, err
	}
	return res.Data, nil
}

func (g *GaeaCCManager) listNamespace() ([]*models.Namespace, error) {
	res, err := http.Get(g.GaeaCCBaseRouter + "list")
	if err != nil {
		log.Fatal("GET request failed:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Read response failed:", err)
		return nil, err
	}

	var resp cc.QueryNamespaceResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal("Unmarshal response failed:", err)
		return nil, err
	}
	return resp.Data, nil
}
