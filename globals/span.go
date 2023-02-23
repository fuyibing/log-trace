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
	"encoding/hex"
	"fmt"
	"github.com/fuyibing/log/conf"
	"sync"
	"time"
)

var (
	poolSpan sync.Pool
)

type (
	// Span
	// is a component of the trace. It represents a single named and timed
	// operation of a workflow that is traced.
	Span interface {

		// NewSpan
		// returns a child span which created from span.
		NewSpan(name string, sos ...SpanOption) Span

		spanGetter
		spanSetter
		spanLogger
	}

	// SpanId
	// is a unique identity of a span.
	SpanId [8]byte

	// SpanOption
	// apply an option to a span.
	SpanOption func(s *span)

	SpanReadonly interface {
		spanGetter
	}

	span struct {
		sync.Mutex

		active       bool
		attributes   map[string]interface{}
		endTime      time.Time
		logs         []*Log
		name         string
		parentSpanId SpanId
		spanId       SpanId
		startTime    time.Time
		trace        Trace
		traceId      TraceId
	}

	spanGetter interface {
		GetAttributes() map[string]interface{}
		GetDuration() time.Duration
		GetEndTime() time.Time
		GetLogs() []*Log
		GetName() string
		GetParentSpanId() SpanId
		GetSpanId() SpanId
		GetStartTime() time.Time
		GetTraceId() TraceId
	}

	spanSetter interface {
		// End
		// completes the Span. The Span is considered complete and ready to be
		// delivered by Collector dispatch.
		End()

		SetAttribute(key string, value interface{})
	}

	spanLogger interface {
		Debug(text string, args ...interface{})
		Info(text string, args ...interface{})
		Warn(text string, args ...interface{})
		Error(text string, args ...interface{})
		Fatal(text string, args ...interface{})
	}
)

// String
// returns the hex string representation form of a SpanId,
// Total length is 16.
func (b SpanId) String() string {
	return hex.EncodeToString(b[:])
}

// /////////////////////////////////////////////////////////////////////////////
// Span: create child span
// /////////////////////////////////////////////////////////////////////////////

func (o *span) NewSpan(name string, sos ...SpanOption) Span {
	if x := poolTrace.Get(); x != nil {
		return x.(*span).
			bind(o.trace, o.traceId, o.spanId).
			before(name).
			apply(sos)
	}

	return (&span{}).init().
		bind(o.trace, o.traceId, o.spanId).
		before(name).
		apply(sos)
}

// /////////////////////////////////////////////////////////////////////////////
// Span: custom operations
// /////////////////////////////////////////////////////////////////////////////

func (o *span) End() {
	// Update span state.
	o.Lock()
	o.active = false
	o.endTime = time.Now()
	o.Unlock()

	// Publish span.
	Collector.PushSpan(o)
}

// /////////////////////////////////////////////////////////////////////////////
// Span: getter group
// /////////////////////////////////////////////////////////////////////////////

func (o *span) GetAttributes() map[string]interface{} { return o.attributes }
func (o *span) GetDuration() time.Duration            { return o.endTime.Sub(o.startTime) }
func (o *span) GetEndTime() time.Time                 { return o.endTime }

func (o *span) GetLogs() []*Log {
	o.Lock()
	defer o.Unlock()
	return o.logs
}

func (o *span) GetName() string         { return o.name }
func (o *span) GetParentSpanId() SpanId { return o.parentSpanId }
func (o *span) GetSpanId() SpanId       { return o.spanId }
func (o *span) GetStartTime() time.Time { return o.startTime }
func (o *span) GetTraceId() TraceId     { return o.traceId }

// /////////////////////////////////////////////////////////////////////////////
// Span: logger setter
// /////////////////////////////////////////////////////////////////////////////

func (o *span) SetAttribute(key string, value interface{}) {
	o.Lock()
	defer o.Unlock()
	o.attributes[key] = value
}

// /////////////////////////////////////////////////////////////////////////////
// Span: logger group
// /////////////////////////////////////////////////////////////////////////////

func (o *span) Debug(text string, args ...interface{}) {
	if conf.Config.DebugOn() {
		o.pushLog(conf.Debug, text, args...)
	}
}

func (o *span) Error(text string, args ...interface{}) {
	if conf.Config.ErrorOn() {
		o.pushLog(conf.Error, text, args...)
	}
}

func (o *span) Fatal(text string, args ...interface{}) {
	if conf.Config.FatalOn() {
		o.pushLog(conf.Fatal, text, args...)
	}
}

func (o *span) Info(text string, args ...interface{}) {
	if conf.Config.InfoOn() {
		o.pushLog(conf.Info, text, args...)
	}
}

func (o *span) Warn(text string, args ...interface{}) {
	if conf.Config.WarnOn() {
		o.pushLog(conf.Warn, text, args...)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Span: init and pool
// /////////////////////////////////////////////////////////////////////////////

func (o *span) after() {
	o.attributes = nil
	o.logs = nil
	o.trace = nil
}

func (o *span) apply(sos []SpanOption) *span {
	for _, so := range sos {
		so(o)
	}
	return o
}

func (o *span) before(name string) *span {
	o.active = true
	o.attributes = make(map[string]interface{})
	o.logs = make([]*Log, 0)
	o.name = name
	o.spanId = Generator.SpanId()
	o.startTime = time.Now()
	return o
}

func (o *span) bind(t Trace, tid TraceId, sid SpanId) *span {
	o.trace = t
	o.traceId = tid
	o.parentSpanId = sid
	return o
}

func (o *span) init() *span {
	o.logs = make([]*Log, 0)
	return o
}

func (o *span) pushLog(level conf.Level, text string, args ...interface{}) {
	o.Lock()
	defer o.Unlock()

	// Append
	// log on the span.
	v := NewLog()
	v.Level = level
	v.Text = fmt.Sprintf(text, args...)
	o.logs = append(o.logs, v)

	// Push
	// base log.
	Collector.PushSpanLog(o, level, text, args...)
}
