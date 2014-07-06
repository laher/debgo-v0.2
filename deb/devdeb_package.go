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
	"log"
	"path/filepath"
)

// BuildDevPackageFunc specifies a function which can build a DevPackage
type BuildDevPackageFunc func(*DevPackage) error

// *DevPackage builds a sources-only '-dev' package, which can be used as a BuildDepends dependency.
// For pure Go packages, this can be cross-platform (architecture == 'all'), but in some cases it might need to be architecture specific
type DevPackage struct {
	*Package
	DebFilePath   string
	BinaryPackage *BinaryPackage
	BuildFunc     BuildDevPackageFunc
}

// Factory for DevPackage
func NewDevPackage(pkg *Package) *DevPackage {
	debPath := filepath.Join(pkg.DestDir, pkg.Name+"-dev_"+pkg.Version+".deb")
	return &DevPackage{Package: pkg,
		DebFilePath: debPath,
		BuildFunc:   BuildDefault}
}

func (ddpkg *DevPackage) InitBinaryPackage() {
	if ddpkg.BinaryPackage == nil {
		//TODO *complete* copy of properties, using reflection?
		devpkg := NewPackage(ddpkg.Name+"-dev", ddpkg.Version, ddpkg.Maintainer)
		devpkg.Description = ddpkg.Description
		devpkg.MaintainerEmail = ddpkg.MaintainerEmail
		devpkg.AdditionalControlData = ddpkg.AdditionalControlData
		devpkg.Architecture = "all"
		devpkg.IsVerbose = ddpkg.IsVerbose
		devpkg.IsRmtemp = ddpkg.IsRmtemp
		ddpkg.BinaryPackage = NewBinaryPackage(devpkg)
	}
	if ddpkg.BinaryPackage.Resources == nil {
		ddpkg.BinaryPackage.Resources = map[string]string{}
	}
}

// Invokes the BuildFunc
func (ddpkg *DevPackage) Build() error {
	if ddpkg.BuildFunc == nil {
		return fmt.Errorf("No build function provided (*DevPackage.BuildFunc)")
	}
	return ddpkg.BuildFunc(ddpkg)
}

// Default build function for Dev packages.
// Implement your own if you prefer
func BuildDefault(ddpkg *DevPackage) error {
	if ddpkg.BinaryPackage == nil {
		ddpkg.InitBinaryPackage()
	}
	destinationGoPathElement := DEVDEB_GO_PATH_DEFAULT
	goPathRoot := getGoPathElement(ddpkg.WorkingDir)
	resources, err := globForSources(goPathRoot, ddpkg.WorkingDir, destinationGoPathElement, []string{ddpkg.TmpDir, ddpkg.DestDir})
	if err != nil {
		return err
	}
	if ddpkg.IsVerbose {
		log.Printf("Resources found: %v", resources)
	}
	for k, v := range resources {
		ddpkg.BinaryPackage.Resources[k] = v
	}
	err = ddpkg.BinaryPackage.Build()
	return err
}
