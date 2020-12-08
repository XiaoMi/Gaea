package models

import (
	"errors"
)

const (
	HeartBeatMode = "heartbeat"
	PollerMode    = "poller"
)

// HealthCheckConfig contains the config for healthCheck.
type HealthCheckConfig struct {
	Open                      bool   `json:"open"`
	Mode                      string `json:"mode"`
	IntervalSeconds           int    `json:"interval"`
	UnhealthyThresholdSeconds int    `json:"unhealthy_threshold"`
}

func (s *HealthCheckConfig) verify() error {
	if s.Open == false {
		return nil
	}

	if s.Mode != HeartBeatMode && s.Mode != PollerMode {
		return errors.New("health check mode should be heartbeat or poller")
	}
	if s.IntervalSeconds <= 0 {
		return errors.New("health check internal seconds should be > 0")
	}
	if s.UnhealthyThresholdSeconds <= 0 {
		return errors.New("unhealthy threshold seconds should be > 0")
	}
	return nil
}
