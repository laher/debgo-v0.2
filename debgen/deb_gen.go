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
	"log"
	"os"
	"path/filepath"
)

//DebGenerator generates source packages using templates and some overrideable behaviours
type DebGenerator struct {
	DebWriter *deb.DebWriter
	BuildParams *BuildParams
	DefaultTemplateStrings map[string]string
	OrigFiles map[string]string
}

//NewDebGenerator is a factory for SourcePackageGenerator.
func NewDebGenerator(debWriter *deb.DebWriter, buildParams *BuildParams) *DebGenerator {
	dgen := &DebGenerator{DebWriter:debWriter, BuildParams:buildParams,
		DefaultTemplateStrings:map[string]string{}, OrigFiles:map[string]string{}}
	return dgen
}

// GenerateAllDefault applies the default build process.
// First it writes each file, then adds them to the .deb as a separate io operation.
func (dgen *DebGenerator) GenerateAllDefault() error {
	if dgen.BuildParams.IsVerbose {
		log.Printf("trying to write control file for %s", dgen.DebWriter.Architecture)
	}
	err := dgen.GenControlArchive()
	if err != nil {
		return err
	}
	err = dgen.GenDataArchive()
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("trying to write .deb file for %s", dgen.DebWriter.Architecture)
	}
	//TODO switch this around
	err = dgen.DebWriter.Build(dgen.BuildParams.TmpDir, dgen.BuildParams.DestDir)
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("Closed deb")
	}
	return err
}

//GenControlArchive generates control archive, using a system of templates or files.
//
//First it attempts to find the file inside BuildParams.Resources.
//If that doesn't exist, it attempts to find a template in templateDir
//Finally, it attempts to use a string-based template.
func (dgen *DebGenerator) GenControlArchive() error {
	archiveFilename := filepath.Join(dgen.BuildParams.TmpDir, dgen.DebWriter.ControlArchive)
	controlTgzw, err := targz.NewWriterFromFile(archiveFilename)
	if err != nil {
		return err
	}
	templateVars := &TemplateData{Package: dgen.DebWriter.Package, Deb: dgen.DebWriter}
	//templateVars.Deb = dgen.DebWriter

	err = dgen.GenControlFile(controlTgzw, templateVars)
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("Wrote control file to control archive")
	}
	// This is where you include Postrm/Postinst etc
	for _, scriptName := range deb.MaintainerScripts {
		resourcePath := filepath.Join(dgen.BuildParams.ResourcesDir, DebianDir, scriptName)
		_, err = os.Stat(resourcePath)
		if err == nil {
			err = TarAddFile(controlTgzw.Writer, resourcePath, scriptName)
			if err != nil {
				return err
			}
		} else {
			templatePath := filepath.Join(dgen.BuildParams.TemplateDir, DebianDir, scriptName+TplExtension)
			_, err = os.Stat(templatePath)
			//TODO handle non-EOF errors
			if err == nil {
				scriptData, err := TemplateFile(templatePath, templateVars)
				if err != nil {
					return err
				}
				err = TarAddBytes(controlTgzw.Writer, scriptData, scriptName, 0755)
				if err != nil {
					return err
				}
			}
		}
	}

	err = controlTgzw.Close()
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("Closed control archive")
	}
	return err
}

// GenDataArchive generates the 'code' archive from files on the file system.
func (dgen *DebGenerator) GenDataArchive() error {
	archiveFilename := filepath.Join(dgen.BuildParams.TmpDir, dgen.DebWriter.DataArchive)
	dataTgzw, err := targz.NewWriterFromFile(archiveFilename)
	if err != nil {
		return err
	}
	err = TarAddFiles(dataTgzw.Writer, dgen.OrigFiles)
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("Added executables")
	}
/*
	// TODO add README.debian automatically
	err = TarAddFiles(dataTgzw.Writer, dgen.DebWriter.Package.MappedFiles)
	if err != nil {
		return err
	}
*/
	if dgen.BuildParams.IsVerbose {
		log.Printf("Added resources")
	}
	err = dataTgzw.Close()
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("Closed data archive")
	}
	return err
}

//Generates the control file based on a template
func (dgen *DebGenerator) GenControlFile(tgzw *targz.Writer, templateVars *TemplateData) error {
	resourcePath := filepath.Join(dgen.BuildParams.ResourcesDir, "debian", "control")
	_, err := os.Stat(resourcePath)
	if err == nil {
		err = TarAddFile(tgzw.Writer, resourcePath, "control")
		return err
	}
	//try template or use a string
	controlData, err := TemplateFileOrString(filepath.Join(dgen.BuildParams.TemplateDir, "control.tpl"), TemplateBinarydebControl, templateVars)
	if err != nil {
		return err
	}
	if dgen.BuildParams.IsVerbose {
		log.Printf("Control file:\n%s", string(controlData))
	}
	err = TarAddBytes(tgzw.Writer, controlData, "control", 0644)
	return err
}
