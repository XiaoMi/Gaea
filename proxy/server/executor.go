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
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/XiaoMi/Gaea/models"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/proxy/plan"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/hack"
)

const (
	// master comments
	masterComment      = "/*master*/"
	masterHint         = "*master*"
	mycatHint          = "/* !mycat:"
	standardMasterHint = "/*+ master */"
	// general query log variable
	gaeaGeneralLogVariable   = "gaea_general_log"
	readonlyVariable         = "read_only"
	globalReadonlyVariable   = "global.read_only"
	TxReadonlyLT5720         = "@@tx_read_only"
	TxReadonlyGT5720         = "@@transaction_read_only"
	SessionTxReadonlyLT5720  = "@@session.tx_read_only"
	SessionTxReadonlyGT5720  = "@@session.transaction_read_only"
	TxIsolationLT5720        = "@@tx_isolation"
	TxIsolationGT5720        = "@@transaction_isolation"
	SessionTxIsolationLT5720 = "@@session.tx_isolation"
	SessionTxIsolationGT5720 = "@@session.transaction_isolation"
	// JdbcInitPrefix jdbc prefix: /* mysql-connector-java (<8.0.30); /* mysql-connector-j-8...(>8.0.30)
	JdbcInitPrefix = "/* mysql-connector-j"

	// multiBackendAddrMark marks the backend addr is one of multi backend addrs
	multiBackendAddrMark = ">"

	// select last_insert_id()
	lastInsetIdMark = "SELECTLAST_INSERT_ID()"
)

// SessionExecutor is bound to a session, so requests are serializable
type SessionExecutor struct {
	manager *Manager

	namespace  string
	user       string
	db         string
	clientAddr string

	status       uint16
	lastInsertID uint64

	collation        mysql.CollationID
	charset          string
	sessionVariables *mysql.SessionVariables

	keepSession bool
	userPriv    int

	txConns          map[string]backend.PooledConnect
	ksConns          map[string]backend.PooledConnect // keep session connections
	nsChangeIndexOld uint32
	savepoints       []string
	txLock           sync.Mutex

	stmtID uint32
	stmts  map[uint32]*Stmt //prepare相关,client端到proxy的stmt

	session             *Session
	serverAddr          net.Addr
	backendAddr         string //记录执行 SQL 后端实例的地址
	backendConnectionId int64  //记录执行 SQL 后端实例的连接ID
	contextNamespace    *Namespace
}

// Response response info
type Response struct {
	RespType int
	Status   uint16
	Data     interface{}
	IsBinary bool
}

const (
	// RespOK means OK message
	RespOK = iota
	// RespResult means Result message
	RespResult
	// RespError means error message
	RespError
	// RespFieldList means field list message
	RespFieldList
	// RespPrepare prepare response message
	RespPrepare
	// RespEOF means EOF message
	RespEOF
	// RespNoop means empty message
	RespNoop
)

// CreateOKResponse create ok response
func CreateOKResponse(status uint16) Response {
	return Response{
		RespType: RespOK,
		Status:   status,
	}
}

// CreateResultResponse create result response
func CreateResultResponse(status uint16, result *mysql.Result, isBinary bool) Response {
	return Response{
		RespType: RespResult,
		Status:   status,
		Data:     result,
		IsBinary: isBinary,
	}
}

// CreateErrorResponse create error response
func CreateErrorResponse(status uint16, err error) Response {
	return Response{
		RespType: RespError,
		Status:   status,
		Data:     err,
	}
}

// CreateFieldListResponse create field list response
func CreateFieldListResponse(status uint16, fl []*mysql.Field) Response {
	return Response{
		RespType: RespFieldList,
		Status:   status,
		Data:     fl,
	}
}

// CreatePrepareResponse create prepare response
func CreatePrepareResponse(status uint16, stmt *Stmt) Response {
	return Response{
		RespType: RespPrepare,
		Status:   status,
		Data:     stmt,
	}
}

// CreateEOFResponse create eof response
func CreateEOFResponse(status uint16) Response {
	return Response{
		RespType: RespEOF,
		Status:   status,
	}
}

// CreateNoopResponse no op response, for ComStmtClose
func CreateNoopResponse() Response {
	return Response{
		RespType: RespNoop,
	}
}

func newSessionExecutor(manager *Manager) *SessionExecutor {

	return &SessionExecutor{
		sessionVariables: mysql.NewSessionVariables(),
		txConns:          make(map[string]backend.PooledConnect),
		ksConns:          make(map[string]backend.PooledConnect),
		stmts:            make(map[uint32]*Stmt),
		status:           initClientConnStatus,
		manager:          manager,
	}
}

// GetNamespace return namespace in session
func (se *SessionExecutor) GetNamespace() *Namespace {
	return se.contextNamespace
}

// GetNamespace return namespace in session
func (se *SessionExecutor) SetContextNamespace() {
	se.contextNamespace = se.GetManagerNamespace()
}

func (se *SessionExecutor) GetManagerNamespace() *Namespace {
	return se.manager.GetNamespace(se.namespace)
}

// GetVariables return variables in session
func (se *SessionExecutor) GetVariables() *mysql.SessionVariables {
	return se.sessionVariables
}

func (se *SessionExecutor) setIntSessionVariable(name string, valueStr string) error {
	if strings.ToLower(valueStr) == mysql.KeywordDefault {
		se.sessionVariables.Delete(name)
		return nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return err
	}
	if err = se.sessionVariables.Set(name, value); err != nil {
		return err
	}
	return nil
}

func (se *SessionExecutor) setStringSessionVariable(name string, valueStr string) error {
	if strings.ToLower(valueStr) == mysql.KeywordDefault {
		se.sessionVariables.Delete(name)
		return nil
	}

	return se.sessionVariables.Set(name, valueStr)
}

func (se *SessionExecutor) setGeneralLogVariable(valueStr string) error {
	v, err := strconv.Atoi(valueStr)
	if err != nil {
		return errors.ErrInvalidArgument
	}
	atomic.StoreUint32(&ProcessGeneralLog, uint32(v))
	return nil
}

// GetLastInsertID return last_inert_id
func (se *SessionExecutor) GetLastInsertID() uint64 {
	return se.lastInsertID
}

// SetLastInsertID store last_insert_id
func (se *SessionExecutor) SetLastInsertID(id uint64) {
	se.lastInsertID = id
}

func (se *SessionExecutor) HandleSet(reqCtx *util.RequestContext, sql string, stmt *ast.SetStmt) (*mysql.Result, error) {
	return se.handleSet(reqCtx, sql, stmt)
}

// GetStatus return session status
func (se *SessionExecutor) GetStatus() uint16 {
	return se.status
}

// SetStatus store status
func (se *SessionExecutor) SetStatus(status uint16) {
	se.status = status
}

// SetCollationID store collation id
func (se *SessionExecutor) SetCollationID(id mysql.CollationID) {
	se.collation = id
}

// SetNamespaceDefaultCollationID store default collation id
func (se *SessionExecutor) SetNamespaceDefaultCollationID() {
	se.collation = se.manager.GetNamespace(se.namespace).GetDefaultCollationID()
}

// GetCollationID return collation id
func (se *SessionExecutor) GetCollationID() mysql.CollationID {
	return se.collation
}

// SetCharset set session charset
func (se *SessionExecutor) SetCharset(charset string) {
	se.charset = charset
}

// SetNamespaceDefaultCharset set session default charset
func (se *SessionExecutor) SetNamespaceDefaultCharset() {
	se.charset = se.manager.GetNamespace(se.namespace).GetDefaultCharset()
}

// GetCharset return charset
func (se *SessionExecutor) GetCharset() string {
	return se.charset
}

// SetDatabase set session database
func (se *SessionExecutor) SetDatabase(db string) {
	se.db = db
}

// GetDatabase return database in session
func (se *SessionExecutor) GetDatabase() string {
	return se.db
}

// IsKeepSession return keepSession flag
func (se *SessionExecutor) IsKeepSession() bool {
	return se.keepSession
}

// ExecuteCommand execute command
func (se *SessionExecutor) ExecuteCommand(cmd byte, data []byte) Response {
	switch cmd {
	case mysql.ComQuit:
		_ = se.manager.statistics.generalLogger.Notice("Quit - conn_id=%d, ns=%s, %s@%s/%s",
			se.session.c.ConnectionID, se.namespace, se.user, se.clientAddr, se.db)
		se.handleRollback(nil)
		// https://dev.mysql.com/doc/internals/en/com-quit.html
		// either a connection close or a OK_Packet, OK_Packet will cause client RST sometimes, but doesn't affect sql execute
		return CreateNoopResponse()
	case mysql.ComQuery: // data type: string[EOF]
		sql := string(data)
		// handle phase
		r, err := se.handleQuery(sql)
		if err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateResultResponse(se.status, r, false)
	case mysql.ComPing:
		if se.IsKeepSession() {
			if err := se.handleKeepSessionPing(); err != nil {
				return CreateErrorResponse(se.status, err)
			}
		}
		return CreateOKResponse(se.status)
	case mysql.ComInitDB:
		db := string(data)
		// handle phase
		err := se.handleUseDB(db)
		if err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateOKResponse(se.status)
	case mysql.ComFieldList:
		fs, err := se.handleFieldList(data)
		if err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateFieldListResponse(se.status, fs)
	case mysql.ComStmtPrepare:
		sql := string(data)
		stmt, err := se.handleStmtPrepare(sql)
		if err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreatePrepareResponse(se.status, stmt)
	case mysql.ComStmtExecute:
		values := make([]byte, len(data))
		copy(values, data)
		r, err := se.handleStmtExecute(values)
		if err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateResultResponse(se.status, r, true)
	case mysql.ComStmtClose: // no response
		if err := se.handleStmtClose(data); err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateNoopResponse()
	case mysql.ComStmtSendLongData: // no response
		values := make([]byte, len(data))
		copy(values, data)
		if err := se.handleStmtSendLongData(values); err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateNoopResponse()
	case mysql.ComStmtReset:
		if err := se.handleStmtReset(data); err != nil {
			return CreateErrorResponse(se.status, err)
		}
		return CreateOKResponse(se.status)
	case mysql.ComSetOption:
		return CreateEOFResponse(se.status)
	default:
		msg := fmt.Sprintf("command %d not supported now", cmd)
		log.Warn("dispatch command failed, error: %s", msg)
		return CreateErrorResponse(se.status, mysql.NewError(mysql.ErrUnknown, msg))
	}
}

func (se *SessionExecutor) getBackendConns(sqls map[string]map[string][]string, fromSlave bool) (pcs map[string]backend.PooledConnect, err error) {
	pcs = make(map[string]backend.PooledConnect)
	backendAddr := ""
	backendConnectionID := int64(0)

	for sliceName := range sqls {
		var pc backend.PooledConnect
		pc, err = se.getBackendConn(sliceName, fromSlave)
		if err != nil {
			return
		}
		pcs[sliceName] = pc
		backendAddr = pc.GetAddr()
		backendConnectionID = pc.GetConnectionID()
	}
	se.backendAddr = backendAddr
	se.backendConnectionId = backendConnectionID
	if len(pcs) > 1 {
		se.backendAddr = multiBackendAddrMark + backendAddr
	}
	return
}

func (se *SessionExecutor) getBackendConn(sliceName string, fromSlave bool) (pc backend.PooledConnect, err error) {
	if se.IsKeepSession() {
		return se.getBackendKsConn(sliceName)
	}
	return se.getBackendNoKsConn(sliceName, fromSlave)
}

func (se *SessionExecutor) getBackendNoKsConn(sliceName string, fromSlave bool) (pc backend.PooledConnect, err error) {
	if !se.isInTransaction() {
		slice := se.GetNamespace().GetSlice(sliceName)
		return slice.GetConn(fromSlave, se.GetNamespace().GetUserProperty(se.user), se.GetNamespace().localSlaveReadPriority)
	}
	return se.getTransactionConn(sliceName)
}

func (se *SessionExecutor) getBackendKsConn(sliceName string) (pc backend.PooledConnect, err error) {
	pc, ok := se.ksConns[sliceName]
	if ok {
		return pc, nil
	}

	slice := se.GetNamespace().GetSlice(sliceName)
	pc, err = slice.GetConn(se.userPriv == models.ReadOnly, se.GetNamespace().GetUserProperty(se.user), se.GetNamespace().localSlaveReadPriority)
	if err != nil {
		log.Warn("get connection from backend failed, error: %s", err.Error())
		return
	}

	if !se.isAutoCommit() {
		if err = pc.SetAutoCommit(0); err != nil {
			pc.Close()
			pc.Recycle()
			return
		}
	}

	if se.isInTransaction() {
		if err = pc.Begin(); err != nil {
			pc.Close()
			pc.Recycle()
			return
		}
	}

	se.ksConns[sliceName] = pc
	return
}

func (se *SessionExecutor) getTransactionConn(sliceName string) (pc backend.PooledConnect, err error) {
	se.txLock.Lock()
	defer se.txLock.Unlock()

	var ok bool
	if pc, ok = se.txConns[sliceName]; ok {
		return
	}

	slice := se.GetNamespace().GetSlice(sliceName) // returns nil only when the conf is error (fatal) so panic is correct
	if pc, err = slice.GetMasterConn(); err != nil {
		return
	}
	// Synchronize session variables before starting the transaction.
	// This step ensures that the session settings like `transaction_read_only` are correctly applied.
	// Setting session variables after `BEGIN` might not affect the transaction as expected,
	// since some session settings need to be established before the transaction starts.
	// pc.SetAutoCommit(0) is equivalent to starting a transaction
	if err = pc.SyncSessionVariables(se.sessionVariables); err != nil {
		pc.Close()
		pc.Recycle()
		return
	}
	if !se.isAutoCommit() {
		if err = pc.SetAutoCommit(0); err != nil {
			pc.Close()
			pc.Recycle()
			return
		}
	} else {
		if err = pc.Begin(); err != nil {
			pc.Close()
			pc.Recycle()
			return
		}
	}
	for _, savepoint := range se.savepoints {
		pc.Execute("savepoint "+savepoint, 0)
	}
	se.txConns[sliceName] = pc
	return
}

func (se *SessionExecutor) recycleBackendConn(pc backend.PooledConnect) {
	if pc == nil {
		return
	}

	if pc.IsClosed() {
		se.recycleTx()
		pc.Recycle()
		return
	}

	// if continueConn set to pc,maybe moreRowsExist or moreResultsExist
	if se.session.continueConn != nil && (pc.MoreRowsExist() || pc.MoreResultsExist()) {
		return
	}

	if se.IsKeepSession() {
		se.session.clearKsConns(se.nsChangeIndexOld)
		return
	}

	if se.isInTransaction() {
		return
	}

	pc.Recycle()
}

func (se *SessionExecutor) recycleContinueConn(pc backend.PooledConnect) {
	if pc == nil {
		return
	}
	if pc.IsClosed() {
		se.recycleTx()
		pc.Recycle()
		return
	}
	if se.IsKeepSession() {
		se.session.clearKsConns(se.nsChangeIndexOld)
		return
	}
	if se.isInTransaction() {
		return
	}
	pc.Recycle()
}

func (se *SessionExecutor) recycleBackendConns(pcs map[string]backend.PooledConnect, rollback bool) {
	if se.isInTransaction() || se.IsKeepSession() {
		return
	}

	for _, pc := range pcs {
		if pc == nil {
			continue
		}
		if rollback {
			pc.Rollback()
		}
		pc.Recycle()
	}
}

// initBackendConn tries to initialize the database connection with the specified database,
// charset, and session variables.
func initBackendConn(pc backend.PooledConnect, phyDB string, charset string, collation mysql.CollationID, sessionVariables *mysql.SessionVariables) error {
	if err := pc.UseDB(phyDB); err != nil {
		return err
	}
	return InitializeSessionVariables(pc, charset, collation, sessionVariables)
}

// InitializeSessionVariables sets the charset and session variables for the pooled connection.
// It attempts to write these settings and handles errors appropriately by closing the connection.
func InitializeSessionVariables(pc backend.PooledConnect, charset string, collation mysql.CollationID, sessionVariables *mysql.SessionVariables) error {
	charsetChanged, err := pc.SetCharset(charset, collation)
	if err != nil {
		return err
	}

	variablesChanged, err := pc.SetSessionVariables(sessionVariables)
	if err != nil {
		return err
	}

	if charsetChanged || variablesChanged {
		if err = pc.WriteSetStatement(); err != nil {
			log.Warn("set charset or session variables failed, address: %s, error: %s", pc.GetAddr(), err.Error())
			// Reset session variables to ensure the next use of the connection does not encounter incorrect settings or character set issues.
			// Resetting helps to address the root causes of session inconsistencies without masking them by simply pc.Close()
			sessionVariables.Reset(err)
			return err
		}
	}

	return nil
}

// executeShardSQLInSlice executes SQL queries in parallel across multiple slices,
// Each slice has an independent timeout. This method ensures that all operations are completed
// Execute within the maximum execution time, and terminate if the time limit is exceeded.
// It aggregates the results of each slice and handles errors and timeouts accordingly.
func (se *SessionExecutor) executeShardSQLInSlice(requestContext *util.RequestContext, pooledConnections map[string]backend.PooledConnect, sliceSQLs map[string]map[string][]string) (retRes []*mysql.Result, retErr error) {
	defer func() {
		// 用来捕获可能出现的panic
		if r := recover(); r != nil {
			log.Warn(
				"Recovered from panic in executeShardSQLInSlice goroutine, error: %v, stack: %s",
				r,
				debug.Stack(),
			)
			retErr = fmt.Errorf("caught executeShardSQLInSlice goroutine panic: %v", r)
		}
	}()

	numberOfSlices := len(pooledConnections)

	// 验证连接数和SQL分片数是否一致
	if numberOfSlices != len(sliceSQLs) {
		log.Warn(
			"Session executeShardSQLInSlice error: number of connections (%d) does not match number of SQL slices (%d)",
			numberOfSlices,
			len(sliceSQLs),
		)
		return nil, fmt.Errorf("number of connections (%d) does not match the number of SQL slices (%d)", numberOfSlices, len(sliceSQLs))
	}

	// 检查是否有分片需要执行
	if numberOfSlices == 0 {
		log.Warn(
			"Session executeShardSQLInSlice error: no SQL slices to execute",
			"connections", pooledConnections,
			"sqls", sliceSQLs,
			"error", errors.ErrNoPlan.Error(),
		)
		return nil, fmt.Errorf("no SQL slices to execute, error: %v", errors.ErrNoPlan)
	}

	// 创建一个带缓冲的通道，用于收集各个分片的执行结果
	executionResultsChan := make(chan sliceResult, numberOfSlices)

	// 使用 WaitGroup 确保所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(numberOfSlices)

	// 错误收集器
	var ec errorCollector

	// 启动所有分片的执行
	for sliceName, pooledConn := range pooledConnections {
		sqlStatementsMap := sliceSQLs[sliceName] // map[string][]string

		// 启动 goroutine 执行当前分片的 SQL
		go func(
			currentSliceName string,
			currentSqlStatementsMap map[string][]string,
			currentPooledConn backend.PooledConnect,
		) {
			defer wg.Done()
			rss, err := se.executeMultipleSQLInSlice(requestContext, currentSliceName, currentSqlStatementsMap, currentPooledConn)
			executionResultsChan <- sliceResult{
				sliceName: currentSliceName,
				results:   rss,
				err:       err,
			}
		}(sliceName, sqlStatementsMap, pooledConn)
	}

	// 启动一个 goroutine 等待所有执行完成后关闭通道
	go func() {
		wg.Wait()
		close(executionResultsChan)
	}()

	// 初始化一个空的 sliceResults 切片，用于存储所有执行结果
	sliceResults := make([]sliceResult, 0)
	for sliceExecutionResult := range executionResultsChan {
		if sliceExecutionResult.err != nil {
			ec.Collect(sliceExecutionResult.err)
			continue
		}
		sliceResults = append(sliceResults, sliceExecutionResult)
	}

	if allErr := ec.Error(); allErr != nil {
		// 直接返回所有错误
		return nil, allErr
	}

	sort.Slice(sliceResults, func(i, j int) bool {
		return sliceResults[i].sliceName < sliceResults[j].sliceName
	})

	// 依次追加每个分片的结果，保持顺序
	aggregatedResults := make([]*mysql.Result, 0)
	for _, result := range sliceResults {
		aggregatedResults = append(aggregatedResults, result.results...)
	}

	return aggregatedResults, nil
}

type sliceResult struct {
	sliceName string
	results   []*mysql.Result
	err       error
}

// executeMultipleSQLInSlice executes multiple SQL statements within a specified logical slice.
// The method takes a map of SQL statements grouped by database names and iterates over each database
// to execute the associated SQL statements. Each SQL execution is performed with proper timeout handling,
// and results are aggregated and returned.
func (se *SessionExecutor) executeMultipleSQLInSlice(requestContext *util.RequestContext, currentSliceName string, sqlStatementsMap map[string][]string, pooledConn backend.PooledConnect) (sliceResults []*mysql.Result, retErr error) {
	defer func() {
		// 用来捕获可能出现的panic
		if r := recover(); r != nil {
			log.Warn(
				"Recovered from panic in executeMultipleSQLInSlice goroutine, slice: %s, sqls: %v, error: %v, stack: %s",
				currentSliceName,
				sqlStatementsMap,
				r,
				debug.Stack(),
			)
			retErr = fmt.Errorf("caught executeMultipleSQLInSlice goroutine panic: %v in slice %s", r, currentSliceName)
		}
	}()
	// 获取命名空间的最大执行时间
	maxExecutionTime := se.manager.GetNamespace(se.namespace).GetMaxExecuteTime()

	var ctx context.Context
	var cancel context.CancelFunc
	if maxExecutionTime > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(maxExecutionTime)*time.Millisecond)
		defer cancel()
	} else {
		// 如果 maxExecutionTime <= 0，不使用超时限制
		ctx = context.Background()
	}

	// 获取并排序数据库名称
	databaseNames := make([]string, 0, len(sqlStatementsMap))
	for dbName := range sqlStatementsMap {
		databaseNames = append(databaseNames, dbName)
	}
	sort.Strings(databaseNames)

	// 遍历每个数据库执行对应的 SQL 语句
	for _, dbName := range databaseNames {
		sqlStatements := sqlStatementsMap[dbName]
		for _, sqlStatement := range sqlStatements {
			// 创建一个通道来接收 Execute 的结果
			execResultChan := make(chan executeResult)
			startTime := time.Now()

			// 1) 为当前这条 SQL 新开一个协程（go routine）去执行
			go func(sql string, begin time.Time) {
				queryResult, execErr := se.executeSingleSQLInSlice(pooledConn, currentSliceName, dbName, sql)
				execResultChan <- executeResult{
					result: queryResult,
					err:    execErr,
					start:  begin,
				}
			}(sqlStatement, startTime)

			// 2) 外层的 for 循环所在的协程用 select 等待 SQL执行结果或 SQL执行超时
			select {
			case <-ctx.Done():
				// SQL 执行超时，记录SQL
				se.manager.RecordBackendSQLMetrics(
					requestContext,
					se,
					currentSliceName,
					dbName,
					sqlStatement,
					pooledConn.GetAddr(),         // 正确传递 backendAddr
					pooledConn.GetConnectionID(), // 正确传递 mysql_connect_id
					startTime,
					fmt.Errorf("executeMultipleSQLInSlice execution SQL context done, maxExecutionTime: %d ms", maxExecutionTime),
				)
				// 触发超时，执行 kill
				if killErr := se.killSliceQueries(pooledConn, currentSliceName, sqlStatement); killErr != nil {
					log.Warn("failed to kill query error: %v", killErr)
				}
				return sliceResults, fmt.Errorf("slice: %s execution timed out, sql : %s", currentSliceName, sqlStatement)
			case execResult := <-execResultChan:
				// SQL 执行成功/失败，记录SQL
				se.manager.RecordBackendSQLMetrics(
					requestContext,
					se,
					currentSliceName,
					dbName,
					sqlStatement,
					pooledConn.GetAddr(),         // 正确传递 backendAddr
					pooledConn.GetConnectionID(), // 正确传递 mysql_connect_id
					execResult.start,
					execResult.err,
				)
				if execResult.err != nil {
					// SQL 执行失败，直接返回错误
					return sliceResults, execResult.err
				}
				sliceResults = append(sliceResults, execResult.result)
			}
		}
	}
	return sliceResults, nil
}

type executeResult struct {
	result *mysql.Result
	err    error
	start  time.Time
}

func (se *SessionExecutor) killSliceQueries(pc backend.PooledConnect, sliceName string, triggerSQL string) error {
	var (
		startTime           = time.Now()
		sessionConnectionId = se.session.c.GetConnectionID()
		address             = pc.GetAddr()
		backendConnectionId = pc.GetConnectionID()
		namespace           = se.namespace
	)
	// 构建公共信息
	msg := fmt.Sprintf("ns=%s, slice=%s, addr=%s, connect_id=%d, mysql_connect_id=%d, trigger_sql=%s",
		namespace,
		sliceName,
		address,             // 后端连接地址
		sessionConnectionId, // 前端连接 ID
		backendConnectionId, // 后端 MySQL 连接 ID
		triggerSQL,          // 触发 kill 的 SQL
	)

	dc, err := se.manager.GetNamespace(namespace).GetSlice(sliceName).GetDirectConn(address)
	if err != nil {
		return fmt.Errorf("failed to get direct connection error: %w, msg: %s", err, msg)
	}
	defer dc.Close()

	// 构造 KILL QUERY SQL
	killSQL := fmt.Sprintf("KILL QUERY %d", backendConnectionId)
	if _, err = dc.Execute(killSQL, 0); err != nil {
		return fmt.Errorf("failed to execute kill sql: %s, error: %w, msg: %s, duration: %s", killSQL, err, msg, time.Since(startTime))
	}
	log.Warn("successfully execute kill sql: %s, msg: %s, duration: %s", killSQL, msg, time.Since(startTime))
	return nil
}

// 安全地收集所有分片执行的错误
type errorCollector struct {
	sync.Mutex
	errors []error
}

// Collect adds a non-nil error to the collector.
func (ec *errorCollector) Collect(err error) {
	if err == nil {
		return
	}
	ec.Lock()
	defer ec.Unlock()
	ec.errors = append(ec.errors, err)
}

// Error returns a single error that combines all collected errors.
func (ec *errorCollector) Error() error {
	ec.Lock()
	defer ec.Unlock()
	if len(ec.errors) == 0 {
		return nil
	}
	var sb strings.Builder
	for _, err := range ec.errors {
		sb.WriteString(err.Error())
		sb.WriteString("; ")
	}
	return fmt.Errorf("%s", strings.TrimSuffix(sb.String(), "; "))
}

// executeUnshardSQLInSlice executes a SQL query on a non-sharded database and table scenario.
// It supports timeout and query kill functionality. Internally, it calls `executeSingleSQLInSlice` through a goroutine.
// Key Functionality:
// - It creates a context with a timeout based on the configured maximum execution time for the namespace.
// - It initializes the backend database connection if not already established.
// - It launches a separate goroutine to run the SQL query through `executeSingleSQLInSlice` on the physical database.
// - The method waits for either the completion of the query execution or for the timeout to occur.
// - In case of a timeout, it records the timeout and tries to kill the query on the backend server./
func (se *SessionExecutor) executeUnshardSQLInSlice(reqCtx *util.RequestContext, pooledConn backend.PooledConnect, phyDb, query string) (retRes *mysql.Result, retErr error) {
	defer func() {
		// 用来捕获initBackendConn可能出现的panic
		if r := recover(); r != nil {
			log.Warn(
				"Recovered from panic in executeUnShardSQLInSlice goroutine, error: %v, stack: %s",
				r,
				debug.Stack(),
			)
			retErr = fmt.Errorf("caught executeUnShardSQLInSlice goroutine panic: %v, db: %s, sql: %s", r, phyDb, query)
		}
	}()

	// 1. 创建上下文，可带超时
	maxExecTime := se.GetNamespace().GetMaxExecuteTime()
	ctx, cancel := context.WithCancel(context.Background())
	if maxExecTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(maxExecTime)*time.Millisecond)
	}
	defer cancel()

	// 用于接收 goroutine 执行结果的通道
	outCh := make(chan executeResult, 1)
	sliceName := "slice-0"
	startTime := time.Now()

	// 安全检查
	if pooledConn == nil {
		return nil, fmt.Errorf("no backend connection for slice: %s, db: %s", sliceName, phyDb)
	}

	// 2. 启动子协程，初始化数据库并执行单条SQL
	go func() {
		r, e := se.executeSingleSQLInSlice(pooledConn, sliceName, phyDb, query)
		outCh <- executeResult{result: r, err: e, start: startTime}
		close(outCh)
	}()

	// 3. 主协程等待执行结果 or 超时
	select {
	case exeRes := <-outCh:
		// (3a) 正常返回或执行出错
		se.manager.RecordBackendSQLMetrics(
			reqCtx,
			se,
			sliceName,
			phyDb,
			query,
			pooledConn.GetAddr(),
			pooledConn.GetConnectionID(),
			startTime,
			exeRes.err,
		)
		return exeRes.result, exeRes.err
	case <-ctx.Done():
		// (3b) 超时处理，先记录一次 "timeout" 指标
		se.manager.RecordBackendSQLMetrics(
			reqCtx,
			se,
			sliceName,
			phyDb,
			query,
			pooledConn.GetAddr(),
			pooledConn.GetConnectionID(),
			startTime,
			fmt.Errorf("timeout after %d ms", maxExecTime),
		)

		// kill 后端SQL
		if killErr := se.killSliceQueries(pooledConn, sliceName, query); killErr != nil {
			log.Warn("failed to kill query, error: %v", killErr)
		}
		// 关闭连接，防止复用有问题
		pooledConn.Close()
		return nil, fmt.Errorf("slice: %s execution timed out, sql: %s", sliceName, query)
	}
}

// executeSingleSQLInSlice is responsible for executing a SQL query on a specific physical database.
// Key Functionality:
// - It ensures the backend connection is initialized and ready for executing the SQL query.
// - It performs the SQL query execution on the physical database and retrieves the result.
// - Any errors encountered during the execution (e.g., invalid SQL, connection issues) are captured and returned.
//
// This method encapsulates the actual query execution logic on the physical database, including error handling and result retrieval.

func (se *SessionExecutor) executeSingleSQLInSlice(pooledConn backend.PooledConnect, sliceName, dbName, sql string) (retRes *mysql.Result, retErr error) {
	// 用于捕获本协程内的 panic
	defer func() {
		if r := recover(); r != nil {
			log.Warn(
				"Recovered from panic in executeSingleSQLInSlice goroutine, slice: %s, db: %s, sql: %s, error: %v, stack: %s",
				sliceName,
				dbName,
				sql,
				r,
				debug.Stack(),
			)
			retErr = fmt.Errorf("caught executeSingleSQLInSlice goroutine panic: %v, sql: %s", r, sql)
		}
	}()
	// 连接初始化
	if initErr := initBackendConn(pooledConn, dbName, se.GetCharset(), se.GetCollationID(), se.GetVariables()); initErr != nil {
		return nil, fmt.Errorf("failed to initialize backendConn, slice: %s, db: %s, error: %v", sliceName, dbName, initErr)
	}
	res, execErr := pooledConn.Execute(sql, se.GetNamespace().GetMaxResultSize())
	return res, execErr
}

func canHandleWithoutPlan(stmtType int) bool {
	return stmtType == parser.StmtShow ||
		stmtType == parser.StmtSet ||
		stmtType == parser.StmtBegin ||
		stmtType == parser.StmtCommit ||
		stmtType == parser.StmtRollback ||
		stmtType == parser.StmtSavepoint ||
		stmtType == parser.StmtUse ||
		stmtType == parser.StmtRelease ||
		stmtType == parser.StmeSRollback ||
		stmtType == parser.StmtLockTables ||
		stmtType == parser.StmtKill
}

const variableRestoreFlag = format.RestoreKeyWordLowercase | format.RestoreNameLowercase
const sqlModeRestoreFlag = format.RestoreStringSingleQuotes

// 获取SET语句中变量的字符串值, 去掉各种引号并转换为小写
func getVariableExprResult(v ast.ExprNode) string {
	s := &strings.Builder{}
	ctx := format.NewRestoreCtx(variableRestoreFlag, s)
	v.Restore(ctx)
	return strings.ToLower(s.String())
}

// 获取 SET 语句中变量的字符串值, 保留引号，不转换大小写
func getSqlModeExprResult(v ast.ExprNode) string {
	s := &strings.Builder{}
	ctx := format.NewRestoreCtx(sqlModeRestoreFlag, s)
	v.Restore(ctx)
	return s.String()
}

func getOnOffVariable(v string) (string, error) {
	if v == "1" || v == "on" {
		return "1", nil
	} else if v == "0" || v == "off" {
		return "0", nil
	} else {
		return "", fmt.Errorf("not an on off string")
	}
}

// extractPrefixCommentsAndRewrite extractPrefixComments and rewrite origin SQL
func extractPrefixCommentsAndRewrite(sql string, version *util.VersionCompareStatus) (trimmed string, comment parser.MarginComments) {
	sql = preRewriteSQL(sql, version)

	//TODO: 优化 tokens 逻辑，所有的 comments 都从 tokens 中获取
	_, comments := parser.SplitMarginComments(sql)
	trimmed = strings.TrimPrefix(sql, comments.Leading)
	trimmed = strings.TrimSuffix(trimmed, comments.Trailing)
	return trimmed, comments
}

// master-slave routing
func checkExecuteFromSlave(reqCtx *util.RequestContext, c *SessionExecutor, sql string) bool {
	stmtType := reqCtx.GetStmtType()
	tokens := reqCtx.GetTokens()
	tokensLen := len(tokens)

	if stmtType != parser.StmtSelect && stmtType != parser.StmtShow {
		return false
	}

	// if user is ReadOnly,then only can
	if !c.GetNamespace().IsAllowWrite(c.user) {
		return true
	}

	// send sql `select ... for update [nowait/skip locked]`
	// or `select ... in share mode [nowait/skip locked]` to master
	if c.GetNamespace().CheckSelectLock {
		if len(tokens) < 2 {
			return true
		}
		lastFirstWord := strings.ToLower(tokens[tokensLen-1])
		lastSecondWord := strings.ToLower(tokens[tokensLen-2])
		if (lastFirstWord == "update" && lastSecondWord == "for") ||
			(lastFirstWord == "mode" && lastSecondWord == "share") ||
			(lastFirstWord == "share" && lastSecondWord == "for") ||
			(lastFirstWord == "nowait" && (lastSecondWord == "share" || lastSecondWord == "update")) ||
			(lastFirstWord == "locked" && lastSecondWord == "skip") {
			return false
		}
	}

	// handle show variables like 'read_only' default to master
	if stmtType == parser.StmtShow && strings.Contains(strings.ToLower(sql), readonlyVariable) {
		return false
	}

	if strings.Contains(sql, "@@") {
		// handle select @@read_only default to master
		if strings.Contains(strings.ToLower(sql), "@@"+readonlyVariable) {
			return false
		}

		// handle select @@global.read_only default to master
		if strings.Contains(strings.ToLower(sql), "@@"+globalReadonlyVariable) {
			return false
		}
	}

	// handle master hint
	if len(tokens) > 1 && util.LowerEqual(tokens[1], masterHint) {
		return false
	}
	// handle master hint
	if len(tokens) > 1 && util.LowerEqual(tokens[tokensLen-1], masterHint) {
		return false
	}

	return c.GetNamespace().IsRWSplit(c.user)
}

// 如果是只读用户, 且SQL是INSERT, UPDATE, DELETE, 则拒绝执行, 返回true
func isSQLNotAllowedByUser(c *SessionExecutor, stmtType int) bool {
	if c.GetNamespace().IsAllowWrite(c.user) {
		return false
	}

	return stmtType == parser.StmtDelete || stmtType == parser.StmtInsert || stmtType == parser.StmtUpdate
}

// 旧版本，这边有个版本对比的函数性能比较差，qps 大时损耗比较严重遂去掉，Contains 比 HasSuffix 性能差，去掉
// preRewriteSQL pre rewite sql with string
func preRewriteSQL(sql string, version *util.VersionCompareStatus) string {
	if !version.LessThanMySQLVersion80 {
		return sql
	}
	// fix jdbc version mismatch gaea version
	if strings.HasPrefix(sql, JdbcInitPrefix) {
		return strings.Replace(sql, TxIsolationGT5720, TxIsolationLT5720, 1)
	}

	// fix `select @@transaction_isolation`
	if strings.HasSuffix(sql, TxIsolationGT5720) {
		return strings.Replace(sql, TxIsolationGT5720, TxIsolationLT5720, 1)
	}
	// fix `select @@session.transaction_isolation`
	if strings.HasSuffix(sql, SessionTxIsolationGT5720) {
		return strings.Replace(sql, SessionTxIsolationGT5720, SessionTxIsolationLT5720, 1)
	}

	// fix `select @@transaction_read_only`
	if strings.HasSuffix(sql, TxReadonlyGT5720) {
		return strings.Replace(sql, TxReadonlyGT5720, TxReadonlyLT5720, 1)
	}
	// fix `select @@session.transaction_read_only`
	if strings.HasSuffix(sql, SessionTxReadonlyGT5720) {
		return strings.Replace(sql, SessionTxReadonlyGT5720, SessionTxReadonlyLT5720, 1)
	}
	return sql
}

func modifyResultStatus(r *mysql.Result, cc *SessionExecutor) {
	if r == nil {
		return
	}
	r.Status = r.Status | cc.GetStatus()
}

func createShowDatabaseResult(dbs []string) *mysql.Result {
	r := new(mysql.Resultset)

	//
	field := &mysql.Field{
		Charset: 33,
		Type:    0xFD, //FIELD_TYPE_VAR_STRING fix: show databases jdbc err
	}
	field.Name = hack.Slice("Database")
	r.Fields = append(r.Fields, field)

	for _, db := range dbs {
		r.Values = append(r.Values, []interface{}{db})
	}

	result := mysql.ResultPool.Get()
	result.AffectedRows = uint64(len(dbs))
	result.Resultset = r

	plan.GenerateSelectResultRowData(result)
	return result
}

func createShowGeneralLogResult() *mysql.Result {
	r := new(mysql.Resultset)

	field := &mysql.Field{}
	field.Name = hack.Slice(gaeaGeneralLogVariable)
	r.Fields = append(r.Fields, field)

	var value string
	if OpenProcessGeneralQueryLog() {
		value = "ON"
	} else {
		value = "OFF"
	}
	r.Values = append(r.Values, []interface{}{value})
	result := mysql.ResultPool.Get()
	result.AffectedRows = 1
	result.Resultset = r

	plan.GenerateSelectResultRowData(result)
	return result
}

func getFromSlave(reqCtx *util.RequestContext) bool {
	slaveFlag := reqCtx.GetFromSlave()
	return slaveFlag == 1
}

// 仅多语句执行时使用
func setContextSQLFingerprint(reqCtx *util.RequestContext, sql string) {
	fingerprint := mysql.GetFingerprint(sql)
	md5sql := mysql.GetMd5(fingerprint)
	reqCtx.SetFingerprint(fingerprint)
	reqCtx.SetFingerprintMD5(md5sql)
}

func getSQLFingerprint(reqCtx *util.RequestContext, sql string) string {
	if reqCtx.GetFingerprint() == "" {
		fingerprint := mysql.GetFingerprint(sql)
		reqCtx.SetFingerprint(fingerprint)
	}
	return reqCtx.GetFingerprint()
}

func getSQLFingerprintMd5(reqCtx *util.RequestContext, sql string) string {
	if reqCtx.GetFingerprintMD5() == "" {
		fingerprint := getSQLFingerprint(reqCtx, sql)
		md5Value := mysql.GetMd5(fingerprint)
		reqCtx.SetFingerprintMD5(md5Value)
	}
	return reqCtx.GetFingerprintMD5()
}

func (se *SessionExecutor) isInTransaction() bool {
	return se.status&mysql.ServerStatusInTrans > 0 || !se.isAutoCommit()
}

func (se *SessionExecutor) isAutoCommit() bool {
	return se.status&mysql.ServerStatusAutocommit > 0
}

func (se *SessionExecutor) handleShow(reqCtx *util.RequestContext, sql string) (*mysql.Result, error) {
	tokens := reqCtx.GetTokens()
	if len(tokens) == 0 {
		return nil, fmt.Errorf("show command is empty")
	}

	// handle show databases;
	if len(tokens) == 2 && strings.ToLower(tokens[1]) == "databases" {
		dbs := se.GetNamespace().GetAllowedDBs()
		return createShowDatabaseResult(dbs), nil
	}
	// readonly && readwrite user send to slave
	if !se.GetNamespace().IsAllowWrite(se.user) || se.GetNamespace().IsRWSplit(se.user) {
		reqCtx.SetFromSlave(1)
	}
	// handle show variables like '%read_only%' default to master
	if strings.Contains(sql, readonlyVariable) && se.GetNamespace().IsAllowWrite(se.user) {
		reqCtx.SetFromSlave(0)
	}
	r, err := se.ExecuteSQL(reqCtx, se.GetNamespace().GetDefaultSlice(), se.db, sql)
	if err != nil {
		return nil, fmt.Errorf("execute sql error, sql: %s, err: %v", sql, err)
	}

	modifyResultStatus(r, se)
	return r, nil
}

func (se *SessionExecutor) handleKill(reqCtx *util.RequestContext, sql string) (*mysql.Result, error) {
	se.manager.statistics.generalLogger.Warn("%s - %dms - ns=%s, %s@%s->%s/%s, connect_id=%d, mysql_connect_id=%d, transaction=%t|%v. err:%s",
		SQLExecStatusIgnore, 0, se.namespace, se.user, se.clientAddr, "", se.db, se.session.c.GetConnectionID(), 0, se.isInTransaction(), sql, "ignore kill sql")
	return mysql.ResultPool.GetWithoutResultSet(), nil
}

func (se *SessionExecutor) handleBegin() error {
	se.txLock.Lock()
	defer se.txLock.Unlock()

	for _, co := range se.txConns {
		if err := co.Begin(); err != nil {
			return err
		}
	}

	// 客户端执行 begin 时后端 MySQL 实际并未执行
	for _, co := range se.ksConns {
		if err := co.Begin(); err != nil {
			return err
		}
	}

	se.status |= mysql.ServerStatusInTrans
	se.savepoints = []string{}
	return nil
}

func (se *SessionExecutor) handleCommit() (err error) {
	if err := se.commit(); err != nil {
		return err
	}
	return nil

}

// handleRollback handle rollback and rollback to savepoint
func (se *SessionExecutor) handleRollback(stmt *ast.RollbackStmt) (err error) {
	if stmt == nil || stmt.Savepoint == "" {
		return se.rollback()
	} else {
		return se.rollbackSavepoint(stmt.Savepoint)
	}
}

func (se *SessionExecutor) commit() (err error) {
	se.txLock.Lock()
	defer se.txLock.Unlock()

	se.status &= ^mysql.ServerStatusInTrans

	for _, pc := range se.txConns {
		if e := pc.Commit(); e != nil {
			err = e
		}
		pc.Recycle()

	}

	for _, pc := range se.ksConns {
		if e := pc.Commit(); e != nil {
			err = e
		}
	}
	se.txConns = make(map[string]backend.PooledConnect)
	se.savepoints = []string{}
	return
}

func (se *SessionExecutor) rollback() (err error) {
	se.txLock.Lock()
	defer se.txLock.Unlock()
	se.status &= ^mysql.ServerStatusInTrans
	for _, pc := range se.txConns {
		err = pc.Rollback()
		pc.Recycle()
	}

	for _, pc := range se.ksConns {
		err = pc.Rollback()
	}
	se.txConns = make(map[string]backend.PooledConnect)
	se.savepoints = []string{}
	return
}

func (se *SessionExecutor) rollbackSavepoint(savepoint string) (err error) {
	se.txLock.Lock()
	defer se.txLock.Unlock()
	for _, pc := range se.txConns {
		_, err = pc.Execute("rollback to "+savepoint, 0)
	}
	for _, pc := range se.ksConns {
		_, err = pc.Execute("rollback to "+savepoint, 0)
	}
	if err == nil && se.isInTransaction() {
		if index := util.ArrayFindIndex(se.savepoints, savepoint); index > -1 {
			se.savepoints = se.savepoints[0:index]
		}
	}
	return
}

// handleSavepoint handle savepoint and release savepoint
func (se *SessionExecutor) handleSavepoint(stmt *ast.SavepointStmt) (err error) {
	se.txLock.Lock()
	defer se.txLock.Unlock()
	if stmt.Release {
		for _, pc := range se.txConns {
			_, err = pc.Execute("release savepoint "+stmt.Savepoint, 0)
		}
		if err == nil && se.isInTransaction() {
			if index := util.ArrayFindIndex(se.savepoints, stmt.Savepoint); index > -1 {
				se.savepoints = se.savepoints[0 : index+1]
			}
		}
	} else {
		for _, pc := range se.txConns {
			_, err = pc.Execute("savepoint "+stmt.Savepoint, 0)
		}
		if err == nil && se.isInTransaction() {
			if util.ArrayFindIndex(se.savepoints, stmt.Savepoint) > -1 {
				se.savepoints = util.ArrayRemoveItem(se.savepoints, stmt.Savepoint)
			}
			se.savepoints = append(se.savepoints, stmt.Savepoint)
		}
	}
	return
}

func (se *SessionExecutor) recycleTx() {
	if !se.isInTransaction() {
		return
	}
	se.txLock.Lock()
	defer se.txLock.Unlock()
	se.txConns = make(map[string]backend.PooledConnect)
}

// handleKQuit close backend connection and recycle, only called when client exit
func (se *SessionExecutor) handleKsQuit() {
	for _, ksConn := range se.ksConns {
		ksConn.Close()
		ksConn.Recycle()
	}
	se.ksConns = make(map[string]backend.PooledConnect)
}

// ExecuteSQL execute sql
func (se *SessionExecutor) ExecuteSQL(reqCtx *util.RequestContext, slice, db, sql string) (*mysql.Result, error) {
	phyDB, err := se.GetNamespace().GetDefaultPhyDB(db)
	if err != nil {
		return nil, err
	}

	pc, err := se.getBackendConn(slice, getFromSlave(reqCtx))
	defer se.recycleBackendConn(pc)

	if err != nil {
		log.Warn("[ns:%s]getBackendConn failed: %v", se.GetNamespace().name, err)
		return nil, fmt.Errorf("getBackendConn failed: %v", err)
	}

	se.backendAddr = pc.GetAddr()
	se.backendConnectionId = pc.GetConnectionID()

	rs, err := se.executeUnshardSQLInSlice(reqCtx, pc, phyDB, sql)
	if err != nil {
		return nil, err
	}

	if pc.MoreRowsExist() || pc.MoreResultsExist() {
		se.session.continueConn = pc
	}
	return rs, nil
}

// ExecuteSQLs len(sqls) must not be 0, or return error
func (se *SessionExecutor) ExecuteSQLs(reqCtx *util.RequestContext, sqls map[string]map[string][]string) ([]*mysql.Result, error) {
	if len(sqls) == 0 {
		return nil, fmt.Errorf("no sql to execute")
	}

	pcs, err := se.getBackendConns(sqls, getFromSlave(reqCtx))
	defer se.recycleBackendConns(pcs, false)
	if err != nil {
		log.Warn("getShardConns failed: %v", err)
		return nil, err
	}

	rs, err := se.executeShardSQLInSlice(reqCtx, pcs, sqls)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
