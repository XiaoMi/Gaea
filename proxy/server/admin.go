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
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"strings"
	"time"

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
	exit struct {
		C chan struct{}
	}
	proxy *Server
	model *models.ProxyInfo

	listener      net.Listener
	adminUser     string
	adminPassword string
	engine        *gin.Engine

	configType          string
	coordinatorAddr     string
	coordinatorUsername string
	coordinatorPassword string
	coordinatorRoot     string
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

	s.exit.C = make(chan struct{})
	s.proxy = proxy
	s.adminUser = cfg.AdminUser
	s.adminPassword = cfg.AdminPassword
	s.configType = cfg.ConfigType
	s.coordinatorAddr = cfg.CoordinatorAddr
	s.coordinatorUsername = cfg.UserName
	s.coordinatorPassword = cfg.Password
	s.coordinatorRoot = cfg.CoordinatorRoot

	s.engine = gin.New()
	l, err := net.Listen(cfg.ProtoType, cfg.AdminAddr)
	if err != nil {
		return nil, err
	}
	s.listener = l
	s.registerURL()
	s.registerMetric()
	s.registerProf()

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
	if err := s.unregisterProxy(); err != nil {
		log.Fatal("unregister proxy failed, %v", err)
		return err
	}
	return nil
}

func (s *AdminServer) registerURL() {
	adminGroup := s.engine.Group("/api/proxy", gin.BasicAuth(gin.Accounts{s.adminUser: s.adminPassword}))
	adminGroup.GET("/ping", s.ping)
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
	if s.configType == models.ConfigFile {
		return nil
	}
	client := models.NewClient(models.ConfigEtcd, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	store := models.NewStore(client)
	defer store.Close()
	if err := store.CreateProxy(s.model); err != nil {
		return err
	}
	return nil
}

func (s *AdminServer) unregisterProxy() error {
	if s.configType == models.ConfigFile {
		return nil
	}
	client := models.NewClient(models.ConfigEtcd, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	store := models.NewStore(client)
	defer store.Close()
	if err := store.DeleteProxy(s.model.Token); err != nil {
		return err
	}
	return nil
}

func (s *AdminServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func (s *AdminServer) prepareConfig(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		c.JSON(selfDefinedInternalError, "missing namespace name")
		return
	}
	client := models.NewClient(models.ConfigEtcd, s.coordinatorAddr, s.coordinatorUsername, s.coordinatorPassword, s.coordinatorRoot)
	defer client.Close()
	err := s.proxy.ReloadNamespacePrepare(name, client)
	if err != nil {
		log.Warn("prepare config of namespace: %s failed, err: %v", name, err)
		c.JSON(selfDefinedInternalError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}

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
	c.JSON(http.StatusOK, "OK")
}

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
	c.JSON(http.StatusOK, "OK")
}

func (s *AdminServer) configFingerprint(c *gin.Context) {
	c.JSON(http.StatusOK, s.proxy.manager.ConfigFingerprint())
}

// getNamespaceSessionSQLFingerprint return namespace sql fingerprint information
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
