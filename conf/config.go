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

		GetExporter() ExporterConfiguration
	}

	ExporterConfiguration interface {
		GetServiceName() string

		GetLogAdapter() string
		GetTraceAdapter() string
	}

	configuration struct {
		Level    Level                  `yaml:"level"`
		Exporter *exporterConfiguration `yaml:"exporter"`

		debugOn, infoOn, warnOn, errorOn, fatalOn bool
	}

	exporterConfiguration struct {
		ServiceName string `yaml:"service-name"` // example: service name

		LogAdapter   string `yaml:"log-adapter"`   // example: term
		TraceAdapter string `yaml:"trace-adapter"` // example: jaeger
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

func (o *configuration) GetExporter() ExporterConfiguration { return o.Exporter }

// /////////////////////////////////////////////////////////////////////////////
// Reporter Configuration
// /////////////////////////////////////////////////////////////////////////////

func (o *exporterConfiguration) GetServiceName() string { return o.ServiceName }

func (o *exporterConfiguration) GetLogAdapter() string   { return o.LogAdapter }
func (o *exporterConfiguration) GetTraceAdapter() string { return o.TraceAdapter }

// /////////////////////////////////////////////////////////////////////////////
// Access: defaults
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) defaults() {
	o.state()

	if o.Exporter == nil {
		o.Exporter = &exporterConfiguration{}
	}
	o.Exporter.defaults()
}

func (o *exporterConfiguration) defaults() {
	// TODO force definitions.
	o.ServiceName = "log-trace"
}

// /////////////////////////////////////////////////////////////////////////////
// Access: constructor
// /////////////////////////////////////////////////////////////////////////////

func (o *configuration) init() *configuration {
	o.scan()
	o.defaults()
	return o
}

func (o *configuration) scan() {
}

func (o *configuration) state() {
	o.debugOn = true
	o.infoOn = true
	o.warnOn = true
	o.errorOn = true
	o.fatalOn = true
}
