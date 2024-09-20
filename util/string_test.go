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

import "testing"

// TestLowerEqual tests the LowerEqual function
func TestLowerEqual(t *testing.T) {
	testCases := []struct {
		name   string
		src    string
		dest   string
		expect bool
	}{
		{"equal_case_sensitive", "hello", "hello", true},
		{"not_equal_length", "hello", "helloo", false},
		{"not_equal_content", "hello", "world", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LowerEqual(tc.src, tc.dest)
			if result != tc.expect {
				t.Errorf("Expected %v but got %v for src=%s and dest=%s", tc.expect, result, tc.src, tc.dest)
			}
		})
	}
}

// TestUpperEqual tests the UpperEqual function
func TestUpperEqual(t *testing.T) {
	testCases := []struct {
		name   string
		src    string
		dest   string
		expect bool
	}{
		{"equal_case_sensitive", "HELLO", "HELLO", true},
		{"not_equal_length", "HELLO", "HELLOO", false},
		{"not_equal_content", "HELLO", "WORLD", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := UpperEqual(tc.src, tc.dest)
			if result != tc.expect {
				t.Errorf("Expected %v but got %v for src=%s and dest=%s", tc.expect, result, tc.src, tc.dest)
			}
		})
	}
}

// TestHasUpperPrefix tests the HasUpperPrefix function
func TestHasUpperPrefix(t *testing.T) {
	testCases := []struct {
		name   string
		src    string
		dest   string
		expect bool
	}{
		{"prefix_match", "HelloWorld", "HELLO", true},
		{"prefix_mismatch", "HelloWorld", "GOD", false},
		{"shorter_source", "Hello", "HELLO", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HasUpperPrefix(tc.src, tc.dest)
			if result != tc.expect {
				t.Errorf("Expected %v but got %v for src=%s and dest=%s", tc.expect, result, tc.src, tc.dest)
			}
		})
	}
}
