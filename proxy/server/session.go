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
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
)

// DefaultCapability means default capability
var DefaultCapability = mysql.ClientLongPassword | mysql.ClientLongFlag |
	mysql.ClientConnectWithDB | mysql.ClientProtocol41 |
	mysql.ClientTransactions | mysql.ClientSecureConnection

//下面的会根据配置文件参数加进去
//mysql.ClientPluginAuth

var baseConnID uint32 = 10000

const initClientConnStatus = mysql.ServerStatusAutocommit

// Session means session between client and proxy
type Session struct {
	sync.Mutex

	c     *ClientConn
	proxy *Server

	manager *Manager

	namespace string

	executor *SessionExecutor

	closed atomic.Value
}

// create session between client<->proxy
func newSession(s *Server, co net.Conn) *Session {
	cc := new(Session)
	tcpConn := co.(*net.TCPConn)

	//SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	//I set this option false.
	tcpConn.SetNoDelay(true)
	cc.c = NewClientConn(mysql.NewConn(tcpConn), s.manager)
	cc.proxy = s
	cc.manager = s.manager

	cc.c.SetConnectionID(atomic.AddUint32(&baseConnID, 1))
	cc.c.proxy = s

	cc.executor = newSessionExecutor(s.manager)
	cc.executor.clientAddr = co.RemoteAddr().String()
	cc.closed.Store(false)
	return cc
}

func (cc *Session) getNamespace() *Namespace {
	return cc.manager.GetNamespace(cc.namespace)
}

// IsAllowConnect check if allow to connect
func (cc *Session) IsAllowConnect() bool {
	ns := cc.getNamespace() // maybe nil, and panic!
	clientHost, _, err := net.SplitHostPort(cc.c.RemoteAddr().String())
	if err != nil {
		log.Warn("[server] Session parse host error: %v", err)
	}
	clientIP := net.ParseIP(clientHost)

	return ns.IsClientIPAllowed(clientIP)
}

// Handshake with client
// step1: server send plain handshake packets to client
// step2: client send handshake response packets to server
// step3: server send ok/err packets to client
func (cc *Session) Handshake() error {
	// First build and send the server handshake packet.
	if err := cc.c.writeInitialHandshakeV10(); err != nil {
		clientHost, _, innerErr := net.SplitHostPort(cc.c.RemoteAddr().String())
		if innerErr != nil {
			log.Warn("[server] Session parse host error: %v", innerErr)
		}
		// filter lvs detect liveness
		hostname, _ := util.HostName(clientHost)
		if len(hostname) > 0 && strings.Contains(hostname, "lvs") {
			return err
		}

		log.Warn("[server] Session writeInitialHandshake error, connId: %d, ip: %s, msg: %s, error: %s",
			cc.c.GetConnectionID(), clientHost, " send initial handshake error", err.Error())
		return err
	}

	info, err := cc.c.readHandshakeResponse()
	if err != nil {
		clientHost, _, innerErr := net.SplitHostPort(cc.c.RemoteAddr().String())
		if innerErr != nil {
			log.Warn("[server] Session parse host error: %v", innerErr)
		}
		// filter lvs detect liveness
		hostname, _ := util.HostName(clientHost)
		if len(hostname) > 0 && strings.Contains(hostname, "lvs") {
			return err
		}

		log.Warn("[server] Session readHandshakeResponse error, connId: %d, ip: %s, msg: %s, error: %s",
			cc.c.GetConnectionID(), clientHost, "read Handshake Response error", err.Error())
		return err
	}

	if err := cc.handleHandshakeResponse(info); err != nil {
		log.Warn("handleHandshakeResponse error, connId: %d, err: %v", cc.c.GetConnectionID(), err)
		return err
	}

	if err := cc.c.writeOK(cc.executor.GetStatus()); err != nil {
		log.Warn("[server] Session readHandshakeResponse error, connId %d, msg: %s, error: %s",
			cc.c.GetConnectionID(), "write ok fail", err.Error())
		return err
	}

	return nil
}

func (cc *Session) handleHandshakeResponse(info HandshakeResponseInfo) error {
	// check and set user
	var password string
	var succ bool
	user := info.User
	if !cc.manager.CheckUser(user) {
		return mysql.NewDefaultError(mysql.ErrAccessDenied, user, cc.c.RemoteAddr().String(), "Yes")
	}
	cc.executor.user = user

	// check password
	if len(info.AuthPlugin) == 0 {
		if len(info.AuthResponse) == 32 {
			succ, password = cc.manager.CheckSha2Password(user, info.Salt, info.AuthResponse)
		} else {
			succ, password = cc.manager.CheckPassword(user, info.Salt, info.AuthResponse)
		}
	} else if info.AuthPlugin == mysql.CachingSHA2Password {
		succ, password = cc.manager.CheckSha2Password(user, info.Salt, info.AuthResponse)
	} else {
		succ, password = cc.manager.CheckPassword(user, info.Salt, info.AuthResponse)
	}

	if !succ {
		return mysql.NewDefaultError(mysql.ErrAccessDenied, user, cc.c.RemoteAddr().String(), "Yes")
	}

	// handle collation
	collationID := info.CollationID
	collationName, ok := mysql.Collations[mysql.CollationID(collationID)]
	if !ok {
		return mysql.NewError(mysql.ErrInternal, "invalid collation")
	}
	charset, ok := mysql.CollationNameToCharset[collationName]
	if !ok {
		return mysql.NewError(mysql.ErrInternal, "invalid collation")
	}
	cc.executor.SetCollationID(mysql.CollationID(collationID))
	cc.executor.SetCharset(charset)

	// set database
	cc.executor.SetDatabase(info.Database)

	// set namespace
	namespace := cc.manager.GetNamespaceByUser(user, password)
	cc.namespace = namespace
	cc.executor.namespace = namespace
	cc.c.namespace = namespace // TODO: remove it when refactor is done
	return nil
}

// Close close session with it's resources
func (cc *Session) Close() {
	if cc.IsClosed() {
		return
	}
	cc.closed.Store(true)
	if err := cc.executor.rollback(); err != nil {
		log.Warn("executor rollback error when Session close: %v", err)
	}
	cc.c.Close()
	log.Debug("client closed, %d", cc.c.GetConnectionID())

	return
}

// IsClosed check if closed
func (cc *Session) IsClosed() bool {
	return cc.closed.Load().(bool)
}

// Run start session to server client request packets
func (cc *Session) Run() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			log.Warn("[server] Session Run panic error, error: %s, stack: %s", err.Error(), string(buf))
		}
		cc.Close()
		cc.proxy.tw.Remove(cc)
		cc.manager.GetStatisticManager().DescSessionCount(cc.namespace)
	}()

	cc.manager.GetStatisticManager().IncrSessionCount(cc.namespace)

	for !cc.IsClosed() {
		cc.c.SetSequence(0)
		data, err := cc.c.ReadEphemeralPacket()
		if err != nil {
			cc.c.RecycleReadPacket()
			return
		}

		cc.proxy.tw.Add(cc.proxy.sessionTimeout, cc, cc.Close)
		cc.manager.GetStatisticManager().AddReadFlowCount(cc.namespace, len(data))

		cmd := data[0]
		data = data[1:]
		rs := cc.executor.ExecuteCommand(cmd, data)
		cc.c.RecycleReadPacket()

		if err = cc.writeResponse(rs); err != nil {
			log.Warn("Session write response error, connId: %d, err: %v", cc.c.GetConnectionID(), err)
			cc.Close()
			return
		}

		if cmd == mysql.ComQuit {
			cc.Close()
		}
	}
}

func (cc *Session) writeResponse(r Response) error {
	switch r.RespType {
	case RespEOF:
		return cc.c.writeEOFPacket(r.Status)
	case RespResult:
		rs := r.Data.(*mysql.Result)
		if rs == nil {
			return cc.c.writeOK(r.Status)
		}
		return cc.c.writeOKResult(r.Status, r.Data.(*mysql.Result))
	case RespPrepare:
		stmt := r.Data.(*Stmt)
		if stmt == nil {
			return cc.c.writeOK(r.Status)
		}
		return cc.c.writePrepareResponse(r.Status, stmt)
	case RespFieldList:
		rs := r.Data.([]*mysql.Field)
		if rs == nil {
			return cc.c.writeOK(r.Status)
		}
		return cc.c.writeFieldList(r.Status, rs)
	case RespError:
		rs := r.Data.(error)
		if rs == nil {
			return cc.c.writeOK(r.Status)
		}
		err := cc.c.writeErrorPacket(rs)
		if err != nil {
			return err
		}
		if rs == mysql.ErrBadConn { // 后端连接如果断开, 应该返回通知Session关闭
			return rs
		}
		return nil
	case RespOK:
		return cc.c.writeOK(r.Status)
	case RespNoop:
		return nil
	default:
		err := fmt.Errorf("invalid response type: %T", r)
		log.Fatal(err.Error())
		return cc.c.writeErrorPacket(err)
	}
}
