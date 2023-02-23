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

package conf

type (
	// Level
	// defined for log record level.
	Level string
)

const (
	Debug Level = "DEBUG"
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
	Fatal Level = "FATAL"
)

var (
	LevelInt = map[Level]int{
		Fatal: 1,
		Error: 2,
		Warn:  3,
		Info:  4,
		Debug: 5,
	}
)

func (l Level) Int() (n int) {
	return LevelInt[l]
}

func (l Level) String() string {
	return string(l)
}
