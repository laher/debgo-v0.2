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

package debgen

import (
	"fmt"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/targz"
	"io/ioutil"
	"log"
	"path/filepath"
)

// Default function for building the source archive
func GenSourceArtifacts(spkg *deb.SourcePackage, build *deb.BuildParams) error {
	//1. Build orig archive.
	err := BuildSourceOrigArchiveDefault(spkg, build)
	if err != nil {
		return err
	}
	//2. Build debian archive.
	err = BuildSourceDebianArchiveDefault(spkg, build)
	if err != nil {
		return err
	}
	//3. Build dsc file.
	err = BuildDscFileDefault(spkg, build)
	if err != nil {
		return err
	}

	return err
}

// Builds <package>.orig.tar.gz
// This contains all the original data.
func BuildSourceOrigArchiveDefault(spkg *deb.SourcePackage, build *deb.BuildParams) error {
	//TODO add/exclude resources to /usr/share
	origFilePath := filepath.Join(build.DestDir, spkg.OrigFileName)
	tgzw, err := targz.NewWriterFromFile(origFilePath)
	if err != nil {
		return err
	}
	err = AddSources(spkg, build.WorkingDir, spkg.Name+"-"+spkg.Version, tgzw, build)
	if err != nil {
		return err
	}
	err = tgzw.Close()
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("Created %s", origFilePath)
	}
	return nil
}

// Builds <package>.debian.tar.gz
// This contains all the control data, changelog, rules, etc
//
func BuildSourceDebianArchiveDefault(spkg *deb.SourcePackage, build *deb.BuildParams) error {
	//set up template
	templateVars := NewTemplateData(spkg.Package)

	// generate .debian.tar.gz (just containing debian/ directory)
	tgzw, err := targz.NewWriterFromFile(filepath.Join(build.DestDir, spkg.DebianFileName))
	templateDir := filepath.Join(build.TemplateDir, "source", "debian")
	//debian/control
	controlData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "control.tpl"), TemplateSourcedebControl, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(controlData, "debian/control", int64(0644))
	if err != nil {
		return err
	}

	//debian/compat
	compatData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "compat.tpl"), deb.DebianCompatDefault, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(compatData, "debian/compat", int64(0644))
	if err != nil {
		return err
	}

	//debian/rules
	rulesData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "rules.tpl"), TemplateDebianRules, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(rulesData, "debian/rules", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/format
	sourceFormatData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "source", "format.tpl"), TemplateDebianSourceFormat, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceFormatData, "debian/source/format", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/options
	sourceOptionsData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "source", "options.tpl"), TemplateDebianSourceOptions, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceOptionsData, "debian/source/options", int64(0644))
	if err != nil {
		return err
	}

	//debian/copyright
	copyrightData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "copyright.tpl"), TemplateDebianCopyright, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(copyrightData, "debian/copyright", int64(0644))
	if err != nil {
		return err
	}

	//debian/changelog
	initialChangelogTemplate := TemplateChangelogHeader + "\n\n" + TemplateChangelogInitialEntry + "\n\n" + TemplateChangelogFooter
	changelogData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
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
	readmeData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, "readme.tpl"), TemplateDebianReadme, templateVars)
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

	if build.IsVerbose {
		log.Printf("Created %s", tgzw.Filename)
	}
	return nil
}

func BuildDscFileDefault(spkg *deb.SourcePackage, build *deb.BuildParams) error {
	//set up template
	templateVars := NewTemplateData(spkg.Package)
	//4. Create dsc file (calculate checksums first)
	cs := new(deb.Checksums)
	err := cs.Add(filepath.Join(build.DestDir, spkg.OrigFileName), spkg.OrigFileName)
	if err != nil {
		return err
	}
	err = cs.Add(filepath.Join(build.DestDir, spkg.DebianFileName), spkg.DebianFileName)
	if err != nil {
		return err
	}
	templateVars.Checksums = cs
	dscData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "source", "dsc.tpl"), TemplateDebianDsc, templateVars)
	if err != nil {
		return err
	}
	dscFilePath := filepath.Join(build.DestDir, spkg.DscFileName)
	err = ioutil.WriteFile(dscFilePath, dscData, 0644)
	if err == nil {
		if build.IsVerbose {
			log.Printf("Wrote %s", dscFilePath)
		}
	}
	return err
}

// TODO: unfinished: need to discover root dir to determine which dirs to pre-make.
func AddSources(spkg *deb.SourcePackage, codeDir, destinationPrefix string, tgzw *targz.Writer, build *deb.BuildParams) error {
	goPathRootTemp := GetGoPathElement(codeDir)
	goPathRoot, err := filepath.EvalSymlinks(goPathRootTemp)
	if err != nil {
		log.Printf("Could not evaluate symlinks for '%s'", goPathRootTemp)
		goPathRoot = goPathRootTemp
	}
	if build.IsVerbose {
		log.Printf("Code dir '%s' (using goPath element '%s')", codeDir, goPathRoot)
	}
	sources, err := GlobForSources(goPathRootTemp, codeDir, GlobGoSources, destinationPrefix, []string{build.TmpDir, build.DestDir})
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
