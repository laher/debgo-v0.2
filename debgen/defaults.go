/*
   Copyright 2013 Am Laher

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package debgen

import (
	"github.com/laher/debgo-v0.2/deb"
)

// Applies go-specific information to packages.
// Includes dependencies, Go Path information.
func ApplyGoDefaults(pkg *deb.Package) {
	if pkg.ExtraData == nil {
		pkg.ExtraData = map[string]interface{}{}
	}
	gpe := []string{}
	ext, ok := pkg.ExtraData["GoPathExtra"]
	if !ok {
		switch t := ext.(type) {
		case []string:
			gpe = t
		}
	}
	pkg.ExtraData["GoPathExtra"] = append(gpe, GoPathExtraDefault)
	pkg.BuildDepends = deb.BuildDependsDefault
	pkg.Depends = deb.DependsDefault
}
