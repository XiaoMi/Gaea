/*
 * @Author: funnyAnt
 * @Date: 2022-05-16 11:52:01
 * @LastEditTime: 2022-05-17 10:56:29
 * @LastEditors:
 * @Description: 分布式死锁检查
 */
package server

import (
	"strconv"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/log"
	"github.com/heimdalr/dag"
)
 
 func doDeadlockCheck(s *Server, slices map[string]*backend.Slice) error {
	 //init dag
	 d := dag.NewDAG()
 
	 addDiVertex := func(value string) error {
		 if v, _ := d.GetVertex(value); v == nil {
			 if err := d.AddVertexByID(value, value); err != nil {
				 return err
			 }
		 }
 
		 return nil
	 }
 
	 addDiEdge := func(from string, to string) (loopId string, err error) {
		 if err = addDiVertex(from); err != nil {
			 return "", err
		 }
 
		 if err = addDiVertex(to); err != nil {
			 return "", err
		 }
 
		 log.Debug("add dig[from:%v,to:%v]", from, to)
		 if err = d.AddEdge(from, to); err != nil {
			 if _, ok := err.(dag.EdgeLoopError); ok {
				 loopId = from
				 return loopId, err
			 }
		 }
 
		 return "", nil
 
	 }
 
	 backend2FrontendMap := fetchBackend2FrontendConn(s)
	 for k, v := range slices {
		 toKill, err := fetchLockWaits(v, func(waiting, blocking string) (string, error) {
			 keyFrom := k + backend2FrontendDelimiter + waiting
			 if v, ok := backend2FrontendMap[keyFrom]; ok {
				 frontendIdFrom := v.c.ConnectionID
				 keyTo := k + backend2FrontendDelimiter + blocking
				 if v, ok := backend2FrontendMap[keyTo]; ok {
					 frontendIdTo := v.c.ConnectionID
					 loopId, err := addDiEdge(strconv.Itoa(int(frontendIdFrom)), strconv.Itoa(int(frontendIdTo)))
					 if loopId != "" {
						 return keyFrom, err
					 }
				 }
			 }
			 return "", nil
		 })
 
		 log.Debug("add dig info[%v]", d.String())
		 if toKill != "" {
			 //kill the frontend
			 killFrontendConn(backend2FrontendMap[toKill])
			 continue
		 }
		 if err != nil {
			 return err
		 }
 
	 }
 
	 return nil
 }

 func fetchLockWaits(slice *backend.Slice, callback func(waiting, blocking string) (string, error)) (string, error) {
	 const MAX_ROWS = 10000
	 sql := "SELECT distinct trx_a.trx_mysql_thread_id AS waiting, trx_b.trx_mysql_thread_id AS blocking " +
		 "FROM information_schema.innodb_lock_waits, information_schema.INNODB_TRX AS trx_a, information_schema.INNODB_TRX AS trx_b " +
		 "WHERE trx_a.trx_id = requesting_trx_id AND trx_b.trx_id = blocking_trx_id;"
 
	 //doQuery
	 conn, err := slice.GetMasterConn()
	 if err != nil {
		 return "", err
	 }
	 result, err := conn.Execute(sql, MAX_ROWS)
	 if err != nil {
		 return "", err
	 }
	 for _, rowValue := range result.Values {
		 waiting := strconv.Itoa(int(rowValue[0].(uint64)))
		 blocking := strconv.Itoa(int(rowValue[1].(uint64)))
		 if toKill, err := callback(waiting, blocking); err != nil {
			 return toKill, err
		 }
	 }
	 return "", nil
 }
 
 func fetchBackend2FrontendConn(s *Server) map[string]*Session {
	 return s.getBackend2FrontendMap()
 }
 
 func killFrontendConn(session *Session) {
	 if session != nil {
		 session.forceClose("deadlock check")
	 }
 }
 