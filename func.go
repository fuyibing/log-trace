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

package log

import (
	"github.com/fuyibing/log/config"
)

// Debug send debug level log to Provider.
func Debug(text string, args ...interface{}) {
	if config.Config.DebugOn() {
		Provider.PushBaseLog(config.Debug, text, args...)
	}
}

// Info send info level log to Provider.
func Info(text string, args ...interface{}) {
	if config.Config.InfoOn() {
		Provider.PushBaseLog(config.Info, text, args...)
	}
}

// Warn send warn level log to Provider.
func Warn(text string, args ...interface{}) {
	if config.Config.WarnOn() {
		Provider.PushBaseLog(config.Warn, text, args...)
	}
}

// Error send error level log to Provider.
func Error(text string, args ...interface{}) {
	if config.Config.ErrorOn() {
		Provider.PushBaseLog(config.Error, text, args...)
	}
}

// Fatal send fatal level log to Provider.
func Fatal(text string, args ...interface{}) {
	if config.Config.FatalOn() {
		Provider.PushBaseLog(config.Fatal, text, args...)
	}
}
