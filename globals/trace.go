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
	"sync"
)

var (
	poolTrace sync.Pool
)

type (
	// Trace
	// is the creator of spans.
	Trace interface {
		// NewSpan
		// returns a span which created from trace.
		NewSpan(name string, sos ...SpanOption) Span

		traceGetter
		traceSetter
	}

	// TraceId
	// is a unique identity of a trace.
	TraceId [16]byte

	trace struct {
		name string

		// spanId
		// initialized from parent.
		spanId SpanId

		// spanId
		// initialize from parent or root.
		traceId TraceId
	}

	traceGetter interface {
		// GetName
		// returns a trace name.
		GetName() string

		// GetSpanId
		// returns a span identify string.
		GetSpanId() SpanId

		// GetTraceId
		// returns a trace identify string.
		GetTraceId() TraceId
	}

	traceSetter interface{}
)

// String
// returns the hex string representation form of a TraceId,
// Total length is 32.
func (b TraceId) String() string {
	return hex.EncodeToString(b[:])
}

// /////////////////////////////////////////////////////////////////////////////
// Trace: exported methods
// /////////////////////////////////////////////////////////////////////////////

func (o *trace) NewSpan(name string, sos ...SpanOption) Span {
	if x := poolSpan.Get(); x != nil {
		return x.(*span).
			bind(o, o.traceId, o.spanId).
			before(name).
			apply(sos)
	}

	return (&span{}).init().
		bind(o, o.traceId, o.spanId).
		before(name).
		apply(sos)
}

// /////////////////////////////////////////////////////////////////////////////
// Trace: getter
// /////////////////////////////////////////////////////////////////////////////

func (o *trace) GetName() string     { return o.name }
func (o *trace) GetSpanId() SpanId   { return o.spanId }
func (o *trace) GetTraceId() TraceId { return o.traceId }

// /////////////////////////////////////////////////////////////////////////////
// Trace: init & pool
// /////////////////////////////////////////////////////////////////////////////

func (o *trace) after() {}

func (o *trace) apply(tos []TraceOption) *trace {
	for _, to := range tos {
		to(o)
	}
	return o
}

func (o *trace) before(name string) *trace {
	o.name = name
	return o
}

func (o *trace) init() *trace {
	return o
}
