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
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/log/xlog"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/stats"
	"github.com/XiaoMi/Gaea/stats/prometheus"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/sync2"
)

// LoadAndCreateManager load namespace config, and create manager
func LoadAndCreateManager(cfg *models.Proxy) (*Manager, error) {
	namespaceConfigs, err := loadAllNamespace(cfg)
	if err != nil {
		log.Warn("init namespace manager failed, %v", err)
		return nil, err

	}

	mgr, err := CreateManager(cfg, namespaceConfigs)
	if err != nil {
		log.Warn("create manager error: %v", err)
		return nil, err
	}
	//globalManager = mgr
	return mgr, nil
}

func loadAllNamespace(cfg *models.Proxy) (map[string]*models.Namespace, error) {
	// get names of all namespace
	root := cfg.CoordinatorRoot
	if cfg.ConfigType == models.ConfigFile {
		root = cfg.FileConfigPath
	}

	client := models.NewClient(cfg.ConfigType, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, root)
	store := models.NewStore(client)
	defer store.Close()
	var err error
	var names []string
	names, err = store.ListNamespace()
	if err != nil {
		log.Warn("list namespace failed, err: %v", err)
		return nil, err
	}

	// query remote namespace models in worker goroutines
	nameC := make(chan string)
	namespaceC := make(chan *models.Namespace)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			client := models.NewClient(cfg.ConfigType, cfg.CoordinatorAddr, cfg.UserName, cfg.Password, root)
			store := models.NewStore(client)
			defer store.Close()
			defer wg.Done()
			for name := range nameC {
				namespace, e := store.LoadNamespace(cfg.EncryptKey, name)
				if e != nil {
					log.Warn("load namespace %s failed, err: %v", name, err)
					// assign extent err out of this scope
					err = e
					return
				}
				// verify namespace config
				e = namespace.Verify()
				if e != nil {
					log.Warn("verify namespace %s failed, err: %v", name, e)
					err = e
					return
				}
				namespaceC <- namespace
			}
		}()
	}

	// dispatch goroutine
	go func() {
		for _, name := range names {
			nameC <- name
		}
		close(nameC)
		wg.Wait()
		close(namespaceC)
	}()

	// collect all namespaces
	namespaceModels := make(map[string]*models.Namespace, 64)
	for namespace := range namespaceC {
		namespaceModels[namespace.Name] = namespace
	}
	if err != nil {
		log.Warn("get namespace failed, err:%v", err)
		return nil, err
	}

	return namespaceModels, nil
}

// Manager contains namespace manager and user manager
type Manager struct {
	reloadPrepared sync2.AtomicBool
	switchIndex    util.BoolIndex
	namespaces     [2]*NamespaceManager
	users          [2]*UserManager
	statistics     *StatisticManager
}

// NewManager return empty Manager
func NewManager() *Manager {
	return &Manager{}
}

// CreateManager create manager
func CreateManager(cfg *models.Proxy, namespaceConfigs map[string]*models.Namespace) (*Manager, error) {
	m := NewManager()

	// init statistics
	statisticManager, err := CreateStatisticManager(cfg, m)
	if err != nil {
		log.Warn("init stats manager failed, %v", err)
		return nil, err
	}
	m.statistics = statisticManager

	current, _, _ := m.switchIndex.Get()

	// init namespace
	m.namespaces[current] = CreateNamespaceManager(namespaceConfigs)

	// init user
	user, err := CreateUserManager(namespaceConfigs)
	if err != nil {
		return nil, err
	}
	m.users[current] = user

	m.startConnectPoolMetricsTask(cfg.StatsInterval)
	return m, nil
}

// Close close manager
func (m *Manager) Close() {
	current, _, _ := m.switchIndex.Get()

	namespaces := m.namespaces[current].namespaces
	for _, ns := range namespaces {
		ns.Close(false)
	}

	m.statistics.Close()
}

// ReloadNamespacePrepare prepare commit
func (m *Manager) ReloadNamespacePrepare(namespaceConfig *models.Namespace) error {
	name := namespaceConfig.Name
	current, other, _ := m.switchIndex.Get()

	// reload namespace prepare
	currentNamespaceManager := m.namespaces[current]
	newNamespaceManager := ShallowCopyNamespaceManager(currentNamespaceManager)
	if err := newNamespaceManager.RebuildNamespace(namespaceConfig); err != nil {
		log.Warn("prepare config of namespace: %s failed, err: %v", name, err)
		return err
	}
	m.namespaces[other] = newNamespaceManager

	// reload user prepare
	currentUserManager := m.users[current]
	newUserManager := CloneUserManager(currentUserManager)
	newUserManager.RebuildNamespaceUsers(namespaceConfig)
	m.users[other] = newUserManager
	m.reloadPrepared.Set(true)

	return nil
}

// ReloadNamespaceCommit commit config
func (m *Manager) ReloadNamespaceCommit(name string) error {
	if !m.reloadPrepared.CompareAndSwap(true, false) {
		err := errors.ErrNamespaceNotPrepared
		log.Warn("commit namespace error, namespace: %s, err: %v", name, err)
		return err
	}

	current, _, index := m.switchIndex.Get()

	currentNamespace := m.namespaces[current].GetNamespace(name)
	if currentNamespace != nil {
		go currentNamespace.Close(true)
	}

	m.switchIndex.Set(!index)

	return nil
}

// DeleteNamespace delete namespace
func (m *Manager) DeleteNamespace(name string) error {
	current, other, index := m.switchIndex.Get()

	// idempotent delete
	currentNamespace := m.namespaces[current].GetNamespace(name)
	if currentNamespace == nil {
		return nil
	}

	// delete namespace of other
	currentNamespaceManager := m.namespaces[current]
	newNamespaceManager := ShallowCopyNamespaceManager(currentNamespaceManager)
	newNamespaceManager.DeleteNamespace(name)
	m.namespaces[other] = newNamespaceManager

	// delete users of other
	currentUserManager := m.users[current]
	newUserManager := CloneUserManager(currentUserManager)
	newUserManager.ClearNamespaceUsers(name)
	m.users[other] = newUserManager

	// switch namespace manager
	m.switchIndex.Set(!index)

	// delay recycle resources of current
	go currentNamespace.Close(true)

	return nil
}

// GetNamespace return specific namespace
func (m *Manager) GetNamespace(name string) *Namespace {
	current, _, _ := m.switchIndex.Get()
	return m.namespaces[current].GetNamespace(name)
}

// CheckUser check if user in users
func (m *Manager) CheckUser(user string) bool {
	current, _, _ := m.switchIndex.Get()
	return m.users[current].CheckUser(user)
}

// CheckPassword check if right password with specific user
func (m *Manager) CheckPassword(user string, salt, auth []byte) (bool, string) {
	current, _, _ := m.switchIndex.Get()
	return m.users[current].CheckPassword(user, salt, auth)
}

// CheckPassword check if right password with specific user
func (m *Manager) CheckSha2Password(user string, salt, auth []byte) (bool, string) {
	current, _, _ := m.switchIndex.Get()
	return m.users[current].CheckSha2Password(user, salt, auth)
}

// GetStatisticManager return proxy status to record status
func (m *Manager) GetStatisticManager() *StatisticManager {
	return m.statistics
}

// GetNamespaceByUser return namespace by user
func (m *Manager) GetNamespaceByUser(userName, password string) string {
	current, _, _ := m.switchIndex.Get()
	return m.users[current].GetNamespaceByUser(userName, password)
}

// ConfigFingerprint return config fingerprint
func (m *Manager) ConfigFingerprint() string {
	current, _, _ := m.switchIndex.Get()
	return m.namespaces[current].ConfigFingerprint()
}

// RecordSessionSQLMetrics record session SQL metrics, like response time, error
func (m *Manager) RecordSessionSQLMetrics(reqCtx *util.RequestContext, se *SessionExecutor, sql string, startTime time.Time, err error) {
	trimmedSql := strings.ReplaceAll(sql, "\n", " ")
	namespace := se.namespace
	ns := m.GetNamespace(namespace)
	if ns == nil {
		log.Warn("record session SQL metrics error, namespace: %s, sql: %s, err: %s", namespace, trimmedSql, "namespace not found")
		return
	}

	var operation string
	if stmtType, ok := reqCtx.Get(util.StmtType).(int); ok {
		operation = parser.StmtType(stmtType)
	} else {
		fingerprint := mysql.GetFingerprint(sql)
		operation = mysql.GetFingerprintOperation(fingerprint)
	}

	// record sql timing
	m.statistics.recordSessionSQLTiming(namespace, operation, startTime)

	// record slow sql
	duration := time.Since(startTime).Nanoseconds() / int64(time.Millisecond)
	if duration > ns.getSessionSlowSQLTime() || ns.getSessionSlowSQLTime() == 0 {
		log.Warn("session slow SQL, namespace: %s, sql: %s, cost: %d ms", namespace, trimmedSql, duration)
		fingerprint := mysql.GetFingerprint(sql)
		md5 := mysql.GetMd5(fingerprint)
		ns.SetSlowSQLFingerprint(md5, fingerprint)
		m.statistics.recordSessionSlowSQLFingerprint(namespace, md5)
	}

	// record error sql
	if err != nil {
		log.Warn("session error SQL, namespace: %s, sql: %s, cost: %d ms, err: %v", namespace, trimmedSql, duration, err)
		fingerprint := mysql.GetFingerprint(sql)
		md5 := mysql.GetMd5(fingerprint)
		ns.SetErrorSQLFingerprint(md5, fingerprint)
		m.statistics.recordSessionErrorSQLFingerprint(namespace, operation, md5)
	}

	if OpenProcessGeneralQueryLog() && ns.openGeneralLog {
		m.statistics.generalLogger.Notice("client: %s, namespace: %s, db: %s, user: %s, cmd: %s, sql: %s, cost: %d ms, succ: %t",
			se.clientAddr, namespace, se.db, se.user, operation, trimmedSql, duration, err == nil)
	}
}

// RecordBackendSQLMetrics record backend SQL metrics, like response time, error
func (m *Manager) RecordBackendSQLMetrics(reqCtx *util.RequestContext, namespace string, sql, backendAddr string, startTime time.Time, err error) {
	trimmedSql := strings.ReplaceAll(sql, "\n", " ")
	ns := m.GetNamespace(namespace)
	if ns == nil {
		log.Warn("record backend SQL metrics error, namespace: %s, backend addr: %s, sql: %s, err: %s", namespace, backendAddr, trimmedSql, "namespace not found")
		return
	}

	var operation string
	if stmtType, ok := reqCtx.Get(util.StmtType).(int); ok {
		operation = parser.StmtType(stmtType)
	} else {
		fingerprint := mysql.GetFingerprint(sql)
		operation = mysql.GetFingerprintOperation(fingerprint)
	}

	// record sql timing
	m.statistics.recordBackendSQLTiming(namespace, operation, startTime)

	// record slow sql
	duration := time.Since(startTime).Nanoseconds() / int64(time.Millisecond)
	if m.statistics.isBackendSlowSQL(startTime) {
		log.Warn("backend slow SQL, namespace: %s, addr: %s, sql: %s, cost: %d ms", namespace, backendAddr, trimmedSql, duration)
		fingerprint := mysql.GetFingerprint(sql)
		md5 := mysql.GetMd5(fingerprint)
		ns.SetBackendSlowSQLFingerprint(md5, fingerprint)
		m.statistics.recordBackendSlowSQLFingerprint(namespace, md5)
	}

	// record error sql
	if err != nil {
		log.Warn("backend error SQL, namespace: %s, addr: %s, sql: %s, cost %d ms, err: %v", namespace, backendAddr, trimmedSql, duration, err)
		fingerprint := mysql.GetFingerprint(sql)
		md5 := mysql.GetMd5(fingerprint)
		ns.SetBackendErrorSQLFingerprint(md5, fingerprint)
		m.statistics.recordBackendErrorSQLFingerprint(namespace, operation, md5)
	}
}

func (m *Manager) startConnectPoolMetricsTask(interval int) {
	if interval <= 0 {
		interval = 10
	}

	go func() {
		t := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-m.GetStatisticManager().closeChan:
				return
			case <-t.C:
				current, _, _ := m.switchIndex.Get()
				for nameSpaceName, _ := range m.namespaces[current].namespaces {
					m.recordBackendConnectPoolMetrics(nameSpaceName)
				}
			}
		}
	}()
}

func (m *Manager) recordBackendConnectPoolMetrics(namespace string) {
	ns := m.GetNamespace(namespace)
	if ns == nil {
		log.Warn("record backend connect pool metrics err, namespace: %s", namespace)
		return
	}

	for sliceName, slice := range ns.slices {
		m.statistics.recordConnectPoolInuseCount(namespace, sliceName, slice.Master.Addr(), slice.Master.InUse())
		m.statistics.recordConnectPoolIdleCount(namespace, sliceName, slice.Master.Addr(), slice.Master.Available())
		m.statistics.recordConnectPoolWaitCount(namespace, sliceName, slice.Master.Addr(), slice.Master.WaitCount())
		for _, slave := range slice.Slave {
			m.statistics.recordConnectPoolInuseCount(namespace, sliceName, slave.Addr(), slave.InUse())
			m.statistics.recordConnectPoolIdleCount(namespace, sliceName, slave.Addr(), slave.Available())
			m.statistics.recordConnectPoolWaitCount(namespace, sliceName, slave.Addr(), slave.WaitCount())
		}
		for _, statisticSlave := range slice.StatisticSlave {
			m.statistics.recordConnectPoolInuseCount(namespace, sliceName, statisticSlave.Addr(), statisticSlave.InUse())
			m.statistics.recordConnectPoolIdleCount(namespace, sliceName, statisticSlave.Addr(), statisticSlave.Available())
			m.statistics.recordConnectPoolWaitCount(namespace, sliceName, statisticSlave.Addr(), statisticSlave.WaitCount())
		}
	}
}

// NamespaceManager is the manager that holds all namespaces
type NamespaceManager struct {
	namespaces map[string]*Namespace
}

// NewNamespaceManager constructor of NamespaceManager
func NewNamespaceManager() *NamespaceManager {
	return &NamespaceManager{
		namespaces: make(map[string]*Namespace, 64),
	}
}

// CreateNamespaceManager create NamespaceManager
func CreateNamespaceManager(namespaceConfigs map[string]*models.Namespace) *NamespaceManager {
	nsMgr := NewNamespaceManager()
	for _, config := range namespaceConfigs {
		namespace, err := NewNamespace(config)
		if err != nil {
			log.Warn("create namespace %s failed, err: %v", config.Name, err)
			continue
		}
		nsMgr.namespaces[namespace.name] = namespace
	}
	return nsMgr
}

// ShallowCopyNamespaceManager copy NamespaceManager
func ShallowCopyNamespaceManager(nsMgr *NamespaceManager) *NamespaceManager {
	newNsMgr := NewNamespaceManager()
	for k, v := range nsMgr.namespaces {
		newNsMgr.namespaces[k] = v
	}
	return newNsMgr
}

// RebuildNamespace rebuild namespace
func (n *NamespaceManager) RebuildNamespace(config *models.Namespace) error {
	namespace, err := NewNamespace(config)
	if err != nil {
		log.Warn("create namespace %s failed, err: %v", config.Name, err)
		return err
	}
	n.namespaces[config.Name] = namespace
	return nil
}

// DeleteNamespace delete namespace
func (n *NamespaceManager) DeleteNamespace(ns string) {
	delete(n.namespaces, ns)
}

// GetNamespace get namespace in NamespaceManager
func (n *NamespaceManager) GetNamespace(namespace string) *Namespace {
	return n.namespaces[namespace]
}

// GetNamespaces return all namespaces in NamespaceManager
func (n *NamespaceManager) GetNamespaces() map[string]*Namespace {
	return n.namespaces
}

// ConfigFingerprint return config fingerprint
func (n *NamespaceManager) ConfigFingerprint() string {
	sortedKeys := make([]string, 0)
	for k := range n.GetNamespaces() {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	h := md5.New()
	for _, k := range sortedKeys {
		h.Write(n.GetNamespace(k).DumpToJSON())
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// UserManager means user for auth
// username+password是全局唯一的, 而username可以对应多个namespace
type UserManager struct {
	users          map[string][]string // key: user name, value: user password, same user may have different password, so array of passwords is needed
	userNamespaces map[string]string   // key: UserName+Password, value: name of namespace
}

// NewUserManager constructor of UserManager
func NewUserManager() *UserManager {
	return &UserManager{
		users:          make(map[string][]string, 64),
		userNamespaces: make(map[string]string, 64),
	}
}

// CreateUserManager create UserManager
func CreateUserManager(namespaceConfigs map[string]*models.Namespace) (*UserManager, error) {
	user := NewUserManager()
	for _, ns := range namespaceConfigs {
		user.addNamespaceUsers(ns)
	}
	return user, nil
}

// CloneUserManager close UserManager
func CloneUserManager(user *UserManager) *UserManager {
	ret := NewUserManager()
	// copy
	for k, v := range user.userNamespaces {
		ret.userNamespaces[k] = v
	}
	for k, v := range user.users {
		users := make([]string, len(v))
		copy(users, v)
		ret.users[k] = users
	}

	return ret
}

// RebuildNamespaceUsers rebuild users in namespace
func (u *UserManager) RebuildNamespaceUsers(namespace *models.Namespace) {
	u.ClearNamespaceUsers(namespace.Name)
	u.addNamespaceUsers(namespace)
}

// ClearNamespaceUsers clear users in namespace
func (u *UserManager) ClearNamespaceUsers(namespace string) {
	for key, ns := range u.userNamespaces {
		if ns == namespace {
			delete(u.userNamespaces, key)

			// delete user password in users
			username, password := getUserAndPasswordFromKey(key)
			var s []string
			for i := range u.users[username] {
				if u.users[username][i] == password {
					s = append(u.users[username][:i], u.users[username][i+1:]...)
				}
			}
			u.users[username] = s
		}
	}
}

func (u *UserManager) addNamespaceUsers(namespace *models.Namespace) {
	for _, user := range namespace.Users {
		key := getUserKey(user.UserName, user.Password)
		u.userNamespaces[key] = namespace.Name
		u.users[user.UserName] = append(u.users[user.UserName], user.Password)
	}
}

// CheckUser check if user in users
func (u *UserManager) CheckUser(user string) bool {
	if _, ok := u.users[user]; ok {
		return true
	}
	return false
}

// CheckPassword check if right password with specific user
func (u *UserManager) CheckPassword(user string, salt, auth []byte) (bool, string) {
	for _, password := range u.users[user] {
		checkAuth := mysql.CalcPassword(salt, []byte(password))
		if bytes.Equal(auth, checkAuth) {
			return true, password
		}
	}
	return false, ""
}

// CheckPassword check if right password with specific user
func (u *UserManager) CheckSha2Password(user string, salt, auth []byte) (bool, string) {
	for _, password := range u.users[user] {
		checkAuth := mysql.CalcCachingSha2Password(salt, password)
		if bytes.Equal(auth, checkAuth) {
			return true, password
		}
	}
	return false, ""
}

// GetNamespaceByUser return namespace by user
func (u *UserManager) GetNamespaceByUser(userName, password string) string {
	key := getUserKey(userName, password)
	if name, ok := u.userNamespaces[key]; ok {
		return name
	}
	return ""
}

func getUserKey(username, password string) string {
	return username + ":" + password
}

func getUserAndPasswordFromKey(key string) (username string, password string) {
	strs := strings.Split(key, ":")
	return strs[0], strs[1]
}

const (
	statsLabelCluster       = "Cluster"
	statsLabelOperation     = "Operation"
	statsLabelNamespace     = "Namespace"
	statsLabelFingerprint   = "Fingerprint"
	statsLabelFlowDirection = "Flowdirection"
	statsLabelSlice         = "Slice"
	statsLabelIPAddr        = "IPAddr"
)

// StatisticManager statistics manager
type StatisticManager struct {
	manager     *Manager
	clusterName string

	statsType     string // 监控后端类型
	handlers      map[string]http.Handler
	generalLogger log.Logger

	sqlTimings                *stats.MultiTimings            // SQL耗时统计
	sqlFingerprintSlowCounts  *stats.CountersWithMultiLabels // 慢SQL指纹数量统计
	sqlErrorCounts            *stats.CountersWithMultiLabels // SQL错误数统计
	sqlFingerprintErrorCounts *stats.CountersWithMultiLabels // SQL指纹错误数统计
	sqlForbidenCounts         *stats.CountersWithMultiLabels // SQL黑名单请求统计
	flowCounts                *stats.CountersWithMultiLabels // 业务流量统计
	sessionCounts             *stats.GaugesWithMultiLabels   // 前端会话数统计

	backendSQLTimings                *stats.MultiTimings            // 后端SQL耗时统计
	backendSQLFingerprintSlowCounts  *stats.CountersWithMultiLabels // 后端慢SQL指纹数量统计
	backendSQLErrorCounts            *stats.CountersWithMultiLabels // 后端SQL错误数统计
	backendSQLFingerprintErrorCounts *stats.CountersWithMultiLabels // 后端SQL指纹错误数统计
	backendConnectPoolIdleCounts     *stats.GaugesWithMultiLabels   //后端空闲连接数统计
	backendConnectPoolInUseCounts    *stats.GaugesWithMultiLabels   //后端正在使用连接数统计
	backendConnectPoolWaitCounts     *stats.GaugesWithMultiLabels   //后端等待队列统计

	slowSQLTime int64
	closeChan   chan bool
}

// NewStatisticManager return empty StatisticManager
func NewStatisticManager() *StatisticManager {
	return &StatisticManager{}
}

// CreateStatisticManager create StatisticManager
func CreateStatisticManager(cfg *models.Proxy, manager *Manager) (*StatisticManager, error) {
	mgr := NewStatisticManager()
	mgr.manager = manager
	mgr.clusterName = cfg.Cluster

	var err error
	if err = mgr.Init(cfg); err != nil {
		return nil, err
	}
	if mgr.generalLogger, err = initGeneralLogger(cfg); err != nil {
		return nil, err
	}
	return mgr, nil
}

type proxyStatsConfig struct {
	Service      string
	StatsEnabled bool
}

func initGeneralLogger(cfg *models.Proxy) (log.Logger, error) {
	c := make(map[string]string, 5)
	c["path"] = cfg.LogPath
	c["filename"] = cfg.LogFileName + "_sql"
	c["level"] = cfg.LogLevel
	c["service"] = cfg.Service
	c["runtime"] = "false"
	return xlog.CreateLogManager(cfg.LogOutput, c)
}

func parseProxyStatsConfig(cfg *models.Proxy) (*proxyStatsConfig, error) {
	enabled, err := strconv.ParseBool(cfg.StatsEnabled)
	if err != nil {
		return nil, err
	}

	statsConfig := &proxyStatsConfig{
		Service:      cfg.Service,
		StatsEnabled: enabled,
	}
	return statsConfig, nil
}

// Init init StatisticManager
func (s *StatisticManager) Init(cfg *models.Proxy) error {
	s.closeChan = make(chan bool, 0)
	s.handlers = make(map[string]http.Handler)
	s.slowSQLTime = cfg.SlowSQLTime
	statsCfg, err := parseProxyStatsConfig(cfg)
	if err != nil {
		return err
	}

	if err := s.initBackend(statsCfg); err != nil {
		return err
	}

	s.sqlTimings = stats.NewMultiTimings("SqlTimings",
		"gaea proxy sql sqlTimings", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	s.sqlFingerprintSlowCounts = stats.NewCountersWithMultiLabels("SqlFingerprintSlowCounts",
		"gaea proxy sql fingerprint slow counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	s.sqlErrorCounts = stats.NewCountersWithMultiLabels("SqlErrorCounts",
		"gaea proxy sql error counts per error type", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	s.sqlFingerprintErrorCounts = stats.NewCountersWithMultiLabels("SqlFingerprintErrorCounts",
		"gaea proxy sql fingerprint error counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	s.sqlForbidenCounts = stats.NewCountersWithMultiLabels("SqlForbiddenCounts",
		"gaea proxy sql error counts per error type", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	s.flowCounts = stats.NewCountersWithMultiLabels("FlowCounts",
		"gaea proxy flow counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFlowDirection})
	s.sessionCounts = stats.NewGaugesWithMultiLabels("SessionCounts",
		"gaea proxy session counts", []string{statsLabelCluster, statsLabelNamespace})

	s.backendSQLTimings = stats.NewMultiTimings("BackendSqlTimings",
		"gaea proxy backend sql sqlTimings", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	s.backendSQLFingerprintSlowCounts = stats.NewCountersWithMultiLabels("BackendSqlFingerprintSlowCounts",
		"gaea proxy backend sql fingerprint slow counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	s.backendSQLErrorCounts = stats.NewCountersWithMultiLabels("BackendSqlErrorCounts",
		"gaea proxy backend sql error counts per error type", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	s.backendSQLFingerprintErrorCounts = stats.NewCountersWithMultiLabels("BackendSqlFingerprintErrorCounts",
		"gaea proxy backend sql fingerprint error counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	s.backendConnectPoolIdleCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolIdleCounts",
		"gaea proxy backend idle connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})
	s.backendConnectPoolInUseCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolInUseCounts",
		"gaea proxy backend in-use connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})
	s.backendConnectPoolWaitCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolWaitCounts",
		"gaea proxy backend wait connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})

	s.startClearTask()
	return nil
}

// Close close proxy stats
func (s *StatisticManager) Close() {
	close(s.closeChan)
}

// GetHandlers return specific handler of stats
func (s *StatisticManager) GetHandlers() map[string]http.Handler {
	return s.handlers
}

func (s *StatisticManager) initBackend(cfg *proxyStatsConfig) error {
	prometheus.Init(cfg.Service)
	s.handlers = prometheus.GetHandlers()
	return nil
}

// clear data to prevent
func (s *StatisticManager) startClearTask() {
	go func() {
		t := time.NewTicker(time.Hour)
		for {
			select {
			case <-s.closeChan:
				return
			case <-t.C:
				s.clearLargeCounters()
			}
		}
	}()
}

func (s *StatisticManager) clearLargeCounters() {
	s.sqlErrorCounts.ResetAll()
	s.sqlFingerprintSlowCounts.ResetAll()
	s.sqlFingerprintErrorCounts.ResetAll()

	s.backendSQLErrorCounts.ResetAll()
	s.backendSQLFingerprintSlowCounts.ResetAll()
	s.backendSQLFingerprintErrorCounts.ResetAll()
}

func (s *StatisticManager) recordSessionSlowSQLFingerprint(namespace string, md5 string) {
	fingerprintStatsKey := []string{s.clusterName, namespace, md5}
	s.sqlFingerprintSlowCounts.Add(fingerprintStatsKey, 1)
}

func (s *StatisticManager) recordSessionErrorSQLFingerprint(namespace string, operation string, md5 string) {
	fingerprintStatsKey := []string{s.clusterName, namespace, md5}
	operationStatsKey := []string{s.clusterName, namespace, operation}
	s.sqlErrorCounts.Add(operationStatsKey, 1)
	s.sqlFingerprintErrorCounts.Add(fingerprintStatsKey, 1)
}

func (s *StatisticManager) recordSessionSQLTiming(namespace string, operation string, startTime time.Time) {
	operationStatsKey := []string{s.clusterName, namespace, operation}
	s.sqlTimings.Record(operationStatsKey, startTime)
}

// millisecond duration
func (s *StatisticManager) isBackendSlowSQL(startTime time.Time) bool {
	duration := time.Since(startTime).Nanoseconds() / int64(time.Millisecond)
	return duration > s.slowSQLTime || s.slowSQLTime == 0
}

func (s *StatisticManager) recordBackendSlowSQLFingerprint(namespace string, md5 string) {
	fingerprintStatsKey := []string{s.clusterName, namespace, md5}
	s.backendSQLFingerprintSlowCounts.Add(fingerprintStatsKey, 1)
}

func (s *StatisticManager) recordBackendErrorSQLFingerprint(namespace string, operation string, md5 string) {
	fingerprintStatsKey := []string{s.clusterName, namespace, md5}
	operationStatsKey := []string{s.clusterName, namespace, operation}
	s.backendSQLErrorCounts.Add(operationStatsKey, 1)
	s.backendSQLFingerprintErrorCounts.Add(fingerprintStatsKey, 1)
}

func (s *StatisticManager) recordBackendSQLTiming(namespace string, operation string, startTime time.Time) {
	operationStatsKey := []string{s.clusterName, namespace, operation}
	s.backendSQLTimings.Record(operationStatsKey, startTime)
}

// RecordSQLForbidden record forbidden sql
func (s *StatisticManager) RecordSQLForbidden(fingerprint, namespace string) {
	md5 := mysql.GetMd5(fingerprint)
	s.sqlForbidenCounts.Add([]string{s.clusterName, namespace, md5}, 1)
}

// IncrSessionCount incr session count
func (s *StatisticManager) IncrSessionCount(namespace string) {
	statsKey := []string{s.clusterName, namespace}
	s.sessionCounts.Add(statsKey, 1)
}

// DescSessionCount decr session count
func (s *StatisticManager) DescSessionCount(namespace string) {
	statsKey := []string{s.clusterName, namespace}
	s.sessionCounts.Add(statsKey, -1)
}

// AddReadFlowCount add read flow count
func (s *StatisticManager) AddReadFlowCount(namespace string, byteCount int) {
	statsKey := []string{s.clusterName, namespace, "read"}
	s.flowCounts.Add(statsKey, int64(byteCount))
}

// AddWriteFlowCount add write flow count
func (s *StatisticManager) AddWriteFlowCount(namespace string, byteCount int) {
	statsKey := []string{s.clusterName, namespace, "write"}
	s.flowCounts.Add(statsKey, int64(byteCount))
}

//record idle connect count
func (s *StatisticManager) recordConnectPoolIdleCount(namespace string, slice string, addr string, count int64) {
	statsKey := []string{s.clusterName, namespace, slice, addr}
	s.backendConnectPoolIdleCounts.Set(statsKey, count)
}

//record in-use connect count
func (s *StatisticManager) recordConnectPoolInuseCount(namespace string, slice string, addr string, count int64) {
	statsKey := []string{s.clusterName, namespace, slice, addr}
	s.backendConnectPoolInUseCounts.Set(statsKey, count)
}

//record wait queue length
func (s *StatisticManager) recordConnectPoolWaitCount(namespace string, slice string, addr string, count int64) {
	statsKey := []string{s.clusterName, namespace, slice, addr}
	s.backendConnectPoolWaitCounts.Set(statsKey, count)
}
