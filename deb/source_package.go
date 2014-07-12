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

// SourcePackage is a cross-platform package with a .dsc file.
type SourcePackage struct {
	Package        *Package
	DscFileName    string
	OrigFileName   string
	DebianFileName string
	MappedFiles    map[string]string
}

// NewSourcePackage is a factory for SourcePackage. Sets up default paths..
// Initialises default filenames, using .tar.gz as the archive type
func NewSourcePackage(pkg *Package) *SourcePackage {
	spkg := &SourcePackage{Package: pkg}
	spkg.DscFileName = pkg.Name + "_" + pkg.Version + ".dsc"
	spkg.OrigFileName = pkg.Name + "_" + pkg.Version + ".orig.tar.gz"
	spkg.DebianFileName = pkg.Name + "_" + pkg.Version + ".debian.tar.gz"
	return spkg
}
