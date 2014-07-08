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

package debgo

import (
	"fmt"
	"github.com/laher/debgo-v0.2/deb"
	"io/ioutil"
	"log"
	"path/filepath"
)


// Default function for building the source archive
func BuildSourcePackageDefault(spkg *deb.SourcePackage, build *deb.BuildParams) error {
	//2. Build orig archive.
	err := BuildSourceOrigArchiveDefault(spkg, build)
	if err != nil {
		return err
	}
	//3. Build debian archive.
	err = BuildSourceDebianArchiveDefault(spkg, build)
	if err != nil {
		return err
	}
	//4. Build dsc file.
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
	tgzw, err := deb.NewTarGzWriter(spkg.OrigFilePath)
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
		log.Printf("Created %s", spkg.OrigFilePath)
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
	tgzw, err := deb.NewTarGzWriter(spkg.DebianFilePath)

	//debian/control
	controlData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "control.tpl"), TEMPLATE_SOURCEDEB_CONTROL, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(controlData, "debian/control", int64(0644))
	if err != nil {
		return err
	}

	//debian/compat
	compatData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "compat.tpl"), deb.DEBIAN_COMPAT_DEFAULT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(compatData, "debian/compat", int64(0644))
	if err != nil {
		return err
	}

	//debian/rules
	rulesData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "rules.tpl"), TEMPLATE_DEBIAN_RULES, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(rulesData, "debian/rules", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/format
	sourceFormatData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "source_format.tpl"), TEMPLATE_DEBIAN_SOURCE_FORMAT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceFormatData, "debian/source/format", int64(0644))
	if err != nil {
		return err
	}

	//debian/source/options
	sourceOptionsData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "source_options.tpl"), TEMPLATE_DEBIAN_SOURCE_OPTIONS, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(sourceOptionsData, "debian/source/options", int64(0644))
	if err != nil {
		return err
	}

	//debian/copyright
	copyrightData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "copyright.tpl"), TEMPLATE_DEBIAN_COPYRIGHT, templateVars)
	if err != nil {
		return err
	}
	err = tgzw.AddBytes(copyrightData, "debian/copyright", int64(0644))
	if err != nil {
		return err
	}

	//debian/changelog
	initialChangelogTemplate := TEMPLATE_CHANGELOG_HEADER + "\n\n" + TEMPLATE_CHANGELOG_INITIAL_ENTRY + "\n\n" + TEMPLATE_CHANGELOG_FOOTER
	changelogData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
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
	readmeData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "readme.tpl"), TEMPLATE_DEBIAN_README, templateVars)
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
		log.Printf("Created %s", spkg.DebianFilePath)
	}
	return nil
}

func BuildDscFileDefault(spkg *deb.SourcePackage, build *deb.BuildParams) error {
	//set up template
	templateVars := NewTemplateData(spkg.Package)
	//4. Create dsc file (calculate checksums first)
	cs := new(deb.Checksums)
	err := cs.Add(spkg.OrigFilePath, filepath.Base(spkg.OrigFilePath))
	if err != nil {
		return err
	}
	err = cs.Add(spkg.DebianFilePath, filepath.Base(spkg.DebianFilePath))
	if err != nil {
		return err
	}
	templateVars.Checksums = cs
	dscData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "dsc.tpl"), TEMPLATE_DEBIAN_DSC, templateVars)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(spkg.DscFilePath, dscData, 0644)
	if err == nil {
		if build.IsVerbose {
			log.Printf("Wrote %s", spkg.DscFilePath)
		}
	}
	return err
}


// TODO: unfinished: need to discover root dir to determine which dirs to pre-make.
func AddSources(spkg *deb.SourcePackage, codeDir, destinationPrefix string, tgzw *deb.TarGzWriter, build *deb.BuildParams) error {
	goPathRootTemp := getGoPathElement(codeDir)
	goPathRoot, err := filepath.EvalSymlinks(goPathRootTemp)
	if err != nil {
		log.Printf("Could not evaluate symlinks for '%s'", goPathRootTemp)
		goPathRoot = goPathRootTemp
	}
	if build.IsVerbose {
		log.Printf("Code dir '%s' (using goPath element '%s')", codeDir, goPathRoot)
	}
	sources, err := globForSources(goPathRootTemp, codeDir, destinationPrefix, []string{build.TmpDir, build.DestDir})
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

