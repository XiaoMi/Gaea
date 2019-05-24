// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package util

import (
	"fmt"
	"net"
	"testing"
)

func TestParseIP(t *testing.T) {
	ip := net.ParseIP("abcdefg")
	fmt.Println(ip)
	fmt.Println(ip.String())
}

func TestCreateIPInfoIPSuccess(t *testing.T) {
	addr := "127.0.0.1"
	info, err := ParseIPInfo(addr)
	if err != nil {
		t.FailNow()
	}
	if info.isIPNet {
		t.FailNow()
	}
	if addr != info.info {
		t.FailNow()
	}
	if addr != info.ip.String() {
		t.FailNow()
	}
}

func TestCreateIPInfoIPError(t *testing.T) {
	addr := "127.255.256.1"
	if _, err := ParseIPInfo(addr); err == nil {
		t.FailNow()
	}
}

func TestCreateIPInfoIPError2(t *testing.T) {
	addr := "abcdefg"
	if _, err := ParseIPInfo(addr); err == nil {
		t.FailNow()
	}
}

func TestCreateIPInfoIPNetSuccess(t *testing.T) {
	addr := "192.168.122.1/24"
	netAddr := "192.168.122.0/24"
	info, err := ParseIPInfo(addr)
	if err != nil {
		t.FailNow()
	}
	if !info.isIPNet {
		t.FailNow()
	}
	if addr != info.info {
		t.FailNow()
	}
	if netAddr != info.ipNet.String() {
		t.FailNow()
	}
}

func TestCreateIPInfoIPNetError(t *testing.T) {
	addr := "192.168.122.1/"
	if _, err := ParseIPInfo(addr); err == nil {
		t.FailNow()
	}
}

func TestCreateIPInfoIPNetError2(t *testing.T) {
	addr := "192.168.122.1/35"
	if _, err := ParseIPInfo(addr); err == nil {
		t.FailNow()
	}
}
