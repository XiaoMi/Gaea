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

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/core"

	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/util"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const (
	selfDefinedInternalError = 800
)

// SQLFingerprint sql fingerprint
type SQLFingerprint struct {
	SlowSQL  map[string]string `json:"slow_sql"`
	ErrorSQL map[string]string `json:"error_sql"`
}

// AdminServer means admin server
type AdminServer struct {
	ctx    context.Context
	cancel context.CancelFunc

	exit struct {
		C chan struct{}
	}
	proxy *Server
	model *models.ProxyInfo
	// proxy config file path
	configFile string

	listener      net.Listener
	adminUser     string
	adminPassword string
	engine        *gin.Engine

	configType          string
	coordinatorAddr     string
	coordinatorUsername string
	coordinatorPassword string
	coordinatorRoot     string

	localNamespaceStoragePath string
}

// NewAdminServer create new admin server
func NewAdminServer(proxy *Server, cfg *models.Proxy) (*AdminServer, error) {
	var err error
	s := new(AdminServer)

	// if error occurs, recycle the resources during creation.
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("NewAdminServer panic: %v", e)
		}

		if err != nil {
			s.Close()
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel
	s.exit.C = make(chan struct{})
	s.proxy = proxy
	s.adminUser = cfg.AdminUser
	s.adminPassword = cfg.AdminPassword
	s.configType = cfg.ConfigType
	s.coordinatorAddr = cfg.CoordinatorAddr
	s.coordinatorUsername = cfg.UserName
	s.coordinatorPassword = cfg.Password
	s.coordinatorRoot = cfg.CoordinatorRoot
	s.localNamespaceStoragePath = cfg.LocalNamespaceStoragePath
	s.configFile = cfg.ConfigFile
	s.engine = gin.New()
	l, err := net.Listen(cfg.ProtoType, cfg.AdminAddr)
	if err != nil {
		return nil, err
	}
	s.listener = l
	s.registerURL()
	s.registerMetric()
	s.registerProf()
	s.registerVersion()

	proxyInfo, err := NewProxyInfo(cfg, s.proxy.Listener().Addr().String())
	if err != nil {
		return nil, err
	}
	s.model = proxyInfo

	if err = s.registerProxy(); err != nil {
		return nil, err
	}

	log.Notice("[server] NewAdminServer, Api Server running, netProto: http, addr: %s", cfg.AdminAddr)
	return s, nil
}

// Run run admin server
func (s *AdminServer) Run() {
	defer s.listener.Close()

	eh := make(chan error, 1)
	go func(l net.Listener) {
		h := http.NewServeMux()
		h.Handle("/", s.engine)
		hs := &http.Server{Handler: h}
		eh <- hs.Serve(s.listener)
	}(s.listener)

	select {
	case <-s.exit.C:
		log.Warn("[%p] admin shutdown", s)
	case err := <-eh:
		log.Fatal("[%p] admin exit on error:%v", s, err)
	}
}

// Close close admin server
func (s *AdminServer) Close() error {
	close(s.exit.C)
	s.cancel()
	if err := s.unregisterProxy(); err != nil {
		log.Fatal("unregister proxy failed, %v", err)
		return err
	}
	return nil
}

func (s *AdminServer) registerURL() {
	adminGroup := s.engine.Group("/api/proxy", gin.BasicAuth(gin.Accounts{s.adminUser: s.adminPassword}))
	adminGroup.GET("/ping", s.ping)
	adminGroup.PUT("/proxyconfig/reload", s.reloadProxyConfig)
	adminGroup.PUT("/config/prepare/:name", s.prepareConfig)
	adminGroup.PUT("/config/commit/:name", s.commitConfig)
	adminGroup.PUT("/namespace/delete/:name", s.deleteNamespace)
	adminGroup.GET("/config/fingerprint", s.configFingerprint)

	adminGroup.GET("/stats/sessionsqlfingerprint/:namespace", s.getNamespaceSessionSQLFingerprint)
	adminGroup.GET("/stats/backendsqlfingerprint/:namespace", s.getNamespaceBackendSQLFingerprint)
	adminGroup.DELETE("/stats/sessionsqlfingerprint/:namespace", s.clearNamespaceSessionSQLFingerprint)
	adminGroup.DELETE("/stats/backendsqlfingerprint/:namespace", s.clearNamespaceBackendSQLFingerprint)

	adminGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	adminGroup.Use(gin.Recovery())
	adminGroup.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
}

// @Summary 获取proxy prometheus指标信息
// @Description 获取gaea proxy prometheus指标信息
// @Security BasicAuth
// @Router /api/metric/metrics [get]
func (s *AdminServer) registerMetric() {
	metricGroup := s.engine.Group("/api/metric", gin.BasicAuth(gin.Accounts{s.adminUser: s.adminPassword}))
	for path, handler := range s.proxy.manager.GetStatisticManager().GetHandlers() {
		log.Debug("[server] AdminServer got metric handler, path: %s", path)
		metricGroup.GET(path, gin.WrapH(handler))
	}
}

func (s *AdminServer) registerProf() {
	profGroup := s.engine.Group("/debug/pprof", gin.BasicAuth(gin.Accounts{s.adminUser: s.adminPassword}))
	profGroup.GET("/", gin.WrapF(pprof.Index))
	profGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	profGroup.GET("/profile", gin.WrapF(pprof.Profile))
	profGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
	profGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
	profGroup.GET("/trace", gin.WrapF(pprof.Trace))
	profGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	profGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	profGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	profGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	profGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	profGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
}

func (s *AdminServer) registerVersion() {
	versionGroup := s.engine.Group("/api/proxy")
	versionGroup.GET("/version", s.ProxyVersion)
}

// NewProxyInfo create proxy information
func NewProxyInfo(cfg *models.Proxy, addr string) (*models.ProxyInfo, error) {
	ipPort, err := util.ResolveAddr(cfg.ProtoType, addr)
	if err != nil {
		return nil, err
	}

	proxyIPPort := strings.Split(cfg.ProxyAddr, ":")
	adminIPPort := strings.Split(cfg.AdminAddr, ":")

	proxyInfo := &models.ProxyInfo{
		StartTime: time.Now().String(),
		ProtoType: cfg.ProtoType,
		ProxyPort: proxyIPPort[1],
		AdminPort: adminIPPort[1],
	}
	tmp := strings.Split(ipPort, ":")
	proxyInfo.IP = tmp[0]
	proxyInfo.Pid = os.Getpid()
	proxyInfo.Pwd, _ = os.Getwd()
	o, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return nil, err
	}
	proxyInfo.Sys = strings.TrimSpace(string(o))

	x, err := generateToken(cfg.ProtoType, addr)
	if err != nil {
		return nil, err
	}
	proxyInfo.Token = x

	return proxyInfo, nil
}

func generateToken(protoType, addr string) (string, error) {
	ipPort, err := util.ResolveAddr(protoType, addr)
	if err != nil {
		return "", err
	}
	return ipPort, nil
}

func (s *AdminServer) registerProxy() error {
	// 如果设定值 s.configType 是为文档 File 就直接回传
	if s.configType == models.ConfigFile {
		return nil
	}
	// 目前设定值 s.configType 可能为 models.ConfigEtcd 和 models.ConfigEtcdV3 两种
	remote, err := models.NewClient(s.configType, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	// register to the remote if the remote is available
	if err == nil {
		store := models.NewStore(remote)
		defer store.Close()
		return store.CreateProxy(s.model)
	}
	log.Warn("[admin server] register proxy failed: remote client is unavailable", err)

	// If the local storage is enabled, start a background task to retry registration
	if len(s.localNamespaceStoragePath) > 0 {
		log.Warn("[admin server] Remote is unavailable, switching to local storage")
		// Start a background task to keep trying to register and alert until successful
		go s.startRegistrationRetry()
		return nil
	}
	// No client available and no local storage, return an error
	return fmt.Errorf("[admin server] register proxy failed: no client available error:'%v'", err)
}

func (s *AdminServer) unregisterProxy() error {
	if s.configType == models.ConfigFile {
		return nil
	}
	remote, err := models.NewClient(s.configType, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	// unregister to the remote if the remote is available
	if err == nil {
		store := models.NewStore(remote)
		defer store.Close()
		return store.DeleteProxy(s.model.Token)
	}
	log.Warn("[admin server] unregister proxy failed: remote client is unavailable", err)
	// if the local is enabled, no error is returned here so that the data can be read from the local normally in the future.
	if len(s.localNamespaceStoragePath) > 0 {
		log.Warn("[admin server] unregister proxy failed: remote is unavailable, switch to local", err)
		return nil
	}
	return fmt.Errorf("[admin server] unregister proxy failed: no client available error:'%v'", err)
}

func (s *AdminServer) startRegistrationRetry() {
	retryInterval := 5 * time.Second
	for {
		select {
		case <-s.ctx.Done():
			// The AdminServer is shutting down
			log.Warn("[admin server] Shutting down registration retry routine")
			return
		case <-time.After(retryInterval):
			log.Warn("[admin server] Attempting to re-register proxy to remote ETCD")
			err := s.tryRegisterProxy()
			if err == nil {
				log.Warn("[admin server] Successfully re-registered proxy to remote ETCD")
				// Registration succeeded, exit the retry loop
				return
			} else {
				// Optionally, alert here
				log.Warn("[admin server] Re-register proxy failed: %v", err)
			}
		}
	}
}

func (s *AdminServer) tryRegisterProxy() error {
	// Attempt to create a new client for the remote ETCD service
	remote, err := models.NewClient(s.configType, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	if err != nil {
		return fmt.Errorf("failed to create remote client: %v", err)
	}
	defer remote.Close()

	store := models.NewStore(remote)
	defer store.Close()

	// Attempt to create the proxy in the remote store
	err = store.CreateProxy(s.model)
	if err != nil {
		return fmt.Errorf("failed to create proxy in remote store: %v", err)
	}

	return nil
}

// @Summary 获取proxy admin接口状态
// @Description  获取proxy admin接口状态
// @Success 200 {string} string "OK"
// @Security BasicAuth
// @Router /api/proxy/ping [get]
func (s *AdminServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

// @Summary reload proxy config
// @Description 通过管理接口, 重载 proxy 配置文件，当前仅支持 log 配置的重载
// @Produce  json
// @Success 200 {string} string "OK"
// @Security BasicAuth
// @Router /api/proxy/proxyconfig/reload [put]
func (s *AdminServer) reloadProxyConfig(c *gin.Context) {
	if err := s.proxy.ReloadProxyConfig(); err != nil {
		c.JSON(selfDefinedInternalError, fmt.Sprintf("reload config file Error:%s", err))
	}
	log.Notice("reload proxy config success")
	c.JSON(http.StatusOK, "OK")
}

// @Summary prepare namespace配置
// @Description 通过管理接口, 二阶段提交, prepare namespace配置
// @Produce  json
// @Param name path string true "namespace name"
// @Success 200 {string} string "OK"
// @Security BasicAuth
// @Router /api/proxy/config/prepare/{name} [put]
func (s *AdminServer) prepareConfig(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(selfDefinedInternalError, "missing namespace name")
		return
	}
	client, err := models.NewClient(s.configType, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	if client == nil || err != nil {
		c.JSON(selfDefinedInternalError, fmt.Errorf("[admin server] prepare config failed: remote client unavailable error:'%v'", err))
		return
	}
	defer client.Close()

	err = s.proxy.ReloadNamespacePrepare(name, client)
	if err != nil {
		log.Warn("[admin server] prepare config of namespace: %s failed, err: %v", name, err)
		c.JSON(selfDefinedInternalError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}

// @Summary commit namespace配置
// @Description 通过管理接口, 二阶段提交, commit namespace配置, 使etcd配置生效
// @Produce  json
// @Param name path string true "namespace name"
// @Success 200 {string} string "OK"
// @Security BasicAuth
// @Router /api/proxy/config/commit/{name} [put]
func (s *AdminServer) commitConfig(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(selfDefinedInternalError, "missing namespace name")
		return
	}
	err := s.proxy.ReloadNamespaceCommit(name)
	if err != nil {
		c.JSON(selfDefinedInternalError, err.Error())
		return
	}
	go s.updateNamespaceLocal(name)
	c.JSON(http.StatusOK, "OK")
}

// updateNamespaceLocal synchronizes a local namespace with its remote version.
// This method is called to ensure the local data is up-to-date with the remote server.
// Note: This method does not return any value and handles all errors internally by logging them.
// It is critical that the `localNamespaceStoragePath` and remote server configurations are properly set before calling this method.
func (s *AdminServer) updateNamespaceLocal(name string) {
	if len(s.localNamespaceStoragePath) == 0 {
		return
	}
	// Get updated Namespace from remote
	remote, err := models.NewClient(s.configType, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	if err != nil {
		log.Warn("[admin server] update namespace local failed: remote client unavailable error:'%v'", err)
		return
	}
	remoteStore := models.NewStore(remote)
	defer remoteStore.Close()
	ns, err := remoteStore.LoadOriginNamespace(name)
	if err != nil {
		log.Warn("[admin server] update namespace local failed: remote store:'%s' error:'%v'", name, err)
		return
	}
	// Synchronize remote and local namespaces
	local, err := models.NewLocalClient(s.localNamespaceStoragePath, s.coordinatorRoot)
	if err != nil {
		log.Warn("[admin server] update namespace local failed: local client unavailable error:'%v'", err)
		return
	}
	localStore := models.NewStore(local)
	defer localStore.Close()
	err = localStore.UpdateNamespace(ns)
	if err != nil {
		log.Warn("[admin server] update namespace local '%s' failed, error:'%v'", name, err)
		return
	}
}

// @Summary 删除namespace配置
// @Description 通过管理接口删除指定namespace配置
// @Produce  json
// @Param name path string true "namespace name"
// @Success 200 {string} string "OK"
// @Security BasicAuth
// @Router /api/proxy/config/delete/{name} [put]
func (s *AdminServer) deleteNamespace(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(selfDefinedInternalError, "missing namespace name")
		return
	}
	// delete namespace
	err := s.proxy.DeleteNamespace(name)
	if err != nil {
		c.JSON(selfDefinedInternalError, err.Error())
		return
	}
	go s.deleteNamespaceLocal(name)
	c.JSON(http.StatusOK, "OK")
}

// deleteNamespaceLocal removes a namespace from the local storage.
// This method is primarily used to ensure that local copies of namespaces are deleted when no longer needed or when updates require removal.
// Note: This method does not return any value and handles all errors internally by logging them.
// It is essential that the local namespace storage path is properly configured before this method is called.
func (s *AdminServer) deleteNamespaceLocal(name string) {
	if len(s.localNamespaceStoragePath) == 0 {
		return
	}
	local, err := models.NewLocalClient(s.localNamespaceStoragePath, s.coordinatorRoot)
	if err != nil {
		log.Warn("[admin server] new local client error", err)
		return
	}

	localStore := models.NewStore(local)
	defer localStore.Close()
	err = localStore.DelNamespace(name)
	if err != nil {
		log.Warn("[admin server] delete namespace local '%s' failed, error:'%v'", name, err)
		return
	}
}

// @Summary 返回配置指纹
// @Description 返回配置指纹, 指纹随配置变化而变化
// @Produce  json
// @Success 200 {string} string "Config Fingerprint"
// @Security BasicAuth
// @Router /api/proxy/config/fingerprint [get]
func (s *AdminServer) configFingerprint(c *gin.Context) {
	c.JSON(http.StatusOK, s.proxy.manager.ConfigFingerprint())
}

// @Summary 获取Porxy 慢SQL、错误SQL信息
// @Description 通过管理接口获取Porxy 慢SQL、错误SQL信息
// @Produce  json
// @Param namespace path string true "namespace name"
// @Success 200 {object} SQLFingerprint
// @Security BasicAuth
// @Router /api/proxy/stats/sessionsqlfingerprint/{namespace} [get]
func (s *AdminServer) getNamespaceSessionSQLFingerprint(c *gin.Context) {
	ns := strings.TrimSpace(c.Param("namespace"))
	namespace := s.proxy.manager.GetNamespace(ns)
	if namespace == nil {
		c.JSON(selfDefinedInternalError, "namespace not found")
		return
	}

	slowSQLFingerprints := namespace.GetSlowSQLFingerprints()
	errSQLFingerprints := namespace.GetErrorSQLFingerprints()
	ret := &SQLFingerprint{SlowSQL: slowSQLFingerprints, ErrorSQL: errSQLFingerprints}

	c.JSON(http.StatusOK, ret)
}

// @Summary 获取后端节点慢SQL、错误SQL信息
// @Description 通过管理接口获取后端节点慢SQL、错误SQL信息
// @Produce  json
// @Param namespace path string true "namespace name"
// @Success 200 {object} SQLFingerprint
// @Security BasicAuth
// @Router /api/proxy/stats/backendsqlfingerprint/{namespace} [get]
func (s *AdminServer) getNamespaceBackendSQLFingerprint(c *gin.Context) {
	ns := strings.TrimSpace(c.Param("namespace"))
	namespace := s.proxy.manager.GetNamespace(ns)
	if namespace == nil {
		c.JSON(selfDefinedInternalError, "namespace not found")
		return
	}

	slowSQLFingerprints := namespace.GetBackendSlowSQLFingerprints()
	errSQLFingerprints := namespace.GetBackendErrorSQLFingerprints()
	ret := &SQLFingerprint{SlowSQL: slowSQLFingerprints, ErrorSQL: errSQLFingerprints}

	c.JSON(http.StatusOK, ret)
}

// @Summary 清空Porxy节点慢SQL、错误SQL信息
// @Description 通过管理接口清空Porxy慢SQL、错误SQL信息
// @Produce  json
// @Param namespace path string true "namespace name"
// @Success 200 {object} SQLFingerprint
// @Security BasicAuth
// @Router /api/proxy/stats/sessionsqlfingerprint/{namespace} [delete]
func (s *AdminServer) clearNamespaceSessionSQLFingerprint(c *gin.Context) {
	ns := strings.TrimSpace(c.Param("namespace"))
	namespace := s.proxy.manager.GetNamespace(ns)
	if namespace == nil {
		c.JSON(selfDefinedInternalError, "namespace not found")
		return
	}

	namespace.ClearSlowSQLFingerprints()
	namespace.ClearErrorSQLFingerprints()

	c.JSON(http.StatusOK, "OK")
}

// @Summary 清空后端节点慢SQL、错误SQL信息
// @Description 通过管理接口清空后端节点慢SQL、错误SQL信息
// @Produce  json
// @Param namespace path string true "namespace name"
// @Success 200 {object} SQLFingerprint
// @Security BasicAuth
// @Router /api/proxy/stats/backendsqlfingerprint/{namespace} [delete]
func (s *AdminServer) clearNamespaceBackendSQLFingerprint(c *gin.Context) {
	ns := strings.TrimSpace(c.Param("namespace"))
	namespace := s.proxy.manager.GetNamespace(ns)
	if namespace == nil {
		c.JSON(selfDefinedInternalError, "namespace not found")
		return
	}

	namespace.ClearBackendSlowSQLFingerprints()
	namespace.ClearBackendErrorSQLFingerprints()

	c.JSON(http.StatusOK, "OK")
}

// @Summary 获取gaea版本信息
// @Description  获取gaea版本信息，2.0版本新增接口
// @Success 200 {string} string "version"
// @Security 不需要鉴权
// @Router /api/proxy/config/version [get]
func (s *AdminServer) ProxyVersion(c *gin.Context) {
	c.JSON(http.StatusOK, core.Info.Version)
}
