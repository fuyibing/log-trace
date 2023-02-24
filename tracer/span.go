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
	"context"
	"encoding/hex"
	"fmt"
	"github.com/fuyibing/log/config"
	"sync"
	"time"
)

type (
	Span interface {
		spanNewer
		spanGetter
		spanSetter
		spanLogger
	}

	SpanId struct {
		bs       [8]byte
		err      error
		security bool
		str      string
	}

	span struct {
		sync.RWMutex

		// attr
		// record span attributes, initialized in constructor.
		attr Attr

		ctx context.Context

		name                 string
		logs                 []*Log
		spanId, parentSpanId SpanId

		// trace
		// recorded belongs to relation.
		trace *trace

		startTime, endTime time.Time
	}

	// spanGetter interface for span reader.
	spanGetter interface {
		// Attr
		// returns an attribute fields.
		GetAttr() Attr

		// GetDuration
		// return span duration.
		GetDuration() time.Duration

		// GetEndTime
		// returns a time of span end.
		GetEndTime() time.Time

		// GetLogs
		// return logs list.
		GetLogs() []*Log

		// GetName
		// returns a span name.
		GetName() string

		// GetParentSpanId
		// returns a span id of parent.
		GetParentSpanId() SpanId

		// GetSpanId
		// returns a span id.
		GetSpanId() SpanId

		// GetStartTime
		// return span started time.
		GetStartTime() time.Time

		// GetTrace
		// returns a trace.
		GetTrace() Trace

		// GetTraceId
		// returns a trace id.
		GetTraceId() TraceId
	}

	// spanLogger interface for log sender.
	spanLogger interface {
		// Debug send debug level log on span.
		Debug(text string, args ...interface{})

		// Error send error level log on span.
		Error(text string, args ...interface{})

		// Fatal send fatal level log on span.
		Fatal(text string, args ...interface{})

		// Info send info level log on span.
		Info(text string, args ...interface{})

		// Warn send warn level log on span.
		Warn(text string, args ...interface{})
	}

	// spanNewer interface for child span creator.
	spanNewer interface {
		// NewSpan
		// returns a child Span.
		NewSpan(name string) Span

		// NewSpanWithContext
		// returns a child Span with context.
		NewSpanWithContext(ctx context.Context, name string) Span
	}

	// spanSetter interface for span setter.
	spanSetter interface {
		End()

		// SetAttr
		// set span attributes, override if exists.
		SetAttr(key string, value interface{}) Span
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Span: newer
// /////////////////////////////////////////////////////////////////////////////

func (o *span) NewSpan(name string) Span {
	v := (&span{}).init(name)
	v.parentSpanId = o.spanId
	v.trace = o.trace
	return v
}

func (o *span) NewSpanWithContext(ctx context.Context, name string) Span {
	v := (&span{}).init(name)
	v.ctx = ctx
	v.parentSpanId = o.spanId
	v.trace = o.trace
	return v
}

// /////////////////////////////////////////////////////////////////////////////
// Span: getter
// /////////////////////////////////////////////////////////////////////////////

// GetAttr
// returns an attribute fields.
func (o *span) GetAttr() Attr { return o.attr }

// GetDuration
// return span duration.
func (o *span) GetDuration() time.Duration { return o.endTime.Sub(o.startTime) }

// GetEndTime
// returns a time of span end.
func (o *span) GetEndTime() time.Time { return o.endTime }

// GetLogs
// return logs list.
func (o *span) GetLogs() []*Log { return o.logs }

// GetName
// returns a span name.
func (o *span) GetName() string { return o.name }

// GetParentSpanId
// returns a span id of parent.
func (o *span) GetParentSpanId() SpanId { return o.parentSpanId }

// GetSpanId
// returns a span id.
func (o *span) GetSpanId() SpanId { return o.spanId }

// GetStartTime
// returns a time of start end.
func (o *span) GetStartTime() time.Time { return o.startTime }

// GetTrace
// returns a trace.
//
// Lock is not necessary on read.
func (o *span) GetTrace() Trace { return o.trace }

// GetTraceId
// returns a trace id.
func (o *span) GetTraceId() TraceId { return o.trace.traceId }

// /////////////////////////////////////////////////////////////////////////////
// Span: setter
// /////////////////////////////////////////////////////////////////////////////

func (o *span) End() {
	o.Lock()
	o.endTime = time.Now()
	o.Unlock()

	o.trace.GetProvider().PushSpan(o)
}

// SetAttr
// set span attributes, override if exists.
func (o *span) SetAttr(key string, value interface{}) Span {
	o.Lock()
	defer o.Unlock()
	o.attr.Add(key, value)
	return o
}

// /////////////////////////////////////////////////////////////////////////////
// Span: logger
// /////////////////////////////////////////////////////////////////////////////

// Debug send debug level log on span.
func (o *span) Debug(text string, args ...interface{}) {
	if config.Config.DebugOn() {
		o.sendLog(config.Debug, text, args...)
	}
}

// Info send info level log on span.
func (o *span) Info(text string, args ...interface{}) {
	if config.Config.InfoOn() {
		o.sendLog(config.Info, text, args...)
	}
}

// Warn send warn level log on span.
func (o *span) Warn(text string, args ...interface{}) {
	if config.Config.WarnOn() {
		o.sendLog(config.Warn, text, args...)
	}
}

// Error send error level log on span.
func (o *span) Error(text string, args ...interface{}) {
	if config.Config.ErrorOn() {
		o.sendLog(config.Error, text, args...)
	}
}

// Fatal send fatal level log on span.
func (o *span) Fatal(text string, args ...interface{}) {
	if config.Config.FatalOn() {
		o.sendLog(config.Fatal, text, args...)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Span: access
// /////////////////////////////////////////////////////////////////////////////

func (o *span) init(name string) *span {
	o.attr = Attr{}
	o.name = name
	o.logs = make([]*Log, 0)
	o.spanId = Identify.NewSpanId()
	o.startTime = time.Now()
	return o
}

func (o *span) sendLog(level config.LoggerLevel, text string, args ...interface{}) {
	// Create log and add
	// into span containers.
	if config.Config.GetTracerWithLog() {
		o.Lock()
		x := NewLog(LogSpan, level)
		x.Text = fmt.Sprintf(text, args...)
		o.logs = append(o.logs, x)
		o.Unlock()
	}

	// Publish to basic.
	o.trace.GetProvider().PushSpanLog(config.Fatal, text, args...)
}

// /////////////////////////////////////////////////////////////////////////////
// SpanId: readonly
// /////////////////////////////////////////////////////////////////////////////

func (o SpanId) Byte() []byte   { return o.bs[:] }
func (o SpanId) Err() error     { return o.err }
func (o SpanId) Security() bool { return o.security }

func (o SpanId) String() string {
	if o.str == "" {
		o.str = hex.EncodeToString(o.bs[:])
	}
	return o.str
}
