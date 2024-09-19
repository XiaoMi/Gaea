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

import "strings"

func LowerEqual(src string, dest string) bool {
	if len(src) != len(dest) {
		return false
	}
	return strings.ToLower(src) == dest
}

func UpperEqual(src string, dest string) bool {
	if len(src) != len(dest) {
		return false
	}
	return strings.ToUpper(src) == dest
}

func HasUpperPrefix(src string, dest string) bool {
	if len(src) < len(dest) {
		return false
	}
	return strings.ToUpper(src[0:len(dest)]) == dest
}
