// Copyright 2019 The Gaea Authors. All Rights Reserved.
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

package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/crypto"
)

// Namespace means namespace model stored in etcd
type Namespace struct {
	IsEncrypt        bool              `json:"is_encrypt"` // true: 加密存储 false: 非加密存储，目前加密Slice、User中的用户名、密码
	Name             string            `json:"name"`
	Online           bool              `json:"online"`
	ReadOnly         bool              `json:"read_only"`
	AllowedDBS       map[string]bool   `json:"allowed_dbs"`
	DefaultPhyDBS    map[string]string `json:"default_phy_dbs"`
	SlowSQLTime      string            `json:"slow_sql_time"`
	BlackSQL         []string          `json:"black_sql"`
	AllowedIP        []string          `json:"allowed_ip"`
	Slices           []*Slice          `json:"slices"`
	ShardRules       []*Shard          `json:"shard_rules"`
	Users            []*User           `json:"users"` // 客户端接入proxy用户，每个用户可以设置读写分离、读写权限等
	DefaultSlice     string            `json:"default_slice"`
	GlobalSequences  []*GlobalSequence `json:"global_sequences"`
	DefaultCharset   string            `json:"default_charset"`
	DefaultCollation string            `json:"default_collation"`
}

// Encode encode json
func (n *Namespace) Encode() []byte {
	return JSONEncode(n)
}

// Verify verify namespace contents
func (n *Namespace) Verify() error {
	if err := n.verifyName(); err != nil {
		return err
	}

	if err := n.verifyAllowDBS(); err != nil {
		return err
	}

	if len(n.Users) == 0 {
		return errors.New("must specify proxy access users")
	}

	if verifySlowSQLTime(n.SlowSQLTime) != nil {
		return errors.New("invalid slow sql time")
	}

	if err := verifyDefaultPhyDB(n.DefaultPhyDBS, n.AllowedDBS); err != nil {
		return fmt.Errorf("verify defaultPhyDBs error: %v", err)
	}

	if err := verifyAllowIps(n.AllowedIP); err != nil {
		return fmt.Errorf("verify allowips error: %v", err)
	}

	if err := mysql.VerifyCharset(n.DefaultCharset, n.DefaultCollation); err != nil {
		return err
	}

	for i, u := range n.Users {
		if u.Namespace == "" {
			u.Namespace = n.Name
		} else if u.Namespace != n.Name {
			return fmt.Errorf("user's namespace name mismatch, user: %s, namespace: %s, %s", u.UserName, n.Name, u.Namespace)
		}

		if err := u.verify(); err != nil {
			return fmt.Errorf("user config error, schema: %s, %v", n.Name, err)
		}

		for j := 0; j < i; j++ {
			if n.Users[j].UserName == u.UserName {
				return fmt.Errorf("user duped, namespace: %s, user: %s", n.Name, u.UserName)
			}
		}
	}

	if len(n.Slices) == 0 {
		return errors.New("empty slices")
	}

	var sliceNames []string
	for i, slice := range n.Slices {
		sliceNames = append(sliceNames, slice.Name)
		if err := slice.verify(); err != nil {
			return fmt.Errorf("slice cfg error, namespace: %s, err: %s", n.Name, err.Error())
		}

		for j := 0; j < i; j++ {
			if n.Slices[j].Name == slice.Name {
				return fmt.Errorf("slice name duped, namespace: %s, slice: %s", n.Name, slice.Name)
			}
		}
	}

	if n.DefaultSlice != "" {
		exist := false

		for _, slice := range n.Slices {
			if slice.Name == n.DefaultSlice {
				exist = true
				break
			}
		}

		if !exist {
			return fmt.Errorf("invalid default slice: %s", n.DefaultSlice)
		}
	}

	rules := make(map[string]map[string]string)
	linkedRuleShards := []*Shard{}
	for _, s := range n.ShardRules {

		for _, slice := range s.Slices {
			if !includeSlice(sliceNames, slice) {
				return fmt.Errorf("shard table[%s] slice[%s] not in the namespace.slices list:[%s]",
					s.Table, slice, strings.Join(s.Slices, ","))
			}
		}
		if s.Type == ShardLinked {
			linkedRuleShards = append(linkedRuleShards, s)
		}
		if err := s.verify(); err != nil {
			return err
		}

		//if the database exist in rules
		if _, ok := rules[s.DB]; ok {
			if _, ok := rules[s.DB][s.Table]; ok {
				return fmt.Errorf("table %s rule in %s duplicate", s.Table, s.DB)
			} else {
				rules[s.DB][s.Table] = s.Type
			}
		} else {
			m := make(map[string]string)
			rules[s.DB] = m
			rules[s.DB][s.Table] = s.Type
		}
	}
	for _, s := range linkedRuleShards {
		tableRules, ok := rules[s.DB]
		if !ok {
			return fmt.Errorf("db of LinkedRule is not found in parent rules")
		}
		dbRuleType, ok := tableRules[s.ParentTable]
		if !ok {
			return fmt.Errorf("parent table of LinkedRule is not found in parent rules")
		}
		if dbRuleType == ShardLinked {
			return fmt.Errorf("LinkedRule cannot link to another LinkedRule")
		}
	}
	return nil
}

func (n *Namespace) verifyName() error {
	if n.isNameExists() {
		return fmt.Errorf("must specify namespace name")
	}
	return nil
}

func (n *Namespace) isNameExists() bool {
	return n.Name != ""
}

func (n *Namespace) verifyAllowDBS() error {
	if n.isAllowedDBSEmpty() {
		return errors.New("must specify usable dbs")
	}
	return nil
}

func (n *Namespace) isAllowedDBSEmpty() bool {
	return len(n.AllowedDBS) == 0
}

// Decrypt decrypt user/password in namespace
func (n *Namespace) Decrypt(key string) (err error) {
	if !n.IsEncrypt {
		return nil
	}
	// Users
	for i := range n.Users {
		n.Users[i].UserName, err = decrypt(key, n.Users[i].UserName)
		if err != nil {
			return
		}
		n.Users[i].Password, err = decrypt(key, n.Users[i].Password)
		if err != nil {
			return
		}
	}
	// Slices
	for i := range n.Slices {
		n.Slices[i].UserName, err = decrypt(key, n.Slices[i].UserName)
		if err != nil {
			return
		}
		n.Slices[i].Password, err = decrypt(key, n.Slices[i].Password)
		if err != nil {
			return
		}
	}

	return nil
}

// Encrypt encrypt user/password in namespace
func (n *Namespace) Encrypt(key string) (err error) {
	n.IsEncrypt = true
	// Users
	for i := range n.Users {
		n.Users[i].UserName, err = encrypt(key, n.Users[i].UserName)
		if err != nil {
			return
		}
		n.Users[i].Password, err = encrypt(key, n.Users[i].Password)
		if err != nil {
			return
		}
	}
	// Slices
	for i := range n.Slices {
		n.Slices[i].UserName, err = encrypt(key, n.Slices[i].UserName)
		if err != nil {
			return
		}
		n.Slices[i].Password, err = encrypt(key, n.Slices[i].Password)
		if err != nil {
			return
		}
	}

	return nil
}

func verifySlowSQLTime(slowSQLTimeStr string) error {
	if slowSQLTimeStr == "" {
		return nil
	}

	slowSQLTime, err := strconv.ParseInt(slowSQLTimeStr, 10, 64)
	if err != nil || slowSQLTime < 0 {
		return errors.New("invalid slow sql time")
	}

	return nil
}

func decrypt(key, data string) (string, error) {
	t, _ := base64.StdEncoding.DecodeString(data)
	origin, err := crypto.DecryptECB(key, t)
	if err != nil {
		return "", err
	}
	return string(origin), nil
}

func encrypt(key, data string) (string, error) {
	tmp, err := crypto.EncryptECB(key, []byte(data))
	if err != nil {
		return "", err
	}
	base64Str := base64.StdEncoding.EncodeToString(tmp)
	return base64Str, nil
}

func verifyDefaultPhyDB(defaultPhyDBs map[string]string, allowedDBs map[string]bool) error {
	// no logic database mode
	if len(defaultPhyDBs) == 0 {
		return nil
	}

	// logic database mode
	for db := range allowedDBs {
		if _, ok := defaultPhyDBs[db]; !ok {
			return fmt.Errorf("db %s have no phy db", db)
		}
	}
	return nil
}

func verifyAllowIps(allowedIP []string) error {
	for _, ipStr := range allowedIP {
		ipStr = strings.TrimSpace(ipStr)
		if len(ipStr) == 0 {
			continue
		}
		_, err := util.ParseIPInfo(ipStr)
		if err != nil {
			return err
		}
	}
	return nil
}
