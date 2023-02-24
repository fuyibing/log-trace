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
	"context"
	"fmt"
	"github.com/fuyibing/log/tracer"
	"os"
)

type (
	Exporter interface {
		Push(log *tracer.Log) error
		Start(ctx context.Context) error
		Stopped() bool
	}

	exporter struct {
		formatter Formatter
	}
)

func New() Exporter {
	return (&exporter{}).init()
}

func (o *exporter) Push(log *tracer.Log) error {
	o.println(log)
	return nil
}

func (o *exporter) Start(_ context.Context) error {
	return nil
}

func (o *exporter) Stopped() bool { return true }

// /////////////////////////////////////////////////////////////////////////////
// Exporter: constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *exporter) init() *exporter {
	o.formatter = (&formatter{}).init()
	return o
}

// println
// print log content on terminal.
func (o *exporter) println(log *tracer.Log) {
	_, _ = fmt.Fprintf(os.Stdout,
		fmt.Sprintf("%s\n",
			o.formatter.Format(log),
		),
	)
}
