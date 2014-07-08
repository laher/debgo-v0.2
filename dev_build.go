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

package debgo

import (
	"github.com/laher/debgo-v0.2/deb"
	"log"
)

// Default build function for Dev packages.
// Implement your own if you prefer
func BuildDevPackageDefault(ddpkg *deb.DevPackage, build *deb.BuildParams) error {
	if ddpkg.BinaryPackage == nil {
		ddpkg.InitBinaryPackage(BuildBinaryArtifactDefault)
	}
	destinationGoPathElement := DEVDEB_GO_PATH_DEFAULT
	goPathRoot := getGoPathElement(build.WorkingDir)
	resources, err := globForSources(goPathRoot, build.WorkingDir, destinationGoPathElement, []string{build.TmpDir, build.DestDir})
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("Resources found: %v", resources)
	}
	for k, v := range resources {
		build.Resources[k] = v
	}
	err = ddpkg.BinaryPackage.Build(build)
	return err
}
