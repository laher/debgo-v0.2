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
	"errors"
	"fmt"
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

	Depends      string // Depends
	BuildDepends string // BuildDepends is only required for "sourcedebs".

	Priority         string
	StandardsVersion string
	Section          string
	Format           string
	Status           string

	// Properties below are mainly for build-related properties rather than metadata

	IsVerbose  bool   // Whether to log debug information
	TmpDir     string // Directory in-which to generate intermediate files & archives
	IsRmtemp   bool   // Delete tmp dir after execution?
	DestDir    string // Where to generate .deb files and source debs (.dsc files etc)
	WorkingDir string // This is the root from which to find .go files, templates, resources, etc

	TemplateDir string            // Optional. Only required if you're using templates
	Resources   map[string]string // Optional. Only if debgo packages your resources automatically. Key is the destination file. Value is the local file

	ExtraData map[string]interface{} // Optional for templates
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

// initialize "template data" object
func (pkg *Package) NewTemplateData() TemplateData {
	templateVars := TemplateData{pkg.Name, pkg.Version, pkg.Maintainer, pkg.MaintainerEmail, pkg.Architecture,
		pkg.Section,
		pkg.Depends,
		pkg.BuildDepends,
		pkg.Priority,
		pkg.Description,
		pkg.StandardsVersion,
		"",
		pkg.Status,
		"",
		pkg.Format,
		pkg.AdditionalControlData,
		pkg.ExtraData,
		nil}
	return templateVars
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

// A factory for a Go Package. Includes dependencies and Go Path information
func NewGoPackage(name, version, maintainer string) *Package {
	pkg := NewPackage(name, version, maintainer)
	pkg.ExtraData = map[string]interface{}{
		"GoPathExtra": []string{GO_PATH_EXTRA_DEFAULT}}
	pkg.BuildDepends = BUILD_DEPENDS_DEFAULT
	pkg.Depends = DEPENDS_DEFAULT
	return pkg
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
	pkg.Priority = PRIORITY_DEFAULT
	pkg.StandardsVersion = STANDARDS_VERSION_DEFAULT
	pkg.Section = SECTION_DEFAULT
	pkg.Format = FORMAT_DEFAULT
	pkg.Status = STATUS_DEFAULT
	return pkg
}
