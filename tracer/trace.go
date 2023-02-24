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
	"github.com/fuyibing/log/config"
	"net/http"
)

type (
	// Trace
	// is the creator of Spans.
	Trace interface {
		traceNewer
		traceGetter
		traceSetter
	}

	// TraceId
	// is a unique identity of a trace.
	TraceId struct {
		bs       [16]byte
		err      error
		security bool
		str      string
	}

	trace struct {
		attr Attr
		ctx  context.Context
		name string

		provider ProviderManager
		spanId   SpanId
		traceId  TraceId
	}

	traceNewer interface {
		NewSpan(name string) Span
		NewSpanWithContext(ctx context.Context, name string) Span
	}

	traceGetter interface {
		GetAttr() Attr
		GetContext() context.Context
		GetName() string
		GetProvider() ProviderManager
		GetSpanId() SpanId
		GetTraceId() TraceId
	}

	traceSetter interface {
		SetAttr(key string, value interface{}) Trace
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Trace: returns a span based on trace
// /////////////////////////////////////////////////////////////////////////////

// NewSpan
// returns a span which created from a trace.
func (o *trace) NewSpan(name string) Span {
	v := (&span{}).init(name)
	v.attr.Copy(o.attr)
	v.parentSpanId = o.spanId
	v.trace = o
	return v
}

// NewSpanWithContext
// returns a span which created from a trace and with specified context.
func (o *trace) NewSpanWithContext(ctx context.Context, name string) Span {
	v := (&span{}).init(name)
	v.attr.Copy(o.attr)
	v.ctx = ctx
	v.parentSpanId = o.spanId
	v.trace = o
	return v
}

// /////////////////////////////////////////////////////////////////////////////
// Trace: getter
// /////////////////////////////////////////////////////////////////////////////

func (o *trace) GetAttr() Attr                { return o.attr }
func (o *trace) GetContext() context.Context  { return o.ctx }
func (o *trace) GetName() string              { return o.name }
func (o *trace) GetProvider() ProviderManager { return o.provider }
func (o *trace) GetSpanId() SpanId            { return o.spanId }
func (o *trace) GetTraceId() TraceId          { return o.traceId }

// /////////////////////////////////////////////////////////////////////////////
// Trace: setter
// /////////////////////////////////////////////////////////////////////////////

func (o *trace) SetAttr(key string, value interface{}) Trace {
	o.attr.Add(key, value)
	return o
}

// /////////////////////////////////////////////////////////////////////////////
// TraceId: readonly
// /////////////////////////////////////////////////////////////////////////////

func (o TraceId) Byte() []byte   { return o.bs[:] }
func (o TraceId) Err() error     { return o.err }
func (o TraceId) Security() bool { return o.security }

func (o TraceId) String() string {
	if o.str == "" {
		o.str = hex.EncodeToString(o.bs[:])
	}
	return o.str
}

// /////////////////////////////////////////////////////////////////////////////
// Trace: access
// /////////////////////////////////////////////////////////////////////////////

func (o *trace) init() *trace {
	return o
}

func (o *trace) useRequest(req *http.Request) {
	var (
		sid = req.Header.Get(config.Config.GetOpenTracingSpanId())
		tid = req.Header.Get(config.Config.GetOpenTracingTraceId())
	)

	o.attr.Add("http.header", req.Header)
	o.attr.Add("http.request.url", req.RequestURI)
	o.attr.Add("http.request.method", req.Method)
	o.attr.Add("http.request.protocol", req.Proto)
	o.attr.Add("http.user.agent", req.UserAgent())

	if tid != "" && sid != "" {
		if o.spanId = Identify.HexSpanId(sid); o.spanId.Err() != nil {
		}

		if o.traceId = Identify.HexTraceId(tid); o.traceId.Err() != nil {
		}
	}
}

func (o *trace) useRoot() {
	o.spanId = Identify.NewEmptySpanId()
	o.traceId = Identify.NewTraceId()
}
