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
	"encoding/base64"
	"fmt"
	"github.com/fuyibing/log/conf"
	"github.com/fuyibing/log/globals"
	"github.com/valyala/fasthttp"

	"net/http"
)

type (
	// Exporter
	// push list span into jaeger server.
	Exporter struct {
		spans     []globals.SpanReadonly
		formatter Formatter
	}
)

// NewExporter
// is a exporter manager.
func NewExporter() *Exporter {
	return (&Exporter{}).init()
}

// Push
// add span into memory.
func (o *Exporter) Push(spans ...globals.SpanReadonly) (err error) {
	return o.push(spans...)
}

func (o *Exporter) Push2(spans ...globals.SpanReadonly) (err error) {
	o.spans = append(o.spans, spans...)

	if count := len(o.spans); count == 8 {
		var buf *bytes.Buffer

		if buf, err = o.formatter.Generate(o.spans); err != nil {
			return
		}

		err = o.upload(buf)
		return
	}

	return
}

func (o *Exporter) push(spans ...globals.SpanReadonly) (err error) {
	var buf *bytes.Buffer

	if buf, err = o.formatter.Generate(spans); err != nil {
		return
	}

	err = o.upload(buf)
	return
}

func (o *Exporter) upload(buf *bytes.Buffer) (err error) {
	var (
		req = fasthttp.AcquireRequest()
		res = fasthttp.AcquireResponse()
	)

	defer func() {
		fasthttp.ReleaseResponse(res)
		fasthttp.ReleaseRequest(req)
	}()

	req = fasthttp.AcquireRequest()
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/vnd.apache.thrift.binary")
	req.SetRequestURI(conf.Config.GetJaeger().GetEndpoint())
	req.SetBodyStream(buf, buf.Len())

	// Basic authorization.
	if user, pass := conf.Config.GetJaeger().GetUsername(), conf.Config.GetJaeger().GetPassword(); user != "" {
		req.Header.Set("Authorization",
			fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(user+":"+pass))),
		)
	}

	err = fasthttp.Do(req, res)
	return
}

// /////////////////////////////////////////////////////////////////////////////
// Collector: access
// /////////////////////////////////////////////////////////////////////////////

func (o *Exporter) init() *Exporter {
	o.formatter = (&formatter{}).init()
	o.spans = make([]globals.SpanReadonly, 0)
	return o
}
