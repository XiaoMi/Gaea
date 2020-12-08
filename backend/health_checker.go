package backend

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	gaeaServerId                 = 3
	sqlUpsertMostRecentHeartbeat = "INSERT INTO gaea_dba_check.heartbeat (id , check_timestamp) VALUES (%d, %d) ON DUPLICATE KEY UPDATE check_timestamp=VALUES(check_timestamp)"
	sqlFetchMostRecentHeartbeat  = "SELECT check_timestamp FROM gaea_dba_check.heartbeat WHERE id=%d"

	sqlFetchMostRecentReplicationStatus = "SHOW SLAVE STATUS"
)

// heartbeatReader reads the heartbeatReader table at a configured interval in order
// to calculate replication lag. It is meant to be run on a replica, and paired
// with a heartbeatWriter on a master.
// Lag is calculated by comparing the most recent timestamp in the heartbeatReader
// table against the current time at read time.
type heartbeatReader struct {
	cp ConnectionPool
}

// InitDBConfig initializes the target name for the heartbeatReader.
func newHeartbeatReader(pool ConnectionPool) *heartbeatReader {
	return &heartbeatReader{cp: pool}
}

// status returns the most recently recorded lag measurement or error encountered.
// fetchMostRecentHeartbeat fetches the most recently recorded heartbeatReader from the heartbeatReader table,
// returning a result with the timestamp of the heartbeatReader.
func (r *heartbeatReader) run() (time.Duration, error) {
	conn, err := r.cp.Get(context.Background())
	if err != nil {
		return 0, err
	}
	defer conn.Recycle()

	qr, err := conn.Execute(fmt.Sprintf(sqlFetchMostRecentHeartbeat, gaeaServerId))
	if err != nil {
		return 0, err
	}

	if qr.Resultset == nil {
		return 0, fmt.Errorf("failed to read heartbeat: result set is nil")
	} else if len(qr.Resultset.Values) != 1 {
		return 0, fmt.Errorf("failed to read heartbeat: writer query did not result in 1 row. Got %d",
			len(qr.Resultset.Values))
	}

	ts, err := qr.Resultset.GetIntByName(0, "check_timestamp")
	if err != nil {
		return 0, fmt.Errorf("failed to read heartbeat check_timestamp err:%v", err)
	}
	lag := time.Now().Sub(time.Unix(ts, 0))
	return lag, nil
}

// heartbeatWriter runs on master tablets and writes heartbeats to the _vt.heartbeat
// table at a regular interval, defined by heartbeat_interval.
type heartbeatWriter struct {
	cp ConnectionPool
}

// InitDBConfig initializes the target name for the heartbeatWriter.
func newHeartbeatWriter(pool ConnectionPool) *heartbeatWriter {
	return &heartbeatWriter{cp: pool}
}

func (w *heartbeatWriter) run() error {
	conn, err := w.cp.Get(context.Background())
	if err != nil {
		return err
	}
	defer conn.Recycle()

	if _, err = conn.Execute(fmt.Sprintf(sqlUpsertMostRecentHeartbeat, gaeaServerId, time.Now().Unix())); err != nil {
		return err
	}
	return nil
}

type poller struct {
	cp ConnectionPool

	mu             sync.Mutex
	lag            time.Duration
	timeRecorded   time.Time
	lastKnownError error
}

// ReplicationStatus holds replication information from SHOW SLAVE STATUS.
type replicationStatus struct {
	masterServerID      uint
	ioThreadRunning     bool
	sqlThreadRunning    bool
	secondsBehindMaster uint
	masterHost          string
	masterPort          int
	masterConnectRetry  int
}

func newPoller(pool ConnectionPool) *poller {
	return &poller{cp: pool}
}

func (p *poller) run() (time.Duration, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	status, err := p.fetchReplicationStatus()
	if err != nil {
		return 0, fmt.Errorf("fetch replication status error, err:%s", err)
	}

	if !(status.ioThreadRunning && status.sqlThreadRunning) {
		if p.timeRecorded.IsZero() {
			return 0, fmt.Errorf("replication is not running")
		}
		return time.Since(p.timeRecorded) + p.lag, nil
	}

	p.lag = time.Duration(status.secondsBehindMaster) * time.Second
	p.timeRecorded = time.Now()
	return p.lag, nil
}

func (p *poller) fetchReplicationStatus() (replicationStatus, error) {
	var status = replicationStatus{}
	conn, err := p.cp.Get(context.Background())
	if err != nil {
		return status, err
	}
	defer conn.Recycle()

	qr, err := conn.Execute(sqlFetchMostRecentReplicationStatus)
	if err != nil {
		return status, err
	} else if qr.Resultset == nil || len(qr.Resultset.Values) == 0 {
		// The query returned no data, meaning the server
		// is not configured as a replica.
		return status, errors.New("no replication status")
	}
	ioThreadRunning, _ := qr.Resultset.GetStringByName(0, "Slave_IO_Running")
	sqlThreadRunning, _ := qr.Resultset.GetStringByName(0, "Slave_SQL_Running")
	secondsBehindMaster, _ := qr.Resultset.GetStringByName(0, "Seconds_Behind_Master")
	parseUint, _ := strconv.ParseUint(secondsBehindMaster, 10, 0)

	// These fields are returned from the underlying DB and cannot be renamed
	status.ioThreadRunning = ioThreadRunning == "Yes" || ioThreadRunning == "Connecting"
	status.sqlThreadRunning = sqlThreadRunning == "Yes"
	status.secondsBehindMaster = uint(parseUint)
	return status, nil
}
