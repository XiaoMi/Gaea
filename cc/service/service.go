// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"fmt"
	"sync"

	"github.com/XiaoMi/Gaea/cc/proxy"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
)

func getCoordinatorRoot(cluster string) string {
	if cluster != "" {
		return "/" + cluster
	}
	return cluster
}

// ListNamespace return names of all namespace
func ListNamespace(cfg *models.CCConfig, cluster string) ([]string, error) {
	client := models.NewClient(models.ConfigEtcd, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, getCoordinatorRoot(cluster))
	mConn := models.NewStore(client)
	defer mConn.Close()
	return mConn.ListNamespace()
}

// QueryNamespace return information of namespace specified by names
func QueryNamespace(names []string, cfg *models.CCConfig, cluster string) (data []*models.Namespace, err error) {
	client := models.NewClient(models.ConfigEtcd, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, getCoordinatorRoot(cluster))
	mConn := models.NewStore(client)
	defer mConn.Close()
	for _, v := range names {
		namespace, err := mConn.LoadNamespace(cfg.EncryptKey, v)
		if err != nil {
			log.Warn("load namespace %s failed, %v", v, err.Error())
			return nil, err
		}
		if namespace == nil {
			log.Warn("namespace %s not found", v)
			return data, nil
		}
		data = append(data, namespace)
	}

	return data, nil
}

// ModifyNamespace create or modify namespace
func ModifyNamespace(namespace *models.Namespace, cfg *models.CCConfig, cluster string) (err error) {
	if err = namespace.Verify(); err != nil {
		return fmt.Errorf("verify namespace error: %v", err)
	}

	// create/modify will save encrypted data default
	if err = namespace.Encrypt(cfg.EncryptKey); err != nil {
		return fmt.Errorf("encrypt namespace error: %v", err)
	}

	// sink namespace
	client := models.NewClient(models.ConfigEtcd, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, getCoordinatorRoot(cluster))
	storeConn := models.NewStore(client)
	defer storeConn.Close()

	if err := storeConn.UpdateNamespace(namespace); err != nil {
		log.Warn("update namespace failed, %s", string(namespace.Encode()))
		return err
	}

	// proxies ready to reload config
	proxies, err := storeConn.ListProxyMonitorMetrics()
	if err != nil {
		log.Warn("list proxies failed, %v", err)
		return err
	}

	// prepare phase
	for _, v := range proxies {
		err := proxy.PrepareConfig(v.IP+":"+v.AdminPort, namespace.Name, cfg)
		if err != nil {
			return err
		}
	}

	// commit phase
	for _, v := range proxies {
		err := proxy.CommitConfig(v.IP+":"+v.AdminPort, namespace.Name, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

// DelNamespace delete namespace
func DelNamespace(name string, cfg *models.CCConfig, cluster string) error {
	client := models.NewClient(models.ConfigEtcd, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, getCoordinatorRoot(cluster))
	mConn := models.NewStore(client)
	defer mConn.Close()

	if err := mConn.DelNamespace(name); err != nil {
		log.Warn("delete namespace %s failed, %s", name, err.Error())
		return err
	}

	proxies, err := mConn.ListProxyMonitorMetrics()
	if err != nil {
		log.Warn("list proxy failed, %s", err.Error())
		return err
	}

	for _, v := range proxies {
		err := proxy.DelNamespace(v.IP+":"+v.AdminPort, name, cfg)
		if err != nil {
			log.Warn("delete namespace %s in proxy %s failed, err: %s", name, v.IP, err.Error())
			return err
		}
	}

	return nil
}

// SQLFingerprint return sql fingerprints of all proxy
func SQLFingerprint(name string, cfg *models.CCConfig, cluster string) (slowSQLs, errSQLs map[string]string, err error) {
	slowSQLs = make(map[string]string, 16)
	errSQLs = make(map[string]string, 16)
	// list proxy
	client := models.NewClient(models.ConfigEtcd, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, getCoordinatorRoot(cluster))
	mConn := models.NewStore(client)
	defer mConn.Close()
	proxies, err := mConn.ListProxyMonitorMetrics()
	if err != nil {
		log.Warn("list proxy failed, %v", err)
		return nil, nil, err
	}
	wg := new(sync.WaitGroup)
	respC := make(chan *proxy.SQLFingerprint, len(proxies))
	// query sql fingerprints concurrently
	for _, p := range proxies {
		wg.Add(1)
		host := p.IP + ":" + p.AdminPort
		go func(host, name string) {
			defer wg.Done()
			r, err := proxy.QueryNamespaceSQLFingerprint(host, name, cfg)
			if err != nil {
				log.Warn("query namespace sql fingerprint failed ,%v", err)
			}
			respC <- r
		}(host, name)
	}
	wg.Wait()
	close(respC)

	for r := range respC {
		if r == nil {
			continue
		}
		for k, v := range r.SlowSQL {
			slowSQLs[k] = v
		}
		for k, v := range r.ErrorSQL {
			errSQLs[k] = v
		}
	}

	return
}

// ProxyConfigFingerprint return fingerprints of all proxy
func ProxyConfigFingerprint(cfg *models.CCConfig, cluster string) (r map[string]string, err error) {
	// list proxy
	client := models.NewClient(models.ConfigEtcd, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, getCoordinatorRoot(cluster))
	mConn := models.NewStore(client)
	defer mConn.Close()
	proxies, err := mConn.ListProxyMonitorMetrics()
	if err != nil {
		log.Warn("list proxy failed, %v", err)
		return nil, err
	}
	wg := new(sync.WaitGroup)
	r = make(map[string]string, len(proxies))
	respC := make(chan map[string]string, len(proxies))
	for _, p := range proxies {
		host := p.IP + ":" + p.AdminPort
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			md5, err := proxy.QueryProxyConfigFingerprint(host, cfg)
			if err != nil {
				log.Warn("query config fingerprint of proxy failed, %s %v", host, err)
			}
			m := make(map[string]string, 1)
			m[host] = md5
			respC <- m
		}(host)
	}
	wg.Wait()
	close(respC)
	for resp := range respC {
		if resp == nil {
			continue
		}
		for k, v := range resp {
			r[k] = v
		}
	}
	return
}
