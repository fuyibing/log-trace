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

package config

type (
	LoggerLevel string

	// LoggerName
	// name of logger exporter.
	LoggerName string
)

const (
	Off   LoggerLevel = "OFF"
	Fatal LoggerLevel = "FATAL"
	Error LoggerLevel = "ERROR"
	Warn  LoggerLevel = "WARN"
	Info  LoggerLevel = "INFO"
	Debug LoggerLevel = "DEBUG"

	LevelDefault = Info
)

const (
	LoggerTerm  LoggerName = "term"
	LoggerFile  LoggerName = "file"
	LoggerKafka LoggerName = "kafka"
)

var (
	levelIntegers = map[LoggerLevel]int{
		Off:   1,
		Fatal: 2,
		Error: 3,
		Warn:  4,
		Info:  5,
		Debug: 6,
	}
)

func (o LoggerLevel) Int() int {
	if i, ok := levelIntegers[o]; ok {
		return i
	}
	return 0
}

func (o LoggerLevel) String() string {
	return string(o)
}
