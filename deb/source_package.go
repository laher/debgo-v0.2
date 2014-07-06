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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Function type to allow customized build process
type BuildSourcePackageFunc func(spkg *SourcePackage) error

// The source package is a cross-platform package with a .dsc file.
type SourcePackage struct {
	*Package
	DscFilePath    string
	OrigFilePath   string
	DebianFilePath string
	BuildFunc      BuildSourcePackageFunc
}

// Factory for a source package. Sets up default paths..
func NewSourcePackage(pkg *Package) *SourcePackage {
	spkg := &SourcePackage{Package: pkg,
		BuildFunc: BuildSourcePackageDefault}
	spkg.InitDefaultFilenames()
	return spkg
}

// Initialises default filenames, using .tar.gz as the archive type
func (spkg *SourcePackage) InitDefaultFilenames() {
	spkg.DscFilePath = filepath.Join(spkg.DestDir, spkg.Name+"_"+spkg.Version+".dsc")
	spkg.OrigFilePath = filepath.Join(spkg.DestDir, spkg.Name+"_"+spkg.Version+".orig.tar.gz")
	spkg.DebianFilePath = filepath.Join(spkg.DestDir, spkg.Name+"_"+spkg.Version+".debian.tar.gz")
}

// TODO: unfinished: need to discover root dir to determine which dirs to pre-make.
func (spkg *SourcePackage) AddSources(codeDir, destinationPrefix string, tgzw *TarGzWriter) error {
	goPathRoot := getGoPathElement(codeDir)
	goPathRootResolved, err := filepath.EvalSymlinks(goPathRoot)
	if err != nil {
		log.Printf("Could not evaluate symlinks for '%s'", goPathRoot)
		goPathRootResolved = goPathRoot
	}
	if spkg.IsVerbose {
		log.Printf("Code dir '%s' (using goPath element '%s')", codeDir, goPathRootResolved)
	}
	return spkg.addSources(goPathRootResolved, codeDir, destinationPrefix, tgzw)
}

// Get sources and append them
func (spkg *SourcePackage) addSources(goPathRoot, codeDir, destinationPrefix string, tgzw *TarGzWriter) error {
	sources, err := globForSources(goPathRoot, codeDir, destinationPrefix, []string{spkg.TmpDir, spkg.DestDir})
	if err != nil {
		return err
	}
	for destName, match := range sources {
		err = tgzw.AddFile(match, destName)
		if err != nil {
			return fmt.Errorf("Error adding go sources (match %s): %v,", match, err)
		}

	}
	return nil
}

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

// Builds 'source package' using default templating technique.
func (spkg *SourcePackage) Build() error {
	//build
	//1. prepare destination
	err := os.MkdirAll(spkg.DestDir, 0777)
	if err != nil {
		return err
	}
	return spkg.BuildFunc(spkg)
}

// Default function for building the source archive
func BuildSourcePackageDefault(spkg *SourcePackage) error {
	//2. Build orig archive.
	err := BuildSourceOrigArchiveDefault(spkg)
	if err != nil {
		return err
	}
	//3. Build debian archive.
	err = BuildSourceDebianArchiveDefault(spkg)
	if err != nil {
		return err
	}
	//4. Build dsc file.
	err = BuildDscFileDefault(spkg)
	if err != nil {
		return err
	}

	return err
}

// Builds <package>.orig.tar.gz
// This contains all the original data.
func BuildSourceOrigArchiveDefault(spkg *SourcePackage) error {
	//TODO add/exclude resources to /usr/share
	tgzw, err := NewTarGzWriter(spkg.OrigFilePath)
	if err != nil {
		return err
	}
	err = spkg.AddSources(spkg.WorkingDir, spkg.Name+"-"+spkg.Version, tgzw)
	if err != nil {
		return err
	}
	err = tgzw.Close()
	if err != nil {
		return err
	}
	if spkg.IsVerbose {
		log.Printf("Created %s", spkg.OrigFilePath)
	}
	return nil
}

// Builds <package>.debian.tar.gz
// This contains all the control data, changelog, rules, etc
//
func BuildSourceDebianArchiveDefault(spkg *SourcePackage) error {
	//set up template
	templateVars := spkg.NewTemplateData()

	// generate .debian.tar.gz (just containing debian/ directory)
	tgzw, err := NewTarGzWriter(spkg.DebianFilePath)

	//debian/control
	controlData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "control.tpl"), TEMPLATE_SOURCEDEB_CONTROL, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(controlData, "debian/control", int64(0644))
	if err != nil {
		return err
	}

	//debian/compat
	compatData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "compat.tpl"), TEMPLATE_DEBIAN_COMPAT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(compatData, "debian/compat", int64(0644))
	if err != nil {
		return err
	}

	//debian/rules
	rulesData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "rules.tpl"), TEMPLATE_DEBIAN_RULES, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(rulesData, "debian/rules", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/format
	sourceFormatData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "source_format.tpl"), TEMPLATE_DEBIAN_SOURCE_FORMAT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceFormatData, "debian/source/format", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/options
	sourceOptionsData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "source_options.tpl"), TEMPLATE_DEBIAN_SOURCE_OPTIONS, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceOptionsData, "debian/source/options", int64(0644))
	if err != nil {
		return err
	}

	//debian/copyright
	copyrightData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "copyright.tpl"), TEMPLATE_DEBIAN_COPYRIGHT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(copyrightData, "debian/copyright", int64(0644))
	if err != nil {
		return err
	}

	//debian/changelog
	initialChangelogTemplate := TEMPLATE_CHANGELOG_HEADER + "\n\n" + TEMPLATE_CHANGELOG_INITIAL_ENTRY + "\n\n" + TEMPLATE_CHANGELOG_FOOTER
	changelogData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(changelogData, "debian/changelog", int64(0644))
	if err != nil {
		return err
	}

	//generate debian/README.Debian
	//TODO: try pulling in README.md etc
	//debian/README.Debian
	readmeData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "readme.tpl"), TEMPLATE_DEBIAN_README, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(readmeData, "debian/README.debian", int64(0644))
	if err != nil {
		return err
	}

	err = tgzw.Close()
	if err != nil {
		return err
	}

	if spkg.IsVerbose {
		log.Printf("Created %s", spkg.DebianFilePath)
	}
	return nil
}

func BuildDscFileDefault(spkg *SourcePackage) error {
	//set up template
	templateVars := spkg.NewTemplateData()
	//4. Create dsc file (calculate checksums first)
	cs := new(Checksums)
	err := cs.Add(spkg.OrigFilePath, filepath.Base(spkg.OrigFilePath))
	if err != nil {
		return err
	}
	err = cs.Add(spkg.DebianFilePath, filepath.Base(spkg.DebianFilePath))
	if err != nil {
		return err
	}
	templateVars.Checksums = cs
	dscData, err := ProcessTemplateFileOrString(filepath.Join(spkg.TemplateDir, "dsc.tpl"), TEMPLATE_DEBIAN_DSC, templateVars)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(spkg.DscFilePath, dscData, 0644)
	if err == nil {
		if spkg.IsVerbose {
			log.Printf("Wrote %s", spkg.DscFilePath)
		}
	}
	return err
}
