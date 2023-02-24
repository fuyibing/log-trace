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

package tracer_jaeger

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/fuyibing/log/config"
	"github.com/fuyibing/log/tracer"
	"github.com/valyala/fasthttp"
	"net/http"
)

type (
	Exporter interface {
		Push(span tracer.Span) error
		Start(ctx context.Context) error
		Stopped() bool
	}

	exporter struct {
		formatter Formatter
		endpoint  string
		username  string
		password  string
	}
)

func New() Exporter {
	return (&exporter{}).init()
}

func (o *exporter) Push(span tracer.Span) (err error) {
	var buf *bytes.Buffer
	if buf, err = o.formatter.Thrift(span); err == nil {
		err = o.Upload(buf)
	}
	return
}

func (o *exporter) Start(ctx context.Context) error {
	return nil
}

func (o *exporter) Stopped() bool { return true }

func (o *exporter) Upload(buf *bytes.Buffer) (err error) {
	var (
		req = fasthttp.AcquireRequest()
		res = fasthttp.AcquireResponse()
	)

	req.SetRequestURI(o.endpoint)
	req.SetBodyStream(buf, buf.Len())
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/x-thrift")

	// Bind authorization.
	if o.username != "" {
		req.Header.Set("Authorization",
			fmt.Sprintf("Basic %s",
				base64.StdEncoding.EncodeToString([]byte(o.username+":"+o.password)),
			),
		)
	}

	// Send request.
	err = fasthttp.Do(req, res)
	return
}

// /////////////////////////////////////////////////////////////////////////////
// Exporter: access
// /////////////////////////////////////////////////////////////////////////////

func (o *exporter) init() *exporter {
	o.formatter = (&formatter{}).init()

	o.endpoint = config.Config.GetJaegerTrace().GetEndpoint()
	o.password = config.Config.GetJaegerTrace().GetPassword()
	o.username = config.Config.GetJaegerTrace().GetUsername()
	return o
}
