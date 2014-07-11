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

// Package is the base unit for this library.
// A *Package contains metadata.
type Package struct {
	Name        string // Package name
	Version     string // Package version
	Description string // Description
	Maintainer  string // Maintainer

	AdditionalControlData map[string]string // Other key/values to go into the Control file.

	Architecture string // Supported values: "all", "x386", "amd64", "armhf". TODO: armel

	Depends      string // Depends
	BuildDepends string // BuildDepends is only required for "sourcedebs".

	Priority         string
	StandardsVersion string
	Section          string
	Format           string
	Status           string
	Other            string
	Source           string

	ExtraData map[string]interface{} // Optional for templates

	MappedFiles map[string]string
}

// NewPackage is a factory for a Package. Name, Version, Maintainer and Description are mandatory.
func NewPackage(name, version, maintainer, description string) *Package {
	pkg := new(Package)
	pkg.Name = name
	pkg.Version = version
	pkg.Maintainer = maintainer
	pkg.Description = description
	pkg.Architecture = "any" //default ...
	pkg.Priority = PriorityDefault
	pkg.StandardsVersion = StandardsVersionDefault
	pkg.Section = SectionDefault
	pkg.Format = FormatDefault
	pkg.Status = StatusDefault
	pkg.MappedFiles = map[string]string{}
	return pkg
}

// GetArches resolves architecture(s) and return as a slice
func (pkg *Package) GetArches() ([]Architecture, error) {
	arches, err := resolveArches(pkg.Architecture)
	return arches, err
}

func (pkg *Package) SetField(key, value string) {
	switch key {
	case "Package":
		pkg.Name = value
	case "Source":
		pkg.Source = value
	case "Version":
		pkg.Version = value
	case "Description":
		pkg.Description = value
	case "Maintainer":
		pkg.Maintainer = value
	case "Architecture":
		pkg.Architecture = value
	case "Depends":
		pkg.Depends = value
	case "BuildDepends":
		pkg.BuildDepends = value
	case "Priority":
		pkg.Priority = value
	case "StandardsVersion":
		pkg.StandardsVersion = value
	case "Section":
		pkg.Section = value
	case "Format":
		pkg.Section = value
	case "Status":
		pkg.Section = value
	case "Other":
		pkg.Section = value
	default:
		pkg.AdditionalControlData[key] = value
	}
}
