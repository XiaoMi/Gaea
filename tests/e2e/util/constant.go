package util

var (
	// gaea 相关命令
	CmdStartGaea   = "cd ./cmd ; ./gaea --config=gaea.ini >/dev/null 2>&1 &"
	CmdStopGaea    = "pkill gaea"
	CmdStartGaeaCC = "cd ./cmd ; ./gaea-cc -c=gaea_cc.ini >/dev/null 2>&1 &"
	CmdStopGaeaCC  = "pkill gaea-cc"
)

type ResultType int

const (
	Equal = iota
	UnEqual
	UnSupport
)
