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

func TestIdentify_NewTraceId(t *testing.T) {
	t1 := Identify.NewTraceId()
	t.Logf("trace 1 id: %v", t1.Byte())
	t.Logf("trace 1 str: %s", t1.String())

	t2 := Identify.HexTraceId(t1.String())
	t.Logf("trace 2 id: %v", t2.Byte())
	t.Logf("trace 2 str: %s", t2.String())
}

func TestIdentify_NewSpanId(t *testing.T) {
	t1 := Identify.NewSpanId()
	t.Logf("span 1 id: %v", t1.Byte())
	t.Logf("span 1 str: %s", t1.String())
	t.Logf("span 1 error: %v", t1.Err())
	t.Logf("span 1 security: %v", t1.Security())

	t2 := Identify.HexSpanId(t1.String())
	t.Logf("span 2 id: %v", t2.Byte())
	t.Logf("span 2 str: %s", t2.String())
	t.Logf("span 2 error: %v", t2.Err())
	t.Logf("span 2 security: %v", t2.Security())
}
