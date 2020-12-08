package backend

import (
	"errors"
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
)

type statManager struct {
	master         replTracker
	slave          []replTracker
	statisticSlave []replTracker

	unhealthyThresholdSeconds int
}

// replTrackerImpl tracks replication lag.
type replTrackerImpl struct {
	open      bool
	mode      string
	isMaster  bool
	unhealthy int

	cp     ConnectionPool
	hw     *heartbeatWriter
	hr     *heartbeatReader
	poller *poller
}

type replTracker interface {
	init()
	status() bool
}

func newStatManager(config *models.HealthCheckConfig, slice *Slice) *statManager {
	var hc statManager

	hc.master = newReplTracker(config, true, slice.Master)
	hc.master.init()

	hc.slave = make([]replTracker, len(slice.Slave))
	hc.statisticSlave = make([]replTracker, len(slice.StatisticSlave))
	for i, slave := range slice.Slave {
		hc.slave[i] = newReplTracker(config, false, slave)
		hc.slave[i].init()
	}
	for i, slave := range slice.StatisticSlave {
		hc.slave[i] = newReplTracker(config, false, slave)
		hc.slave[i].init()
	}
	return &hc
}

func (c *statManager) getMaster() bool {
	return c.master.status()
}

func (c *statManager) getSlave(index int) bool {
	return c.slave[index].status()
}

func (c *statManager) getStatisticSlave(index int) bool {
	return c.statisticSlave[index].status()
}

// newReplTracker creates a new replTrackerImpl.
func newReplTracker(config *models.HealthCheckConfig, isMaster bool, pool ConnectionPool) replTracker {
	return &replTrackerImpl{
		open:      config.Open,
		mode:      config.Mode,
		cp:        pool,
		isMaster:  isMaster,
		unhealthy: config.UnhealthyThresholdSeconds,
	}
}

// InitDBConfig initializes the target name.
func (rt *replTrackerImpl) init() {
	if !rt.open {
		return
	}

	switch rt.mode {
	case models.HeartBeatMode:
		rt.hr = newHeartbeatReader(rt.cp)
		if rt.isMaster {
			rt.hw = newHeartbeatWriter(rt.cp)
		}
	case models.PollerMode:
		if !rt.isMaster {
			rt.poller = newPoller(rt.cp)
		}
	}
	return
}

func (rt *replTrackerImpl) status() bool {
	lag, err := rt.fetchStatus()
	if err != nil {
		log.Warn("get instance status error, instance:%s err:%s", rt.cp.Addr(), err)
		return false
	}
	if lag > time.Second*time.Duration(rt.unhealthy) {
		log.Warn("going unhealthy due to high replication lag, instance:%s lag:%v", rt.cp.Addr(), lag)
		return false
	}
	return true
}

// Status reports the replication status.
func (rt *replTrackerImpl) fetchStatus() (time.Duration, error) {
	switch {
	case !rt.open:
		return 0, nil
	case rt.mode == models.HeartBeatMode:
		if rt.isMaster {
			if err := rt.hw.run(); err != nil {
				log.Warn("execute heartbeat sql error, master addr:%s, error:%s", rt.cp.Addr(), err)
			}
		}
		return rt.hr.run()
	case rt.mode == models.PollerMode:
		if rt.isMaster {
			return 0, nil
		}
		return rt.poller.run()
	}
	return 0, errors.New(fmt.Sprintf("invalid health check mode:%s", rt.mode))
}
