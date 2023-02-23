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

package conf

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
		InfoOn() bool
		WarnOn() bool
		ErrorOn() bool
		FatalOn() bool

		GetPid() int

		GetExporter() ExporterConfiguration
		GetJaeger() JaegerReportConfiguration
	}

	ExporterConfiguration interface {
		GetServiceName() string
		GetServiceVersion() string

		GetLogAdapter() string
		GetTraceAdapter() string
	}

	JaegerReportConfiguration interface {
		GetEndpoint() string
		GetUsername() string
		GetPassword() string
	}

	configuration struct {
		Level    Level                        `yaml:"level"`
		Exporter *exporterConfiguration       `yaml:"exporter"`
		Jaeger   *jaegerReporterConfiguration `yaml:"jaeger"`

		processId int

		debugOn,
		infoOn,
		warnOn,
		errorOn,
		fatalOn bool
	}

	exporterConfiguration struct {
		ServiceName    string `yaml:"service-name"`    // example: service name
		ServiceVersion string `yaml:"service-version"` // example: 1.0

		LogAdapter   string `yaml:"log-adapter"`   // example: term
		TraceAdapter string `yaml:"trace-adapter"` // example: jaeger

	}

	jaegerReporterConfiguration struct {
		Endpoint string `yaml:"endpoint"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

// /////////////////////////////////////////////////////////////////////////////
// Configuration: log state
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) DebugOn() bool { return o.debugOn }
func (o *configuration) InfoOn() bool  { return o.infoOn }
func (o *configuration) WarnOn() bool  { return o.warnOn }
func (o *configuration) ErrorOn() bool { return o.errorOn }
func (o *configuration) FatalOn() bool { return o.fatalOn }

// /////////////////////////////////////////////////////////////////////////////
// Configuration: child
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) GetPid() int                          { return o.processId }
func (o *configuration) GetExporter() ExporterConfiguration   { return o.Exporter }
func (o *configuration) GetJaeger() JaegerReportConfiguration { return o.Jaeger }

// /////////////////////////////////////////////////////////////////////////////
// Reporter Configuration
// /////////////////////////////////////////////////////////////////////////////

func (o *exporterConfiguration) GetServiceName() string    { return o.ServiceName }
func (o *exporterConfiguration) GetServiceVersion() string { return o.ServiceVersion }
func (o *exporterConfiguration) GetLogAdapter() string     { return o.LogAdapter }
func (o *exporterConfiguration) GetTraceAdapter() string   { return o.TraceAdapter }

// /////////////////////////////////////////////////////////////////////////////
// Jaeger reporter Configuration
// /////////////////////////////////////////////////////////////////////////////

func (o *jaegerReporterConfiguration) GetEndpoint() string { return o.Endpoint }
func (o *jaegerReporterConfiguration) GetUsername() string { return o.Username }
func (o *jaegerReporterConfiguration) GetPassword() string { return o.Password }

// /////////////////////////////////////////////////////////////////////////////
// Access: defaults
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) defaults() {
	o.state()

	if o.Exporter == nil {
		o.Exporter = &exporterConfiguration{}
	}
	o.Exporter.defaults()

	if o.Jaeger == nil {
		o.Jaeger = &jaegerReporterConfiguration{}
	}
	o.Jaeger.defaults()
}

func (o *exporterConfiguration) defaults() {
}

func (o *jaegerReporterConfiguration) defaults() {
}

// /////////////////////////////////////////////////////////////////////////////
// Access: constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) init() *configuration {
	o.processId = os.Getpid()
	o.scan()
	o.defaults()
	return o
}

func (o *configuration) scan() {
	for _, s := range []string{"config/log.yaml", "../config/log.yaml"} {
		if buf, err := os.ReadFile(s); err == nil {
			if yaml.Unmarshal(buf, o) == nil {
				break
			}
		}
	}
}

func (o *configuration) state() {
	var (
		l  = Level(strings.ToUpper(o.Level.String()))
		n  = l.Int()
		ni = n > 0
	)

	o.debugOn = ni && n >= Debug.Int()
	o.infoOn = ni && n >= Info.Int()
	o.warnOn = ni && n >= Warn.Int()
	o.errorOn = ni && n >= Error.Int()
	o.fatalOn = ni && n >= Fatal.Int()
}
