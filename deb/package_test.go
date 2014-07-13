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

package deb_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"testing"
)

func TestCopy(t *testing.T) {
	pkg := deb.NewPackage("a", "1", "me", "desc")
	npkg := deb.Copy(pkg)
	if pkg == npkg {
		t.Errorf("Copy returned the same reference - not a copy")
	}
	if pkg.Name != npkg.Name {
		t.Errorf("Copy didn't copy the same Name value")
	}
	t.Logf("Original: %+v", pkg)
	t.Logf("Copy:     %+v", npkg)
}
