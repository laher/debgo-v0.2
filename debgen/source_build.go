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
	defer tgzw.Close()
	if err != nil {
		return err
	}
	err = TarAddFiles(tgzw.Tw, spkg.MappedFiles)
	if err != nil {
		return err
	}
	err = TarAddFiles(tgzw.Tw, spkg.Package.MappedFiles)
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
	defer tgzw.Close()
	resourceDir := filepath.Join(build.TemplateDir, "source", DebianDir)
	templateDir := filepath.Join(build.TemplateDir, "source", DebianDir)

	for debianFile, defaultTemplateStr := range SourceDebianFiles {
		debianFilePath := strings.Replace(debianFile, "/", string(os.PathSeparator), -1) //fixing source/options, source/format for local files
		resourcePath := filepath.Join(resourceDir, debianFilePath)
		_, err = os.Stat(resourcePath)
		if err == nil {
			err = TarAddFile(tgzw.Tw, resourcePath, debianFile)
			if err != nil {
				return err
			}
		} else {
			controlData, err := ProcessTemplateFileOrString(filepath.Join(templateDir, debianFilePath+TplExtension), defaultTemplateStr, templateVars)
			if err != nil {
				return err
			}
			err = TarAddBytes(tgzw.Tw, controlData, DebianDir+"/"+debianFile, int64(0644))
			if err != nil {
				return err
			}
		}
	}

	// postrm/postinst etc from main store
	for _, scriptName := range deb.MaintainerScripts {
		resourcePath := filepath.Join(build.ResourcesDir, DebianDir, scriptName)
		_, err = os.Stat(resourcePath)
		if err == nil {
			err = TarAddFile(tgzw.Tw, resourcePath, scriptName)
			if err != nil {
				return err
			}
		} else {
			templatePath := filepath.Join(build.TemplateDir, DebianDir, scriptName+TplExtension)
			_, err = os.Stat(templatePath)
			//TODO handle non-EOF errors
			if err == nil {
				scriptData, err := ProcessTemplateFile(templatePath, templateVars)
				if err != nil {
					return err
				}
				err = TarAddBytes(tgzw.Tw, scriptData, scriptName, 0755)
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
