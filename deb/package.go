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

import ("fmt"
	"errors"
	"os"
)

// *Package is the base unit for this library.
// A *Package contains metadata and some information about building.
type Package struct {
	Name            string // Package name
	Version         string // Package version
	Description     string // Description
	Maintainer      string // Maintainer
	MaintainerEmail string // Maintainer Email

	AdditionalControlData map[string]string // Other key/values to go into the Control file.

	Architecture string // Supported values: "all", "x386", "amd64", "armel"

	Depends string // Depends
	BuildDepends string // BuildDepends is only required for "sourcedebs".

	Priority         string 
	StandardsVersion string
	Section          string
	Format           string
	Status           string

	// Properties below are mainly for build-related properties rather than metadata

	IsVerbose bool // Whether to log debug information
	TmpDir     string // Directory in-which to generate intermediate files & archives
	IsRmtemp   bool // Delete tmp dir after execution?
	DestDir    string // Where to generate .deb files and source debs (.dsc files etc)
	WorkingDir string // This is the root from which to find .go files, templates, resources, etc

	TemplateDir string // Optional. Only required if you're using templates
	Resources map[string]string // Optional. Only if debgo packages your resources automatically. Key is the destination file. Value is the local file

	GoPathExtra []string
}

func resolveArches(arches string) ([]Architecture, error) {
	if arches == "any" || arches == "" {
		return []Architecture{Arch_i386, Arch_armel, Arch_amd64}, nil
	} else if arches == string(Arch_i386) {
		return []Architecture{Arch_i386}, nil
	} else if arches == string(Arch_armel) {
		return []Architecture{Arch_armel}, nil
	} else if arches == string(Arch_amd64) {
		return []Architecture{Arch_amd64}, nil
	}
	return nil, fmt.Errorf("Architecture %s not supported", arches)
}

//Resolve architecture(s) and return as a slice
func (pkg *Package) GetArches() ([]Architecture, error) {
	arches, err := resolveArches(pkg.Architecture)
	return arches, err
}

//Initialise build process (make Temp and Dest directories)
func (pkg *Package) Init() error {
	//make tmpDir
	if pkg.TmpDir == "" {
		return errors.New("Temp directory not specified")
	}
	err := os.MkdirAll(pkg.TmpDir, 0755)
	if err != nil {
		return err
	}
	//make destDir
	if pkg.DestDir == "" {
		return errors.New("Destination directory not specified")
	}
	err = os.MkdirAll(pkg.DestDir, 0755)
	if err != nil {
		return err
	}
	return err
}

func (pkg *Package) NewTemplateData() TemplateData {
	templateVars := newTemplateData(pkg.Name, pkg.Version, pkg.Maintainer, pkg.MaintainerEmail, pkg.Version, pkg.Architecture, pkg.Description, pkg.Depends, pkg.BuildDepends, pkg.Priority, pkg.Status, pkg.StandardsVersion, pkg.Section, pkg.Format, pkg.GoPathExtra, pkg.AdditionalControlData)
	return templateVars
}

// A factory for a Package. Name, Version and Maintainer are mandatory.
func NewPackage(name, version, maintainer string) *Package {
	pkg := new(Package)
	pkg.Name = name
	pkg.Version = version
	pkg.Maintainer = maintainer

	pkg.TmpDir = TEMP_DIR_DEFAULT
	pkg.DestDir = DIST_DIR_DEFAULT
	pkg.IsRmtemp = true
	pkg.WorkingDir = WORKING_DIR_DEFAULT

	pkg.GoPathExtra = []string{GO_PATH_EXTRA_DEFAULT}

	pkg.BuildDepends = BUILD_DEPENDS_DEFAULT
	pkg.Priority = PRIORITY_DEFAULT
	pkg.StandardsVersion = STANDARDS_VERSION_DEFAULT
	pkg.Section = SECTION_DEFAULT
	pkg.Format = FORMAT_DEFAULT
	pkg.Status = STATUS_DEFAULT
	pkg.Depends = DEPENDS_DEFAULT
	return pkg
}
