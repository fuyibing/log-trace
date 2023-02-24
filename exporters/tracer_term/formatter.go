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

package tracer_term

import (
	"fmt"
	"github.com/fuyibing/log/tracer"
)

type (
	Formatter interface {
		Format(span tracer.Span) []string
	}

	formatter struct {
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Formatter: access
// /////////////////////////////////////////////////////////////////////////////

func (o *formatter) Format(span tracer.Span) (list []string) {
	list = []string{
		fmt.Sprintf("Span [%s][%vus] | %s | %s",
			span.GetSpanId().String(),
			span.GetDuration().Microseconds(),
			span.GetName(),
			span.GetAttr().JSON(),
		),
	}

	for k, v := range span.GetTrace().GetProvider().GetAttr() {
		list = append(list, fmt.Sprintf("     [%s] : {%v}", k, v))
	}

	for i, log := range span.GetLogs() {
		if i == 0 {
			list = append(list, "     +--- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ----")
		}
		list = append(list,
			fmt.Sprintf("     + [%-15s][%5s] %s",
				log.Time.Format("15:04:05.999999"),
				log.Level.String(),
				log.Text,
			),
		)
	}

	return
}

func (o *formatter) init() *formatter {
	return o
}
