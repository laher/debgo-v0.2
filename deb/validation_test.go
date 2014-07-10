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
	"log"
	"testing"
)

var (
	validVersions = []string{"1.2.3a", "123-x"}
	badVersions   = []string{"12!3", "a123-x"}
)

func ExampleValidateVersion() {
	v := "1.0.1-git123"
	err := deb.ValidateVersion(v)
	if err != nil {
		log.Fatalf("Version validation broken for %v", v)
	}

}

func TestValidateVersion(t *testing.T) {
	for _, v := range validVersions {
		err := deb.ValidateVersion(v)
		if err != nil {
			t.Fatalf("Version validation broken for %v", v)
		}
	}
	for _, v := range badVersions {
		err := deb.ValidateVersion(v)
		if err == nil {
			t.Fatalf("Bad Version not detected for %v", v)
		}
	}
}
