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
)

var _ = check.Suite(&testTypeSuite{})

type testTypeSuite struct{}

func (s *testTypeSuite) TestFlags(c *check.C) {
	c.Assert(HasNotNullFlag(NotNullFlag), check.IsTrue)
	c.Assert(HasUniKeyFlag(UniqueKeyFlag), check.IsTrue)
	c.Assert(HasNotNullFlag(NotNullFlag), check.IsTrue)
	c.Assert(HasNoDefaultValueFlag(NoDefaultValueFlag), check.IsTrue)
	c.Assert(HasAutoIncrementFlag(AutoIncrementFlag), check.IsTrue)
	c.Assert(HasUnsignedFlag(UnsignedFlag), check.IsTrue)
	c.Assert(HasZerofillFlag(ZerofillFlag), check.IsTrue)
	c.Assert(HasBinaryFlag(BinaryFlag), check.IsTrue)
	c.Assert(HasPriKeyFlag(PriKeyFlag), check.IsTrue)
	c.Assert(HasMultipleKeyFlag(MultipleKeyFlag), check.IsTrue)
	c.Assert(HasTimestampFlag(TimestampFlag), check.IsTrue)
	c.Assert(HasOnUpdateNowFlag(OnUpdateNowFlag), check.IsTrue)
}
