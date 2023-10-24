package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/XiaoMi/Gaea/cc"
	"github.com/XiaoMi/Gaea/models"
)

const (
	// 注意这里的端口要和gaea_cc.ini 中的addr=0.0.0.0:23306
	gaeaCCBaseRouter    = "http://localhost:23306/api/cc/namespace/"
	gaeaCCAdminUsername = "admin"
	gaeaCCAdminPassword = "admin"
)

type NamespaceRegisterManager struct {
	Namespaces    map[string]*models.Namespace
	GaeaCCManager *GaeaCCManager
}

func NewNamespaceRegisterManger() *NamespaceRegisterManager {
	return &NamespaceRegisterManager{
		Namespaces: map[string]*models.Namespace{},
		GaeaCCManager: &GaeaCCManager{
			GaeaCCAdminUser:     gaeaCCAdminUsername,
			GaeaCCAdminPassword: gaeaCCAdminPassword,
			GaeaCCBaseRouter:    gaeaCCBaseRouter,
		},
	}
}

func (nr *NamespaceRegisterManager) AddNamespace(m *models.Namespace) error {
	// 存在就替换
	err := m.Verify()
	if err != nil {
		return err
	}
	if _, ok := nr.Namespaces[m.Name]; ok {
		return fmt.Errorf("namespace has exists")
	}
	nr.Namespaces[m.Name] = m
	return nil
}

func (nr *NamespaceRegisterManager) unRegisterNamespaceByKey(key string) error {
	// 在 ETCD 中删除 key
	err := nr.GaeaCCManager.deleteNamespace(key)
	delete(nr.Namespaces, key)
	if err != nil {
		return fmt.Errorf("error deleting from ETCD: %v", err)
	}
	return nil
}

// Register namespaces by gaeacc
func (nr *NamespaceRegisterManager) RegisterAllNamespaces() error {
	for _, namespace := range nr.Namespaces {
		err := nr.GaeaCCManager.modifyNamespace(namespace)
		if err != nil {
			return fmt.Errorf("error writing to ETCD: %v", err)
		}
	}
	return nil
}

// Unregister namespaces by gaeacc
func (nr *NamespaceRegisterManager) UnRegisterNamespaces() error {
	for _, namespace := range nr.Namespaces {
		err := nr.unRegisterNamespaceByKey(namespace.Name)
		if err != nil {
			return err
		}
	}
	return nil
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
	req, err := http.NewRequest("PUT", g.GaeaCCBaseRouter+"modify", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
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
	router := g.GaeaCCBaseRouter + "delete/" + key
	req, err := http.NewRequest("PUT", router, nil)
	if err != nil {
		return fmt.Errorf("create http request %s error: %w", router, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("create http request %s error: %w", router, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *GaeaCCManager) getNamespace(key string) ([]*models.Namespace, error) {
	router := g.GaeaCCBaseRouter + "detail/" + key
	req, err := http.NewRequest("GET", router, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request %s error: %w", router, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request to %s error: %w", router, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	var res cc.QueryNamespaceResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	return res.Data, nil
}

func (g *GaeaCCManager) listNamespace() ([]*models.Namespace, error) {
	res, err := http.Get(g.GaeaCCBaseRouter + "list")
	if err != nil {
		return nil, fmt.Errorf("GET request failed:%s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed:%s", err.Error())
	}

	var resp cc.QueryNamespaceResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response failed:%s", err.Error())
	}
	return resp.Data, nil
}
