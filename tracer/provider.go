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
	"fmt"
	"github.com/fuyibing/log/config"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	// Provider
	// singleton instance for global provider manager..
	Provider ProviderManager
)

type (
	// ProviderManager
	// is a global provider.
	ProviderManager interface {
		providerGetter
		providerPusher
		providerSetter
	}

	provider struct {
		sync.RWMutex

		attr    Attr
		cancel  context.CancelFunc
		ctx     context.Context
		started bool

		loggerExporter        LoggerExporter
		loggerExporterEnabled bool

		tracerExporter        TracerExporter
		tracerExporterEnabled bool
	}

	providerGetter interface {
		GetAttr() Attr
		NewTrace(name string) Trace
		NewTraceWithContext(ctx context.Context, name string) Trace
		NewTraceWithRequest(name string, request *http.Request) Trace
	}

	providerPusher interface {
		PushBaseLog(level config.LoggerLevel, text string, args ...interface{})
		PushSpan(span Span)
		PushSpanLog(level config.LoggerLevel, text string, args ...interface{})
	}

	providerSetter interface {
		SetAttr(key string, value interface{}) ProviderManager
		SetLoggerExporter(logger LoggerExporter)
		SetTracerExporter(exporter TracerExporter)
		Start(ctx context.Context) error
		Stop() bool
	}
)

// GetAttr
// returns an attribute fields.
func (o *provider) GetAttr() Attr { return o.attr }

// NewTrace
// returns a trace with background context.
func (o *provider) NewTrace(name string) Trace {
	return o.NewTraceWithContext(context.Background(), name)
}

// NewTraceWithContext
// returns a trace with specified context.
func (o *provider) NewTraceWithContext(ctx context.Context, name string) Trace {
	tr := traceCreator(ctx, o, name)

	// Use trace as root.
	tr.useRoot()
	return tr
}

// NewTraceWithRequest
// returns a trace with http request.
func (o *provider) NewTraceWithRequest(name string, req *http.Request) Trace {
	tr := traceCreator(req.Context(), o, name)

	// Use trace base on parent http request.
	tr.useRequest(req)
	return tr
}

// /////////////////////////////////////////////////////////////////////////////
// Provider: pusher
// /////////////////////////////////////////////////////////////////////////////

func (o *provider) PushBaseLog(level config.LoggerLevel, text string, args ...interface{}) {
	log := NewLog(LogInternal, level)
	log.Text = fmt.Sprintf(text, args...)

	// Push
	// log on provider.
	if o.loggerExporterEnabled {
		_ = o.loggerExporter.Push(log)
	}
}

func (o *provider) PushSpan(span Span) {
	if o.tracerExporterEnabled {
		if err := o.tracerExporter.Push(span); err != nil {
			o.PushBaseLog(config.Error, "push span error: %v", err)
		}
	}
}

func (o *provider) PushSpanLog(level config.LoggerLevel, text string, args ...interface{}) {
	log := NewLog(LogSpan, level)
	log.Text = fmt.Sprintf(text, args...)

	// Push
	// log on provider.
	if o.loggerExporterEnabled {
		_ = o.loggerExporter.Push(log)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Provider: setter
// /////////////////////////////////////////////////////////////////////////////

func (o *provider) SetAttr(key string, value interface{}) ProviderManager {
	o.attr.Add(key, value)
	return o
}

func (o *provider) SetLoggerExporter(e LoggerExporter) {
	o.loggerExporter = e
	o.loggerExporterEnabled = e != nil
}

func (o *provider) SetTracerExporter(e TracerExporter) {
	o.tracerExporter = e
	o.tracerExporterEnabled = e != nil
}

func (o *provider) Start(ctx context.Context) error {
	o.Lock()

	// Returns an error
	// if started already.
	if o.started {
		o.Unlock()
		return fmt.Errorf("provider started already")
	}

	// Lock
	// provider status.
	o.ctx, o.cancel = context.WithCancel(ctx)
	o.started = true
	o.initService()
	o.Unlock()

	go func(call, end func()) {
		call()
		end()
	}(o.start, o.startEnd)

	return nil
}

func (o *provider) Stop() bool {
	o.Lock()

	// Wait
	// provider stopped.
	if o.started {
		// Send cancel signal.
		if o.ctx != nil && o.ctx.Err() == nil {
			o.cancel()
		}

		o.Unlock()

		// Waiting
		// stopped state.
		for {
			if func() bool {
				o.RLock()
				defer o.RUnlock()
				return o.started
			}() {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			break
		}
	}

	return true
}

// /////////////////////////////////////////////////////////////////////////////
// Provider: access
// /////////////////////////////////////////////////////////////////////////////

func (o *provider) init() *provider {
	o.attr = Attr{}
	return o.initRuntime()
}

func (o *provider) initRuntime() *provider {
	// Process info.
	//
	// - runtime.pid: process id
	// - runtime.version: go version
	o.attr.Add("runtime.pid", os.Getpid())
	o.attr.Add("runtime.version", runtime.Version())

	// Host name.
	if s, se := os.Hostname(); se == nil {
		o.attr.Add("runtime.host", s)
	}

	// Host addr, IPv4.
	if l, le := net.InterfaceAddrs(); le == nil {
		ls := make([]string, 0)
		for _, la := range l {
			if ipn, ok := la.(*net.IPNet); ok && !ipn.IP.IsLoopback() {
				if ipn.IP.To4() != nil {
					ls = append(ls, ipn.IP.String())
				}
			}
		}
		o.attr.Add("runtime.addr", strings.Join(ls, ", "))
	}

	return o
}

func (o *provider) initService() *provider {
	o.attr.Add("service.name", config.Config.GetServiceName())
	o.attr.Add("service.port", config.Config.GetServicePort())
	o.attr.Add("service.version", config.Config.GetServiceVersion())
	return o
}

// /////////////////////////////////////////////////////////////////////////////
// Provider: blocking exporters
// /////////////////////////////////////////////////////////////////////////////

func (o *provider) debugger(text string, args ...interface{}) {
	if config.Config.DebugOn() {
		_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf(text, args...)+"\n")
	}
}

func (o *provider) start() {
	// Start
	// in 3 coroutines.
	wait := &sync.WaitGroup{}
	wait.Add(3)
	go func() { o.startLogger(); wait.Done() }()
	go func() { o.startProvider(); wait.Done() }()
	go func() { o.startTracer(); wait.Done() }()
	wait.Wait()
}

func (o *provider) startEnd() {
	// Revert
	// provider state.
	o.Lock()
	o.ctx = nil
	o.cancel = nil
	o.started = false
	o.Unlock()

	o.debugger("end process")
}

func (o *provider) startLogger() {
	if o.loggerExporterEnabled {
		if err := o.loggerExporter.Start(o.ctx); err != nil {
			o.debugger("end logger: %v", err)
		} else {
			o.debugger("end logger")
		}
	}
}

func (o *provider) startProvider() {
	for {
		select {
		case <-o.ctx.Done():
			return
		}
	}
}

func (o *provider) startTracer() {
	if o.tracerExporterEnabled {
		if err := o.tracerExporter.Start(o.ctx); err != nil {
			o.debugger("end tracer: %v", err)
		} else {
			o.debugger("end tracer")
		}
	}
}
