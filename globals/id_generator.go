// author: wsfuyibing <websearch@163.com>
// date: 2023-02-22

package globals

import (
	cr "crypto/rand"
	eb "encoding/binary"
	mr "math/rand"
	"sync"
)

var (
	Generator IdGenerator
)

type (
	IdGenerator interface {
		SpanId() SpanId
		TraceId() TraceId
	}

	idg struct {
		sync.Mutex
		data   int64
		err    error
		random *mr.Rand
	}
)

func (o *idg) SpanId() SpanId {
	o.Lock()
	defer o.Unlock()

	s := SpanId{}
	o.random.Read(s[:])
	return s
}

func (o *idg) TraceId() TraceId {
	o.Lock()
	defer o.Unlock()

	s := TraceId{}
	o.random.Read(s[:])
	return s
}

func (o *idg) init() *idg {
	o.err = eb.Read(cr.Reader, eb.LittleEndian, &o.data)
	o.random = mr.New(mr.NewSource(o.data))
	return o
}
