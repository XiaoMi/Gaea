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

func ArrayFindIndex(items []string, findItem string) int {
	for index, item := range items {
		if item == findItem {
			return index
		}
	}
	return -1
}

func ArrayRemoveItem(items []string, findItem string) []string {
	newArray := make([]string, 0)
	for _, item := range items {
		if item == findItem {
			continue
		}
		newArray = append(newArray, item)
	}
	return newArray
}
