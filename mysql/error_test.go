// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"github.com/pingcap/check"
	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ = check.Suite(&testSQLErrorSuite{})

type testSQLErrorSuite struct {
}

func (s *testSQLErrorSuite) TestSQLError(c *check.C) {
	e := NewErrf(ErrNoDB, "no db error")
	c.Assert(len(e.Error()), check.Greater, 0)

	e = NewErrf(0, "customized error")
	c.Assert(len(e.Error()), check.Greater, 0)

	e = NewDefaultError(ErrNoDB)
	c.Assert(len(e.Error()), check.Greater, 0)

	e = NewDefaultError(0, "customized error")
	c.Assert(len(e.Error()), check.Greater, 0)
}

func TestRespError(t *testing.T) {
	v := 0
	var writeResp = func(interface{}) error {
		v += 1
		return nil
	}
	var writeRespErr = func(interface{}) error {
		v += 1
		return errors.New("write resp error")
	}
	err1 := genRespErr(errors.New("normal error"), writeResp)
	assert.Equal(t, err1, nil)
	assert.Equal(t, v, 1)
	err2 := genRespErr(NewSessionCloseError("close no resp"), writeResp)
	assert.NotEqual(t, err2, nil)
	assert.Equal(t, err2.Error(), "close no resp")
	assert.Equal(t, v, 1)
	err3 := genRespErr(NewSessionCloseRespError("close resp"), writeResp)
	assert.NotEqual(t, err3, nil)
	assert.Equal(t, err3.Error(), "close resp")
	assert.Equal(t, v, 2)
	err4 := genRespErr(errors.New("normal error"), writeRespErr)
	assert.NotEqual(t, err4, nil)
	assert.Equal(t, err4.Error(), "write resp error")
	assert.Equal(t, v, 3)
	err5 := genRespErr(NewSessionCloseError("close no resp"), writeRespErr)
	assert.NotEqual(t, err5, nil)
	assert.Equal(t, err5.Error(), "close no resp")
	assert.Equal(t, v, 3)
	err6 := genRespErr(NewSessionCloseRespError("close resp"), writeRespErr)
	assert.NotEqual(t, err6, nil)
	assert.Equal(t, err6.Error(), "close resp")
	assert.Equal(t, v, 4)
}

func genRespErr(err interface{}, writeResp func(interface{}) error) error {
	switch err.(type) {
	case *SessionCloseRespError:
		writeResp(err)
		return err.(*SessionCloseRespError)
	case *SessionCloseNoRespError:
		return err.(*SessionCloseNoRespError)
	case error:
		return writeResp(err.(error))
	}
	return nil
}
