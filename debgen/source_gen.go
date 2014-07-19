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
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/targz"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)


//SourcePackageGenerator generates source packages using templates and some overrideable behaviours
type SourcePackageGenerator struct {
	SourcePackage *deb.SourcePackage
	BuildParams *BuildParams
	DefaultTemplateStrings map[string]string
	//DebianFiles map[string]string
	OrigFiles map[string]string
}

//NewSourcePackageGenerator is a factory for SourcePackageGenerator.
func NewSourcePackageGenerator(sourcePackage *deb.SourcePackage, buildParams *BuildParams) *SourcePackageGenerator {
	spgen := &SourcePackageGenerator{SourcePackage:sourcePackage, BuildParams:buildParams}
	spgen.DefaultTemplateStrings = defaultTemplateStrings()
	return spgen
}

// ApplyDefaultsPureGo overrides some template variables for pure-Go packages
func (spgen *SourcePackageGenerator) ApplyDefaultsPureGo() {
	spgen.DefaultTemplateStrings["debian/rules"] = TemplateDebianRulesForGo
}

// Get the default templates for source packages
func defaultTemplateStrings() map[string]string {
	//defensive copy
	templateStringsSource := map[string]string{}
	for k, v := range TemplateStringsSourceDefault {
		templateStringsSource[k] = v
	}
	return templateStringsSource
}

// GenerateAll builds all the artifacts using the default behaviour.
// Note that we could implement alternative methods in future (e.g. using a GenDiffArchive)
func (spgen *SourcePackageGenerator) GenerateAllDefault() error {
	//1. Build orig archive.
	err := spgen.GenOrigArchive()
	if err != nil {
		return err
	}
	//2. Build debian archive.
	err = spgen.GenDebianArchive()
	if err != nil {
		return err
	}
	//3. Build dsc file, including checksums
	err = spgen.GenDscFile()
	if err != nil {
		return err
	}
	return err
}

// GenOrigArchive builds <package>.orig.tar.gz
// This contains the original upstream source code and data.
func (spgen *SourcePackageGenerator) GenOrigArchive() error {
	//TODO add/exclude resources to /usr/share
	origFilePath := filepath.Join(spgen.BuildParams.DestDir, spgen.SourcePackage.OrigFileName)
	tgzw, err := targz.NewWriterFromFile(origFilePath)
	defer tgzw.Close()
	if err != nil {
		return err
	}
	err = TarAddFiles(tgzw.Writer, spgen.OrigFiles)
	if err != nil {
		return err
	}
	err = tgzw.Close()
	if err != nil {
		return err
	}
	if spgen.BuildParams.IsVerbose {
		log.Printf("Created %s", origFilePath)
	}
	return nil
}

// GenDebianArchive builds <package>.debian.tar.gz
// This contains all the control data, changelog, rules, etc
func (spgen *SourcePackageGenerator) GenDebianArchive() error {
	//set up template
	templateVars := NewTemplateData(spgen.SourcePackage.Package)

	// generate .debian.tar.gz (just containing debian/ directory)
	tgzw, err := targz.NewWriterFromFile(filepath.Join(spgen.BuildParams.DestDir, spgen.SourcePackage.DebianFileName))
	defer tgzw.Close()
	resourceDir := filepath.Join(spgen.BuildParams.TemplateDir, "source", DebianDir)
	templateDir := filepath.Join(spgen.BuildParams.TemplateDir, "source", DebianDir)

	for debianFile, defaultTemplateStr := range spgen.DefaultTemplateStrings {
		debianFilePath := strings.Replace(debianFile, "/", string(os.PathSeparator), -1) //fixing source/options, source/format for local files
		resourcePath := filepath.Join(resourceDir, debianFilePath)
		_, err = os.Stat(resourcePath)
		if err == nil {
			err = TarAddFile(tgzw.Writer, resourcePath, debianFile)
			if err != nil {
				return err
			}
		} else {
			controlData, err := TemplateFileOrString(filepath.Join(templateDir, debianFilePath+TplExtension), defaultTemplateStr, templateVars)
			if err != nil {
				return err
			}
			err = TarAddBytes(tgzw.Writer, controlData, DebianDir+"/"+debianFile, int64(0644))
			if err != nil {
				return err
			}
		}
	}

	// postrm/postinst etc from main store
	for _, scriptName := range deb.MaintainerScripts {
		resourcePath := filepath.Join(spgen.BuildParams.ResourcesDir, DebianDir, scriptName)
		_, err = os.Stat(resourcePath)
		if err == nil {
			err = TarAddFile(tgzw.Writer, resourcePath, scriptName)
			if err != nil {
				return err
			}
		} else {
			templatePath := filepath.Join(spgen.BuildParams.TemplateDir, DebianDir, scriptName+TplExtension)
			_, err = os.Stat(templatePath)
			//TODO handle non-EOF errors
			if err == nil {
				scriptData, err := TemplateFile(templatePath, templateVars)
				if err != nil {
					return err
				}
				err = TarAddBytes(tgzw.Writer, scriptData, scriptName, 0755)
				if err != nil {
					return err
				}
			}
		}
	}

	err = tgzw.Close()
	if err != nil {
		return err
	}

	if spgen.BuildParams.IsVerbose {
		log.Printf("Created %s", filepath.Join(spgen.BuildParams.DestDir, spgen.SourcePackage.DebianFileName))
	}
	return nil
}

func (spgen *SourcePackageGenerator) GenDscFile() error {
	//set up template
	templateVars := NewTemplateData(spgen.SourcePackage.Package)
	//4. Create dsc file (calculate checksums first)
	cs := new(deb.Checksums)
	err := cs.Add(filepath.Join(spgen.BuildParams.DestDir, spgen.SourcePackage.OrigFileName), spgen.SourcePackage.OrigFileName)
	if err != nil {
		return err
	}
	err = cs.Add(filepath.Join(spgen.BuildParams.DestDir, spgen.SourcePackage.DebianFileName), spgen.SourcePackage.DebianFileName)
	if err != nil {
		return err
	}
	templateVars.Checksums = cs
	dscData, err := TemplateFileOrString(filepath.Join(spgen.BuildParams.TemplateDir, "source", "dsc.tpl"), TemplateDebianDsc, templateVars)
	if err != nil {
		return err
	}
	dscFilePath := filepath.Join(spgen.BuildParams.DestDir, spgen.SourcePackage.DscFileName)
	err = ioutil.WriteFile(dscFilePath, dscData, 0644)
	if err == nil {
		if spgen.BuildParams.IsVerbose {
			log.Printf("Wrote %s", dscFilePath)
		}
	}
	return err
}
