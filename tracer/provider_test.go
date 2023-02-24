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
	"testing"
)

func TestProvider_NewTrace(t *testing.T) {
	t.Logf("provider attr: %v", Provider.GetAttr().JSON())

	tr := Provider.NewTrace("example")
	t.Logf("Trace: %s - [%s][err=%v][sec=%v]", tr.GetName(), tr.GetTraceId().String(), tr.GetTraceId().Err(), tr.GetTraceId().Security())
	t.Logf("Span: [%s][err=%v][sec=%v]", tr.GetSpanId().String(), tr.GetSpanId().Err(), tr.GetSpanId().Security())

	s1 := tr.NewSpan("span 1")
	s1.SetAttr("key", "value")
	t.Logf("attr: %v", s1.GetAttr().JSON())
}
