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
	"testing"
)

func TestConfiguration_With(t *testing.T) {
	Config.With(
		ServiceName("my app"),
		ServicePort(3721),
		ServiceVersion("2.3.4"),
	)

	buf, _ := yaml.Marshal(Config)
	// buf, _ := json.Marshal(Config)
	t.Logf("config: %s", buf)
}
