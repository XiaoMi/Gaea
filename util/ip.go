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
	"errors"
	"fmt"
	"net"
	"strings"
)

// IPInfo ip information
type IPInfo struct {
	info    string
	isIPNet bool
	ip      net.IP
	ipNet   net.IPNet
}

// ParseIPInfo parse ip
func ParseIPInfo(v string) (IPInfo, error) {
	if ip, ipNet, err := net.ParseCIDR(v); err == nil {
		return IPInfo{
			info:    v,
			isIPNet: true,
			ip:      ip,
			ipNet:   *ipNet,
		}, nil
	}

	if ip := net.ParseIP(v); ip != nil {
		return IPInfo{
			info:    v,
			isIPNet: false,
			ip:      ip,
		}, nil
	}

	return IPInfo{}, errors.New("invalid ip address")
}

// Info return information of ip
func (t *IPInfo) Info() string {
	return t.info
}

// Match check if ip matched
func (t *IPInfo) Match(ip net.IP) bool {
	if t.isIPNet {
		return t.ipNet.Contains(ip)
	}
	return t.ip.Equal(ip)
}

func parseAllowIps(allowIpsStr string) ([]IPInfo, error) {
	if len(allowIpsStr) == 0 {
		return make([]IPInfo, 0, 10), nil
	}
	ipVec := strings.Split(allowIpsStr, ",")
	allowIpsList := make([]IPInfo, 0, 10)
	for _, ipStr := range ipVec {
		if ip, err := ParseIPInfo(strings.TrimSpace(ipStr)); err == nil {
			allowIpsList = append(allowIpsList, ip)
		}
	}
	return allowIpsList, nil
}

func resolveAddr(network string, address string, replaceZeroAddr bool) (string, error) {
	switch network {
	default:
		return "", fmt.Errorf("invalid network '%s'", network)

	case "unix", "unixpacket":
		return address, nil

	case "tcp", "tcp4", "tcp6":
		addr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return "", err
		}
		if addr.Port != 0 {
			//是否为空地址
			if !isZeroIPAddr(addr) {
				return addr.String(), nil
			}
			if replaceZeroAddr {
				ifaddrs, err := net.InterfaceAddrs()
				if err != nil {
					return "", err
				}
				for _, ifaddr := range ifaddrs {
					switch in := ifaddr.(type) {
					case *net.IPNet:
						if in.IP.IsGlobalUnicast() {
							addr.IP = in.IP
							return addr.String(), nil
						}
					}
				}
			}
		}
		return "", fmt.Errorf("invalid address '%s'", addr.String())
	}
}

// ResolveAddr return real ip by net interface
func ResolveAddr(network string, locAddr string) (string, error) {
	return resolveAddr(network, locAddr, true)
}

func isZeroIPAddr(addr *net.TCPAddr) bool {
	if ipv4 := addr.IP.To4(); ipv4 != nil {
		return net.IPv4zero.Equal(ipv4)
	} else if ipv6 := addr.IP.To16(); ipv6 != nil {
		return net.IPv6zero.Equal(ipv6)
	}
	return false
}

// HostName return hostname by ip
func HostName(ip string) (hostname string, err error) {
	hostName, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}
	return hostName[0], err
}
