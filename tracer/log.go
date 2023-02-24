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
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-02-24

package tracer

import (
	"github.com/fuyibing/log/config"
	"time"
)

type (
	// Log
	// an custom message.
	Log struct {
		Level config.LoggerLevel
		Time  time.Time
		Text  string
		Type  LogType
	}

	LoggerManager interface {
		Push(log *Log)
	}

	LogType int
)

const (
	_ LogType = iota
	LogInternal
	LogSpan
)

func NewLog(t LogType, l config.LoggerLevel) *Log {
	return &Log{
		Level: l,
		Time:  time.Now(),
		Type:  t,
	}
}
