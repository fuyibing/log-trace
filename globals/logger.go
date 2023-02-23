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
)

var (
	Logger LoggerManager
)

type (
	LoggerManager interface {
		Debug(text string, args ...interface{})
		Info(text string, args ...interface{})
		Warn(text string, args ...interface{})
		Error(text string, args ...interface{})
		Fatal(text string, args ...interface{})
	}

	logger struct {
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Logger: operations
// /////////////////////////////////////////////////////////////////////////////

func (o *logger) Debug(text string, args ...interface{}) {
	if conf.Config.DebugOn() {
		Collector.PushBaseLog(conf.Debug, text, args...)
	}
}

func (o *logger) Info(text string, args ...interface{}) {
	if conf.Config.InfoOn() {
		Collector.PushBaseLog(conf.Info, text, args...)
	}
}

func (o *logger) Warn(text string, args ...interface{}) {
	if conf.Config.WarnOn() {
		Collector.PushBaseLog(conf.Warn, text, args...)
	}
}

func (o *logger) Error(text string, args ...interface{}) {
	if conf.Config.ErrorOn() {
		Collector.PushBaseLog(conf.Error, text, args...)
	}
}

func (o *logger) Fatal(text string, args ...interface{}) {
	if conf.Config.FatalOn() {
		Collector.PushBaseLog(conf.Fatal, text, args...)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Logger: access
// /////////////////////////////////////////////////////////////////////////////

func (o *logger) init() *logger {
	return o
}
