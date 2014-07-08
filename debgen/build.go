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

// A factory for a Go Package. Includes dependencies and Go Path information
func NewGoPackage(name, version, maintainer string) *deb.Package {
	pkg := deb.NewPackage(name, version, maintainer)
	pkg.ExtraData = map[string]interface{}{
		"GoPathExtra": []string{GO_PATH_EXTRA_DEFAULT}}
	pkg.BuildDepends = deb.BUILD_DEPENDS_DEFAULT
	pkg.Depends = deb.DEPENDS_DEFAULT
	return pkg
}
