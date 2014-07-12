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

// DevPackage builds a '-dev' package, which itself should just contain the sources to use this package as a BuildDepends dependency.
// For pure Go packages, this can be cross-platform (architecture == 'all'), but in some cases it might need to be architecture specific
type DevPackage struct {
	Package *Package
	Dev     *Package
}

// NewDevPackage is a factory for DevPackage
func NewDevPackage(pkg *Package) *DevPackage {
	ddpkg := &DevPackage{Package: pkg} //		DebFilePath: debPath,
	return ddpkg
}

// InitBinaryPackage initialises the binary package for -dev.deb packages
func (ddpkg *DevPackage) InitDev() {
	if ddpkg.Dev == nil {
		//TODO *complete* copy of properties, using reflection?
		devpkg := NewPackage(ddpkg.Package.Name+"-dev", ddpkg.Package.Version, ddpkg.Package.Maintainer, ddpkg.Package.Description)
		devpkg.AdditionalControlData = ddpkg.Package.AdditionalControlData
		devpkg.Architecture = "all"
		ddpkg.Dev = devpkg
	}
}
