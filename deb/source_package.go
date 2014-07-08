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
	"os"
	"path/filepath"
)

// Function type to allow customized build process
type BuildSourcePackageFunc func(spkg *SourcePackage, build *BuildParams) error

// The source package is a cross-platform package with a .dsc file.
type SourcePackage struct {
	*Package
	DscFilePath    string
	OrigFilePath   string
	DebianFilePath string
	BuildFunc      BuildSourcePackageFunc
}

// Factory for a source package. Sets up default paths..
func NewSourcePackage(pkg *Package, buildFunc BuildSourcePackageFunc) *SourcePackage {
	spkg := &SourcePackage{Package: pkg,
		BuildFunc: buildFunc}
//	spkg.InitDefaultFilenames()
	return spkg
}

// Initialises default filenames, using .tar.gz as the archive type
func (spkg *SourcePackage) InitDefaultFilenames(build *BuildParams) {
	spkg.DscFilePath = filepath.Join(build.DestDir, spkg.Name+"_"+spkg.Version+".dsc")
	spkg.OrigFilePath = filepath.Join(build.DestDir, spkg.Name+"_"+spkg.Version+".orig.tar.gz")
	spkg.DebianFilePath = filepath.Join(build.DestDir, spkg.Name+"_"+spkg.Version+".debian.tar.gz")
}
/*
func (spkg *SourcePackage) CopySourceRecurse(codeDir, destDir string) (err error) {
	if spkg.IsVerbose {
		log.Printf("Globbing %s", codeDir)
	}
	//get all files and copy into destDir
	matches, err := filepath.Glob(filepath.Join(codeDir, "*.go"))
	if err != nil {
		return err
	}
	if len(matches) > 0 {
		err = os.MkdirAll(destDir, 0777)
		if err != nil {
			return err
		}
	}
	for _, match := range matches {
		if spkg.IsVerbose {
			log.Printf("copying %s into %s", match, filepath.Join(destDir, filepath.Base(match)))
		}
		r, err := os.Open(match)
		if err != nil {
			return err
		}
		defer func() {
			err := r.Close()
			if err != nil {
				panic(err)
			}
		}()
		w, err := os.Create(filepath.Join(destDir, filepath.Base(match)))
		if err != nil {
			return err
		}
		defer func() {
			err := w.Close()
			if err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}
	}
	fis, err := ioutil.ReadDir(codeDir)
	for _, fi := range fis {
		if spkg.IsVerbose {
			log.Printf("Comparing fi.Name %s with tmpdir %v", fi.Name(), spkg.Package)
		}
		if fi.IsDir() && fi.Name() != spkg.TmpDir {
			err = spkg.CopySourceRecurse(filepath.Join(codeDir, fi.Name()), filepath.Join(destDir, fi.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
*/
// Builds 'source package' using default templating technique.
func (spkg *SourcePackage) Build(build *BuildParams) error {
	//build
	//1. prepare destination
	err := os.MkdirAll(build.DestDir, 0777)
	if err != nil {
		return err
	}
	return spkg.BuildFunc(spkg, build)
}

