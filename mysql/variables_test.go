// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetEqualsWith(t *testing.T) {
	// Test case 1: dst has variables not present in s
	t.Run("Dst has additional variables", func(t *testing.T) {
		// Create src and dst variables
		s := NewSessionVariables()
		dst := NewSessionVariables()

		// Adding a variable that only exists in dst
		dst.Set("transaction_isolation", "READ-COMMITTED")

		// Store the original value of dst for comparison
		originalDstVariables := make(map[string]interface{})
		for key, v := range dst.variables {
			originalDstVariables[key] = v.Get()
		}

		// Check if src has been updated
		changed, err := s.SetEqualsWith(dst)
		assert.Nil(t, err, "No error should occur during variable synchronization")
		assert.True(t, changed, "The 'changed' flag should be true as there are updates and movements of variables")
		assert.Nil(t, err, "No error should occur when setting variables")
		assert.True(t, changed, "The 'changed' flag should be true as variables were added")
		assert.Equal(t, "READ-COMMITTED", s.variables["transaction_isolation"].Get(), "The value of 'transaction_isolation' should match the dst")

		// Check if dst has not been altered
		for key, originalValue := range originalDstVariables {
			currentVar, exists := dst.variables[key]
			assert.True(t, exists, fmt.Sprintf("%s should still exist in dst", key))
			assert.Equal(t, originalValue, currentVar.Get(), fmt.Sprintf("The value of %s in dst should not change", key))
		}
		for key := range dst.variables {
			_, exists := originalDstVariables[key]
			assert.True(t, exists, fmt.Sprintf("No new variables should be added to dst, found %s", key))
		}
	})

	// Test case 2: s has variables not present in dst (should be moved to unused)
	t.Run("S has additional variables", func(t *testing.T) {
		// Create src and dst variables
		s := NewSessionVariables()
		dst := NewSessionVariables()

		// Adding a variable that only exists in s
		s.Set("max_execution_time", int64(1000))

		// Store the original value of dst for comparison
		originalDstVariables := make(map[string]interface{})
		for key, v := range dst.variables {
			originalDstVariables[key] = v.Get()
		}

		// Check if src has been updated
		changed, err := s.SetEqualsWith(dst)
		assert.Nil(t, err, "No error should occur during variable synchronization")
		assert.True(t, changed, "The 'changed' flag should be true as there are updates and movements of variables")
		_, existsInVariables := s.variables["max_execution_time"]
		_, existsInUnused := s.unused["max_execution_time"]
		assert.False(t, existsInVariables, "The 'max_execution_time' should not exist in active variables")
		assert.True(t, existsInUnused, "The 'max_execution_time' should exist in unused variables")

		// Check if dst has not been altered
		for key, originalValue := range originalDstVariables {
			currentVar, exists := dst.variables[key]
			assert.True(t, exists, fmt.Sprintf("%s should still exist in dst", key))
			assert.Equal(t, originalValue, currentVar.Get(), fmt.Sprintf("The value of %s in dst should not change", key))
		}
		for key := range dst.variables {
			_, exists := originalDstVariables[key]
			assert.True(t, exists, fmt.Sprintf("No new variables should be added to dst, found %s", key))
		}

	})

	// Test case 3: Both s and dst have the same variables but different values
	t.Run("Same variables, different values", func(t *testing.T) {
		// Create src and dst variables
		s := NewSessionVariables()
		dst := NewSessionVariables()

		// Adding a variable that exists in src and dst
		s.Set("unique_checks", int64(0))
		dst.Set("unique_checks", int64(1))

		// Store the original value of dst for comparison
		originalDstVariables := make(map[string]interface{})
		for key, v := range dst.variables {
			originalDstVariables[key] = v.Get()
		}

		// Check if src has been updated
		changed, err := s.SetEqualsWith(dst)
		assert.Nil(t, err, "No error should occur when updating variables")
		assert.True(t, changed, "The 'changed' flag should be true as variables values were updated")
		assert.Equal(t, int64(1), s.variables["unique_checks"].Get(), "The value of 'unique_checks' should be updated to match dst")

		// Check if dst has not been altered
		for key, originalValue := range originalDstVariables {
			currentVar, exists := dst.variables[key]
			assert.True(t, exists, fmt.Sprintf("%s should still exist in dst", key))
			assert.Equal(t, originalValue, currentVar.Get(), fmt.Sprintf("The value of %s in dst should not change", key))
		}
		for key := range dst.variables {
			_, exists := originalDstVariables[key]
			assert.True(t, exists, fmt.Sprintf("No new variables should be added to dst, found %s", key))
		}

	})
	// Test case 4: S and dst have partially overlapping variables
	t.Run("Partially overlapping variables", func(t *testing.T) {
		// Create src and dst variables
		s := NewSessionVariables()
		dst := NewSessionVariables()

		// Adding variables to src (s) and dst
		s.Set("key1", "value1")     // unique to s
		s.Set("key2", "value2-old") // overlaps with dst, different value
		s.Set("key3", "value3")     // overlaps with dst, same value

		dst.Set("key2", "value2-new") // overlaps with s, updated value
		dst.Set("key3", "value3")     // overlaps with s, same value
		dst.Set("key4", "value4")     // unique to dst

		// Store the original value of dst for comparison
		originalDstVariables := make(map[string]interface{})
		for key, v := range dst.variables {
			originalDstVariables[key] = v.Get()
		}

		// Check if src has been updated
		changed, err := s.SetEqualsWith(dst)
		assert.Nil(t, err, "No error should occur when syncing variables")
		assert.True(t, changed, "The 'changed' flag should be true as there are variable updates and movements")

		_, key1ExistsInVariables := s.variables["key1"]
		_, key1ExistsInUnused := s.unused["key1"]
		assert.False(t, key1ExistsInVariables, "key1 should not exist in active variables")
		assert.True(t, key1ExistsInUnused, "key1 should exist in unused variables")

		assert.Equal(t, "value2-new", s.variables["key2"].Get(), "key2 should be updated to 'value2-new'")
		assert.Equal(t, "value3", s.variables["key3"].Get(), "key3 should remain unchanged with 'value3'")

		_, key4ExistsInVariables := s.variables["key4"]
		assert.True(t, key4ExistsInVariables, "key4 should be added to active variables")
		assert.Equal(t, "value4", s.variables["key4"].Get(), "key4 should have the value 'value4'")

		// Check if dst has not been altered
		for key, originalValue := range originalDstVariables {
			currentVar, exists := dst.variables[key]
			assert.True(t, exists, fmt.Sprintf("%s should still exist in dst", key))
			assert.Equal(t, originalValue, currentVar.Get(), fmt.Sprintf("The value of %s in dst should not change", key))
		}
		for key := range dst.variables {
			_, exists := originalDstVariables[key]
			assert.True(t, exists, fmt.Sprintf("No new variables should be added to dst, found %s", key))
		}
	})

	// Test case 5: S and dst have partially overlapping variables of type bool
	t.Run("Partially overlapping variables with bool", func(t *testing.T) {
		// Create src and dst variables
		s := NewSessionVariables()
		dst := NewSessionVariables()

		// Adding bool variables to src (s) and dst
		s.Set("flag1", true)  // unique to s
		s.Set("flag2", false) // overlaps with dst, different value
		s.Set("flag3", true)  // overlaps with dst, same value

		dst.Set("flag2", true)  // overlaps with s, updated value
		dst.Set("flag3", true)  // overlaps with s, same value
		dst.Set("flag4", false) // unique to dst

		// Store the original value of dst for comparison
		originalDstVariables := make(map[string]interface{})
		for key, v := range dst.variables {
			originalDstVariables[key] = v.Get()
		}

		// Check if src has been updated
		changed, err := s.SetEqualsWith(dst)
		assert.Nil(t, err, "No error should occur when syncing variables")
		assert.True(t, changed, "The 'changed' flag should be true as there are variable updates and movements")

		_, flag1ExistsInVariables := s.variables["flag1"]
		_, flag1ExistsInUnused := s.unused["flag1"]
		assert.False(t, flag1ExistsInVariables, "flag1 should not exist in active variables")
		assert.True(t, flag1ExistsInUnused, "flag1 should exist in unused variables")

		assert.Equal(t, true, s.variables["flag2"].Get(), "flag2 should be updated to 'true'")
		assert.Equal(t, true, s.variables["flag3"].Get(), "flag3 should remain unchanged with 'true'")

		_, flag4ExistsInVariables := s.variables["flag4"]
		assert.True(t, flag4ExistsInVariables, "flag4 should be added to active variables")
		assert.Equal(t, false, s.variables["flag4"].Get(), "flag4 should have the value 'false'")

		// Check if dst has not been altered
		for key, originalValue := range originalDstVariables {
			currentVar, exists := dst.variables[key]
			assert.True(t, exists, fmt.Sprintf("%s should still exist in dst", key))
			assert.Equal(t, originalValue, currentVar.Get(), fmt.Sprintf("The value of %s in dst should not change", key))
		}
		for key := range dst.variables {
			_, exists := originalDstVariables[key]
			assert.True(t, exists, fmt.Sprintf("No new variables should be added to dst, found %s", key))
		}
	})

	// Test case6 : S and dst have partially overlapping variables of type int64
	t.Run("Partially overlapping variables with int", func(t *testing.T) {
		// Create src and dst variables
		s := NewSessionVariables()
		dst := NewSessionVariables()

		// Adding int64 variables to src (s) and dst
		s.Set("key1", 100) // unique to s
		s.Set("key2", 200) // overlaps with dst, different value
		s.Set("key3", 300) // overlaps with dst, same value

		dst.Set("key2", 250) // overlaps with s, updated value
		dst.Set("key3", 300) // overlaps with s, same value
		dst.Set("key4", 400) // unique to dst
		// Store the original value of dst for comparison
		originalDstVariables := make(map[string]interface{})
		for key, v := range dst.variables {
			originalDstVariables[key] = v.Get()
		}

		// Check if src has been updated
		changed, err := s.SetEqualsWith(dst)
		assert.Nil(t, err, "No error should occur when syncing variables")
		assert.True(t, changed, "The 'changed' flag should be true as there are variable updates and movements")

		_, key1ExistsInVariables := s.variables["key1"]
		_, key1ExistsInUnused := s.unused["key1"]
		assert.False(t, key1ExistsInVariables, "key1 should not exist in active variables")
		assert.True(t, key1ExistsInUnused, "key1 should exist in unused variables")

		assert.Equal(t, 250, s.variables["key2"].Get(), "key2 should be updated to '250'")
		assert.Equal(t, 300, s.variables["key3"].Get(), "key3 should remain unchanged with '300'")

		_, key4ExistsInVariables := s.variables["key4"]
		assert.True(t, key4ExistsInVariables, "key4 should be added to active variables")
		assert.Equal(t, 400, s.variables["key4"].Get(), "key4 should have the value '400'")

		// Check if dst has not been altered
		for key, originalValue := range originalDstVariables {
			currentVar, exists := dst.variables[key]
			assert.True(t, exists, fmt.Sprintf("%s should still exist in dst", key))
			assert.Equal(t, originalValue, currentVar.Get(), fmt.Sprintf("The value of %s in dst should not change", key))
		}
		for key := range dst.variables {
			_, exists := originalDstVariables[key]
			assert.True(t, exists, fmt.Sprintf("No new variables should be added to dst, found %s", key))
		}
	})
}
