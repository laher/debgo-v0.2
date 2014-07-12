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
	"fmt"
	"github.com/laher/debgo-v0.2/deb"
)

// Default build function for Dev packages.
// Implement your own if you prefer
func GenDevArtifact(ddpkg *deb.Package, build *deb.BuildParams) error {
	artifacts, err := deb.GetDebs(ddpkg)
	if err != nil {
		return err
	}
	for arch, artifact := range artifacts {
		err = GenDeb(artifact, build)
		if err != nil {
			return fmt.Errorf("Error building for '%s': %v", arch, err)
		}
	}
	return err
}
