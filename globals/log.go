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
// date: 2023-02-23

package globals

import (
	"github.com/fuyibing/log/conf"
	"sync"
	"time"
)

var (
	poolLogs sync.Pool
)

type (
	// Log
	// is a item of custom content.
	Log struct {
		Level conf.Level
		Text  string
		Time  time.Time
		Type  LogType
	}

	// LogFormatter
	// generate log as target format string.
	LogFormatter func(log *Log) string

	LogType int
)

const (
	_ LogType = iota

	LogTypeInternal
	LogTypeSpan
)

func NewLog() (log *Log) {
	if x := poolLogs.Get(); x != nil {
		log = x.(*Log)
	} else {
		log = &Log{}
		log.init()
	}

	log.before()
	return
}

// /////////////////////////////////////////////////////////////////////////////
// Log: access
// /////////////////////////////////////////////////////////////////////////////

func (o *Log) after() {
	o.Text = ""
}

func (o *Log) before() {
	o.Time = time.Now()
}

func (o *Log) init() {}
