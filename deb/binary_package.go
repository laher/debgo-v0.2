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

// BuildDevPackageFunc specifies a function which can build a DevPackage
type BuildBinaryArtifactFunc func(*BinaryArtifact, *BuildParams) error

// *BinaryPackage specifies functionality for building binary '.deb' packages.
// This encapsulates a Package plus information about platform-specific debs and executables
type BinaryPackage struct {
	*Package
	ExeDest string
}

// NewBinaryPackage is a factory for BinaryPackage
func NewBinaryPackage(pkg *Package) *BinaryPackage {
	return &BinaryPackage{Package: pkg, ExeDest: "/usr/bin"}
}

// GetArtifacts gets and returns an artifact for each architecture.
// Returns an error if the package's architecture is un-parseable
func (pkg *BinaryPackage) GetArtifacts() (map[Architecture]*BinaryArtifact, error) {
	arches, err := pkg.GetArches()
	if err != nil {
		return nil, err
	}
	ret := map[Architecture]*BinaryArtifact{}
	for _, arch := range arches {
		archArtifact := NewBinaryArtifact(pkg, arch)
		ret[arch] = archArtifact
	}
	return ret, nil
}
