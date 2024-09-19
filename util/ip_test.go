// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/stretchr/testify/assert"
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

func TestGetInstanceDatacenter(t *testing.T) {
	testCases := []struct {
		name     string
		addr     string
		hasErr   bool
		expectDc string
	}{
		{
			"test get datacenter success",
			"c3-mysql01.bj:3306",
			false,
			"c3",
		},
		{
			"test get datacenter format error too much hyphens",
			"c3-mysql01.bj-3306",
			true,
			"",
		},
		{
			"test get datacenter format error too much dots",
			"c3.mysql01.bj.3306",
			true,
			"",
		},
		{
			"test get datacenter format error no enough hyphens",
			"c3-mysql01.bj-3306",
			true,
			"",
		},
		{
			"test get datacenter format error no enough dots",
			"c3.mysql01-bj:3306",
			false,
			"c3.mysql01",
		},
		{
			"test get datacenter format error no colons",
			"c3.mysql01.bj",
			true,
			"",
		},
	}

	for k, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			dc, err := GetInstanceDatacenter(tt.addr)
			t.Logf("Testing Case:%d addr: %s, dc: %v, err: %v", k, tt.addr, dc, err)
			if tt.hasErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expectDc, dc)
		})
	}
}

func TestGetHostDatacenter(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		hasErr   bool
		expectDc string
	}{
		{
			"test get datacenter success",
			"c3-mysql01.bj",
			false,
			"c3",
		},
		{
			"test get datacenter success local",
			"MacBook-Pro-2.local",
			false,
			"MacBook",
		},
		{
			"test get datacenter format error too much dots",
			"c3.mysql01.bj",
			true,
			"",
		},
		{
			"test get datacenter format error no enough dots",
			"c3.mysql01-bj",
			false,
			"c3.mysql01",
		},
		{
			"test get datacenter format error no colons",
			"c3.mysql01.bj",
			true,
			"",
		},
	}

	for k, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			dc, err := GetHostDatacenter(tt.host)
			t.Logf("Testing Case:%d addr: %s, dc: %v, err: %v", k, tt.host, dc, err)
			if tt.hasErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expectDc, dc)
		})
	}
}
