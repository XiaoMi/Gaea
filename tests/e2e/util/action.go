// Copyright 2024 The Gaea Authors. All Rights Reserved.
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
