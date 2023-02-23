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

package jaeger

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/fuyibing/log/conf"
	"github.com/fuyibing/log/exports/traces/jaeger/thrift"
	"github.com/fuyibing/log/globals"
	"strconv"
)

type (
	Formatter interface {
		Generate(spans []globals.SpanReadonly) (res *bytes.Buffer, err error)
	}

	// Formatter
	// generate span as thrift format body.
	formatter struct {
	}
)

func (o *formatter) Generate(gs []globals.SpanReadonly) (res *bytes.Buffer, err error) {
	var (
		batch = &Batch{
			Process: &Process{
				ServiceName: conf.Config.GetExporter().GetServiceName(),
			},
			Spans: make([]*Span, 0),
		}
	)

	// Process/Resource fields.
	if tags := o.makeTags(map[string]interface{}{
		"service.name":    conf.Config.GetExporter().GetServiceName(),
		"service.version": conf.Config.GetExporter().GetServiceVersion(),
		"process.id":      conf.Config.GetPid(),
	}); len(tags) > 0 {
		batch.Process.Tags = tags
	}

	// Span lists.
	for _, s := range gs {
		batch.Spans = append(batch.Spans, o.build(s))
	}

	return o.serialize(batch)
}

// /////////////////////////////////////////////////////////////////////////////
// Formatter: access and initialize
// /////////////////////////////////////////////////////////////////////////////

func (o *formatter) build(s globals.SpanReadonly) (span *Span) {
	var (
		tid = s.GetTraceId()
		sid = s.GetSpanId()
		pid = s.GetParentSpanId()
	)

	// Basic fields.
	span = &Span{}
	span.TraceIdHigh = int64(binary.BigEndian.Uint64(tid[0:8]))
	span.TraceIdLow = int64(binary.BigEndian.Uint64(tid[8:16]))
	span.SpanId = int64(binary.BigEndian.Uint64(sid[:]))
	span.ParentSpanId = int64(binary.BigEndian.Uint64(pid[:]))
	span.OperationName = s.GetName()
	span.Flags = 1
	span.StartTime = s.GetStartTime().UnixMicro()
	span.Duration = s.GetDuration().Microseconds()

	// Span tags.
	if tags := o.makeTags(s.GetAttributes()); len(tags) > 0 {
		span.Tags = tags
	}

	// Span logs.
	if logs := o.makeLogs(s.GetLogs()); len(logs) > 0 {
		span.Logs = logs
	}

	return
}

func (o *formatter) init() *formatter {
	return o
}

func (o *formatter) serialize(data thrift.TStruct) (res *bytes.Buffer, err error) {
	var (
		buf = thrift.NewTMemoryBuffer()
	)

	if err = data.Write(context.Background(), thrift.NewTBinaryProtocolConf(buf, &thrift.TConfiguration{})); err != nil {
		return
	}

	res = buf.Buffer
	return
}

func (o *formatter) makeLogs(lists []*globals.Log) (logs []*Log) {
	logs = make([]*Log, 0)
	for _, item := range lists {
		logs = append(logs, &Log{
			Timestamp: item.Time.UnixMicro(),
			Fields: o.makeTags(map[string]interface{}{
				item.Level.String(): fmt.Sprintf("[%-26s] %s",
					item.Time.Format("2006-01-02 15:04:05.999999"),
					item.Text,
				),
			}),
		})
	}
	return
}

func (o *formatter) makeTags(mapper map[string]interface{}) (tags []*Tag) {
	tags = make([]*Tag, 0)
	for k, v := range mapper {
		switch x := v.(type) {
		case bool:
			tags = append(tags, &Tag{Key: k, VType: TagType_BOOL, VBool: &x})
		case float64, float32:
			if n, ne := strconv.ParseFloat(fmt.Sprintf("%v", v), 64); ne == nil {
				tags = append(tags, &Tag{Key: k, VType: TagType_DOUBLE, VDouble: &n})
			}
		case int, int8, int16, int32, int64, uint, uint8, uint32, uint64:
			if n, ne := strconv.ParseInt(fmt.Sprintf("%v", v), 0, 64); ne == nil {
				tags = append(tags, &Tag{Key: k, VType: TagType_LONG, VLong: &n})
			}
		case string:
			tags = append(tags, &Tag{Key: k, VType: TagType_STRING, VStr: &x})
		}
	}
	return
}
