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
	"fmt"
)
/*
// Generate artifacts for Go-specific packages
func GenGoDevArtifact(ddpkg *deb.DevPackage, build *deb.BuildParams, sourcesDir string) error {
	destinationGoPathElement := DEVDEB_GO_PATH_DEFAULT
	sourcesRelativeTo := GetGoPathElement(sourcesDir)
	return GenDevArtifact(ddpkg, build, sourcesDir, sourcesRelativeTo, destinationGoPathElement)

}
*/

// Default build function for Dev packages.
// Implement your own if you prefer
func GenDevArtifact(ddpkg *deb.DevPackage, build *deb.BuildParams) error {
	if ddpkg.BinaryPackage == nil {
		ddpkg.InitBinaryPackage()
	}
	artifacts, err := ddpkg.BinaryPackage.GetArtifacts()
	if err != nil {
		return err
	}
	for arch, artifact := range artifacts {
		err = GenBinaryArtifact(artifact, build)
		if err != nil {
			return fmt.Errorf("Error building for '%s': %v", arch, err)
		}
	}
	return err
}
