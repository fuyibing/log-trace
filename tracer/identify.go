// author: wsfuyibing <websearch@163.com>
// date: 2023-02-22

package tracer

import (
	cr "crypto/rand"
	eb "encoding/binary"
	"encoding/hex"
	mr "math/rand"
	"sync"
)

var (
	Identify IdentifyGenerator
)

type (
	// IdentifyGenerator
	// generate ids of a trace.
	IdentifyGenerator interface {
		HexSpanId(s string) SpanId
		HexTraceId(s string) TraceId
		NewEmptySpanId() SpanId
		NewSpanId() SpanId
		NewTraceId() TraceId
	}

	identify struct {
		sync.Mutex
		data   int64
		err    error
		random *mr.Rand
	}
)

// HexSpanId
// returns a SpanId with hex string from parent.
func (o *identify) HexSpanId(s string) SpanId {
	var (
		d []byte
		v = SpanId{security: false, str: s}
	)
	if d, v.err = hex.DecodeString(s); v.err == nil {
		copy(v.bs[:], d)
	}
	return v
}

// HexTraceId
// returns a TraceId with hex string from parent.
func (o *identify) HexTraceId(s string) TraceId {
	var (
		d []byte
		v = TraceId{security: false, str: s}
	)
	if d, v.err = hex.DecodeString(s); v.err == nil {
		copy(v.bs[:], d)
	}
	return v
}

// NewSpanId
// returns a SpanId with rand strategy.
func (o *identify) NewSpanId() SpanId {
	o.Lock()
	defer o.Unlock()

	s := SpanId{security: true}
	o.random.Read(s.bs[:])

	s.str = hex.EncodeToString(s.bs[:])
	return s
}

// NewEmptySpanId
// returns an empty SpanId without rand strategy.
func (o *identify) NewEmptySpanId() SpanId {
	s := SpanId{security: true}
	s.bs = [8]byte{0}
	s.str = hex.EncodeToString(s.bs[:])
	return s
}

// NewTraceId
// returns a TraceId with rand strategy.
func (o *identify) NewTraceId() TraceId {
	o.Lock()
	defer o.Unlock()

	t := TraceId{security: true}
	o.random.Read(t.bs[:])

	t.str = hex.EncodeToString(t.bs[:])
	return t
}

func (o *identify) init() *identify {
	o.err = eb.Read(cr.Reader, eb.LittleEndian, &o.data)
	o.random = mr.New(mr.NewSource(o.data))
	return o
}
