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
	"context"
	"github.com/fuyibing/log"
	"github.com/fuyibing/log/exporters/logger_term"
	"github.com/fuyibing/log/exporters/tracer_term"
)

func init() {

	log.Provider.SetLoggerExporter(logger_term.New())
	log.Provider.SetTracerExporter(tracer_term.New())
	// log.Provider.SetTracerExporter(tracer_jaeger.New())

	ctx := context.Background()
	_ = log.Provider.Start(ctx)
}

func main() {
	ch := make(chan bool)

	go func() {
		// time.Sleep(time.Second * 5)
		log.Provider.Stop()
		ch <- true
	}()

	go span()

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	log.Fatal("fatal")

	for {
		select {
		case <-ch:
			return
		}
	}
}

func span() {
	tr := log.Provider.NewTrace("trace")

	s1 := tr.NewSpan("span 1")
	defer s1.End()

	s1.Info("span 1: info")
	s1.Fatal("span 1: fatal")

}
