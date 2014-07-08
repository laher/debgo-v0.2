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
)

// *Package is the base unit for this library.
// A *Package contains metadata.
type Package struct {
	Name            string // Package name
	Version         string // Package version
	Description     string // Description
	Maintainer      string // Maintainer
	MaintainerEmail string // Maintainer Email

	AdditionalControlData map[string]string // Other key/values to go into the Control file.

	Architecture string // Supported values: "all", "x386", "amd64", "armhf". TODO: armel

	Depends      string // Depends
	BuildDepends string // BuildDepends is only required for "sourcedebs".

	Priority         string
	StandardsVersion string
	Section          string
	Format           string
	Status           string

	ExtraData map[string]interface{} // Optional for templates
}

// A factory for a Package. Name, Version and Maintainer are mandatory.
func NewPackage(name, version, maintainer string) *Package {
	pkg := new(Package)
	pkg.Name = name
	pkg.Version = version
	pkg.Maintainer = maintainer
	pkg.Priority = PRIORITY_DEFAULT
	pkg.StandardsVersion = STANDARDS_VERSION_DEFAULT
	pkg.Section = SECTION_DEFAULT
	pkg.Format = FORMAT_DEFAULT
	pkg.Status = STATUS_DEFAULT
	return pkg
}

//Resolve architecture(s) and return as a slice
func (pkg *Package) GetArches() ([]Architecture, error) {
	arches, err := resolveArches(pkg.Architecture)
	return arches, err
}

func (pkg *Package) Validate() error {
	if pkg.Name == "" {
		return fmt.Errorf("Name property is required")
	}
	if pkg.Version == "" {
		return fmt.Errorf("Version property is required")
	}
	if pkg.Maintainer == "" {
		return fmt.Errorf("Maintainer property is required")
	}
	return nil
}
