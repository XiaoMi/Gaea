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
	"net"
	"runtime"
	"strconv"
	"time"

	"fmt"

	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/sync2"
)

var (
	timeWheelUnit       = time.Second * 1
	timeWheelBucketsNum = 3600
)

// Server means proxy that serve client request
type Server struct {
	closed         sync2.AtomicBool
	listener       net.Listener
	sessionTimeout time.Duration
	tw             *util.TimeWheel
	adminServer    *AdminServer
	manager        *Manager
	EncryptKey     string
	ServerVersion  string
	AuthPlugin     string
}

// NewServer create new server
func NewServer(cfg *models.Proxy, manager *Manager) (*Server, error) {
	var err error
	s := new(Server)

	// init key
	s.EncryptKey = cfg.EncryptKey
	s.manager = manager
	s.ServerVersion = cfg.ServerVersion
	s.AuthPlugin = cfg.AuthPlugin
	if len(s.AuthPlugin) > 0 {
		DefaultCapability |= mysql.ClientPluginAuth
	}

	// if error occurs, recycle the resources during creation.
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("NewServer panic: %v", e)
		}

		if err != nil {
			s.Close()
		}
	}()

	s.closed = sync2.NewAtomicBool(false)

	s.listener, err = net.Listen(cfg.ProtoType, cfg.ProxyAddr)
	if err != nil {
		return nil, err
	}

	st := strconv.Itoa(cfg.SessionTimeout)
	st = st + "s"
	s.sessionTimeout, err = time.ParseDuration(st)
	if err != nil {
		return nil, err
	}

	s.tw, err = util.NewTimeWheel(timeWheelUnit, timeWheelBucketsNum)
	if err != nil {
		return nil, err
	}
	s.tw.Start()

	// create AdminServer
	adminServer, err := NewAdminServer(s, cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("NewAdminServer error, quit. error: %s", err.Error()))
		return nil, err
	}
	s.adminServer = adminServer

	log.Notice("server start succ, netProtoType: %s, addr: %s", cfg.ProtoType, cfg.ProxyAddr)
	return s, nil
}

// Listener return proxy's listener
func (s *Server) Listener() net.Listener {
	return s.listener
}

func (s *Server) onConn(c net.Conn) {
	cc := newSession(s, c) //新建一个conn
	defer func() {
		err := recover()
		if err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)] //获得当前goroutine的stacktrace
			log.Warn("[server] onConn panic error, remoteAddr: %s, stack: %s", c.RemoteAddr().String(), string(buf))
		}

		// close session finally
		cc.Close()
	}()

	if err := cc.Handshake(); err != nil {
		log.Warn("[server] onConn error: %s", err.Error())
		if err != mysql.ErrBadConn {
			cc.c.writeErrorPacket(err)
		}
		return
	}

	// must invoke after handshake
	if allowConnect := cc.IsAllowConnect(); allowConnect == false {
		err := mysql.NewError(mysql.ErrAccessDenied, "ip address access denied by gaea")
		cc.c.writeErrorPacket(err)
		return
	}

	// added into time wheel
	s.tw.Add(s.sessionTimeout, cc, cc.Close)
	log.Notice("Connected conn_id=%d, %s@%s (%s) namespace:%s capability: %d",
		cc.c.ConnectionID,
		cc.executor.user,
		cc.executor.clientAddr,
		cc.executor.db,
		cc.executor.namespace,
		cc.c.capability)
	cc.Run()
}

// Run proxy run and serve client request
func (s *Server) Run() error {
	// start AdminServer first
	go s.adminServer.Run()

	// start Server
	s.closed.Set(false)
	for s.closed.Get() != true {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Warn("[server] listener accept error: %s", err.Error())
			continue
		}

		go s.onConn(conn)
	}

	return nil
}

// Close close proxy server
func (s *Server) Close() error {
	if s.adminServer != nil {
		s.adminServer.Close()
	}

	s.closed.Set(true)
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}

	s.manager.Close()
	return nil
}

// ReloadNamespacePrepare config change prepare phase
func (s *Server) ReloadNamespacePrepare(name string, client models.Client) error {
	// get namespace conf from etcd
	log.Notice("prepare config of namespace: %s begin", name)
	store := models.NewStore(client)
	namespaceConfig, err := store.LoadNamespace(s.EncryptKey, name)
	if err != nil {
		return err
	}

	if err = s.manager.ReloadNamespacePrepare(namespaceConfig); err != nil {
		log.Warn("Manager ReloadNamespacePrepare error: %v", err)
		return err
	}

	log.Notice("prepare config of namespace: %s end", name)
	return nil
}

// ReloadNamespaceCommit config change commit phase
// commit namespace does not need lock
func (s *Server) ReloadNamespaceCommit(name string) error {
	log.Notice("commit config of namespace: %s begin", name)

	if err := s.manager.ReloadNamespaceCommit(name); err != nil {
		log.Warn("Manager ReloadNamespaceCommit error: %v", err)
		return err
	}

	log.Notice("commit config of namespace: %s end", name)
	return nil
}

// DeleteNamespace delete namespace in namespace manager
func (s *Server) DeleteNamespace(name string) error {
	log.Notice("delete namespace begin: %s", name)

	if err := s.manager.DeleteNamespace(name); err != nil {
		log.Warn("Manager DeleteNamespace error: %v", err)
		return err
	}

	log.Notice("delete namespace end: %s", name)
	return nil
}
