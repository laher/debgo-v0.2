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

package deb

import (
	"fmt"
//	"path/filepath"
)

// BuildDevPackageFunc specifies a function which can build a DevPackage
type BuildDevPackageFunc func(*DevPackage, *BuildParams) error

// *DevPackage builds a sources-only '-dev' package, which can be used as a BuildDepends dependency.
// For pure Go packages, this can be cross-platform (architecture == 'all'), but in some cases it might need to be architecture specific
type DevPackage struct {
	*Package
	//DebFilePath   string
	BinaryPackage *BinaryPackage
	BuildFunc     BuildDevPackageFunc
}

// Factory for DevPackage
func NewDevPackage(pkg *Package, buildFunc BuildDevPackageFunc) *DevPackage {
	//debPath := filepath.Join(pkg.DestDir, pkg.Name+"-dev_"+pkg.Version+".deb")
	ddpkg := &DevPackage{ Package: pkg,
	//	DebFilePath: debPath,
		BuildFunc: buildFunc}
	return ddpkg
}

func (ddpkg *DevPackage) InitBinaryPackage(buildFunc BuildBinaryArtifactFunc) {
	if ddpkg.BinaryPackage == nil {
		//TODO *complete* copy of properties, using reflection?
		devpkg := NewPackage(ddpkg.Name+"-dev", ddpkg.Version, ddpkg.Maintainer)
		devpkg.Description = ddpkg.Description
		devpkg.MaintainerEmail = ddpkg.MaintainerEmail
		devpkg.AdditionalControlData = ddpkg.AdditionalControlData
		devpkg.Architecture = "all"
		//devpkg.IsVerbose = ddpkg.IsVerbose
		//devpkg.IsRmtemp = ddpkg.IsRmtemp
		ddpkg.BinaryPackage = NewBinaryPackage(devpkg, buildFunc)
	}
/*
	if ddpkg.BinaryPackage.Resources == nil {
		ddpkg.BinaryPackage.Resources = map[string]string{}
	}
*/
}

// Invokes the BuildFunc
func (ddpkg *DevPackage) Build(build *BuildParams) error {
	if ddpkg.BuildFunc == nil {
		return fmt.Errorf("No build function provided (*DevPackage.BuildFunc)")
	}
	return ddpkg.BuildFunc(ddpkg, build)
}

