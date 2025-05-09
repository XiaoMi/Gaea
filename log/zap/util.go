// Copyright 2025 The Gaea Authors. All Rights Reserved.
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

package zap

import "strings"

const (
	LogBlockStrategy   = iota // 阻塞等待
	LogDiscardStrategy        // 不阻塞等待
)

// LevelFromStr get log level from level string
func getLogStrategyFromStr(s string) int {
	strategyLower := strings.ToLower(s)
	switch strategyLower {
	case "block":
		return LogBlockStrategy
	case "discard":
		return LogDiscardStrategy
	default:
		return LogBlockStrategy
	}
}
