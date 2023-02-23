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
	"fmt"
	"github.com/fuyibing/log/conf"
	"net/http"
)

var (
	Collector CollectorManager
)

type (
	CollectorManager interface {
		NewTrace(name string) Trace
		NewTraceFromRequest(name string, request *http.Request) Trace

		collectorLogger
		collectorExporter
	}

	collector struct {
		LogExporter  LogExporter
		SpanExporter SpanExporter
	}

	collectorExporter interface {
		SetLogExporter(v LogExporter)
		SetSpanExporter(v SpanExporter)
	}

	collectorLogger interface {
		PushBaseLog(level conf.Level, text string, args ...interface{})
		PushSpan(span Span)
		PushSpanLog(span Span, level conf.Level, text string, args ...interface{})
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Collector: operate
// /////////////////////////////////////////////////////////////////////////////

func (o *collector) NewTrace(name string) Trace {
	return o.empty(name,
		WithTraceId(Generator.TraceId()),
	)
}

func (o *collector) NewTraceFromRequest(name string, request *http.Request) Trace {
	return nil
}

// /////////////////////////////////////////////////////////////////////////////
// Collector: push
// /////////////////////////////////////////////////////////////////////////////

func (o *collector) PushBaseLog(level conf.Level, text string, args ...interface{}) {
	// Ignore
	// if exporter not registered.
	if o.LogExporter == nil {
		return
	}

	// Build
	// log fields.
	v := NewLog()
	v.Level = level
	v.Text = fmt.Sprintf(text, args...)
	v.Type = LogTypeInternal

	// Push immediately.
	o.LogExporter.Push(v)
}

func (o *collector) PushSpan(span Span) {
	// Ignore
	// if exporter not registered.
	if o.SpanExporter == nil {
		return
	}

	// Push immediately.
	if err := o.SpanExporter.Push(span); err != nil {
		println("push error: ", err.Error(), conf.Config.GetJaeger().GetEndpoint())
	}
}

func (o *collector) PushSpanLog(span Span, level conf.Level, text string, args ...interface{}) {
	// Ignore
	// if exporter not registered.
	if o.LogExporter == nil {
		return
	}

	// Build
	// log fields.
	v := NewLog()
	v.Level = level
	v.Text = fmt.Sprintf(text, args...)
	v.Type = LogTypeSpan

	// Push immediately.
	o.LogExporter.Push(v)
}

// /////////////////////////////////////////////////////////////////////////////
// Collector: config
// /////////////////////////////////////////////////////////////////////////////

func (o *collector) SetLogExporter(v LogExporter)   { o.LogExporter = v }
func (o *collector) SetSpanExporter(v SpanExporter) { o.SpanExporter = v }

// /////////////////////////////////////////////////////////////////////////////
// Collector: access
// /////////////////////////////////////////////////////////////////////////////

func (o *collector) empty(name string, tos ...TraceOption) Trace {
	if x := poolTrace.Get(); x != nil {
		return x.(*trace).
			before(name).
			apply(tos)
	}

	return (&trace{}).
		init().
		before(name).
		apply(tos)
}

func (o *collector) init() *collector { return o }
