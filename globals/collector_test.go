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

package globals

import (
	"fmt"
	"github.com/fuyibing/log/conf"
	"testing"
)

func TestCollector_New(t *testing.T) {
	tr := Collector.NewTrace("new trace " + conf.Config.GetExporter().GetServiceName())
	t.Logf("[tid=%s][sid=%v] %s", tr.GetTraceId(), tr.GetSpanId(), tr.GetName())

	sp := tr.NewSpan("top span")
	defer sp.End()

	spanIterate(sp, 0)

	// r1 := tr.NewSpan("root span 1")
	// defer r1.End()
	//
	// t.Logf("[tid=%s][pid=%s][sid=%s] %s", r1.GetTraceId(), r1.GetParentSpanId(), r1.GetSpanId(), r1.GetName())
	//
	// r2 := r1.NewSpan("child span 1")
	// defer r2.End()
	//
	// t.Logf("[tid=%s][pid=%s][sid=%s] %s", r2.GetTraceId(), r2.GetParentSpanId(), r2.GetSpanId(), r2.GetName())
}

func spanIterate(s Span, index int) {
	if index > 5 {
		return
	}

	sp := s.NewSpan(fmt.Sprintf("span %d", index))
	defer sp.End()

	sp.Info("span info: %d", index)
	spanIterate(sp, index+1)
}
