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

var ()

// This is the default build process for a BuildArtifact
func GenDeb(bdeb *deb.Deb, build *BuildParams) error {
	if build.IsVerbose {
		log.Printf("trying to write control file for %s", bdeb.Architecture)
	}

	_, err := GenControlArchive(bdeb, build)
	if err != nil {
		return err
	}
	_, err = GenDataArchive(bdeb, build)
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("trying to write .deb file for %s", bdeb.Architecture)
	}
	err = bdeb.Build(build.TmpDir, build.DestDir)
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("Closed deb")
	}
	return err
}

func GenControlArchive(bdeb *deb.Deb, build *BuildParams) (string, error) {
	archiveFilename := filepath.Join(build.TmpDir, bdeb.DebianArchive)
	controlTgzw, err := targz.NewWriterFromFile(archiveFilename)
	if err != nil {
		return "", err
	}
	templateVars := &TemplateData{Package: bdeb.Package, Deb: bdeb}
	//templateVars.Deb = bdeb

	err = GenControlFile(controlTgzw, templateVars, build)
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Wrote control file to control archive")
	}
	// This is where you include Postrm/Postinst etc
	for _, scriptName := range deb.MaintainerScripts {
		resourcePath := filepath.Join(build.ResourcesDir, DebianDir, scriptName)
		_, err = os.Stat(resourcePath)
		if err == nil {
			err = TarAddFile(controlTgzw.Tw, resourcePath, scriptName)
			if err != nil {
				return "", err
			}
		} else {
			templatePath := filepath.Join(build.TemplateDir, DebianDir, scriptName+TplExtension)
			_, err = os.Stat(templatePath)
			//TODO handle non-EOF errors
			if err == nil {
				scriptData, err := TemplateFile(templatePath, templateVars)
				if err != nil {
					return "", err
				}
				err = TarAddBytes(controlTgzw.Tw, scriptData, scriptName, 0755)
				if err != nil {
					return "", err
				}
			}
		}
	}

	err = controlTgzw.Close()
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Closed control archive")
	}
	return controlTgzw.Filename, err
}

func GenDataArchive(bdeb *deb.Deb, build *BuildParams) (string, error) {
	archiveFilename := filepath.Join(build.TmpDir, bdeb.DataArchive)
	dataTgzw, err := targz.NewWriterFromFile(archiveFilename)
	if err != nil {
		return "", err
	}
	err = TarAddFiles(dataTgzw.Tw, bdeb.MappedFiles)
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Added executables")
	}
	// TODO add README.debian automatically
	err = TarAddFiles(dataTgzw.Tw, bdeb.Package.MappedFiles)
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Added resources")
	}
	err = dataTgzw.Close()
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Closed data archive")
	}
	return dataTgzw.Filename, err
}

func GenControlFile(tgzw *targz.Writer, templateVars *TemplateData, build *BuildParams) error {
	resourcePath := filepath.Join(build.ResourcesDir, "DEBIAN", "control")
	_, err := os.Stat(resourcePath)
	if err == nil {
		err = TarAddFile(tgzw.Tw, resourcePath, "control")
		return err
	}
	//try template or use a string
	controlData, err := TemplateFileOrString(filepath.Join(build.TemplateDir, "control.tpl"), TemplateBinarydebControl, templateVars)
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("Control file:\n%s", string(controlData))
	}
	err = TarAddBytes(tgzw.Tw, controlData, "control", 0644)
	return err
}
