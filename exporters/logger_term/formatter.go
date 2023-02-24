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

package logger_term

import (
	"fmt"
	"github.com/fuyibing/log/config"
	"github.com/fuyibing/log/tracer"
)

var (
	colors = map[config.LoggerLevel][]int{
		config.Debug: {37, 0},  // Text: gray, Background: white
		config.Info:  {34, 0},  // Text: blue, Background: white
		config.Warn:  {33, 0},  // Text: yellow, Background: white
		config.Error: {31, 0},  // Text: red, Background: white
		config.Fatal: {33, 41}, // Text: yellow, Background: red
	}
)

type (
	Formatter interface {
		Format(log *tracer.Log) string
	}

	formatter struct {
	}
)

// Format
// generate log as string used on terminal.
func (o *formatter) Format(log *tracer.Log) (text string) {
	text = fmt.Sprintf("[%-15s][%5s] %s",
		log.Time.Format("15:04:05.999999"),
		log.Level.String(),
		log.Text,
	)

	if c, ok := colors[log.Level]; ok {
		text = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m",
			0x1B, 0, c[1], c[0], text, 0x1B,
		)
	}

	return
}

// /////////////////////////////////////////////////////////////////////////////
// Formatter: constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *formatter) init() *formatter {
	return o
}
