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

package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var (
	Config Configuration
)

type (
	Configuration interface {
		DebugOn() bool
		ErrorOn() bool
		FatalOn() bool
		GetJaegerTrace() JaegerTraceConfiguration
		GetLoggerLevel() LoggerLevel
		GetLoggerName() LoggerName
		GetOpenTracingSample() string
		GetOpenTracingSpanId() string
		GetOpenTracingTraceId() string
		GetServiceName() string
		GetServicePort() int
		GetServiceVersion() string
		GetTracerName() TracerName
		GetTracerTopic() string
		GetTracerWithLog() bool
		InfoOn() bool
		SetLoggerLevel(level LoggerLevel)
		SetLoggerName(name LoggerName)
		SetTracerName(name TracerName)
		WarnOn() bool
		With(opts ...Option)
	}

	JaegerTraceConfiguration interface {
		GetEndpoint() string
		GetUsername() string
		GetPassword() string
	}

	configuration struct {
		OpenTracingSample  string `yaml:"open-tracing-sample"`
		OpenTracingSpanId  string `yaml:"open-tracing-span-id"`
		OpenTracingTraceId string `yaml:"open-tracing-trace-id"`
		ServiceName        string `yaml:"service-name"`
		ServicePort        int    `yaml:"service-port"`
		ServiceVersion     string `yaml:"service-version"`

		LoggerLevel LoggerLevel `yaml:"logger-level"`
		LoggerName  LoggerName  `yaml:"logger-name"`

		// TracerName
		// config trace exporter name.
		TracerName TracerName `yaml:"tracer-name"`

		// TracerTopic
		// name for trace storage.
		TracerTopic string `yaml:"tracer-topic"`

		// TracerWithLog
		// whether to join the log when reporting Trace.
		TracerWithLog bool `yaml:"tracer-with-log"`

		JaegerTrace *jaegerTraceConfiguration `yaml:"jaeger-trace"`

		debugOn, infoOn, warnOn, errorOn, fatalOn bool
	}

	jaegerTraceConfiguration struct {
		Endpoint string `yaml:"endpoint"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Configuration: open tracing
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) DebugOn() bool                            { return o.debugOn }
func (o *configuration) ErrorOn() bool                            { return o.errorOn }
func (o *configuration) FatalOn() bool                            { return o.fatalOn }
func (o *configuration) GetJaegerTrace() JaegerTraceConfiguration { return o.JaegerTrace }
func (o *configuration) GetLoggerLevel() LoggerLevel              { return o.LoggerLevel }
func (o *configuration) GetLoggerName() LoggerName                { return o.LoggerName }
func (o *configuration) GetOpenTracingSample() string             { return o.OpenTracingSample }
func (o *configuration) GetOpenTracingSpanId() string             { return o.OpenTracingSpanId }
func (o *configuration) GetOpenTracingTraceId() string            { return o.OpenTracingTraceId }
func (o *configuration) GetServiceName() string                   { return o.ServiceName }
func (o *configuration) GetServicePort() int                      { return o.ServicePort }
func (o *configuration) GetServiceVersion() string                { return o.ServiceVersion }
func (o *configuration) GetTracerName() TracerName                { return o.TracerName }
func (o *configuration) GetTracerTopic() string                   { return o.TracerTopic }
func (o *configuration) GetTracerWithLog() bool                   { return o.TracerWithLog }
func (o *configuration) InfoOn() bool                             { return o.infoOn }
func (o *configuration) SetLoggerLevel(level LoggerLevel)         { o.LoggerLevel = level; o.resetState() }
func (o *configuration) SetLoggerName(name LoggerName)            { o.LoggerName = name }
func (o *configuration) SetTracerName(name TracerName)            { o.TracerName = name }
func (o *configuration) WarnOn() bool                             { return o.warnOn }

func (o *configuration) With(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Jaeger Trace Configuration
// /////////////////////////////////////////////////////////////////////////////

func (o *jaegerTraceConfiguration) GetEndpoint() string { return o.Endpoint }
func (o *jaegerTraceConfiguration) GetUsername() string { return o.Username }
func (o *jaegerTraceConfiguration) GetPassword() string { return o.Password }

// /////////////////////////////////////////////////////////////////////////////
// Access: initialize
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) scan() {
	for _, s := range []string{"config/log.yaml", "../config/log.yaml"} {
		if buf, err := os.ReadFile(s); err == nil {
			if err = yaml.Unmarshal(buf, o); err == nil {
				return
			}
		}
	}
}

func (o *configuration) init() *configuration {
	o.scan()
	o.initDefaults()
	o.initChildren()
	return o
}

func (o *configuration) initDefaults() {
	// OpenTracing.
	// for: OpenTelemetry
	if o.OpenTracingSample == "" {
		o.OpenTracingSample = DefaultOpenTracingSample
	}
	if o.OpenTracingSpanId == "" {
		o.OpenTracingSpanId = DefaultOpenTracingSpanId
	}
	if o.OpenTracingTraceId == "" {
		o.OpenTracingTraceId = DefaultOpenTracingTraceId
	}

	// Default topic name.
	if o.TracerTopic == "" {
		o.TracerTopic = "log-trace"
	}

	// Init
	// log level state.
	if level := LoggerLevel(strings.ToUpper(o.LoggerLevel.String())); level.Int() >= Off.Int() {
		o.SetLoggerLevel(level)
	} else {
		o.SetLoggerLevel(LevelDefault)
	}
}

func (o *configuration) initChildren() {
	if o.JaegerTrace == nil {
		o.JaegerTrace = &jaegerTraceConfiguration{}
	}
	o.JaegerTrace.initDefaults()
}

func (o *configuration) resetState() {
	// Level compare.
	i := o.LoggerLevel.Int()
	is := i > Off.Int()

	// Level state.
	o.debugOn = is && i >= Debug.Int()
	o.infoOn = is && i >= Info.Int()
	o.warnOn = is && i >= Warn.Int()
	o.errorOn = is && i >= Error.Int()
	o.fatalOn = is && i >= Fatal.Int()
}

// /////////////////////////////////////////////////////////////////////////////
// Access: initialize
// /////////////////////////////////////////////////////////////////////////////

func (o *jaegerTraceConfiguration) initDefaults() {}
