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

package sequence

import "fmt"

// Sequence is interface of global sequences with different types
type Sequence interface {
	GetPKName() string
	NextSeq() (int64, error)
}

type SequenceManager struct {
	sequences map[string]map[string]Sequence
}

func NewSequenceManager() *SequenceManager {
	return &SequenceManager{
		sequences: make(map[string]map[string]Sequence),
	}
}

func (s *SequenceManager) SetSequence(db, table string, seq Sequence) error {
	if _, ok := s.sequences[db]; !ok {
		s.sequences[db] = make(map[string]Sequence)
	}
	if _, ok := s.sequences[db][table]; ok {
		return fmt.Errorf("already set sequence, db: %s, table: %s", db, table)
	}

	s.sequences[db][table] = seq
	return nil
}

func (s *SequenceManager) GetSequence(db, table string) (Sequence, bool) {
	dbSeq, ok := s.sequences[db]
	if !ok {
		return nil, false
	}
	seq, ok := dbSeq[table]
	return seq, ok
}
