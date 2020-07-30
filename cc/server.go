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

package cc

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/XiaoMi/Gaea/cc/service"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
)

// Server admin server
type Server struct {
	cfg *models.CCConfig

	engine   *gin.Engine
	listener net.Listener

	exitC chan struct{}
}

// RetHeader response header
type RetHeader struct {
	RetCode    int    `json:"ret_code"`
	RetMessage string `json:"ret_message"`
}

// NewServer constructor of Server
func NewServer(addr string, cfg *models.CCConfig) (*Server, error) {
	srv := &Server{cfg: cfg, exitC: make(chan struct{})}
	srv.engine = gin.New()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	srv.listener = l
	srv.registerURL()
	return srv, nil
}

func (s *Server) registerURL() {
	api := s.engine.Group("/api/cc", gin.BasicAuth(gin.Accounts{s.cfg.AdminUserName: s.cfg.AdminPassword}))
	api.Use(gin.Recovery())
	api.Use(gzip.Gzip(gzip.DefaultCompression))
	api.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
	api.GET("/namespace/list", s.listNamespace)
	api.GET("/namespace", s.queryNamespace)
	api.GET("/namespace/detail/:name", s.detailNamespace)
	api.PUT("/namespace/modify", s.modifyNamespace)
	api.PUT("/namespace/delete/:name", s.delNamespace)
	api.GET("/namespace/sqlfingerprint/:name", s.sqlFingerprint)
	api.GET("/proxy/config/fingerprint", s.proxyConfigFingerprint)
}

// ListNamespaceResp list names of all namespace response
type ListNamespaceResp struct {
	RetHeader *RetHeader `json:"ret_header"`
	Data      []string   `json:"data"`
}

// return names of all namespace
func (s *Server) listNamespace(c *gin.Context) {
	var err error
	r := &ListNamespaceResp{RetHeader: &RetHeader{RetCode: -1, RetMessage: ""}}
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	r.Data, err = service.ListNamespace(s.cfg, cluster)
	if err != nil {
		log.Warn("list names of all namespace failed, %v", err)
		r.RetHeader.RetMessage = err.Error()
		c.JSON(http.StatusOK, r)
		return
	}
	r.RetHeader.RetCode = 0
	r.RetHeader.RetMessage = "SUCC"
	c.JSON(http.StatusOK, r)
	return
}

// QueryReq query namespace request
type QueryReq struct {
	Names []string `json:"names"`
}

// QueryNamespaceResp query namespace response
type QueryNamespaceResp struct {
	RetHeader *RetHeader          `json:"ret_header"`
	Data      []*models.Namespace `json:"data"`
}

func (s *Server) queryNamespace(c *gin.Context) {
	var err error
	var req QueryReq
	h := &RetHeader{RetCode: -1, RetMessage: ""}
	r := &QueryNamespaceResp{RetHeader: h}

	err = c.BindJSON(&req)
	if err != nil {
		log.Warn("queryNamespace got invalid data, err: %v", err)
		h.RetMessage = err.Error()
		c.JSON(http.StatusBadRequest, r)
		return
	}
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	r.Data, err = service.QueryNamespace(req.Names, s.cfg, cluster)
	if err != nil {
		log.Warn("query namespace failed, %v", err)
		c.JSON(http.StatusOK, r)
		return
	}

	h.RetCode = 0
	h.RetMessage = "SUCC"
	c.JSON(http.StatusOK, r)
	return
}

func (s *Server) detailNamespace(c *gin.Context) {
	var err error
	var names []string
	h := &RetHeader{RetCode: -1, RetMessage: ""}
	r := &QueryNamespaceResp{RetHeader: h}

	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		h.RetMessage = "input name is empty"
		c.JSON(http.StatusOK, h)
		return
	}

	names = append(names, name)
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	r.Data, err = service.QueryNamespace(names, s.cfg, cluster)
	if err != nil {
		log.Warn("query namespace failed, %v", err)
		c.JSON(http.StatusOK, r)
		return
	}

	h.RetCode = 0
	h.RetMessage = "SUCC"
	c.JSON(http.StatusOK, r)
	return
}

func (s *Server) modifyNamespace(c *gin.Context) {
	var err error
	var namespace models.Namespace
	h := &RetHeader{RetCode: -1, RetMessage: ""}

	err = c.BindJSON(&namespace)
	if err != nil {
		log.Warn("modifyNamespace failed, err: %v", err)
		c.JSON(http.StatusBadRequest, h)
		return
	}
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	err = service.ModifyNamespace(&namespace, s.cfg, cluster)
	if err != nil {
		log.Warn("modifyNamespace failed, err: %v", err)
		h.RetMessage = err.Error()
		c.JSON(http.StatusOK, h)
		return
	}

	h.RetCode = 0
	h.RetMessage = "SUCC"
	c.JSON(http.StatusOK, h)
	return
}

func (s *Server) delNamespace(c *gin.Context) {
	var err error
	h := &RetHeader{RetCode: -1, RetMessage: ""}
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		h.RetMessage = "input name is empty"
		c.JSON(http.StatusOK, h)
		return
	}
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	err = service.DelNamespace(name, s.cfg, cluster)
	if err != nil {
		h.RetMessage = fmt.Sprintf("delete namespace faild, %v", err.Error())
		c.JSON(http.StatusOK, h)
		return
	}

	h.RetCode = 0
	h.RetMessage = "SUCC"
	c.JSON(http.StatusOK, h)
	return
}

type sqlFingerprintResp struct {
	RetHeader *RetHeader        `json:"ret_header"`
	ErrSQLs   map[string]string `json:"err_sqls"`
	SlowSQLs  map[string]string `json:"slow_sqls"`
}

func (s *Server) sqlFingerprint(c *gin.Context) {
	var err error
	r := &sqlFingerprintResp{RetHeader: &RetHeader{RetCode: -1, RetMessage: ""}}
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		r.RetHeader.RetMessage = "input name is empty"
		c.JSON(http.StatusOK, r)
		return
	}
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	r.SlowSQLs, r.ErrSQLs, err = service.SQLFingerprint(name, s.cfg, cluster)
	if err != nil {
		r.RetHeader.RetMessage = err.Error()
		c.JSON(http.StatusOK, r)
		return
	}
	r.RetHeader.RetCode = 0
	r.RetHeader.RetMessage = "SUCC"
	c.JSON(http.StatusOK, r)
	return
}

type proxyConfigFingerprintResp struct {
	RetHeader *RetHeader        `json:"ret_header"`
	Data      map[string]string `json:"data"` // key: ip:port value: md5 of config
}

func (s *Server) proxyConfigFingerprint(c *gin.Context) {
	var err error
	r := &proxyConfigFingerprintResp{RetHeader: &RetHeader{RetCode: -1, RetMessage: ""}}
	cluster := c.DefaultQuery("cluster", s.cfg.DefaultCluster)
	r.Data, err = service.ProxyConfigFingerprint(s.cfg, cluster)
	if err != nil {
		r.RetHeader.RetMessage = err.Error()
		c.JSON(http.StatusOK, r)
		return
	}
	r.RetHeader.RetCode = 0
	r.RetHeader.RetMessage = "SUCC"
	c.JSON(http.StatusOK, r)
	return
}

func (s *Server) Run() {
	defer s.listener.Close()

	errC := make(chan error)

	go func(l net.Listener) {
		h := http.NewServeMux()
		h.Handle("/", s.engine)
		hs := &http.Server{Handler: h}
		errC <- hs.Serve(l)
	}(s.listener)

	select {
	case <-s.exitC:
		log.Notice("server exit.")
		return
	case err := <-errC:
		log.Fatal("gaea cc serve failed, %v", err)
		return
	}

}

func (s *Server) Close() {
	s.exitC <- struct{}{}
	return
}
