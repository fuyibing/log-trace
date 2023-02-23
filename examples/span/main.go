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

package main

import (
	"fmt"
	"github.com/fuyibing/log"
	"github.com/fuyibing/log/exports/logs/term"
	"github.com/fuyibing/log/exports/traces/jaeger"
	"github.com/fuyibing/log/globals"
)

func main() {
	log.Collector.SetLogExporter(term.NewExporter())
	log.Collector.SetSpanExporter(jaeger.NewExporter())

	tr := log.Collector.NewTrace("trace")
	sp := tr.NewSpan("span")
	defer sp.End()

	sp.Debug("span debug")
	sp.Info("span info")
	sp.Warn("span warning")
	sp.Error("span error")
	sp.Fatal("span fatal")

	spanIterate(sp, 0)

	s2 := tr.NewSpan("span from trace")
	defer s2.End()

	s2.SetAttribute("author", "fuyibing")
	s2.SetAttribute("email", "websearch@163.com")
	s2.SetAttribute("mobile", "13966013721")
	s2.Fatal("s2 fatal")
}

func spanIterate(s globals.Span, index int) {
	if index > 5 {
		return
	}

	sp := s.NewSpan(fmt.Sprintf("span %d", index))
	defer sp.End()

	sp.Info("span info: %d", index)
	spanIterate(sp, index+1)
}
