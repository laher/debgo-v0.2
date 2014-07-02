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

type SourcePackage struct {
	*Package
	DscFilePath    string
	OrigFilePath   string
	DebianFilePath string
}

func NewSourcePackage(pkg *Package) *SourcePackage {
	dscPath := filepath.Join(pkg.DestDir, pkg.Name+"_"+pkg.Version+".dsc")
	origFilePath := filepath.Join(pkg.DestDir, pkg.Name+"_"+pkg.Version+".orig.tar.gz")
	debianFilePath := filepath.Join(pkg.DestDir, pkg.Name+"_"+pkg.Version+".debian.tar.gz")
	return &SourcePackage{Package: pkg,
		DscFilePath:    dscPath,
		OrigFilePath:   origFilePath,
		DebianFilePath: debianFilePath}
}

// TODO: unfinished: need to discover root dir to determine which dirs to pre-make.
func (pkg *SourcePackage) AddSources(codeDir, destinationPrefix string, tgzw *TarGzWriter) error {
	goPathRoot := getGoPathElement(codeDir)
	goPathRootResolved, err := filepath.EvalSymlinks(goPathRoot)
	if err != nil {
		log.Printf("Could not evaluate symlinks for '%s'", goPathRoot)
		goPathRootResolved = goPathRoot
	}
	log.Printf("Code dir '%s' (using goPath element '%s')", codeDir, goPathRootResolved)
	return pkg.addSources(goPathRootResolved, codeDir, destinationPrefix, tgzw)
}



// Get sources and append them
func (pkg *SourcePackage) addSources(goPathRoot, codeDir, destinationPrefix string, tgzw *TarGzWriter) error {
	sources, err := globForSources(goPathRoot, codeDir, destinationPrefix, []string{pkg.TmpDir, pkg.DestDir})
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
	/*
	//1. Glob for files in this dir
	//log.Printf("Globbing %s", codeDir)
	matches, err := filepath.Glob(filepath.Join(codeDir, "*.go"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		absMatch, err := filepath.Abs(match)
		if err != nil {
			return fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
		}
		relativeMatch, err := filepath.Rel(goPathRoot, absMatch)
		if err != nil {
			return fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
		}
		destName := filepath.Join(destinationPrefix, relativeMatch)
		err = tgzw.AddFile(match, destName)
		if err != nil {
			return fmt.Errorf("Error adding go sources (match %s): %v,", match, err)
		}
	}

	//2. Recurse into subdirs
	fis, err := ioutil.ReadDir(codeDir)
	for _, fi := range fis {
		if fi.IsDir() && fi.Name() != pkg.TmpDir {
			err := pkg.addSources(goPathRoot, filepath.Join(codeDir, fi.Name()), destinationPrefix, tgzw)
			//sources = append(sources, additionalItems...)
			if err != nil {
				return err
			}
		}
	}
	return err
	*/
}

func (pkg *SourcePackage) CopySourceRecurse(codeDir, destDir string) (err error) {
	log.Printf("Globbing %s", codeDir)
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
		//TODO copy files
		log.Printf("copying %s into %s", match, filepath.Join(destDir, filepath.Base(match)))
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
		log.Printf("Comparing fi.Name %s with tmpdir %v", fi.Name(), pkg.Package)
		if fi.IsDir() && fi.Name() != pkg.TmpDir {
			err = pkg.CopySourceRecurse(filepath.Join(codeDir, fi.Name()), filepath.Join(destDir, fi.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
// prepare folders and debian/ files.
// (everything except copying source)
func SdebPrepare(workingDirectory, appName, maintainer, version, arches, description, buildDepends string, metadataDeb map[string]interface{}) (err error) {
	//make temp dir & subfolders
	tmpDir := filepath.Join(workingDirectory, DIRNAME_TEMP)
	debianDir := filepath.Join(tmpDir, "debian")
	err = os.MkdirAll(filepath.Join(debianDir, "source"), 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tmpDir, "src"), 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tmpDir, "bin"), 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tmpDir, "pkg"), 0777)
	if err != nil {
		return err
	}
	//write control file and related files
	tpl, err := template.New("rules").Parse(TEMPLATE_DEBIAN_RULES)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(debianDir, "rules"))
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	err = tpl.Execute(file, appName)
	if err != nil {
		return err
	}
	sdebControlFile := getSdebControlFileContent(appName, maintainer, version, arches, description, buildDepends, metadataDeb)
	ioutil.WriteFile(filepath.Join(debianDir, "control"), sdebControlFile, 0666)
	//copy source into folders
	//call dpkg-build, if available
	return err
}
*/

func (pkg *SourcePackage) BuildWithDefaults() error {
	//build

	//1. prepare destination
	err := os.MkdirAll(pkg.DestDir, 0777)
	if err != nil {
		return err
	}

	err = pkg.BuildOrigArchive()
	if err != nil {
		return err
	}
	err = pkg.BuildDebianArchive()
	if err != nil {
		return err
	}
	err = pkg.BuildDscFile()
	if err != nil {
		return err
	}

	return err
}

func (pkg *SourcePackage) BuildOrigArchive() error {
	//2. generate orig.tar.gz

	//TODO add/exclude resources to /usr/share
	tgzw, err := NewTarGzWriter(pkg.OrigFilePath)
	if err != nil {
		return err
	}
	err = pkg.AddSources(pkg.WorkingDir, pkg.Name+"-"+pkg.Version, tgzw)
	if err != nil {
		return err
	}
	err = tgzw.Close()
	if err != nil {
		return err
	}
	log.Printf("Created %s", pkg.OrigFilePath)
	return nil
}

func (pkg *SourcePackage) BuildDebianArchive() error {
	//set up template
	templateVars := pkg.NewTemplateData()

	//3. generate .debian.tar.gz (just containing debian/ directory)
	tgzw, err := NewTarGzWriter(pkg.DebianFilePath)

	//debian/control
	controlData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "control.tpl"), TEMPLATE_SOURCEDEB_CONTROL, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(controlData, "debian/control", int64(0644))
	if err != nil {
		return err
	}

	//debian/compat
	compatData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "compat.tpl"), TEMPLATE_DEBIAN_COMPAT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(compatData, "debian/compat", int64(0644))
	if err != nil {
		return err
	}

	//debian/rules
	rulesData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "rules.tpl"), TEMPLATE_DEBIAN_RULES, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(rulesData, "debian/rules", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/format
	sourceFormatData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "source_format.tpl"), TEMPLATE_DEBIAN_SOURCE_FORMAT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceFormatData, "debian/source/format", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/options
	sourceOptionsData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "source_options.tpl"), TEMPLATE_DEBIAN_SOURCE_OPTIONS, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceOptionsData, "debian/source/options", int64(0644))
	if err != nil {
		return err
	}

	//debian/copyright
	copyrightData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "copyright.tpl"), TEMPLATE_DEBIAN_COPYRIGHT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(copyrightData, "debian/copyright", int64(0644))
	if err != nil {
		return err
	}

	//debian/changelog
	/*(slightly different)
	var changelogData []byte
	if pkg.Changelog != nil {
		rdr, err := pkg.Changelog.GetReader()
		if err != nil {
			return err
		}
		changelogData, err = ioutil.ReadAll(rdr)
		if err != nil {
			return err
		}
	}

		_, err = os.Stat(changelogFilename)
		if os.IsNotExist(err) {
			initialChangelogTemplate := TEMPLATE_CHANGELOG_HEADER + "\n\n" + TEMPLATE_CHANGELOG_INITIAL_ENTRY + "\n\n" + TEMPLATE_CHANGELOG_FOOTER
			changelogData, err = ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			changelogData, err = ioutil.ReadFile(changelogFilename)
			if err != nil {
				return err
			}
		}
	*/
	//if len(changelogData) == 0 {
	initialChangelogTemplate := TEMPLATE_CHANGELOG_HEADER + "\n\n" + TEMPLATE_CHANGELOG_INITIAL_ENTRY + "\n\n" + TEMPLATE_CHANGELOG_FOOTER
	changelogData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(changelogData, "debian/changelog", int64(0644))
	if err != nil {
		return err
	}
	//}

	//generate debian/README.Debian
	//TODO: try pulling in README.md etc
	//debian/README.Debian
	readmeData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "readme.tpl"), TEMPLATE_DEBIAN_README, templateVars)
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
	log.Printf("Created %s", pkg.DebianFilePath)
	return nil
}

func (pkg *SourcePackage) BuildDscFile() error {
	//set up template
	templateVars := pkg.NewTemplateData()
	//4. Create dsc file (calculate checksums first)
	cs := new(Checksums)
	err := cs.Add(pkg.OrigFilePath, filepath.Base(pkg.OrigFilePath))
	if err != nil {
		return err
	}
	err = cs.Add(pkg.DebianFilePath, filepath.Base(pkg.DebianFilePath))
	if err != nil {
		return err
	}
	templateVars.Checksums = cs
	dscData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "dsc.tpl"), TEMPLATE_DEBIAN_DSC, templateVars)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(pkg.DscFilePath, dscData, 0644)
	if err == nil {
		log.Printf("Wrote %s", pkg.DscFilePath)
	}
	return err
}
