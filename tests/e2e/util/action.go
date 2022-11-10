package util

import (
	"fmt"
	"os/exec"
)

// StartGaeaDefault gaea的一些操作
func StartGaeaDefault() error {
	cmd := exec.Command("sh", "-c", CmdStartGaea)
	out, err := cmd.CombinedOutput()
	fmt.Println("start gaea out:", string(out))
	return err
}

// StopGaeaDefault gaea的一些操作
func StopGaeaDefault() error {
	cmd := exec.Command("sh", "-c", CmdStopGaea)
	_, err := cmd.CombinedOutput()
	return err
}
