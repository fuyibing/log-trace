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

package main

import (
	"github.com/fuyibing/log"
	"github.com/fuyibing/log/exporters/logger_term"
	"github.com/fuyibing/log/exporters/tracer_jaeger"
	"time"
)

func init() {
	initLogger()
	initTrace()
}

func initLogger() {
	v := logger_term.New()
	log.Provider.SetLoggerExporter(v)
}

func initTrace() {
	v := tracer_jaeger.New()
	log.Provider.SetTracerExporter(v)
}

func main() {
	mainBase()
	mainTrace()
}

func mainBase() {
	log.Debug("base debug")
	log.Info("base info")
	log.Warn("base warn")
	log.Error("base error")
	log.Fatal("base fatal")
}

func mainTrace() {
	tr := log.Provider.NewTrace("trace")

	p1 := tr.NewSpan("span1")
	defer p1.End()

	p1.Error("error on span 1")
	p1.Info("End span 1")

	p2 := p1.NewSpan("span1")
	defer p2.End()

	p2.Info("info on span 2")
	time.Sleep(time.Millisecond * 10)

}
