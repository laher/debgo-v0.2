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

import "fmt"

// a package contains metadata only
type Package struct {
	Name            string
	Version         string
	Description     string
	Maintainer      string
	MaintainerEmail string

	AdditionalControlData        map[string]string

	Architecture string

//	ExecutablePaths map[string][]string
//	OtherFiles      map[string]string

	IsVerbose bool

	Depends      string

	//only required for sourcedebs
	BuildDepends string


	Priority string
	StandardsVersion string
	Section string
	Format string
	Status string

	TemplateDir  string

	IsRmtemp   bool
	TmpDir     string
	DestDir    string
	WorkingDir string
}

func resolveArches(arches string) ([]string, error) {
	if arches == "any" || arches == "" {
		return []string{"i386", "armel", "amd64"}, nil
	}
	if arches != "i386" && arches != "armel" && arches != "amd64" {
		return nil, fmt.Errorf("Architecture %s not supported", arches)
	}
	return []string{arches}, nil
}

//Resolve architecture(s) and return as a slice
func (pkg *Package) GetArches() ([]string, error) {
	arches, err := resolveArches(pkg.Architecture)
	return arches, err
}


func (pkg *Package) NewTemplateData() TemplateData {
	templateVars := newTemplateData(pkg.Name, pkg.Version, pkg.Maintainer, pkg.MaintainerEmail, pkg.Version, pkg.Architecture, pkg.Description, pkg.Depends, pkg.BuildDepends, pkg.Priority, pkg.Status, pkg.StandardsVersion, pkg.Section, pkg.Format, pkg.AdditionalControlData)
	return templateVars
}


// A factory for a Package. Name, Version and Maintainer are mandatory.
func NewPackage(name, version, maintainer string) *Package {
	pkg := new(Package)
	pkg.Name = name
	pkg.Version = version
	pkg.Maintainer = maintainer

	pkg.TmpDir = "_test/tmp"
	pkg.DestDir = "_test/dist"
	pkg.IsRmtemp = true
	pkg.WorkingDir = "."

	pkg.BuildDepends = BUILD_DEPENDS_DEFAULT
	pkg.Priority = PRIORITY_DEFAULT
	pkg.StandardsVersion = STANDARDS_VERSION_DEFAULT
	pkg.Section = SECTION_DEFAULT
	pkg.Format = FORMAT_DEFAULT
	pkg.Status = STATUS_DEFAULT
	pkg.Depends = DEPENDS_DEFAULT
	return pkg
}
