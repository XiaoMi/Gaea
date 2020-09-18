package server

import "sync/atomic"

// Process global variables.
var (
	ProcessGeneralLog uint32
)

func OpenProcessGeneralQueryLog() bool {
	return atomic.LoadUint32(&ProcessGeneralLog) == 1
}
