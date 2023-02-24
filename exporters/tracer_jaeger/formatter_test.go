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

package tracer_jaeger

import (
	"github.com/fuyibing/log"
	"testing"
)

func TestExporter_Upload(t *testing.T) {

	tr := log.Provider.NewTrace("trace")
	tr.SetAttr("trace-key", "attribute value of trace")

	s1 := tr.NewSpan("span")
	defer s1.End()

	s1.Debug("span 1. debug")
	s1.Info("span 1. info")

	s2 := s1.NewSpan("span 2 of span 1")
	defer s2.End()
	s2.SetAttr("span-key-2", "attribute value of span 2")
	s2.Debug("span 2. debug")
	s2.Info("span 2. info")

	ex := &formatter{}

	buf, err := ex.Thrift(s1, s2)
	if err != nil {
		t.Errorf("generate error: %v", err)
	}

	Exp := (&exporter{}).init()

	err = Exp.Upload(buf)
	t.Logf("batch: %v", err)
}
