package util

import (
	"os/exec"
)

// StartGaeaDefault gaea的一些操作
func StartGaeaDefault() error {
	cmd := exec.Command("sh", "-c", CmdStartGaea)
	_, err := cmd.CombinedOutput()
	return err
}

// StartGaeaCCDefault gaeacc的一些操作
func StartGaeaCCDefault() error {
	cmd := exec.Command("sh", "-c", CmdStartGaeaCC)
	_, err := cmd.CombinedOutput()
	return err
}

// StopGaeaDefault gaea的一些操作
func StopGaeaDefault() error {
	cmd := exec.Command("sh", "-c", CmdStopGaea)
	_, err := cmd.CombinedOutput()
	return err
}

// StopGaeaCCDefault gaeacc的一些操作
func StopGaeaCCDefault() error {
	cmd := exec.Command("sh", "-c", CmdStopGaeaCC)
	_, err := cmd.CombinedOutput()
	return err
}
