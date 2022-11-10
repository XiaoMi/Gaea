package util

var (
	// etcd 相关命令

	// gaea 相关命令
	CmdStartGaea = "cd cmd ; ./gaea --config=gaea.ini >/dev/null 2>&1 &"
	CmdStopGaea  = "pkill gaea"
)
