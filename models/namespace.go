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

	"github.com/XiaoMi/Gaea/mysql"
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
func (p *Namespace) Encode() []byte {
	return JSONEncode(p)
}

// Verify verify namespace contents
func (p *Namespace) Verify() error {
	if p.Name == "" {
		return errors.New("must specify namespace name")
	}

	if len(p.AllowedDBS) == 0 {
		return errors.New("must specify usable dbs")
	}

	if len(p.Users) == 0 {
		return errors.New("must specify proxy access users")
	}

	if verifySlowSQLTime(p.SlowSQLTime) != nil {
		return errors.New("invalid slow sql time")
	}

	if err := mysql.VerifyCharset(p.DefaultCharset, p.DefaultCollation); err != nil {
		return err
	}

	for i, u := range p.Users {
		if u.Namespace == "" {
			u.Namespace = p.Name
		} else if u.Namespace != p.Name {
			return fmt.Errorf("user's namespace name mismatch, user: %s, namespace: %s, %s", u.UserName, p.Name, u.Namespace)
		}

		if err := u.verify(); err != nil {
			return fmt.Errorf("user config error, schema: %s, %v", p.Name, err)
		}

		for j := 0; j < i; j++ {
			if p.Users[j].UserName == u.UserName {
				return fmt.Errorf("user duped, namespace: %s, user: %s", p.Name, u.UserName)
			}
		}
	}

	if len(p.Slices) == 0 {
		return errors.New("empty slices")
	}

	for i, n := range p.Slices {
		if err := n.verify(); err != nil {
			return fmt.Errorf("slice cfg error, namespace: %s, err: %s", p.Name, err.Error())
		}

		for j := 0; j < i; j++ {
			if p.Slices[j].Name == n.Name {
				return fmt.Errorf("slice name duped, namespace: %s, slice: %s", p.Name, n.Name)
			}
		}
	}

	if p.DefaultSlice != "" {
		exist := false

		for _, n := range p.Slices {
			if n.Name == p.DefaultSlice {
				exist = true
				break
			}
		}

		if !exist {
			return fmt.Errorf("invalid default slice: %s", p.DefaultSlice)
		}
	}

	for _, s := range p.ShardRules {
		if err := s.verify(); err != nil {
			return err
		}
	}

	return nil
}

// Decrypt decrypt user/password in namespace
func (p *Namespace) Decrypt(key string) (err error) {
	if !p.IsEncrypt {
		return nil
	}
	// Users
	for i := range p.Users {
		p.Users[i].UserName, err = decrypt(key, p.Users[i].UserName)
		if err != nil {
			return
		}
		p.Users[i].Password, err = decrypt(key, p.Users[i].Password)
		if err != nil {
			return
		}
	}
	// Slices
	for i := range p.Slices {
		p.Slices[i].UserName, err = decrypt(key, p.Slices[i].UserName)
		if err != nil {
			return
		}
		p.Slices[i].Password, err = decrypt(key, p.Slices[i].Password)
		if err != nil {
			return
		}
	}

	return nil
}

// Encrypt encrypt user/password in namespace
func (p *Namespace) Encrypt(key string) (err error) {
	p.IsEncrypt = true
	// Users
	for i := range p.Users {
		p.Users[i].UserName, err = encrypt(key, p.Users[i].UserName)
		if err != nil {
			return
		}
		p.Users[i].Password, err = encrypt(key, p.Users[i].Password)
		if err != nil {
			return
		}
	}
	// Slices
	for i := range p.Slices {
		p.Slices[i].UserName, err = encrypt(key, p.Slices[i].UserName)
		if err != nil {
			return
		}
		p.Slices[i].Password, err = encrypt(key, p.Slices[i].Password)
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
